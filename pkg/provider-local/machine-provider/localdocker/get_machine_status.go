// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package localdocker

import (
	"context"
	"fmt"

	"github.com/containerd/errdefs"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	"k8s.io/klog/v2"

	apiv1alpha1 "github.com/gardener/gardener/pkg/provider-local/machine-provider/api/v1alpha1"
)

func (d *localDriver) GetMachineStatus(ctx context.Context, req *driver.GetMachineStatusRequest) (*driver.GetMachineStatusResponse, error) {
	if req.MachineClass.Provider != apiv1alpha1.Provider {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("requested for Provider '%s', we only support '%s'", req.MachineClass.Provider, apiv1alpha1.Provider))
	}

	klog.V(3).Infof("Machine status request has been received for %q", req.Machine.Name)
	defer klog.V(3).Infof("Machine status request has been processed for %q", req.Machine.Name)

	providerSpec, err := validateProviderSpecAndSecret(req.MachineClass, req.Secret)
	if err != nil {
		return nil, err
	}

	name := containerName(req.Machine.Name)

	inspectResponse, err := d.dockerClient.ContainerInspect(ctx, name)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("container %s not found: %v", name, err))
		}
		return nil, fmt.Errorf("failed to inspect container %s: %w", name, err)
	}

	return &driver.GetMachineStatusResponse{
		ProviderID: name,
		NodeName:   name,
		Addresses:  addressesFromInspectResponse(inspectResponse, providerSpec.DockerNetworkName),
	}, nil
}
