---
title: "GDS"
date: 2022-12-21T18:25:10-06:00
lastmod: 2022-12-21T18:25:10-06:00
description: "Configuring the GDS"
weight: 10
---

The GDS is a collection of three APIs and several background services backed by a global TrtlDB data store. A single GDS node implements the TRISA Directory gRPC API, the secure TRISA Members gRPC API, and the Admin REST API that powers the Admin UI. It also runs the CertMan process that issues and collects certificates from Sectigo. Its environment variables are all prefixed with the `GDS_` tag. The primary configuration is as follows:

| EnvVar           | Type   | Default           | Description                                                                                                      |
|------------------|--------|-------------------|------------------------------------------------------------------------------------------------------------------|
| GDS_DIRECTORY_ID | string | vaspdirectory.net | The network ID the GDS serves, either vaspdirectory.net or trisatest.net (or .dev)                               |
| GDS_SECRET_KEY   | string |                   | A base64 encoded random 32 byte array that is used for salts and random seeds - required.                        |
| GDS_MAINTENANCE  | bool   | false             | Sets the server to maintenance mode, which will respond to requests with Unavailable except for status requests. |
| GDS_LOG_LEVEL    | string | info              | The verbosity of logging, one of trace, debug, info, warn, error, fatal, or panic.                               |
| GDS_CONSOLE_LOG  | bool   | false             | If true will print human readable logs instead of JSON logs for machine consumption.                             |                          |

### GDS API

The primary GDS API (the gRPC TRISA Directory service) is configured as follows:

| EnvVar          | Type   | Default | Description                                            |
|-----------------|--------|---------|--------------------------------------------------------|
| GDS_API_ENABLED | bool   | true    | If false, disables the TRISA Directory service         |
| GDS_BIND_ADDR   | string | :4433   | The IP address and port to bind the GDS gRPC server to |

Note that the enabled flag is only respected if `GDS_MAINTENANCE` is `false`, otherwise maintenance mode supersedes service enabled flags.

### Admin API

The Admin API (a REST API that powers the Admin UI) is configured as follows:

| EnvVar                                   | Type              | Default                                | Description                                                                                            |
|------------------------------------------|-------------------|----------------------------------------|--------------------------------------------------------------------------------------------------------|
| GDS_ADMIN_ENABLED                        | bool              | true                                   | If false, disables the Admin API service                                                               |
| GDS_ADMIN_BIND_ADDR                      | string            | :4434                                  | The IP address and port to bind the Admin API http server to                                           |
| GDS_ADMIN_MODE                           | string            | release                                | Sets the Gin mode, one of debug, release, or test.                                                     |
| GDS_ADMIN_ALLOW_ORIGINS                  | []string          | http://localhost,http://localhost:3000 | A list of allowed origins for the CORS middleware to accept.                                           |
| GDS_ADMIN_COOKIE_DOMAIN                  | string            |                                        | The domain to set secure cookies for (particularly for CSRF and authentication)                        |
| GDS_ADMIN_AUDIENCE                       | string            |                                        | The audience to set and verify in JWT tokens issued by the Admin API                                   |
| GDS_ADMIN_OAUTH_GOOGLE_AUDIENCE          | string            |                                        | The audience from the Google OAuth config to verify Google login tokens                                |
| GDS_ADMIN_OAUTH_AUTHORIZED_EMAIL_DOMAINS | []string          |                                        | The list of authorized email domains to allow access to the admin UI for (e.g. trisa.io)               |
| GDS_ADMIN_TOKEN_KEYS                     | map[string]string |                                        | A mapping of key id (ksuid/ulid) to the path to an RSA signing key in PEM format for JWT token signing |

Note that the enabled flag is only respected if `GDS_MAINTENANCE` is `false`, otherwise maintenance mode supersedes service enabled flags.

### Members API

The Members API is a gRPC API that is secured by TRISA verified mTLS and is configured as follows:

| EnvVar                | Type   | Default | Description                                                           |
|-----------------------|--------|---------|-----------------------------------------------------------------------|
| GDS_MEMBERS_ENABLED   | bool   | true    | If false, disables the Members API service                            |
| GDS_MEMBERS_BIND_ADDR | string | :4435   | The IP address and port to bind the Members gRPC server to            |
| GDS_MEMBERS_INSECURE  | bool   | false   | If set do not enable mTLS authentication                              |
| GDS_MEMBERS_CERTS     | string |         | The path to the mTLS server-side certs for server auth                |
| GDS_MEMBERS_CERT_POOL | string |         | The path to the mTLS public cert pools to accept incoming connections |

