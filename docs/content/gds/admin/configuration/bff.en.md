---
title: "BFF"
date: 2022-12-22T11:52:54-06:00
lastmod: 2022-12-22T11:52:54-06:00
description: "Configuring the BFF for vaspdirectory.net"
weight: 20
---

The BFF (backend for frontend) is the backend API that powers the GDS UI at [vaspdirectory.net](https://vaspdirectory.net). The GDS UI is intended to give users a single access portal to both the MainNet and TestNet as well as to support non-GDS features like collaborators and TRISA Service Providers. Because of this it sits in the middle of multiple services, including both the MainNet and TestNet GDS services and multiple data sources. Its environment variables are all prefixed with the `GDS_BFF_` tag. The primary configuration is as follows:

| EnvVar                | Type     | Default               | Description                                                                                                      |
|-----------------------|----------|-----------------------|------------------------------------------------------------------------------------------------------------------|
| GDS_BFF_MAINTENANCE   | bool     | false                 | Sets the server to maintenance mode, which will respond to requests with Unavailable except for status requests. |
| GDS_BFF_BIND_ADDR     | string   | :4437                 | The IP address and port to bind the BFF http server on.                                                          |
| GDS_BFF_MODE          | string   | release               | Sets the Gin mode, one of debug, release, or test.                                                               |
| GDS_BFF_LOG_LEVEL     | string   | info                  | The verbosity of logging, one of trace, debug, info, warn, error, fatal, or panic.                               |
| GDS_BFF_CONSOLE_LOG   | bool     | false                 | If true will print human readable logs instead of JSON logs for machine consumption.                             |
| GDS_BFF_ALLOW_ORIGINS | []string | http://localhost:3000 | A list of allowed origins for the CORS middleware to accept.                                                     |
| GDS_BFF_REGISTER_URL  | string   |                       | The base URL to direct users to for registration signup (no trailing slash) - used in email templates.           |
| GDS_BFF_LOGIN_URL     | string   |                       | The base URL to direct users to for login (no trailing slash) - used in email templates.                         |
| GDS_BFF_COOKIE_DOMAIN | string   |                       | The domain to set secure cookies for (particularly for CSRF and authentication).                                  |
| GDS_BFF_SERVE_DOCS    | bool     | false                 | If true, OpenAPI documentation is compiled and served alongside the BFF API.                                     |

### User Cache Config

The BFF interacts with Auth0 to fetch data about users. To reduce the number of Auth0 network lookups an expiring LRU cache is used to store user information for a fixed amount of time while bounding the amount of space used by the cache. The configuration for the user cache is as follows:

| EnvVar                        | Type     | Default | Description                                                        |
|-------------------------------|----------|---------|--------------------------------------------------------------------|
| GDS_BFF_USER_CACHE_ENABLED    | bool     | false   | Enable user caching to reduce lookups with the Auth0 API.          |
| GDS_BFF_USER_CACHE_SIZE       | uint     | 16384   | The size in bytes to limit the LRU cache to.                       |
| GDS_BFF_USER_CACHE_EXPIRATION | duration | 8h      | How long to keep records in the cache before forcing a new lookup. |

### Auth0 Config

The BFF uses Auth0 for authentication and authorization and must connect to the Auth0 Management API in order to manage users. The Auth0 client is configured as follows:

| EnvVar                       | Type     | Default | Description                                                                                                                                                                        |
|------------------------------|----------|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| GDS_BFF_AUTH0_DOMAIN         | string   |         | The tenant domain provided by the Auth0 application or API (domain only, no scheme or path).                                                                                        |
| GDS_BFF_AUTH0_ISSUER         | string   |         | Set to the custom domain if enabled in Auth0 (ensuring the trailing slash is set if required by the Auth0 configuration) - this will confirm the issuer from the Auth0 JWT tokens. |
| GDS_BFF_AUTH0_AUDIENCE       | string   |         | The audience to verify for the Auth0 API configuration (usually the unique name of the API).                                                                                       |
| GDS_BFF_AUTH0_PROVIDER_CACHE | duration | 5m      | Configures the JWKS caching provider to fetch public keys from Auth0 for JWT token validation.                                                                                     |
| GDS_BFF_AUTH0_CLIENT_ID      | string   |         | The Client ID for the management API specified by Auth0.                                                                                                                           |
| GDS_BFF_AUTH0_CLIENT_SECRET  | string   |         | The Client Secret for the management API specified by Auth0.                                                                                                                       |
| GDS_BFF_AUTH0_TESTING        | bool     | false   | If true a mock authenticator is used for testing purposes.                                                                                                                         |

### Network Configuration

The network configuration enables the BFF to connect to the GDS Database, Directory API, and Members API for both the MainNet and TestNet networks. The required configuration and configuration values is the same for both networks but the environment variables are prefixed with `GDS_BFF_TESTNET_` and `GDS_BFF_MAINNET_` respectively.

The configuration without the prefixes is specified below, at the end of each section, we will provide an exhaustive list of environment variables that are required to fully configure the BFF for connecting to both the TestNet and MainNet GDS services.

#### GDS Database Connection

The GDS Database connection is a store connection that is similar to the BFF Database connection described in the following section. Configuration for the TestNet and MainNet database connections is as follows:

| EnvVar                   | Type   | Default | Description                                                                          |
|--------------------------|--------|---------|--------------------------------------------------------------------------------------|
| DATABASE_URL             | string |         | Required, the DSN to connect to the database on (see below for details).              |
| DATABASE_REINDEX_ON_BOOT | bool   | false   | When the server starts, instead of loading indexes from disk, recreate and save them. |
| DATABASE_INSECURE        | bool   | false   | If set do not connect to the TrtlDB with mTLS authentication.                        |
| DATABASE_CERT_PATH       | string |         | The path to the mTLS client-side certs for database auth.                             |
| DATABASE_POOL_PATH       | string |         | The path to the mTLS public cert pools to accept server connections.                  |

Note that only a `trtl://` database url should be used for network database connections and that `DATABASE_REINDEX_ON_BOOT` should _always_ be `false` for the BFF.

List of environment variables:

- `GDS_BFF_TESTNET_DATABASE_URL`
- `GDS_BFF_TESTNET_DATABASE_INSECURE`
- `GDS_BFF_TESTNET_DATABASE_CERT_PATH`
- `GDS_BFF_TESTNET_DATABASE_POOL_PATH`
- `GDS_BFF_MAINNET_DATABASE_URL`
- `GDS_BFF_MAINNET_DATABASE_INSECURE`
- `GDS_BFF_MAINNET_DATABASE_CERT_PATH`
- `GDS_BFF_MAINNET_DATABASE_POOL_PATH`

#### GDS Directory API Configuration

The BFF connects to the GDS TRISA Directory API to perform operations like registration submission, contact verification, and verification status lookups. The configuration for the directory API clients is as follows:

| EnvVar             | Type     | Default | Description                                                         |
|--------------------|----------|---------|---------------------------------------------------------------------|
| DIRECTORY_INSECURE | bool     | false   | If false does not connect to the directory API using TLS.            |
| DIRECTORY_ENDPOINT | string   |         | The endpoint (host:port) to connect to the directory API on.         |
| DIRECTORY_TIMEOUT  | duration | 10s     | The connection timeout for directory API request and dial contexts.  |

List of environment variables:

- `GDS_BFF_TESTNET_DIRECTORY_INSECURE`
- `GDS_BFF_TESTNET_DIRECTORY_ENDPOINT`
- `GDS_BFF_TESTNET_DIRECTORY_TIMEOUT`
- `GDS_BFF_MAINNET_DIRECTORY_INSECURE`
- `GDS_BFF_MAINNET_DIRECTORY_ENDPOINT`
- `GDS_BFF_MAINNET_DIRECTORY_TIMEOUT`

#### GDS Members API Configuration

The BFF connects to the secure GDS Members API to give logged in users access to the complete directory including listing verified members and updating their member record. The configuration for the members API client is as follows:

| EnvVar                 | Type     | Default | Description                                                         |
|------------------------|----------|---------|---------------------------------------------------------------------|
| MEMBERS_ENDPOINT       | string   |         | The endpoint (host:port) to connect to the members API on.           |
| MEMBERS_TIMEOUT        | duration | 10s     | The connection timeout for members API request and dial contexts.    |
| MEMBERS_MTLS_INSECURE  | bool     | false   | If false does not connect to the members API using mTLS.             |
| MEMBERS_MTLS_CERT_PATH | string   |         | The path to the mTLS client-side certs for members API auth.         |
| MEMBERS_MTLS_POOL_PATH | string   |         | The path to the mTLS public cert pools to accept server connections. |

List of environment variables:

- `GDS_BFF_TESTNET_MEMBERS_ENDPOINT`
- `GDS_BFF_TESTNET_MEMBERS_TIMEOUT`
- `GDS_BFF_TESTNET_MEMBERS_MTLS_INSECURE`
- `GDS_BFF_TESTNET_MEMBERS_MTLS_CERT_PATH`
- `GDS_BFF_TESTNET_MEMBERS_MTLS_POOL_PATH`
- `GDS_BFF_MAINNET_MEMBERS_ENDPOINT`
- `GDS_BFF_MAINNET_MEMBERS_TIMEOUT`
- `GDS_BFF_MAINNET_MEMBERS_MTLS_INSECURE`
- `GDS_BFF_MAINNET_MEMBERS_MTLS_CERT_PATH`
- `GDS_BFF_MAINNET_MEMBERS_MTLS_POOL_PATH`

### Database

The BFF makes use of a globally replicated key-value store for persistence of BFF-specific data structures such as organizations and announcements. By default it uses the TrtlDB for this, but for testing or smaller deployments it can use a local LevelDB database instead.

Generally speaking, the BFF uses the same TrtlDB as the GDS MainNet instance, and we've ensured there are no namespace conflicts to prevent this. The BFF primary store can also be independent of the GDS stores if necessary.

The primary database store is configured as follows:

| EnvVar                       | Type   | Default | Description                                                                          |
|------------------------------|--------|---------|--------------------------------------------------------------------------------------|
| GDS_BFF_DATABASE_URL             | string |         | Required, the DSN to connect to the database on (see below for details).              |
| GDS_BFF_DATABASE_REINDEX_ON_BOOT | bool   | false   | When the server starts, instead of loading indexes from disk, recreate and save them. |
| GDS_BFF_DATABASE_INSECURE        | bool   | false   | If set do not connect to the TrtlDB with mTLS authentication.                         |
| GDS_BFF_DATABASE_CERT_PATH       | string |         | The path to the mTLS client-side certs for database auth.                             |
| GDS_BFF_DATABASE_POOL_PATH       | string |         | The path to the mTLS public cert pools to accept server connections.                  |

The `GDS_BFF_DATABASE_URL` is a standard DSN with a scheme, host, path, and query parameters. In the case of the BFF the scheme can be either `trtl://` or `leveldb://`.

When connecting to a TrtlDB, the host and port need to be specified with a trailing slash, e.g. `trtl://localhost:4436/`. Connecting via mTLS is only relevant when connecting to a TrtlDB. If `GDS_BFF_DATABASE_INSECURE` is `false`, then the `GDS_BFF_DATABASE_CERT_PATH` and `GDS_BFF_DATABASE_POOL_PATH` are required.

When connecting to a LevelDB, the path to the directory on disk where the leveldb should be stored must be passed to the DSN. To specify a relative path, use three slashes: `leveldb:///relpath/to/db`, to specify an absolute path use four: `leveldb:////abspath/to/db`.

### Email and SendGrid

The BFF uses [SendGrid](https://sendgrid.com/) to send email notifications and to enable communication workflows with users and admins. Configure BFF to use SendGrid as follows:

| EnvVar                | Type   | Default                                           | Description                                                              |
|-----------------------|--------|---------------------------------------------------|--------------------------------------------------------------------------|
| GDS_BFF_SERVICE_EMAIL | string | TRISA Directory Service <admin@vaspdirectory.net> | The email address used as the sender for all emails from the BFF system. |
| SENDGRID_API_KEY      | string |                                                   | API Key to authenticate to SendGrid with.                                |
| GDS_BFF_EMAIL_TESTING | bool   | false                                             | Use email in testing mode rather than send live emails.                   |
| GDS_BFF_EMAIL_STORAGE | string | ""                                                | Directory to store test emails for "mark one eyeball" review.             |

SendGrid is considered **enabled** if the SendGrid API Key is set. The service and admin email addresses are required if SendGrid is enabled.

### Sentry

The BFF uses [Sentry](https://sentry.io/) to assist with error monitoring and performance tracing. Configure BFF to use Sentry as follows:

| EnvVar                       | Type    | Default     | Description                                                                                       |
|------------------------------|---------|-------------|---------------------------------------------------------------------------------------------------|
| GDS_BFF_SENTRY_DSN               | string  |             | The DSN for the Sentry project. If not set then Sentry is considered disabled.                    |
| GDS_BFF_SENTRY_SERVER_NAME       | string  |             | Optional - a server name to tag Sentry events with.                                               |
| GDS_BFF_SENTRY_ENVIRONMENT       | string  |             | The environment to report (e.g. development, staging, production). Required if Sentry is enabled. |
| GDS_BFF_SENTRY_RELEASE           | string  | {{version}} | Specify the release version for Sentry tracking. By default this will be the package version.   |
| GDS_BFF_SENTRY_TRACK_PERFORMANCE | bool    | false       | Enable performance tracing to Sentry with the specified sample rate.                              |
| GDS_BFF_SENTRY_SAMPLE_RATE       | float64 | 0.2         | The percentage of transactions to trace (0.0 to 1.0).                                             |
| GDS_BFF_SENTRY_DEBUG             | bool    | false       | Set Sentry to debug mode for testing.                                                             |

Sentry is considered **enabled** if a DSN is configured. Performance tracing is only enabled if Sentry is enabled *and* track performance is set to true. If Sentry is enabled, an environment is required, otherwise the configuration will be invalid.