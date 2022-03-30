---
title: Secure Envelopes
date: 2022-03-03T15:36:09-06:00
lastmod: 2022-03-03T15:36:09-06:00
description: "Working with Secure Envelopes and cryptography in the TRISA protocol."
weight: 21
---

![Secure Envelopes](/img/secure_envelopes.png)


The primary data structure for a TRISA exchange is the `SecureEnvelope` -- a wrapper for compliance payload data that facilitates peer-to-peer trust in compliance information exchanges. The design of TRISA secure envelopes had the following requirements:

1. **Privacy**: only the recipient of the secure envelope should be able to open the secure envelope, even outside of the context of an RPC that is secured by transport layer encryption.
2. **Non-repudiation**: the secure envelope should cryptographically guarantee that the compliance payload is valid and has not been modified or tampered with since the original exchange.
3. **Long-term data storage**: the secure envelope should be encrypted at rest and can be "erased" by deleting its associated private keys when the compliance statue is over (usually 5-7 years).
4. **Security**: secure envelopes should prevent statistical attacks on the encrypted payload data while using public-key cryptography for ease of key management.

These requirements are essential to understanding how to successfully engage with TRISA peers, therefore understanding the `SecureEnvelope` is the first step to being able to implement the TRISA protocol.

## Working with Secure Envelopes

There are two basic workflows for secure envelopes: creating and sealing an envelope to send to a counterparty, or unsealing and parsing a received secure envelope.

### Creating a Secure Envelope

**Prerequisites**:

1. You should have constructed an appropriate TRISA `Payload` that contains an `identity` (an IVMS 101 `IdentityPayload`), a `transaction` (a TRISA generic transaction) and a `sent_at` timestamp (RFC-3339 formatted).
2. You should have the _public sealing key_ of the receipient. You can obtain this key either via the `KeyExchange` RPC or by requesting the key from the directory service.

**Steps**:

1. Create a new envelope with a uuid4 envelope ID and the current timestamp
2. Marshal the `Payload` protocol buffers into an array of bytes.
3. Generate an encryption key and encrypt the payload bytes.
4. Generate an hmac secret and sign the encrypted payload
5. Use the public sealing key of the recipient to encrypt the encryption key and hmac secret.
6. Mark the envelope as sealed and populate all required metadata.

### Opening a Secure Envelope

1. Use the `public_key_signature` to identify the private key required to decrypt the encryption key and hmac secret.
2. Use your _private sealing key_ to decrypt the encryption key and HMAC secret.
3. Use the hmac secret and hmac algorithm to verify the encrypted payload has not been tampered with by ensuring the hmac you generate is identical to the hmac on the envelope.
4. Use the encryption key and encryption algorithm to decrypt the payload.
5. Unmarshal the payload into a TRISA `Payload` object.
6. Unmarshal the `identity` and `transaction` payloads and verify that you can parse them into datastructures you can use for your compliance workflow.

### Envelope States

As you can see from the above workflows, envelopes can be in one of three states:

1. **Sealed**: the encryption key and hmac secret on the envelope are encrypted with the public key of the recipient. Only the recipient can open the envelope in this state.
2. **Unsealed**: the encryption key and hmac secret are in the clear and can be used to decrypt the payload and verify the HMAC.
3. **Clear**: the payload has been decrypted and can be unmarshaled into a TRISA `Payload` protocol buffer.

Generally speaking, when working with secure envelopes, "sealing" an envelope moves it through the following states:

```
Payload --> Clear --> Unsealed --> Sealed
```

Conversely opening an envelope moves it through the following states:

```
Sealed --> Unsealed --> Clear --> Payload
```

Maintaining envelopes in these various states can be useful to different applications. For example, an unsealed or clear envelope can be used to move data inside of an application while maintaining associated TRISA metdata. A sealed envelope can be used for long-term storage and with proper key-management, ensure erasure of the envelope simply by deleting the keys.

## The Anatomy of a Secure Envelope

A `SecureEnvelope` contains envelope metadata, cryptographic metadata, and an encrypted payload and HMAC signature. There are two types of complete envelope:

1. An error envelope containing complete envelope metadata and a TRISA error. This envelope is sent as a rejection or transfer control message back to the sender from the recipient.
2. A payload envelope that has complete envelope and cryptographic metadata, as well as an encrypted payload and HMAC signature, but without an error. Payload envelopes can be either "sealed" or "unsealed" as described above.

