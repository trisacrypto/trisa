---
title: Domain Verification
date: 2023-03-28T14:30:51-05:00
lastmod: 2023-03-28T14:30:51-05:00
description: "Domain Verification"
weight: 70
---

The [TRISA Certificate Authority]({{% relref "/ca" %}}) issues [x509 certificates](https://sectigo.com/resource-library/what-is-x509-certificate) for mTLS authentication in its peer-to-peer network. These certificates are dependent on the verified ownership of a domain name that is used both in the _common name_ of the certificate (CN) as well as the _subject alternative name_ (SAN) extension of the certificate. Your TRISA node must be hosted at a domain name that is in the SAN field of the certificate otherwise mTLS connections will fail.

{{% notice note %}}
The domain that you host your TRISA node, e.g. `trisa.example.com` must be the common name of your TRISA Identity Certificates and _must_ match the endpoint of your TRISA directory record, e.g. `trisa.example.com:443`. If not, TRISA peers will be unable to connect to your TRISA node using mTLS. If you are using multiple domain names, please contact TRISA support for assistance.
{{% /notice %}}

Because the domain names assigned to the certificate are part of the secure handshake in mTLS; you must prove that you own the specified domain name before certificates may be issued. This prevents bad actors from attempting to pose as another entity in the TRISA network. There are two methods of domain verification: email verification (the default) and DNS verification.

## Email Verification

By default, when you register for a TRISA account a secure token is sent to your email address with a link to validate the email. When you click on the link the token is verified in the directory service; proving that you do have control over the email account specified in the registration. Multiple verification emails may be sent - please click on them all!

If the domain of your email address matches the common name requested for the TRISA certificates, then this validation process also proves that you have control of the domain.

If your email address domain does not match the domain of the common name; then you may contact TRISA support in order to change your email address to one that does. A verification email will be resent and you can complete the domain verification with the new email address.

## DNS Verification

If you do not have an email address with a domain that matches the common name you have requested for the TRISA certificates; then please get in touch with TRISA support to perform manual DNS verification. Note that this process can be challenging and you will need access to your registrar or DNS host and the permissions to add DNS records for the specified domain.

After contacting support, you will receive a [whisper](https://whisper.rotational.dev/) link that contains a verification code that will likely look something like:

```txt
TRISA-DOMAIN-VERIFICATION-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

Once you've obtained the verification code, perform the following steps:

1. Log into your registrar or your DNS management utility
2. Edit the DNS records for the domain you wish to verify
3. Create a new TXT record for the TRISA domain whose value is the verification code
4. Contact TRISA support and let them know they can verify the DNS record

Consider the example where you wish to create a verification record for `trisa.example.com`. In this example, `example.com` is the root domain whose DNS records you need to edit. Create a `TXT` record for the name `trisa.example.com` (note there will likely be a field called **Host**, **Hostname**, or **Alias** where you will enter `trisa`). There are many types of records that you can create, e.g. `A` records or `CNAME` records, but it is important that you create a `TXT` record for verification. Make sure you copy the complete verification code without any spaces before or after the code into the value field (which may be called **Data**, **Answer**, or **Destination**).

Note that this process will soon be automated in the future to make it easier for TRISA members to verify the domains that they control.