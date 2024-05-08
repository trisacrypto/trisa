---
title: TRISA Envoy
date: 2021-04-23T01:35:35-04:00
lastmod: 2022-08-10T13:22:02-04:00
description: "Using the TRISA Envoy self-hosted node"
weight: 12
---

Built to help compliance teams handle travel rule exchanges efficiently and economically, Envoy is an open source, secure, peer-to-peer messaging platform. Designed by TRISA engineers and compliance experts, Envoy offers a mechanism for Travel Rule compliance by providing [IVMS101](https://www.intervasp.org/) data exchanges using the [TRISA](https://trisa.io) and [TRP](https://www.openvasp.org/) protocols.

It is important to understand what Envoy is, and also, what Envoy _is not_.

#### Key Benefits

- **Intuitive UI**: Send and received TRISA & TRP messages with a user-friendly interface.
- **Efficient Administration**: Simplify setup and maintenance with our admin tools.
- **Secure Decentralized Messaging**: Protect customer privacy and IVMS101 data with encrypted peer-to-peer communication.
- **Compliance-Friendly API**: A straightforward JSON REST API for seamless TRISA & TRP interactions.

#### Not Included

- Does not include _on chain analytics or automated KYC/AML checks_.
- Future plug-ins may be available for on-chain analytics or KYC checks.
- Strictly a **secure messaging service** designed to meet FATF requirements.
- Does not solve the **Sunrise** problem (nothing does yet).

## Implementation Options

You have three options to deploy or manage your Envoy node:

1. **DIY**: Envoy is open source ([MIT License](https://github.com/trisacrypto/envoy/blob/main/LICENSE)). You're free to download, install, integrate, host and support your own node. You're also free to fork and modify the node for your own use cases. Provided you have the technical capabilities and availability, you can have complete flexibility and control of your deployment!
2. **One-time Integration Service**: For a one-time fee, Rotational Labs will install and configure your Envoy node in your environment while you host, maintain, and support the node on an ongoing basis. Rotational will provide some training to get started with your node, and upgrade emails when it is time to update the Envoy version.
3. **Managed Service**: If you're looking to get something up and running now, Rotational Labs will install, configure, host, maintain, and support an Envoy node for you.

This documentation is focused on the DIY option. If you'd like more information on the one-time integration service or managed services, please [schedule a demo](https://rtnl.link/p2WzzmXDuSu) with Rotational Labs!

## Data Storage and Security

Because Envoy is intended to exchange and store compliance information that is by nature PII (Personally Identifiable Information), security and privacy is a top-of-mind concern. To that end, Envoy works to ensure data is stored in a secure and protected fashion.

### Data Storage

Nodes must store compliance data locally on disk in a backend store. Currently only sqlite3 is supported, but Postgres and other databases may be options in the future. When configuring an Envoy node, ensure:

1. Enough disk is provisioned for **long-term** travel rule data storage.
2. Ensure that the disk is independently **secured and encrypted**, particularly if you are hosting your Envoy service on a shared archticture or the cloud.
3. If you're using an external database, ensure it is not accessible from the public Internet.

### Security

Envoy employs TRISA cryptographic security standards for data in-flight and at-rest. Data exchanges are **secured by mTLS** when conducted over the TRISA network and if available for a TRP exchange. All TRP exchanges require valid TLS certificates to establish a secure connection.

Travel Rule PII is stored as secure envelopes with original TRISA cryptography. Even if the transfer occurs over TRP, a secure envelope is created to store the information on disk. Secure envelopes use multi-stage strong encryption, encrypting data symmetrically with AES-256, signing the data with HMAC-256, then encrypting the encryption keys and hmac secrets via asymmetric encryption using RSA-OAEP-256.

One you delete the private keys used to seal secure envelopes -- the data encrypted with those keys is effectively deleted.

Some indexing information is extracted for lookups and and search, but is not PII.

### Key Management

Public key management is required for effective use of Envoy. Envoy can store key material locally on the same disk as its database, but this is not recommended. Instead we recommend the use of a key store or vault such as Google Secret Manager or Hashicorp Vault.

### Authentication and Access

We strongly recommend that the internal UI and API is restricted to specific IP addresses or accessed from within a VPN. This will prevent brute force attacks on your Envoy node.

Passwords and API Key secrets are stored on the Envoy disk as [Argon2](https://github.com/P-H-C/phc-winner-argon2) hashes that cannot be reversed. Argon2 is a password-hashing function that includes a time and memory cost to balance defense against cracking attacks. Even using Argon2, it is still recommended that passwords are routinely changed and that strong passwords are generated. API keys should also be revoked if not in use, or recycled on a routine basis.