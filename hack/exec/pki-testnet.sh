#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

PKI_PROFILE="testnet"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
source "${REPO_ROOT}/hack/lib/init.sh"

pki::init::ca

pki::issue::subca 1
pki::issue::subca 2

pki::issue::end-entity::local server ../server-csr.json subca1

pki::issue::end-entity::local vasp1 ../vasp1-csr.json subca1
pki::issue::end-entity::local vasp2 ../vasp2-csr.json subca1
pki::issue::end-entity::local vasp3 ../vasp3-csr.json subca2