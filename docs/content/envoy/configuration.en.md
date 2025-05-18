---
title: Configuration
date: 2024-05-08T12:14:45-05:00
lastmod: 2024-05-08T12:14:45-05:00
description: "Configuring Envoy"
weight: 70
---

For the latest and most up to date description of the Envoy configuration, ask Envoy directly! You can do this using the Envoy docker image as follows:

```
$ docker run trisa/envoy:latest envoy config
```

This will print out a table of the configuration options, default values, and descriptions. If you'd prefer it in list form, run:

```
$ docker run trisa/envoy:latest envoy config --list
```

## Configuration Values

Envoy is configured via the environment and for local development, also supports using `.env` files in the working directory for loading environment variables. We recommend configuring Envoy using the deployment mechanism of your choice. For example, if you're running the binary using `systemd`, then the environment should be defined in your `.service` using `Environment` or an `EnvironmentFile`. If you're using Kubernetes or Docker, then the environment variables should be added to the manifest of your deployment.

A list of the primary environment variables and their configuration are as follows:

| EnvVar | Type | Default | Description |
|---|---|---|---|
| TRISA_MAINTENANCE | bool | false | If true, the node will start in maintenance mode and will respond Unavailable to requests |
| TRISA_ORGANIZATION | string | Envoy | Specify the display name of the organization using the Envoy node for the web UI and interactive docs |
| TRISA_MODE | string | release | Specify the mode of the API/UI server (release, debug, or testing) |
| TRISA_LOG_LEVEL | string | info | Specify the verbosity of logging (trace, debug, info, warn, error, fatal, panic) |
| TRISA_CONSOLE_LOG | bool | false | If true, logs colorized human readable output instead of json |
| TRISA_DATABASE_URL | string | sqlite3:///trisa.db | DSN containing the backend database configuration |
| TRISA_SEARCH_THRESHOLD | float | 0.0 | Specify a threshold for fuzzy search from 0.0 (any match) to 1.0 (exact matches only) |
| TRISA_ENDPOINT | string |  | The endpoint of the TRISA node as defined by the mTLS certificates (to create travel addresses) |
| TRISA_TRP_ENDPOINT | string |  | If enabled, the endpoint of the TRP node as assigned by the mTLS certificates (to create travel addresses) |

### Web UI/API Configuration

These configuration values influence the behavior of the internal web UI and API.

| EnvVar | Type | Default | Description |
|---|---|---|---|
| TRISA_WEB_ENABLED | bool | true | If false, both the web UI and API are disabled |
| TRISA_WEB_API_ENABLED | bool | true | If false, the API will return unavailable when accessed |
| TRISA_WEB_UI_ENABLED | bool | true | If false, the web UI will return unavailable when accessed |
| TRISA_WEB_BIND_ADDR | string | :8000 | The IP address and port to bind the web server on |
| TRISA_WEB_ORIGIN | string | http://localhost:8000 | The origin (url) of the web UI for creating API endpoints |
| TRISA_WEB_DOCS_NAME | string |  | The display name for the API docs server in the Swagger app (by default the organization name) |
| TRISA_WEB_AUTH_KEYS | map |  | Optional static RSA key configuration for signing access and refresh tokens. Should be a comma separated map of keyID:path. |
| TRISA_WEB_AUTH_AUDIENCE | string | http://localhost:8000 | The value for the `aud` (audience) claim in JWT tokens issued by the API |
| TRISA_WEB_AUTH_ISSUER | string | http://localhost:8000 | The value for the `iss` (issuer) claim in JWT tokens issued by the API |
| TRISA_WEB_AUTH_COOKIE_DOMAIN | string | localhost | Limit cookies for the UI to the specified domain (exclude any port information) |
| TRISA_WEB_AUTH_ACCESS_TOKEN_TTL | duration | 1h | The amount of time before an access token expires |
| TRISA_WEB_AUTH_REFRESH_TOKEN_TTL | duration | 2h | The amount of time before refresh tokens expire |
| TRISA_WEB_AUTH_TOKEN_OVERLAP | duration | -15m | The amount of overlap between the access and refresh tokens, the more negative the duration the more the overlap |

