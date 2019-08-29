SHELL := /bin/bash

.PHONY: all
all: build test

.PHONY: build
build:
	hack/exec/build.sh

.PHONY: test
test:
	hack/exec/test.sh

.PHONY: gazelle
gazelle:
	hack/exec/gazelle.sh

.PHONY: dependencies
dependencies:
	hack/exec/dependencies.sh
