// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package localdocker

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	machinev1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"

	apiv1alpha1 "github.com/gardener/gardener/pkg/provider-local/machine-provider/api/v1alpha1"
	"github.com/gardener/gardener/pkg/provider-local/machine-provider/api/validation"
)

func (d *localDriver) CreateMachine(ctx context.Context, req *driver.CreateMachineRequest) (*driver.CreateMachineResponse, error) {
	if req.MachineClass.Provider != apiv1alpha1.Provider {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("requested for provider '%s', we only support '%s'", req.MachineClass.Provider, apiv1alpha1.Provider))
	}

	klog.V(3).Infof("Machine creation request has been received for %q", req.Machine.Name)
	defer klog.V(3).Infof("Machine creation request has been processed for %q", req.Machine.Name)

	providerSpec, err := validateProviderSpecAndSecret(req.MachineClass, req.Secret)
	if err != nil {
		return nil, err
	}

	klog.V(3).Infof("Ensure machine image %q is available", providerSpec.Image)

	if _, err := d.dockerClient.ImageInspect(ctx, providerSpec.Image); err != nil {
		_, err = d.dockerClient.ImagePull(ctx, providerSpec.Image, image.PullOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to pull image %q: %w", providerSpec.Image, err)
		}
	}

	name := containerName(req.Machine.Name)
	dockerNetwork := "kind"

	createResponse, err := d.dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Hostname: name,
			Image:    providerSpec.Image,
			Healthcheck: &container.HealthConfig{
				Test:     []string{"CMD", "sh", "-c", "/usr/bin/kubectl --kubeconfig /var/lib/kubelet/kubeconfig-real get no $NODE_NAME"},
				Interval: time.Second,
				Timeout:  time.Second,
				Retries:  3,
			},
			Env: []string{
				fmt.Sprintf("NODE_NAME=%s", name),
				"WAIT_FOR_USERDATA=true",
			},
			Labels: map[string]string{
				labelKeyProvider:      apiv1alpha1.Provider,
				labelKeyApp:           labelValueMachine,
				labelKeyMachine:       req.Machine.Name,
				labelMachineNamespace: req.Machine.Namespace,
			},
			Volumes: map[string]struct{}{
				"/etc/machine":                    {},
				"/var/lib/containerd":             {},
				"/var/lib/local-path-provisioner": {},
				"/var/lib/kubelet":                {},
				"/run":                            {},
				"/tmp":                            {},
			},
		},
		&container.HostConfig{
			NetworkMode:   container.NetworkMode(dockerNetwork),
			Privileged:    true,
			RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyOnFailure},
			Binds: []string{
				"/var/run/docker.sock:/var/run/docker.sock",
				"/lib/modules:/lib/modules:ro",
			},
			DNS: []string{"172.18.255.53", "fd00:ff::53"},
		},
		nil,
		nil,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create container %s: %w", name, err)
	}

	klog.V(1).Infof("Container %q successfully created", name)

	if err := d.dockerClient.ContainerStart(ctx, createResponse.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container %s: %w", createResponse.ID, err)
	}

	klog.V(1).Infof("Container %q successfully started", name)

	var inspectResponse container.InspectResponse

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	if err := wait.PollUntilContextCancel(timeoutCtx, 5*time.Second, true, func(ctx context.Context) (bool, error) {
		var err error

		inspectResponse, err = d.dockerClient.ContainerInspect(ctx, createResponse.ID)
		if err != nil {
			klog.V(1).Infof("error inspecting container %s: %v", createResponse.ID, err)
			return false, nil
		}

		if inspectResponse.State.Running {
			return true, nil
		}

		return false, nil
	}); err != nil {
		// will be retried with short retry by machine controller
		return nil, status.Error(codes.DeadlineExceeded, fmt.Sprintf("container %q is in phase %q, failed waiting for phase running: %v", name, inspectResponse.State.Status, err))
	}

	klog.V(1).Infof("Container %q is running", name)

	if err := d.copyResolvConf(ctx, createResponse.ID, name); err != nil {
		return nil, err
	}

	if err := d.copyUserData(ctx, createResponse.ID, name, req.Secret.Data["userData"]); err != nil {
		return nil, err
	}

	return &driver.CreateMachineResponse{
		ProviderID: name,
		NodeName:   name,
		Addresses:  addressesFromInspectResponse(inspectResponse, dockerNetwork),
	}, nil
}

