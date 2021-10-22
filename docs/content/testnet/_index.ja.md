---
title: TestNet
date: 2021-06-14T11:20:10-05:00
lastmod: 2021-10-22T12:41:44-05:00
description: "The TRISA TestNet"
weight: 30
---

TRISA テストネットはTRISAピアツーピアプロトコルのデモンストレーションを提供するために設立され、TRISA統合を容易にする「ロボットVASP」サービスをホストし、公開鍵交換とピア検出を容易にするプライマリTRISAディレクトリサービスの場所です。

{{% figure src="/img/testnet_architecture.png" %}}

TRISA テストネットは次のサービスで構成されています。

- [TRISA Directory Service](https://vaspdirectory.net) - TRISAグローバルディレクトリサービスを探索し、TRISAメンバーになるために登録するためのユーザーインターフェイス
- [TestNet Demo](https://vaspbot.net) - テストネットで実行される「ロボット」VASP間のTRISA相互作用を示すデモサイト

テストネットはTRISAメンバーがTRISAサービスを統合するための利便性として実装された3つのロボットVASPまたはrVASPもホストします。 プライマリrVASPはアリスであり、デモ目的のセカンダリはボブです。検証されていないTRISAメンバーとの相互作用をテストするために、「邪悪な」rVASPもあります。

