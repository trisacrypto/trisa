.PHONY: all
all: build build-docker test

.PHONY: build
build:
	hack/exec/build.sh

.PHONY: build-docker
build-docker:
	hack/exec/build-docker.sh

.PHONY: bake-tooling
bake-tooling:
	hack/exec/bake-tooling.sh

.PHONY: test
test:
	hack/exec/test.sh

.PHONY: gazelle
gazelle:
	hack/exec/gazelle.sh

.PHONY: dependencies
dependencies:
	hack/exec/dependencies.sh

.PHONY: docs-dev
docs-dev:
	hack/exec/docs-dev.sh

.PHONY: docs-generate
docs-generate:
	hack/exec/docs-generate.sh

.PHONY: pki-dev-init
pki-dev-init:
	hack/exec/pki-dev-init.sh
	ls -l hack/etc/pki/dev/out

.PHONY: docker-run
docker-run:
	docker run -it --rm bazel/cmd/trisa:docker

.PHONY: docker-release
docker-release:
	skaffold build -p latest
