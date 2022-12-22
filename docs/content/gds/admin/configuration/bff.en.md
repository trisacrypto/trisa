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
| GDS_BFF_COOKIE_DOMAIN | string   |                       | The domain to set secure cookies for (particularly for CSRF and authentication)                                  |
| GDS_BFF_SERVE_DOCS    | bool     | false                 | If true, OpenAPI documentation is complied and served alongside the BFF API.                                     |
| GDS_BFF_USER_CACHE    | bool     | false                 | If true, user information from auth0 is cached in an expiring LRU to reduce the number of look ups.              |

### Auth0 Config

The BFF uses Auth0 for authentication and authorization and must connect to the Auth0 Management API in order to manage users. The Auth0 client is configured as follows:

| EnvVar                       | Type     | Default | Description                                                                                                                                                                        |
|------------------------------|----------|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| GDS_BFF_AUTH0_DOMAIN         | string   |         | The tenant domain provided by the Auth0 application or API (domain only, no scheme or path)                                                                                        |
| GDS_BFF_AUTH0_ISSUER         | string   |         | Set to the custom domain if enabled in Auth0 (ensuring the trailing slash is set if required by the Auth0 configuration) - this will confirm the issuer from the Auth0 JWT tokens. |
| GDS_BFF_AUTH0_AUDIENCE       | string   |         | The audience to verify for the Auth0 API configuration (usually the unique name of the API).                                                                                       |
| GDS_BFF_AUTH0_PROVIDER_CACHE | duration | 5m      | The maximum duration to cache user and application data on Auth0 to prevent repeated data queries.                                                                                 |
| GDS_BFF_AUTH0_CLIENT_ID      | string   |         | The Client ID for the management API specified by Auth0.                                                                                                                           |
| GDS_BFF_AUTH0_CLIENT_SECRET  | string   |         | The Client Secret for the management API specified by Auth0.                                                                                                                       |
| GDS_BFF_AUTH0_TESTING        | bool     | false   | If true a mock authenticator is used for testing purposes.                                                                                                                         |

### Network Configuration



### Database

The BFF makes use of a globally replicated key-value store for persistence of BFF-specific data structures such as organizations and announcements. By default it uses the TrtlDB for this, but for testing or smaller deployments it can use a local LevelDB database instead.

Generally speaking, the BFF uses the same TrtlDB as the GDS MainNet instance, and we've ensured there are no namespace conflicts to prevent this. The BFF primary store can also be independent of the GDS stores if necessary.

The primary database store is configured as follows:

| EnvVar                       | Type   | Default | Description                                                                          |
|------------------------------|--------|---------|--------------------------------------------------------------------------------------|
| GDS_BFF_DATABASE_URL             | string |         | Required, the DSN to connect to the database on (see below for details)              |
| GDS_BFF_DATABASE_REINDEX_ON_BOOT | bool   | false   | When the server starts, instead of loading indexes from disk, recreate and save them |
| GDS_BFF_DATABASE_INSECURE        | bool   | false   | If set do not connect to the TrtlDB with mTLS authentication                         |
| GDS_BFF_DATABASE_CERT_PATH       | string |         | The path to the mTLS client-side certs for database auth                             |
| GDS_BFF_DATABASE_POOL_PATH       | string |         | The path to the mTLS public cert pools to accept server connections                  |

The `GDS_BFF_DATABASE_URL` is a standard DSN with a scheme, host, path, and query parameters. In the case of the BFF the scheme can be either `trtl://` or `leveldb://`.

When connecting to a TrtlDB, the host and port need to be specified with a trailing slash, e.g. `trtl://localhost:4436/`. Connecting via mTLS is only relevant when connecting to a TrtlDB. If `GDS_BFF_DATABASE_INSECURE` is `false`, then the `GDS_BFF_DATABASE_CERT_PATH` and `GDS_BFF_DATABASE_POOL_PATH` are required.

When connecting to a LevelDB, the path to the directory on disk where the leveldb should be stored must be passed to the DSN. To specify a relative path, use three slashes: `leveldb:///relpath/to/db`, to specify an absolute path use four: `leveldb:////abspath/to/db`.

### Email and SendGrid

The BFF uses [SendGrid](https://sendgrid.com/) to send email notifications and to enable communication workflows with users and admins. Configure BFF to use SendGrid as follows:

| EnvVar                | Type   | Default                                           | Description                                                              |
|-----------------------|--------|---------------------------------------------------|--------------------------------------------------------------------------|
| GDS_BFF_SERVICE_EMAIL | string | TRISA Directory Service <admin@vaspdirectory.net> | The email address used as the sender for all emails from the BFF system. |
| SENDGRID_API_KEY      | string |                                                   | API Key to authenticate to SendGrid with.                                |
| GDS_BFF_EMAIL_TESTING | bool   | false                                             | Use email in testing mode rather than send live emails                   |
| GDS_BFF_EMAIL_STORAGE | string | ""                                                | Directory to store test emails for "mark one eyeball" review             |

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