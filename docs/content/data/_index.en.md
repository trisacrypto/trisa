---
title: Working with TRISA Data
date: 2022-08-10T16:46:31-04:00
lastmod: 2022-08-10T16:46:31-04:00
description: "Working with Data in TRISA"
weight: 20
---

This section of the documentation contains resources for developers who are working with TRISA data, such as:

- **IVMS101**: TRISA uses the IVMS101 standard to describe participants in cryptographic transactions. Learn more in our documentation about working with [IVMS101]({{< ref "/data/ivms" >}}).

- **SecureEnvelopes**: The primary data structure for a TRISA exchange is the `SecureEnvelope`, a wrapper for compliance payload data that facilitates peer-to-peer trust in compliance information exchanges. Learn more in our documentation about creating and parsing [Secure Envelopes]({{< ref "/data/envelopes" >}}).

- **Data Payloads**: A TRISA `Payload` contains information to be securely exchanged for Travel Rule compliance. The payload is serialized and encrypted to be sent in a `SecureEnvelope`. Learn more in our documentation about different types of [Payloads]({{< ref "/data/payloads" >}}) in TRISA.

- **Signing and Sealing Keys**: Your TRISA node will need to handle keys in a variety of formats, such as x.509 certificates on disk or marshaled data when sending keys in TRISA key exchanges. The Key Handler package provides helpful utilities for managing public/private key pairs used for sealing and unsealing `SecureEnvelopes`. Learn more in our documentation about the [Key Handler package]({{< ref "/data/keys" >}}).