Note that the enabled flag is only respected if `GDS_MAINTENANCE` is `false`, otherwise maintenance mode supersedes service enabled flags. If `GDS_MEMBERS_INSECURE` is `false`, then `GDS_MEMBERS_CERTS` and `GDS_MEMBERS_CERT_POOL` are required.

### Database

The GDS makes use of a globally replicated key-value store for persistence. By default it uses the TrtlDB for this, but for testing or smaller deployments it can use a local LevelDB database instead. The database store is configured as follows:

| EnvVar                       | Type   | Default | Description                                                                          |
|------------------------------|--------|---------|--------------------------------------------------------------------------------------|
| GDS_DATABASE_URL             | string |         | Required, the DSN to connect to the database on (see below for details)              |
| GDS_DATABASE_REINDEX_ON_BOOT | bool   | false   | When the server starts, instead of loading indexes from disk, recreate and save them |
| GDS_DATABASE_INSECURE        | bool   | false   | If set do not connect to the TrtlDB with mTLS authentication                         |
| GDS_DATABASE_CERT_PATH       | string |         | The path to the mTLS client-side certs for database auth                             |
| GDS_DATABASE_POOL_PATH       | string |         | The path to the mTLS public cert pools to accept server connections                  |

The `GDS_DATABASE_URL` is a standard DSN with a scheme, host, path, and query parameters. In the case of the GDS the scheme can be either `trtl://` or `leveldb://`.

When connecting to a TrtlDB, the host and port need to be specified with a trailing slash, e.g. `trtl://localhost:4436/`. Connecting via mTLS is only relevant when connecting to a TrtlDB. If `GDS_DATABASE_INSECURE` is `false`, then the `GDS_DATABASE_CERT_PATH` and `GDS_DATABASE_POOL_PATH` are required.

When connecting to a LevelDB, the path to the directory on disk where the leveldb should be stored must be passed to the DSN. To specify a relative path, use three slashes: `leveldb:///relpath/to/db`, to specify an absolute path use four: `leveldb:////abspath/to/db`.

### Email and SendGrid

