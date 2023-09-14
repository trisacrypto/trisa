---
title: "TRP Bridge"
date: 2023-09-13T16:49:59-05:00
lastmod: 2023-09-13T16:49:59-05:00
description: "Integrating a TRP server into your TRISA node"
weight: 10
---

The TRP Bridge implements [`http.Handler`](https://pkg.go.dev/net/http#Handler) objects for each TRP request type and which can be passed to a basic [`http.Server`](https://pkg.go.dev/net/http#Server) and easily built into your TRISA node. A minimal example that simply logs incoming requests and return no error is below:

```golang
package main

import (
	"log"
	"net/http"

	"github.com/trisacrypto/trisa/pkg/openvasp"
)

// TRPHandler implements both the InquiryHandler and the Confirmation Handler interfaces
type TRPHandler struct{}


// OnInquiry implements the InquiryHandler interface and is used to respond to TRP
// transfer inquiry requests that initiate or conclude the TRP protocol.
func (t *TRPHandler) OnInquiry(in *openvasp.Inquiry) (*openvasp.InquiryResolution, error) {
    log.Printf(
        "received trp inquiry with request identifier %q\n",
        in.TRP.RequestIdentifier
    )
    return nil, nil
}

// OnConfirmation implements the ConfirmationHandler interface and is used to respond to
// TRP callbacks from the beneficiary VASP.
func (t *TRPHandler) OnConfirmation(in *openvasp.Confirmation) error {
    log.Printf(
        "received trp confirmation with request identifier %q\n",
        in.TRP.RequestIdentifier
    )
    return nil
}

func main() {
    // Create a new handler object
    handler := &TRPHandler{}

    // Create a mux to route requests to different paths to different handlers.
    mux := http.NewServeMux()
    mux.Handle("/transfers", openvasp.TransferInquiry(handler))
    mux.Handle("/confirm", openvasp.TransferConfirmation(handler))

    log.Printf(
        "waiting for TRP requests with API version %s at http://localhost:8080\n",
        openvasp.APIVersion
    )

    // Serve the TRP API server on port 8080. In production applications you would
    // likely configure mTLS and TLS termination at the server first.
    http.ListenAndServe(":8080", mux)
}
```

## Transfer Inquiry

The [`openvasp.TransferInquiry`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#TransferInquiry) function accepts an object that implements the [`openvasp.InquiryHandler`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#InquiryHandler) interface and returns an [`http.Handler`](https://pkg.go.dev/net/http#Handler) that wraps the [`InquiryHandler`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#InquiryHandler) to handle [TRP Transfer Inquiry `POST` requests](https://gitlab.com/OpenVASP/travel-rule-protocol/-/blob/master/core/specification.md?ref_type=heads#detailed-protocol-flow). The [`InquiryHandler`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#InquiryHandler) is defined as follows:

```golang
type InquiryHandler interface {
	OnInquiry(*Inquiry) (*InquiryResolution, error)
}
```

The HTTP handler returned performs the following operations when an incoming HTTP `POST` request is received:

1. Validates the incoming TRP request
2. Parses the TRP [`Inquiry`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#Inquiry) object along with any extensions.
3. Calls the handlers `OnInquiry` method
4. Returns success or failure based on the returned value of `OnInquiry`.

To return a failure condition from the `OnInquiry()` function, users may return an [`openvasp.StatusError`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#StatusError) that specifies the HTTP status code and message to return. This is useful particularly to return `404` errors if the Travel Address is incorrect or no beneficiary account exists at the endpoint. If a generic `error` is returned, then the handler will return a [500 Internal Server Error](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/500) along with the `err.Error()` text.

```golang
func (t *TRPHandler) OnInquiry(in *openvasp.Inquiry) (*openvasp.InquiryResolution, error) {

    // Lookup Travel Address and return a 404 error if beneficiary is not found.
    if notFound {
        return nil, &openvasp.StatusError{Code: http.StatusNotFound}
    }

}
```

To return a succesful [200 OK](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/200) response, return an [`openvasp.InquiryResolution`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#InquiryResolution) object or simply `nil, nil`.

The semantics are as follows:

- If the resolution response is `nil` or contains just the `Version`, then the counterparty expects a subsequent `POST` request to the callback in the request.
- The inquiry can be automatically approved by returning the `Approved` field without the `Version` or the `Rejected` fields (these must be zero valued).
- The inquiry can be automatically rejected by specifying a `Rejected` reason without the `Version` or the `Approved` fields (these must be zero valued).


## Transfer Confirmation


## Setting up a Server

## API Checks

The [`openvasp.APIChecks`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#Inquiry) middleware validates TRP requests and parses header information from the request. Both the `TransferInquiry` and the `TransferConfirmation` handlers implement this middleware.

The checks that are performed include:

1. Ensure the HTTP Method is `POST` otherwise a [405 Method Not Allowed](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/405) error is returned.
2. Ensure the TRP API Version header is set and the version is compatible with the implemented TRP version otherwise a [400 Bad Request](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/400) error is returned.
3. Ensures that the currently implemented TRP API Version header is set on the outgoing response.
4. Checks that there is a request identifier header on the request, otherwise a [400 Bad Request](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/400) error is returned.
5. Ensures that the request identifier is echoed back on the outgoing response.
6. Enforces that the `Content-Type` header is specified and that the content type is `application/json`, otherwise a [415 Unsupported Media Type](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/415) error is returned.

## ParseTRPInfo

The [`openvasp.ParseTRPInfo`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#ParseTRPInfo) function parses TRP-specific headers from the request and adds them to a [`openvasp.TRPInfo`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/openvasp#TRPInfo) struct that can be used for TRP processing. A mapping of the headers to the parsed fields is as follows:

| Field             | Header             | Type     | Description                                                                      |
|-------------------|--------------------|----------|----------------------------------------------------------------------------------|
| Address           |                    | string   | The Travel Address, LNURL, or URL of the request                                 |
| APIVersion        | api-version        | string   | Defaults to the APIVersion of the package                                        |
| RequestIdentifier | api-extensions     | string   | A unique identifier representing the specific transfer (used as the envelope ID) |
| APIExtensions     | request-identifier | []string | The comma separated names of any extensions used in the request                  |