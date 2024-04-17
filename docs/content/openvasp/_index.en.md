---
title: OpenVASP/TRP Integration
date: 2022-06-17T09:15:46-04:00
lastmod: 2022-06-17T16:53:39-04:00
description: "Documentation on how to integrate TRISA with the OpenVASP/TRP Protocol"
weight: 45
---

The [OpenVASP Association](https://www.openvasp.org/) implements an open protocol for travel rule called [TRP](https://www.openvasp.org/trp) that is similar to TRISA but puts more focus on ease of information sharing and human-in-the-loop review than on verified counterparties and cryptography. TRISA and OpenVASP have worked together over the past several months to create a "bridge" that will allow partial compatibility between the two protocols. This page describes how implementers of TRISA nodes might integrate the bridge in their node to respond to TRP requests.

For more information on the TRP protocol, please refer to the [TRP Documentation](https://gitlab.com/OpenVASP/travel-rule-protocol/-/blob/master/core/specification.md?ref_type=heads).

## TRP Workflow

![TRP Workflow](/img/trp_flow.png)

The TRP workflow uses HTTPS `POST` requests with JSON payloads to facilitate information exchange. The initial request endpoint is defined with an LNURL or Travel Address endpoint that the beneficiary must request from their VASP and send to the originator. Subsequent request endpoints are defined with callback URLs in the request itself. HTTP error codes and JSON payloads are used to communicate success or failure back to the counterparty.

In principle, then, a TRISA node must add an HTTP service to its node to accept and respond to these POST requests. TRISA has implemented [handlers and middleware](https://pkg.go.dev/github.com/trisacrypto/trisa@v0.4.0/pkg/openvasp) in the `trisacrypto/trisa` go package to make it easier to add a service to your TRISA node and to translate TRISA data structures into TRP data structures.

### Next Steps:

1. [Integrating a TRP Bridge Handler into your TRISA node]()
2. [Making Outgoing TRP Requests]()

## Considerations

As a TRISA node implementer, you have registered for mTLS certificates with TRISA's Global Directory Service and went through a rigorous KYC process to be verified as a VASP that must exchange PII information to comply with the Travel Rule. You are probably used to using the GDS to lookup counterparty endpoints and you've probably experienced significant time and effort implementing key management for the cryptographic requirements that TRISA uses for non-repudiation and secure storage. Counterparties that implement TRP do not have these same requirements.

Therefore, as you implement your TRP integration, you need to consider the following policies and integration standards.

1. **mTLS is a core part of TRP, however, TRP does not specify a Certificate Authority.** Your implementation must consider whether it wants to perform TRP only with TRISA issued certificates or if it is willing to allow other public CAs such as Verisign or Google.
2. **There is no directory of TRP VASPs, TRP discovery is facilitated by Travel Addresses.** TRP uses [Travel Addresses](https://www.21analytics.ch/blog/how-the-trp-travel-address-solves-the-fatf-travel-rule/) to solve counterparty identification. The TRISA bridge is able to parse both LNURLs and the newer Travel Address format, however for complete TRP integration, you will have to supply your accounts with Travel Addresses so that other TRP implementers can reach you as a beneficiary counterparty.
3. **TRP only supports transport-level cryptography, not payload-level cryptography.** There are three levels of cryptography supported by the TRISA bridge: no-cryptography, partial, insecure TRISA envelopes (encrypted but not sealed), and full TRISA compatibility. The first level is plain-vanill TRP and the second two levels are implemented using [TRP Extensions](https://gitlab.com/OpenVASP/travel-rule-protocol/-/tree/master/extensions).

### For TRP Implementers

Welcome, thank you for checking out TRISA! The best thing you can do for integration is to register for TRISA mTLS certificates on [vaspdirectory.net](https://vaspdirectory.net). If you perform mTLS with TRISA VASPs using TRISA certificates, that will help a lot in establishing trust and verification for consideration #1 above.

If you're interested in implementing the extensions for key exchange and parsing secure envelopes, you're more than welcome to use the Golang code in this library to get you started! If you're implementation is in another language, [please let us know](https://github.com/trisacrypto/trisa/issues) so that we can create library code to help you implement secure PII transfers.

Please see the [Getting Started Guide]() for more on how to implement TRISA-specific protocol details using the extensions.

### Policies

Given the above considerations, TRISA implementers will have to consider the following policies before TRP integration:

1. Allow certificate authorities other than the TRISA authority for mTLS?
2. Allow native TRP transfers without signatures or payload cryptography?
3. Allow [TRP message signing](https://gitlab.com/OpenVASP/travel-rule-protocol/-/blob/master/extensions/message-signing.md?ref_type=heads) for non-repudiation?
4. Require either TRISA or TRP message signing for non-repudiation?
5. Require public key exchange and secure envelope extension?

The answer to these policies considerations will determine how permissive your node is to accepting different types of transfers, but at the same time create more transfer cases that need to be handled by your node.

### Protocol Comparison

| Feature              | TRISA                                                                 | OpenVASP TRP                                 |
|----------------------|-----------------------------------------------------------------------|----------------------------------------------|
| Governance           | Delaware non-profit and technical working group                       | Swiss non-profit and technical working group |
| Non Repudiation      | Signed Envelopes                                                      | Signed JSON                                  |
| Data at Rest         | Secure Envelopes                                                      |                                              |
| Transport Encryption | mTLS 1.3                                                              | mTLS 1.3                                     |
| Exceptions/Errors    | Error codes                                                           | HTTP protocol errors                         |
| Message              | Protocol Buffers                                                      | JSON                                         |
| Data Types           | IVMS101, PayString, Generic                                           | IVMS101                                      |
| Authentication       | X.509 KYV Certs                                                       |                                              |
| Transport Protocol   | gRPC                                                                  | HTTPS                                        |
| Addressing           | VASP/Account                                                          | VASP/Account Travel Address  (LNURL)         |
| Discovery            | Sender provided, receiver verified, blockchain analytics, round robin | Sender provided                              |
| Onboarding           | TRIXO Questionnaire/GDS                                               |                                              |