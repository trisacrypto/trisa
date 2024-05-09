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
| TRISA_MODE | string | release | Specify the mode of the API/UI server (release, debug, or testing) |
| TRISA_LOG_LEVEL | string | info | Specify the verbosity of logging (trace, debug, info, warn, error, fatal, panic) |
| TRISA_CONSOLE_LOG | bool | false | If true, logs colorized human readable output instead of json |
| TRISA_DATABASE_URL | string | sqlite3:///trisa.db | DSN containing the backend database configuration |
| TRISA_ENDPOINT | string |  | The endpoint of the TRISA node as defined by the mTLS certificates (to create travel addresses) |

### Web UI/API Configuration

These configuration values influence the behavior of the internal web UI and API.

| EnvVar | Type | Default | Description |
|---|---|---|---|
| TRISA_WEB_ENABLED | bool | true | If false, both the web UI and API are disabled |
| TRISA_WEB_BIND_ADDR | string | :8000 | The IP address and port to bind the web server on |
| TRISA_WEB_ORIGIN | string | http://localhost:8000 | The origin (url) of the web UI for creating API endpoints |
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

### TRISA Directory Configuration

The following configuration influences how the Envoy node connects to the TRISA Global Directory Service.

If you're running a TestNet node, then ensure the values point to `trisatest.net` (e.g. `api.trisatest.net:443`), if you're running a MainNet node, then ensure the values point to `vaspdirectory.net` (the default values).

| EnvVar | Type | Default | Description |
|---|---|---|---|
| TRISA_NODE_DIRECTORY_INSECURE | bool | false | If true, do not connect to the directory using TLS (only useful for local development) |
| TRISA_NODE_DIRECTORY_ENDPOINT | string | api.vaspdirectory.net:443 | The endpoint of the public GDS service |
| TRISA_NODE_DIRECTORY_MEMBERS_ENDPOINT | string | members.vaspdirectory.net:443 | The endpoint of the private members GDS service |
| TRISA_DIRECTORY_SYNC_ENABLED | bool | true | If false, then the background directory sync service will not run |
| TRISA_DIRECTORY_SYNC_INTERVAL | duration | 6h | The interval that the node will synchronize counterparties with the GDS |