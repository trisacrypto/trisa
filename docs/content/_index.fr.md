---
title: Accueil
date: 2020-12-24T07:58:37-05:00
lastmod: 2021-10-13T14:33:07-05:00
description: "Documentation pour les développeurs TRISA"
weight: 0
---

# TRISA

[![Go Reference](https://pkg.go.dev/badge/github.com/trisacrypto/trisa/pkg.svg)](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg)
[![Go Report Card](https://goreportcard.com/badge/github.com/trisacrypto/trisa)](https://goreportcard.com/report/github.com/trisacrypto/trisa)

L'objectif de l'architecture de partage d'informations sur les règles de voyage (TRISA) est de permettre la mise en conformité avec les règles de voyage du GAFI et du FinCEN pour les informations d'identité des transactions en crypto-monnaies sans modifier les protocoles de base de la blockchain, et sans encourir de coûts de transaction accrus ni modifier les flux de transaction P2P en monnaie virtuelle. Le protocole et la spécification TRISA sont définis par le [groupe de travail TRISA](https://trisa.io) ; pour en savoir plus sur la spécification, [veuillez lire la version actuelle du livre blanc TRISA](https://trisa.io/trisa-whitepaper/).

Ce site contient la documentation pour développeurs du protocole TRISA et la mise en œuvre de référence qui se trouve sur [github.com/trisacrypto/trisa](https://github.com/trisacrypto/trisa). Le protocole TRISA est défini comme une [API gRPC](https://grpc.io/) conçue pour faciliter le service pair à pair, performants et indépendants du langage entre fournisseurs de services et d'actifs virtuels (VASP) qui doivent mettre en œuvre des solutions de conformité aux règles de voyage. L'API et le format d'échange de messages sont définis via les [tampons de protocole](https://developers.google.com/protocol-buffers), qui se trouvent dans le [répertoire `protos`](https://github.com/trisacrypto/trisa/tree/main/proto) du référentiel. De plus, une implémentation de référence dans le [langage de programmation Go](https://golang.org/) a été mise à disposition dans le [répertoire `pkg`](https://github.com/trisacrypto/trisa/tree/main/proto) du référentiel. Dans le futur, d'autres implémentations seront disponibles sous forme de code de bibliothèque pour des langages spécifiques, dans le [répertoire `lib`](https://github.com/trisacrypto/trisa/tree/main/lib) du référentiel.

La version v1 de TRISA est en cours de développement actif, plus de documentation sera bientôt disponible !
