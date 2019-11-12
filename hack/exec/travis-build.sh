#!/usr/bin/env bash

# Dockerized build stage used by Travis

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
source "${REPO_ROOT}/hack/lib/init.sh"

trisa::bazel::exec build //cmd/trisa
trisa::bazel::exec test //pkg/...
trisa::bazel::exec run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/trisa:docker -- --norun