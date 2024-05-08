---
title: Configuration
date: 2024-05-08T12:14:45-05:00
lastmod: 2024-05-08T12:14:45-05:00
description: "Configuring Envoy"
weight: 20
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

This application is configured via the environment. The following environment
variables can be used:

```
TRISA_MAINTENANCE
  [description] if true, the node will start in maintenance mode
  [type]        True or False
  [default]     false
  [required]
TRISA_MODE
  [description] specify the mode of the server (release, debug, testing)
  [type]        String
  [default]     release
  [required]
TRISA_LOG_LEVEL
  [description] specify the verbosity of logging (trace, debug, info, warn, error, fatal panic)
  [type]        LevelDecoder
  [default]     info
  [required]
TRISA_CONSOLE_LOG
  [description] if true logs colorized human readable output instead of json
  [type]        True or False
  [default]     false
  [required]
TRISA_DATABASE_URL
  [description] dsn containing backend database configuration
  [type]        String
  [default]     sqlite3:///trisa.db
  [required]
TRISA_WEB_TRISA_MAINTENANCE
  [description] if true sets the web UI to maintenance mode; inherited from parent
  [type]        True or False
  [default]
  [required]
TRISA_WEB_ENABLED
  [description] if false, the web UI server will not be run
  [type]        True or False
  [default]     true
  [required]
TRISA_WEB_BIND_ADDR
  [description] the ip address and port to bind the web server on
  [type]        String
  [default]     :8000
  [required]
TRISA_WEB_ORIGIN
  [description] origin (url) of the web ui for creating endpoints and CORS access
  [type]        String
  [default]     http://localhost:8000
  [required]
TRISA_WEB_TRISA_ENDPOINT
  [description] trisa endpoint as assigned to the mTLS certificates for the trisa node
  [type]        String
  [default]
  [required]
TRISA_WEB_AUTH_KEYS
  [description] optional static key configuration as a map of keyID to path on disk
  [type]        Comma-separated list of String:String pairs
  [default]
  [required]    false
TRISA_WEB_AUTH_AUDIENCE
  [description] value for the aud jwt claim
  [type]        String
  [default]     http://localhost:8000
  [required]
TRISA_WEB_AUTH_ISSUER
  [description] value for the iss jwt claim
  [type]        String
  [default]     http://localhost:8000
  [required]
TRISA_WEB_AUTH_COOKIE_DOMAIN
  [description] limit cookies to the specified domain (exclude port)
  [type]        String
  [default]     localhost
  [required]
TRISA_WEB_AUTH_ACCESS_TOKEN_TTL
  [description] the amount of time before an access token expires
  [type]        Duration
  [default]     1h
  [required]
TRISA_WEB_AUTH_REFRESH_TOKEN_TTL
  [description] the amount of time before a refresh token expires
  [type]        Duration
  [default]     2h
  [required]
TRISA_WEB_AUTH_TOKEN_OVERLAP
  [description] the amount of overlap between the access and refresh token
  [type]        Duration
  [default]     -15m
  [required]
TRISA_NODE_TRISA_MAINTENANCE
  [description] if true sets the TRISA node to maintenance mode; inherited from parent
  [type]        True or False
  [default]
  [required]
TRISA_NODE_TRISA_ENDPOINT
  [description] trisa endpoint as assigned to the mTLS certificates for the trisa node
  [type]        String
  [default]
  [required]
TRISA_NODE_ENABLED
  [description] if false, the TRISA node server will not be run
  [type]        True or False
  [default]     true
  [required]
TRISA_NODE_BIND_ADDR
  [description]
  [type]        String
  [default]     :8100
  [required]
TRISA_NODE_POOL
  [description]
  [type]        String
  [default]
  [required]    false
TRISA_NODE_CERTS
  [description]
  [type]        String
  [default]
  [required]    false
TRISA_NODE_KEY_EXCHANGE_CACHE_TTL
  [description]
  [type]        Duration
  [default]     24h
  [required]
TRISA_NODE_DIRECTORY_INSECURE
  [description] if true, do not connect using TLS
  [type]        True or False
  [default]     false
  [required]
TRISA_NODE_DIRECTORY_ENDPOINT
  [description] the endpoint of the public GDS service
  [type]        String
  [default]     api.vaspdirectory.net:443
  [required]    true
TRISA_NODE_DIRECTORY_MEMBERS_ENDPOINT
  [description] the endpoint of the members only GDS service
  [type]        String
  [default]     members.vaspdirectory.net:443
  [required]    true
TRISA_DIRECTORY_SYNC_ENABLED
  [description] if false, the sync background service will not be run
  [type]        True or False
  [default]     true
  [required]
TRISA_DIRECTORY_SYNC_INTERVAL
  [description] the interval synchronization is run
  [type]        Duration
  [default]     6h
  [required]
```