package client

import (
	"bytes"
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/trisacrypto/trisa/pkg/openvasp"
	"github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

const (
	UserAgent      = "TRISA TRP Bridge Client/v1"
	Accept         = "application/json"
	AcceptLanguage = "en-US,en"
	AcceptEncode   = "gzip, deflate, br"
)

// Creates a new http request for use to a TRP server with the appropriate headers set.
// The headers are set by the info packet or by the defaults in the current version of
// this package if not set. The url the request is set to is extracted from the address
// info and can be either a travel rule address, an LNURL, or a standard URL.
//
// For convenience the request identifier and a tracing id can be set in the context
// of the request. However, the request identifier in the address info will take
// priority. If there is no request identifier specified, then a UUID is generated.
//
// If the body is not nil, this method sets the content type header to application/json
// as most requests require JSON payloads. If you need to send a different content type,
// ensure you reset the header after creating the request.
func NewRequest(ctx context.Context, method string, ta *trp.Info, body io.Reader) (req *http.Request, err error) {
	var url string
	if url, err = ta.URL(); err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	// Set defaults on the info packet as required
	if ta.APIVersion == "" {
		ta.APIVersion = openvasp.APIVersion
	}

	// Get the request identifier from the context or generate a new one if one is not set.
	if ta.RequestIdentifier == "" {
		if ta.RequestIdentifier, _ = RequestIDFromContext(ctx); ta.RequestIdentifier == "" {
			ta.RequestIdentifier = uuid.NewString()
		}
	}

	// Create the http request
	if req, err = http.NewRequestWithContext(ctx, method, url, body); err != nil {
		return nil, fmt.Errorf("could not create request: %s", err)
	}

	// Client specific headers
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", Accept)
	req.Header.Set("Accept-Language", AcceptLanguage)
	req.Header.Set("Accept-Encoding", AcceptEncode)

	// OpenVASP required headers
	req.Header.Set(openvasp.APIVersionHeader, ta.APIVersion)
	req.Header.Set(openvasp.RequestIdentifierHeader, ta.RequestIdentifier)

	if body != nil {
		req.Header.Set(openvasp.ContentTypeHeader, openvasp.ContentTypeValue)
	}

	if len(ta.APIExtensions) > 0 {
		req.Header.Set(openvasp.APIExtensionsHeader, strings.Join(ta.APIExtensions, ", "))
	}

	// If there is a tracing ID on the context, set it on the request
	var tracingID string
	if tracingID, _ = TracingFromContext(ctx); tracingID != "" {
		req.Header.Add("X-Request-ID", tracingID)
	}

	return req, nil
}

// Creates a new request with the default headers set and the client's API version and
// extensions. Use the NewJSONRequest and NewTextRequest methods handle specific data
// types and body encoding. See NewRequest for more information on the headers and
// context. This method also manages cookies from the cookiejar if needed for CSRF.
func (c *Client) NewRequest(ctx context.Context, method string, ta *trp.Info, body io.Reader) (req *http.Request, err error) {
	// Set defaults on the info packet as required
	if ta.APIVersion == "" {
		ta.APIVersion = c.apiVersion
	}

	if len(ta.APIExtensions) == 0 {
		ta.APIExtensions = c.extensions
	}

	if req, err = NewRequest(ctx, method, ta, body); err != nil {
		return nil, err
	}

	// Add CSRF protection if its available
	if c.client.Jar != nil {
		cookies := c.client.Jar.Cookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "csrf_token" {
				req.Header.Add("X-CSRF-TOKEN", cookie.Value)
			}
		}
	}

	return req, nil
}

// Create a new JSON payload, marshaling the given data into a JSON object.
func (c *Client) NewJSONRequest(ctx context.Context, method string, ta *trp.Info, data interface{}) (req *http.Request, err error) {
	var body io.ReadWriter
	switch {
	case data == nil:
		body = nil
	default:
		body = &bytes.Buffer{}
		if err = json.NewEncoder(body).Encode(data); err != nil {
			return nil, fmt.Errorf("could not serialize request data as json: %s", err)
		}
	}

	// NOTE: the NewRequest function sets the Accept and Content-Type headers to JSON.
	return c.NewRequest(ctx, method, ta, body)
}

// Create a new text payload, wrapping the given text unmarshaler in a reader.
func (c *Client) NewTextRequest(ctx context.Context, method string, ta *trp.Info, data interface{}) (req *http.Request, err error) {
	var body io.ReadWriter
	switch v := data.(type) {
	case encoding.TextMarshaler:
		var text []byte
		if text, err = v.MarshalText(); err != nil {
			return nil, fmt.Errorf("could not marshal text: %s", err)
		}
		body = bytes.NewBuffer(text)
	case fmt.Stringer:
		body = bytes.NewBufferString(v.String())
	case string:
		body = bytes.NewBufferString(v)
	case []byte:
		body = bytes.NewBuffer(v)
	default:
		body = bytes.NewBufferString(fmt.Sprintf("%v", v))
	}

	if req, err = c.NewRequest(ctx, method, ta, body); err != nil {
		return nil, err
	}

	// Ensure that the content type and accept headers are set to text/plain.
	req.Header.Set(openvasp.ContentTypeHeader, openvasp.MIMEPlainText)
	req.Header.Set("Accept", openvasp.MIMEPlainText)

	return req, nil
}
