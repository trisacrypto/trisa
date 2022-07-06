---
title: Key Handler Package
date: 2022-07-06T13:09:50-04:00
lastmod: 2022-07-06T13:09:50-04:00
description: "Describing Key Handler Package"
weight: 50
---

The Key Handler Package provides interfaces and handlers for managing public/private key pairs used for sealing and unsealing secure envelopes. This package is not intended for use with symmetric keys that are used to encrypt payloads. Instead, this key management package makes PKS simpler by standardizing how keys (both private and public keys) are managed, serialized, and stored.

TRISA strongly recommends that sealing/unsealing keys be distinct from identity
certificates for improved security and to support differing retirement criteria. Some organizations may choose to use unique sealing/unsealing keys with each unique counterparty.

## `Key` Interface

The `Key` interface is a generic interface to either a private key pair or to a public key that has been shared in a TRISA key-exchange. The primary use of this top-level interface is serializing and deserializing keys with the marshaler interface and creating a unified mechanism to manage keys on disk.

```golang
type Key interface {
    PublicKey
    PrivateKey
    KeyMarshaler

    IsPrivate() bool
}
```

`PublicKey` provides the sealing public key algorithm to identify the key type and an identifier of the public key for the key management. The public key object is used to seal an envelope, typically an RSA Public Key, and attempts to parse the keys that may have been sent in the exchange with the `SigningKey`. `SigningKey`  provides metadata for decoding a PEM encoded PKIX public key for RSA encryption and transaction signing.

```golang
type PublicKey interface {
    KeyIdentifier

    SealingKey() (interface{}, error)

    Proto() (*api.SigningKey, error)
}
```

`PrivateKey` provides the key object that can be used to unseal an envelope, typically an RSA Private Key.

```golang
type PrivateKey interface {
    UnsealingKey() (interface{}, error)
}

```

`KeyMarshaler` provides the protocol buffer marshaled data for the most compact storage through the `Marshal` method. Then, it provides the protocol buffer marshaled data and loads the sealing key through the `Unmarshal` method.

```golang
type KeyMarshaler interface {
    Marshal() ([]byte, error)
    Unmarshal(data []byte) error
}
```
