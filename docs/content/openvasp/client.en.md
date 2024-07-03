---
title: "TRP Client"
date: 2023-09-13T16:49:59-05:00
lastmod: 2023-09-13T16:49:59-05:00
description: "Making requests to other TRP nodes"
weight: 20
---

TRISA provides an [`openvasp.Client`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#Client) to make communicating with counterparty VASPs simpler and ensure that the TRP protocol is being completed successfully. To create a client:

```golang
client := openvasp.NewClient()
```

{{% notice note %}}
Client options such as mTLS configuration are coming soon!
{{% /notice %}}

## Inquiries

To send a Travel Rule inquiry, use the [`client.Inquiry()`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#Client.Inquiry) method.

The client parses the routing for the `Inquiry` using the `Inquiry.TRP` field; at the very least, the `Address` and `RequestIdentifier` fields must be specified:

```golang
type TRPInfo struct {
	Address           string   // Address can be a Travel Rule Address, LNURL, or URL
	APIVersion        string   // Defaults to the APIVersion of the package
	RequestIdentifier string   // A unique identifier representing the specific transfer
	APIExtensions     []string // The names of any extensions uses in the request
}
```

The `APIVersion` field is automatically populated with the default version, and any extensions are populated from the `Inquiry` itself, so these can be ignored.

The `Address` field can be one of:

- [LNURL](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp/lnurl)
- [Travel Address](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp/traddr)
- Plain HTTPS URL

Note that the TRISA library has packages for creating and parsing LNURLs and Travel Addresses so that you do not have to import additional dependencies for your code.

The client returns a [`TravelRuleResponse`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#TravelRuleResponse), which you can use as though it were a regular `http.Response` object, including use cases like checking the `http.StatusCode`. However, you can also parse the `InquiryResolution` response as follows:

```golang
func main() {
    client := openvasp.NewClient()

    inquiry := &openvasp.Inquiry{
        TRP: &openvasp.TRP{
            Address: "ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
            RequestIdentifier: "129f6013-9125-4beb-9e86-bb20c440e164"
        },
        Asset: &openvasp.Asset{},
        Amount: 0.001,
        Callback: "https://originator.com/confirm?t=i",
        IVMS101: &ivms101.IdentityPayload{},
    }

    rep, err := client.Inquiry(inquiry)
    if err != nil {
        log.Fatal(err)
    }

    if rep.StatusCode == http.StatusOK {
        resolution, err := rep.InquiryResolution()
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println(resolution)

    } else {
        log.Fatal(rep.Status)
    }

}
```

## Confirmation

To send a Transfer Confirmation, use the [`client.Confirmation()`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#Client.Confirmation) method.

The client parses the routing for the `Confirmation` using the `Confirmation.TRP` field; at the very least, the `Address` and `RequestIdentifier` fields must be specified:

```golang
type TRPInfo struct {
	Address           string   // Address can be a Travel Rule Address, LNURL, or URL
	APIVersion        string   // Defaults to the APIVersion of the package
	RequestIdentifier string   // A unique identifier representing the specific transfer
	APIExtensions     []string // The names of any extensions uses in the request
}
```

The `APIVersion` field is automatically populated with the default version, and any extensions are populated from the `Confirmation` itself, so these can be ignored.

The client returns a `TravelRuleResponse` however, a simple [204 No Content](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/204) response is expected for a confirmation, so there is no JSON body to parse.

```golang
func main() {
    client := openvasp.NewClient()

    confirm := &openvasp.Confirmation{
        TRP: &openvasp.TRP{
            Address: "https://beneficiary.com/confirm?t=i",
            RequestIdentifier: "129f6013-9125-4beb-9e86-bb20c440e164",
        },
        Canceled: "the transaction could not be completed",
    }

    rep, err := client.Confirmation(inquiry)
    if err != nil {
        log.Fatal(err)
    }

    if rep.StatusCode == http.StatusNoContent {
        fmt.Println("success!")
    } else {
        log.Fatal(rep.Status)
    }
}
```