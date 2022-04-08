---
title: TRISA CLI
date: 2022-04-02T12:09:09-05:00
lastmod: 2022-04-02T12:09:09-05:00
description: "Using the TRISA command line interface for development"
weight: 15
---

The TRISA command line client is a utility that assists TRISA integrators and developers testing their TRISA service. Advanced users may also use the TRISA client to execute TRISA requests for compliance purposes, although this is not recommended for extensive use. To install the latest version of the TRISA CLI, you must have [Go installed](https://go.dev/doc/install) on your computer. The following command will install the latest version of the CLI:

```
$ go install github.com/trisacrypto/trisa/cmd/trisa@main
```

{{% notice note %}}
We are currently working on a release mechanism that will automatically build the CLI for a variety of platforms so that in the future you will not need to have Go installed.
Stay tuned!
{{% /notice %}}

## Configuration

Before you can start using the TRISA CLI, you must first configure your environment to ensure that you can successfully connect to a remote peer or the directory service.

**Prerequisites**:

1. The `trisa` command installed and on your `$PATH`
2. Your [testnet certificates]({{< ref "/gds/registration" >}}) that include both the trust chain and private key.

The TRISA CLI command is configured via flags specified for each command or by setting environment variables in your shell with the configuration. The CLI also supports the use of [.env](https://platform.sh/blog/2021/we-need-to-talk-about-the-env/) files in the current working directory for configuration. To see what CLI flags should be specified use `trisa --help`. An example `.env` configuration file is as follows:

```ini
# The endpoint to the TRISA node that you'd like to connect to. The endpoint can be
# found using the directory service lookup command.
TRISA_ENDPOINT=example.com:443

# Directory service you'd like to connect to. You can specify a short name such as
# "testnet" or "mainnet" or the endpoint of the directory service to connect to. The
# configured directory is trisatest.net by default.
TRISA_DIRECTORY=testnet

# Path to your TRISA identity certificates that include the private key. This can be the
# original .zip file sent by Sectigo or the unzipped .p12 file; in which case the
# PKCS12 password must also be supplied. If you've decrypted it manually it should be in
# PEM encoded format with the .pem or .crt extension.
TRISA_CERTS=path/to/certs.pem

# If you've split your certs into the public trust chain without private keys and a
# private key file, then specify the path to the trust chain (optional).
TRISA_TRUST_CHAIN=path/to/chain.pem

# If the certs are PKCS12 encrypted then specify the password for decryption (optional).
TRISA_CERTS_PASSWORD=supersecret
```

The simplest way to get started with TRISA is to copy and paste the above snippet into a `.env` file in your current directory, then modifying the values as necessary.

## Creating Secure Envelopes

## Interacting with TRISA Peers

The primary use of the `trisa` CLI is to execute TRISA RPC requests to a TRISA node.

### Transfers

Send a secure envelope to the TRISA node.

{{% notice note %}}
The TRISA CLI command currently does not implement the `TRISANetwork/TransferStream` birdirectional streaming RPC and does not have plans to implement this in the CLI. If you would like an implementation of streaming from the command line, please open an issue on our [GitHub repository](https://github.com/trisacrypto/trisa/issues).
{{% /notice %}}

### Key Exchanges

Send a key exchange request to get the public sealing key of the node. By default, this command stores the sealing key in the configuration directory and uses it to seal secure envelopes when making transfers.

### Address Confirmation

{{% notice warning %}}
Address confirmation has not yet been fully defined by the TRISA working group. The TRISA technical subcommittee is currently working on the Address Confirmation protocol, so stay tuned for more information!
{{% /notice %}}

### Health Checks

Send a health check request to check the status of the TRISA node.

## Interacting with the GDS

The TRISA CLI supports some basic interactions with the TRISA Global Directory Service (GDS) based on your initial configuration.

### Lookup

Lookup a TRISA VASP by common name or VASP ID.

### Search

Search for a TRISA VASP by name or by website.

### List

Return a list of currently verified VASPs (requires mTLS certificates).