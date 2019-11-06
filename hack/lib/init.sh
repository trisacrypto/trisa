#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${REPO_ROOT}/hack/lib/artifacts.sh"
source "${REPO_ROOT}/hack/lib/bazel.sh"
source "${REPO_ROOT}/hack/lib/docs.sh"
source "${REPO_ROOT}/hack/lib/tooling.sh"
source "${REPO_ROOT}/hack/lib/pki.sh"
source "${REPO_ROOT}/hack/lib/demo.sh"
