---
title: Testing and Deployment
date: 2021-06-14T11:20:10-05:00
lastmod: 2022-07-07T15:34:15-04:00
description: "Describing the Integration and Development with TestNet"
weight: 30
---

The TRISA TestNet has been established to provide a demonstration of the TRISA peer-to-peer protocol, host "robot VASP" services to facilitate TRISA integration, and facilitate public key exchange and peer discovery via TRISA's Global Directory Service (GDS).  The TestNet instance is designed for testing, and the registration process is streamlined in the TestNet to facilitate quick integration. The Testnet enables you to test and validate transactions that share sensitive information safely and securely.

For reference, the [TRISA Protocol documentation]({{% relref "api/protocol" %}}) provides additional information on enabling the peer-to-peer exchange of compliance information.

The TRISA TestNet is comprised of several services, including:

- A TestNet [Certificate Authority]({{% relref "joining-trisa/ca" %}}) that issues TestNet Identity Certificates (*note that these are distinct from MainNet certificates and not interchangeable*).
- [TRISA Directory Service](https://trisa.directory/) - a user interface to explore the TRISA Global Directory Service and register to become a TRISA member
- [TestNet Demo](https://vaspbot.com) - a demo site to show TRISA interactions between “robot” VASPs that run in the TestNet

The TestNet also hosts three ["robot VASPs" (rVASPs)]({{% relref "testing/rvasps" %}}) that have been implemented as a convenience for TRISA members to integrate their TRISA services and validate the compliance solution safely. The primary rVASP is Alice, a secondary for demo purposes is Bob, and an "evil" rVASP to test interactions with non-verified TRISA members.

The TestNet also provides a [command line utility]({{% relref "testing/trisa-cli" %}}) for interacting with the API for administrative and debugging purposes, using the testnet certificates.

![TestNet Architecture](/img/testnet_architecture.png)

## Joining the TestNet

The following steps are required to join the TestNet:

1. [Register](https://trisa.directory/guide) with the GDS to create your TRISA Account with your VASP email address, where you can opt-in for TestNet Certificates. During registration, you can add collaborators within your organization.

2. Complete the VASP Verification form and due diligence process. Once approved, you will gain access to the TestNet.

3. Set up your TRISA node and implement the [TRISA API]({{% relref "/api" %}}).