The GDS uses [SendGrid](https://sendgrid.com/) to send email notifications and to enable communication workflows with users and admins. Configure GDS to use SendGrid as follows:

| EnvVar                 | Type   | Default                                | Description                                                                                 |
|------------------------|--------|----------------------------------------|---------------------------------------------------------------------------------------------|
| GDS_SERVICE_EMAIL      | string | TRISA Directory Service                | The email address used as the sender for all emails from the GDS system.                    |
| GDS_ADMIN_EMAIL        | string | TRISA Admins                           | The email address to send admin emails to from the server.                                  |
| SENDGRID_API_KEY       | string |                                        | API Key to authenticate to SendGrid with.                                                   |
| GDS_DIRECTORY_ID       | string | vaspdirectory.net                      | (Reused) The network ID the GDS serves, either vaspdirectory.net or trisatest.net (or .dev) |
| GDS_VERIFY_CONTACT_URL | string | https://vaspdirectory.net/verify       | The base URL to include in emails for contact verification                                  |
| GDS_ADMIN_REVIEW_URL   | string | https://admin.vaspdirectory.net/vasps/ | The base URL to include in emails to link to a new VASP registration                        |
| GDS_EMAIL_TESTING      | bool   | false                                  | Use email in testing mode rather than send live emails                                      |
| GDS_EMAIL_STORAGE      | string | ""                                     | Directory to store test emails for "mark one eyeball" review                                |

SendGrid is considered **enabled** if the SendGrid API Key is set. The service and admin email addresses are required if SendGrid is enabled.

### CertMan

The CertMan is a background process that runs on the GDS and manages certificate requests, interacts with Sectigo, and finalizes the VASP registration process. It is configured as follows:

| EnvVar                          | Type     | Default           | Description                                                                                                    |
|---------------------------------|----------|-------------------|----------------------------------------------------------------------------------------------------------------|
| GDS_CERTMAN_ENABLED             | bool     | true              | If false, disables the CertMan background process.                                                             |
| GDS_CERTMAN_REQUEST_INTERVAL    | duration | 10m               | The interval between certificate request processing runs                                                       |
| GDS_CERTMAN_REISSUANCE_INTERVAL | duration | 24h               | The interval between certificate reissuance processing runs                                                    |
| GDS_CERTMAN_STORAGE             | string   |                   | The path to a directory on disk where CertMan temporarily downloads certificates (otherwise a tmp dir is used) |
| GDS_DIRECTORY_ID                | string   | vaspdirectory.net | (Reused) The network ID the GDS serves, either vaspdirectory.net or trisatest.net (or .dev)                    |
| SECTIGO_USERNAME                | string   |                   | The username to authenticate with the Sectigo API                                                              |
| SECTIGO_PASSWORD                | string   |                   | The password to authenticate with the Sectigo API                                                              |
| SECTIGO_PROFILE                 | string   |                   | The certificate profile to use (see below for details)                                                         |
| SECTIGO_ENDPOINT                | string   |                   | The endpoint to connect to Sectigo on (leave blank for production)                                             |
| SECTIGO_TESTING                 | bool     |                   | If Sectigo is in testing mode it will not make actual Sectigo requests                                         |

Note that CertMan also needs valid email (SendGrid) and Google Secrets configurations if it is enabled.

The Sectigo username, password, and profile are required if the Sectigo config is not in testing mode. The profile must be a valid certificate profile, which can be either "CipherTrace EE" for the TestNet or "CipherTrace End Entity Certificate" for the MainNet. These profiles determine what parameters must be sent to Sectigo to make certificate requests, and what is populated on the certificate (including what intermediate authority is used) when it is generated.

Do not set the `SECTIGO_ENDPOINT` unless you're in testing, development, or staging mode and are pointing it to a local cathy server.

### Backups

The Backup manager is only available when the GDS is using the LevelDB database store. It is a background process that clones the database into a zipped folder on disk that can be downloaded for backup purposes.

| EnvVar              | Type     | Default | Description                                                      |
|---------------------|----------|---------|------------------------------------------------------------------|
| GDS_BACKUP_ENABLED  | bool     | false   | If true, enables the backup background process.                  |
| GDS_BACKUP_INTERVAL | duration | 24h     | The interval between database backups                            |
| GDS_BACKUP_STORAGE  | string   |         | The path on disk to store database backups (required if enabled) |
| GDS_BACKUP_KEEP     | int      | 1       | The number of backups to keep before cleaning up old backups     |

Backups should not be enabled when using the TrtlDB database store! If the backups are enabled then the storage path is required.

### Google Secrets

GDS uses Google Secret Manager to store certificates and PKCS12 passwords and other sensitive information to ensure that it is encrypted and secure. Access to Google Secret Manager is configured as follows:

| EnvVar                         | Type   | Default | Description                                                                                          |
|--------------------------------|--------|---------|------------------------------------------------------------------------------------------------------|
| GOOGLE_APPLICATION_CREDENTIALS | string |         | Path to the JSON credentials for the Google Service Account that has access to Google Secret Manager |
| GOOGLE_PROJECT_NAME            | string |         | Name of the Google Project for API access to Google Secret Manager                                   |
| GDS_SECRETS_TESTING            | bool   | false   | If set to true, uses a local in-memory "secret store" for testing and development                    |

Note that the `GOOGLE_APPLICATION_CREDENTIALS` and `GOOGLE_PROJECT_NAME` are required if `GDS_SECRETS_TESTING` is `false`.

### Sentry

The GDS uses [Sentry](https://sentry.io/) to assist with error monitoring and performance tracing. Configure GDS to use Sentry as follows:

| EnvVar                       | Type    | Default     | Description                                                                                       |
|------------------------------|---------|-------------|---------------------------------------------------------------------------------------------------|
| GDS_SENTRY_DSN               | string  |             | The DSN for the Sentry project. If not set then Sentry is considered disabled.                    |
| GDS_SENTRY_SERVER_NAME       | string  |             | Optional - a server name to tag Sentry events with.                                               |
| GDS_SENTRY_ENVIRONMENT       | string  |             | The environment to report (e.g. development, staging, production). Required if Sentry is enabled. |
| GDS_SENTRY_RELEASE           | string  | {{version}} | Specify the release of Ensign for Sentry tracking. By default this will be the package version.   |
| GDS_SENTRY_TRACK_PERFORMANCE | bool    | false       | Enable performance tracing to Sentry with the specified sample rate.                              |
| GDS_SENTRY_SAMPLE_RATE       | float64 | 0.2         | The percentage of transactions to trace (0.0 to 1.0).                                             |
| GDS_SENTRY_DEBUG             | bool    | false       | Set Sentry to debug mode for testing.                                                             |

Sentry is considered **enabled** if a DSN is configured. Performance tracing is only enabled if Sentry is enabled *and* track performance is set to true. If Sentry is enabled, an environment is required, otherwise the configuration will be invalid.