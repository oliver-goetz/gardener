// Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package garden

import (
	"fmt"

	"github.com/Masterminds/semver"
	kubernetesclientset "k8s.io/client-go/kubernetes"
	"k8s.io/utils/clock"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	operatorv1alpha1 "github.com/gardener/gardener/pkg/apis/operator/v1alpha1"
	versionutils "github.com/gardener/gardener/pkg/utils/version"
)

// ControllerName is the name of this controller.
const ControllerName = "garden"

// AddToManager adds Reconciler to the given manager.
func (r *Reconciler) AddToManager(mgr manager.Manager) error {
	if r.RuntimeClient == nil {
		r.RuntimeClient = mgr.GetClient()
	}
	if r.Clock == nil {
		r.Clock = clock.RealClock{}
	}
	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor(ControllerName + "-controller")
	}
	if r.GardenNamespace == "" {
		r.GardenNamespace = v1beta1constants.GardenNamespace
	}

	if r.RuntimeVersion == nil {
		kubernetesClient, err := kubernetesclientset.NewForConfig(mgr.GetConfig())
		if err != nil {
			return fmt.Errorf("failed creating Kubernetes client: %w", err)
		}

		serverVersion, err := kubernetesClient.DiscoveryClient.ServerVersion()
		if err != nil {
			return fmt.Errorf("failed getting server version for runtime cluster: %w", err)
		}

		r.RuntimeVersion, err = semver.NewVersion(versionutils.Normalize(serverVersion.GitVersion))
		if err != nil {
			return fmt.Errorf("failed parsing version %q for runtime cluster: %w", serverVersion.GitVersion, err)
		}
	}

	return builder.
		ControllerManagedBy(mgr).
		Named(ControllerName).
		For(&operatorv1alpha1.Garden{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: pointer.IntDeref(r.Config.Controllers.Garden.ConcurrentSyncs, 0),
			RecoverPanic:            true,
		}).
		Complete(r)
}
