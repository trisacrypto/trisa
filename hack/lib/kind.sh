#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KIND_CLUSTER=kind
KIND_CTX="kind-${KIND_CLUSTER}@kind"
KIND_IMAGE=kindest/node:v1.14.6
KIND_CONFIG=trisa-cluster-one-worker.yaml

# Start our local kind k8s cluster.
kind::cluster::start() {
    kind create cluster \
        --name ${KIND_CLUSTER} \
        --image ${KIND_IMAGE} \
        --config ${REPO_ROOT}/hack/etc/kind/${KIND_CONFIG}

    # Workaround until skaffold supports newer kind 0.6.0 ctx. When available, we can
    # also rename our `kind` cluster to `trisa` as KIND_CLUSTER name.
    kubectl config unset contexts.${KIND_CTX}
    kubectl config rename-context kind-${KIND_CLUSTER} ${KIND_CTX}

    kubectl config use-context ${KIND_CTX}
    kubectl cluster-info
}

# Destroy our local kind k8s cluster.
kind::cluster::destroy() {
    kind delete cluster --name ${KIND_CLUSTER}
}

# Ensure kubectl context is pointing to kind cluster
kind::ensure-ctx() {
    local current=$(kubectl config view -o template --template='{{ index . "current-context" }}')

    if [ "${KIND_CTX}" != "${current}" ]; then
        echo "Incorrect kubectl context: '${current}', expecting '${KIND_CTX}'"
        exit 1
    fi

    echo "Detected correct kubectl context '${KIND_CTX}'"
}

# Setup VASP configs collecting the server certificates and private keys and create the trust
# chain for kustomize to be able to pick them up when deploying to k8s.
kind::vasp::prepare-certificates() {
    cp hack/etc/pki/dev/out/vasp1.pem hack/etc/k8s/vasps/vasp1/server.crt
	cp hack/etc/pki/dev/out/vasp1-key.pem hack/etc/k8s/vasps/vasp1/server.key

    cp hack/etc/pki/dev/out/vasp2.pem hack/etc/k8s/vasps/vasp2/server.crt
	cp hack/etc/pki/dev/out/vasp2-key.pem hack/etc/k8s/vasps/vasp2/server.key

    cp hack/etc/pki/dev/out/vasp3.pem hack/etc/k8s/vasps/vasp3/server.crt
	cp hack/etc/pki/dev/out/vasp3-key.pem hack/etc/k8s/vasps/vasp3/server.key
	
	cat hack/etc/pki/dev/out/root.pem > hack/etc/k8s/vasps/vasp1/trust.chain
	cat hack/etc/pki/dev/out/subca1.pem >> hack/etc/k8s/vasps/vasp1/trust.chain
	cat hack/etc/pki/dev/out/subca2.pem >> hack/etc/k8s/vasps/vasp1/trust.chain

	cp hack/etc/k8s/vasps/vasp1/trust.chain hack/etc/k8s/vasps/vasp2
	cp hack/etc/k8s/vasps/vasp1/trust.chain hack/etc/k8s/vasps/vasp3
}