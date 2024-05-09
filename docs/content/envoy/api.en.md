---
title: API Guide
date: 2024-05-08T14:42:51-05:00
lastmod: 2024-05-08T14:42:51-05:00
description: "Guide to using the Envoy API"
weight: 90
---

{{% notice style="tip" title="API Credentials" icon="rocket" %}}
In order to use the API you will have to generate an API Key Client ID and Secret either from the user interface, or by using the command line tool as described in ["creating api keys"]({{< relref "deploy.md#creating-api-keys" >}}).
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