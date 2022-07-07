---
title: Key Handler Package
date: 2022-07-07T15:43:12-04:00
lastmod: 2022-07-07T15:43:12-04:00
description: "Describing Key Handler Package"
weight: 50
---

The Key Handler Package provides interfaces and handlers for managing public/private key pairs used for sealing and unsealing secure envelopes (often referred to as _sealing_ or _signing_ keys). TRISA nodes must handle keys in a variety of formats such as x.509 certificates on disk or marshaled data when sending keys in TRISA key exchanges. The key management package makes PKS simpler by standardizing how keys (both private and public keys) are managed, serialized, and stored.

Godoc Package Reference: [github.com/trisacrypto/trisa/pkg/trisa/keys](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys)

TRISA strongly recommends that sealing/unsealing keys be distinct from identity
certificates for improved security and to support differing retirement criteria. Some organizations may choose to use unique sealing/unsealing keys with each unique counterparty. Others may choose to issue new sealing/unsealing keys every month, deleting keys after the compliance window is over to "erase" private information (which cannot be decrypted without the keys). Either way, when multiple keys are involved, some key management system is needed. The handlers in this package ensure that keys are managed consistently for long running systems.

{{% notice note %}}
The key handler package is not intended to help or handle the symmetric keys that are used to encrypt payloads. For more information on symmetric cryptography, please see the [github.com/trisacrypto/trisa/pkg/trisa/crypto](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/crypto) package.
{{% /notice %}}


## Key Interface

The [`Key`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#Key) interface is a generic interface to either a private key pair or to a public key that has been shared in a TRISA key-exchange. The primary use of this top-level interface is serializing and deserializing keys with the marshaler interface and creating a unified mechanism to manage keys on disk.

```golang
type Key interface {
    PublicKey
    PrivateKey
    KeyMarshaler

    IsPrivate() bool
}
```

The [`PublicKey`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#PublicKey) interface provides access and management of a public key, either as part of a keypair or a stand-alone public key. Critically, this interface allows you to identify the key type using a signature-based identifier of the public key for the key management. This identifier should be added to secure envelopes to ensure the envelope cryptography can be matched with the correct keys.

Public keys are used to seal envelopes, typically using the RSA public key algorithm. When part of a private keypair, public keys should be serialized to send to counterparties during key exchange into a `trisa.api.v1beta.SigningKey` message. The `SigningKey` message provides metadata for decoding a PEM encoded PKIX public key for RSA encryption and transaction signing.

```golang
type PublicKey interface {
    KeyIdentifier

    SealingKey() (interface{}, error)

    Proto() (*api.SigningKey, error)
}
```

The [`PrivateKey`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#PrivateKey) interface provides access to the private key object that can be used to unseal an envelope, typically an RSA Private Key. While all keys managed by a TRISA node have a public key component, not all keys managed by a TRISA node will have the private component. Private keypairs belong to the node itself and TRISA recommends slightly different key management for these keys, ensuring robustness and security of storage, using a key manager such as Vault, KMS, Kubernetes Secrets, etc.

```golang
type PrivateKey interface {
    UnsealingKey() (interface{}, error)
}
```

The [`KeyMarshaler`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#KeyMarshaler) interface provides mechanisms for marshaling and unmarshaling keys either from disk or during key exchange. Key storage and management is discussed further in the next section.

```golang
type KeyMarshaler interface {
    Marshal() ([]byte, error)
    Unmarshal(data []byte) error
}
```

## Key Management

Currently, the `keys` package wraps two types of objects:

1. An [x.509 Certificate](https://en.wikipedia.org/wiki/X.509) either as a key pair or a stand alone certificate. These types of certificates and private keys are what the TRISA GDS issues to users.
2. A `trisa.api.v1beta1.SigningKey` protocol buffer message sent during a key exchange and containing _only_ public key information.

Keys can be instantiated from objects using one of the following methods:

- [`FromCertificate`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#FromCertificate): use when your key manager returns certificates directly.
- [`FromGDSLookup`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#FromGDSLookup): used to load public sealing keys from the GDS.
- [`FromProvider`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#FromProvider): providers handle compression and PKCS12 encrypted keys.
- [`FromSigningKey`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#FromSigningKey): used to load a public sealing key after a key exchange.
- [`FromX509KeyPair`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#FromX509KeyPair): used to load a key pair directly from a KMS.

Alternatively they can be parsed from raw data using the [`ParseKeyExchangedata`]() function which tries a variety of parsing techniques from PEM encoded data to raw certificate material, to protocol buffer unmarshaling.

When storing private keypairs on disk, TRISA recommends using the `Marshal` and `Unmarshal` functions of the [`Certificate`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#Certificate) object which creates PEM encoded keys. Private keypairs should be stored securely, using a dedicated key management system or encrypted disk.

To handle PKCS12 encrypted certificates sent from the GDS via email, use the TRISA trust [`Serializer`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trust#Serializer) to decrypt a trust [`Provider`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trust#Provider) which can then be used to collect the `Certificate` object. At this point the certificates are decrypted and can be marshaled and unmarshaled into secure storage.

Public [`Exchange`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trisa/keys#Exchange) keys received from counterparties should be cached in memory for a short duration and key exchanges should be conducted routinely to ensure the correct keys are used. TRISA recommends caching public keys from counterparties for at most 1 hour or for the duration of a `TransferStream` RPC. When creating secure envelopes in batch, if the counterparty is unreachable for a key exchange, you may request the sealing certificate of the counterparty from the GDS using the `Lookup` RPC.