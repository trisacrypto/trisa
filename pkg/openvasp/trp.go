package openvasp

import (
	"github.com/trisacrypto/trisa/pkg/ivms101"
	"github.com/trisacrypto/trisa/pkg/slip0044"
)

// OpenVASP Application Headers
const (
	APIVersionHeader        = "api-version"
	APIExtensionsHeader     = "api-extensions"
	RequestIdentifierHeader = "request-identifier"
	ContentTypeHeader       = "content-type"
)

// OpenVASP Application Header Values
const (
	APIVersion       = "3.1.0"
	ContentTypeValue = "application/json; charset=utf-8"
	ContentMediaType = "application/json"
)

// TRPInfo contains metadata information from the TRP API Headers.
type TRPInfo struct {
	LNURL             string
	APIVersion        string
	RequestIdentifier string
	APIExtensions     []string
}

// Inquiry defines a Travel Rule Protocol payload that contains information about the
// transaction and the originator and beneficiary of the transaction.
type Inquiry struct {
	TRP        *TRPInfo                `json:"-"`
	Asset      slip0044.CoinType       `json:"asset"`
	Amount     float64                 `json:"amount"`
	Callback   string                  `json:"callback"`
	IVMS101    ivms101.IdentityPayload `json:"IVMS101"`
	Extensions map[string]interface{}  `json:"extensions,omitempty"`
}

// The TransactionPayload extension is used to provide information about the
// transaction on the blockchain or network so it can be linked to the identity
// information in the TRP payload.
type TransactionPayload struct {
	// Transaction ID on the blockchain or network.
	TxID string `json:"txid,omitempty"`

	// Crypto address of the originator and beneficiary.
	Originator  string `json:"originator,omitempty"`
	Beneficiary string `json:"beneficiary,omitempty"`

	// Amount and asset type of the transaction.
	Amount    float64 `json:"amount,omitempty"`
	AssetType string  `json:"asset_type,omitempty"`

	// The blockchain or network of the transaction.
	Network string `json:"network,omitempty"`

	// The RFC3339 timestamp of the transaction.
	Timestamp string `json:"timestamp,omitempty"`

	// Tags and extra JSON data about the transaction.
	Tag       string                 `json:"tag,omitempty"`
	ExtraJSON map[string]interface{} `json:"extra_json,omitempty"`
}

// The SealedTRISAEnvelope extension is used to faciliate the TRISA protocol by providing a
// JSON serialized version of the secure envelope that contains the transaction.
type SealedTRISAEnvelope struct {
	Envelope string `json:"envelope"`
}

// The UnsealedTRISAEnvelope extension is used to provide an unsealed version of a
// secure envelope where the key is unencrypted, allowing any party to decrypt the
// payload and access the identity and transaction information.
type UnsealedTRISAEnvelope struct {
	// Transaction ID generated by the originator.
	Id string `json:"id"`

	// The encrypted payload containing the IVMS101 identity and generic transaction.
	Payload []byte `json:"payload"`

	// Encryption key used to encrypt the payload, in this struct the key is unencrypted.
	EncryptionKey       []byte `json:"encryption_key"`
	EncryptionAlgorithm string `json:"encryption_algorithm"`

	// HMAC of the payload to ensure integrity.
	HMAC          []byte `json:"hmac"`
	HMACSecret    []byte `json:"hmac_secret"`
	HMACAlgorithm string `json:"hmac_algorithm"`
}

// InquiryResolution is used to approve or reject a TRP Transfer Inquiry either
// automatically in direct response to the inquiry request or via the callback URL
// specified in the request. One of "approved", "rejected", or "version" should be
// specified to ensure unambiguous results are returned to the caller.
type InquiryResolution struct {
	Version  string    `json:"version,omitempty"`  // the API version of the request
	Approved *Approval `json:"approved,omitempty"` // payment address and callback
	Rejected string    `json:"rejected,omitempty"` // human readable comment (must be specified to reject)
}

// Approval is used to accept a TRP Transfer Inquiry.
type Approval struct {
	Address  string `json:"address"`  // some payment address
	Callback string `json:"callback"` // some implementation defined URL for transfer confirmation
}

// Confirmation JSON data is sent in response to a TransferInquiry via a POST to the
// callback URl. Only one of txid or canceled should be specified. The txid should be
// specified only if the transaction has been broadcasted. Canceled is used to indicate
// that the transfer will not move forward with a human readable comment.
type Confirmation struct {
	TRP      *TRPInfo `json:"-"`
	TXID     string   `json:"txid,omitempty"`     // some asset-specific tx identifier
	Canceled string   `json:"canceled,omitempty"` // human readable comment or null
}

func (c Confirmation) Validate() error {
	if c.TXID == "" && c.Canceled == "" {
		return ErrEmptyConfirmation
	}

	if c.TXID != "" && c.Canceled != "" {
		return ErrAmbiguousConfirmation
	}

	return nil
}
