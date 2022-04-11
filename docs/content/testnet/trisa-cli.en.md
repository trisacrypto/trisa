---
title: TRISA CLI
date: 2022-04-02T12:09:09-05:00
lastmod: 2022-04-10T09:32:16-05:00
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

The first step when using the TRISA CLI is to create some payload data that can be sealed inside of secure envelopes for TRISA envelopes. At a minimum, there are two JSON files that you need to create or provide for the payload:

1. An _identity_ payload containing [IVMS 101](https://github.com/trisacrypto/trisa/blob/347f88d55df4d4e0167ad4e005721b638991ecef/proto/ivms101/identity.proto#L46-L53) data.
2. A _transaction_ payload containing a [Transaction](https://github.com/trisacrypto/trisa/blob/347f88d55df4d4e0167ad4e005721b638991ecef/proto/trisa/data/generic/v1beta1/transaction.proto#L11) or [Pending](https://github.com/trisacrypto/trisa/blob/347f88d55df4d4e0167ad4e005721b638991ecef/proto/trisa/data/generic/v1beta1/transaction.proto#L30) message.

The identity payload is the compliance information required for the transfer and the transaction payload is used to identify the transaction on the chain and associate it with the identity information. For ease of data entry, these files may be specified as JSON files and the protocol buffer payload created using the `trisa make` command.

```
$ trisa make -i identity.json -t transaction.json -o envelope.json
```

With no other arguments, this command creates an unsealed envelope that does not have an envelope ID nor sent at and received at timestamps in the payload. The documentation refers to this kind of secure envelope as a "payload template" in the rest of the documentation because it can be loaded by the `trisa seal` or `trisa transfer` commands to update the envelope ID, timestamps, before sealing the envelope with the public keys of the recipient.

To create a complete envelope or a fully sealed envelope


## Sealing

## Opening

By default the `trisa open` command is used to unseal an envelope and save it as an unsealed envelope for further processing. This command can also be used to extract the payload or check if an incoming envelope has an error on it. To extract an unsealed envelope and save it to disk:

```
$ trisa open -in envelope.json -out unsealed_envelope.json -key private.pem
```

If you add the `-payload` flag, then the payload will be decrypted and saved to disk; adding the `-error` flag will extract an error and save it to disk. If the `envelope.json` is an unsealed envelope, then the `-key` flag can be omitted. If the `-out` flag is ommitted, the contents will be printed to disk. For example to simply view the payload of a sealed envelope:

```
$ trisa open -in envelope.json -key private.pem -payload
```

Or to view an error on the envelope:

```
$ trisa open -in envelope.json -error
```

Note that no private key is required for errors since errors are not encrypted.

## Interacting with TRISA Peers

The primary use of the `trisa` CLI is to execute TRISA RPC requests to a TRISA node. A general workflow is as follows:

1. Identify the peer endpoint using the Directory Service lookup or search functionality.
2. Create a secure envelope or payload template to prepare to send to the remote peer.
3. Perform a key exchange with the remote peer and save the sealing keys.
4. Seal the secure envelope or payload template with the remote peer's sealing keys.
5. Execute a transfer and save the response envelope.

This workflow generally mirrors the workflow of live TRISA compliance operations, though many of the steps are manual to facilitate integration and development.

### Transfers

Send a secure envelope to the remote TRISA peer and receive a secure envelope in exchange. Transfers are the central compliance exchange mechanism in the TRISA protocol. If you have already created and sealed an envelope, saving it to `outgoing.json` you can transfer it as follows:

```
$ trisa transfer -i outgoing.json -o response.json
```

This will execute the TRISA transfer and save the response, including TRISA error envelopes, to disk at the specified path. If the extension of the output path is `.json` then the envelope is marshaled to `.json` format, if it is the `.pb` extension it will be saved as a raw protocol buffer. If the `-o` flag is not supplied, then the JSON response will be printed to the command line. If you would like the decrypted payload printed, then you must provide the private sealing key:

```
$ trisa transfer -i outgoing.json -k private.pem
```

If both an output path and the private key are provided then a JSON file is produced with the unsealed envelope that can be read using the `open` command or resent using the `seal` command.

You can also use a secure envelope payload template to seal and transfer an envelope in one step instead of using the intermediate `seal` command:

```
$ trisa transfer -i outgoing.json -s public_sealing_key.pem
```

See [sealing secure envelopes]({{< relref "#sealing" >}}) for more information on the command line arguments that can be used to adapt secure envelopes before sending them.

If you would like to send an error-only secure envelope to the recipient, then you must supply the envelope ID, error code, and error message as follows:

```
$ trisa transfer -I envelope-id-foo -C COMPLIANCE_CHECK_FAIL -E "something went wrong"
```

Note that sending an error-only secure envelope is usually a response to an incoming message. This mechanism is used primarily to test a server's handling of an asynchronous transfer workflow.

{{% notice note %}}
The TRISA CLI command currently does not implement the `TRISANetwork/TransferStream` birdirectional streaming RPC and does not have plans to implement this in the CLI. If you would like an implementation of streaming from the command line, please open an issue on our [GitHub repository](https://github.com/trisacrypto/trisa/issues).
{{% /notice %}}

### Key Exchanges

Send a key exchange request to get the public sealing key of the node. Key management is a somewhat complex topic, and the TRISA CLI attempts to do the simplest possible thing to enable testing and development. A key exchange requires you to send your public sealing keys to the remote node, and they will return keys to you. Prior to a transfer, a key exchange must be completed so that you have the sealing keys to create a secure envelope and so that the remote has your public keys to send a response.

By default, the TRISA CLI will simply use your TRISA mTLS identity certificates as the keys for a key exchange. The simplest exchange is therefore:

```
$ trisa exchange -o peer_sealing_keys.pem
```

The `-o` flag saves the keys to disk at the specified path, so that you can use the keys later on to make secure envelopes or conduct transfers. If the `-o` flag is ommitted, the JSON data of the response, an `SigningKey` protocol buffer message will be printed to `stdout`. There are several formats that the keys can be saved in: a path with a `.json` or `.pb` extension will save the protocol buffer message to disk in the specified format; a path with a `.pem` or `.crt` extension will save the keys as PEM encoded public key.

To send alternative keys to the remote peer in a key exchange, you may use the `-i` flag to specify the path of keys to send. If the input path ends in `.json` or `.pb` it is parsed as a `SigningKey` protocol buffer message in JSON format or raw protobuf format respectively. If the input path ends in `.pem` or `.crt` it is parsed as a PEM encoded public key or x.509 certificate. Note that the PEM encoded format, the first `PUBLIC KEY` or `CERTIFICATE` block that is found is used for parsing.

To generate your own RSA keys to send to the remote server for key exchange, use the following commands:

```
$ openssl genrsa -out private.pem 4096
$ openssl rsa -in private.pem -pubout -out public.pem
```

This will create two files, `private.pem` that contains your private keys and `public.pem` which contains your public keys. Send the public key to the remote TRISA peer as follows:

```
$ trisa exchange -i public.pem -o peer_sealing_keys.pem
```

Ensure that you keep the `private.pem` file so that you can decrypt transfers that follow; it is likely that the remote TRISA node will use the key you just exchanged in sending outgoing transfers and preparing responses to your transfers. The only way to decrypt that data is with the private key!

### Address Confirmation

{{% notice warning %}}
Address confirmation has not yet been fully defined by the TRISA working group. The TRISA technical subcommittee is currently working on the Address Confirmation protocol, so stay tuned for more information!
{{% /notice %}}

### Health Checks

Send a health check request to check the status of the TRISA node.

The health check RPC is primariliy for the directory service to monitor if the TRISA network is online, however it is also useful for debugging or diagnosing connection issues (e.g. is the remote peer offline or are my certificates invalid). A simple health check request is as follows:

```
$ trisa status
```

Use the `--insecure` flag to connect without mTLS credentials, though not all TRISA peers will support an insecure status check.

## Interacting with the GDS

The TRISA CLI supports some basic interactions with the TRISA Global Directory Service (GDS) based on your initial configuration.

### Lookup

Lookup a TRISA VASP by common name or VASP ID. The following requests will lookup Alice on the TestNet by both common name and ID:

```
$ trisa lookup -n api.alice.vaspbot.net
```

```
$ trisa lookup -i 7a96ca2c-2818-4106-932e-1bcfd743b04c
```

Lookups also support the registered directory argument if looking up a VASP that is a member of the network but was issued certificates from a different directory service. If omitted, by default the directory service will lookup the VASP record as though it was the registered directory.

### Search

Search for a TRISA VASP by name or by website. You can specify multiple names and websites to expand your search. E.g. to search for "Alice" and "Bob" on the TestNet:

```
$ trisa search -n alice -n bob
```

If this returns too many results you may specify either category or country filters for the results. Country filters are inclusive and should be ISO Alpha-2 country codes:

```
$ trisa search -n alice -n bob -c US -c SG
```

{{% notice tip %}}
Categories are case sensitive and websites must be full URLs for the search to work. If you're not getting any results for a website search, try adding the `http://` prefix or removing any paths from the URL. If you're not getting any results for the name, try using a prefix of the name that is greater than 3 characters.
{{% /notice %}}

Categories that may be helpful in filtering:

| VASP Categories | Business Categories   |
|-----------------|-----------------------|
| Exchange        | PRIVATE_ORGANIZATION  |
| DEX             | GOVERNMENT_ENTITY     |
| P2P             | BUSINESS_ENTITY       |
| Kiosk           | NON_COMMERCIAL_ENTITY |
| Custodian       |                       |
| OTC             |                       |
| Fund            |                       |
| Project         |                       |
| Gambling        |                       |
| Miner           |                       |
| Mixer           |                       |
| Individual      |                       |
| Other           |                       |