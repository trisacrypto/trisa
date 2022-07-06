---
title: Best Practices
date: 2021-07-06T15:59:09-05:00
lastmod: 2021-07-06T16:31:17-05:00
description: "Best Practices Page"
weight: 200
---

## Long Term Data Storage
TRISA recommends that Travel Rule information transfers be stored post-transfer as encrypted Secure Envelope protocol buffers. The unsealing keys for these envelopes should be stored separately, and should be deleted once the compliance period has ended, rendering the Envelopes un-openable. This is commonly referred to as "deletion by erasure". Compliance periods differ by region and jurisdiction but are typically between 5 and 7 years.

## Key Management
TRISA recommends that TRISA implementers use the [Key Handler package](https://github.com/trisacrypto/trisa/tree/main/pkg/trisa/handler) for key management. (Documentation coming soon)

## Throughput
Depending on the volume of Travel Rule transactions that your organization executes on a regular basis, you may wish to consider using the bidirectional streaming mode, `TransferStream`, which will support more throughput.

## Deployment
For deployment, TRISA implementers should be prepared to install their Identity Certificates (which will need to be reissued and reinstalled on an annual basis) to support routine mTLS termination. 

Implementers must also ensure that their existing databases are configured to support responding to routine Travel Rule information requests. Additional storage may need to be configured to store the results of these requests (e.g. encrypted [Secure Envelopes](https://trisa.dev/secure-envelopes/)).