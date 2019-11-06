---
title: "Development"
draft: false
weight: 20
---

## Build System

The TRISA server is currently written in Go. The code repository is setup with Bazel for the compilation
and dockerization. Installing Bazel is the primary requirement to get started with development. Bazel will
install all the necessary tooling and dependencies completely sandboxed on your system. There is no real
need to install Go, protobuf or any other requirements. As everything is sandboxed, Bazel will not disturb
any existing tooling on your local machine either.

### Bazel

On OSX Bazel can be installed as follows:

1. Make sure you don't have bazel intalled using core brew: `brew uninstall bazel`
2. Install bazel using `brew tap bazelbuild/tap` followed by `brew install bazelbuild/tap/bazel`

Verify your installation using `bazel --version`.

For other platforms, consult the [Bazel Installation Instructions](https://docs.bazel.build/versions/master/install.html).

### Additional Requirements

* Ensure `docker` and `docker-compose` are installed
* A regular build environment with `make` is advized as that will make it easier to consume the convenience targets we have setup.

## Building the Code

**NOTE: it can take a while the first time bazel runs as it needs to download and compile the dependencies**

The `Makefile` has some additional documentation for each available target. The primary targets for building are:

* `make build`
* `make test`