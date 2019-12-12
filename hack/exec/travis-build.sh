#!/usr/bin/env bash

# Dockerized build stage used by Travis

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
source "${REPO_ROOT}/hack/lib/init.sh"

# Same as `make all`.
hack/exec/build.sh
hack/exec/test.sh
hack/exec/build-docker.sh

# Startup demo environment.
hack/exec/pki-dev-init.sh
hack/exec/demo-init.sh
hack/exec/demo-run.sh

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

# Cleanup
hack/exec/demo-stop.sh
