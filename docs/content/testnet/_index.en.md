---
title: TestNet
date: 2021-06-14T11:20:10-05:00
lastmod: 2022-04-10T09:32:16-05:00
description: "The TRISA TestNet"
weight: 30
---

The TRISA TestNet has been established to provide a demonstration of the TRISA peer-to-peer protocol, host "robot VASP" services to facilitate TRISA integration, and is the location of the primary TRISA Directory Service that facilitates public key exchange and peer discovery.

{{% figure src="/img/testnet_architecture.png" %}}

The TRISA TestNet is comprised of the following services:

- [TRISA Directory Service](https://trisatest.net) - a user interface to explore the TRISA Global Directory Service and register to become a TRISA member
- [TestNet Demo](https://vaspbot.net) - a demo site to show TRISA interactions between “robot” VASPs that run in the TestNet

The TestNet also hosts three robot VASPs or rVASPs that have been implemented as a convenience for TRISA members to integrate their TRISA services. The primary rVASP is Alice, a secondary for demo purposes is Bob, and to test interactions with non-verified TRISA members, there is also an "evil" rVASP.