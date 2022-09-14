---
title: "Registration"
date: 2021-07-22T09:21:59-04:00
lastmod: 2021-07-22T09:21:59-04:00
description: "Registering a VASP with the Directory Service"
weight: 20
---

Before you can integrate the TRISA protocol into your VASP software, you must [register](https://vaspdirectory.net/guide) with the TRISA Global Directory Service (GDS).

The TRISA Global Directory Service (GDS) provides public key and TRISA remote peer connection information for registered VASPs. For more detailed information about the directory, see the documentation on the [GDS]({{< ref "/gds" >}}).

Once you have registered and been verified, you will receive Identity Certificates. The public key in these certificates will be made available to other VASPs via the GDS.

When registering with the GDS, you will need to provide the `address:port` endpoint where your VASP implements the TRISA Network service. This address will be registered with the GDS and utilized by other VASPs when your VASP is identified as the beneficiary VASP.

For integration purposes, when you [register](https://vaspdirectory.net/guide) with the GDS, you can opt for either MainNet or TestNet Certificates, or both. The TestNet instance is designed for [testing]({{< ref "/testing" >}}), and the registration process is streamlined in the TestNet to facilitate quick integration. The MainNet is design for production Travel Rule implementations. It is recommended to register for both MainNet and TestNet, specifying different endpoints to reduce confusion for your VASP counterparties.

### Directory Service Registration

To start your registration, visit [https://vaspdirectory.net/](https://vaspdirectory.net/guide). You will first need to create an account, and then log in using that account to start the registration process. Note that you can use this website to enter your registration details on a field-by-field basis, or to upload a JSON document containing your registration details.

One of the key pieces of information you'll need is your TRIXO Form. Below is an excerpt of some of the key fields in the TRIXO form, which provides information about transaction thresholds, currency types, and applicable regulators. Frequently, several people at an organization (e.g. legal, technical, administrative points-of-contact) need to collaborate to complete the needed information. To see the TRIXO form in full, see the [TRIXO documentation]({{< ref "/joining-trisa/trixo" >}}).

```json

 "trixo": {
    "primary_national_jurisdiction": "USA",
    "primary_regulator": "FinCEN",
    "other_jurisdictions": [],
    "financial_transfers_permitted": "no",
    "has_required_regulatory_program": "yes",
    "conducts_customer_kyc": true,
    "kyc_threshold": "1.00",
    "kyc_threshold_currency": "USD",
    "must_comply_travel_rule": true,
    "applicable_regulations": [
      "FATF Recommendation 16"
    ],
    "compliance_threshold": "3000.00",
    "compliance_threshold_currency": "USD",
    "must_safeguard_pii": true,
    "safeguards_pii": true
  }
```

The final step of registration will be a [pkcs12 password]({{< ref "/joining-trisa/pkcs12" >}}), which you must keep to decrypt the Identity Certificates that will be sent via email.

This registration will result in an email being sent to all the technical contacts specified through the webform or in the JSON file. The emails will guide you through the remainder of the registration process. Once you’ve completed the registration steps, TRISA administrators will receive your registration for review.

Once the administrators have reviewed and approved the registration, you will receive [pkcs12 password]({{< ref "/joining-trisa/pkcs12" >}})-protected, compressed Identity Certificate via email and your VASP will be publicly visible in the GDS.

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