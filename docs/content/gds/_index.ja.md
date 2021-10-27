---
title: "Directory Service"
date: 2021-06-14T11:21:42-05:00
lastmod: 2021-10-22T12:31:23-05:00
description: "The Global TRISA Directory Service (GDS)"
weight: 50
---

グローバルTRISAディレクトリサービス（GDS）は、次のようにTRISAメンバー間のピアツーピア交換を容易にします。

- TRISAエンドポイントを見つけるための検出サービスを提供する
- 交換を検証するためにmTLS証明書を発行する
- 検証用の証明書とKYC情報を提供する

ディレクトリサービスとの相互作用は、TRISAプロトコルによって指定されます。 現在、TRISA組織は、TRISAネットワークに代わってGDSをホストしています。 このドキュメントでは、ディレクトリサービスのTRISA実装と、それとのTRISA固有の相互作用について説明します。