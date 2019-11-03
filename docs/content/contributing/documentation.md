---
title: "Documentation"
draft: false
weight: 50
---

## Overview

All TRISA related documentation is maintained at [Github](https://github.com/trisacrypto/trisa/tree/master/docs).
From there the documentation website is automatically generated and hosted on [Github Pages](https://trisacrypto.github.io/).

For quick edits to fix typos or make minor changes, there is a convenience link at the bottom of each page to directly
edit the content on Github directly. From there you can follow the prompts to submit your pull request for review.

## Prerequisites

The TRISA repository contains the necessary tooling to run the documentation website on your machine. This will allow
you to make updates to the content and inspect the results before submitting your changes. 

You only need to have [Git](https://git-scm.com/) and [Docker](https://docs.docker.com/) installed on your system:

* [Docker Engine Community for Linux](https://docs.docker.com/install/)
* [Install Docker Desktop on Mac](https://docs.docker.com/docker-for-mac/install/)
* [Install Docker Desktop on Windows](https://docs.docker.com/docker-for-windows/install/)

## Setup your git directory

Browse to the [TRISA repository](https://github.com/trisacrypto/trisa) and click on `fork` on the right top. Github
will create the fork for you and will be available on `https://github.com/<YOUR-USERNAME>/trisa`.

Follow these steps to setup your git directory using your fork (`origin`) and the main repository (`upstream`). The
example commands are executed from your home directory.

```
git clone git@github.com:<YOUR-USERNAME>/trisa.git/trisa
cd trisa
git remote add upstream git@github.com:trisacrypto/trisa.git
git fetch upstream
git submodule init
git submodule update
```

## Running the doc website

Before making any changes, make sure you have the latest version from the upstream repository. To do so,
perform the following steps from within the `trisa` directory.

```
git fetch upstream
git rebase upstream/master
```

To run and preview the documentation website on your machine:

* Run `make docs-dev` to start the web server
* Browse to `http://127.0.0.1:1313`

The development web server will automatically reload when you make changes to the content files which are
located under `docs/content`. If the preview becomes cluttered, you can always termianted the web server
using `CTRL+C` and running the `make docs-dev` command again to reinitialize.


## Formatting

We are making us of [Gohugo](https://gohugo.io/) to generate the content using markdown `md` files. You can
use existing pages as an example to format your text. Additional information is available on the
[Gohugo documentation section](https://gohugo.io/content-management/).

This website is powered by the [techdoc theme](https://themes.gohugo.io/hugo-theme-techdoc/). The
[example website](https://themes.gohugo.io//theme/hugo-theme-techdoc/sample/build-in-shortcodes/) contains
a rich set of available formatting options. The code for these pages can be found under
`docs/themes/hugo-theme-techdoc/exampleSite` for your reference.  

## Creating your pull request

Once you are happy with your changes, you can stop the documentation web server and push the changes to
your fork. Perform the following steps from your `trisa` directory:

```
git add docs/content
git commit -m "Your commit message, i.e. adding new diagram"
git push origin master:my-doc-changes
```

Now you can browse to `https://github.com/<YOUR-USERNAME>/trisa` and click on "Create Pull Request" from
the yellow popup box. Make sure to create the pull request against the `master` branch of the upstream
`trisacrypto/trisa` repository.

Your pull request will be reviewed and merged when approved. When your pull request is merged, the changes
are automatically pushed the the [TRISA website](https://trisacrypto.github.io/).
