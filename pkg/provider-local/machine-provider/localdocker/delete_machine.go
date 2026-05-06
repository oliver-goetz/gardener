// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package localdocker

import (
	"context"
	"fmt"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	apiv1alpha1 "github.com/gardener/gardener/pkg/provider-local/machine-provider/api/v1alpha1"
)

func (d *localDriver) DeleteMachine(ctx context.Context, req *driver.DeleteMachineRequest) (*driver.DeleteMachineResponse, error) {
	if req.MachineClass.Provider != apiv1alpha1.Provider {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("requested for Provider '%s', we only support '%s'", req.MachineClass.Provider, apiv1alpha1.Provider))
	}

	klog.V(3).Infof("Machine deletion request has been received for %q", req.Machine.Name)
	defer klog.V(3).Infof("Machine deletion request has been processed for %q", req.Machine.Name)

	name := containerName(req.Machine.Name)

	inspectResponse, err := d.dockerClient.ContainerInspect(ctx, name)
	if err != nil {
		if errdefs.IsNotFound(err) {
			klog.V(1).Infof("Container %q not found, nothing to do", name)
			return &driver.DeleteMachineResponse{}, nil
		}
		return nil, fmt.Errorf("failed to inspect container %s: %w", name, err)
	}

	klog.V(1).Infof("Container %q found", name)

	if inspectResponse.State.Running {
		klog.V(1).Infof("Container %q is running, stopping it", name)

		err := d.dockerClient.ContainerStop(
			ctx,
			containerName(req.Machine.Name),
			container.StopOptions{
				Timeout: ptr.To(600),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("could not stop container %q: %w", containerName(req.Machine.Name), err)
		}
	}

	klog.V(1).Infof("Removing container %q", name)

	err = d.dockerClient.ContainerRemove(
		ctx,
		containerName(req.Machine.Name),
		container.RemoveOptions{
			RemoveVolumes: true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not remove container %q: %w", containerName(req.Machine.Name), err)
	}

	klog.V(1).Infof("Container %q removed", containerName(req.Machine.Name))

	return &driver.DeleteMachineResponse{}, nil
}
