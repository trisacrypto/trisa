---
title: "Sectigo"
date: 2020-12-24T07:58:37-05:00
lastmod: 2020-12-24T07:58:37-05:00
description: "Directory Service interactions with the Sectigo CA API"
weight: 10
---

[![Go Reference](https://pkg.go.dev/badge/github.com/trisacrypto/testnet/pkg/sectigo.svg)](https://pkg.go.dev/github.com/trisacrypto/testnet/pkg/sectigo)


The TRISA Directory Service issues certificates using the Sectigo Certificate Authority via its IoT Portal. Because the directory service must collect public key material in order to facilitate an initial trusted handshake for mTLS, it uses the Sectigo IoT Manager API as part of the VASP registration and verification process. The `github.com/trisacrypto/testnet/pkg/sectigo` package is a Go library for interacting with the API, implementing the endpoints and methods required by the directory service. The TestNet also provides a command line utility for interacting with the API for administrative and debugging purposes. This documentation describes the command line utility, which also gives an overview of how to use the API directly to issue and revoke certificates.

Reference material:

- [Package Documentation](https://pkg.go.dev/github.com/trisacrypto/testnet/pkg/sectigo)
- [IoT Manager API Documentation](https://support.sectigo.com/Com_KnowledgeDetailPage?Id=kA01N000000bvCJ)
- [IoT Manager Portal](https://iot.sectigo.com)

## Getting Started

To install the `sectigo` CLI utility, either download a pre-compiled binary from the [releases on GitHub](https://github.com/trisacrypto/testnet/releases) or install locally using:

```
$ go get github.com/trisacrypto/testnet/cmd/sectigo
```

This will add the `sectigo` command to your `$PATH`.

## Authentication

The first step is authentication, you should set your username and password in the `$SECTIGO_USERNAME` and `$SECTIGO_PASSWORD` environment variables (alternatively you can pass them as parameters on the command line). To verify your authentication status you can use:

```
$ sectigo auth
```

The API authenticates by username and password then returns acess and refresh tokens which are stored in a local cache file. To see where your cache is stored:

```
$ sectigo auth --cache
```

If you'd like to check your credentials state, e.g. if the access tokens are valid, refreshable, or expired, use:

```
$ sectigo auth --debug
```

## Authorities and Profiles

To begin to interact with certificates you need to list the authorities and profiles that your user account has access to.

```
$ sectigo authorities
[
  {
    "id": 1,
    "ecosystemId": 100,
    "signerCertificateId": 0,
    "ecosystemName": "TRISA",
    "balance": 10,
    "enabled": true,
    "profileId": 42,
    "profileName": "TRISA Profile"
  }
]
```

The authority displays the methods and profiles that certificates are created under. Here the `profileId` field is very important for use in subsequent calls. You can also view how many licenses have been ordered/issued across all authorities as follows:

```
$ sectigo licenses
{
  "ordered": 2,
  "issued": 2
}
```

To get detail information for a profile, use the profile ID with the following command:

```
$ sectigo profiles -i 42
```

This will return the raw profile configuration. Before creating certificates with the authority, you'll need to know the required profile parameters:

```
$ sectigo profile -i 42 --params
```

## Creating Certificates

You can request a certificate to be created with the `commonName` and `pkcs12Password` params as follows (note for profiles that require other params, you'll have to use the code base directly and implement your own method):

```
$ sectigo create -a 42 -d example.com -p secrtpasswrd -b "example.com certs"
{
  "batchId": 24,
  "orderNumber": 1024,
  "creationDate": "2020-12-10T16:35:32.805+0000",
  "profile": "TRISA Profile",
  "size": 1,
  "status": "CREATED",
  "active": false,
  "batchName": "example.com certs",
  "rejectReason": "",
  "generatorParametersValues": null,
  "userId": 10,
  "downloadable": true,
  "rejectable": true
}
```

The `-a` flag specifies the authority, but should be a profile id. The domain must be a valid domain. If you don't specify a password, one is generated for you and printed on the CLI before exit. The `-b` flag gives a human readable name for the batch creation. The return data shows detail about the batch certificate job that was created; you can fetch the data to keep checking on the status as follows:

```
$ sectigo batches -i 24
```

You can also get processing information for the batch:

```
$ sectigo batches -i 24 --status
```

Once the batch is created, it's time to download the certificates in a ZIP file:

```
$ sectigo download -i 24 -o certs/
```

This will download the batch file (usually batchId.zip, 24.zip in this case) to the `certs/` directory. Unzip the certs then decrypt the .pem file as follows:

```
$ unzip certs/24.zip
$ openssl pkcs12 -in certs/example.com.p12 -out certs/example.com.pem -nodes
```

For more on working with the PKCS12 file, see [Export Certificates and Private Key from a PKCS#12 File with OpenSSL](https://www.ssl.com/how-to/export-certificates-private-key-from-pkcs12-file-with-openssl/).

## Managing Certificates

You can search for a certificate by name or serial number, but mostly commonly you search by the domain or common name to get the serial number:

```
$ sectigo find -n example.com
```

Once you've obtained the serial number you can revoke the certificate as follows:

```
$ sectigo revoke -p 42 -r "cessation of operation" -s 12345
```

This command expects the profile id that issued the certificate with the `-p` flag, an [RFC 5280 reason code](https://tools.ietf.org/html/rfc5280#section-5.3.1) passed via the `-r` flag (unspecified by default), and the serial number of the certificate using the `-s` flag. If this command doesn't error, then the certificate has been successfully revoked.

The RFC 5280 reasons are:

- "unspecified"
- "keycompromise"
- "ca compromise"
- "affiliation changed"
- "superseded"
- "cessation of operation"
- "certificate hold"
- "remove from crl"
- "privilege withdrawn"
- "aa compromise"

Note that the reason is whitespace and case insensitive.