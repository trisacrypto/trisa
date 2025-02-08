package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
		APIVersion:        openvasp.APIVersion,
		RequestIdentifier: "",
		APIExtensions:     nil,
	}

	if _, err = info.URL(); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) Inquiry(ctx context.Context, in *trp.Inquiry) (out *trp.Resolution, err error) {
	// Ensure the the travel rule query parameter is set
	var uri string
	if uri, err = in.Info.URL(); err != nil {
		return nil, err
	}

	// Parse the URL to add the query parameters
	u, _ := url.Parse(uri)
	query := u.Query()
	query.Set("t", "i")
	u.RawQuery = query.Encode()
	in.Info.Address = u.String()

	return out, nil
}

func (c *Client) Resolve(ctx context.Context, in *trp.Resolution) (err error) {
	return nil
}

func (c *Client) Confirm(ctx context.Context, in *trp.Confirmation) (err error) {
	return nil
}

//===========================================================================
// Client Methods - Discoverability
//===========================================================================

func (c *Client) Version(ctx context.Context, address string) (out *discoverability.Version, err error) {
	return out, nil
}

func (c *Client) Uptime(ctx context.Context, address string) (out discoverability.Uptime, err error) {
	return out, nil
}

func (c *Client) Extensions(ctx context.Context, address string) (out *discoverability.Extensions, err error) {
	return out, nil
}

// // ParseAddress parses a travel rule address, LNURL, or URL.
// func ParseAddress(address string) (*url.URL, error) {
// 	switch {
// 	case strings.HasPrefix(address, "lnurl1"), strings.HasPrefix(address, "LNURL1"):
// 		return lnurl.Decode(address)

// 	case strings.HasPrefix(address, "ta"):
// 		return traddr.DecodeURL(address)

// 	default:
// 		var u *url.URL
// 		if u, err = url.Parse(address); err != nil {
// 			return "", err
// 		}

// 		if u.Scheme != "" {
// 			query := u.Query()
// 			query.Set("t", "i")
// 			u.RawQuery = query.Encode()
// 			return u.String(), nil
// 		}
// 	}
// 	return nil, trp.ErrUnknownTravelAddress
// }
