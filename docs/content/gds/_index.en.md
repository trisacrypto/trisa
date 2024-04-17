---
title: "Global Directory Service"
date: 2021-06-14T11:21:42-05:00
lastmod: 2022-06-22T16:45:54-04:00
description: "The Global TRISA Directory Service (GDS)"
weight: 25
---

The TRISA organization hosts the TRISA Global Directory Service (GDS) on behalf of the TRISA network. The TRISA protocol specifies interactions with GDS. GDS facilitates peer-to-peer exchanges between TRISA members as follows:

- By issuing mTLS certificates to verify exchanges
- By providing discovery services for finding TRISA endpoints
- By providing certificate and KYCV (Know Your Counterparty VASP) information for verification

It is important to note that TRISA is a peer-to-peer network with no centralized authority for collecting or exchanging Travel Rule data.

GDS serves as a Certificate Authority for TRISA exchanges. After VASPs submit required information and are verified, GDS issues them Identity Certificates which VASPs can use to establish mTLS connections with counterparties, thereby securing communications that contain private Originator and Beneficiary data. These certificates are issued after extended validation and prove that the VASP is a trusted member of the TRISA network. In this way, GDS does not control the exchange; rather it helps to confirm the identities of parties involved in Travel Rule information exchanges.

GDS also serves as a decentralized store of member information, including member node (aka endpoint) addresses, TRIXO form details, and public keys. TRISA members can access the directory listing of other verified members, search, and lookup VASP counterparties. In this way, GDS helps members to make informed compliance decisions before sending or receiving large sums of virtual assets.

GDS is replicated across multiple continents. The servers hosting GDS are in three regions: US, EU, and Singapore. The servers are decentralized and geo-replicated to ensure that the GDS is consistent, available, and fault-tolerant. TRISA plans to expand to more regions in the future.

GDS also manages the certificate revocation list (CRL) to maintain the network over time. The directory issues sealing keys and manages revocation and reissuance of certificates.

This documentation describes the TRISA implementation of the directory service and TRISA-specific interactions with it. For details about working with the GDS API, visit the [API documentation]({{% ref "/gds/api" %}}).

## Networks

TRISA currently operates two directory services: a TestNet (trisatest.net) and the MainNet (vaspdirectory.net). The [TestNet]({{% ref "/testing" %}}) is intended to facilitate development and integration and should not be used for actual compliance exchanges. The MainNet is separated from the TestNet with a completely different certificate authority, and certificates issued to TestNet nodes cannot be used to connect to MainNet nodes and vice-versa.

Connect to the GDS and register for certificates with the following endpoints/urls:

| Directory         | Network | Website                   | gRPC Endpoint               |
|-------------------|---------|---------------------------|-----------------------------|
| trisatest.net     | TestNet | https://trisatest.net     | `api.trisatest.net:443`     |
| vaspdirectory.net | MainNet | https://vaspdirectory.net | `api.vaspdirectory.net:443` |

## Registered Directories

TRISA supports the idea of different directory services that can interoperate by exchanging VASP records with each other. A directory service by definition is a system that has an intermediate certificate authority under one of the TRISA root authority networks (e.g. TestNet or MainNet) and can issue leaf certificates via the intermediate authority. Directory services exchange records with each other to facilitate lookups.

Currently the only registered directories are the TRISA hosted directory services.