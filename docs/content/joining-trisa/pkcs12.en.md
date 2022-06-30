---
title: Anatomy of a PKCS12 file 
date: 2022-06-17T09:15:46-04:00
lastmod: 2022-06-17T16:53:39-04:00
description: "Anatomy of a PKCS12 file"
weight: 50
---

When your registration has been reviewed and approved, you will receive mTLS certificates from the Global Directory Service (GDS). These will be sent via an email containing PKCS12 password protected certificates.

What to expect after registration and during certificate issuance:
- If the registration form is valid: GDS will return a PKCS12 password. This is the password you will use to decrypt your certificates.
{{% notice note %}}
This is the only time the PKCS12 password is made available during the certificate issuance process.
{{% /notice %}}
- Once approved, TRISA Certificate Authority (CA) will generate a complete certificate including private keys and encrypt it using the PKCS12 password.
- After registration of the public keys in the directory service, GDS will then email the encrypted certificate as a ZIP file to the technical contact, or first available contact on the registration form. The content of that ZIP file is the PKCS12 file.

## Anatomy of the PKCS12 file

A PKCS12 file is a binary format for storing a certificate chain and a private key in a single encryptable file. The PKCS12 is one of the family of standards called Public-Key Cryptography Standards (PKCS). The PKCS12 contains the X.509 digital certificate that carries a copy of the public key of the subject (i.e. the VASP or entity being issued the certificate), as well as other standard fields. 

These certificates issued by the TRISA CA are used to identify and verify originating VASPs and beneficiary VASPs. The PKCS12 may include 2 certificates: the client's certificate and one or more issuer (CAs) certificates. In this case, it will include the VASP's certificate and the TRISA CA certificate. 

The PKCS12 file below contains the X.509 Certificate chain between the client (the VASP) and the issuer (TRISA CA).

![PKCS12 file](/img/pkcs12.png)

The client certificate also includes additional information about the VASP: the Bag Attribute and LocalKeyId (which is associated with an X.509 Public Key Certificate and its matching asymmetric private key in this PKCS12 file), and the subject and issuer fields. The subject and issuer fields are information from the certificate itself and is extracted by `openssl` and may include Country (C), State (S), Locality (L), Street Address (street), Postal Code (postal_code), O (Organization), and Common Name (CN) of both subject and issuer.

Similarly, the TRISA CA certificate also includes the subject and issuer fields, for the TRISA CA, the entities are the same.
 
The Private Key also contains the Bag Attribute and LocalKeyId. 

## Accessing and Saving the Certificates

The ZIP file contains a PKCS12 file (e.g. `<common_name>.p12`) and can be unzipped as follows:

```
$ unzip <common_name>.zip
```

The PKCS12 file can be decrypted and viewed on screen in PEM format using `openssl`. The `-info` flag provides information about the PKCS12 structure, the `-in` flag specifies the input filename, and the `-nodes` flag ensures private keys are not encrypted. You will then be prompted to enter the PKCS12 password provided during registration.

```
$ openssl pkcs12 -info -in <common_name>.p12 -nodes 
``` 

You can save the contents of the PKCS12 file (e.g. a `.pem` file), using a similar command mentioned above by using the `-out` flag to specify an output file. Where necessary, excluding the `-nodes` flag will encrypt the output file.

```
$ openssl pkcs12 -info -in <common_name>.p12 -out <common_name>.pem -nodes
``` 

You can also save the certificates or the private keys separately, using the `-nocerts` flag or the `-nokeys` flag, respectively.

```
$ openssl pkcs12 -info -in <common_name>.p12 -out <common_name>.key -nodes -nocerts
$ openssl pkcs12 -info -in <common_name>.p12 -out <common_name>.crt -nodes -nokeys
```

Or only save the client certificates or CA certificates, using the `-clcerts` flag or the `-cacerts` flag, respectively.

```
$ openssl pkcs12 -info -in <common_name>.p12 -nodes -clcerts
$ openssl pkcs12 -info -in <common_name>.p12 -nodes -cacerts
```

**Further Reading**
- [OpenSSL documentation](https://www.openssl.org/docs/man1.1.1/man1/openssl-pkcs12.html)
- [TRISA White Paper](https://trisa.io/trisa-whitepaper/)
