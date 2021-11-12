---
title: テストネット統合のためのrVASPの操作
date: 2021-06-14T12:50:23-05:00
lastmod: 2021-10-22T14:49:33-05:00
description: "テストネット統合のためのrVASPの操作"
weight: 10
---

TestNetは、3つの便利な「ロボットVASP」（rVASP）サービスをホストして、TRISA TestNetとの統合とテストを容易にします。 これらのサービスは次のとおりです。

- Alice(`api.alice.vaspbot.net:443`): TRISAメッセージのトリガーと受信に使用されるプライマリ統合rVASP。
- Bob(`api.bob.vaspbot.net:443`): Aliceとの交換を表示するデモrVASP。
- Evil(`api.evil.vaspbot.net:443`): 認証されていない相互作用をテストするために使用される、TRISA TestNetメンバーではない「Evil」rVASP。

注：rVASPは現在、主にデモ用に構成されており、統合の目的でより堅牢にするための作業が開始されています。変更がないか、このドキュメントを定期的に確認してください。rVASPコードまたは動作にバグに気付いた場合は、[問題を開いてください](https://github.com/trisacrypto/testnet/issues)。

## rVASPの使用開始

rVASPを使用してTRISAサービスを開発するには、次の2つの方法があります。

1. rVASPをトリガーして、サービスにTRISA交換メッセージを送信できます。
2. 有効な（または無効な）rVASPカスタマーを使用してTRISAメッセージをrVASPに送信できます。

rVASPには、偽のウォレットアドレスを持つ偽の顧客のデータベースが組み込まれています。 TRISAメッセージまたはトリガーされた転送への応答には、発信者/受取人の顧客がrVASPに対して有効である必要があります。例えば。顧客のウォレットアドレスが有効なAliceアドレスであり、Aliceが受取人である場合、rVASPは顧客の偽のKYCデータで応答しますが、そうでない場合は、TRISAエラーコードを返します。

次のAlice、Bob、およびEvilの「顧客」の表は、各rVASPと対話するためのリファレンスとして使用できます。

| VASP                  | "Crypto Wallet"                    | Email                 |
|-----------------------|------------------------------------|-----------------------|
| api.bob.vaspbot.net   | 18nxAxBktHZDrMoJ3N2fk9imLX8xNnYbNh | robert@bobvasp.co.uk  |
| api.bob.vaspbot.net   | 1LgtLYkpaXhHDu1Ngh7x9fcBs5KuThbSzw | george@bobvasp.co.uk  |
| api.bob.vaspbot.net   | 14WU745djqecaJ1gmtWQGeMCFim1W5MNp3 | larry@bobvasp.co.uk   |
| api.bob.vaspbot.net   | **1Hzej6a2VG7C8iCAD5DAdN72cZH5THSMt9** | fred@bobvasp.co.uk    |
| api.alice.vaspbot.net | **19nFejdNSUhzkAAdwAvP3wc53o8dL326QQ** | sarah@alicevasp.us    |
| api.alice.vaspbot.net | 1ASkqdo1hvydosVRvRv2j6eNnWpWLHucMX | mary@alicevasp.us     |
| api.alice.vaspbot.net | 1MRCxvEpBoY8qajrmNTSrcfXSZ2wsrGeha | alice@alicevasp.us    |
| api.alice.vaspbot.net | 14HmBSwec8XrcWge9Zi1ZngNia64u3Wd2v | jane@alicevasp.us     |
| api.evil.vaspbot.net  | **1AsF1fMSaXPzz3dkBPyq81wrPQUKtT2tiz** | gambler@evilvasp.gg   |
| api.evil.vaspbot.net  | 1PFTsUQrRqvmFkJunfuQbSC2k9p4RfxYLF | voldemort@evilvasp.gg |
| api.evil.vaspbot.net  | 172n89jLjXKmFJni1vwV5EzxKRXuAAoxUz | launderer@evilvasp.gg |
| api.evil.vaspbot.net  | 182kF4mb5SW4KGEvBSbyXTpDWy8rK1Dpu  | badnews@evilvasp.gg   |

すべてのrVASPデータは、フェイカーツールを使用して生成され、現実的で一貫性のあるテストデータとフィクスチャを生成し、完全に架空のものであることに注意してください。 たとえば、Alice VASP（偽の米国企業）のレコードは、主に北米などにあります。

あなたがトラベラーの顧客である場合、上記の太字のアドレスにはいくつかのアトリビューションデータが関連付けられており、トラベラーベースのrVASPインタラクションの候補として適しています。

### 予選

このドキュメントは、最新の「TRISA ネットワーク」サービスを実行しているサービスがあり、そのサービスがTRISA TestNetに登録されており、TestNet証明書が正しくインストールされていることを前提としています。 [TRISA統合の概要を参照してください]({{< ref "integration/_index.md" >}}) 詳細については。 **警告**: rVASPはTRISAネットワークに参加せず、検証済みのTRISA  TestNet mTLS接続にのみ応答します。

rVASP APIと対話するには、次のいずれかを実行できます。

1. `rvasp`CLIツールを使用します
2. rVASPプロトコルバッファを使用し、APIと直接対話します

CLIツールをインストールするには、[TestNetリリースページ](https://github.com/trisacrypto/testnet/releases) で適切なアーキテクチャの `rvasp`実行可能ファイルをダウンロードし、[TestNetリポジトリ](https//github.com/trisacrypto/testnet/) そして、 `cmd/rvasp`バイナリをビルドするか、次のように`goget`でインストールします。

```
$ go get github.com/trisacrypto/testnet/...
```

[rVASPプロトコルバッファ](https://github.com/trisacrypto/testnet/tree/main/proto/rvasp/v1)を使用するには, それらをTestNetリポジトリから複製またはダウンロードしてから、`protoc`を使用して好みの言語にコンパイルします。

### メッセージを送信するためのrVASPのトリガー

rVASP管理エンドポイントは、開発および統合の目的でrVASPと直接対話するために使用されます。 このエンドポイントは、前述のTRISAエンドポイントとは異なることに注意してください。

- Alice: `admin.alice.vaspbot.net:443`
- Bob: `admin.bob.vaspbot.net:443`
- Evil: `admin.evil.vaspbot.net:443`

コマンドラインツールを使用してメッセージをトリガーするには、次のコマンドを実行します。

```
$ rvasp transfer -e admin.alice.vaspbot.net:443 \
        -a mary@alicevasp.us \
        -d 0.3 \
        -B trisa.example.com \
        -b cryptowalletaddress \
        -E
```

このメッセージは、`-e`または` --endpoint`フラグを使用してアリスrVASPにメッセージを送信し、 `-a`または` --account`フラグを使用して発信元アカウントが"mary@alicevasp.us"であることを指定します 。 元のアカウントは、受取人に送信するIVMS101データを決定するために使用されます。 `-d`または` --amount`フラグは、送信する「アリスコイン」の量を指定します。

次の2つの部分は重要です。 `-E`または` --external-demo`フラグは、別のrVASPとのデモ交換を実行するのではなく、サービスへの要求をトリガーするようにrVASPに指示します。 このフラグは必須です！ 最後に、 `-B`または` --beneficiary-vasp`フラグは、rVASPがリクエストを送信する場所を指定します。このフィールドは、TRISA TestNet ディレクトリサービスで検索できる必要があります。例えば共通名、または検索可能な場合はVASPの名前にする必要があります。

`$RVASP_ADDR`および` $RVASP_CLIENT_ACCOUNT`環境変数を設定してそれぞれ `-e`および` -a`フラグを指定できることに注意してください。

プロトコルバッファを直接使用するには、次の `TransferRequest`で` TRISAIntegration`サービス `Transfer`RPCを使用します。

```json
{
    "account": "mary@alicevasp.us",
    "amount": 0.3,
    "beneficiary": "cryptowalletaddress",
    "beneficiary_vasp": "trisa.example.com",
    "check_beneficiary": false,
    "external_demo": true
}
```

これらの値は、コマンドラインプログラムの値とまったく同じ仕様です。

### TRISAメッセージをrVASPに送信する

rVASPは、トランザクションペイロードとして`trisa.data.generic.v1beta1.Transaction`を、IDペイロードとして` ivms101.IdentityPayload`を想定しています。 IDペイロードの受取人情報を入力する必要はありませんrVASPは受取人の入力に応答しますが、IDペイロードをnullにすることはできません。rVASPの解析および検証コマンドを利用するために、偽のIDデータを指定することをお勧めします。

トランザクションペイロードで、上記の表のrVASP受取人と一致する受取人ウォレットを指定していることを確認してください。例えば 使用：

```json
{
    "txid": "1234",
    "originator": "anydatahereisfine",
    "beneficiary": "1MRCxvEpBoY8qajrmNTSrcfXSZ2wsrGeha",
    "amount": 0.3,
    "network": "TestNet",
    "timestamp": "2021-06-14T16:41:52-05:00"
}
```

`txid`または` originator`文字列を指定でき、`network`および` timestamp`フィールドは無視されます。

ディレクトリサービスまたは直接キー交換を使用してrVASPRSA公開鍵をフェッチし、エンベロープ暗号化として「AES256-GCM」および「HMAC-SHA256」を使用して、封印されたエンベロープを作成します。 次に、 `TRISANetwork`サービス` Transfer` RPCを使用して、封印された封筒をrVASPに送信します。

TODO：まもなく `trisa`コマンドラインプログラムが利用可能になります。 CLIプログラムを使用して、メッセージがリリースされたらすぐに送信する方法をここで指定します。
