#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
source "${REPO_ROOT}/hack/lib/init.sh"

go mod tidy

trisa::bazel::exec run //:gazelle -- fix
trisa::bazel::exec run //:gazelle -- update-repos -from_file=go.mod
