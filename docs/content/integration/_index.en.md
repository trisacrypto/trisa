---
title: Integration Overview
date: 2021-04-23T01:35:35-04:00
lastmod: 2021-04-23T01:35:35-04:00
description: "Describes how to integrate the TRISA protocol in the TestNet"
weight: 20
---

## TRISA Integration Overview

1. Register with a TRISA Directory Service
2. Implement the TRISA Network protocol
3. Implement the TRISA Health protocol

## VASP Directory Service Registration

### Registration Overview

Before you can integrate the TRISA protocol into your VASP software, you must register with a TRISA Directory Service (DS).  The TRISA DS provides public key and TRISA remote peer connection information for registered VASPs.

Once you have registered with the TRISA DS, you will receive a KYV certificate.  The public key in the KYV certificate will be made available to other VASPs via the TRISA DS.

When registering with the DS, you will need to provide the `address:port` endpoint where your VASP implements the TRISA Network service. This address will be registered with the DS and utilized by other VASPs when your VASP is identified as the beneficiary VASP.

For integration purposes, a hosted TestNet TRISA DS instance is available for testing.  The registration process is streamlined in the TestNet to facilitate quick integration.  It is recommended to start the production DS registration while integrating with the TestNet.


### Directory Service Registration

To start registration with the TRISA DS, visit website at [https://vaspdirectory.net/](https://vaspdirectory.net/)

You can select the "Register" tab to begin registration. Note that you can use this website to enter your registration details on a field-by-field basis, or to upload a JSON document containing your registration details.

This registration will result in an email being sent to all the technical contacts specified in the JSON file.  The emails will guide you through the remainder of the registration process.  Once you’ve completed the registration steps, TRISA TestNet administrators will receive your registration for review.

Once TestNet administrators have reviewed and approved the registration, you will receive a KYV certificate via email and your VASP will be publicly visible in the TestNet DS.


## Implementing the TRISA P2P Protocol


### Prerequisites

To begin setup, you’ll need the following:



*   KYV certificate (from TRISA DS registration)
*   The public key used for the CSR to obtain your certificate
*   The associated private key
*   The host name of the TRISA directory service
*   Ability to bind to the address:port that is associated with your VASP in the TRISA directory service.


### Integration Overview

Integrating the TRISA protocol involves both a client component and server component.

The client component will interface with a TRISA Directory Service (DS) instance to lookup other VASPs that integrate the TRISA messaging protocol.  The client component is utilized for outgoing transactions from your VASP to verify the receiving VASP is TRISA compliant.

The server component receives requests from other VASPs that integrate the TRISA protocol and provides responses to their requests.  The server component provides callbacks that must be implemented so that your VASP can return information satisfying the TRISA Network protocol.

Currently, a reference implementation of the TRISA Network protocol is available in Go.

[https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go](https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go)

Integrating VASPs must run their own implementation of the protocol.  If a language beside Go is required, client libraries may be generated from the protocol buffers that define the TRISA Network protocol.

Integrators are expected to integrate incoming transfer requests and key exchanges and may optionally also integrate outgoing transfer requests and key exchanges.

### Integration Notes

The TRISA Network protocol defines how data is transferred between participating VASPs.  The recommended format for data transferred for identifying information is the IVMS101 data format.  It is the responsibility of the implementing VASP to ensure the identifying data sent/received satisfies the FATF Travel Rule.

The result of a successful TRISA transaction results in a key and encrypted data that satisfies the FATF Travel Rule.  TRISA does not define how this data should be stored once obtained.  It is the responsibility of the implementing VASP to handle the secure storage of the resulting data for the transaction.

