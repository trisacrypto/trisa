#!/usr/bin/env bash
# Build logic for dockerized tooling.

set -o errexit
set -o nounset
set -o pipefail

# Reference list of dockerized tools
TOOLING_GOHUGO=trisacrypto/tooling:gohugo
TOOLING_CFSSL=trisacrypto/tooling:cfssl
TOOLING_BAZEL=trisacrypto/tooling:bazel

# Switch to PR built images when running in non-fork pull_request mode on Travis.
# This will use the newly built images to test the CI execution.
TRAVIS_PULL_REQUEST_SLUG=${TRAVIS_PULL_REQUEST_SLUG:-}
if [ "${TRAVIS_PULL_REQUEST_SLUG}" == "trisacrypto/trisa" ]; then
    TOOLING_GOHUGO=${TOOLING_GOHUGO}-pr-${TRAVIS_PULL_REQUEST}
    TOOLING_CFSSL=${TOOLING_CFSSL}-pr-${TRAVIS_PULL_REQUEST}
    TOOLING_BAZEL=${TOOLING_BAZEL}-pr-${TRAVIS_PULL_REQUEST}
fi

# Bake all dockerized tooling
tooling::bake() {
    for dir in $(find ${REPO_ROOT}/hack/tooling -type d -mindepth 1); do
        tooling::skaffold ${dir}
    done
}

# Some testing to make sure our tooling containers work.
tooling::test() {
    echo "gohugo --> $(docker run -it --rm trisacrypto/tooling:gohugo version)"
    echo "cfssl -> $(docker run -it --rm trisacrypto/tooling:cfssl cfssl version)"
    echo "skaffold --> $(docker run -it --rm --entrypoint skaffold trisacrypto/tooling:bazel version)"
    docker run -it --rm --entrypoint bazel ${TOOLING_BAZEL} info
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

# Travis dockerized Bazel environment
tooling::travis::run() {

    local args=""
    local bazelrc=".bazelrc-travis"

    # Attach docker credentials if injected from Travis.
    if [ -f "/home/travis/.docker/config.json" ]; then
        args="${args} -v /home/travis/.docker/config.json:/home/bazel/.docker/config.json"
    fi

    # Google Credentials if injected from Travis.
    if [ ! -z "${GOOGLE_CREDENTIALS:-}" ]; then
        echo ${GOOGLE_CREDENTIALS} > .remote-cache-sa.json
        bazelrc=".bazelrc-travis-cache"
    fi

    # Pass TRAVIS_ env vars to container.
    env | grep TRAVIS_ > travis.env

    docker run --rm -it ${args} \
        -w /workspace \
        -v $(pwd):/workspace \
        -v $(pwd)/${bazelrc}:/home/bazel/.bazelrc \
        -v /var/run/docker.sock:/var/run/docker.sock \
        --env-file travis.env \
        ${TOOLING_BAZEL} -c "${*}"

    # Cleanup remote cache secret if any
    if [ -f ".remote-cache-sa.json" ]; then
        rm -f .remote-cache-sa.json
    fi
}
