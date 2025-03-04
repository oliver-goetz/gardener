//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	configv1alpha1 "github.com/gardener/gardener/pkg/gardenlet/apis/config/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	componentbaseconfigv1alpha1 "k8s.io/component-base/config/v1alpha1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConditionThreshold) DeepCopyInto(out *ConditionThreshold) {
	*out = *in
	out.Duration = in.Duration
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConditionThreshold.
func (in *ConditionThreshold) DeepCopy() *ConditionThreshold {
	if in == nil {
		return nil
	}
	out := new(ConditionThreshold)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControllerConfiguration) DeepCopyInto(out *ControllerConfiguration) {
	*out = *in
	in.Garden.DeepCopyInto(&out.Garden)
	in.GardenCare.DeepCopyInto(&out.GardenCare)
	in.GardenletDeployer.DeepCopyInto(&out.GardenletDeployer)
	in.NetworkPolicy.DeepCopyInto(&out.NetworkPolicy)
	in.VPAEvictionRequirements.DeepCopyInto(&out.VPAEvictionRequirements)
	in.Extension.DeepCopyInto(&out.Extension)
	in.ExtensionCare.DeepCopyInto(&out.ExtensionCare)
	in.ExtensionRequiredRuntime.DeepCopyInto(&out.ExtensionRequiredRuntime)
	in.ExtensionRequiredVirtual.DeepCopyInto(&out.ExtensionRequiredVirtual)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControllerConfiguration.
func (in *ControllerConfiguration) DeepCopy() *ControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionCareControllerConfiguration) DeepCopyInto(out *ExtensionCareControllerConfiguration) {
	*out = *in
	if in.SyncPeriod != nil {
		in, out := &in.SyncPeriod, &out.SyncPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	if in.ConditionThresholds != nil {
		in, out := &in.ConditionThresholds, &out.ConditionThresholds
		*out = make([]ConditionThreshold, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionCareControllerConfiguration.
func (in *ExtensionCareControllerConfiguration) DeepCopy() *ExtensionCareControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ExtensionCareControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionControllerConfiguration) DeepCopyInto(out *ExtensionControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionControllerConfiguration.
func (in *ExtensionControllerConfiguration) DeepCopy() *ExtensionControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ExtensionControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionRequiredRuntimeControllerConfiguration) DeepCopyInto(out *ExtensionRequiredRuntimeControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionRequiredRuntimeControllerConfiguration.
func (in *ExtensionRequiredRuntimeControllerConfiguration) DeepCopy() *ExtensionRequiredRuntimeControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ExtensionRequiredRuntimeControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionRequiredVirtualControllerConfiguration) DeepCopyInto(out *ExtensionRequiredVirtualControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionRequiredVirtualControllerConfiguration.
func (in *ExtensionRequiredVirtualControllerConfiguration) DeepCopy() *ExtensionRequiredVirtualControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ExtensionRequiredVirtualControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GardenCareControllerConfiguration) DeepCopyInto(out *GardenCareControllerConfiguration) {
	*out = *in
	if in.SyncPeriod != nil {
		in, out := &in.SyncPeriod, &out.SyncPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	if in.ConditionThresholds != nil {
		in, out := &in.ConditionThresholds, &out.ConditionThresholds
		*out = make([]ConditionThreshold, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GardenCareControllerConfiguration.
func (in *GardenCareControllerConfiguration) DeepCopy() *GardenCareControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(GardenCareControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GardenControllerConfig) DeepCopyInto(out *GardenControllerConfig) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	if in.SyncPeriod != nil {
		in, out := &in.SyncPeriod, &out.SyncPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	if in.ETCDConfig != nil {
		in, out := &in.ETCDConfig, &out.ETCDConfig
		*out = new(configv1alpha1.ETCDConfig)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GardenControllerConfig.
func (in *GardenControllerConfig) DeepCopy() *GardenControllerConfig {
	if in == nil {
		return nil
	}
	out := new(GardenControllerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GardenletDeployerControllerConfig) DeepCopyInto(out *GardenletDeployerControllerConfig) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GardenletDeployerControllerConfig.
func (in *GardenletDeployerControllerConfig) DeepCopy() *GardenletDeployerControllerConfig {
	if in == nil {
		return nil
	}
	out := new(GardenletDeployerControllerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkPolicyControllerConfiguration) DeepCopyInto(out *NetworkPolicyControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	if in.AdditionalNamespaceSelectors != nil {
		in, out := &in.AdditionalNamespaceSelectors, &out.AdditionalNamespaceSelectors
		*out = make([]v1.LabelSelector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkPolicyControllerConfiguration.
func (in *NetworkPolicyControllerConfiguration) DeepCopy() *NetworkPolicyControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(NetworkPolicyControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeTolerationConfiguration) DeepCopyInto(out *NodeTolerationConfiguration) {
	*out = *in
	if in.DefaultNotReadyTolerationSeconds != nil {
		in, out := &in.DefaultNotReadyTolerationSeconds, &out.DefaultNotReadyTolerationSeconds
		*out = new(int64)
		**out = **in
	}
	if in.DefaultUnreachableTolerationSeconds != nil {
		in, out := &in.DefaultUnreachableTolerationSeconds, &out.DefaultUnreachableTolerationSeconds
		*out = new(int64)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeTolerationConfiguration.
func (in *NodeTolerationConfiguration) DeepCopy() *NodeTolerationConfiguration {
	if in == nil {
		return nil
	}
	out := new(NodeTolerationConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OperatorConfiguration) DeepCopyInto(out *OperatorConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.RuntimeClientConnection = in.RuntimeClientConnection
	out.VirtualClientConnection = in.VirtualClientConnection
	in.LeaderElection.DeepCopyInto(&out.LeaderElection)
	in.Server.DeepCopyInto(&out.Server)
	if in.Debugging != nil {
		in, out := &in.Debugging, &out.Debugging
		*out = new(componentbaseconfigv1alpha1.DebuggingConfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.FeatureGates != nil {
		in, out := &in.FeatureGates, &out.FeatureGates
		*out = make(map[string]bool, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Controllers.DeepCopyInto(&out.Controllers)
	if in.NodeToleration != nil {
		in, out := &in.NodeToleration, &out.NodeToleration
		*out = new(NodeTolerationConfiguration)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OperatorConfiguration.
func (in *OperatorConfiguration) DeepCopy() *OperatorConfiguration {
	if in == nil {
		return nil
	}
	out := new(OperatorConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OperatorConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Server) DeepCopyInto(out *Server) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Server.
func (in *Server) DeepCopy() *Server {
	if in == nil {
		return nil
	}
	out := new(Server)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerConfiguration) DeepCopyInto(out *ServerConfiguration) {
	*out = *in
	out.Webhooks = in.Webhooks
	if in.HealthProbes != nil {
		in, out := &in.HealthProbes, &out.HealthProbes
		*out = new(Server)
		**out = **in
	}
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = new(Server)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServerConfiguration.
func (in *ServerConfiguration) DeepCopy() *ServerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ServerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VPAEvictionRequirementsControllerConfiguration) DeepCopyInto(out *VPAEvictionRequirementsControllerConfiguration) {
	*out = *in
	if in.ConcurrentSyncs != nil {
		in, out := &in.ConcurrentSyncs, &out.ConcurrentSyncs
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VPAEvictionRequirementsControllerConfiguration.
func (in *VPAEvictionRequirementsControllerConfiguration) DeepCopy() *VPAEvictionRequirementsControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(VPAEvictionRequirementsControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}
