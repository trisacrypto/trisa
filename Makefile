.PHONY: all
all: build build-docker test

# Build TRISA server.
.PHONY: build
build:
	hack/exec/build.sh

# Build docker container for TRISA server. The resulting images will be pushed
# to the local docker instance as "bazel/cmd/trisa:docker".
.PHONY: build-docker
build-docker:
	hack/exec/build-docker.sh

# Bake dokerized tooling (cfssl, gohugo).
.PHONY: bake-tooling
bake-tooling:
	hack/exec/bake-tooling.sh

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