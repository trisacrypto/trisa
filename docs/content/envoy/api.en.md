---
title: API Guide
date: 2024-05-08T14:42:51-05:00
lastmod: 2024-05-08T14:42:51-05:00
description: "Guide to using the Envoy API"
weight: 90
---

{{% notice style="tip" title="API Credentials" icon="rocket" %}}
In order to use the API you will have to generate an API Key Client ID and Secret either from the user interface, or by using the command line tool as described in ["creating api keys"]({{% relref "deploy.md#creating-api-keys" %}}).
{{% /notice %}}

## Authentication and Authorization

All API endpoints require authorization using `Bearer` tokens in the `Authorization` header. The `Bearer` token is a time-limited JWT token that is signed by the Envoy server. The JWT token contains the valid claims that the API key has: most importantly the _permissions_ assigned to the API Key. If the endpoint requires a permission that the API key does not have then a 403 error is returned.

To obtain the JWT bearer token, you must first authenticate with the server using your client ID and secret with a `POST` request to `/v1/authenticate`.

```
POST /v1/authenticate HTTP/1.1
Host: [ENDPOINT]
Accept: application/json
Content-Type: application/json

{
    "client_id": "[CLIENT_ID]",
    "client_secret": "[CLIENT_SECRET]"
}
```

The response will contain a JWT `access_token` and `refresh_token`. The `access_token` should be used as the `Authorization: Bearer` token for all API requests:

```
GET /v1/counterparties HTTP/1.1
Host: [ENDPOINT]
Accept: application/json
Content-Type: application/json
Authorization: Bearer [ACCESS_TOKEN]
```

By default the access token expires after 1 hour and the refresh token expires after 2 hours, becoming available 15 minutes before the access token expires (though these durations can all be configured).

To reauthenticate while the `refresh_token` is still valid, `POST` the `refresh_token` to the `/v1/reauthenticate` endpoint to obtain new access and refresh tokens without having to resupply the client ID and secret.

```
POST /v1/reauthenticate HTTP/1.1
Host: [ENDPOINT]
Accept: application/json
Content-Type: application/json

{
    "refresh_token": "[REFRESH_TOKEN]"
}
```

### Permissions

The following table contains the permissions available to API keys:

| Permission | Description |
|---|---|
| `users:manage` | Can create, edit, and delete users |
| `users:view` | Can view users registered on the node |
| `apikeys:manage` | Can create apikeys and view associated secret |
| `apikeys:view` | Can view apikeys created on the node |
| `apikeys:revoke` | Can revoke apikeys and delete them |
| `counterparties:manage` | Can create, edit, and delete counterparties |
| `counterparties:view` | Can view counterparty details |
| `accounts:manage` | Can create, edit, and delete accounts and crypto addresses |
| `accounts:view` | Can view accounts and crypto addresses |
| `travelrule:manage` | Can create, accept, reject, and archive transactions and send secure envelopes |
| `travelrule:delete` | Can delete transactions and associated secure envelopes |
| `travelrule:view` | Can view travel rule transactions and secure envelopes |
| `config:manage` | Can manage the configuration of the node |
| `config:view` | Can view the configuration of the node |
| `pki:manage` | Can create and edit certificates and sealing keys |
| `pki:delete` | Can delete certificates and sealing keys |
| `pki:view` | Can view certificates and sealing keys |

## Sending Travel Rule Messages

The basic method to create and send travel rule messages to counterparties is to use the "preparing transactions" workflow described by the API documentation. Using this workflow you will:

1. Create and validate an IVMS101 payload using the `/v1/transactions/prepare` endpoint
2. Create the transaction and send the message using the `/v1/transactions/send-prepared` endpoint.

This two step workflow performs a lot of work on your behalf, saving the multiple required calls you would need to different REST objects and endpoints to perform the same activities, namely:

- Lookup or create a counterparty record
- Create IVMS101 Legal Person records for the counterparty and your VASP
- Create a Transaction record for the transfer
- Send an outgoing message, get a response, &amp; create secure envelopes for the messages.

The `/v1/transactions/prepare` endpoint takes a simplified data representation and requires only the following data:

1. Routing information to identify the counterparty
2. Identity information about the originator and beneficiary accounts
3. Transfer information to identify the transaction on the chain.

The response that is returned from this endpoint includes:

1. The Routing information, possibly updated
2. A TRISA identity payload including Originator and Beneficiary VASP records
3. A Transaction payload

If there are no validation errors, the response can then be directly posted to the `/v1/transactoins/send-prepared` endpoint in order to create and send the transfer.

### Routing

The routing information for the prepare and send-prepared endpoints determines the counterparty to send the travel rule message to, as well as the protocol to use. Different protocols require different counterparty identification mechanisms.

{{% notice style="warning" title="Backwards Compatibility" icon="gear" %}}
Note that the new v1.0.0 release of Envoy has a backwards incompatible change with the v0.30.1 prepare/send-prepared endpoints. Previously only a travel address was used for counterparty identification, but now a more complex routing object is used.
{{% /notice %}}

