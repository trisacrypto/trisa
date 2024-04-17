---
title: TRISA-Entwicklerdokumentation
date: 2020-12-24T07:58:37-05:00
lastmod: 2021-10-08T15:17:08-05:00
description: "TRISA-Entwicklerdokumentation"
weight: 0
---

[![Go Reference](https://pkg.go.dev/badge/github.com/trisacrypto/trisa/pkg.svg)](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg)

[![Go Report Card](https://goreportcard.com/badge/github.com/trisacrypto/trisa)](https://goreportcard.com/report/github.com/trisacrypto/trisa)

Das Ziel der Travel Rule Information Sharing Architecture (TRISA) ist es, die Einhaltung der FATF- und FinCEN-Reiseregeln für Kryptowährungs-Transaktionsidentitätsinformationen zu ermöglichen, ohne die Kernprotokolle der Blockchain zu verändern und ohne erhöhte Transaktionskosten zu verursachen oder die Peer-to-Peer-Transaktionsflüsse virtueller Währungen zu verändern. Das TRISA-Protokoll und die Spezifikation werden von der [TRISA Working Group](https://trisa.io) definiert; um mehr über die Spezifikation zu erfahren, [lesen Sie bitte die aktuelle Version des TRISA-Whitepapers](https://trisa.io/trisa-whitepaper/).

Diese Website enthält die Entwicklerdokumentation für das TRISA-Protokoll und die Referenzimplementierung, die unter [github.com/trisacrypto/trisa](https://github.com/trisacrypto/trisa) zu finden ist. Das TRISA-Protokoll ist als [gRPC-API](https://grpc.io/) definiert, um sprachunabhängige, hochleistungsfähige Peer-to-Peer-Dienste zwischen Anbietern virtueller Asset Services (VASPs) zu ermöglichen, die Lösungen zur Einhaltung der Travel Rule implementieren müssen. Sowohl die API als auch das Nachrichtenaustauschformat werden über [protocol buffers](https://developers.google.com/protocol-buffers) definiert, die im [`protos` directory](https://github.com/trisacrypto/trisa/tree/main/proto) des Repository zu finden sind. Darüber hinaus wurde eine Referenzimplementierung in der [Go Programmiersprache](https://golang.org/) im [`pkg` directory](https://github.com/trisacrypto/trisa/tree/main/proto) des Repositorys zur Verfügung gestellt. In Zukunft werden weitere Implementierungen als Bibliothekscode für bestimmte Sprachen im [`lib` directory](https://github.com/trisacrypto/trisa/tree/main/lib) des Repositorys zur Verfügung gestellt.

Die Version v1 von TRISA befindet sich in aktiver Entwicklung, weitere Dokumentation folgt in Kürze!