func (d *localDriver) copyFileToContainer(ctx context.Context, containerID, destDir, fileName string, content []byte, mode int64) error {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	if err := tw.WriteHeader(&tar.Header{Name: fileName, Mode: mode, Size: int64(len(content))}); err != nil {
		return fmt.Errorf("failed to write %s tar header: %w", fileName, err)
	}
	if _, err := tw.Write(content); err != nil {
		return fmt.Errorf("failed to write %s tar content: %w", fileName, err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("failed to close %s tar writer: %w", fileName, err)
	}
	if err := d.dockerClient.CopyToContainer(ctx, containerID, destDir, &buf, container.CopyToContainerOptions{}); err != nil {
		return fmt.Errorf("failed to copy %s to container: %w", fileName, err)
	}
	return nil
}

func (d *localDriver) copyResolvConf(ctx context.Context, containerID, containerName string) error {
	content := []byte("nameserver 172.18.255.53\nnameserver fd00:ff::53\noptions ndots:0\n")
	if err := d.copyFileToContainer(ctx, containerID, "/tmp", "resolv.conf", content, 0644); err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to copy resolv.conf to container %q: %v", containerName, err))
	}

	execCfg, err := d.dockerClient.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		Cmd: []string{"cp", "/tmp/resolv.conf", "/etc/resolv.conf"},
	})
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to create exec for resolv.conf copy in container %q: %v", containerName, err))
	}
	if err := d.dockerClient.ContainerExecStart(ctx, execCfg.ID, container.ExecStartOptions{}); err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to exec resolv.conf copy in container %q: %v", containerName, err))
	}

	klog.V(1).Infof("Custom /etc/resolv.conf successfully written to container %q", containerName)
	return nil
}

func (d *localDriver) copyUserData(ctx context.Context, containerID, containerName string, userData []byte) error {
	if err := d.copyFileToContainer(ctx, containerID, "/etc/machine", "userdata", userData, 0700); err != nil {
		return status.Error(codes.Aborted, fmt.Sprintf("failed to copy user data to container %q: %v", containerName, err))
	}

	klog.V(1).Infof("Userdata successfully copied to container %q", containerName)
	return nil
}

func validateProviderSpecAndSecret(machineClass *machinev1alpha1.MachineClass, secret *corev1.Secret) (*apiv1alpha1.ProviderSpec, error) {
	if machineClass == nil {
		return nil, status.Error(codes.Internal, "MachineClass ProviderSpec is nil")
	}

	var providerSpec *apiv1alpha1.ProviderSpec
	if err := json.Unmarshal(machineClass.ProviderSpec.Raw, &providerSpec); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	validationErr := validation.ValidateProviderSpec(providerSpec, secret, field.NewPath("providerSpec"))
	if validationErr.ToAggregate() != nil && len(validationErr.ToAggregate().Errors()) > 0 {
		err := fmt.Errorf("error while validating ProviderSpec: %v", validationErr.ToAggregate().Error())
		klog.V(2).Infof("Validation of AWSMachineClass failed %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return providerSpec, nil
}

func addressesFromInspectResponse(inspectResponse container.InspectResponse, dockerNetworkName string) []corev1.NodeAddress {
	var addresses []corev1.NodeAddress

	if inspectResponse.NetworkSettings == nil || inspectResponse.NetworkSettings.Networks[dockerNetworkName] == nil {
		return addresses
	}

	network := inspectResponse.NetworkSettings.Networks[dockerNetworkName]

	if network.IPAddress != "" {
		addresses = append(addresses, corev1.NodeAddress{
			Type:    corev1.NodeInternalIP,
			Address: network.IPAddress,
		})
	}

	if network.GlobalIPv6Address != "" {
		addresses = append(addresses, corev1.NodeAddress{
			Type:    corev1.NodeInternalIP,
			Address: network.GlobalIPv6Address,
		})
	}

	return addresses
}
