---
title: "访问成员"
date: 2022-02-19T13:25:50-05:00
lastmod: 2022-01-03T14:02:37-05:00
description: "访问其它TRISA成员"
weight: 70
---

`TRISAMembers`服务是一项实验服务，为经过验证的TRISA成员提供额外的、安全的目录服务访问。只有TRISA网络的成员（如，已获得有效身份证明的成员）可使用该服务访问关于网络的详细目录服务数据。*注意：一旦通过验证，这项服务将被移入官方的TRISA规范。

## `TRISAMembers`服务

本节描述了`TRISAMembers`端点的protocol buffers，可以在[这里](https://github.com/trisacrypto/directory/blob/main/proto/gds/members/v1alpha1/members.proto) 找到。 该文件可以编译成您选择的语言（[Golang 的示例](https://github.com/trisacrypto/directory/tree/main/pkg/gds/members/v1alpha1))。 *注意：如果您还没有下载并安装protocol buffer编译器，则您需下载并安装一个。*

目前，`TRISAMembers`服务只包括一个RPC &mdash；`List` RPC。将来可能会提供其他实验性的安全 RPC。

```proto
service TRISAMembers {
    // List all verified VASP members in the Directory Service.
    rpc List(ListRequest) returns (ListReply) {};
}
```

## 列表中的认证会员

`List` RPC 返回一个所有经验证的TRISA成员分页列表，方便TRISA查找。RPC期望输入一个`ListRequest`并返回一个`ListReply`。

### `ListRequest`

`ListRequest`可用来管理VASP列表请求的分页。如果有超过指定页面大小的结果，`ListReply`将返回一个页面标记，可以用来获取下一页（只要不修改原始请求的参数，如任何过滤器或分页参数）。

`page_size`指定每页的结果数量，在不同的页面请求之间不能改变；它的默认值是100。`page_token`指定了获取下一页结果的页面标记。

```proto
message ListRequest {
    int32 page_size = 1;
    string page_token = 2;
}
```

### `ListReply`

`ListReply`返回一个VASP详细信息的简短列表，用于促进点对点密钥交换或对目录服务进行更详细的查询。

`vasps`是一个VASP的列表（见下面`VASPMember`的定义），`next_page_token`，如果指定的话，是存在另一页结果的通知。

```proto
message ListReply {
    repeated VASPMember vasps = 1;
    string next_page_token = 2;
}
```

### `VASPMember`

一个`VASPMember`包含足够的信息，以方便点对点的交流或对目录服务进行更详细的查询。`ListReply`将包含一个没有、一个或多个`VASPMembers`的列表。

```proto
message VASPMember {
    // 目录服务中唯一能识别VASP的组件
    string id = 1;
    string registered_directory = 2;
    string common_name = 3;

    // 连接到远程VASP以执行TRISA请求的地址
    string endpoint = 4;

    // 用于促进搜索和匹配的额外详情
    string name = 5;
    string website = 6;
    string country = 7;
    trisa.gds.models.v1beta1.BusinessCategory business_category = 8;
    repeated string vasp_categories = 9;
    string verified_on = 10;
}
```

## 连接mTLS


要使用`TRISAMembers`服务，您必须使用您在注册时获得的TRISA身份证书，用[mTLS](https://grpc.io/docs/guides/auth/) 进行认证。

[身份验证](https://grpc.io/docs/guides/auth) 上的gRPC文档提供了在各种语言中使用mTLS连接的代码样本，其中包括[Java](https://grpc.io/docs/guides/auth/#java), [C++](https://grpc.io/docs/guides/auth/#c), [Golang](https://grpc.io/docs/guides/auth/#go), [Ruby](https://grpc.io/docs/guides/auth/#ruby), 和 [Python](https://grpc.io/docs/guides/auth/#python)。


举个例子，如果您使用Golang连接到服务目录，您将使用[`tls`](https://pkg.go.dev/crypto/tls)，[`x509`](https://pkg.go.dev/crypto/x509), 和[`credentials`](https://pkg.go.dev/google.golang.org/grpc/credentials) 库，从您电脑上的安全位置加载您的TRISA身份证书，并构建TLS凭证，与服务器相互验证连接。最后，你将使用编译的成员protocol buffer代码来创建一个成员客户端。*注：protocol buffer的定义在本页前面有更全面的描述。

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

*注意，目前有2个TRISA目录；TRISA [TestNet](https://trisatest.net/)，允许用户对TRISA交互进行实验，以及[VASP Directory](https://vaspdirectory.net/)，是用于TRISA交易的生产网络。  如果您已经注册了TestNet并拥有TestNet证书，那么您将传递到拨号功能的端点将是`members.trisatest.net:443`。 或者，如果您想访问VASP会员目录，且已经是一位注册会员，您可以使用终端：`members.vaspdirectory.net:443`*

