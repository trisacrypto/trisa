---
title: Using the GDS API
date: 2022-07-05T11:54:25-04:00
lastmod: 2022-07-05T11:54:25-04:00
description: "How to use the GDS API as a client"
weight: 20
---


The `TRISADirectory` service, defined in the [`proto/trisa/gds` protocol buffers](https://github.com/trisacrypto/trisa/tree/main/proto/trisa/gds), specifies how TRISA clients should interact with the directory service to facilitate peer-to-peer transfers.

This service includes four types of RPCs:
- RPCs designed to support interactions with the Global Directory Service during a TRISA information transfer: `Lookup` and `Search` (*Note: These are the primary Directory Service interactions.*)
- RPCs designed to support the Directory Service registration and verification processes: `Register` and `VerifyContact`
- RPCs designed to supply the entity review and TRISA verification status of a given VASP: `Verification`
- RPCs designed to enable routine health check and status requests: `Status` (*Note: these mirror the `TRISAHealth` service*)

```proto
service TRISADirectory {
    rpc Lookup(LookupRequest) returns (LookupReply) {}
    rpc Search(SearchRequest) returns (SearchReply) {}
    rpc Register(RegisterRequest) returns (RegisterReply) {}
    rpc VerifyContact(VerifyContactRequest) returns (VerifyContactReply) {}
    rpc Verification(VerificationRequest) returns (VerificationReply) {}
    rpc Status(HealthCheck) returns (ServiceState) {}
}
```

## `Lookup`

The `Lookup` RPC is designed to provide VASP certification information. *Note that only verified VASPs are returned using the `Lookup` RPC.*

This RPC expects a `LookupRequest` message, which provides either the VASP's unique ID in the Global Directory Service, or the domain of their TRISA implementation endpoint (*Note: this should be the common name of their certificate*). If both `id` and `common_name` are supplied, the `id` field will be prioritized.

The `registered_directory` field is the URL of the directory that registered the VASP (e.g. either `"trisatest.net"` or `"vaspdirectory.net"`). If omitted, this value defaults to the directory currently being queried.

```proto
message LookupRequest {
    string id = 1;
    string registered_directory = 2;
    string common_name = 3;
}
```

{{% notice note %}}
To use other name fields such as the legal business name, you must use the `Search` RPC.
{{% /notice %}}


A `LookupRequest` returns a `LookupReply` message, which contains summary information intended to facilitate peer-to-peer verification and public key exchange.

```proto
message LookupReply {
    // If no error is specified, the lookup was successful
    Error error = 1;

    // The uniquely identifying components of the VASP in the directory service
    string id = 2;
    string registered_directory = 3;
    string common_name = 4;

    // The endpoint to connect to for the TRISA P2P protocol (addr:port)
    string endpoint = 5;

    // Certificate information if the VASP is available and verified. The identity
    // certificate is used to establish mTLS connections, the signing certificate is
    // used on a per-transaction basis.
    trisa.gds.models.v1beta1.Certificate identity_certificate = 6;
    trisa.gds.models.v1beta1.Certificate signing_certificate = 7;

    // Other VASP information to facilitates P2P exchanges
    string name = 8;
    string country = 9;
    string verified_on = 10;
}
```

## `Search`

The `Search` RPC allows more flexible search to identify member VASPs. These requests are primarily used to locate a beneficiary VASP in order to begin the TRISA peer-to-peer protocol. *Note that only verified VASPs are returned using the `Search` RPC.*

The `Search` RPC expects a `SearchRequest` and returns a `SearchReply`. The `SearchRequest` message enables searching by (case insensitive) `name` and `website` (using `OR` logic). The `name` field can be the legal, short, or DBA name of the VASP. It could also be the common name of the certificate issued to that VASP, though in this case, it is recommended to use the `Lookup` RPC instead. The `website` must be a parseable URL.

The remaining fields are search filters on which to condition the search. Only VASPs in the specified `country` or `category` (both of which can be a single value or multiple values) will be returned.

```proto
message SearchRequest {
    repeated string name = 1;
    repeated string website = 2;
    repeated string country = 7;
    repeated trisa.gds.models.v1beta1.BusinessCategory business_category = 8;
    repeated string vasp_category = 9;
}
```

A `SearchRequest` returns a `SearchReply` message containing potentially multiple `Result` messages. *Note: if no `error` is returned, the search was successful, even if no results were returned.*


```proto
message SearchReply {
    message Result {
        // The uniquely identifying components of the VASP in the directory service
        string id = 1;
        string registered_directory = 2;
        string common_name = 3;

        // Address to connect to the remote VASP on to perform a TRISA request
        string endpoint = 4;
    }

    Error error = 1;
    repeated Result results = 2;
}
```


## `Verification`

The `Verification` RPC enables VASPs to check on the status of a VASP, including its verification status. If the queried TRISA directory service performs health check monitoring, this RPC will also perform a health check and return the service status.

The `Verification` RPC expects as input a `VerificationRequest` message, which provides either the VASP's unique ID in the Global Directory Service, or the domain of their TRISA implementation endpoint (*Note: this should be the commone name of their certificate). If both `id` and `common_name` are supplied, the `id` field will be prioritized.

The `registered_directory` field is the URL of the directory that registered the VASP (e.g. either `"trisatest.net"` or `"vaspdirectory.net"`). If omitted, this value defaults to the directory currently being queried.

```proto

message VerificationRequest {
    string id = 1;
    string registered_directory = 2;
    string common_name = 3;
}
```

{{% notice note %}}
While `VerificationRequest` expects the same parameters as the `LookupRequest` message described earlier in this page, a TRISA directory service may refuse to return all or part of the status request.
{{% /notice %}}

A `VerificationRequest` returns a `VerificationReply` message:

```proto
message VerificationReply {
    // Status information
    trisa.gds.models.v1beta1.VerificationState verification_status = 1;
    trisa.gds.models.v1beta1.ServiceState service_status = 2;
    string verified_on = 3;  // Should be an RFC 3339 Timestamp
    string first_listed = 4; // Should be an RFC 3339 Timestamp
    string last_updated = 5; // Should be an RFC 3339 Timestamp
    string revoked_on = 6; // Should be an RFC 3339 Timestamp
}
```

## `Status`

The `Status` RPC is designed to enable routine health checks to verify the health of the Directory Service.

This RPC expects a `HealthCheck` message and returns a `ServiceState` message. The system is obliged to respond with the closest matching `status` in a best-effort fashion. Alerts will be triggered on service status changes if the system does not respond and the previous system state was not `UNKNOWN`.:

```proto
message HealthCheck {
    // The number of failed health checks that proceeded the current check.
    uint32 attempts = 1;

    // The timestamp of the last health check, successful or otherwise.
    string last_checked_at = 2;
}

message ServiceState {
    enum Status {
        UNKNOWN = 0;
        HEALTHY = 1;
        UNHEALTHY = 2;
        DANGER = 3;
        OFFLINE = 4;
        MAINTENANCE = 5;
    }

    // Current service status as defined by the receiving system.
    Status status = 1;

    // Suggest to the directory service when to check the health status again.
    string not_before = 2;
    string not_after = 3;
}
```

## Programmatic Registration

Two of the `TRISADirectory` are designed to support programatic registration for third party service providers that would like to use the GDS to issue certs: `Register` and `VerifyContact`.

### `Register`

Registration requests submitted via the `Register` RPC are validated to ensure they contain correct information and then are sent through the verification process, creating or updating a VASP as needed.

The `Register` RPC expects a `RegisterRequest` message containing:

1. The `LegalPerson` IVMS 101 `entity` to enable VASP KYC information exchange. This is the IVMS 101 data that should be exchanged in the TRISA P2P protocol as the Originator, Intermediate, or Beneficiary VASP fields. A complete and valid identity record with country of registration is required.
2. Technical, legal, billing, and administrative `contacts` for the VASP.
3. The Travel Rule implementation `trisa_endpoint` where other TRISA peers should connect. This should be an `addr:port` combination, (e.g. `trisa.vaspbot.net:443`).
4. The VASP's `common_name`, which should be the domain name to issue certificates for and should match the domain in the `trisa_endpoint`. If this field is omitted, the `common_name` is inferred from the `trisa_endpoint`.
5. Business information including `website`, `business_category`, `vasp_categories`, and the company's `established_on` date (in YYYY-MM-DD format).
6. The VASP's `trixo` questionnaire data. For more information, see the [TRIXO documentation]({{% ref "/joining-trisa/trixo" %}}).

```proto
message RegisterRequest {
    ivms101.LegalPerson entity = 1;
    trisa.gds.models.v1beta1.Contacts contacts = 2;
    string trisa_endpoint = 3;
    string common_name = 4;
    string website = 5;
    trisa.gds.models.v1beta1.BusinessCategory business_category = 6;
    repeated string vasp_categories = 7;
    string established_on = 8;
    trisa.gds.models.v1beta1.TRIXOQuestionnaire trixo = 9;
}
```

A `RegisterRequest` returns a `RegisterReply` message containing verification metadata as well as a `pkcs12password` that must be used to decrypt the emailed certifications. For more information, see the [PKCS12 password documentation]({{% ref "/joining-trisa/pkcs12" %}}). Do not lose or share this password!

```proto
message RegisterReply {
     // If the registration was successful, no error will be returned
    Error error = 1;

    // Unique identifiers for the VASP created by the registration.
    // Use these identifiers for status lookup requests and any follow-on interactions.
    string id = 2;
    string registered_directory = 3;
    string common_name = 4;

    // The verification status of the VASP entity.
    trisa.gds.models.v1beta1.VerificationState status = 5;
    string message = 6;

    // Used to decrypt the emailed certificates in PKCS12 format
    string pkcs12password = 7;
}
```

### `VerifyContact`

The `VerifyContact` RPC expects a `VerifyContactRequest` contains the VASP's unique ID in the Directory Service. The RPC returns as `VerifyContactReply` with the verification status of the queried VASP.

```proto
message VerifyContactRequest {
    string id = 1;
    string token = 2;
}

message VerifyContactReply {
    // If no error is specified, the verification request was successful.
    Error error = 1;

    // The verification status of the VASP entity.
    trisa.gds.models.v1beta1.VerificationState status = 2;
    string message = 3;
}
```