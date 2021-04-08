#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

# Path to Envoy
ENVOY=${ENVOY:-/usr/local/bin/envoy}

# are there no input args set? then run our default
if [[ $# -eq 0 ]]; then
  echo "using default"
  ${ENVOY} -c  hack/bootstrap.yaml # --drain-time-s 1  # -l debug
else #take the first arg as the yaml to run
  echo "using bootstrap file: $1"
  ${ENVOY} -c  $1 # --drain-time-s 1  # -l debug
fi
