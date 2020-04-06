.PHONY: all
all: build test build-docker

# Build TRISA server.
.PHONY: build
build:
	hack/exec/build.sh

# Build docker container for TRISA server. The resulting images will be pushed
# to the local docker instance as "bazel/cmd/trisa:docker".
.PHONY: build-docker
build-docker:
	hack/exec/build-docker.sh

# Bake dokerized tooling.
.PHONY: bake-tooling
bake-tooling:
	hack/exec/bake-tooling.sh

# Pull dokerized tooling.
.PHONY: pull-tooling
pull-tooling:
	hack/exec/pull-tooling.sh

# Run test suites.
.PHONY: test
test:
	hack/exec/test.sh

# Update/rebuild BUILD.bazel files for Go dependency management.
.PHONY: gazelle
gazelle:
	hack/exec/gazelle.sh

# Run 'go mod' to update vendor. Implicitly runs the gazelle target as well.
.PHONY: dependencies
dependencies:
	hack/exec/dependencies.sh

# Start local documentation service, see https://trisacrypto.github.io/contributing/documentation/.
.PHONY: docs-dev
docs-dev:
	hack/exec/docs-dev.sh

# Generate documentation website, usually run from CI only.
.PHONY: docs-generate
docs-generate:
	hack/exec/docs-generate.sh

# Initialize the PKI environment for local development. This create a new root CA,
# two subordinate CA's and 3 VASP certificates.
.PHONY: pki-dev-init
pki-dev-init:
	hack/exec/pki-dev-init.sh
	ls -l hack/etc/pki/dev/out

# Execute TRISA server in docker container. Requires the "build-docker" target to
# be execute first to produce the docker image.
.PHONY: docker-run
docker-run:
	docker run -it --rm bazel/cmd/trisa:docker

# Publish docker release as trisacrypto/trisa:latest.
.PHONY: docker-release
docker-release:
	skaffold build -p latest

# Initialize the demo environment. This requires "pki-dev-init" to be executed as
# it relies on the PKI setup to generate the demo system.
.PHONY: demo-init
demo-init:
	hack/exec/demo-init.sh

# Run the 3 demo VAPs locally.
.PHONY: demo-run
demo-run:
	hack/exec/demo-run.sh

# Rebuild TRISA server and restart the running demo VASPs.
.PHONY: demo-rebuild
demo-rebuild:
	hack/exec/demo-rebuild.sh

# Stop all running (demo) TRISA server processes.
.PHONY: demo-stop
demo-stop:
	hack/exec/demo-stop.sh

# Demo using docker and docker-compose only. This does not require setting up any
# build system and solely relies the dockerized tooling and published trisa image.
#
# This target should not be used for development as it will blow away the PKI dev
# setup, regenerate the configs and will not rely on any locally made code changes.
#
# Once the VASP containers are up and running, message exchanges can be triggered
# by hitting the admin port on each VASP server. The admin ports for each VASP are
# follows:
#
#	vasp1 --> 8591
#	vasp2 --> 8592
#	vasp3 --> 8593
#
# The TRISA mesh ports are as follows:
#
#	vasp1 --> 8091
#	vasp2 --> 8092
#	vasp3 --> 8093
#
# To trigger a transaction exchange from VASP1 to VASP2:
# 	curl -ks "https://127.0.0.1:8591/send?target=vasp3:8093" > /dev/null
#
# Replace the ?target parameter with any of the following:
#
#	vasp1:8591
#	vasp2:8592
#	vasp3:8593
#
.PHONY: demo-docker
demo-docker:
	hack/exec/demo-docker.sh

# Cleanup running VASP containers.
.PHONY: demo-docker-cleanup
demo-docker-cleanup:
	hack/exec/demo-docker-cleanup.sh

# Start local k8s cluster. To make use of (local) Kubernetes deployments the following
# binaries are required to be available: kubectl, skaffold and kind.
.PHONY: k8s-cluster-start
k8s-cluster-start:
	hack/exec/k8s-cluster-start.sh

# Destroy local k8s cluster.
.PHONY: k8s-cluster-destroy
k8s-cluster-destroy:
	hack/exec/k8s-cluster-destroy.sh

# Startup the VASPs in k8s.
.PHONY: k8s-vasps-run
k8s-vasps-run:
	hack/exec/k8s-vasps-run.sh

# Stop the VASPs in k8s.
.PHONY: k8s-vasps-delete
k8s-vasps-delete:
	hack/exec/k8s-vasps-delete.sh