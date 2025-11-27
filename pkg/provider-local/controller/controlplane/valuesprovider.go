// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package controlplane

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"net"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/controlplane/genericactuator"
	extensionssecretsmanager "github.com/gardener/gardener/extensions/pkg/util/secret/manager"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/component/kubernetes/apiserver"
	"github.com/gardener/gardener/pkg/provider-local/charts"
	localimagevector "github.com/gardener/gardener/pkg/provider-local/imagevector"
	"github.com/gardener/gardener/pkg/provider-local/local"
	"github.com/gardener/gardener/pkg/utils/chart"
	kubernetesutils "github.com/gardener/gardener/pkg/utils/kubernetes"
	"github.com/gardener/gardener/pkg/utils/secrets"
	secretsutils "github.com/gardener/gardener/pkg/utils/secrets"
	secretsmanager "github.com/gardener/gardener/pkg/utils/secrets/manager"
)

const caNameControlPlane = "ca-" + local.Name + "-controlplane"

// NewValuesProvider creates a new ValuesProvider for the generic actuator.
func NewValuesProvider(client client.Client) genericactuator.ValuesProvider {
	return &valuesProvider{client: client}
}

type valuesProvider struct {
	genericactuator.NoopValuesProvider
	client client.Client
}

func getSecretConfigs(namespace string) []extensionssecretsmanager.SecretConfigWithOptions {
	return []extensionssecretsmanager.SecretConfigWithOptions{
		{
			Config: &secretsutils.CertificateSecretConfig{
				Name:       caNameControlPlane,
				CommonName: caNameControlPlane,
				CertType:   secretsutils.CACert,
			},
			Options: []secretsmanager.GenerateOption{secretsmanager.Persist()},
		},
		{
			Config: &secretsutils.CertificateSecretConfig{
				Name:       local.Name + "-dummy-server",
				CommonName: local.Name + "-dummy-server",
				DNSNames:   kubernetesutils.DNSNamesForService(local.Name+"-dummy-server", namespace),
				CertType:   secretsutils.ServerCert,
			},
			Options: []secretsmanager.GenerateOption{secretsmanager.SignedByCA(caNameControlPlane)},
		},
		{
			Config: &secretsutils.CertificateSecretConfig{
				Name:                        local.Name + "-dummy-client",
				CommonName:                  "extensions.gardener.cloud:" + local.Name + ":dummy-client",
				Organization:                []string{"extensions.gardener.cloud:" + local.Name + ":dummy"},
				CertType:                    secretsutils.ClientCert,
				SkipPublishingCACertificate: true,
			},
			Options: []secretsmanager.GenerateOption{secretsmanager.SignedByCA(caNameControlPlane)},
		},
		{
			Config: &secretsutils.BasicAuthSecretConfig{
				Name:           local.Name + "-dummy-auth",
				Format:         secretsutils.BasicAuthFormatNormal,
				PasswordLength: 32,
			},
			Options: []secretsmanager.GenerateOption{secretsmanager.Validity(time.Hour)},
		},
	}
}

var (
	controlPlaneShootChart = &chart.Chart{
		Name:       "shoot-system-components",
		EmbeddedFS: charts.ChartShootSystemComponents,
		Path:       charts.ChartPathShootSystemComponents,
		SubCharts: []*chart.Chart{
			{
				Name: "local-path-provisioner",
				Images: []string{
					localimagevector.ImageNameLocalPathProvisioner,
					localimagevector.ImageNameLocalPathHelper,
				},
			},
		},
	}

	storageClassChart = &chart.Chart{
		Name:       "shoot-storageclasses",
		EmbeddedFS: charts.ChartShootStorageClasses,
		Path:       charts.ChartPathShootStorageClasses,
	}
	controlPlaneChart = &chart.Chart{
		Name:       "shoot-control-plane-components",
		EmbeddedFS: charts.ChartShootControlPlaneComponents,
		Path:       charts.ChartPathShootControlPlaneComponents,
	}
)

