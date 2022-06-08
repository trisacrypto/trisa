---
title: FAQ
date: 2021-06-14T15:59:09-05:00
lastmod: 2022-06-08T12:34:54-04:00
description: "Frequently asked questions"
weight: 10
---

 
### How does TRISA work?
The Travel Rule Information Sharing Architecture (TRISA) is a peer-to-peer network that enables compliance with the Financial Action Task Force (FATF) and the Financial Crimes Enforcement Network (FinCEN) Travel Rule for transaction identity information.
 
TRISA uses a trusted Certified Authority (CA) model, adopted from the International Organization for Standardization (ISO) and used to secure most e-commerce and governmental communications, to actively verify the identity of Virtual Asset Service Providers (VASPs). Certificate authorities offer a root of trust to anchor identities to a chain of trusted entities. The architecture catalogs verified VASPs’ public key addresses, which can be used to open a secure line of communication between VASPs to transmit users’ Personally Identifiable Information (PII). The CA model safeguards against a single point of failure, provides protection from attacks, and is scalable to accommodate the growing crypto landscape.
 
TRISA creates a mechanism for VASPs to comply with, where VASPs can host their own TRISA node (server) to exchange Travel Rule data packets. Travel Rule data is encrypted in flight and at rest using [secure envelopes](https://www.trisa.io/wp-content/uploads/2022/03/Secure-Envelopes.png). There is no centralized authority collecting or exchanging Travel Rule data.
 
TRISA operates as the Trusted VASP Certificate Authority for the network of trusted VASPs, also known as the Global Directory Service (GDS). GDS facilitates peer-to-peer exchanges between TRISA members by issuing mutual Transport Layer Security (mTLS) certificates, providing discovery services for finding TRISA endpoints, and providing certificate and Know your Client (KYC) information for verification.
 
For the peer-to-peer exchange protocol, TRISA uses [protocol buffers](https://github.com/trisacrypto/trisa/blob/main/proto/trisa/data/generic/v1beta1/transaction.proto) for IVMS messages, natural and legal person identification, secure envelopes, and more.
 
The [TRISA white paper](https://www.trisa.io/trisa-whitepaper/) further describes the Travel Rule compliance and how it implements the form of a Remote Procedure Call (RPC) specification for inter-VASP communication, and a reference library implemented in the Go programming language.
 
TRISA maintains open source [Github repositories](https://github.com/trisacrypto) available to VASPs as they build, test, and deploy their TRISA node into production.
 
### How is TRISA decentralized?
TRISA is a peer-to-peer network of trusted counter-parties seeking to securely exchange Personally Identifiable Information to comply with the Travel Rule. TRISA’s protocol enables VASPs to communicate directly with each other and not through a centralized hub.
 
However, TRISA centralizes a small component of the protocol by acting as a certificate authority. This allows VASPs to understand exactly who the counter-party is and what their compliance requirements are.
 
### Who's in TRISA? How can I look up TRISA members?
TRISA is open to organizations that offer virtual asset or digital asset services. Organizations must have a legitimate business purpose to join TRISA. Member organizations may be:
- Virtual Asset Service Providers (VASPs)
- Crypto Asset Service Providers (CASPs)
- Money Service Businesses (MSBs)
- Traditional financial services institutions
- Regulatory bodies
 
For more information about joining TRISA, please review [TRISA’s VASP Verification Process](https://vaspdirectory.net/getting-started).
 
You can look up TRISA members in [TRISA’s Global Directory Service](https://vaspdirectory.net/). To do so, you will need to know the VASP’s TRISA Endpoint or Common Name.
 
### Why does my common name have to match my endpoint?
The Common Name typically matches the Endpoint, without the port number at the end (e.g. trisa.myvasp.com) and is used to identify the subject in the X.509 certificate.
 
### What are secure envelopes?
The Secure Envelope can wrap the identity and blockchain transaction payload, applying additional encryption and digital signatures for verification.
 
The [Secure Envelope documentation](https://trisa.dev/secure-envelopes/) discusses its implementation further.
 
### How does mTLS work?
mutual Transport Layer Security (mTLS) is a mechanism for mutual authentication between services or servers. Also known as two-way authentication, it ensures that the parties at each end of a connection are who they claim to be.
 
In flight, Secure Envelopes are exchanged over an mTLS encrypted channel created by the identity certificates issued by the TRISA network. Only TRISA members can communicate via this channel. At rest, the Secure Envelopes use a combination of asymmetric (public/private key cryptography using TRISA issued or peer-to-peer exchanged signing keys) and symmetric cryptography so that the VASPs can securely store the envelope at rest in a backend of their choice while maintaining full repudiation of the exchange.
 
### What is IVMS 101?
interVASP Messaging Standard (IVMS 101) is an internationally recognized standard that helps with language encodings, numeric identification systems, phonetic name pronunciations, and standardized country codes (ISO 3166).
 
TRISA uses the interVASP Messaging Standard (IVMS101) standard to describe two primary types: a Natural Person to define a human user; and a Legal Person to define an organization or legal entity.
 
Depending on your business details, specific fields may be required. For more information on IVMS 101, please see the documentation at [intervasp.org](https://intervasp.org/).
 
Members can also validate their information through the [IVMS 101 Validator](https://ivmsvalidator.com/).
 
### How do I figure out where to connect to the counter-party? How do I get counter-party IVMS 101 info?
The Global Directory can be used by VASPs to determine if the counter-party is a VASP and if that VASP is in a jurisdiction that enforces travel rule legislation, thus allowing Alliance members to make informed compliance decisions before sending or receiving large sums of virtual assets. The Global Directory will be augmented with:
- Additional entity details
- Regulator/licensing information
- Privacy information
- Travel rule sharing end-point and protocols registered VASPs.
 
As part of the protocol, Beneficiary VASP verifies counter-party PII information (and may update it as necessary) and returns a signed secure envelope with the receipt timestamp using the same message identifier.
 
Counter-party IVMS101 information is exchanged via gRPC. Each message is wrapped in an encrypted envelope. The goal of this is to be able to store the messages as-is and maintain full repudiation. Multiple encryption schemes are to be implemented in the future.
 
### Why gRPC?
The [gRPC library](https://grpc.io/docs/) is an open-source Remote Procedure Call (RPC) framework that can run in any environment and can establish a secure communication.
 
gRPC includes a bidirectional streaming mode that allows long running connections with high throughput messaging between nodes. This mode is optional in TRISA but can be used to support batch messaging and increase the performance of TRISA messaging.
 
### What's the difference between TestNet and MainNet?
The [TestNet](https://trisa.dev/testnet/) is a sandbox environment that allows VASPs to test securely sharing the cryptocurrency transaction details required to meet the FATF Travel Rule requirements. The TestNet includes [“robot” VASPs](http://localhost:1313/testnet/rvasps/) that give users the ability to interact with the TestNet by simulating transactions to see how secure transactions are conducted. Once a VASP completes testing, the VASP can request a certificate for MainNet, where live transactions take place. It’s important to note that the reason that there are two networks is because those networks are issued from different intermediate certificate authorities. A VASP that has been issued a TestNet certificate cannot connect to a node that is running on MainNet and vice versa. In other words, the MainNet certificate authority will not recognize TestNet certificates.
