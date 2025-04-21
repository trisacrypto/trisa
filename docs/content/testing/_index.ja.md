---
title: TestNet
date: 2021-06-14T11:20:10-05:00
lastmod: 2021-10-22T12:41:44-05:00
description: "The TRISA TestNet"
weight: 30
---

TRISA TestNetはTRISAピアツーピアプロトコルのデモンストレーションを提供するために設立され、TRISA統合を容易にする「ロボットVASP」サービスをホストし、公開鍵交換とピア検出を容易にするプライマリTRISAディレクトリサービスの場所です。

{{% figure src="/img/testnet_architecture.png" %}}

TRISA TestNetは次のサービスで構成されています。

- [TRISA Directory Service](https://trisa.directory) - TRISAグローバルディレクトリサービスを探索し、TRISAメンバーになるために登録するためのユーザーインターフェイス
- [TestNet Demo](https://vaspbot.com) - TestNetで実行される「ロボット」VASP間のTRISA相互作用を示すデモサイト

TestNetはTRISAメンバーがTRISAサービスを統合するための利便性として実装された3つのロボットVASPまたはrVASPもホストします。 プライマリrVASPはAliceであり、デモ目的のセカンダリはBobです。検証されていないTRISAメンバーとの相互作用をテストするために、「Evil」rVASPもあります。
