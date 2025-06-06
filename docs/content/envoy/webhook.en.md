---
title: Webhook Guide
date: 2024-08-11T10:52:58-05:00
lastmod: 2024-08-11T10:52:58-05:00
description: "Handling incoming messages with webhooks"
weight: 95
---

{{% notice style="tip" title="Enabling the Webhook" icon="anchor" %}}
To enable the Envoy webhook, you must first specify a webhook endpoint using the `$TRISA_WEBHOOK_URL` as specified in the ["configuration values"]({{% relref "envoy/configuration.md#configuration-values" %}}). This URL should be available for HTTP POST requests from the Envoy node. If this value is not specified, Envoy will not make any webhook callbacks.
{{% /notice %}}

When Envoy receives an incoming travel rule message from a remote counterparty it must return a response to that request (for both TRP and TRISA requests). Envoy may handle the request automatically using accept/reject policies or it may simply send a pending message back as the default response. Envoy also supports handling incoming messages by using a webhook callback for another compliance system to determine handling.

When the webhook is enabled, on every incoming travel rule message Envoy will:

1. Decrypt the message (returning an error if decryption is unsuccessful)
2. Validate the message (returning an error if the message is invalid)
3. Make an http POST request to your webhook with the decrypted message as JSON
4. Validate the JSON response from the webhook
5. Encrypt the outgoing message and send back to the recipient
6. Store the secure envelopes locally

The webhook will make an POST request to your endpoint supplying the data descrbed in the webhook request in JSON format. Envoy expects a status code 200 response from your endpoint with a reply in UTF-8 encoded JSON.

A couple of important things to note:

- You must provide a valid response back to Envoy or an error will be returned to the counterparty and logged on the server.
- If the webhook is unavailable, Envoy will return service unavailable to the counterparty.
- Envoy has a request timeout of 30 seconds, however the counterparty may have a shorter duration timeout; if this is the case the request will be canceled.

## Webhook Request

The webhook POST request will contain some metadata about the incoming secure envelope, the counterparty, and will contain _one of_ an `error` or a `payload` as follows:

```json
{
    "transaction_id": "",
    "timestamp": "",
    "counterparty": {},
    "hmac_signature": "",
    "public _key_signature": "",
    "transfer_state": "",
    "error": {},
    "payload": {}
}
```

| Field | Type | Description |
|---|---|---|
| transaction_id | uuid | the id of both the Envoy transaction and the secure envelope id |
| timestamp | timestamp | the RFC3339Nano timestamp of the secure envelope |
| counterparty | object | the information about the counterparty that Envoy has in its database |
| hmac_signature | base64 | if the secure envelope has a valid hmac_signature, this is the base64 encoded signature (omitted if not set) |
| public_key_signature | string | used to identify the public keys that sealed the incoming secure envelope (omitted if not set) |
| transfer_state | string | the TRISA transfer state of the secure envelope |
| error | object | if the incoming envelope is a rejection or repair request the error will be set (omitted if payload is set) |
| payload | object | if the incoming envelope contains a travel rule information payload, this will be the decrypted info (omitted if error is set) |

Details of the `counterparty`, `error`, and `payload` objects are described below.

## Webhook Reply

The reply to the POST request must have a 200 status code and the following data serialized as UTF-8 encoded JSON. Note that either the `error` field or the `payload` field may be specified _but not both_.

```json
{
    "transaction_id": "",
    "error": {},
    "payload": {},
    "transfer_action": ""
}
```

| Field | Type | Description |
|---|---|---|
| transaction_id | uuid | the transaction_id must match the transaction_id of the request |
| error | object | if you're rejection or requesting a repair, specify the details in this field (must be omitted if payload is set) |
| payload | object | the payload to return to the counterparty (may simply be an echo of the original payload) (must be omitted if error is set) |
| transfer_action | string | the TRISA transfer action you'd like specified in the reply to the counterparty (see below for values) |