### Envelope Metadata

| Field        | Definition                                                                                                                                                                                                                                                                                                            |
|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `id`         | also referred to as the "envelope id" - this is a unique ID generated by the originator and must be identical on all envelopes refering to the same virtual asset transaction. The originator usually generates this ID as a [UUID4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)). |
| `timestamp`  | a nanosecond resolution RFC-3339 formatted timestamp (e.g. [`RFC3339Nano`](https://pkg.go.dev/time#RFC3339Nano)) that is used to order messages with the same ID.                                                                                                                                                     |
| `error`      | a TRISA-specific error that is intended to help facilitate compliance exchanges. TRISA [error codes](https://github.com/trisacrypto/trisa/blob/main/proto/trisa/api/v1beta1/errors.proto) are used to reject compliance data or to request a fix and retry in a follow-on secure envelope.                            |

### Cryptographic Metadata

| Field                  | Definition                                                                                                                                                                                                                                                                                                                                                                   |
|------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `encryption_key`       | the key used to encrypt the compliance payload. Keys are generated by the originator using a crypto random method. To prevent statistical attacks, a new key is generated for every transfer (e.g. at the same time as the envelope ID). If the envelope is sealed the  `encryption_key`  is encrypted using the public sealing key of the recipient.                        |
| `encryption_algorithm` | a string that describes the algorithm used to encrypt the compliance payload. This string should provide enough information for the recipient to be able to decrypt the payload. The default in the reference implementation is  `"AES256-GCM"`  which describes the use of [ AES Galois Counter Mode ]( https://datatracker.ietf.org/doc/html/rfc5288 ) with a 32 byte key. |
| `hmac_secret`          | the secret used to calculate the HMAC signature. This secret is generated by the originator using a crypto random method, and is generated for every transfer. If the envelope is sealed then the  `hmac_secret`  is encrypted using the public sealing key of the recipient.                                                                                                |
| `hmac_algorithm`       | a string that describes the algorithm used to compute the HMAC signature. The default in the reference implementation is  `"HMAC-SHA256"`  which describes the use of the [ HMAC ]( https://en.wikipedia.org/wiki/HMAC ) algorithm with a [ SHA-256 ]( https://en.wikipedia.org/wiki/SHA-2 ) secure hashing function.                                                        |
| `sealed`               | a boolean that describes the state of the envelope. If true, this means that the  `encryption_key`  and  `hmac_secret`  have been encrypted using the public sealing key of the recipient.                                                                                                                                                                                   |
| `public_key_signature` | the signature of the public key used to seal the envelope, a helper for the recipient to identify the private key required to unseal the envelope.                                                                                                                                                                                                                           |

### Payload

| Field     | Definition                                                                                    |
|-----------|-----------------------------------------------------------------------------------------------|
| `payload` | the ` Payload ` protocol buffer marshaled to bytes and encrypted using the ` encryption_key`. |
| `hmac`    | the HMAC signature computed from the encrypted payload bytes and the ` hmac_secret`.          |

The `Payload` protocol buffer has the following fields:

| Field         | Definition                                                                                                                                                                                                                                                                                                                                                                                                                |
|---------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `identity`    | a protobuf [ any ]( https://developers.google.com/protocol-buffers/docs/proto3#any ) that contains the compliance identity information of the originator and the beneficiary. Although this can be any message type, it should be an [ IVMS101 Identity Payload ]( https://intervasp.org ).                                                                                                                               |
| `transaction` | a protobuf [ any ]( https://developers.google.com/protocol-buffers/docs/proto3#any )  that contains information used to identify the associated transaction on the blockchain or to send control flow messages and handling-specific instructions. Use one of the TRISA defined [generic transaction data structures](https://github.com/trisacrypto/trisa/blob/main/proto/trisa/data/generic/v1beta1/transaction.proto). |
| `sent_at`     | The RFC-3339 formatted timestamp that the originator sent the first compliance message to the beneficiary. This timestamp is part of the compliance payload for non-repudiation purposes.                                                                                                                                                                                                                                 |
| `received_at` | The RFC-3339 formatted timestamp when the beneficiary accepted the compliance message and returned the completed payload to the originator. This timestamp is part of the compliance payload for non-repudiation purposes.                                                                                                                                                                                                |