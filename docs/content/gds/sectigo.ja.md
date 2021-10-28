---
title: "Sectigo"
date: 2020-12-24T07:58:37-05:00
lastmod: 2022-10-22T12:34:20-05:00
description: "ディレクトリサービスとSectigoCAAPIとのやり取り"
weight: 50
---

[![Go Reference](https://pkg.go.dev/badge/github.com/trisacrypto/directory/pkg/sectigo.svg)](https://pkg.go.dev/github.com/trisacrypto/directory/pkg/sectigo)

TRISAディレクトリサービスは、IoTポータル経由でSectigo認証局を使用して証明書を発行します。 ディレクトリサービスは、mTLSの最初の信頼できるハンドシェイクを容易にするために公開鍵マテリアルを収集する必要があるため、VASP登録および検証プロセスの一部としてSectigo IoT マネジャーAPIを使用します。 `github.com/trisacrypto/directory/pkg/sectigo`パッケージは、APIと対話し、ディレクトリサービスに必要なエンドポイントとメソッドを実装するための行けライブラリです。テストネットは、管理およびデバッグの目的でAPIと対話するためのコマンドラインユーティリティも提供します。 このドキュメントでは、コマンドラインユーティリティについて説明します。このユーティリティでは、APIを直接使用して証明書を発行および取り消す方法の概要も説明しています。

参考資料:

- [パッケージドキュメント](https://pkg.go.dev/github.com/trisacrypto/directory/pkg/sectigo)
- [IoT マネジャーAPI ドキュメンテーション](https://support.sectigo.com/Com_KnowledgeDetailPage?Id=kA01N000000bvCJ)
- [IoT Manager Portal](https://iot.sectigo.com)

## 取得開始

`sectigo` CLIユーティリティをインストールするには、[GitHubのリリース](https://github.com/trisacrypto/directory/releases) からコンパイル済みのバイナリをダウンロードします またはローカルにインストール使用：

```
$ go get github.com/trisacrypto/directory/cmd/sectigo
```

これにより、 `sectigo`コマンドが` $PATH`に追加されます。

## 認証

最初のステップは認証です。ユーザー名とパスワードを `$SECTIGO_USERNAME`と` $SECTIGO_PASSWORD`環境変数に設定する必要があります（または、コマンドラインでパラメーターとして渡すこともできます）。認証ステータスを確認するには、次を使用できます。

```
$ sectigo auth
```

APIはユーザー名とパスワードで認証し、ローカルキャッシュファイルに保存されているアクセストークンと更新トークンを返します。 キャッシュが保存されている場所を確認するには：

```
$ sectigo auth --cache
```

クレデンシャルの状態を確認したい場合、例： アクセストークンが有効、更新可能、または期限切れの場合は、次を使用します。

```
$ sectigo auth --debug
```

## 当局とプロフィール

証明書の操作を開始するには、ユーザーアカウントがアクセスできる権限とプロファイルを一覧表示する必要があります。

```
$ sectigo authorities
[
  {
    "id": 1,
    "ecosystemId": 100,
    "signerCertificateId": 0,
    "ecosystemName": "TRISA",
    "balance": 10,
    "enabled": true,
    "profileId": 42,
    "profileName": "TRISA Profile"
  }
]
```

権限は、証明書が作成されるメソッドとプロファイルを表示します。 ここで、 `profileId`フィールドは、後続の呼び出しで使用するために非常に重要です。 次のように、すべての機関で注文/発行されたライセンスの数を確認することもできます。

```
$ sectigo licenses
{
  "ordered": 2,
  "issued": 2
}
```

プロファイルの詳細情報を取得するには、次のコマンドでprofileIdを使用します。

```
$ sectigo profiles -i 42
```

これにより、生のプロファイル構成が返されます。 権限で証明書を作成する前に、必要なプロファイルパラメータを知っておく必要があります。

```
$ sectigo profile -i 42 --params
```

## 証明書の作成

次のように、`commonName`および` pkcs12Password`パラメーターを使用して証明書の作成を要求できます（他のパラメーターを必要とするプロファイルについては、コードベースを直接使用して独自のメソッドを実装する必要があります）。

```
$ sectigo create -a 42 -d example.com -p secrtpasswrd -b "example.com certs"
{
  "batchId": 24,
  "orderNumber": 1024,
  "creationDate": "2020-12-10T16:35:32.805+0000",
  "profile": "TRISA Profile",
  "size": 1,
  "status": "CREATED",
  "active": false,
  "batchName": "example.com certs",
  "rejectReason": "",
  "generatorParametersValues": null,
  "userId": 10,
  "downloadable": true,
  "rejectable": true
}
```

`-a`フラグは権限を指定しますが、プロファイルIDである必要があります。 ドメインは有効なドメインである必要があります。パスワードを指定しない場合、パスワードが生成され、終了する前にCLIに出力されます。 `-b`フラグは、バッチ作成に人間が読める形式の名前を付けます。戻りデータには、作成されたバッチ証明書ジョブに関する詳細が表示されます。次のように、データをフェッチしてステータスをチェックし続けることができます。

```
$ sectigo batches -i 24
```

バッチの処理情報を取得することもできます。

```
$ sectigo batches -i 24 --status
```

バッチが作成されたら、ジップファイルで証明書をダウンロードします。

```
$ sectigo download -i 24 -o certs/
```

これにより、バッチファイル（通常はbatchId.zip、この場合は24.zip）が `certs/`ディレクトリにダウンロードされます。 証明書を解凍してから、次のように.pemファイルを復号化します。

```
$ unzip certs/24.zip
$ openssl pkcs12 -in certs/example.com.p12 -out certs/example.com.pem -nodes
```

PKCS12ファイルの操作の詳細については、を参照してください[オープンSSLを使用してPKCS＃12ファイルから証明書と秘密鍵をエクスポートする](https://www.ssl.com/how-to/export-certificates-private-key-from-pkcs12-file-with-openssl/).

## CSRのアップロード

証明書を作成する代わりに、証明書署名要求（CSR）をアップロードすることもできます。 このメカニズムは、秘密鍵の素材をネットワーク経由で送信する必要がなく、秘密鍵を安全なハードウェアに残すことができるため、多くの場合望ましいものです。

コマンドラインで `openssl`を使用してCSRを生成するには、最初に現在の作業ディレクトリに` trisa.conf`という名前の構成ファイルを作成し、`example.com`をTRISAエンドポイントをホストする予定のドメインに置き換えます。

```conf
[req]
distinguished_name = dn_req
req_extensions = v3ext_req
prompt = no
default_bits = 4096
[dn_req]
CN = example.com
O = [Organization]
L = [City]
ST = [State or Province (完全にスペルアウトされ、略語はありません)]
C = [2 digit country code]
[v3ext_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment, nonRepudiation
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = example.com
```

証明書の構成を慎重に入力してください。この情報は正しくなければならず、証明書を再発行せずに変更することはできません。また、構成のエントリの後にスペースがないことを確認してください。

次に、次のコマンドを実行し、`example.com`をTRISAエンドポイントとして使用するドメイン名に置き換えます。

```
$ openssl req -new -newkey rsa:4096 -nodes -sha384 -config trisa.conf \
  -keyout example.com.key -out example.com.csr
```

次に、CLIプログラムを使用して次のようにCSRをアップロードできます。

```
$ sectigo upload -p 42 <common_name>.csr
{
  "batchId": 24,
  "orderNumber": 1024,
  "creationDate": "2020-12-10T16:35:32.805+0000",
  "profile": "TRISA Profile",
  "size": 1,
  "status": "CREATED",
  "active": false,
  "batchName": "example.com certs",
  "rejectReason": "",
  "generatorParametersValues": null,
  "userId": 10,
  "downloadable": true,
  "rejectable": true
}
```

`-p`フラグは、CSRバッチ要求を使用するプロファイルを指定し、有効なプロファイルIDである必要があります。アップロードされたCSRは、標準のベギン/終了セパレーターを使用して、PEM形式の複数のCSRを含む単一のテキストファイルにすることができます。

##証明書の管理

名前またはシリアル番号で証明書を検索できますが、ほとんどの場合、ドメインまたは一般名で検索してシリアル番号を取得します。

```
$ sectigo find -n example.com
```

シリアル番号を取得したら、次のように証明書を取り消すことができます。

```
$ sectigo revoke -p 42 -r "cessation of operation" -s 12345
```

このコマンドは、 `-p`フラグが付いた証明書を発行したプロファイルID、[RFC 5280理由コード](https://tools.ietf.org/html/rfc5280#section-5.3.1) が`を介して渡されることを想定しています。 -r`フラグ（デフォルトでは指定されていません）、および `-s`フラグを使用した証明書のシリアル番号。このコマンドでエラーが発生しない場合、証明書は正常に取り消されています。

RFC5280の理由は次のとおりです。

- "unspecified"
- "keycompromise"
- "ca compromise"
- "affiliation changed"
- "superseded"
- "cessation of operation"
- "certificate hold"
- "remove from crl"
- "privilege withdrawn"
- "aa compromise"

理由は空白と大文字と小文字を区別しないことに注意してください。
