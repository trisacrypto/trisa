#!/usr/bin/env bash
# Build logic for dockerized tooling.

set -o errexit
set -o nounset
set -o pipefail

# Reference list of dockerized tools
TOOLING_GOHUGO=trisacrypto/tooling:gohugo
TOOLING_CFSSL=trisacrypto/tooling:cfssl

# Bake all dockerized tooling
tooling::bake() {
    for dir in $(find ${REPO_ROOT}/hack/tooling -type d -mindepth 1); do
        tooling::skaffold ${dir}
    done
}

# Build image using skaffold
tooling::skaffold() {
    local dir=${1}

    if [ ! -f "${dir}/skaffold.yaml" ]; then
        echo "no skaffold.yaml found in ${dir}"
        return
    fi

    cd ${dir} && skaffold build
}
