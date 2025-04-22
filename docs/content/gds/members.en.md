---
title: "Accessing Members"
date: 2022-01-03T14:02:37-05:00
lastmod: 2022-01-03T14:02:37-05:00
description: "Accessing Other TRISA Members"
weight: 70
---

The `TRISAMembers` service is an experimental service that provides extra, secure access to the directory service for verified TRISA members. Only members of the TRISA network (e.g. members who have been issued valid identity certificates) can access detailed directory service data about the network using this service. *Note: Once validated, this service will be moved into the official TRISA specification.*

## The `TRISAMembers` Service

This section describes the protocol buffers for the `TRISAMembers` endpoint, which can be found [here](https://github.com/trisacrypto/directory/blob/main/proto/gds/members/v1alpha1/members.proto). This file can be compiled into the language of your choice ([example for Go](https://github.com/trisacrypto/directory/tree/main/pkg/gds/members/v1alpha1)). *Note: You will need to download and install the protocol buffer compiler if you have not already.*

Currently, the `TRISAMembers` service is comprised of three RPCs &mdash; the `List`, `Summary`, and `Details` RPCs. Other experimental, secure RPCs may be made available in the future.

```proto
service TRISAMembers {
    // List all verified VASP members in the Directory Service.
    rpc List(ListRequest) returns (ListReply) {};

    // Get a short summary of the verified VASP members in the Directory Service.
    rpc Summary(SummaryRequest) returns (SummaryReply) {};

    // Get details for a VASP member in the Directory Service.
    rpc Details(DetailsRequest) returns (MemberDetails) {};
}
```

## Listing Verified Members

The `List` RPC returns a paginated list of all _verified_ TRISA members to facilitate TRISA peer lookups. The RPC expects as input a `ListRequest` and returns a `ListReply`.

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
## Summarizing Verified VASP Members in the Directory Service

The `Summary` RPC returns a brief summary of information about verified VASP members in the Directory Service. The RPC expects as input a `SummaryRequest` and returns a `SummaryReply`.

### `SummaryRequest`

A `SummaryRequest` allows the caller to request VASP summary information from the Global Directory Service (GDS). If desired, the caller can indicate a start date to return a count of how many new members have been added since a previous query. The caller can also include a VASP ID to return that VASP's record in the summary.

```proto
message SummaryRequest {
    // The start date for determining how many members are new - optional
    string since = 1;

    // Include your VASP ID to return details about your VASP record in the summary - optional
    string member_id = 2;
}
```

### `SummaryReply`

A `SummaryReply` returns summary info about the members in the Directory Service. This information includes a count of registered VASPs and certificates issued within the Directory Service. `SummaryReply` also provides a count of new members and details about a particular member VASP, if requested.

```proto
message SummaryReply {
    // Counts of VASPs and certificates
    int32 vasps = 1;
    int32 certificates_issued = 2;
    int32 new_members = 3;

    // Details for the requested VASP
    VASPMember member_info = 4;
}
```

## Getting TRIXO and IVMS101 `LegalPerson` Details

The `Details` RPC returns a detailed record for the requested VASP, including the `VASPMember` details described earlier in this page, as well as the TRIXO form and IVMS101 record for the VASP. The RPC expects as input a `DetailsRequest` and returns a `MemberDetails` message.

### `DetailsRequest`

A `DetailsRequest` allows the caller to specify the VASP member to retrieve details for from the Global Directory Service (GDS). The `member_id` is expected to be the VASP's TRISA member ID, a unique identifier to the GDS.

```proto
message DetailsRequest {
    string member_id = 1;
}
```

### `MemberDetails`

A `MemberDetails` message provides details about the requested VASP member, which includes not only their high level `VASPMember` summary, but also the IVMS101 legal person record and responses to the TRIXO questionnaire (both of which are a required component of all TRISA certificate requests).


```proto
message MemberDetails {
    // Summary information about the VASP member
    VASPMember member_summary = 1;

    // The IVMS101 legal person identifying the VASP member
    ivms101.LegalPerson legal_person = 2;

    // The TRIXO questionnaire used to register the VASP
    trisa.gds.models.v1beta1.TRIXOQuestionnaire trixo = 3;
}
```

## Connecting with mTLS

To use the `TRISAMembers` service, you must authenticate with [mTLS](https://grpc.io/docs/guides/auth/) using the TRISA identity certificates you were granted during registration.

The gRPC documentation on [authentication](https://grpc.io/docs/guides/auth) provides code samples for connecting using mTLS in a variety of languages, including [Java](https://grpc.io/docs/guides/auth/#java), [C++](https://grpc.io/docs/guides/auth/#c), [Golang](https://grpc.io/docs/guides/auth/#go), [Ruby](https://grpc.io/docs/guides/auth/#ruby), and [Python](https://grpc.io/docs/guides/auth/#python).

For example, if you were using Golang to connect to the Directory Service, you would use the [`tls`](https://pkg.go.dev/crypto/tls), [`x509`](https://pkg.go.dev/crypto/x509), and [`credentials`](https://pkg.go.dev/google.golang.org/grpc/credentials) libraries to load your TRISA identity certificates from their secure location on your computer and construct TLS credentials to mutually validate the connection with the server. Finally you would use the compiled members protocol buffer code to create a members client. *Note: the protocol buffer definitions are described more fully earlier in this page.*

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

*Note that there are currently two TRISA directories; the TRISA [TestNet]({{% ref "/testing" %}}), which allows users to experiment with the TRISA interactions, and the [VASP Directory](https://trisa.directory/), which is the production network for TRISA transactions. If you have registered for the TestNet and have TestNet certificates, the endpoint you will pass into the dialing function will be `members.testnet.directory:443`. Alternatively, if you wish to access members of the VASP Directory and are already a registered member, you will use the endpoint `members.trisa.directory:443`.*
