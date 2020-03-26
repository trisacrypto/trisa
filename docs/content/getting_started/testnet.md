---
title: "TRISA TESTNET"
draft: false
weight: 40
---

## Introduction

For testing purposes to develop new clients and integrations, we have a test PKI available and some VASP nodes
running the latest version of TRISA server. This allows for full integration testing without disturbing the
production TRISA network mesh.

Note that the TESTNET PKI is currently in flux, so the root and issuing CA's may still change. To retrieve your
test certificate and private key, visit our [TESTNET Portal Page](http://testnet.trisa.io). You
will need to login using your Github account. You can generate as many test certificates as you want, there
are currently no rate limits implemented.

## Test VASP nodes

The following VASP nodes are available for testing:

* vasp1.trisa.ciphertrace.com
* vasp2.trisa.ciphertrace.com
* vasp3.trisa.ciphertrace.com

Each VASP uses port 8888 for the gRPC peer-to-peer communication and port 9999 for its admin endpoints.
