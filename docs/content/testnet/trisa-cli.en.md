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

Before you can start using the TRISA CLI, you must first install a profile that specifies the configuration details used by the command.

**Prerequisites**:

1. The `trisa` command installed and on your `$PATH`
2. Your [testnet certificates]({{< ref "/gds/registration" >}}) that include both the trust chain and private key.

The TRISA CLI command is configured via a collection of _profiles_. A profile manages how the command connects to the remote TRISA peer including what certificates are used for mTLS authentication, what endpoint is used to connect, and how to store and manage public sealing keys exchanged with the remote peer. The TRISA CLI may manage multiple profiles if you're connecting to multiple TRISA nodes, but it requires at least one profile before the CLI can be used.

### Installation

To create your first profile, execute the install command as follows:

```
$ trisa install
```

This command is essentially an alias for `trisa profile --create default`. It will take you through an interactive prompt to configure your first profile.

### Managing Profiles

The `trisa profiles` command allows you to manage the profiles you've created or installed. To view the configuration of the currently active profile, simply run:

```
$ trisa profile
```

To list the available profiles that have been created:

```
$ trisa profiles --list
```

{{% notice tip %}}
The `trisa profile` and `trisa profiles` commands are aliases, you may use both interchangely.
{{% /notice %}}

To activate a different profile, use the activate command with the name of the profile you wish to use as follows:

```
$ trisa profile --activate [name]
```

If you wish to create a new profile with a short name, use the create command:

```
$ trisa profile --create [name]
```

This will run the interactive create profile script that you used when you installed your first profile.

### Configuration Directory

All profiles, certificates, and sealing keys are stored in an operating-specific configuration directory (e.g. `${HOME}/.config/trisa` for Linux/BSD operating systems). To view the configuration directory run:

```
$ trisa profile --path
```

The `trisa` CLI is configured by a single YAML file called `config.yaml` that holds all profiles, but expects that all related files (such as certificates or sealing keys) are in the same directory as the configuration file.

If you would prefer to use a different directory for configuration, specify the `$TRISA_CONF_DIR` environment variable. If you'd prefer to specify a different configuration file, use the `$TRISA_CONF` environment variable.

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