---
title: FAQ
date: 2021-06-14T15:59:09-05:00
lastmod: 2022-06-08T12:34:54-04:00
description: "Frequently asked questions"
weight: 10
---

## What does the "TRISA" acronym mean?

TRISA stands for Travel Rule Information Sharing Architecture.

## What is TRISA?

TRISA is a peer-to-peer network that enables compliance with the Financial Action Task Force (FATF) and the Financial Crimes Enforcement Network (FinCEN) Travel Rule for transaction identity information.

TRISA is also a protocol that enables the secure exchange of PII between members.

Read more about TRISA in the [TRISA whitepaper](https://trisa.io/trisa-whitepaper/).

## How does TRISA work?

TRISA creates a mechanism for VASPs to exchange Travel Rule data packets by performing remote procedure calls using [gRPC]({{< relref "reference/faq#grpc" >}}) via [mTLS]({{< relref "reference/faq#mtls" >}}) connection. VASPs format the data as protocol buffers that contain details such as [IVMS101]({{< relref "reference/faq#ivms" >}}) identity and transaction payload data.

TRISA maintains open source [Github repositories](https://github.com/trisacrypto) available to VASPs as they build, test, and deploy their TRISA node into production.

## What is the Global Directory Service? {##gds}

The Global Directory Service (GDS) is a geo-distributed store of VASP endpoint and certificate details. VASPs participating in the TRISA protocol use the GDS to look up the TRISA endpoint of their intended counterparty.

Instructions for working with the Global Directory Service can be found in the documentation of the [GDS](https://trisa.dev/gds/).

## How is TRISA decentralized?

TRISA is a peer-to-peer network. The TRISA protocol enables VASPs to communicate directly with each other and not through a centralized hub. Peers in the TRISA network are responsible for coordinating their own information exchanges using the TRISA protocol and the Global Directory Service.

TRISA centralizes a small component of the protocol by acting as a certificate authority. This allows VASPs to establish trust with verified counterparties and understand their compliance requirements.

## How does TRISA safeguard PII?

TRISA safeguards Personally Identifiable Information (PII) in flight and at rest.

In flight, [secure envelopes]({{< relref "reference/faq#envelope" >}}) are exchanged over an mTLS encrypted channel created using the identity certificates issued by the TRISA network. TRISA members can use each other's public key addresses to open a secure line of communication to transmit users’ PII.

TRISA uses a trusted Certified Authority (CA) model, and only verified VASPs are granted certificates from the CA. Certificate authorities offer a root of trust to anchor identities to a chain of trusted entities. The CA model safeguards against a single point of failure, provides protection from attacks, and is scalable to accommodate the growing crypto landscape.

Using the TRISA protocol, VASPs exchange secure envelopes formatted as protocol buffer messages. These messages use a combination of asymmetric (public/private key cryptography using TRISA issued or peer-to-peer exchanged signing keys) and symmetric cryptography so that the VASPs can securely store the envelope at rest in a backend of their choice while maintaining full repudiation of the exchange.

## Who's in TRISA?

TRISA is open to organizations that offer virtual asset or digital asset services. Organizations must have a legitimate business purpose to join TRISA. Member organizations may be:
- Virtual Asset Service Providers (VASPs)
- Crypto Asset Service Providers (CASPs)
- Money Service Businesses (MSBs)
- Traditional financial services institutions
- Regulatory bodies

For more information about joining TRISA, please review [TRISA’s VASP Verification Process](https://vaspdirectory.net/getting-started).

## How can I look up TRISA members?

If you are a member, you can look up TRISA members using the [Members Endpoint](https://trisa.dev/gds/members/).

If you know a VASP’s TRISA Endpoint or Common Name, you can look it up in TRISA’s Global Directory via https://vaspdirectory.net.

## Why does my common name have to match my endpoint?

The Common Name typically matches the Endpoint (e.g. `trisa.myvasp.com:443`), without the port number at the end (e.g. `trisa.myvasp.com`) and is used to identify the subject in the X.509 certificate.

## What are secure envelopes? {##envelope}

Secure Envelopes wrap the identity and blockchain transaction payloads, applying additional encryption and digital signatures for verification.

![secure envelopes](/img/secure_envelopes.png).

The [Secure Envelope documentation](https://trisa.dev/secure-envelopes/) discusses its implementation further.

## How does mTLS work? {##mtls}

Mutual Transport Layer Security (mTLS) is a mechanism for mutual authentication between services or servers. Also known as two-way authentication, it ensures that the parties at each end of a connection are who they claim to be.

## What is IVMS 101? {##ivms}

The interVASP Messaging Standard (IVMS 101) is an internationally recognized standard that helps with language encodings, numeric identification systems, phonetic name pronunciations, and standardized country codes (ISO 3166).

TRISA uses the interVASP Messaging Standard (IVMS101) standard to describe Originators and Beneficiaries in terms of Natural Persons (human users) and Legal Persons ( organizations or legal entities), as well as necessary information about parties such as geographic addresses.

Depending on your business details, specific fields may be required. For more information on IVMS 101, please see the documentation at [intervasp.org](https://intervasp.org/).

There is an [IVMS 101 Validator](https://ivmsvalidator.com/) which can be used to double check the formatting of an IVMS101 message to ensure it is correct.

For help marshaling and unmarshaling [IVMS101 identity payloads]({{< relref "ivms/" >}}), see the documentation about the [`ivms101` package in `trisa`](https://github.com/trisacrypto/trisa/tree/main/pkg/ivms101).

## How do I figure out where to connect to the counterparty? How do I get counterparty IVMS 101 info?

As part of the protocol, the Originator can use the [Global Directory Service]({{< relref "reference/faq#gds" >}}) to lookup the counterpoint endpoint, and sends a secure envelope providing their IVMS101 details. The Beneficiary can then verify and store the counterparty PII information needed for compliance. Next the Beneficiary can return a new secure envelope with their IVMS101 details so that the Originator can store the information for compliance.

## Why gRPC? {##grpc}

The [gRPC library](https://grpc.io/docs/) is an open-source Remote Procedure Call (RPC) framework that can run in any environment and can establish a secure communication.

gRPC includes a bidirectional streaming mode that allows long running connections with high throughput messaging between nodes. This mode is optional in TRISA but can be used to support batch messaging and increase the performance of TRISA messaging.

The use of gRPC also facilitates convenient encryption at rest to ensure that PII is safeguarded.

## What's the difference between TestNet and MainNet?

The [TestNet](https://trisa.dev/testnet/) is a sandbox environment that allows VASPs to test securely sharing the cryptocurrency transaction details required to meet the FATF Travel Rule requirements. The TestNet includes [“robot” VASPs](http://localhost:1313/testnet/rvasps/) that give users the ability to interact with the TestNet by simulating transactions to see how secure transactions are conducted. Once a VASP completes testing, the VASP can switch to MainNet, where live transactions take place.

It’s important to note that the reason that there are two networks is because those networks are issued from different intermediate certificate authorities. A VASP that has been issued a TestNet certificate cannot connect to a node that is running on MainNet and vice versa. In other words, the MainNet certificate authority will not recognize TestNet certificates. When you [submit a request for TRISA certificates](https://vaspdirectory.net/certificate/registration), you may simultaneously request certificates for TestNet and MainNet.
