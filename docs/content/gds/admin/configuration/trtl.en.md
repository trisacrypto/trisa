---
title: "TrtlDB"
date: 2022-12-23T12:34:23-06:00
lastmod: 2022-12-23T12:34:23-06:00
description: "Configuring the Trtl Database"
weight: 40
---

TrtlDb is a globally replicated key-value store that uses anti-entropy for eventually consistent replication that minimizes egress costs, maximizes availability, and is partition tolerant. It exposes unary and streaming APIs for data accesses and is secured via mTLS connections. Its environment variables are prefixed with the `TRTL_` tag. The primary configuration is as follows:

| EnvVar               | Type   | Default | Description                                                                                                      |
|----------------------|--------|---------|------------------------------------------------------------------------------------------------------------------|
| TRTL_MAINTENANCE     | bool   | false   | Sets the server to maintenance mode, which will respond to requests with Unavailable except for status requests. |
| TRTL_BIND_ADDR       | string | :4436   | The IP address and port to bind the trtl grpc server to.                                                         |
| TRTL_METRICS_ADDR    | string | :7777   | The IP address and port to bind the prometheus metrics http server to.                                           |
| TRTL_METRICS_ENABLED | bool   | true    | If false, disables the prometheus metrics http server.                                                           |
| TRTL_LOG_LEVEL       | string | info    | The verbosity of logging, one of trace, debug, info, warn, error, fatal, or panic.                               |
| TRTL_CONSOLE_LOG     | bool   | false   | If true will print human readable logs instead of JSON logs for machine consumption.                             |

### Database

Internally, trtl uses a log-structured merge tree embedded database such as LevelDB to store its pages on disk. This embedded database is configured as follows:

| EnvVar                       | Type   | Default | Description                                                                          |
|------------------------------|--------|---------|--------------------------------------------------------------------------------------|
| TRTL_DATABASE_URL             | string |         | Required, the DSN to connect to the database on (see below for details).             |
| TRTL_DATABASE_REINDEX_ON_BOOT | bool   | false   | When the server starts, instead of loading indexes from disk, recreate and save them. |

The `TRTL_DATABASE_URL` is a standard DSN with a scheme, host, path, and query parameters. Trtl should always have access to a locally embedded database such as LevelDB. When connecting to a LevelDB, the path to the directory on disk where the leveldb should be stored must be passed to the DSN. To specify a relative path, use three slashes: `leveldb:///relpath/to/db`, to specify an absolute path use four: `leveldb:////abspath/to/db`.

### Backups

The Backup manager is a background process that clones the embedded database pages into a zipped folder on disk that can be downloaded for backup purposes. Backups run periodically and can be run without interupting database processing since the backup takes a live snapshot of the database. Backups are configured as follows:

| EnvVar              | Type     | Default | Description                                                      |
|---------------------|----------|---------|------------------------------------------------------------------|
| TRTL_BACKUP_ENABLED  | bool     | false   | If true, enables the backup background process.                  |
| TRTL_BACKUP_INTERVAL | duration | 24h     | The interval between database backups.                            |
| TRTL_BACKUP_STORAGE  | string   |         | The path on disk to store database backups (required if enabled). |
| TRTL_BACKUP_KEEP     | int      | 1       | The number of backups to keep before cleaning up old backups.     |

### Sentry

TrtlDB uses [Sentry](https://sentry.io/) to assist with error monitoring and performance tracing. Configure TrtlDB to use Sentry as follows:

| EnvVar                       | Type    | Default     | Description                                                                                       |
|------------------------------|---------|-------------|---------------------------------------------------------------------------------------------------|
| TRTL_SENTRY_DSN               | string  |             | The DSN for the Sentry project. If not set then Sentry is considered disabled.                    |
| TRTL_SENTRY_SERVER_NAME       | string  |             | Optional - a server name to tag Sentry events with.                                               |
| TRTL_SENTRY_ENVIRONMENT       | string  |             | The environment to report (e.g. development, staging, production). Required if Sentry is enabled. |
| TRTL_SENTRY_RELEASE           | string  | {{version}} | Specify the release version for Sentry tracking. By default this will be the package version.   |
| TRTL_SENTRY_TRACK_PERFORMANCE | bool    | false       | Enable performance tracing to Sentry with the specified sample rate.                              |
| TRTL_SENTRY_SAMPLE_RATE       | float64 | 0.2         | The percentage of transactions to trace (0.0 to 1.0).                                             |
| TRTL_SENTRY_DEBUG             | bool    | false       | Set Sentry to debug mode for testing.                                                             |

Sentry is considered **enabled** if a DSN is configured. Performance tracing is only enabled if Sentry is enabled *and* track performance is set to true. If Sentry is enabled, an environment is required, otherwise the configuration will be invalid.