---
title: "TrtlDB"
date: 2022-12-23T12:34:23-06:00
lastmod: 2022-12-23T12:34:23-06:00
description: "Configuring the Trtl Database"
weight: 40
---

TrtlDb is a globally replicated key-value store that uses anti-entropy for eventually consistent replication that minimizes egress costs, maximizes availability, and is partition tolerant. It exposes unary and streaming APIs for data access and is secured via mTLS connections. Its environment variables are prefixed with the `TRTL_` tag. The primary configuration is as follows:

| EnvVar               | Type   | Default | Description                                                                                                      |
|----------------------|--------|---------|------------------------------------------------------------------------------------------------------------------|
| TRTL_MAINTENANCE     | bool   | false   | Sets the server to maintenance mode, which will respond to requests with Unavailable except for status requests. |
| TRTL_BIND_ADDR       | string | :4436   | The IP address and port to bind the trtl grpc server to.                                                         |
| TRTL_METRICS_ADDR    | string | :7777   | The IP address and port to bind the prometheus metrics http server to.                                           |
| TRTL_METRICS_ENABLED | bool   | true    | If false, disables the prometheus metrics http server.                                                           |
| TRTL_LOG_LEVEL       | string | info    | The verbosity of logging, one of trace, debug, info, warn, error, fatal, or panic.                               |
| TRTL_CONSOLE_LOG     | bool   | false   | If true, will print human readable logs instead of JSON logs for machine consumption.                             |

### Database

Internally, trtl uses a log-structured merge tree embedded database such as LevelDB to store its pages on disk. This embedded database is configured as follows:

| EnvVar                       | Type   | Default | Description                                                                          |
|------------------------------|--------|---------|--------------------------------------------------------------------------------------|
| TRTL_DATABASE_URL             | string |         | Required, the DSN to connect to the database on (see below for details).             |
| TRTL_DATABASE_REINDEX_ON_BOOT | bool   | false   | When the server starts, instead of loading indexes from disk, recreate and save them. |

The `TRTL_DATABASE_URL` is a standard DSN with a scheme, host, path, and query parameters. Trtl should always have access to a locally embedded database such as LevelDB. When connecting to a LevelDB, the path to the directory on disk where the leveldb should be stored must be passed to the DSN. To specify a relative path, use three slashes: `leveldb:///relpath/to/db`, to specify an absolute path use four: `leveldb:////abspath/to/db`.

### Replica Config

TrtlDB is a globally replicated database and the replica config provides metadata about the TrtlDB instance for replication purposes while also setting parameters for replication. Replication is configured as follows:

| EnvVar                       | Type     | Default | Description                                                                                                           |
|------------------------------|----------|---------|-----------------------------------------------------------------------------------------------------------------------|
| TRTL_REPLICA_ENABLED         | bool     | true    | If false, disables replication so that the trtldb acts as a single node.                                              |
| TRTL_REPLICA_PID             | uint64   |         | The precedence ID of the node (must be unique in the system) - lower IDs have precedence in consistency tie breakers. |
| TRTL_REPLICA_REGION          | string   |         | The region that the replica is assigned to for provenance and geographic compliance.                                  |
| TRTL_REPLICA_NAME            | string   |         | A unique name that identifies the replica name across peers.                                                          |
| TRTL_REPLICA_GOSSIP_INTERVAL | duration | 1m      | The mean interval between gossip (synchronization) sessions between trtl peers.                                       |
| TRTL_REPLICA_GOSSIP_SIGMA    | duration | 5s      | The standard deviation of the jittered interval between gossip sessions.                                              |

If `TRTL_REPLICA_ENABLED` is `true` then `TRTL_REPLICA_PID`, `TRTL_REPLICA_REGION`, and `TRTL_REPLICA_NAME` are required to identify the unique replica in the system.

Replication occurs with bilateral anti-entropy, meaning that after a jittered interval, each replica randomly selects a peer to synchronize with. This is also referred to as a gossip protocol. The `TRTL_REPLICA_GOSSIP_INTERVAL` and `TRTL_REPLICA_GOSSIP_SIGMA` describe a normal distribution of randomly selected synchronization intervals, providing jitter so that the network is not bursty.

### Replica Identification Strategy

When deployed in a kubernetes cluster, a trtl process running in a container in a pod must identify itself to determine what replica it is so that it can configure itself correctly (particularly if the replica is part of a stateful set). The replica strategy configuration defines how a trtl process bootstraps itself as follows:

| EnvVar                             | Type   | Default | Description                                                                                                              |
|------------------------------------|--------|---------|--------------------------------------------------------------------------------------------------------------------------|
| TRTL_REPLICA_STRATEGY_HOSTNAME_PID | bool   | false   | Set to true to use the hostname PID strategy (will be the first strategy tried).                                         |
| TRTL_REPLICA_HOSTNAME              | string |         | Set the hostname of the process from the environment for the hostname-pid strategy.                                      |
| TRTL_REPLICA_STRATEGY_FILE_PID     | string |         | Set to the path of a PID file, if not empty uses the file PID strategy (second strategy after hostname-pid).              |
| TRTL_REPLICA_STRATEGY_JSON_CONFIG  | string |         | Set to the path of a JSON configuration file, if not empty uses the JSON config strategy (third strategy after file-pid). |

The strategies allow the process to identify its PID - either by processing a hostname (e.g. trtl-10) or by reading the pid from a file. The remainder of the replica can then be configured from the JSON file rather than directly from the environment.

### MTLS Config

Connections to the TrtlDB are secured and authenticated using mTLS. The mTLS configuration is as follows:

| EnvVar               | Type   | Default | Description                                                                               |
|----------------------|--------|---------|-------------------------------------------------------------------------------------------|
| TRTL_INSECURE        | bool   | false   | If true, the server will start without mTLS configured.                                    |
| TRTL_MTLS_CHAIN_PATH | string |         | The path to the pool file with valid certificate authorities to authenticate clients for. |
| TRTL_MTLS_CERT_PATH  | string |         | The path to the certificates with the private key for the trtl server.                    |

If `TRTL_INSECURE` is `false` then the `TRTL_MTLS_CHAIN_PATH` and the `TRTL_MTLS_CERT_PATH` are both required.

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
| TRTL_SENTRY_SAMPLE_RATE       | float64 | 0.85        | The percentage of transactions to trace (0.0 to 1.0).                                             |
| TRTL_SENTRY_REPORT_ERRORS     | bool    | false       | If true sends gRPC errors to Sentry as exceptions (useful for development or staging)             |
| TRTL_SENTRY_DEBUG             | bool    | false       | Set Sentry to debug mode for testing.                                                             |

Sentry is considered **enabled** if a DSN is configured. Performance tracing is only enabled if Sentry is enabled *and* track performance is set to true. If Sentry is enabled, an environment is required, otherwise the configuration will be invalid.

Note that the `sentry.Config` object has a field `Repanic` that should not be set by the user. This field is used to manage panics in chained interceptors.