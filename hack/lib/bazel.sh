#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

BAZEL_MODE=${BAZEL_MODE:-host}
BAZEL_RC_FILE=${BAZEL_RC_FILE:-".bazelrc"}

trisa::bazel::exec() {
    if [ "${BAZEL_MODE}" = "docker" ]; then
        trisa::bazel::docker::exec --bazelrc=${BAZEL_RC_FILE} ${*}
    else
        bazel --bazelrc=${BAZEL_RC_FILE} ${*}
    fi
}

trisa::bazel::info::workspace() {
    printf $(trisa::bazel::exec info workspace)
}

trisa::bazel::info::bazel-bin() {
    printf $(trisa::bazel::exec info bazel-bin)
}

trisa::bazel::docker::exec() {
    docker run --rm \
    -w /workspace \
    -v $(pwd):/workspace \
    -v /var/run/docker.sock:/var/run/docker.sock \
    gcr.io/cloud-builders/bazel ${*}
}