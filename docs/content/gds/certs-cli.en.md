---
title: "Certificates CLI"
date: 2022-06-28T15:45:48-04:00
lastmod: 2022-06-28T19:34:20-04:00
description: "Using the Certificates Tool and Local Disk Certificate Management and Test CA"
weight: 80
---

The Certificates command line client under the Directory Service is a utility designed to create local disk certificate management and test Certificate Authority (CA) for testing purposes. The code is located under `cmd/certs`. To install the latest version of the Certificates CLI, you must have [Go installed](https://go.dev/doc/install) on your computer. The following command will install the latest version of the CLI:

```
$ go install github.com/trisacrypto/directory/cmd/certs@latest
```

Alternatively, to install the Certificate CLI from the Directory Project root, run the following:

```
$ go install ./cmd/certs
```

## Configuration

Before using the `certs` CLI, you must confirm your environment has the `certs` program in your go workspace `bin` and is present in your `$PATH`.

```
$ ls $GOPATH/bin
certs                                            
```

The `certs` CLI has two command categories: for the CA and for the client certs. To see the CLI commands and which flags should be specified, use `certs --help`. 

## Creating a Test Certificate Authority (CA) Certs 

The first step to creating the local disk certificate management is creating the CA certs that contain a randomly generated certificate and RSA key pair. You can use `certs init` to create CA certs if they do not already exist. The `-c` flag specifies the output location to write the CA certs; without the flag, the default location of the CA certs will be in `fixtures/certs/ca.gz`. If you want to overwrite existing CA keys, you can use the `-f` flag.

```
$ certs init -c path/to/ca.gz
```

## Generating the Certificates

After creating the local CA, it then can issue "self-signed" certificates by that local CA, e.g., certificates signed by the keys in `ca.gz`, as follows:

```
$ certs issue -c ca.gz -n "Common_Name" -O "Organization" 
```

The `issue` command requires the `-n` flag (which specifies the common name for the certificate) and the `-O` flag (which specifies the name of the organization that the certificates are for.) The default output location is `<common_name.gz>`. However, there are other flags to specify additional information for the organization.

```
-n value, --name value          common name for the certificate
-O value, --organization value  name of organization to issue certificates for
-C value, --country value       country of the organization
-p value, --province value      province or state of the organization
-l value, --locality value      locality or city of the organization
-a value, --address value       street address of the organization
-z value, --postcode value      postal code of the organization       
```

You can also [PKCS12]({{%relref "joining-trisa/pkcs12" %}}) encrypt the issued cert with the `-P` flag to specify the password and the `-o` flag to specify the output location using the `.p12` extension.

```
$ certs issue -c ca.gz -o certs.p12 -P password -n "Common_Name" -O "Organization" 
```

## Decrypting Files

The `certs` CLI also has a command to decrypt PKCS12 files and can be used as follows:
```
$ certs decrypt <common_name.p12> -p password
```

The `decrypt` command requires the `-p` flag (which specifies the password of the encrypted file). And default output is `<common_name.gz>`, but you can specify the output location using the `-o` flag.

## Creating Trust Chain Pools

Additionally, the `certs` CLI has a command to create a trust chain pool from certificates that were created on your local disk. It saves the pool into a zip file that you can specify with the `-o` flag. The pools are used to manage the public trust certificates of peers in the network. The pool maps common names to the provider certificates and ensures that only public providers without private keys are stored in the pool. 

```
$ certs pool <common_name1.pem> <common_name2.pem> -o pool.zip
```
While public providers are used in providers pools to facilitate mTLS clients, providers with keys (private providers) are used to instantiate mTLS servers. Such that, if two pairs of certs signed with the same `ca.gz` will enable you to establish an mTLS connection between two nodes.
