---
title: Home
date: 2020-12-24T07:58:37-05:00
lastmod: 2021-06-14T11:17:15-05:00
description: "TRISA Developer Documentation"
weight: 0
---

# TRISA

[![Go Reference](https://pkg.go.dev/badge/github.com/trisacrypto/testnet/pkg.svg)](https://pkg.go.dev/github.com/trisacrypto/testnet/pkg)
[![Go Report Card](https://goreportcard.com/badge/github.com/trisacrypto/trisa)](https://goreportcard.com/report/github.com/trisacrypto/trisa)

The goal of the Travel Rule Information Sharing Architecture (TRISA) is to enable
compliance with the FATF and FinCEN Travel Rules for cryptocurrency transaction
identity information without modifying core blockchain protocols, and without
incurring increased transaction costs or modifying virtual currency peer-to-peer
transaction flows. The TRISA protocol and specification is defined by the [TRISA Working Group](https://trisa.io); to learn more about the specification, [please read the current version of the TRISA whitepaper](https://trisa.io/trisa-whitepaper/).

This site contains the developer documentation for the TRISA protocol and reference implementation which can be found at [github.com/trisacrypto/trisa](https://github.com/trisacrypto/trisa). The TRISA protocol is defined as a [gRPC API](https://grpc.io/) to facilitate language-agnostic, high-performance, peer-to-peer services between Virtual Asset Service Providers (VASPs) that must implement Travel Rule compliance solutions. Both the API and message interchange format are defined via [protocol buffers](https://developers.google.com/protocol-buffers), which can be found in the [`protos` directory](https://github.com/trisacrypto/trisa/tree/main/proto) of the repository. In addition, a reference implementation in the [Go programming language](https://golang.org/) has been made available in the [`pkg` directory](https://github.com/trisacrypto/trisa/tree/main/proto) of the repository. In the future, other implementations will be made available as library code for specific languages, found in the [`lib` directory](https://github.com/trisacrypto/trisa/tree/main/lib) of the repository.

The v1 release of TRISA is under active development, more documentation is coming soon!