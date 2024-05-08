---
title: Deploying Envoy
date: 2024-05-08T13:56:37-05:00
lastmod: 2024-05-08T13:56:37-05:00
description: "A quick guide on deploying your Envoy node"
weight: 60
---

This guide assumes that you're ready to deploy your Envoy node and that you've already obtained either TRISA TestNet or MainNet certificates as described by [the Joining TRISA guide]({{< ref "joining-trisa" >}}). If you haven't already, please go to the [TRISA Global Directory Service (vaspdirectory.net)](https://vaspdirectory.net/) to register for your certificates!

{{% notice style="note" title="Local Development" icon="code" %}}
If you'd like information about how to run Envoy locally using [Docker Compose](https://docs.docker.com/compose/) and self-signed keys generated using `openssl` please go to the repository at [trisacrypto/envoy](https://github.com/trisacrypto/envoy) and follow the instructions in the `README.md`.
{{% /notice %}}

The general/top-level steps to deploy an Envoy node are as follows:

1. Obtain and decrypt TRISA certificates
2. Setup a deployment environment (e.g. a cloud instance or kubernetes cluster)
3. [Configure]({{< relref "configuration.md" >}}) the Envoy node via the environment
4. Deploy your Envoy node using one of the instructions below
5. Ensure that you can reach your node at port 443
6. Configure DNS to point your TRISA endpoint at your node
7. Create users/api keys to access the internal UI/API

## Deploying Envoy

There are many ways to deploy Envoy and a lot depends on your internal infrastructure or cloud service provider. This guide provides examples for deploying Envoy using a Kubernetes cluster (the default way that we deploy our services) or using `systemd` on Ubuntu to run your Envoy service on a cloud instance.

### Using Kubernetes

Coming Soon!

### Using systemd

Coming Soon!

### Compiling and Installing Envoy

Envoy is written in the [Go programming language](https://pkg.go.dev/github.com/trisacrypto/envoy) and so can be compiled using the [go toolchain](https://go.dev/doc/tutorial/compile-install). If you have `go` installed on your computer, it may be possible for you to simply run:

```
$ go install github.com/trisacrypto/envoy/cmd/envoy@latest
```

To compile and install the latest version of `envoy` on your `$PATH`. Or if you prefer to install a specific version of Envoy:

```
$ go install github.com/trisacrypto/envoy/cmd/envoy@v0.11.0
```

The complicating factor is that `CGO` is required to compile Envoy, which means you may have to set the `CGO_ENABLED=1` environment variable. Additionally you'll have to have either `clang` or `gcc` installed to compile the dependent packages for your architecture if they cannot be installed using go modules.

### Building the Docker Image

You can build the Envoy Docker image using the `Dockerfile` in the root of the [trisacrypto/envoy](https://github.com/trisacrypto/envoy) repository. After cloning the repository and changing into its root directory, run:

```
$ docker build -t [TAG] .
```

If you'd like to build the image for a different platform, you can use `docker buildx` as follows:

```
$ docker buildx build -t [TAG] --platform linux/amd64,linux/arm64 .
```

Feel free to push the resulting image to your container registry of choice; or just use the default Docker images on DockerHub!

## Accessing Envoy

Although you can create users and API keys using the user interface there is a bit of a chicken and egg problem: how do you create a user to log in to the user interface to create a user? Additionally, if you've disabled the UI, how do you create API keys for the internal API? Luckily the `envoy` command has some helper tools to do this on your behalf.

### Running the Envoy Command

To create users and API keys, wherever you run the `envoy` command will need to have the correct configuration to reach the database that backs Envoy. In practice, this typically means that you need to run the `envoy` command on the instance where you deployed it. For example, if you've deployed using Kubernetes you could run:

```
$ kubectl -n [NAMESPACE] exec -it [POD] -- envoy -h
```

Alternatively you will have to SSH into the instance you're running, and ensure that `envoy` is on your `$PATH` and that you have the correct permissions to execute it.

### Creating Users

Use the following command to create a user:

```
$ envoy createuser -n [NAME] -e [EMAIL] -r [ROLE]
```

This command will create the specified user and will print out the password that you can use to login to the user interface with.

Role can be one of:

- `admin`: has full access to the UI
- `compliance`: can perform compliance related actions but not create users/apikeys
- `observer`: only has read-only access to the Envoy node

### Creating API Keys

Use the following command to create an API key that has all permissions:

```
$ envoy createapikey all
```

Alternative, you can specify which permissions you want the API key to have by listing them each in a space delimited form:

```
$ envoy createapikey users:manage users:view
```

The list of the permissions you can add to an API key can be found in the [API guide permissions table]({{< relref "api.en.md#permissions" >}}).

