---
title: TRISA API
date: 2022-06-29T19:53:56-04:00
lastmod: 2022-08-10T16:16:40-04:00
description: "Describing the TRISA API"
weight: 10
---

The TRISA Network Protocol implementation exposes several RPCs needed to build TRISA into your organization's application, which are defined within [the protocol buffers](https://github.com/trisacrypto/trisa/tree/main/proto). All TRISA members must implement both services described by the TRISA protocol (the `TRISANetwork` service and the `TRISAHealth` service) to ensure that exchanges are conducted correctly and securely.

The protocol buffers can be compiled into the language of your choice. This section describes the protocol buffers for the `TRISANetwork` endpoint and the protocol buffer API in [Go](https://github.com/trisacrypto/trisa/tree/main/pkg/trisa/api/v1beta1).

{{% notice note %}}
You will need to download and install the protocol buffer compiler if you have not already.
{{% /notice %}}

## The `TRISANetwork` Service
The `TRISANetwork` service defines the peer-to-peer interactions between Virtual Asset Providers (VASPs) necessary to conduct compliance information exchanges.

The `TRISANetwork Service` has two primary RPCs, `Transfer` and `TransferStream`. `Transfer` and `TransferStream` allow VASPs to exchange compliance information before conducting a virtual asset transaction. `KeyExchange` allows public keys to exchange so that transaction envelopes can be encrypted and signed.

```proto
service TRISANetwork {
    rpc Transfer(SecureEnvelope) returns (SecureEnvelope) {}
    rpc TransferStream(stream SecureEnvelope) returns (stream SecureEnvelope) {}
    rpc ConfirmAddress(Address) returns (AddressConfirmation) {}
    rpc KeyExchange(SigningKey) returns (SigningKey) {}
}
```

{{% notice note %}}
The `ConfirmAddress` RPC is not currently implemented.
{{% /notice %}}

## `Transfer` and `TransferStream` RPCs

The `Transfer` and `TransferStream` RPCs conduct the information exchange before a virtual asset transaction. The RPCs enable an originating VASP to send an encrypted transaction envelope to the beneficiary VASP containing a unique ID for the transaction, the encrypted transaction bundle, and metadata associated with the transaction cipher. In response, the beneficiary can validate the transaction request, then return the beneficiary's transaction information using the same unique transaction ID.

The `Transfer` RPC is a unary RPC for simple, single transactions. The `TransferStream` RPC is a bidirectional streaming RPC for high throughput transaction workloads.

### `SecureEnvelope`

A `SecureEnvelope` is the encrypted transaction envelope that is the outer layer of the TRISA information exchange protocol and facilitates the secure storage of Know Your Client (KYC) data in a transaction. The envelope specifies a unique id to reference the transaction out-of-band (e.g., in the blockchain layer). It provides the necessary information so only the originator and the beneficiary can decrypt the transaction data. For more information about Secure Envelopes, [this section]({{% relref "data/envelopes" %}}) of the documentation further describes this primary data structure for the TRISA exchange.

A `SecureEnvelope` message contains different types of metadata. [The Anatomy of Secure Envelope]({{% relref "data/envelopes#the-anatomy-of-a-secure-envelope" %}}) section of this documentation further describes the envelope metadata, cryptographic metadata, and an encrypted payload and HMAC signature within the `SecureEnvelope`.

```proto
message SecureEnvelope {
    string id = 1;
    bytes payload = 2;
    bytes encryption_key = 3;
    string encryption_algorithm = 4;
    bytes hmac = 5;
    bytes hmac_secret = 6;
    string hmac_algorithm = 7;
    Error error = 9;
    string timestamp = 10;
    bool sealed = 11;
    string public_key_signature = 12;
```

## `ConfirmAddress` RPC

{{% notice note %}}
Address confirmation was initially described in the TRISA whitepaper as a mechanism to allow an originator VASP to establish that a beneficiary VASP has control of a crypto wallet address before sending transaction information with sensitive PII data. However, the details have not yet been defined, so the `ConfirmAddress` RPC is not currently implemented.
{{% /notice %}}

## `KeyExchange` RPC

The `KeyExchange` RPC allows VASPs to exchange public signing keys to facilitate transaction signatures if they have not already obtained them from the directory service.

### SigningKey

`SigningKey` provides metadata for decoding a PEM encoded PKIX public key for RSA encryption and transaction signing. The SigningKey is a lightweight version of the certificate information stored in the [Directory Service](https://vaspdirectory.net/).

```proto
message SigningKey {
    // x.509 metadata for reference without parsing the key
    int64 version = 1;
    bytes signature = 2;
    string signature_algorithm = 3;
    string public_key_algorithm = 4;

    // Validity information
    string not_before = 8;
    string not_after = 9;
    bool revoked = 10;

    // The serialized public key to PKIX, ASN.1 DER form.
    bytes data = 11;
}
```

## The `TRISAHealth` Service and `Status` RPC

The `TRISAHealth` service contains the `Status` RPC, which is optional but highly recommended for VASP members to implement. The Status endpoint allows TRISA members and the TRISA Directory Service to perform health checks with a VASP's TRISA Node and report the service conditions of the TRISA network. Because a down TRISA node will prevent travel rule compliant virtual asset transactions, the health service is intended to quickly identify network problems and notify members as quickly as possible.

{{% notice note %}}
The `TRISAHealth` service must also be behind mTLS so that the health check service can verify the identity certificates used for the TRISANetwork service.
{{% /notice %}}

```proto
service TRISAHealth {
    rpc Status(HealthCheck) returns (ServiceState) {}
}
```
### `HealthCheck`
`HealthCheck` specifies `attempts`, which is the number of failed health checks that proceeded the current check, and  `last_checked`, which is the timestamp of the last health check, successful or otherwise.

```proto
message HealthCheck {
    uint32 attempts = 1;
    string last_checked_at = 2;
}
```

### `ServiceState`
`ServiceState` returns the `status`, which is the Current service status as defined by the receiving system. The system must respond with the closest matching status in a best-effort fashion. Alerts will be triggered on service status changes if the system does not respond and the previous system state was not unknown. `not_before` and `not_after` are also returned; they suggest to the directory service when to recheck the health status.

```proto
message ServiceState {
    enum Status {
        UNKNOWN = 0;
        HEALTHY = 1;
        UNHEALTHY = 2;
        DANGER = 3;
        OFFLINE = 4;
        MAINTENANCE = 5;
    }

    Status status = 1;

    // When to check the health status again.
    string not_before = 2;
    string not_after = 3;
}
```

## TRISA API in Go

The implementation of the TRISA Protocol Buffers in Go is compiled using [`protoc`](https://grpc.io/docs/protoc-installation/) when `go generate ./...` is executed in the root of the repository. The [compiled files](https://github.com/trisacrypto/trisa/tree/main/pkg/trisa/api/v1beta1) in the TRISA repository contain the TRISA Network Protocol implemented in Go.

The `TRISANetworkServer` is the server API for `TRISANetwork`, while the `TRISANetworkClient` is the client API for the `TRISANetwork` service. Both contain the `Transfer`, `TransferStream`, `ConfirmAddress`, and  `KeyExchange` methods described above as [RPCs]({{% relref "api/api#the-trisanetwork-service" %}}) under the `TRISANetwork` service.

```golang
type TRISANetworkClient interface {
	Transfer(ctx context.Context, in *SecureEnvelope, opts ...grpc.CallOption) (*SecureEnvelope, error)
	TransferStream(ctx context.Context, opts ...grpc.CallOption) (TRISANetwork_TransferStreamClient, error)
	ConfirmAddress(ctx context.Context, in *Address, opts ...grpc.CallOption) (*AddressConfirmation, error)
	KeyExchange(ctx context.Context, in *SigningKey, opts ...grpc.CallOption) (*SigningKey, error)
}

type TRISANetworkServer interface {
	Transfer(context.Context, *SecureEnvelope) (*SecureEnvelope, error)
	TransferStream(TRISANetwork_TransferStreamServer) error
	ConfirmAddress(context.Context, *Address) (*AddressConfirmation, error)
	KeyExchange(context.Context, *SigningKey) (*SigningKey, error)
}
```
For further information, a [reference implementation](https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go) of the TRISA Network protocol is available in Go in the [TRISA TestNet Repository](https://github.com/trisacrypto/testnet),