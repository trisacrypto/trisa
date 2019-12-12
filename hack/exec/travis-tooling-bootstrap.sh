#!/usr/bin/env bash

# Dockerized build stage used by Travis

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
source "${REPO_ROOT}/hack/lib/init.sh"

# Build tooling in bazel container context as skaffold is available there.
tooling::travis::run-latest hack/exec/bake-tooling.sh

# Quick test on the resulting containers.
tooling::test
