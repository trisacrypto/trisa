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

For help marshaling and unmarshaling [IVMS101 identity payloads]({{< relref "ivms/" >}}), see the documentation about the [`ivms101` package in `trisa`](https://github.com/trisacrypto/trisa/tree/main/pkg/ivms101).

You can use the online [IVMS101 Validator](https://ivmsvalidator.com/) produced by [21 Analytics](https://www.21analytics.ch/) to ensure your message is properly structured IVMS101.

Examples of identity payloads, described both as [protocol buffer messages](https://github.com/trisacrypto/trisa/blob/main/pkg/ivms101/testdata/identity_payload.pb.json) and as [JSON](https://github.com/trisacrypto/trisa/blob/main/pkg/ivms101/testdata/identity_payload.json), are available in the [`trisa`](https://github.com/trisacrypto/trisa) reference implementation.

## Transaction Payloads

The `transaction` field in a TRISA `Payload` is a protobuf message intended to contain information to identify the associated transaction on the blockchain. It may also be used to send control flow messages and handling-specific instructions. As with the `identity` payload, a `transaction` is defined as an [`any`](https://developers.google.com/protocol-buffers/docs/proto3#any); meaning that technically, it can be *any* message type.

To ensure compatibility with fellow TRISA members and convenient message parsing, use one of the TRISA defined [generic transaction data structures](https://github.com/trisacrypto/trisa/blob/main/proto/trisa/data/generic/v1beta1/transaction.proto) described below:

### A `Transaction` Message

A `Transaction` is a generic message for TRISA transaction payloads. The goal of this payload is to provide enough information to link Travel Rule compliance information in the `identity` payload with a transaction on the blockchain or network.

```proto
message Transaction {
    string txid = 1;           // Transaction ID or hash unique to chain
                               // Used to notify beneficiary VASP of the transaction
                               // sent by the originating VASP.

    string originator = 2;     // Crypto address of originator
    string beneficiary = 3;    // Crypto address of beneficiary
    double amount = 4;         // Amount of transaction
    string network = 5;        // Chain of transaction/network ticker (e.g.: ETH, BTC)
    string timestamp = 6;      // Transaction timestamp (RFC 3339)
    string extra_json = 7;     // Extra data (JSON-formatted)

    string asset_type = 8;     // Token ticker, for identifying the token on chain.
                               // For native tokens, set to the same as network ticker
                               // (e.g.: ETH, BTC, USDT, etc)

    string tag = 9;            // Optional address memo/destination-tag
}
```

Below are some examples of how the `Transaction` message might be used in the context of different types of transactions and chains.

1. Blockchain without smart contract (e.g.: BTC)
    ```json
    {
      "txid": "05d9dc3fcbf48771c8ee9e95200877ef08e2766a780d4e44eee397633eb164d0",
      "originator": "14HmBSwec8XrcWge9Zi1ZngNia64u3Wd2v",
      "beneficiary": "14WU745djqecaJ1gmtWQGeMCFim1W5MNp3",
      "amount": 0.00206412,
      "network": "BTC",
      "timestamp": "2022-01-30T16:14:00Z",
      "extra_json": "{\"value_when_transacted\": \"USD $77.86\"}",
      "asset_type": "BTC",
      "tag": ""
   }
   ```

2. Blockchain supporting a smart contract (e.g.: ETH)
    - Native token
        ```json
        {
          "txid": "0x3e23a5165fd5c1c0f95cfc85c1419959a21e3a1a057040328fe9d3ffd7f2f991",
          "originator": "0x829bd824b016326a401d083b33d092293333a830",
          "beneficiary": "0x3d9d22647690d9b2b4d95ed6e527628746153323",
          "amount": 0.2008,
          "network": "ETH",
          "timestamp": "2022-06-30T03:16:58Z",
          "extra_json": "{\"value_when_transacted\": \"USD $218.872\"}",
          "asset_type": "ETH",
          "tag": ""
        }
        ```

    - Custom token (e.g.: ERC20)
        ```json
        {

          "txid": "0x6286b5688bcc789c2d681c01beb3e49ac870de15461cb5a90d14b8a161e84236",
          "originator": "0xea0b5f97c3843175ee67ddb237e294a9144c0a68",
          "beneficiary": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
          "amount": 255,
          "network": "ETH",
          "timestamp": "2022-06-30T03:38:34Z",
          "extra_json": "{\"value_when_transacted\": \"USD $255\"}",
          "asset_type": "USDT",
          "tag": ""
        }
        ```

3. Blockchain supporting destination tags (e.g.: XRP)
    ```json
    {
      "txid": "BFD895E1D93FB6E25A3BE38A2E62B6D88753F502B4C6E55F297981538538A2F2",
      "originator": "rKKB4S7jysrevAd1BBqbAwE7mitXY59Zsc",
      "beneficiary": "rEb8TK3gBgk5auZkwc6sHnwrGVJH8DuaLh",
      "amount": 100,
      "network": "XRP",
      "timestamp": "2022-01-29T16:14:00Z",
      "extra_json": "{\"value_when_transacted\": \"USD $32.53\"}",
      "asset_type": "XRP",
      "tag": "311041419"
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