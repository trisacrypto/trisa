---
title: Robot VASPs
date: 2021-06-14T12:50:23-05:00
lastmod: 2021-06-14T12:50:23-05:00
description: "Working with rVASPs for TestNet integration"
weight: 10
---

The TestNet hosts three convenience "robot VASPs" (rVASPs) services to facilitate integration and testing with the TRISA TestNet. These services are as follows:

- Alice (`api.alice.vaspbot.net:443`): the primary integration rVASP used to trigger and receive TRISA messages.
- Bob (`api.bob.vaspbot.net:443`): a demo rVASP to view exchanges with Alice.
- Evil (`api.evil.vaspbot.net:443`): an "evil" rVASP that is not a TRISA TestNet member, used to test non-authenticated interactions.

Note: the rVASPs are currently primarily configured for demos and work has begun on making them more robust for integration purposes; please check in with this documentation regularly for any changes. If you notice any bugs in the rVASP code or behavior, [please open an issue](https://github.com/trisacrypto/testnet/issues).

## Getting Started with rVASPs

There are two ways that you can use the rVASPs to develop your TRISA service:

1. You can trigger the rVASP to send a TRISA exchange message to your service.
2. You can send a TRISA message to the rVASP with a valid (or invalid) rVASP customer.

The rVASPs have a built-in database of fake customers with fake wallet addresses. Their response to TRISA messages or to a triggered transfer requires the originator/beneficiary customer to be valid for the rVASP. E.g. if the customer wallet address is a valid Alice address and Alice is the beneficiary, the rVASP will respond with the customer's fake KYC data but if not, it will return an TRISA error code.

The following table of "customers" for Alice, Bob, and Evil can be used as a reference for interacting with each rVASP:

| VASP                  | "Crypto Wallet"                    | Email                 | Originator Policy | Beneficiary Policy |
|-----------------------|------------------------------------|-----------------------|-------------------|--------------------|
| api.bob.vaspbot.net   | **18nxAxBktHZDrMoJ3N2fk9imLX8xNnYbNh** | robert@bobvasp.co.uk  | SendPartial       | SyncRepair         |
| api.bob.vaspbot.net   | 1LgtLYkpaXhHDu1Ngh7x9fcBs5KuThbSzw | george@bobvasp.co.uk  | SendFull          | SyncRequire        |
| api.bob.vaspbot.net   | 14WU745djqecaJ1gmtWQGeMCFim1W5MNp3 | larry@bobvasp.co.uk   | SendFull          | AsyncRepair        |
| api.bob.vaspbot.net   | _1Hzej6a2VG7C8iCAD5DAdN72cZH5THSMt9_ | fred@bobvasp.co.uk    | SendError     | AsyncReject        |
| api.alice.vaspbot.net | **1ASkqdo1hvydosVRvRv2j6eNnWpWLHucMX** | mary@alicevasp.us     | SendPartial       | SyncRepair         |
| api.alice.vaspbot.net | 1MRCxvEpBoY8qajrmNTSrcfXSZ2wsrGeha | alice@alicevasp.us    | SendFull          | SyncRequire        |
| api.alice.vaspbot.net | 14HmBSwec8XrcWge9Zi1ZngNia64u3Wd2v | jane@alicevasp.us     | SendFull          | AsyncRepair        |
| api.alice.vaspbot.net | _19nFejdNSUhzkAAdwAvP3wc53o8dL326QQ_ | sarah@alicevasp.us    | SendError     | AsyncReject        |
| api.evil.vaspbot.net  | **1PFTsUQrRqvmFkJunfuQbSC2k9p4RfxYLF** | voldemort@evilvasp.gg | SendPartial       | SyncRepair         |
| api.evil.vaspbot.net  | 172n89jLjXKmFJni1vwV5EzxKRXuAAoxUz | launderer@evilvasp.gg | SendFull          | SyncRequire        |
| api.evil.vaspbot.net  | 182kF4mb5SW4KGEvBSbyXTpDWy8rK1Dpu  | badnews@evilvasp.gg   | SendFull          | AsyncRepair        |
| api.evil.vaspbot.net  | _1AsF1fMSaXPzz3dkBPyq81wrPQUKtT2tiz_ | gambler@evilvasp.gg   | SendError     | AsyncReject        |

Note that all rVASP data was generated using a Faker tool to produce realistic/consistent test data and fixtures and is completely fictional. For example, the records for Alice VASP (a fake US company) are primarily in North America, etc.

If you're a Traveler customer, the bold addresses above have some attribution data associated with them, and they're a good candidate for Traveler-based rVASP interactions.

In order to support multiple behaviors at once, such as synchronous and asynchronous responses, each wallet address is configured with an originator policy and a beneficiary policy. The originator policy defines how the rVASP will behave when *initiating* transfers to a remote VASP. The beneficiary policy defines how the rVASP will behave when *responding* to transfers from a remote VASP. These policies are described in depth below.

### Preliminaries

This documentation assumes that you have a service that is running the latest `TRISANetwork` service and that it has been registered in the TRISA TestNet and correctly has TestNet certificates installed. See [ TRISA Integration Overview]({{< ref "getting-started/_index.md" >}}) for more information. **WARNING**: the rVASPs do not participate in the TRISA production network, they will only respond to verified TRISA TestNet mTLS connections.

To interact with the rVASP API, you may either:

1. Use the `rvasp` CLI tool
2. Use the rVASP protocol buffers and interact with the API directly

To install the CLI tool, either download the `rvasp` executable for the appropriate architecture at the [TestNet releases page](https://github.com/trisacrypto/testnet/releases), clone [the TestNet repository](https://github.com/trisacrypto/testnet/) and build the `cmd/rvasp` binary or install with `go install` as follows:

```
$ go install github.com/trisacrypto/testnet/cmd/rvasp@latest
```

To use the [rVASP protocol buffers](https://github.com/trisacrypto/testnet/tree/main/proto/rvasp/v1), clone or download them from the TestNet repository then compile them to your preferred language using `protoc`.

### Triggering an rVASP to send a message

The rVASP admin endpoints are used to interact with the rVASP directly for development and integration purposes. Note that this endpoint is different than the TRISA endpoint, which was described above.

- Alice: `admin.alice.vaspbot.net:443`
- Bob: `admin.bob.vaspbot.net:443`
- Evil: `admin.evil.vaspbot.net:443`

To use the command line tool to trigger a message, run the following command:

```
$ rvasp transfer -e admin.alice.vaspbot.net:443 \
        -a mary@alicevasp.us \
        -d 0.3 \
        -B trisa.example.com \
        -b cryptowalletaddress \
        -E
```

This message sends the Alice rVASP a message using the `-e` or `--endpoint` flag, and specifies that the originating account should be "mary@alicevasp.us" using the `-a` or `--account` flag. The originating account is used to determine what IVMS101 data to send to the beneficiary. The `-d` or `--amount` flag specifies the amount of "AliceCoin" to send. Finally, the `-b` or `--beneficiary` flag specifies the crypto wallet address of the beneficiary you intend as the recipient.

The next two parts are critical. The `-E` or `--external-demo` flag tells the rVASP to trigger a request to your service rather than to perform a demo exchange with another rVASP. This flag is required! Finally, the `-B` or `--beneficiary-vasp` flag specifies where the rVASP will send the request. This field should be able to be looked up in the TRISA TestNet directory service; e.g. it should be your common name or the name of your VASP if it is searchable.

Note that you can set the `$RVASP_ADDR` and `$RVASP_CLIENT_ACCOUNT` environment variables to specify the `-e` and `-a` flags respectively.

To use the protocol buffers directly, use the `TRISAIntegration` service `Transfer` RPC with the following `TransferRequest`:

```json
{
    "account": "mary@alicevasp.us",
    "amount": 0.3,
    "beneficiary": "cryptowalletaddress",
    "beneficiary_vasp": "trisa.example.com",
    "check_beneficiary": false,
    "external_demo": true
}
```

These values have the same exact specification as the ones in the command line program.

#### Responding to a TRISA message from an rVASP

Envelopes received from rVASPs can be decrypted using the private unsealing key paired with the public key you exchanged with the rVASP. The payload that you receive from the rVASP will be determined by the originator policy associated with the originator address in the transaction payload.

It is then up to your TRISA node to determine how to handle the payload. Your options are:

1. Reject: return an error in the secure envelope without a payload
2. Pending: return a `trisa.data.generic.v1beta1.Pending` in the transaction payload
3. Accept: ensure the identity payload is complete and return payload with `received_at`

The rVASP handles each type of response appropriately. If a reject message is returned, the rVASP fails the transaction; if accept it "executes" the transaction.

The pending message initiates an asynchronous transaction. The transaction is placed into an "await" state until the rVASP receives a follow-up reject or accept response with the same envelope id.

#### Originator Policies

An _originator policy_ determines how the rVASP constructs an outgoing payload and handles the response from the transfer. The originator policy is loaded when the `rvasp` command is used to instruct the rVASP to send an outgoing message either to another rVASP or to your own TRISA node. Originator policies are mapped to wallet addresses, and the rVASP identifies the originator address using the wallet address in the transaction payload. A description of the originator policies follow.

##### SendPartial

For the `SendPartial` policy, the rVASP sends an envelope containing a transaction payload and a partial identity payload containing a full originator identity and a partial beneficiary identity. For acceptance, the recipient must complete the identity payload by filling in the beneficiary identity.

##### SendFull

For the `SendFull` policy, the rVASP sends an envelope containing a transaction payload and a complete identity payload. For acceptance, the recipient must reseal and echo the envelope back with a `received_at` timestamp.

##### SendError

For the `SendError` policy, the rVASP will simply send an error envelope with the `ComplianceCheckFail` TRISA error code. This is useful for simluating asynchronous rejections or cancellations from the rVASP.

### Sending a TRISA message to an rVASP

The rVASP expects a `trisa.data.generic.v1beta1.Transaction` as the transaction payload and an `ivms101.IdentityPayload` as the identity payload. The transaction payload must contain a beneficiary address that matches the rVASP beneficiaries from the table above; e.g. use:

```json
{
    "originator": "anydatahereisfine",
    "beneficiary": "1MRCxvEpBoY8qajrmNTSrcfXSZ2wsrGeha",
    "amount": 0.3,
    "network": "TestNet"
}
```

You may specify any `originator` "wallet address" string and the `network` and `asset_type` fields are ignored if you would like to simulate other transaction information.

The identity payload must also be not null and must be valid IVMS101, but may be partially or completely filled in depending on the corresponding beneficiary policy.

- Partial Identity: Contains only the originator and originating vasp identities
- Complete Identity: Contains the originator, benefciiary, and originating and beneficiary vasp identities.

Create a sealed envelope either using the directory service or direct key exchange to fetch the rVASP RSA public keys and using `AES256-GCM` and `HMAC-SHA256` as the envelope cryptography. Then use the `TRISANetwork` service `Transfer` RPC to send the sealed envelope to the rVASP.

See [Secure Envelopes]({{< ref "secure-envelopes" >}}) for more on how to compose a valid secure envelope for transfers and the [TRISA CLI]({{< ref "testnet/trisa-cli" >}}) for more on using a command line application for sending transfers.

#### Beneficiary Policies

The _beneficiary policy_ determines how an rVASP responds to an incoming TRISA transfer message to its TRISA endpoint. Beneficiary policies are mapped to wallet addresses, and the rVASP identifies the beneficiary address using the wallet address in the transaction payload. A description of the beneficiary policies follow.

##### SyncRepair

For the `SyncRepair` policy, the identity payload does not have to include the beneficiary identity, although it must be not null. The rVASP will respond synchronously by sending an accept response containing a `received_at` timestamp and the complete beneficiary identity.

##### SyncRequire

For the `SyncRequire` policy, the identity payload must contain a complete beneficiary identity. The rVASP will respond synchronously by sending an accept response containing a `received_at` timestamp in the payload. If the beneficiary information is not complete or incorrect, the rVASP will respond with a rejection error.

##### AsyncRepair

For the `AsyncRepair` policy, the identity payload does not have to include the beneficiary identity. The rVASP will respond with a `trisa.data.generic.v1beta1.Pending` message containing `reply_not_before` and `reply_not_after` timestamps which specifies the beginning of an asynchronous transaction. In order to continue the transaction handshake, you should be ready to receive a `Transfer` RPC request from the rVASP within the time window containing a `SecureEvelope` with the same `Id`. In order to continue the transaction, you must respond with a resealed envelope containing a `received_at` timestamp in the payload, and then send a new `Transfer` request to the rVASP containing any final transaction details (`txid`, etc.). The rVASP will respond with another pending message which will initiate a final asynchronous handshake. Once the final `Transfer` request is received from the rVASP, the envelope should be resealed and echoed again to complete the transaction. The entire `AsyncRepair` workflow between two rVASPs is displayed below, where Alice is acting as the originator and Bob is acting as the beneficiary.

{{< mermaid >}}
sequenceDiagram
autonumber
    Alice->>Bob: Transfer() Partial Identity Info + Partial Transaction
    activate Bob
    Bob-->>Alice: Pending (5-10 min)
    activate Alice
    Bob->>Alice: Transfer() Full Identity Info + received_at
    deactivate Bob
    Alice->>Bob: Echo Payload
    deactivate Alice
    Alice->>Bob: Transfer() Full Identity Info + Full Transaction
    activate Bob
    Bob-->>Alice: Pending (5-10 min)
    activate Alice
    Bob->>Alice: Transfer() Full Identity Info + Full Transaction + received_at
    deactivate Bob
    Alice->>Bob: Echo Payload
    deactivate Alice
{{< /mermaid >}}

##### AsyncReject

The `AsyncReject` policy simulates a synchronous rejected transaction. In this policy, the identity payload does not have to include the beneficiary identity, and the rVASP will respond with a pending message as in the `AsyncRepair` policy. The difference is that the rVASP will send an error envelope within the time window containing a TRISA rejection, indicating that the transaction has been rejected.