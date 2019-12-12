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
    mkdir -p ${ARTIFACTS}/bin
    cp -f $(trisa::bazel::info::bazel-bin)/cmd/trisa/${PLATFORM}_amd64_stripped/trisa ${ARTIFACTS}/bin
}

trisa::demo::vasp::config-gen() {
    ${ARTIFACTS}/bin/trisa config generate --path=${ARTIFACTS}/demo/vasp1 --config ${ARTIFACTS}/demo/vasp1/config.yaml --listen=":8091" --listen-admin=":8591"
    ${ARTIFACTS}/bin/trisa config generate --path=${ARTIFACTS}/demo/vasp2 --config ${ARTIFACTS}/demo/vasp2/config.yaml --listen=":8092" --listen-admin=":8592"
    ${ARTIFACTS}/bin/trisa config generate --path=${ARTIFACTS}/demo/vasp3 --config ${ARTIFACTS}/demo/vasp3/config.yaml --listen=":8093" --listen-admin=":8593"
}

trisa::demo::vasp::config-gen-docker() {
    docker run -it --rm -v ${ARTIFACTS}/demo/vasp1:/etc/trisa trisacrypto/trisa:latest config generate --listen=":8091" --listen-admin=":8591"
    docker run -it --rm -v ${ARTIFACTS}/demo/vasp2:/etc/trisa trisacrypto/trisa:latest config generate --listen=":8092" --listen-admin=":8592"
    docker run -it --rm -v ${ARTIFACTS}/demo/vasp3:/etc/trisa trisacrypto/trisa:latest config generate --listen=":8093" --listen-admin=":8593"
}

trisa::demo::start::vasps() {
    trisa::demo::stop::vasps
    cd ${ARTIFACTS}/demo/vasp1 && ../../bin/trisa server --config config.yaml &> ${DEMO_LOGS}/vasp1.log &
    cd ${ARTIFACTS}/demo/vasp2 && ../../bin/trisa server --config config.yaml &> ${DEMO_LOGS}/vasp2.log &
    cd ${ARTIFACTS}/demo/vasp3 && ../../bin/trisa server --config config.yaml &> ${DEMO_LOGS}/vasp3.log &
    sleep 5
    echo "VASP servers started. Log output can be found under ${ARTIFACTS}/demo/logs"
    ls -l ${ARTIFACTS}/demo/logs
}

trisa::demo::stop::vasps() {
    killall trisa &> /dev/null || true
}

trisa::demo::docker::up() {
    docker-compose --project-directory . -f hack/etc/demo/docker-compose.yml up
}

trisa::demo::docker::down() {
    docker-compose --project-directory . -f hack/etc/demo/docker-compose.yml down
}