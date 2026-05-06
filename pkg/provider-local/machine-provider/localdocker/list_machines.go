// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package localdocker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	"k8s.io/klog/v2"

	apiv1alpha1 "github.com/gardener/gardener/pkg/provider-local/machine-provider/api/v1alpha1"
)

func (d *localDriver) ListMachines(ctx context.Context, req *driver.ListMachinesRequest) (*driver.ListMachinesResponse, error) {
	if req.MachineClass.Provider != apiv1alpha1.Provider {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("requested for Provider '%s', we only support '%s'", req.MachineClass.Provider, apiv1alpha1.Provider))
	}

	klog.V(3).Infof("Machine list request has been received for %q", req.MachineClass.Name)
	defer klog.V(3).Infof("Machine list request has been processed for %q", req.MachineClass.Name)

	summaries, err := d.dockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	machineList := make(map[string]string, len(summaries))
	for _, summary := range summaries {
		if summary.Labels[labelKeyProvider] != apiv1alpha1.Provider ||
			summary.Labels[labelKeyApp] != labelValueMachine ||
			summary.Labels[labelMachineNamespace] != req.MachineClass.Namespace {
			continue
		}

		machineList[summary.Names[0]] = machineName(summary.Names[0])
	}

	return &driver.ListMachinesResponse{MachineList: machineList}, nil
}
