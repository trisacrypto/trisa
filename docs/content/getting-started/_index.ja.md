---
title: TRISAプロトコルをTestNetに統合する方法について説明します
date: 2021-04-23T01:35:35-04:00
lastmod: 2021-10-22T12:37:52-04:00
description: "TRISAプロトコルをTestNetに統合する方法について説明します"
weight: 5
---

## TRISA統合の概要

1. TRISAディレクトリサービスに登録する
2. TRISAネットワークプロトコルを実装します
3. TRISA健康プロトコルを実装します

## VASPディレクトリサービス登録

### 登録の概要

TRISAプロトコルをVASPソフトウェアに統合する前に、TRISAディレクトリサービス（DS）に登録する必要があります。TRISA DSは登録されたVASPの公開鍵およびTRISAリモートピア接続情報を提供します。

TRISA DSに登録すると、KYV証明書を受け取ります。KYV証明書の公開鍵はTRISADを介して他のVASPで利用できるようになります。

DSに登録するときは、VASPがTRISAネットワークサービスを実装する「address：port」エンドポイントを指定する必要があります。 このアドレスはDSに登録され、VASPが受益者VASPとして識別されたときに他のVASPによって利用されます。

統合の目的で、ホストされたテストネットTRISADSインスタンスをテストに使用できます。 テストネットでは登録プロセスが合理化され、迅速な統合が容易になります。 テストネットと統合しながら、本番DS登録を開始することをお勧めします。


### ディレクトリサービス登録

TRISA DSへの登録を開始するには、次のウェブサイトにアクセスしてください。 [https://trisa.directory/](https://trisa.directory/)

「登録」タブを選択して登録を開始できます。このウェブサイトを使用して、フィールドごとに登録の詳細を入力したり、登録の詳細を含むJSONドキュメントをアップロードしたりできることに注意してください。

この登録により、JSONファイルで指定されたすべての技術担当者に電子メールが送信されます。電子メールは、登録プロセスの残りの部分を案内します。登録手順を完了すると、TRISAテストネット管理者はレビューのために登録を受け取ります。

テストネット管理者が登録を確認して承認すると、電子メールでKYV証明書を受け取り、VASPがテストネットDSに公開されます。


## TRISAP2Pプロトコルの実装


### 前提条件

セットアップを開始するには、次のものが必要です。



* KYV証明書（TRISA DS登録から）
* CSRが証明書を取得するために使用する公開鍵
* 関連する秘密鍵
* TRISAディレクトリサービスのホスト名
* TRISAディレクトリサービスでVASPに関連付けられているaddress：portにバインドする機能。


### 統合の概要

TRISAプロトコルの統合には、クライアントコンポーネントとサーバーコンポーネントの両方が含まれます。

クライアントコンポーネントは、TRISAディレクトリサービス（DS）インスタンスとインターフェイスして、TRISAメッセージングプロトコルを統合する他のVASPを検索します。 クライアントコンポーネントは、VASPからの送信トランザクションに使用され、受信VASPがTRISAに準拠していることを確認します。

サーバーコンポーネントは、TRISAプロトコルを統合する他のVASPから要求を受信し、それらの要求への応答を提供します。 サーバーコンポーネントは、VASPがTRISAネットワークプロトコルを満たす情報を返すことができるように実装する必要があるコールバックを提供します。

現在、TRISAネットワークプロトコルのリファレンス実装はGoで利用できます。

[https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go](https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go)

VASPを統合するには、プロトコルの独自の実装を実行する必要があります。 行け以外の言語が必要な場合は、TRISAネットワークプロトコルを定義するプロトコルバッファからクライアントライブラリを生成できます。

インテグレータは、着信転送要求と鍵交換を統合することが期待されており、オプションで、発信転送要求と鍵交換を統合することもできます。

### 統合ノート

TRISAネットワークプロトコルは、参加しているVASP間でデータを転送する方法を定義します。 情報を識別するために転送されるデータの推奨フォーマットは、IVMS101データフォーマットです。 送受信された識別データがFATFトラベルルールを満たしていることを確認するのは、実装するVASPの責任です。

TRISAトランザクションが成功すると、FATFトラベルルールを満たすキーと暗号化されたデータが得られます。 TRISAは、取得後にこのデータを保存する方法を定義していません。 トランザクションの結果データの安全なストレージを処理するのは、実装するVASPの責任です。
