// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	. "github.com/gardener/gardener/pkg/apis/settings/v1alpha1"
)

var _ = Describe("OpenIDConnectPreset defaulting", func() {
	It("should default OpenIDConnectPreset correctly", func() {
		obj := &OpenIDConnectPreset{}
		expected := &OpenIDConnectPreset{
			Spec: OpenIDConnectPresetSpec{
				Server: KubeAPIServerOpenIDConnect{
					// string literal are used to be sure that the test fails
					// if the constant values are changed.
					UsernameClaim: ptr.To("sub"),
					SigningAlgs:   []string{"RS256"},
				},
				ShootSelector: &metav1.LabelSelector{},
			},
		}
		SetObjectDefaults_OpenIDConnectPreset(obj)

		Expect(obj).To(Equal(expected))
	})

	It("should not default OpenIDConnectPreset if it is already set", func() {
		obj := &OpenIDConnectPreset{
			Spec: OpenIDConnectPresetSpec{
				Server: KubeAPIServerOpenIDConnect{
					// string literal are used to be sure that the test fails
					// if the constant values are changed.
					UsernameClaim: ptr.To("usr"),
					SigningAlgs:   []string{"alg1", "alg2"},
				},
				ShootSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
			},
		}
		expected := obj.DeepCopy()
		SetObjectDefaults_OpenIDConnectPreset(obj)

		Expect(obj).To(Equal(expected))
	})

})