### TRISA Node Configuration

Configuration values for the public facing TRISA node.

| EnvVar | Type | Default | Description |
|---|---|---|---|
| TRISA_NODE_ENABLED | bool | true | If false, the TRISA node server will not be run |
| TRISA_NODE_BIND_ADDR | string | :8100 | The ip address and port to bind the TRISA node server on |
| TRISA_NODE_POOL | path |  | The path to TRISA x509 certificate pool; this allows you to define what certificate authorities you're willing to accept using mTLS (optional) |
| TRISA_NODE_CERTS | path |  | The path to your TRISA identify certificates and private key for establishing mTLS connections to TRISA peer counterparties |
| TRISA_NODE_KEY_EXCHANGE_CACHE_TTL | duration | 24h | The duration to cache public keys exchanged with remote TRISA nodes before performing another key exchange |

### Webhook Configuration

If you would like to configure the Envoy node to send incoming travel rule requests to a webhook, you can configure those details below. For more information on the webhook and authentication, please see the ["webhook guide"]({{% relref "envoy/webhook.md" %}})

| EnvVar | Type | Default | Description |
|---|---|---|---|
| TRISA_WEBHOOK_URL | string |  | Specify a callback webhook that incoming travel rule messages will be posted to |
| TRISA_WEBHOOK_USE_MTLS | bool | false | Set to true to require the webhook client to use mTLS to authenticate to the server |
| TRISA_WEBHOOK_CERTS | string |  | Specify the path to the webhook client certificates and private key (TRISA certs used by default) |
| TRISA_WEBHOOK_POOL | string |  | Specify the path to the webhook client certificate authority pool (TRISA pool used by default) |
| TRISA_WEBHOOK_AUTH_KEY_ID | string |  | Used to identify the shared secret for HMAC authorization headers (required if secret is set) |
| TRISA_WEBHOOK_AUTH_KEY_SECRET | string |  | Specify a hexadecimal encoded 32 byte shared secret for HMAC authorization (required if key id is set) |
| TRISA_WEBHOOK_REQUIRE_SERVER_AUTH | bool | false | If true, the client will expect the webhook server to send a Server-Authorization header with HMAC token |

### TRISA Directory Configuration

The following configuration influences how the Envoy node connects to the TRISA Global Directory Service.

If you're running a TestNet node, then ensure the values point to `testnet.directory` (e.g. `api.testnet.directory:443`), if you're running a MainNet node, then ensure the values point to `trisa.directory` (the default values).

| EnvVar | Type | Default | Description |
|---|---|---|---|
| TRISA_NODE_DIRECTORY_INSECURE | bool | false | If true, do not connect to the directory using TLS (only useful for local development) |
| TRISA_NODE_DIRECTORY_ENDPOINT | string | api.trisa.directory:443 | The endpoint of the public GDS service |
| TRISA_NODE_DIRECTORY_MEMBERS_ENDPOINT | string | members.trisa.directory:443 | The endpoint of the private members GDS service |
| TRISA_DIRECTORY_SYNC_ENABLED | bool | true | If false, then the background directory sync service will not run |
| TRISA_DIRECTORY_SYNC_INTERVAL | duration | 6h | The interval that the node will synchronize counterparties with the GDS |

### Sunrise Configuration

To enable the Sunrise protocol use the following configuration and ensure that you also update the email configuration for the node to send outgoing emails.

