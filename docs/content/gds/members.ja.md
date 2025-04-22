---
title: "メンバーへのアクセス"
date: 2022-01-03T14:02:37-05:00
lastmod: 2022-02-25T09:51:50-05:00
description: "他のトリサメンバーがアクセス"
weight: 70
---

あの `TRISAMembers` このサービスは実験的なものであり、検証済みのトリッサメンバーにディレクトリサービスへの追加の安全なアクセスを提供します。ネットワークのメンバーのみが、このサービスを使用してネットワークに関する詳細なディレクトリサービスデータにアクセスできます。
*注：検証されると、このサービスは公式のトリッサ仕様に移行されます。.*

## あの `TRISAMembers` サービス

このセクションでは、protocol buffers のために`TRISAMembers` 見つけることができるエンドポイント。[ここ](https://github.com/trisacrypto/directory/blob/main/proto/gds/members/v1alpha1/members.proto). このファイルは、選択した言語にコンパイルできます。 ([ゴランの例](https://github.com/trisacrypto/directory/tree/main/pkg/gds/members/v1alpha1)). *ノート: ダウンロードしてインストールする必要があります protocol buffer まだ行っていない場合はコンパイラ.*

現在、`TRISAMembers` サービスは単一のRPCのみで構成されます&mdash; あの `List` RPC. その他の実験的で安全RPCs 将来的に利用可能になる可能性があります。.

```proto
service TRISAMembers {
    // List all verified VASP members in the Directory Service.
    rpc List(ListRequest) returns (ListReply) {};
}
```

## 確認済みメンバーのリスト

あの `List` RPC トリッサピアルックアップを容易にするために、検証されたすべてのトリッサメンバーのページ付けされたリストを返します。あの  RPC 入力として期待します`ListRequest` とを返します `ListReply`.

### `ListRequest`

あ `ListRequest` 仮想資産サービスプロバイダーリスト要求のページ付けを管理するために使用できます。 指定されたページサイズよりも多くの結果がある場合`ListReply次のページをフェッチするために使用できるページトークンを返します。
`page_size` ページごとの結果の数を指定し、ページ要求間で変更することはできません。 デフォルト値は100です。 あの  `page_token` 結果の次のページをフェッチするためのページトークンを指定します。

```proto
message ListRequest {
    int32 page_size = 1;
    string page_token = 2;
}
```

### `ListReply`

あ `ListReply` ピアツーピアキー交換またはディレクトリサービスに対するより詳細なルックアップを容易にすることを目的とした仮想アセットサービスプロバイダーの詳細の簡略リストを返します。

`vasps` 仮想資産サービスプロバイダーのリストです（の定義を参照してください）`VASPMember` 以下、および`next_page_token`, 指定されている場合、結果の別のページが存在することの通知です

```proto
message ListReply {
    repeated VASPMember vasps = 1;
    string next_page_token = 2;
}
```

### `VASPMember`

あ`VASPMember` ディレクトリサービスに対するピアツーピア交換またはより詳細なルックアップを容易にするのに十分な情報が含まれています。 あの`ListReply` なし、1つ、または複数のリストが含まれます。`VASPMembers`.

```proto
message VASPMember {
    // ディレクトリサービス内の仮想アセットサービスプロバイダーの一意に識別するコンポーネント
    string id = 1;
    string registered_directory = 2;
    string common_name = 3;

    // トリッサ要求を実行するためにリモート仮想資産サービスプロバイダーに接続するためのアドレス。
    string endpoint = 4;

    // 検索と照合を容易にするために使用される追加の詳細。
    string name = 5;
    string website = 6;
    string country = 7;
    trisa.gds.models.v1beta1.BusinessCategory business_category = 8;
    repeated string vasp_categories = 9;
    string verified_on = 10;
}
```

## mTLSで接続する

を使用するには`TRISAMembers` サービス、あなたはで認証する必要があります[mTLS](https://grpc.io/docs/guides/auth/) 登録時に付与されたトリッサ識別証明書を使用します。
[認証に関するgRPCドキュメント](https://grpc.io/docs/guides/auth) を含むさまざまな言語でmTLSを使用して接続するためのコードサンプルを提供します[Java](https://grpc.io/docs/guides/auth/#java), [C++](https://grpc.io/docs/guides/auth/#c), [Golang](https://grpc.io/docs/guides/auth/#go), [Ruby](https://grpc.io/docs/guides/auth/#ruby), と [Python](https://grpc.io/docs/guides/auth/#python).

たとえば、ディレクトリサービスへの接続にを使用している場合は、[`tls`](https://pkg.go.dev/crypto/tls), [`x509`](https://pkg.go.dev/crypto/x509), と[`credentials`](https://pkg.go.dev/google.golang.org/grpc/credentials) ライブラリを使用して、コンピュータ上の安全な場所からトリッサID証明書をロードし、サーバーとの接続を相互に検証するためのTLSクレデンシャルを構築します。 最後に、コンパイルされたメンバーを使用します protocol buffer メンバークライアントを作成するコード。 *注：protocol buffer 定義については、このページの前半で詳しく説明しています。

```golang
import (
    "crypto/tls"
    "crypto/x509"

    members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
    "google.golang.org/grpc/credentials"
)

func (m *MyProfile) Connect() (_ members.TRISAMembersClient, err error){
    config := &tls.Config{
		Certificates: []tls.Certificate{m.Cert}, // m.Cert is your TRISA certificate parsed into a *x509.Certificate
		MinVersion:   tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
		ClientAuth: tls.RequireAndVerifyClientCert, // this will ensure an mTLS connection
		ClientCAs:  m.Pool, // m.Pool is a *x509.CertPool that must contain the RA and IA public certs from your TRISA certificate chain
	}

    var creds *credentials.TransportCredentials
    if creds, err = credentials.NewTLS(config); err != nil {
        return nil, err
    }

    var cc *grpc.ClientConn
    if cc, err = grpc.NewClient(m.Endpoint, grpc.WithTransportCredentials(creds)); err != nil {
        return nil, err
    }

    return members.NewTRISAMembersClient(cc), nil
}
```

*現在、2つのtrisaディレクトリがあることに注意してください。TRISA [TestNet](https://testnet.directory/), これにより、ユーザーはトリッサの相互作用と仮想アプリケーションサービスプロバイダーディレクトリを試すことができます。[VASP Directory](https://trisa.directory/), これは、トリッサ取引の生産ネットワークです。 テストネットに登録し、証明書を持っている場合、ダイヤル機能に渡すエンドポイントは次のようになります。`members.testnet.directory:443`. または、仮想アプリケーションサービスプロバイダーディレクトリのメンバーにアクセスする場合で、すでに登録済みのメンバーである場合は、エンドポイントを使用します。`members.trisa.directory:443`.*
