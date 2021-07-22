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

To generate a CSR run the following command:

```
$ openssl req -nodes -newkey rsa:4096 -keyout <common_name>.key -out <common_name>.csr
```

Please carefully fill out the prompts for your certificate, this information must be correct and cannot be changed without reissuing the certificate.
