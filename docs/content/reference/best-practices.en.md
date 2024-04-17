---
title: Best Practices
date: 2022-07-06T15:41:17-04:00
lastmod: 2022-08-10T16:34:11-04:00
description: "Best Practices for TRISA Implementers"
weight: 10
---

## Long Term Data Storage
TRISA recommends that Travel Rule information transfers be stored post-transfer as encrypted Secure Envelope protocol buffers. The unsealing keys for these envelopes should be stored separately, and should be deleted once the compliance period has ended, rendering the Envelopes un-openable. This is commonly referred to as "deletion by erasure". Note that compliance periods differ by region and jurisdiction but are typically between 5 and 7 years.

## Key Management
TRISA recommends that TRISA implementers use the [Key Handler package]({{% relref "data/keys" %}}) for key management.

## Throughput
Depending on the volume of Travel Rule transactions that your organization executes on a regular basis, you may wish to consider using the bidirectional streaming mode, `TransferStream`, which will support more throughput.

## Deployment
For deployment, TRISA implementers should be prepared to install their Identity Certificates (which will need to be reissued and reinstalled on an annual basis) to support routine and possibly long-running mTLS connections.

Implementers must also ensure that their existing databases are configured to support responding to routine Travel Rule information requests. Additional storage may be needed to store the results of these requests (e.g. encrypted [Secure Envelopes]()).