// GetControlPlaneShootChartValues returns the values for the control plane shoot chart applied by the generic actuator.
func (vp *valuesProvider) GetControlPlaneShootChartValues(
	_ context.Context,
	_ *extensionsv1alpha1.ControlPlane,
	_ *extensionscontroller.Cluster,
	_ secretsmanager.Reader,
	_ map[string]string,
) (map[string]any, error) {
	return map[string]any{}, nil
}

func (vp *valuesProvider) GetControlPlaneChartValues(ctx context.Context, cp *extensionsv1alpha1.ControlPlane, cluster *extensionscontroller.Cluster, secretsReader secretsmanager.Reader, checksums map[string]string, scaledDown bool) (map[string]any, error) {
	const talosSecretName = "talos"

	values := map[string]any{}

	var talosSecret corev1.Secret
	if err := vp.client.Get(ctx, client.ObjectKey{Name: talosSecretName, Namespace: cp.Namespace}, &talosSecret); err != nil {
		if errors.IsNotFound(err) {
			values, err = vp.createValues(ctx, cluster)
			if err != nil {
				return nil, fmt.Errorf("unable to create talos bootstrap values: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unable to get talos bootstrap secret: %w", err)
		}
	}

	// TODO!!!!!!!!

	return values, nil
}

func (vp *valuesProvider) createValues(ctx context.Context, cluster *extensionscontroller.Cluster) (map[string]any, error) {
	values := map[string]any{}

	talosToken, err := genToken(6, 16)
	if err != nil {
		return nil, err
	}

	values["talosToken"] = talosToken

	var (
		internalDomain string
		externalDomain string
	)

	for _, address := range cluster.Shoot.Status.AdvertisedAddresses {
		switch address.Name {
		case "internal":
			internalDomain = strings.TrimPrefix(address.URL, "https://")
		case "external":
			externalDomain = strings.TrimPrefix(address.URL, "https://")
		}
	}

	var istioLBs corev1.ServiceList

	if err := vp.client.List(ctx, &istioLBs, client.MatchingLabels{"app": "istio-ingressgateway"}); err != nil {
		return nil, fmt.Errorf("unable to list istio load balancer services: %w", err)
	}

	var ips []net.IP
	for _, svc := range istioLBs.Items {
		if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
			for _, clusterIP := range svc.Spec.ClusterIPs {
				if clusterIP != corev1.ClusterIPNone {
					ips = append(ips, net.ParseIP(clusterIP))
				}
			}
			for _, ingress := range svc.Status.LoadBalancer.Ingress {
				if ingress.IP != "" {
					ips = append(ips, net.ParseIP(ingress.IP))
				}
			}
		}
	}

	cert, ca, dir, err := secrets.SelfGenerateTLSServerCertificate("talos",
		[]string{"localhost", internalDomain, externalDomain}, ips)

	return values, nil
}

func randBootstrapTokenString(length int) (string, error) {
	// validBootstrapTokenChars defines the characters a bootstrap token can consist of
	const validBootstrapTokenChars = "0123456789abcdefghijklmnopqrstuvwxyz"

	// len("0123456789abcdefghijklmnopqrstuvwxyz") = 36 which doesn't evenly divide
	// the possible values of a byte: 256 mod 36 = 4. Discard any random bytes we
	// read that are >= 252 so the bytes we evenly divide the character set.
	const maxByteValue = 252

	var (
		b     byte
		err   error
		token = make([]byte, length)
	)

	reader := bufio.NewReaderSize(rand.Reader, length*2)

	for i := range token {
		for {
			if b, err = reader.ReadByte(); err != nil {
				return "", err
			}

			if b < maxByteValue {
				break
			}
		}

		token[i] = validBootstrapTokenChars[int(b)%len(validBootstrapTokenChars)]
	}

	return string(token), err
}

func genToken(lenFirst, lenSecond int) (string, error) {
	var err error

	tokenTemp := make([]string, 2)

	tokenTemp[0], err = randBootstrapTokenString(lenFirst)
	if err != nil {
		return "", err
	}

	tokenTemp[1], err = randBootstrapTokenString(lenSecond)
	if err != nil {
		return "", err
	}

	return tokenTemp[0] + "." + tokenTemp[1], nil
}
