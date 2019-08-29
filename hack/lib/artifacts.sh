#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ARTIFACTS=${REPO_ROOT}/artifacts
mkdir -p ${ARTIFACTS}

trisa::artifacts::clear() {
   rm -rf ${ARTIFACTS}/*
}
