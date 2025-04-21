---
title: "Deployment"
date: 2022-12-21T11:14:19-06:00
lastmod: 2022-12-21T11:14:19-06:00
description: "Deploying the GDS System"
weight: 10
---

TRISA currently maintains two side-by-side networks that are deployed in Kubernetes clusters in North America, Germany, and Singapore. The "MainNet" (also referred to as the TRISA network) is the production TRISA service where peers exchange compliance information for real transactions. The "TestNet" is a mirror network that is setup to allow peers to develop their TRISA nodes and to ensure that they are configured correctly before engaging in production transactions. For more on using the TestNet, please see the [Testing documentation]({{% ref "testing" %}}).

The current network architecture is as follows:

![TestNet Architecture](/img/system_diagram.png)

## Containers and Images

All TRISA related containers are hosted publicly on [DockerHub](https://hub.docker.com/repositories/trisa). Containers are also pushed to GCR for rapid loading into our GKE kubernetes environment.

The primary container images are as follows:

- [trisa/gds](https://hub.docker.com/repository/docker/trisa/gds): The primary GDS service that runs the TRISA Directory API, the TRISA Members API, the Admin API, and the Certificate Manager process.
- [trisa/gds-user-ui](https://hub.docker.com/repository/docker/trisa/gds-user-ui): The primary GDS UI React application compiled for production and served using an nginx container.
- [trisa/bff](https://hub.docker.com/repository/docker/trisa/gds-bff): The Backend for Frontend (BFF) that allows the GDS UI to access both the MainNet and TestNet GDS services.
- [trisa/gds-admin-ui](https://hub.docker.com/repository/docker/trisa/gds-admin-ui) and [trisa/gds-testnet-admin-ui](https://hub.docker.com/repository/docker/trisa/gds-testnet-admin-ui): the Admin UI React application compiled for production and served using an nginx container.
- [trisa/trtl](https://hub.docker.com/repository/docker/trisa/trtl) and [trisa/trtl-init](https://hub.docker.com/repository/docker/trisa/trtl-init): distributed TrtlDB containers optimized for TRISA cross-region deployment.

Helper images that may be deployed for maintenance or development:

- [trisa/maintenance](https://hub.docker.com/repository/docker/trisa/maintenance): a basic UI that presents a "down for maintenance" screen.
- [trisa/placeholder](https://hub.docker.com/repository/docker/trisa/placeholder): a basic UI that presents a placeholder image and some text.
- [trisa/grpc-proxy](https://hub.docker.com/repository/docker/trisa/grpc-proxy): an envoy proxy configured for using grpc-web with the GDS.
- [trisa/docs-redirect](https://hub.docker.com/repository/docker/trisa/docs-redirect): a simple nginx redirect API for ensuring users are connected to the correct documentation site.

TestNet images:

- [trisa/rvasp](https://hub.docker.com/repository/docker/trisa/rvasp): the base image for all rVASPs that connect to a Postgres database.
- [trisa/rvasp-migrate](https://hub.docker.com/repository/docker/trisa/rvasp-migrate): a container job that migrates the rVASP Postgres database.
- [trisa/trtlsim](https://hub.docker.com/repository/docker/trisa/trtlsim): a data generating utility that simulates TRISA usage for load testing the trtl database.

Staging images (user interface images configured for a staging environment):

- [trisa/gds-staging-user-ui](https://hub.docker.com/repository/docker/trisa/gds-staging-user-ui): GDS UI configured for vaspdirectory.dev
- [trisa/gds-staging-testnet-admin-ui](https://hub.docker.com/repository/docker/trisa/gds-staging-testnet-admin-ui): GDS Admin UI configured for admin.trisatest.dev
- [trisa/gds-staging-admin-ui](https://hub.docker.com/repository/docker/trisa/gds-staging-admin-ui): GDS Admin UI configured for admin.vaspdirectory.dev
- [trisa/cathy](https://hub.docker.com/repository/docker/trisa/cathy): a fake certificate authority for integration testing and development without issuing real certificates.

Deprecated images:

- [trisa/gds-ui](https://hub.docker.com/repository/docker/trisa/gds-ui): GDS v1.4.1 single GDS user interface.
- [trisa/gds-testnet-ui](https://hub.docker.com/repository/docker/trisa/gds-testnet-ui): GDS v1.4.1 single GDS user interface configured for testnet.directory.
- [trisa/rvasp-alice](https://hub.docker.com/repository/docker/trisa/rvasp-alice): the sqlite configured Alice rVASP.
- [trisa/rvasp-bob](https://hub.docker.com/repository/docker/trisa/rvasp-bob): the sqlite configured Bob rVASP.
- [trisa/rvasp-evil](https://hub.docker.com/repository/docker/trisa/rvasp-evil): the sqlite configured Evil rVASP.
