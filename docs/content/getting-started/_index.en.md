---
title: Getting Started
date: 2021-04-23T01:35:35-04:00
lastmod: 2022-06-08T14:51:13-04:00
description: "Describes how to integrate the TRISA protocol in the TestNet"
weight: 20
---


## Getting Started with TRISA

There are three key steps to getting started integrating TRISA into your organization's Travel Rule solution:

1. Register with a TRISA Directory Service
2. Implement the TRISA Network protocol
3. Implement the TRISA Health protocol

This page will walk through each of these three steps, providing links to other resources as necessary.

## VASP Directory Service Registration

Before you can integrate the TRISA protocol into your VASP software, you must [register](https://vaspdirectory.net/certificate/registration) with the TRISA Global Directory Service (GDS).

### Registration Overview

The TRISA Global Directory Service (GDS) provides public key and TRISA remote peer connection information for registered VASPs. For more detailed information about the directory, see the documentation on the [GDS]({{< ref "/gds" >}}).

Once you have registered with the GDS and been verified, you will receive Identity Certificates. The public key in these certificates will be made available to other VASPs via the GDS.

When registering with the GDS, you will need to provide the `address:port` endpoint where your VASP implements the TRISA Network service. This address will be registered with the GDS and utilized by other VASPs when your VASP is identified as the beneficiary VASP.

For integration purposes, when you [register](https://vaspdirectory.net/certificate/registration) with the GDS, you can opt for either MainNet or TestNet Certificates, or both. The TestNet instance is designed for testing, and the registration process is streamlined in the TestNet to facilitate quick integration. The MainNet is design for production Travel Rule implementations. It is recommended to register for both MainNet and TestNet, specifying different endpoints to reduce confusion for your VASP counterparties.

For more complete details, visit the documentation on [registration]({{< ref "/joining-trisa/registration" >}}).

### Directory Service Registration

To start your registration, visit [https://vaspdirectory.net/](https://vaspdirectory.net/certificate/registration). Note that you can use this website to enter your registration details on a field-by-field basis, or to upload a JSON document containing your registration details. Several people at your organization (e.g. legal, technical, administrative points-of-contact) may need to collaborate to complete the needed information. The final step of registration will be a [pkcs12 password]({{< ref "/joining-trisa/pkcs12" >}}), which you must keep to decrypt the Identity Certificates that will be sent via email.

This registration will result in an email being sent to all the technical contacts specified through the webform or in the JSON file. The emails will guide you through the remainder of the registration process. Once you’ve completed the registration steps, TRISA administrators will receive your registration for review.

Once the administrators have reviewed and approved the registration, you will receive [pkcs12 password]({{< ref "/joining-trisa/pkcs12" >}})-protected, compressed Identity Certificate via email and your VASP will be publicly visible in the GDS.


## Implementing the TRISA P2P Protocol

### Prerequisites

To begin setup, you’ll need the following:

*   Identity Certificates (from TRISA GDS registration)
*   The public key used for the CSR to obtain your certificate
*   The associated private key
*   The host name of the TRISA directory service
*   Ability to bind to the `address:port` that is associated with your VASP in the TRISA directory service.

### Integration Overview

Integrating VASPs must run their own implementation of the protocol. Integrators are expected to integrate incoming transfer requests and key exchanges and may optionally also integrate outgoing transfer requests and key exchanges.

Integrating the TRISA protocol involves both a client component and server component.

The client component will interface with the TRISA Global Directory Service (GDS) instance to lookup other VASPs that integrate the TRISA messaging protocol. The client component is utilized for outgoing transactions from your VASP to verify the receiving VASP is TRISA compliant.

The server component receives requests from other VASPs that integrate the TRISA protocol and provides responses to their requests. The server component provides callbacks that must be implemented so that your VASP can return information satisfying the TRISA Network protocol.

Currently, a reference implementation of the TRISA Network protocol is available in Go.

[https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go](https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go)

If a language beside Go is required, client libraries may be generated from the [protocol buffers](https://github.com/trisacrypto/trisa/tree/main/proto) that define the TRISA Network protocol.

### Integration Notes

The TRISA Network protocol defines how data is transferred between participating VASPs. The recommended format for data transferred for identifying information is the IVMS101 data format. It is the responsibility of the implementing VASP to ensure the identifying data sent/received satisfies the FATF Travel Rule.

The result of a successful TRISA transaction results in a key and encrypted data that satisfies the FATF Travel Rule. TRISA does not define how this data should be stored once obtained. It is the responsibility of the implementing VASP to handle the secure storage of the resulting data for the transaction.

Some other considerations you will need to make to be prepared to fully integrate TRISA include:

1. How will your TRISA endpoint integrate with your existing backend systems?
2. How will you handle key management (e.g. your own private keys as well as the public keys of your counterparties)?
3. Are you prepared to store [envelopes]({{< ref "/secure-envelopes" >}}) securely and in compliance with privacy regulations once you have received them from your counterparties?


