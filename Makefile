.PHONY: all
all: build test

.PHONY: build
build:
	hack/exec/build.sh

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