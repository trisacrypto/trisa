#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DEMO_LOGS=${ARTIFACTS}/demo/logs

trisa::demo::init() {

    local pki=${REPO_ROOT}/hack/etc/pki/dev/out

    if [ ! -f ${pki}/vasp1.pem ]; then
        echo "Run 'make pki-dev-init' first"
        exit 1
    fi

    trisa::artifacts::clear

    mkdir -p ${ARTIFACTS}/bin
    
    mkdir -p ${ARTIFACTS}/demo/vasp1
    mkdir -p ${ARTIFACTS}/demo/vasp2
    mkdir -p ${ARTIFACTS}/demo/vasp3

    mkdir -p ${DEMO_LOGS}

    cat ${pki}/root.pem > ${ARTIFACTS}/demo/vasp1/trust.chain
    cat ${pki}/subca1.pem >> ${ARTIFACTS}/demo/vasp1/trust.chain
    cat ${pki}/subca2.pem >> ${ARTIFACTS}/demo/vasp1/trust.chain

    cp ${ARTIFACTS}/demo/vasp1/trust.chain ${ARTIFACTS}/demo/vasp2/trust.chain
    cp ${ARTIFACTS}/demo/vasp1/trust.chain ${ARTIFACTS}/demo/vasp3/trust.chain

    cp -f ${pki}/vasp1.pem ${ARTIFACTS}/demo/vasp1/server.crt
    cp -f ${pki}/vasp1-key.pem ${ARTIFACTS}/demo/vasp1/server.key

    cp -f ${pki}/vasp2.pem ${ARTIFACTS}/demo/vasp2/server.crt
    cp -f ${pki}/vasp2-key.pem ${ARTIFACTS}/demo/vasp2/server.key

    cp -f ${pki}/vasp3.pem ${ARTIFACTS}/demo/vasp3/server.crt
    cp -f ${pki}/vasp3-key.pem ${ARTIFACTS}/demo/vasp3/server.key
}

trisa::demo::build() {
    trisa::bazel::exec build //cmd/trisa
    # TODO: linux support, darwin harded for now
    cp -f $(trisa::bazel::info::bazel-bin)/cmd/trisa/darwin_amd64_stripped/trisa ${ARTIFACTS}/bin
}

trisa::demo::vasp::config-gen() {
    ${ARTIFACTS}/bin/trisa config generate --path=${ARTIFACTS}/demo/vasp1 --config ${ARTIFACTS}/demo/vasp1/config.yaml --listen=":8091" --listen-admin=":8591"
    ${ARTIFACTS}/bin/trisa config generate --path=${ARTIFACTS}/demo/vasp2 --config ${ARTIFACTS}/demo/vasp2/config.yaml --listen=":8092" --listen-admin=":8592"
    ${ARTIFACTS}/bin/trisa config generate --path=${ARTIFACTS}/demo/vasp3 --config ${ARTIFACTS}/demo/vasp3/config.yaml --listen=":8093" --listen-admin=":8593"
}

trisa::demo::start::vasps() {
    trisa::demo::stop::vasps
    ${ARTIFACTS}/bin/trisa server --config ${ARTIFACTS}/demo/vasp1/config.yaml &> ${DEMO_LOGS}/vasp1.log &
    ${ARTIFACTS}/bin/trisa server --config ${ARTIFACTS}/demo/vasp2/config.yaml &> ${DEMO_LOGS}/vasp2.log &
    ${ARTIFACTS}/bin/trisa server --config ${ARTIFACTS}/demo/vasp3/config.yaml &> ${DEMO_LOGS}/vasp3.log &
}

trisa::demo::stop::vasps() {
    killall trisa &> /dev/null || true
}