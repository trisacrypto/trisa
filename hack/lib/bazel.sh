#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

BAZEL_MODE=${BAZEL_MODE:-host}

trisa::bazel::exec() {
    if [ "${BAZEL_MODE}" = "docker" ]; then
        trisa::bazel::docker::exec ${*}
    else
        bazel ${*}
    fi
}

trisa::bazel::info::workspace() {
    printf $(trisa::bazel::exec info workspace)
}

trisa::bazel::info::bazel-bin() {
    printf $(trisa::bazel::exec info bazel-bin)
}

trisa::bazel::docker::exec() {
    docker run --rm -v $(pwd):/workspace -w /workspace gcr.io/cloud-builders/bazel ${*}
}