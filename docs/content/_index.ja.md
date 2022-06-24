---
title: ホームページ
date: 2020-12-24T07:58:37-05:00
lastmod: 2021-10-15T14:35:53-05:00
description: "TRISA デベロッパードキュメント"
weight: 0
---

## TRISA

[![Go Reference](https://pkg.go.dev/badge/github.com/trisacrypto/trisa/pkg.svg)](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg)
[![Go Report Card](https://goreportcard.com/badge/github.com/trisacrypto/trisa)](https://goreportcard.com/report/github.com/trisacrypto/trisa)

トラベルルール情報共有アーキテクチャ（TRISA）の目標は、コアブロックチェーンプロトコルを変更したり、トランザクションコストを増加させたり、仮想通貨のピアツーピアトランザクションを変更したりすることなく、暗号通貨トランザクションID情報のFATFおよびFinCENトラベルルールへの準拠を可能にすることです流れ。TRISAプロトコルと仕様は、[TRISAワーキンググループ](https://trisa.io) によって定義されています。 仕様の詳細については、[現在のバージョンのTRISAホワイトペーパーをお読みください](https://trisa.io/trisa-whitepaper/)。

このサイトには、[github.com/trisacrypto/trisa](https://github.com/trisacrypto/trisa) にあるTRISAプロトコルとリファレンス実装の開発者向けドキュメントが含まれています。 TRISAプロトコルは[gRPC API](https://grpc.io/) として定義されており、トラベルルールを実装する必要のある仮想資産サービスプロバイダー（VASP）間で、言語に依存しない高性能のピアツーピアサービスを促進しますコンプライアンスソリューション。 APIとメッセージ交換の形式はどちらも、[`protos`ディレクトリ](https://github.com/trisacrypto/trisa/tree/main/proto) にある[protocol buffers](https://developers.google.com/protocol-buffers) を介して定義されます。リポジトリの 。さらに、[Goプログラミング言語](https://golang.org/) のリファレンス実装が [`pkg`ディレクトリ](https://github.com/trisacrypto/trisa/tree/main/proto) で利用できるようになりました。 リポジトリの。将来的には、リポジトリの[`lib`ディレクトリ](https://github.com/trisacrypto/trisa/tree/main/lib) にある特定の言語のライブラリコードとして他の実装が利用できるようになる予定です。

TRISA v1リリースは活発に開発されており、さらに多くのドキュメントが間もなく公開されます。


