---
title: "User Interfaces"
date: 2022-12-23T14:04:26-06:00
lastmod: 2022-12-23T14:04:26-06:00
description: "Configuring the React Apps"
weight: 30
---

The GDS user interfaces that power our user and administrative web applications are built and deployed using React and webpack. Because webpack compiles the app into a bundle of HTML, JavaScript, and CSS -- any configuration must be available at build time (e.g. when containers are compiled). This makes configuring our React apps and user interfaces slightly more complicated.

Right now, we use a continuous integration and deployment service to build our containers (e.g. GitHub Actions). During the container build process, environment variables are injected into the container via [build args](https://vsupalov.com/docker-arg-env-variable-guide/), which are in turn processed by the node build script. All environment variables and build args must be prefixed by `REACT_APP_`.

Note that front-end applications _should not_ have any secret configurations. There are, however, sensitive configurations (such as Google Analytics tags, Sentry DSNs, etc.). Any sensitive configurations should be stored securely (e.g. using GitHub Secrets).

### GDS User UI

The build arguments for the GDS User Interface (vaspdirectory.net) are as follows:

| EnvVar                       | Type   | Default   | Description                                                                          |
|------------------------------|--------|-----------|--------------------------------------------------------------------------------------|
| REACT_APP_TRISA_BASE_URL     | string |           | The base URL of the BFF API endpoint, e.g. https://bff.vaspdirectory.net/v1/.        |
| REACT_APP_VERSION_NUMBER     | string |           | The semvar build version of the app (usually parsed from the git tag).               |
| REACT_APP_GIT_REVISION       | string |           | The seven-digit prefix of the git hash of the commit being built.                    |
| REACT_APP_ANALYTICS_ID       | string |           | The Google Analytics tag (e.g. G-XXXXXXXXXX).                                        |
| REACT_APP_SENTRY_DSN         | string |           | The DSN for configuring React to send errors to Sentry.                              |
| REACT_APP_SENTRY_ENVIRONMENT | string | $NODE_ENV | The environment for Sentry logging (not required except for staging or development). |

The GDS User Interface uses Auth0 for authentication. Front-end Auth0 configuration is as follows:

| EnvVar                       | Type   | Default | Description                                                                                |
|------------------------------|--------|---------|--------------------------------------------------------------------------------------------|
| REACT_APP_AUTH0_DOMAIN       | string |         | The domain (or custom domain) to connect to Auth0 on (e.g. auth.vaspdirectory.net).         |
| REACT_APP_AUTH0_CLIENT_ID    | string |         | The ClientID of the Auth0 app as configured in the Auth0 dashboard.                        |
| REACT_APP_AUTH0_REDIRECT_URI | string |         | The callback URI for the application to receive Auth0 redirects after authentication.      |
| REACT_APP_AUTH0_SCOPE        | string |         | The required Auth0 scope (usually 'openid profile email')                                  |
| REACT_APP_AUTH0_AUDIENCE     | string |         | The audience of the tokens, usually the ID of the API (e.g. https://bff.vaspdirectory.net) |

### GDS Admin UI

The build arguments for the GDS Admin UI (admin.vaspdirectory.net and admin.trisatest.net) are as follows:

| EnvVar                       | Type   | Default   | Description                                                                          |
|------------------------------|--------|-----------|--------------------------------------------------------------------------------------|
| REACT_APP_GDS_API_ENDPOINT   | string |           | The base URL of the Admin API endpoint, e.g. https://api.admin.vaspdirectory.net/v2. |
| REACT_APP_GDS_IS_TESTNET     | bool   | false     | True if the Admin UI is managing the TestNet, false if MainNet.                      |
| REACT_APP_VERSION_NUMBER     | string |           | The semvar build version of the app (usually parsed from the git tag).               |
| REACT_APP_GOOGLE_CLIENT_ID   | string |           | The Google Client ID for Google OAuth2 authentication.                               |
| REACT_APP_GIT_REVISION       | string |           | The seven-digit prefix of the git hash of the commit being built.                    |
| REACT_APP_ANALYTICS_ID       | string |           | The Google Analytics tag (e.g. G-XXXXXXXXXX).                                        |
| REACT_APP_SENTRY_DSN         | string |           | The DSN for configuring React to send errors to Sentry.                              |
| REACT_APP_SENTRY_ENVIRONMENT | string | $NODE_ENV | The environment for Sentry logging (not required except for staging or development). |