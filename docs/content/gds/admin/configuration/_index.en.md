---
title: "Configuration"
date: 2022-12-21T18:09:44-06:00
lastmod: 2022-12-21T18:09:44-06:00
description: "Configuration Guide for GDS Services"
weight: 20
---

TRISA GDS and TestNet services are primarily configured using environment variables and will respect [dotenv files](https://github.com/joho/godotenv) in the current working directory. The canonical reference of the configuration for a GDS service is the `config` package of that service (described below). This documentation enumerates the most important configuration variables, their default values, and any hints or warnings about how to use them.

> **Required Configuration**
>
> If a configuration parameter does not have a default value that means it is required and must be specified by the user! If the configuration parameter does have a default value then that environment variable does not have to be set.

## Configuration Documentation

- [GDS Node Configuration]({{% ref "gds/admin/configuration/gds" %}})
- [BFF Service Configuration]({{% ref "gds/admin/configuration/bff" %}})
- [React Apps Configuration]({{% ref "gds/admin/configuration/ui" %}})
- [TrtlDB Configuration]({{% ref "gds/admin/configuration/trtl" %}})