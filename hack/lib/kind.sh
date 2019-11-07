#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

KIND_CLUSTER=kind
KIND_IMAGE=kindest/node:v1.14.6
KIND_CONFIG=trisa-cluster-one-worker.yaml

# Start our local kind k8s cluster.
kind::cluster::start() {
    kind create cluster \
        --name ${KIND_CLUSTER} \
        --image ${KIND_IMAGE} \
        --config ${REPO_ROOT}/hack/etc/kind/${KIND_CONFIG}

    kind::cluster::config::merge

    echo
    echo "***** Please ignore the message above to set KUBECONFIG *****"
    echo
    echo "Your configuration is already in place and the kubectl context"
    echo "has been set to 'kubernetes-admin@kind' for you."
    echo

    kubectl config use-context kubernetes-admin@kind
    kubectl cluster-info
}

# Destroy our local kind k8s cluster.
kind::cluster::destroy() {
    kind delete cluster --name ${KIND_CLUSTER}
}

# Some trickery to merge our local cluster config into the users existing kube config.
# This will cause issues if there are overlaps on the context configuration. Users
# are on their own if that's the case (which mostly means they know what they are doing).
kind::cluster::config::merge() {

    # Some housekeeping
    mkdir -p ${HOME}/.kube
    touch ${HOME}/.kube/config

    # Make a backup of current config
    cp ${HOME}/.kube/config ${HOME}/.kube/config-hack

    # Use KUBECONFIG to merge config-hack and our new one together
    export KUBECONFIG="$(kind get kubeconfig-path --name=${KIND_CLUSTER}):${HOME}/.kube/config-hack"
    kubectl config view --raw > ${HOME}/.kube/config
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