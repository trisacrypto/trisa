#!/usr/bin/env bash

# Dockerized build stage used by Travis

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
source "${REPO_ROOT}/hack/lib/init.sh"

# Same as `make all`.
hack/exec/build.sh
hack/exec/test.sh
hack/exec/build-docker.sh

# Build demo binary for later consumption.
trisa::artifacts::clear
trisa::demo::build
