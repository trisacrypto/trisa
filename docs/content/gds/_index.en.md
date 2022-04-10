---
title: "Directory Service"
date: 2021-06-14T11:21:42-05:00
lastmod: 2022-04-10T09:32:16-05:00
description: "The Global TRISA Directory Service (GDS)"
weight: 50
---

The Global TRISA Directory Service (GDS) facilitates peer-to-peer exchanges between TRISA members as follows:

- By providing discovery services for finding TRISA endpoints
- By issuing mTLS certificates to verify exchanges
- By providing certificate and KYC information for verification

Interactions with a Directory Service are specified by the TRISA protocol. Currently, the TRISA organization hosts the GDS on behalf of the TRISA network. This documentation describes the TRISA implementation of the directory service and TRISA-specific interactions with it.

## Networks

TRISA currently operates two directory services: a TestNet (trisatest.net) and the MainNet (vaspdirectory.net). The [TestNet]({{< ref "/testnet" >}}) is intended to facilitate development and integration and should not be used for actual compliance exchanges. The MainNet is separated from the TestNet with a completely different certificate authority, and certificates issued to TestNet nodes cannot be used to connect to MainNet nodes and vice-versa.

Connect to the GDS and register for certificates with the following endpoints/urls:

| Directory         | Network | Website                   | gRPC Endpoint               |
|-------------------|---------|---------------------------|-----------------------------|
| trisatest.net     | TestNet | https://trisatest.net     | `api.trisatest.net:443`     |
| vaspdirectory.net | MainNet | https://vaspdirectory.net | `api.vaspdirectory.net:443` |

## Registered Directories

TRISA supports the idea of different directory services that can interoperate by exchanging VASP records with each other. A directory service by definition is a system that has an intermediate certificate authority under one of the TRISA root authority networks (e.g. TestNet or MainNet) and can issue leaf certificates via the intermediate authority. Directory services exchange records with each other to faciliate lookups.

Currently the only registered directories are the TRISA hosted directory services.