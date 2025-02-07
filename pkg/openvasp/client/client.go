package client

import (
	"context"
	"net/url"

	"github.com/trisacrypto/trisa/pkg/openvasp"
	"github.com/trisacrypto/trisa/pkg/openvasp/extensions/discoverability"
	"github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

type Client struct{}

var (
	_ trp.Client             = &Client{}
	_ discoverability.Client = &Client{}
)

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
	var uri *url.URL
	if uri, err = in.Info.URL(); err != nil {
		return nil, err
	}

	query := uri.Query()
	query.Set("t", "i")
	uri.RawQuery = query.Encode()
	in.Info.Address = uri.String()

	return out, nil
}

func (c *Client) Resolve(ctx context.Context, in *trp.Resolution) (err error) {
	return nil
}

func (c *Client) Confirm(ctx context.Context, in *trp.Confirmation) (err error) {
	return nil
}

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
