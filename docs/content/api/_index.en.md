---
title: TRISA Protocol and API
date: 2022-07-06T15:08:52-04:00
lastmod: 2022-07-06T15:08:52-04:00
description: "Navigating the Open Source TRISA Code"
weight: 15
---

This section of the documentation describes the TRISA Protocol and API.

###  Navigating the Open Source TRISA Code

The [`trisa`](https://github.com/trisacrypto/trisa) repository contains a gRPC implementation of the TRISA protocol as described by the [white paper](https://trisa.io/trisa-whitepaper/), which leverages [protocol buffers](https://grpc.io/) and Golang.

The `proto` folder contains the core RPC definitions, including:
 - the interVASP Messaging Standard (IVMS) message definitions, which serve as the basis for how two VASP peers should mutually describe entities involved in cryptographic transfers, including names, locations, and government identifiers. This is the spec that will allow originators to identify themselves to beneficiaries and to request information from those beneficiaries to meet the legal requirements of their regulators.
 - the TRISA Network's service definitions, essentially how the different parts of the API work &mdash; from the exchange of keys (to ensure both peers have the requisite details to exchange information) to the transfer of "secure envelopes" (cryptographically sealed protocol buffer messages that can only be decrypted by the intended recipient). The `trisa` subfolder also contains generic message types for transactions that are intended to provide maximum flexibility for a wide range of TRISA use cases.

The `pkg` folder contains the reference implementation code, including compiled code generated from the protocol buffer definitions in the `proto` folder[^1].
 - The `iso3166` folder contains language codes.
 - The `ivms101` folder extends the generated protobuf code with JSON loading utilities, validation helpers, short constants, etc.
 - The `trisa` folder contains structs and methods for a range of TRISA-related tasks, such as performing cryptography and making mTLS connections.

 The `lib` folder is intended to showcase utility code similar to that in the `pkg` folder, but for languages other than Go. If you work in a language other than go, this would be a great place to start your contribution!

[^1]: Note that these compiled files are compiled for Golang; but this is certainly not the only option. Those interested in building implementation code in a different language should look to the `lib` folder, which currently contains placeholder folders but is intended to showcase such other implementations (including compiled protocol buffer code for these other languages).

Another integral part of the TRISA protocol is the [Global Directory Service]({{% ref "/gds" %}}), which serves as a look-up tool for TRISA members to identify peers with which they wish to exchange information. For RPC definitions and implementation code related to the Global Directory Service, visit the companion [directory repository](https://github.com/trisacrypto/directory).



