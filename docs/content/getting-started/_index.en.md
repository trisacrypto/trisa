---
title: Getting Started
date: 2021-04-23T01:35:35-04:00
lastmod: 2022-08-10T13:22:02-04:00
description: "Getting Started with TRISA"
weight: 5
---

Welcome! If you're new to TRISA, this guide will provide a path to help get you started.

- If you just want to get started as quickly as possible, consider simply using the [**TRISA Envoy self-hosted node**]({{% ref "envoy.md" %}}).
- If you are a **hands-on developer**, you may wish to skip to the section for [**technical implementers**]({{% relref "getting-started#dev" %}}).
- If you are an **administrator or technical leader** at your organization, you may wish to skip to the [**administrative guide**]({{% relref "getting-started#admin" %}}).

## Getting Started with TRISA

There are three key steps to getting started integrating TRISA into your organization's Travel Rule solution:

1. *Register with the TRISA Directory Service*: Before you can integrate the TRISA protocol into your VASP software, you must [create an account and register](https://trisa.directory/guide) with the TRISA Global Directory Service (GDS). This process is typically done by an administrator or technical leader at the organization and takes 1-2 business days. For more complete details, visit the documentation on [registration]({{% ref "/joining-trisa/registration" %}}).
2. *Implement the TRISA Protocol using TestNet*: To integrate TRISA, you will need to implement a TRISA node that can act as both a server and a client in TRISA exchanges. This process will require a team of developers and generally takes a few months. The TRISA [TestNet]({{% ref "/testing" %}}) is designed to support you with the development and testing of your TRISA node. It will enable you to perform live tests and validate transactions that share sensitive information safely and securely. Check out the [implementation overview]({{% ref "/getting-started#overview" %}}) for more details about how to get started.
3. *Implement your MainNet TRISA Node*: Once you have fully tested your implementation using the TestNet and RobotVASPs, you can quickly switch to the production TRISA Network by installing MainNet certificates to your TRISA node. Note that this may require registering for MainNet certificates if those were not requested in step 3.

This page will provide resources for getting started for both VASP administrators as well as technical implementation teams.

## For Administrators {##admin}

If you're an administrator or technical leader whose organization is using (or planning to adopt) the TRISA protocol, this section is for you!

These are the key portions of the documentation you will need to get started:
- [Joining TRISA]({{% ref "/joining-trisa" %}})
- [Registering with the Global Directory Service]({{% ref "/joining-trisa/registration" %}})
- [TRISA FAQ]({{% ref "/reference/faq" %}})
- [TRISA Glossary]({{% ref "/reference/glossary" %}})
- [External links and resources for TRISA]({{% ref "/reference" %}})

Integrating TRISA will require a team of engineers capable of implementing the TRISA protocol. When considering setting up your own server to host your own TRISA node, you must consider items necessary that may incur significant costs and resources, such as the server itself, long-term data storage solution, and developer time to configure and test. If your organization does not have access to a technical team or resources, you may instead choose to integrate with a 3rd-party TRISA solution. A list of some of the commercial TRISA solutions is available in [this guide](https://trisa.io/regulators-guide/).


## For Technical Implementers {##dev}

If you're a developer whose organization is using (or planning to adopt) the TRISA protocol, this section is for you!

These are the key portions of the documentation you will need to get started:
- [TRISA Protocol and API]({{% ref "/api" %}})
- [Working with TRISA Data]({{% ref "/data" %}})
- [Testing and Deployment]({{% ref "/testing" %}})

### Prerequisites

To begin setup, you’ll need the following:

*   Identity Certificates (from TRISA GDS registration)
*   The public key used for the CSR to obtain your certificate
*   The associated private key
*   The host name of the TRISA directory service
*   Ability to bind to the `address:port` that is associated with your VASP in the TRISA directory service.

### Integration Overview {##overview}

Integrating VASPs must run their own implementation of the protocol. Integrators are expected to integrate incoming transfer requests and key exchanges and may optionally also integrate outgoing transfer requests and key exchanges.

Integrating the TRISA protocol involves both a client component and server component.

The client component will interface with the TRISA Global Directory Service (GDS) instance to lookup other VASPs that integrate the TRISA messaging protocol. The client component is utilized for outgoing transactions from your VASP to verify the receiving VASP is TRISA compliant.

The server component receives requests from other VASPs that integrate the TRISA protocol and provides responses to their requests. The server component provides callbacks that must be implemented so that your VASP can return information satisfying the TRISA Network protocol.

Currently, a reference implementation of the TRISA Network protocol is available in Go.

[https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go](https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go)

If a language beside Go is required, client libraries may be generated from the [protocol buffers](https://github.com/trisacrypto/trisa/tree/main/proto) that define the TRISA Network protocol.

### Integration Notes

The TRISA Network protocol defines how data is transferred between participating VASPs. The recommended format for data transferred for identifying information is the [IVMS101 data format]({{% relref "data/ivms" %}}). It is the responsibility of the implementing VASP to ensure the identifying data sent/received satisfies the FATF Travel Rule.

The result of a successful TRISA transaction results in a key and encrypted data that satisfies the FATF Travel Rule. TRISA does not define how this data should be stored once obtained. It is the responsibility of the implementing VASP to handle the secure storage of the resulting data for the transaction.

Some other considerations you will need to make to be prepared to fully integrate TRISA include:

1. How will your TRISA endpoint integrate with your existing backend systems?
2. How will you handle key management (e.g. your own private keys as well as the public keys of your counterparties)?
3. Are you prepared to store [envelopes]({{% ref "/data/envelopes" %}}) securely and in compliance with privacy regulations once you have received them from your counterparties?

For more considerations, see our [Best Practices]({{% ref "/reference/best-practices" %}}) documentation.
