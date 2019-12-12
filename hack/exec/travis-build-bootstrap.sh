#!/usr/bin/env bash

# Dockerized build stage used by Travis

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
source "${REPO_ROOT}/hack/lib/init.sh"

tooling::travis::run hack/exec/travis-build.sh

# Initialize PKI using two subcas and 3 VASPS.
pki::init::ca
pki::issue::subca 1
pki::issue::subca 2
pki::issue::end-entity::local vasp1 ../vasp1-csr.json subca1
pki::issue::end-entity::local vasp2 ../vasp2-csr.json subca1
pki::issue::end-entity::local vasp3 ../vasp3-csr.json subca2

# Demo environment initialization. Ensure no artifacts are cleaned as
# the trisa binary has already been built in the tooling container context.
trisa::demo::init
trisa::demo::vasp::config-gen
trisa::demo::start::vasps

# Generate some traffic.
curl -ks "https://127.0.0.1:8591/send?target=vasp3:8093" > /dev/null
curl -ks "https://127.0.0.1:8592/send?target=vasp3:8093" > /dev/null
curl -ks "https://127.0.0.1:8593/send?target=vasp1:8091" > /dev/null

# Dump traffic logs.
echo "*** VASP1 logs ***"
cat artifacts/demo/logs/vasp1.log
echo "*** VASP2 logs ***"
cat artifacts/demo/logs/vasp2.log
echo "*** VASP3 logs ***"
cat artifacts/demo/logs/vasp3.log

# Cleanup.
trisa::demo::stop::vasps
