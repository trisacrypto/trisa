package keys

import (
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/keys/signature"
	"google.golang.org/protobuf/proto"
)

// Create a Key from an *trisa.SigningKey received during the trisa.KeyExchange RPC. If
// the signing key contains an invalid or unparseable public key an error is returned.
func FromSigningKey(msg *api.SigningKey) (Key, error) {
	exchange := &Exchange{msg: msg}
	if _, err := exchange.SealingKey(); err != nil {
		return nil, err
	}
	return exchange, nil
}

// Exchange wraps a *trisa.SigningKey protocol buffer that is sent during a
// trisa.KeyExchange RPC and implements the Key interface. Exchange Keys are public keys
// only (they do not and cannot contain a private key) and are only used for sealing
// secure envelopes before a trisa.Transfer RPC.
type Exchange struct {
	msg    *api.SigningKey
	pubkey interface{}
}

// Ensure the Exchange implements the  Key interface
var _ Key = &Exchange{}

// IsPrivate always returns false for an Exchange key - no private keys are available.
func (e *Exchange) IsPrivate() bool {
	return false
}

// SealingKey returns the public key used to seal envelopes, usually an *rsa.PublicKey.
// This method attempts to parse the keys that may have been sent in the exchange,
// either a raw PKIX key or an x509 certificate, in PEM encoding or not. If the sealing
// key cannot be parsed an error is returned.
func (e *Exchange) SealingKey() (_ interface{}, err error) {
	if e.pubkey == nil {
		if e.pubkey, err = ParseKeyExchangeData(e.msg.Data); err != nil {
			return nil, err
		}
	}
	return e.pubkey, nil
}

func (e *Exchange) Proto() (*api.SigningKey, error) {
	return e.msg, nil
}

// UnsealingKey always returns an error for an Exchange key - no private keys are available.
func (e *Exchange) UnsealingKey() (interface{}, error) {
	return nil, ErrNoPrivateKey
}

// PublicKeyAlgorithm refers to the public key algorithm of the x509.Certificate and
// may be used to determine which cipher to use. Typically is "RSA".
func (e *Exchange) PublicKeyAlgorithm() string {
	return e.msg.PublicKeyAlgorithm
}

// PublicKeySignature returns a unique identifier that can be used to manage public keys
// and associate them with their counterpart private keys for unsealing.
func (e *Exchange) PublicKeySignature() (_ string, err error) {
	var pubkey interface{}
	if pubkey, err = e.SealingKey(); err != nil {
		return "", err
	}
	return signature.New(pubkey)
}

// Marshal simply returns the protocol buffer marshaled data for the most compact storage.
func (e *Exchange) Marshal() ([]byte, error) {
	return proto.Marshal(e.msg)
}

// Unmarshal the protocol buffer marshaled data and load the sealing key. If the sealing
// key is invalid or unparsable this method returns an error.
func (e *Exchange) Unmarshal(data []byte) (err error) {
	if err = proto.Unmarshal(data, e.msg); err != nil {
		return err
	}

	if _, err = e.SealingKey(); err != nil {
		return err
	}
	return nil
}
