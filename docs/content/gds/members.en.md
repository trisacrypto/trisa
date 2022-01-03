---
title: "Accessing Members"
date: 2022-01-03T14:02:37-05:00
lastmod: 2022-01-03T14:02:37-05:00
description: "Accessing Other TRISA Members"
weight: 70
---

The `TRISAMembers` service is an experimental service that provides extra access to the directory service for verified TRISA members. Only members of the TRISA network (e.g. members who have been issued valid identity certificates) can access other members of the network using this service. *Note: Once validated, this service will be moved into the official TRISA specification.*

## The `TRISAMembers` Service

This section describes the protocol buffers for the `TRISAMembers` endpoint, which can be found [here](https://github.com/trisacrypto/directory/blob/main/proto/gds/members/v1alpha1/members.proto). This file can be compiled into the language of your choice ([example for Go](https://github.com/trisacrypto/directory/tree/main/pkg/gds/members/v1alpha1)). *Note: You will need to download and install the protocol buffer compiler if you have not already.*

The `TRISAMembers` service expects as input a `ListRequest` and returns a `ListReply`.

```proto
service TRISAMembers {
    // List all verified VASP members in the Directory Service.
    rpc List(ListRequest) returns (ListReply) {};
}
```

### `ListRequest`

A `ListRequest` can be used to manage pagination for the VASP list request. If there are more results than the specified page size, the `ListReply` will return a page token that can be used to fetch the next page (so long as the parameters of the original request are not modified, e.g. any filters or pagination parameters).

The `page_size` specifies the number of results per page and cannot change between page requests; it's default value is 100. The `page_token` specifies the page token to fetch the next page of results.

```proto
message ListRequest {
    int32 page_size = 1;
    string page_token = 2;
}
```

### `ListReply`

A `ListReply` returns an abbreviated listing of VASP details intended to facilitate peer-to-peer key exchanges or more detailed lookups against the Directory Service.

The `vasps` are a list of VASPs (see definition of `VASPMember` below), and the `next_page_token`, if specified, is notification that another page of results exists.

```proto
message ListReply {
    repeated VASPMember vasps = 1;
    string next_page_token = 2;
}
```

### `VASPMember`

A `VASPMember` contains enough information to facilitate peer-to-peer exchanges or more detailed lookups against the Directory Service. The `ListReply` will contain a list of none, one, or multiple `VASPMembers`.

```proto
message VASPMember {
    // The uniquely identifying components of the VASP in the directory service
    string id = 1;
    string registered_directory = 2;
    string common_name = 3;

    // Address to connect to the remote VASP on to perform a TRISA request
    string endpoint = 4;

    // Extra details used to facilitate searches and matching
    string name = 5;
    string website = 6;
    string country = 7;
    trisa.gds.models.v1beta1.BusinessCategory business_category = 8;
    repeated string vasp_categories = 9;
    string verified_on = 10;
}
```

## Connecting with mTLS

To use the `TRISAMembers` service, you must authenticate with [mTLS](https://grpc.io/docs/guides/auth/) using the TRISA identity certificates you were granted during registration.

The gRPC documentation on [authentication](https://grpc.io/docs/guides/auth) provides code samples for connecting using mTLS in a variety of languages, including [Java](https://grpc.io/docs/guides/auth/#java), [C++](https://grpc.io/docs/guides/auth/#c), [Golang](https://grpc.io/docs/guides/auth/#go), [Ruby](https://grpc.io/docs/guides/auth/#ruby), and [Python](https://grpc.io/docs/guides/auth/#python).

For example, if you were using Golang to connect to the Global Directory Service, you would use the [`tls`](https://pkg.go.dev/crypto/tls), [`x509`](https://pkg.go.dev/crypto/x509), and [`credentials`](https://pkg.go.dev/google.golang.org/grpc/credentials) libraries to load your TRISA identity certificates from their secure location on your computer and construct TLS credentials to mutually validate the connection with the server. Finally you would use the compiled members protocol buffer code to create a members client; the protocol buffer definitions are described more fully earlier in this page.

```golang
import (
    "crypto/tls"
    "crypto/x509"

    members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
    "google.golang.org/grpc/credentials"
)

func (m *MyProfile) Connect() (_ members.TRISAMembersClient, err error){
    config := &tls.Config{
		Certificates: []tls.Certificate{m.Cert}, // m.Cert is your TRISA certificate
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
		ClientCAs:  m.Pool, // m.Pool is a *x509.CertPool
	}

    var creds *credentials.TransportCredentials
    if creds, err = credentials.NewTLS(config); err != nil {
        return nil, err
    }

    var cc *grpc.ClientConn
    if cc, err = grpc.Dial(m.Endpoint, grpc.WithTransportCredentials(creds)); err != nil {
        return nil, err
    }

    return members.NewTRISAMembersClient(cc), nil
}
```

*Note that there are currently two TRISA directories; the TRISA [TestNet](https://trisatest.net/), which allows users to experiment with the TRISA interactions, and the [VASP Directory](https://vaspdirectory.net/), which is the production network for TRISA transactions. If you have registered for the TestNet and have TestNet certificates, the endpoint you will pass into the dialing function will be `members.testnet:443`. Alternatively, if wish to access members of the VASP Directory and are already a registered member, you will use the endpoint `members.vaspdirectory.net:443`.*