| EnvVar | Type | Default | Description |
|---|---|---|---|
| TRISA_SUNRISE_ENABLED | bool | true | Used to disable sunrise access which will cause external sunrise pages to return a 404; both this and email need to be enabled for Sunrise |
| TRISA_SUNRISE_TRISA_WEB_ORIGIN | string |  | The URL to send sunrise requests to (by default the same as TRISA_WEB_ORIGIN) |
| TRISA_SUNRISE_INVITE_ENDPOINT | string | /sunrise/verify | The endpoint to verify an incoming Sunrise request |
| TRISA_SUNRISE_REQUIRE_OTP | true |  | If true, Sunrise verification will require an additional OTP step to access PII data |

### TRP Node Configuration

Configuration values for the publically facing TRP server.

| EnvVar | Type | Default | Description |
|---|---|---|---|
| TRISA_TRP_ENABLED | bool | true | If false, the TRP node server will not be run |
| TRISA_TRP_BIND_ADDR | string | :8200 | The ip address and port to bind the TRISA node server on |
| TRISA_TRP_IDENTITY_VASP_NAME | string |  | Specify the name of your VASP for TRP identity requests |
| TRISA_TRP_IDENTITY_LEI | string |  | Specify the LEI of your VASP to respond to a TRP identity request |
| TRISA_TRP_USE_MTLS | bool | true | If true, the TRP server will require mTLS authentication |
| TRISA_TRP_POOL | path |  | The path to TRP x509 certificate pool; this allows you to define what certificate authorities you're willing to accept using mTLS (optional) |
| TRISA_TRP_CERTS | path |  | The path to your TRP identify certificates and private key for establishing mTLS connections to TRISA peer counterparties |

### Email Configuration

Configure either SMTP _or_ SendGrid so that the Envoy node can send emails for Sunrise messages, forgot password resets, etc. If email is not enabled, the Sunrise protocol will be disabled.

| EnvVar | Type | Default | Description |
|---|---|---|---|
| TRISA_EMAIL_SENDER | string |  | The email address that messages are sent from by the Envoy node (e.g. compliance@envoy.local) |
| TRISA_EMAIL_SENDER_NAME | string |  | The name of the sender, usually the name of the VASP or compliance team |
| TRISA_EMAIL_SUPPORT_EMAIL | string |  | An email address to refer support requests to, will appear on error pages |
| TRISA_EMAIL_COMPLIANCE_EMAIL | string |  | An email address to refer compliance requests to in case an originator counterparty does not use TRISA |
| TRISA_EMAIL_TESTING | bool | false | Sets the emailer to testing mode and ensures no live emails are sent |
| TRISA_EMAIL_SMTP_HOST | string |  | If configuring SMTP, the host without the port (e.g. smtp.example.com) |
| TRISA_EMAIL_SMTP_PORT | int | 587 | The port to access the SMTP on |
| TRISA_EMAIL_SMTP_USERNAME | string |  | A username to authenticate to the SMTP server with |
| TRISA_EMAIL_SMTP_PASSWORD | string |  | A password to authenticate to the SMTP server with |
| TRISA_EMAIL_SMTP_USE_CRAM_MD5 | bool | false | Enables CRAM-MD5 auth to your SMTP server as defined in RFC 2195 instead of simple authentication |
| TRISA_EMAIL_SMTP_POOL_SIZE | int | 2 | The SMTP connection pool size for concurrent email sending |
| TRISA_EMAIL_SENDGRID_API_KEY | string |  | If configuring SendGrid, add the your API key to access the SendGrid API |

### Region Info

Envoy nodes support some provenance features when deployed in a geographically replicated fashion. If you would like to configure your node with hosting information (even just for debugging using the about page on the node), you may set the following environment variables:

| EnvVar | Type | Default | Description |
|---|---|---|---|
| REGION_INFO_ID | int32 |  | the 7 digit region identifier code |
| REGION_INFO_NAME | string |  | the name of the region |
| REGION_INFO_COUNTRY | string |  | the alpha-2 country code of the region |
| REGION_INFO_CLOUD | string |  | the cloud service provider |
| REGION_INFO_CLUSTER | string |  | the name of the cluster the node is hosted in |
