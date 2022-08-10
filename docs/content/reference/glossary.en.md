---
title: Glossary
date: 2021-06-14T15:59:09-05:00
lastmod: 2022-08-10T16:33:59-04:00
description: "TRISA Glossary and Terminology"
weight: 30
---

## General Terminology

- **Originator**: The initiator of a blockchain transaction and therefore also the initiator of a TRISA transfer. The "originator" can refer to the originating VASP, the originating customer of the VASP or both.

- **Beneficiary**: The recipient of a blockchain transaction and therefore also the recipient of a TRISA transfer. The "beneficiary" can refer to the beneficiary VASP, the beneficiary customer of the VASP or both.

- **Natural Person**: A "Natural Person" is a human user. Originating and beneficiary customers are considered Natural Persons in accordance with the [IVMS101 specification](https://intervasp.org/).

- **Legal Person**: A "Legal Person" is an organization or legal entity. Originating and beneficiary VASPs are considered Legal Persons in accordance with the [IVMS101 specification](https://intervasp.org/).

- **Local vs Remote VASP**: A reference to the source of peer-to-peer traffic in an information exchange. The "local VASP" usually refers to the service you are running, while the "remote VASP" usually refers to some other VASP in the TRISA network. Local vs. Remote is often used interchangeably with originator vs. beneficiary. If you are initiating the transaction then the local VASP is the originator and the remote VASP is the beneficiary. If you are receiving a transaction then the local VASP is the beneficiary and the remote VASP is the originator.

- **Travel Rule**: Record-keeping rules for transfers between financial institutions that allow law enforcement agencies to prevent illicit finance (e.g. money laundering or the financing of terrorism).

- **VASP**: Virtual Asset Service Provider. A legal entity (usually a business) that manages and transfers virtual assets and are required by the Travel Rule to conduct information exchanges. Compliance exchanges in the TRISA network are between VASPs.

- **KYC**: Know Your Customer (or "KYC") is a due diligence practice used by VASPs to verify the identities and potential risks of their clients, users, and other counterparties. KYC procedures are intended to help banks and other financial institutions combat money laundering and other financial crimes.

- **AML**: Anti-Money Laundering (or "AML") refers to laws and regulations intended to stop criminals from disguising illegally obtained funds as legitimate income.



## Cryptographic Terminology

- **mTLS**: Mutual Transport Layer Security is an encryption protocol that authenticates both the client and the server in a network connection and encrypts communications between the parties so that data cannot be read in flight. mTLS is an extention of TLS (formerly SSL) that requires both sides of a network connection to have a certificate that establishes their identity and which can be used to encrypt packets sent on the channel.

- **Symmetric Cryptography**: Both encryption and decryption of data are performed using a single secret key that must be shared by the sender and the recipient. Shared secrets introduce the problem of how to share the secret key, however symmetric cryptographic algorithms are usually faster and better for bulk encryption of larger amounts of data. Generally, secrets are shared by asymmetric-key encryption and data encrypted using symmetric encryption.

- **Secret Key Cryptography**: See _symmetric cryptography_.

- **Public Key Cryptography**: A cryptographic method that uses a pair of related keys. Each pair consists of a _public key_, which can be shared with others, and a _private key_, which must not be shared with anyone but the owner. In practice, data can be encrypted with a public key but only decrypted with the private key.

- **Asymmetric Cryptography**: See _public key cryptography_.

- **Digital Signature**: A mathematical method that produces a _signature_ of data, e.g. some other piece of data that summarizes or describes the original data, usually via a hashing method. If the original data changes, its digital signature will change, therefore digital signatures are generally used as proof that the original data has not been tampered with, particularly if the signature has been generated cryptographically. For example, if a certificate has a signature that is signed with the private key of the certificate authority, anyone with the CA's public key can verify the signature of the certificate, ensuring it was the certificate produced by the CA.

- **HMAC**: Hash Message Authentication Code. A secure method of producing a digital signature for data that uses a cryptographic hash function and a secret key. HMACs are used to verify that data has not been modified or changed, see also _digital signature_.

- **Certificate**: Usually a reference to an X.509 certificate, a standard format for public key infrastructure. An X.509 certificate is a digital document that securely associates cryptographic key pairs with identities. Certificates are signed by a certificate authority to guarantee their provenance and contain subject information and other metadata concerning , the keypair, and its usage. Generally a reference to a certificate refers to the public key -- the part of the certificate that is shared for authentication purposes. However when certificates are issued, they are issued as a public/private key pair.

- **Identity Certificate**: a TRISA specific term that refers to the certificate issued by the TRISA CA to a VASP entity that they should use to connect to other VASPs in the TRISA network via mTLS.

- **Sealing Certificate (or keys)**: a TRISA specific term that refers to a key pair that is used to seal secure envelopes in the TRISA protocol. Sealing keys may be certificates issued by the TRISA CA, or they may be keys generated by the VASP and exchanged during the TRISA protocol.

- **Certificate Authority**: an entity that issues and revokes certificates and whose public keys are used to establish trust in the identities its issued certificates provide. Certificate authorities usually control cryptographic hardware and the "root of trust" -- a digital key pair that is used to generate and sign intermediate and leaf certificates for public key infrastructure.