If the `transaction_id` does not match the ID of the request, then Envoy will fail and return an error to the counterparty.

The `transfer_action` can be one of "PENDING", "REVIEW", "REPAIR", "ACCEPTED", "REJECTED", or "COMPLETED". See [TRISA workflows]({{% relref "../api/workflows.md" %}}) for more about how these states determine responses to incoming messages. A basic summary is as follows:

Use `"PENDING"` if a compliance officer on your team needs to review the travel rule request, make sure you include a `pending` and `identity` in your reply.

Use `"REVIEW"` if the counterparty's compliance officer needs to review the travel rule request, make sure you include a `transaction` and `identity` in your reply. This is generally used if you've made changes to the original payload.

Use `"REPAIR"` if you need the counterparty to make a change to the payload, e.g. if it is missing fields required for your jurisdiction, make sure you've included an `error` in your reply that describes the changes that need to be made.

Use `"ACCEPTED"` if you're ready for the counterparty to make the on-chain transaction and you're satisfied with the compliance exchange. Make sure to echo back the `transaction` and `identity` in your reply.

Use `"REJECTED"` if you do not want the counterparty to proceed with the on-chain transaction; make sure to include the `error` in your reply with the rejection information.

Use `"COMPLETED"` if the compliance exchange is accepted _and_ the on-chain transaction has already been conducted and the payload contains a proper, chain-specific transaction ID in the payload. Make sure you echo back the `transaction` and `identity` in your reply.

## API Objects

This section describes nested objects in the request and reply.

### Error

Errors are mutually exclusive with payloads, meaning a response or reply to the webhook will have either an `error` field or a `payload` field but should not have both.

An error describes either a _repair_ or a _rejection_ (determined by the `retry` boolean flag). A rejection (`retry=false`) implies that the counterparty should stop the on-chain transaction as the compliance exchange has not been satisfied. A repair (`retry=true`) implies that the counterparty needs to amend the identity payload, e.g. because there are missing fields that are required for compliance in your jurisdiction.

The fields for the error are as follows:

```json
{
    "code": 0,
    "message": "",
    "retry": false
}
```

| Field | Type | Description |
|---|---|---|
| code | int32 | the error code from the TRISA api package |
| message | string | a detailed message about why the request is being rejected or what needs to be repaired |
| retry | bool | should be false if this is a rejection, true if a repair is being requested |

View the available [error codes in the TRISA API docs]({{% relref "../api/errors.md" %}}).

### Payload

Payloads are mutually exclusive with errors, meaning a response or reply to the webhook will have either a `payload` field or an `error` field but should not have both.

Additionally, payloads have mutually exclusive fields `pending` and `transaction` - a payload can have one or the other, but not both.

The payload fields are as follows:

```json
{
    "identity": {},
    "pending": {},
    "transaction": {},
    "sent_at": "",
    "received_at": ""
}
```

| Field | Type | Description |
|---|---|---|
| identity | object | an IVMS101 payload, see IVMS101 docs for more information |
| pending | object | a TRISA pending message, see TRISA docs for more information |
| transaction | object | a TRISA generic transaction, see TRISA docs for more information |
| sent_at | string | the RFC3339 encoded timestamp of when the compliance exchange was initiated |
| received_at | string | the RFC3339 encoded timestamp of when the compliance exchange was approved by the counterparty |

For more information about IVMS101, please see: [Working with IVMS101]({{% relref "../data/ivms.md" %}}).

For more information about the pending and transaction generic payloads, please see: [TRISA Data Payloads]({{% relref "../data/payloads.md" %}}).

### Counterparty

The counterparty describes what the Envoy node knows about the remote counterparty. This information may be more complete if the counterparty record was created by the TRISA directory service or stored locally on the node. If the record was created by a TRP record, this information may be fairly sparse.

