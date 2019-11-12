#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

BAZEL_RC_FILE=${BAZEL_RC_FILE:-".bazelrc"}

trisa::bazel::exec() {
    bazel --bazelrc=${BAZEL_RC_FILE} ${*}
}

trisa::bazel::info::workspace() {
    printf $(trisa::bazel::exec info workspace)
}

trisa::bazel::info::bazel-bin() {
    printf $(trisa::bazel::exec info bazel-bin)
}
