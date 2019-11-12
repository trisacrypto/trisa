#!/usr/bin/env bash

# Dockerized build stage used by Travis

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
source "${REPO_ROOT}/hack/lib/init.sh"

# Force usage of cache config as we can't trap the bazel command from
# skaffold. Also it seems only build arguments can be passed so no
# cookies for using --bazelrc=.bazelrc-travis-cache.
cp .bazelrc-travis-cache .bazelrc

tooling::travis::run skaffold build -v info -p ${1}