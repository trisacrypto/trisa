#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Force dev profile only
PKI_PROFILE="dev"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
source "${REPO_ROOT}/hack/lib/init.sh"

# pki-dev-init is required, this multirootca is just for local testing purposes.
pki::issue::server subca1
pki::server
