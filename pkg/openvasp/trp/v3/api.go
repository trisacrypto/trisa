/*
See: https://gitlab.com/OpenVASP/travel-rule-protocol/-/blob/master/core/specification.md
*/
package trp

import (
	"context"
	"net/url"
	"strings"

	"github.com/trisacrypto/trisa/pkg/ivms101"
	"github.com/trisacrypto/trisa/pkg/openvasp/lnurl"
	"github.com/trisacrypto/trisa/pkg/openvasp/traddr"
	"github.com/trisacrypto/trisa/pkg/slip0044"
)

type Client interface {
	Identity(ctx context.Context, address string) (*Identity, error)
	Inquiry(ctx context.Context, in *Inquiry) (*Resolution, error)
	Resolve(ctx context.Context, in *Resolution) error
	Confirm(ctx context.Context, in *Confirmation) error
}

//===========================================================================
// TRP Info Headers
//===========================================================================

// Info describes the TRP specific metadata that are in the headers of the request.
// This data is not serialized with the JSON payload.
type Info struct {
	Address           string   `json:"-"` // Address can be a Travel Rule Address, LNURL, or URL
	APIVersion        string   `json:"-"` // Defaults to the APIVersion of the package
	RequestIdentifier string   `json:"-"` // A unique identifier representing the specific transfer
	APIExtensions     []string `json:"-"` // The names of any extensions uses in the request
}

func (i *Info) URL() (_ *url.URL, err error) {
	switch {
	case i.Address == "":
		return nil, ErrUnknownTravelAddress

	case strings.HasPrefix(i.Address, "lnurl1"), strings.HasPrefix(i.Address, "LNURL1"):
		var uri string
		if uri, err = lnurl.Decode(i.Address); err != nil {
			return nil, err
		}
		return url.Parse(uri)

	case strings.HasPrefix(i.Address, "ta"):
		var uri string
		if uri, err = traddr.DecodeURL(i.Address); err != nil {
			return nil, err
		}
		return url.Parse(uri)

	default:
		// Attempt to parse the URL
		var u *url.URL
		if u, err = url.Parse(i.Address); err != nil {
			return nil, err
		}

		// The URL must have a scheme to be valid
		if u.Scheme != "" {
			return u, nil
		}
	}

	return nil, ErrUnknownTravelAddress
}

//===========================================================================
// Identity Discovery
//===========================================================================

// Identity is returned on the identity discoverability endpoint.
// This is a top-level data structure: it is either sent or received.
type Identity struct {
	Info *Info  `json:"-"`              // Request metadata and headers
	Name string `json:"name,omitempty"` // company name as incorporated in the commercial registry as a string
	LEI  string `json:"lei,omitempty"`  // the Legal Entity Indentifier (see gleif.org) as string
	X509 string `json:"x509,omitempty"` // x509 certificate in PEM format
}

//===========================================================================
// Inquiries
//===========================================================================

// Inquiry defines a Travel Rule Protocol payload that contains information about the
// transaction and the originator and beneficiary of the transaction.
// This is a top-level data structure: it is either sent or received.
type Inquiry struct {
	Info       *Info                    `json:"-"` // Request metadata and headers
	Asset      *Asset                   `json:"asset"`
	Amount     float64                  `json:"amount"`
	Callback   string                   `json:"callback"`
	IVMS101    *ivms101.IdentityPayload `json:"IVMS101"`
	Extensions map[string]interface{}   `json:"extensions,omitempty"`
}

// Asset is used to describe the virtual or crypto asset being transferred.
// This is a sub-data structure: it is nested within an Inquiry.
type Asset struct {
	DTI     string            `json:"dti,omitempty"`      // digital token identifier as per Digital Token Identifier Foundation
	SLIP044 slip0044.CoinType `json:"slip0044,omitempty"` // registered coin types defined by BIP-0044
}

//===========================================================================
// Resolutions
//===========================================================================

// A response to a Travel Rule Inquiry either directly to a request or as a callback.
// This is a top-level data structure: it is either sent or received.
type Resolution struct {
	Info     *Info     `json:"-"`                  // Request metadata and headers
	Version  string    `json:"version,omitempty"`  // the API version of the request
	Approved *Approval `json:"approved,omitempty"` // payment address and callback
	Rejected string    `json:"rejected,omitempty"` // human readable comment (must be specified to reject)
}

// Approval is used to accept a TRP Transfer Inquiry.
// This is a sub-data structure: it is nested within a Resolution.
type Approval struct {
	Address  string `json:"address"`  // some payment address
	Callback string `json:"callback"` // some implementation defined URL for transfer confirmation
}

//===========================================================================
// Confirmation
//===========================================================================

// Confirmation is sent after an approval callback to indicate that the transaction
// has taken place on the chain or has been canceled.
// This is a top-level data structure: it is either sent or received.
type Confirmation struct {
	Info     *Info  `json:"-"`                  // Request metadata and headers
	TXID     string `json:"txid,omitempty"`     // some asset-specific tx identifier
	Canceled string `json:"canceled,omitempty"` // human readable comment or null
}
