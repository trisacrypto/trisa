---
title: "ディレクトリサービスへのVASPの登録"
date: 2021-07-22T09:21:59-04:00
lastmod: 2021-10-22T12:32:23-04:00
description: "ディレクトリサービスへのVASPの登録"
weight: 20
---

TRISAまたはテストネットネットワークに参加するには、TRISAグローバルディレクトリサービス（GDS）または管轄区域固有のディレクトリサービスの1つに登録する必要があります。 ディレクトリサービスへの登録には、次の2つのワークフローがあります。

1. ネットワークが信頼できるメンバーシップを維持していることを確認するためのKYVレビュープロセス
2. ネットワークでのmTLS認証のための証明書の発行

近日公開：登録フォーム、メールによる確認、レビュープロセスの詳細。

##証明書の発行

現在、登録が確認および承認されたときにGDSからmTLS証明書を受信するメカニズムは2つあります。

1. 電子メールで送信されたPKCS12暗号化証明書
2. 証明書署名要求（CSR）

登録を送信するときに、これらのオプションの1つを選択する必要があります。 登録が送信された後は、オプションを切り替えることはできません。

### PKCS12暗号化メールの添付ファイル

最初のメカニズムは最も簡単です- 登録時に電子メールオプションを選択し、CSRフィールドを省略してください。 登録フォームが有効な場合、GDSはPKCS12パスワードを返します。 **このパスワードを紛失しないでください。証明書の発行プロセス中に使用できるようになるのはこのパスワードのみです**。

レビューの承認時に、GDS CAは秘密鍵を含む完全な証明書を生成し、PKCS12パスワードを使用して暗号化します。 ディレクトリサービスに公開鍵を登録した後、GDSは暗号化された証明書をZIPファイルとして技術担当者、または登録フォームで最初に利用可能な連絡先に電子メールで送信します。

電子メールの添付ファイルを解凍すると、`<common_name> .p12`という名前のファイルが見つかります。このファイルを復号化して、次のように証明書を抽出できます。

```
$ openssl pkcs12 -in <common_name>.p12 -out <common_name>.pem -nodes
```

また、.zipファイルを復号化または抽出せずに直接使用することもできます[`github.com/trisacrypto/trisa/pkg/trust`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trust#NewSerializer) モジュール。

### 証明書署名リクエスト

証明書を作成する代わりに、証明書署名要求（CSR）をアップロードすることもできます。このメカニズムは、秘密鍵の素材をネットワーク経由で送信する必要がなく、秘密鍵を安全なハードウェアに残しておくことができるため、多くの場合に適しています。

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
ST = [State or Province (完全に綴られており、略語はありません)]
C = [2 digit country code]
[v3ext_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment, nonRepudiation
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = example.com
```

証明書の構成を慎重に入力してください。この情報は正しくなければならず、証明書を再発行せずに変更することはできません。 また、構成のエントリの後にスペースがないことを確認してください。

次に、次のコマンドを実行し、`example.com`をTRISAエンドポイントとして使用するドメイン名に置き換えます。

```
$ openssl req -new -newkey rsa:4096 -nodes -sha384 -config trisa.conf \
  -keyout example.com.key -out example.com.csr
```

これで、秘密鍵が `example.com.key`＆mdash;に保存されます。 **この秘密鍵を安全に保つ**＆mdash; これは、mTLSサービスのmTLS接続に必要であり、TRISAネットワークで信頼を確立します。

`example.com.csr`ファイルには、証明書署名要求が含まれています。 `----- BEGIN CERTIFICATE REQUEST -----`と `----- END CERTIFICATE REQUEST -----`を含むこのファイルの内容をコピーして登録要求に貼り付けます。
