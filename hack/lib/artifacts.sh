#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ARTIFACTS=${REPO_ROOT}/artifacts
ARTIFACTS_CLEAR_MODE="clear"

mkdir -p ${ARTIFACTS}

trisa::artifacts::clear() {
   if [ "${ARTIFACTS_CLEAR_MODE}" == "clear" ]; then
      rm -rf ${ARTIFACTS}/*
   fi
}
