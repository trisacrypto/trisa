#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

trisa::docs::dev() {
    docker run --rm -it -v ${REPO_ROOT}/docs:/docs -w /docs -p 1313:1313 \
        --user $(id -u) ${TOOLING_GOHUGO} server --bind 0.0.0.0
}

trisa::docs::generate() {
    docker run --rm -it -v ${REPO_ROOT}/docs:/docs -w /docs \
        --user $(id -u) ${TOOLING_GOHUGO}
}
