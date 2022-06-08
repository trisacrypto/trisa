---
title: Data Payloads
date: 2022-06-08T14:54:48-04:00
lastmod: 2022-06-08T14:54:48-04:00
description: "Preparing and formatting data for payloads."
weight: 22
---

A `Payload` contains information to be securely exchanged for Travel Rule compliance, including (1) identity details, (2) transaction information for both counterparties to uniquely identify the transaction on the chain, and (3) timestamps for use in regulatory
non-repudiation. The payload is serialized and encrypted to be sent in a
`SecureEnvelope`. Payloads are digitally signed to ensure that they have not been tampered with in transmission.

The formal definition of a `Payload` is described on this page and available in the form of protocol buffers in the [`trisa` code repository](https://github.com/trisacrypto/trisa/tree/main/proto).

As shown below, the internal message types of the `Payload` are purposefully generic to allow flexibility with the data needs for different exchanges:

```proto
message Payload {
    google.protobuf.Any identity = 1;
    google.protobuf.Any transaction = 2;
    string sent_at = 3;
    string received_at = 4;
}
```

## Identity Payload

The `identity` field in a TRISA `Payload` is a protobuf message intended to contain compliance identity information of the Originator and the Beneficiary. It is defined as an [`any`](https://developers.google.com/protocol-buffers/docs/proto3#any); this means that technically, it can be *any* message type. However, to encourage maximum compatibility between yourself and fellow TRISA members, we strongly recommend the use of [IVMS101](https://intervasp.org).

Examples of identity payloads, described both as [protocol buffer messages](https://github.com/trisacrypto/trisa/blob/main/pkg/ivms101/testdata/identity_payload.pb.json) and as [JSON](https://github.com/trisacrypto/trisa/blob/main/pkg/ivms101/testdata/identity_payload.json) are available in the [`trisa`](https://github.com/trisacrypto/trisa) reference implementation.

You can use the online [IVMS101 Validator](https://ivmsvalidator.com/) produced by [21Analytics](https://www.21analytics.ch/) to ensure your message is properly structured IVMS101.


## Transaction Payloads

The `transaction` field in a TRISA `Payload` is a protobuf message intended to contain information to identify the associated transaction on the blockchain. It may also be used to send control flow messages and handing-specific instructions. As with the `identity` payload, a `transaction` is defined as an [`any`](https://developers.google.com/protocol-buffers/docs/proto3#any); meaning that technically, it can be *any* message type.

To ensure compatibility with fellow TRISA members and convenient message parsing, use one of the TRISA defined [generic transaction data structures](https://github.com/trisacrypto/trisa/blob/main/proto/trisa/data/generic/v1beta1/transaction.proto) described below:

### A `Transaction` Message

A `Transaction` is a generic message for TRISA transaction payloads. The goal of this payload is to provide enough information to link Travel Rule compliance information in the `identity` payload with a transaction on the blockchain or network. All fields are
optional.

```proto
message Transaction {
    string txid = 1;              // transaction ID unique to chain
    string originator = 2;        // crypto address of originator
    string beneficiary = 3;       // crypto address of beneficiary
    double amount = 4;            // amount of transaction
    string network = 5;           // chain of transaction
    string timestamp = 6;         // transaction timestamp (RFC 3339)
    string extra_json = 7;        // extra data (JSON-formatted)
    string asset_type = 8;        // asset type (for multi-asset chains)
    string tag = 9;               // optional memo/destination-tag
}
```

An [example](https://github.com/trisacrypto/trisa/blob/a2a71ed0b32b04c9859b5a9f17efae8d2d4791d8/pkg/trisa/envelope/testdata/payload/transaction.json) `Transaction` message can be found in the [`trisa`](https://github.com/trisacrypto/trisa) reference implementation.

### A `Pending` Message

A `Pending` message is a control flow message to support asynchronous TRISA transfers. Pending messages can be returned as an intermediate response during a compliance transfer if further processing is required before a response can be sent. The `Pending` message should provide information to the originator about when they can expect a response via the `reply_not_before` and `reply_not_after` timestamps. The `Pending` message should also provide collation information such as the `envelope_id` and original transaction so that the response message can be matched to the original request.

```proto
message Pending {
    string envelope_id = 1;       // TRISA envelope ID
    string received_by = 2;       // recipient or recipient VASP name
    string received_at = 3;       // when request was received (RFC3339)
    string message = 4;           // optional message for counterparty
    string reply_not_after = 5;   // when response will be returned (RFC3339)
    string reply_not_before = 6;  // response will not be sent before (RFC3339)
    string extra_json = 7;        // extra data (JSON-formatted)
    Transaction transaction = 15; // original transaction for reference
}
```

An [example](https://github.com/trisacrypto/trisa/blob/a2a71ed0b32b04c9859b5a9f17efae8d2d4791d8/pkg/trisa/envelope/testdata/payload/pending.json) `Pending` message can be found in the [`trisa`](https://github.com/trisacrypto/trisa) reference implementation.

## Timestamps

The `sent_at` and `received_at` timestamps are RFC-3339 formatted timestamps intended for use in regulatory non-repudiation.

The `sent_at` timestamp marks the time the Originator sent the first compliance message to the Beneficiary.

The `received_at`  timestamp when the Beneficiary accepted the compliance message and returned the completed payload to the Originator.