#!/usr/bin/env bash
# PKI Management for local and testing using the TRISA Testnet

set -o errexit
set -o nounset
set -o pipefail

# The PKI_PROFILE defines from which directory the PKI configs are read
# under hack/etc/pki. Defaults to "dev" for local development.
PKI_DIR=${REPO_ROOT}/hack/etc/pki/${PKI_PROFILE:-dev}
PKI_OUT=${PKI_DIR}/out

# Initialize a new root CA.
pki::init::ca() {
    mkdir -p ${PKI_OUT}
    rm -rf ${PKI_OUT}/*

    echo "Generating CA from ${PKI_DIR}"

    pki::cfssl::run "cfssl gencert -initca ../root-csr.json | cfssljson -bare root"
}

# Issue a new issuing CA (aka intermediate/subordinate CA).
pki::issue::subca() {
    local number=${1}

    # Generate CSR from template
    sed -e "s/%%NAME%%/Issuing CA ${number}/" ${PKI_DIR}/subca-csr.json > ${PKI_OUT}/subca${number}-csr.json

    # Setup API keys and generate config from template
    local key=$(pki::generate-key)
    echo ${key} > ${PKI_OUT}/subca${number}-key.api
    sed -e "s/%%KEY%%/${key}/" ${PKI_DIR}/subca-config.json > ${PKI_OUT}/subca${number}-config.json

    # Generate private key and CSR
    pki::cfssl::run "cfssl genkey subca${number}-csr.json | cfssljson -bare subca${number}"

    # Sign CSR from root CA
    pki::cfssl::run "cfssl sign -ca root.pem -ca-key root-key.pem --config subca${number}-config.json subca${number}.csr | cfssljson -bare subca${number}"

    # Certificate chain
    cat ${PKI_OUT}/subca${number}.pem > ${PKI_OUT}/subca${number}-chain.pem
    cat ${PKI_OUT}/root.pem >> ${PKI_OUT}/subca${number}-chain.pem

    # Attach profle to server
    cat << EOF >> ${PKI_OUT}/server.ini
[subca${number}]
private = file://subca${number}-key.pem
certificate = subca${number}.pem
config = subca${number}-config.json
EOF
}

# Issue server certificate
pki::issue::server() {
    local subca=${1}

    pki::cfssl::run "cfssl genkey ../server-csr.json | cfssljson -bare server"
    pki::cfssl::run "cfssl sign -ca ${subca}.pem -ca-key ${subca}-key.pem --config ../end-entity-config.json server.csr | cfssljson -bare server"
}

# Issue end-entity certicate locally
pki::issue::end-entity::local() {
    local name=${1}
    local csr=${2}
    local subca=${3}

    pki::cfssl::run "cfssl genkey ${csr} | cfssljson -bare ${name}"
    pki::cfssl::run "cfssl sign -ca ${subca}.pem -ca-key ${subca}-key.pem --config ../end-entity-config.json ${name}.csr | cfssljson -bare ${name}"
}

pki::issue::end-entity::remote() {
    echo "not implemented"
}

# Run a local cfssl server using multirootca config. This requires the server keys to be generated to secure the key exchange.
pki::server() {
    pki::cfssl::run "multirootca -a 0.0.0.0:8765 -roots server.ini -tls-cert server.pem -tls-key server-key.pem" "--name cfssl-server -p 8765:8765"
}

# Generate random 16 bytes hex string
pki::generate-key() {
    echo $(hexdump -n 16 -e '4/4 "%08X" ' /dev/random)
}

# Wrapper running cfssl through docker.
pki::cfssl::run() {
    local cmd=${1}
    local dopts=${2:-}

    docker run -it --rm --user $(id -u):$(id -g) ${dopts} \
        -v ${PKI_DIR}:/ca -w /ca/out \
        ${TOOLING_CFSSL} /bin/bash -c "${cmd}"
}