```json
{
    "id": "",
    "source": "",
    "directory_id": "",
    "registered_directory": "",
    "protocol": "",
    "common_name": "",
    "endpoint": "",
    "travel_address": "",
    "name": "",
    "website": "",
    "country": "",
    "business_category": "",
    "vasp_categories": [],
    "verified_on": "",
    "ivms101": "",
    "created": "",
    "modified": ""
}
```

| Field | Type | Description |
|---|---|---|
| id | ulid | the Envoy id of the counterparty for api lookups |
| source | string | the source of the counterparty info (e.g. user, gds, x509) |
| directory_id | uuid | the id of the counterparty in the TRISA directory |
| registered_directory | string | the id of the directory the counterparty is listed in |
| protocol | string | the preferred protocol of the counterparty (trp or trisa) |
| common_name | string | the common name associated with the counterparty's certificates |
| endpoint | string | the uri endpoint of the counterparty to send travel rule requests |
| travel_address | string | a travel address identifying the generic endpoint of the counterparty |
| name | string | the display name of the counterparty |
| website | string | the website of the counterparty |
| country | string | the alpha-2 country code of the counterparty |
| business_category | string | the business category of the counterparty (e.g. private, government, etc) |
| vasp_categories | []string | the type of virtual asset services the counterparty provides (e.g. Exchange or Miner) |
| verified_on | string | the RFC3339 timestamp of when the counterparty was verified by the TRISA network |
| ivms101 | IVMS 101 | the ivms101 record of the legal person representing the counterparty |
| created | string | the RFC3339nano timestamp of when the counterparty was added to your Envoy node |
| modified | string | the RFC3339nano timestamp of when the counterparty was lasted modified on your Envoy node |

## Authentication

Webhook requests can be authenticated in two ways:

1. mTLS using non-TRISA certificates
2. HMAC signature verification with a shared secret

Either of these options or both (or neither) can be used!

### mTLS

To configure mTLS on the webhook client; you'll need to specify the following configuration values:

- `$TRISA_WEBHOOK_USE_MTLS` should be set to `"true"`.
- `$TRISA_WEBHOOK_CERTS` should point to the path of the mTLS client certificate and private key.
- `$TRISA_WEBHOOK_POOL` should point to an x509 certificate authority pool for the certs.

Note that if `$TRISA_WEBHOOK_CERTS` and `$TRISA_WEBHOOK_POOL` are unset and `$TRISA_WEBHOOK_USE_MTLS` is `"true"` then by default your TRISA certificates are used as the client certificates for the webhook call. You can obtain TRISA certificates for your webhook server by requesting them using the GDS service.

### HMAC Shared Secret

