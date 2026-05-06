// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package infrastructuredocker

import (
	"context"
	"fmt"
	"strings"

	dockerclient "github.com/docker/docker/client"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/infrastructure"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
)

type actuator struct {
	// runtimeClient uses provider-local's in-cluster config, e.g., for the seed/bootstrap cluster it runs in.
	// It's used to interact with extension objects. By default, it's also used as the provider client to interact with
	// infrastructure resources, unless a kubeconfig is specified in the cloudprovider secret.
	runtimeClient client.Client

	dockerClient *dockerclient.Client
}

// NewActuator creates a new Actuator that updates the status of the handled Infrastructure resources.
func NewActuator(mgr manager.Manager, dockerclient *dockerclient.Client) infrastructure.Actuator {
	return &actuator{
		runtimeClient: mgr.GetClient(),
		dockerClient:  dockerclient,
	}
}

func (a *actuator) Reconcile(ctx context.Context, log logr.Logger, infrastructure *extensionsv1alpha1.Infrastructure, cluster *extensionscontroller.Cluster) error {
	//if cluster.Shoot.Spec.Networking == nil || cluster.Shoot.Spec.Networking.Nodes == nil {
	//	return fmt.Errorf("shoot specification does not contain node network CIDR required for VPN tunnel")
	//}
	//
	//// Get existing docker networks
	//networks, err := a.dockerClient.NetworkList(ctx, network.ListOptions{})
	//if err != nil {
	//	return fmt.Errorf("could not list docker networks: %w", err)
	//}
	//
	//networkExists := false
	//for _, network := range networks {
	//	if network.Name == GetDockerNetworkName(*cluster.Shoot.Spec.Networking.Nodes) {
	//		networkExists = true
	//	}
	//}
	//
	//if !networkExists {
	//	createOptions := network.CreateOptions{
	//		Driver: "bridge",
	//		Scope:  "local",
	//		IPAM: &network.IPAM{
	//			Driver: "default",
	//			Config: []network.IPAMConfig{
	//				{Subnet: *cluster.Shoot.Spec.Networking.Nodes},
	//			},
	//		},
	//	}
	//	for _, ipfamily := range cluster.Shoot.Spec.Networking.IPFamilies {
	//		switch ipfamily {
	//		case corev1beta1.IPFamilyIPv4:
	//			createOptions.EnableIPv4 = ptr.To(true)
	//		case corev1beta1.IPFamilyIPv6:
	//			createOptions.EnableIPv6 = ptr.To(true)
	//		}
	//	}
	//
	//	//var gatewayIP string
	//	//if strings.Contains(*cluster.Shoot.Spec.Networking.Nodes, ":") {
	//	//	gatewayIP = strings.Split(*cluster.Shoot.Spec.Networking.Nodes, "/")[0] + "1"
	//	//} else {
	//	//	withoutMask := strings.Split(*cluster.Shoot.Spec.Networking.Nodes, "/")[0]
	//	//	gatewayIP = withoutMask[:strings.LastIndex(withoutMask, ".")+1] + "1"
	//	//}
	//
	//	_, err := a.dockerClient.NetworkCreate(ctx, GetDockerNetworkName(*cluster.Shoot.Spec.Networking.Nodes), createOptions)
	//	if err != nil {
	//		return fmt.Errorf("could not create docker network: %w", err)
	//	}
	//}

	patch := client.MergeFrom(infrastructure.DeepCopy())
	infrastructure.Status.Networking = &extensionsv1alpha1.InfrastructureStatusNetworking{}
	if nodes := cluster.Shoot.Spec.Networking.Nodes; nodes != nil {
		infrastructure.Status.Networking.Nodes = []string{*nodes}
		// The egress CIDRs of local nodes are hard to define and depends on the traffic destination.
		// Traffic to the seed cluster will have nodes IPs as source IPs (i.e., the machine pod IPs).
		// Traffic to other containers in the kind network or the outside world will be NATed.
		// For now, we only report the nodes CIDR here to test the feature that propagates it back to the shoot status.
		infrastructure.Status.EgressCIDRs = []string{*nodes}
	}
	if pods := cluster.Shoot.Spec.Networking.Pods; pods != nil {
		infrastructure.Status.Networking.Pods = []string{*pods}
	}
	if services := cluster.Shoot.Spec.Networking.Services; services != nil {
		infrastructure.Status.Networking.Services = []string{*services}
	}

	return a.runtimeClient.Status().Patch(ctx, infrastructure, patch)
}

func (a *actuator) Delete(ctx context.Context, log logr.Logger, infrastructure *extensionsv1alpha1.Infrastructure, _ *extensionscontroller.Cluster) error {
	return nil
}

func (a *actuator) Migrate(context.Context, logr.Logger, *extensionsv1alpha1.Infrastructure, *extensionscontroller.Cluster) error {
	return nil
}

func (a *actuator) ForceDelete(ctx context.Context, log logr.Logger, infrastructure *extensionsv1alpha1.Infrastructure, cluster *extensionscontroller.Cluster) error {
	return a.Delete(ctx, log, infrastructure, cluster)
}

func (a *actuator) Restore(ctx context.Context, log logr.Logger, infrastructure *extensionsv1alpha1.Infrastructure, cluster *extensionscontroller.Cluster) error {
	return a.Reconcile(ctx, log, infrastructure, cluster)
}

// GetDockerNetworkName return the docker network name of the shoot node network
func GetDockerNetworkName(nodeCIDR string) string {
	return fmt.Sprintf("provider-local-%s", strings.NewReplacer(".", "-", ":", "-", "/", "-").Replace(nodeCIDR))
}