A routing object has the following fields:

- `protocol`: identify the protocol used to send the message to and identify the counterparty
- `counterparty_id`: the ULID of the counterparty in the database (TRISA)
- `travel_address`: the OpenVASP Travel Address of the counterparty (TRISA, TRP)
- `email`: The email address of the counterparty (Sunrise)
- `counterparty`: The legal name of the counterparty (Sunrise)

The following protocols are currently supported by Envoy:

1. TRISA: send a secure envelope to a counterparty identified by ID or travel address.
2. TRP: send a TRP request to a counterparty identified by a travel address.
3. Sunrise: send an email with a secure link to a counterparty identified by email address and the name of the counterparty.

Examples of prepare messages for each protocol are below:

#### TRISA

To send a prepare/send-prepared message using the TRISA protocol you can identify the counterparty by ID:

```json
{
    "routing": {
        "protocol": "trisa",
        "counterparty_id": "01JPMJG3ZFSX753ZAKRWRZPWVC"
    },
    "originator": {},
    "beneficiary": {},
    "transfer": {}
}
```

Or you can use the Travel Address for the TRISA node:

```json
{
    "routing": {
        "protocol": "trisa",
        "travel_address": "taLg4sBFp3cWhB9wN7qsiUF8pxo7JXtVShYkv5ix1wG2kX5y4pRiJ3TRHmeD8H67TLLm5wHyDktVw1onfDeQfESumf91mjRTMbi"
    },
    "originator": {},
    "beneficiary": {},
    "transfer": {}
}
```

Note that you must specify either `counterparty_id` or `travel_address`, not both, and `email` and `counterparty` should not be specified.

#### Sunrise

To send a sunrise message, create a counterparty with the email address of the compliance officer (using the full name, bracket email address format) and the name of the counterparty to send the message to:

```json
{
    "routing": {
        "protocol": "sunrise",
        "counterparty": "AliceCoin, LLC",
        "email": "Jane Smith <jane@alicevasp.com>"
    },
    "originator": {},
    "beneficiary": {},
    "transfer": {}
}
```

If the counterparty already exists in your Envoy database, it's easier to use only the email address so that the counterparty can be looked up successfully:

```json
{
    "routing": {
        "protocol": "sunrise",
        "email": "jane@alicevasp.com"
    },
    "originator": {},
    "beneficiary": {},
    "transfer": {}
}
```

#### TRP

To send a TRP message, the travel address is required:

```json
{
    "routing": {
        "protocol": "trp",
        "travel_address": "taLg4sBFp3cWhB9wN7qsiUF8pxo7JXtVShYkv5ix1wG2kX5y4pRiJ3TRHmeD8H67TLLm5wHyDktVw1onfDeQfESumf91mjRTMbi"
    },
    "originator": {},
    "beneficiary": {},
    "transfer": {}
}
```

### Originator and Beneficiary Persons

The `originator` and `beneficiary` fields are `Person` objects that are simplifications of the IVMS101 Natural Person object and are used to create IVMS101 natural persons for the identity payload. They both have the same exact fields.

**NOTE**: If either the originator or beneficiary account holder is a business (e.g. a Legal Person), you'll have to use raw IVMS101 with the `/v1/send-prepared` endpoint.

An example of a JSON person with all fields is as follows:

```json
{
    "crypto_address": "mvr5YZBdAuV8sgexCHL4CRkbCvTV7odT1i",
    "forename": "James",
    "surname": "Bond",
    "country_of_residence": "GB",
    "customer_id": "674907513",
    "identification": {
      "type_code": "DRLC",
      "number": "BOND9211110JA9OB",
      "country": "GB",
      "dob": "1920-11-11",
      "birth_place": "Wattenscheid, Germany"
    },
    "addresses": [
      {
        "address_type": "HOME",
        "address_lines": [
          "1 High Street",
          "Apt. 007",
          "London SE1 2QH"
        ],
        "country": "GB"
      }
    ]
  }
```

Not all of these fields are required for all jurisdictions! We recommend that you send the minimum amount of PII as required by your jursidiction and use the review/repair workflow to request any additional data required for compliance.

### Transfer

The transfer object is used to identify the transaction on the blockchain. A fully populated JSON representation is as follows:

```json
"transfer": {
    "amount": 0.0007,
    "network": "BTC",
    "asset_type": "",
    "transaction_id": "",
    "tag": ""
}
```

The `network` field should be either the SLIP-0044 short code for the virtual asset or the DTI (Digital Token Identifier) to allow the counterparty to identify the blockchain network being used. The `asset_type` and `tag` (sometimes referred to as "memo") fields are chain-specific fields and are not required in all cases. The `transaction_id` is the hash of the block or ledger record and is only available once the transaction on the chain has been completed.