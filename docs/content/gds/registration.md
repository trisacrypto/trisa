---
title: "Registration"
date: 2021-07-22T09:21:59-04:00
lastmod: 2021-07-22T09:21:59-04:00
description: "Registering a VASP with the Directory Service"
weight: 20
---

To join the TRISA or TestNet networks, you must register with the TRISA Global Directory Service (GDS) or one of the jurisdiction-specific directory services. Registering with the directory service engages two workflows:

1. A KYV review process to ensure the network maintains trusted membership
2. Certificate issuance for mTLS authentication in the network

Coming soon: more details on the registration form, email verification, and the review process.

## Certificate Issuance

There are currently two mechanisms to receive mTLS certificates from the GDS when your registration has been reviewed and approved.

1. Emailed PKCS12 Encrypted Certificates
2. Certificate Signing Request (CSR)

You must select one of these options _when you submit_ your registration; after your registration is submitted you will not be able to switch between options.

### PKCS12 Encrypted Email Attachment

The first mechanism is the easiest &mdash; simply select the email option during registration and omit the CSR fields. If the registration form is valid, the GDS will return a PKCS12 password. **Do not lose this password, it is the only time it is made available during the certificate issuance process**.

Upon review approval, the GDS CA will generate a complete certificate including private keys and encrypt it using the PKCS12 password. After registering the public keys in the directory service, the GDS will then email the encrypted certificate as a ZIP file to the technical contact, or first available contact on the registration form.

After unzipping the email attachment, you should find a file named `<common_name>.p12`; you can decrypt this file to extract the certificates as follows:

```
$ openssl pkcs12 -in <common_name>.p12 -out <common_name>.pem -nodes
```

You can also directly use the .zip file without decrypting or extracting it via the [`github.com/trisacrypto/trisa/pkg/trust`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trust#NewSerializer) module.

### Certificate Signing Requests

An alternative to certificate creation is to upload a certificate signing request (CSR). This mechanism is often preferable because it means that no private key material has to be transmitted accross the network and the private key can remain on secure hardware.

To generate a CSR using `openssl` on the command line, first create a configuration file named `trisa.conf` in your current working directory, replacing `example.com` with the domain you plan to host your TRISA endpoint on:

```conf
[req]
distinguished_name = dn_req
req_extensions = v3ext_req
prompt = no
default_bits = 4096
[dn_req]
CN = example.com
O = [Organization]
L = [City]
ST = [State or Province (fully spelled out, no abbreviations)]
C = [2 digit country code]
[v3ext_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment, nonRepudiation
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = example.com
```

Please carefully fill out the configuration for your certificate, this information must be correct and cannot be changed without reissuing the certificate. Also make sure that there are no spaces after the entries in the configuration!

Then run the following command, replacing `example.com` with the domain name you will be using as your TRISA endpoint:

```
$ openssl req -new -newkey rsa:4096 -nodes -sha384 -config trisa.conf \
  -keyout example.com.key -out example.com.csr
```

Your private key is now stored in `example.com.key` &mdash; **keep this private key safe** &mdash; it is required for mTLS connections in your mTLS service and establishes trust on the TRISA network.

The `example.com.csr` file contains your certificate signing request. Copy and paste the contents of this file including the `-----BEGIN CERTIFICATE REQUEST-----` and `-----END CERTIFICATE REQUEST-----` into your registration request.