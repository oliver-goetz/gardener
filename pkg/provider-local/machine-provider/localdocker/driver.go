// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package localdocker

import (
	"context"
	"strings"

	dockerclient "github.com/docker/docker/client"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("machine-provider-local")

const (
	fieldOwner            = client.FieldOwner("machine-controller-manager-provider-local")
	labelKeyApp           = "app"
	labelKeyProvider      = "machine-provider"
	labelKeyMachine       = "machine"
	labelValueMachine     = "machine"
	labelMachineNamespace = "machine-namespace"
	machinePrefix         = "machine-"
)

// NewDriver returns an empty AWSDriver object
func NewDriver(client client.Client, dockerClient *dockerclient.Client) driver.Driver {
	return &localDriver{runtimeClient: client, dockerClient: dockerClient}
}

type localDriver struct {
	runtimeClient client.Client

	dockerClient *dockerclient.Client
}

// GenerateMachineClassForMigration is not implemented.
func (d *localDriver) GenerateMachineClassForMigration(_ context.Context, _ *driver.GenerateMachineClassForMigrationRequest) (*driver.GenerateMachineClassForMigrationResponse, error) {
	return &driver.GenerateMachineClassForMigrationResponse{}, nil
}

// InitializeMachine is not implemented.
func (*localDriver) InitializeMachine(context.Context, *driver.InitializeMachineRequest) (*driver.InitializeMachineResponse, error) {
	return nil, status.Error(codes.Unimplemented, "InitializeMachine is not yet implemented")
}

func containerName(machineName string) string {
	return machinePrefix + machineName
}

func machineName(containerName string) string {
	return strings.TrimPrefix(containerName, machinePrefix)
}
