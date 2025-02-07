package client

import (
	"context"

	"github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

type Client struct{}

func (c *Client) Inquiry(ctx context.Context, in *trp.Inquiry) (*trp.Resolution, error) {
	// Ensure the the travel rule query parameter is set
	// query := u.Query()
	// query.Set("t", "i")

	// u.RawQuery = query.Encode()
	return nil, nil
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
