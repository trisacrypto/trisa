package client

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/trisacrypto/trisa/pkg/openvasp"
	"github.com/trisacrypto/trisa/pkg/openvasp/extensions/discoverability"
	"github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

// New creates a TRP client that uses travel addresses to send TRP protocol requests.
// The client implements the TRP v3 API protocol as well as the Discoverability
// extension for querying the version, uptime, and extensions of a TRP server.
func New(opts ...ClientOption) (client *Client, err error) {
	// Create the default client before applying options
	client = &Client{
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Timeout:       30 * time.Second,
		},
		apiVersion: openvasp.APIVersion,
		extensions: nil,
	}

	// Create cookie jar
	if client.client.Jar, err = cookiejar.New(nil); err != nil {
		return nil, fmt.Errorf("could not create cookiejar: %w", err)
	}

	// Apply options and configuration from user.
	// Options can expect that a client has already been configured with defaults.
	for _, opt := range opts {
		if err = opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

type Client struct {
	client     *http.Client
	apiVersion string
	extensions []string
}

// Ensure the Client implements the TRPv3 and Discoverability Interfaces
var (
	_ trp.Client             = &Client{}
	_ discoverability.Client = &Client{}
)

// Returns the default API version used in requests.
func (c *Client) APIVersion() string {
	return c.apiVersion
}

// Returns the default API extensions used in requests.
func (c *Client) APIExtensions() []string {
	return c.extensions
}

//===========================================================================
// Client Methods - TRPv3
//===========================================================================

func (c *Client) Identity(ctx context.Context, address string) (out *trp.Identity, err error) {
	info := &trp.Info{
		Address:           address,
		APIVersion:        c.apiVersion,
		RequestIdentifier: "",
		APIExtensions:     c.extensions,
	}

	// Convert the endpoint to the Identity endpoint specified
	var uri *url.URL
	if uri, err = info.ParseURL(); err != nil {
		return nil, err
	}

	uri.Path = trp.IdentityEndpoint
	info.Address = uri.String()

	out = &trp.Identity{}
	if _, err = c.Get(ctx, info, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Inquiry(ctx context.Context, in *trp.Inquiry) (out *trp.Resolution, err error) {
	// Ensure the the travel rule query parameter is set
	var uri *url.URL
	if uri, err = in.Info.ParseURL(); err != nil {
		return nil, err
	}

	// Parse the URL to add the query parameters
	query := uri.Query()
	query.Set("t", "i")
	uri.RawQuery = query.Encode()
	in.Info.Address = uri.String()

	// Send the inquiry request
	out = &trp.Resolution{}
	if _, err = c.Post(ctx, in.Info, in, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) Resolve(ctx context.Context, in *trp.Resolution) (err error) {
	// Expects a 204 response
	if _, err = c.Post(ctx, in.Info, in, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) Confirm(ctx context.Context, in *trp.Confirmation) (err error) {
	// Expects a 204 response
	if _, err = c.Post(ctx, in.Info, in, nil); err != nil {
		return err
	}
	return nil
}

//===========================================================================
// Client Methods - Discoverability
//===========================================================================

func (c *Client) Version(ctx context.Context, address string) (out *discoverability.Version, err error) {
	info := &trp.Info{
		Address:           address,
		APIVersion:        c.apiVersion,
		RequestIdentifier: "",
		APIExtensions:     c.extensions,
	}

	// Convert the endpoint to the Identity endpoint specified
	var uri *url.URL
	if uri, err = info.ParseURL(); err != nil {
		return nil, err
	}

	uri.Path = discoverability.VersionEndpoint
	info.Address = uri.String()

	out = &discoverability.Version{}
	if _, err = c.Get(ctx, info, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Uptime(ctx context.Context, address string) (out discoverability.Uptime, err error) {
	info := &trp.Info{
		Address:           address,
		APIVersion:        c.apiVersion,
		RequestIdentifier: "",
		APIExtensions:     c.extensions,
	}

	// Convert the endpoint to the Identity endpoint specified
	var uri *url.URL
	if uri, err = info.ParseURL(); err != nil {
		return 0, err
	}

	uri.Path = discoverability.UptimeEndpoint
	info.Address = uri.String()

	var req *http.Request
	if req, err = c.NewTextRequest(ctx, http.MethodGet, info, nil); err != nil {
		return 0, err
	}

	if _, err = c.Do(req, &out); err != nil {
		return 0, err
	}
	return out, nil
}

func (c *Client) Extensions(ctx context.Context, address string) (out *discoverability.Extensions, err error) {
	info := &trp.Info{
		Address:           address,
		APIVersion:        c.apiVersion,
		RequestIdentifier: "",
		APIExtensions:     c.extensions,
	}

	// Convert the endpoint to the Identity endpoint specified
	var uri *url.URL
	if uri, err = info.ParseURL(); err != nil {
		return nil, err
	}

	uri.Path = discoverability.ExtensionsEndpoint
	info.Address = uri.String()

	out = &discoverability.Extensions{}
	if _, err = c.Get(ctx, info, out); err != nil {
		return nil, err
	}
	return out, nil
}

//===========================================================================
// Client Methods
//===========================================================================

// Executes an HTTP GET request expecting a JSON response (out) from the server. If
// out is nil or the status code is 204 then the body response will not be parsed.
func (s *Client) Get(ctx context.Context, ta *trp.Info, out interface{}) (rep *http.Response, err error) {
	var req *http.Request
	if req, err = s.NewJSONRequest(ctx, http.MethodGet, ta, nil); err != nil {
		return nil, err
	}
	return s.Do(req, out)
}

// Executes an HTTP POST request expecting a JSON payload (in) and accepts a JSON
// payload (out) from the server (though the response will be handled based on the
// content-type the server responds with). Both in and out can be nil to skip handling.
func (s *Client) Post(ctx context.Context, ta *trp.Info, in, out interface{}) (rep *http.Response, err error) {
	var req *http.Request
	if req, err = s.NewJSONRequest(ctx, http.MethodPost, ta, in); err != nil {
		return nil, err
	}
	return s.Do(req, out)
}

// Executes an http request, performs error checking, and deserializes the response
// data into the specified out interface. The out interface should implement either
// encoding.TextUnmarshaler in the case of text/plain responses or will be decoded using
// json in the case of application/json responses.
func (s *Client) Do(req *http.Request, out interface{}) (rep *http.Response, err error) {
	if rep, err = s.client.Do(req); err != nil {
		return nil, fmt.Errorf("could not execute request: %w", err)
	}
	defer rep.Body.Close()

	// Attempt to parse the content type from the response
	var contentType, mediaType string
	if contentType = rep.Header.Get(openvasp.ContentTypeHeader); contentType != "" {
		mediaType, _, _ = mime.ParseMediaType(contentType)
	}

	// Handle Status Errors
	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		// Create a status error from the response
		serr := &trp.StatusError{
			Code:    rep.StatusCode,
			Message: http.StatusText(rep.StatusCode),
		}

		// Attempt to deserialize the status error from the body
		switch mediaType {
		case openvasp.MIMEJSON:
			// Unmarshal a JSON response from the server
			jerr := &trp.StatusError{}
			if err = json.NewDecoder(rep.Body).Decode(jerr); err == nil {
				serr = jerr
			}

		case openvasp.MIMEPlainText:
			// Unmarshal a text response from the server
			var msg []byte
			if msg, err = io.ReadAll(rep.Body); err != nil {
				// If we can't read the body, just return the status error.
				return nil, serr
			}

			if msgs := strings.TrimSpace(string(msg)); msgs != "" {
				serr.Message = msgs
			}
		}

		return rep, serr
	}

	// Read the body into the specified interface.
	if out != nil && rep.StatusCode != http.StatusNoContent {
		switch mediaType {
		case openvasp.MIMEJSON:
			if err = json.NewDecoder(rep.Body).Decode(out); err != nil {
				return nil, fmt.Errorf("could not deserialize json response: %w", err)
			}

		case openvasp.MIMEPlainText:
			var data []byte
			if data, err = io.ReadAll(rep.Body); err != nil {
				return nil, fmt.Errorf("could not read response body: %w", err)
			}

			switch v := out.(type) {
			case encoding.TextUnmarshaler:
				if err = v.UnmarshalText(data); err != nil {
					return nil, fmt.Errorf("could not unmarshal text response: %w", err)
				}
			case encoding.BinaryUnmarshaler:
				if err = v.UnmarshalBinary(data); err != nil {
					return nil, fmt.Errorf("could not unmarshal binary response: %w", err)
				}
			default:
				return nil, fmt.Errorf("could not unmarshal text response to %T", v)
			}

		case "":
			return nil, fmt.Errorf("could not identify media type from response header %q", contentType)

		default:
			return nil, fmt.Errorf("unsupported content type: %s", mediaType)
		}
	}

	return rep, nil
}