Simple, cryptographic authentication for webhook requests can performed with a shared secret using a mechanism similar to [Hawk](https://mohawk.readthedocs.io/en/latest/usage.html#sending-a-request) and [AWS4 HMAC-256](https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-auth-using-authorization-header.html). This method adds authentication variables to the HTTP request in the `Authorization` header including a cryptographically signed HMAC that can be used to verify both the client and the server.

The following is an example of the `Authorization` header value. Line breaks are added to this example for readibility:

```
Authorization: HMAC
sig=zFeIxdyVnJtpMUExK7HoL37VN4tF6sMQZPEr58MBpMQ,
nonce=wm0eHegl4lD0Uw_sSYYQCw,
headers=x-transfer-id,x-transfer-timestamp
kid=01JT4B3R5Z6AHJXV87QHPPKRBM
```

The following table describes the various components of the `Authorization` header value in the preceding example:

| Component | Description |
|---|---|
| HMAC | The authentication scheme being used; indicates that this is an Envoy HMAC-256 signed header.  |
| sig | A 32 byte HMAC signature using a SHA-256 hash and a secret key. The signature is expressed as base64 encoded characters. |
| nonce | 16 cryptographically random bytes expressed as base64 encoded characters. The nonce helps prevent reply attacks by adding random data to the signature. |
| headers | A semicolon-separated, ordered list of request headers used to compute the signature. The list includes header names only, and the header names must be in lowercase. |
| kid | The ID (usually a ULID) that identifies the secret key being used to compute the signature. |

If the webhook is configured with HMAC authentication, then it will send the signed header with every webhook request. It is up to the webhook server to verify the HMAC token (or not) and to respond with a 401 status code if the token is not verified.

Additionally, server authentication is optionally available that requires the server to add a response header, `Server-Authorization` with a reciprocal HMAC code and a _new_ nonce using the same shared key. If this option is configured, the webhook client will require a verified token before accepting a response from the webhook.

{{% notice style="orange" title="Use TLS" icon="lock" %}}
Even if HMAC shared secret authorization is used, it is strongly recommended that TLS is used to secure and encrypt the connection between the webhook client and server.
{{% /notice %}}

#### Verification

In order to verify the `Authorization` header, the following steps are required:

1. Parse the value of the `Authorization` header.

    - If the header does not contain an `HMAC ` prefix, then it is invalid; otherwise strip this prefix.
    - Split the remaining header by `,` and trim any whitespace from the front or back of the split strings.
    - Create key/value pairs by splitting each part of the header on `=` using a splitn algorithm that allows at most two parts (e.g. split on the first `=`). The key is the left side and the value is the right side.
    - Ensure that `sig`, `nonce`, `headers`, and `kid` are present in the keys and contain values; any other key value pairs can be ignored.
    - An example of parsing code in GoLang can be found [here](https://github.com/trisacrypto/envoy/blob/main/pkg/webhook/auth.go#L110-L168).

2. Create a byte array of data to be signed.

    - The beginning of the array is the base64 decoded `nonce` value.
    - Split the `headers` value by `,` and then for each header in the order specified, append the utf-8 bytes of the unmodified header value to the data array.

3. Sign the data using the key identified by the key id.

    - Use the Key ID to identify the shared secret
    - Use HMAC with SHA-256 hashing to create a signed token of the data

4. Compare the bytes of the signed data with the signature.

    - Base64 decode the `sig` value
    - Ensure the bytes of `sig` are qual to the bytes of the constructed signature.

#### Base64 Encoding

The [Base64 encoding method](https://developer.mozilla.org/en-US/docs/Glossary/Base64) used is the URL-safe base64 encoding method that replaces `+/` with `-_` and does not use `=` padding characters to ensure that it does not interfere with the HTTP headers.

#### Signing a Response

You can optionally configure the Envoy webhook client to verify the response sent back from your server with the same shared secret.

To create the authorization key:

1. Generate a random 16 byte nonce
2. Select headers and append their UTF-8 bytes to the nonce
3. Use the secret key to create an HMAC with SHA-256 hash
4. Base64 encode the HMAC sum and the nonce

Create a `Server-Authorization` header using the signature, nonce, comma-separated lower case headers, and the keyID and send back as part of the reply.

#### Configuration

In order to configure HMAC authentication you need to set the following environment variables:

1. `$TRISA_WEBHOOK_AUTH_KEY_ID`: the ULID representing the key for key management purposes; if set the auth key secret must also be set.
2. `$TRISA_WEBHOOK_AUTH_KEY_SECRET`: the _hexadecimal_ encoded 32 byte secret used for securing the HMAC token; if set the auth key id must also be set.

You can generate a key ID and secret with:

```
$ envoy hmackey
```

Optionally, to require server authentication (which is not enforced by default) you can set `$TRISA_WEBHOOK_REQUIRE_SERVER_AUTH` to `"true"`; this will cause the webhook client to verify the `Server-Authorization` header before accepting the response.

Note that all of the configuration values above require `$TRISA_WEBHOOK_URL` to be set, which enables webhook usage; otherwise the configuration is ignored.

See the ["configuration values"]({{% relref "envoy/configuration.md#configuration-values" %}}) section for more on configuring Envoy from the environment.