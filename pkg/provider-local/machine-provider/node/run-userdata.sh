#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

userdata_path="${1:-/etc/machine/userdata}"
# if WAIT=true, wait until the userdata file exists:
if [[ "${WAIT_FOR_USERDATA:-}" == "true" ]]; then
  echo "Waiting for userdata file to be created at $userdata_path..."
  while [[ ! -f "$userdata_path" ]]; do
    sleep 1
  done
fi
if [[ -f "$userdata_path" ]]; then
  echo "Executing userdata at $userdata_path"
  "$userdata_path"
fi
