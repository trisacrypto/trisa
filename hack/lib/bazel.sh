#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

trisa::bazel::exec() {
    bazel ${*}
}

trisa::bazel::info::workspace() {
    printf $(trisa::bazel::exec info workspace)
}

trisa::bazel::info::bazel-bin() {
    printf $(trisa::bazel::exec info bazel-bin)
}
