#!/usr/bin/env bash
# SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o pipefail

COMMAND="${1:-up}"
VALID_COMMANDS=("up" "down")

SCENARIO="${SCENARIO:-bootstrap}"
declare -A SCENARIO_LEVEL=(
  [setup]=1     # Only starts the kind cluster and renders the manifests.
  [bootstrap]=2 # Like 'setup', but also runs `gardenadm bootstrap` and exports the kubeconfig for the self-hosted shoot
  [connect]=3   # Like 'join', but also deploys Gardener into the self-hosted shoot and runs `gardenadm connect` to deploy gardenlet which registers the Shoot
  [full]=4      # Like 'connect', but also registers the self-hosted shoot as a seed via a ManagedSeed
)

if [[ -z "${SCENARIO_LEVEL[$SCENARIO]+x}" ]]; then
  echo "Error: Invalid scenario '${SCENARIO}'. Valid options are: ${!SCENARIO_LEVEL[*]}." >&2
  exit 1
fi

level="${SCENARIO_LEVEL[$SCENARIO]}"

up() {
  if ! kind get clusters | grep gardener-local &>/dev/null; then
    make kind-up
  fi

  make gardenadm-up SCENARIO=managed-infra

  # Run `gardenadm bootstrap` and export the kubeconfig for the self-hosted shoot
  if ((level >= 2)); then
    make gardenadm
    local generated_dir="$(dirname "$0")/../dev-setup/gardenadm/resources/generated"
    export IMAGEVECTOR_OVERWRITE="$generated_dir/.imagevector-overwrite.yaml"
    export IMAGEVECTOR_OVERWRITE_CHARTS="$generated_dir/.imagevector-overwrite-charts.yaml"

    # if we can connect to the self-hosted shoot cluster, we assume that the bootstrap has already been done and skip
    # it. This allows re-running the script without re-bootstrapping, which takes a long time.
    if ! kubectl --kubeconfig="$KUBECONFIG_SELFHOSTEDSHOOT_CLUSTER" cluster-info &>/dev/null || ((level == 2)); then
      KUBECONFIG="$KUBECONFIG_RUNTIME_CLUSTER" "$(dirname "$0")/../bin/gardenadm" bootstrap -d "$generated_dir/managed-infra"
    fi

    tmp_exposure
    KUBECONFIG="$KUBECONFIG_RUNTIME_CLUSTER" ./hack/usage/generate-kubeconfig.sh self-hosted-shoot >"$KUBECONFIG_SELFHOSTEDSHOOT_CLUSTER"
  fi

  # Deploy Gardener into the self-hosted shoot and run `gardenadm connect` to deploy gardenlet which registers the Shoot
  if ((level >= 3)); then
    make gardenadm-up SCENARIO=connect-managed-infra # deploys gardener-operator, the 'Garden' resource, and waits for reconciliation
    connect_command="$(KUBECONFIG=$KUBECONFIG_VIRTUAL_GARDEN_CLUSTER "$(dirname "$0")/../bin/gardenadm" token create --print-connect-command --shoot-namespace garden --shoot-name root)"
    exec_controlplane "export IMAGEVECTOR_OVERWRITE=/var/lib/gardenadm/imagevector-overwrite.yaml; export IMAGEVECTOR_OVERWRITE_CHARTS=/var/lib/gardenadm/imagevector-overwrite-charts.yaml; /opt/bin/$connect_command"
  fi

  # Register the self-hosted shoot as a seed via a ManagedSeed
  if ((level >= 4)); then
    make seed-up KUBECONFIG="$KUBECONFIG_SELFHOSTEDSHOOT_CLUSTER"
  fi
}

down() {
  if kubectl --kubeconfig "$KUBECONFIG_VIRTUAL_GARDEN_CLUSTER" -n garden get managedseed root &>/dev/null; then
    make seed-down KUBECONFIG="$KUBECONFIG_SELFHOSTEDSHOOT_CLUSTER"
  fi

  make gardenadm-down SCENARIO=connect-managed-infra

  make kind-down
}

exec_controlplane() {
  local cmd="$*"
  local pod="$(kubectl --kubeconfig="$KUBECONFIG_RUNTIME_CLUSTER" get pod -n shoot--garden--root -o name | grep "control-plane" | cut -d/ -f2)"
  kubectl --kubeconfig="$KUBECONFIG_RUNTIME_CLUSTER" exec -it -n shoot--garden--root "$pod" -- sh -c "$cmd"
}

# TODO(maboehm): Remove after SelfHostedShoot Exposure is working
tmp_exposure() {
  local machine="$(kubectl --kubeconfig="$KUBECONFIG_RUNTIME_CLUSTER" get po -n shoot--garden--root -l app=machine -o yaml | yq '.items[]|select(.metadata.name|contains("control-plane")).metadata.labels.machine')"

  export KUBECONFIG="$KUBECONFIG_RUNTIME_CLUSTER"
  cat <<EOF | kubectl apply -o yaml -f -
apiVersion: v1
kind: Service
metadata:
  name: control-plane
  namespace: shoot--garden--root
spec:
  selector:
    machine: $machine
  ports:
    - protocol: TCP
      port: 443
      targetPort: 443
  type: LoadBalancer
EOF
  kubectl wait -n shoot--garden--root service/control-plane --for=jsonpath='{.status.loadBalancer.ingress}'
  ip="$(kubectl get svc -n shoot--garden--root control-plane -o yaml | yq '.status.loadBalancer.ingress[0].ip')"
  patch='{"spec":{"values":["'$ip'"]},"metadata":{"annotations":{"gardener.cloud/operation":"reconcile"}}}'
  exec_controlplane "KUBECONFIG=/etc/kubernetes/admin.conf kubectl patch dnsrecord -n kube-system root-external -o yaml --type=merge -p '$patch'"

  KUBECONFIG="$KUBECONFIG_RUNTIME_CLUSTER" ./hack/usage/generate-kubeconfig.sh self-hosted-shoot >"$KUBECONFIG_SELFHOSTEDSHOOT_CLUSTER"
  count=0
  until kubectl --kubeconfig="$KUBECONFIG_SELFHOSTEDSHOOT_CLUSTER" cluster-info --request-timeout=1s &>/dev/null; do
    sleep 1
    count=$((count + 1))
    if ((count > 60)); then
      echo "Error: Timed out waiting for the self-hosted shoot cluster to become available after control-plane exposure." >&2
      exit 1
    fi
  done
}

case "$COMMAND" in
up)
  up
  ;;

down)
  down
  ;;
*)
  echo "Error: Invalid command '${COMMAND}'. Valid options are: ${VALID_COMMANDS[*]}." >&2
  exit 1
  ;;
esac
