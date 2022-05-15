package keys

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/keys/signature"
	"github.com/trisacrypto/trisa/pkg/trust"
)

// Create a Key from an x509 certificate or sealing certificate issued by TRISA.
func FromCertificate(cert *x509.Certificate) (Key, error) {
	return &Certificate{certs: cert}, nil
}

// Create a Key from an x509 certificate and a private key pair.
func FromX509KeyPair(cert *x509.Certificate, privateKey interface{}) (Key, error) {
	return &Certificate{certs: cert, privateKey: privateKey}, nil
}

// FromProvider returns a Key from a deserialized trust provider. This is the most
// common mechanism of loading a certificate from disk and is the most flexible.
func FromProvider(certs *trust.Provider) (_ Key, err error) {
	key := &Certificate{
		privateKey: certs.GetKey(),
	}

	if key.certs, err = certs.GetLeafCertificate(); err != nil {
		return nil, err
	}
	return key, nil
}

// Certificate is a wrapper for an x509 certificate (containing the public key) and
// optionally for the private key associated with the certificate. This data structure
// is used to create exchange keys using TRISA issued sealing-certificates. The public
// key in the certificate will be used to send keys for sealing and the private key for
// unsealing certificates.
type Certificate struct {
	certs      *x509.Certificate
	privateKey interface{}
}

// Ensure the Certificate implements the Key interface.
var _ Key = &Certificate{}

// IsPrivate returns ture if there is a private key associated with the certificate.
func (c *Certificate) IsPrivate() bool {
	return c.privateKey != nil
}

// SealingKey is public key used to seal envelopes, usually an *rsa.PublicKey.
func (c *Certificate) SealingKey() (interface{}, error) {
	return c.certs.PublicKey, nil
}

// Proto returns the protocol buffer message with the public sealing key to exchange
// with a remote counterparty to begin a TRISA transfer.
func (c *Certificate) Proto() (msg *api.SigningKey, err error) {
	msg = &api.SigningKey{
		Version:            int64(c.certs.Version),
		Signature:          c.certs.Signature,
		SignatureAlgorithm: c.certs.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: c.PublicKeyAlgorithm(),
		NotBefore:          c.certs.NotBefore.Format(time.RFC3339),
		NotAfter:           c.certs.NotAfter.Format(time.RFC3339),
	}

	// TODO: should we marshal into PEM encoded form?
	if msg.Data, err = x509.MarshalPKIXPublicKey(c.certs.PublicKey); err != nil {
		return nil, fmt.Errorf("could not marshal pkix public key: %s", err)
	}

	return msg, nil
}

// UnsealingKey is the private key used to unseal envelopes, usually an *rsa.PrivateKey.
// If no private key is available (IsPrivate() is false) then an error is returned.
func (c *Certificate) UnsealingKey() (interface{}, error) {
	if c.IsPrivate() {
		return c.privateKey, nil
	}
	return nil, ErrNoPrivateKey
}

// PublicKeyAlgorithm refers to the public key algorithm of the x509.Certificate and
// may be used to determine which cipher to use. Typically is "RSA".
func (c *Certificate) PublicKeyAlgorithm() string {
	return c.certs.PublicKeyAlgorithm.String()
}

// PublicKeySignature returns a unique identifier that can be used to manage public keys
// and associate them with their counterpart private keys for unsealing.
func (c *Certificate) PublicKeySignature() (string, error) {
	return signature.New(c.certs.PublicKey)
}

// Marshal returns a PEM encoded Certificate Block along with a PEM encoded private key
// block if one is available. This data is useful for storing Keys on disk or in a
// database, but should not be used to transfer or send keys.
func (c *Certificate) Marshal() (_ []byte, err error) {
	var b bytes.Buffer
	if err = pem.Encode(&b, &pem.Block{Type: trust.BlockCertificate, Bytes: c.certs.Raw}); err != nil {
		return nil, err
	}

	if c.privateKey != nil {
		// Add a new line to separate the certificate and the private key
		if _, err = b.WriteString("\n"); err != nil {
			return nil, err
		}

		var pkcs8 []byte
		if pkcs8, err = x509.MarshalPKCS8PrivateKey(c.privateKey); err != nil {
			return nil, err
		}
		if err = pem.Encode(&b, &pem.Block{Type: trust.BlockPrivateKey, Bytes: pkcs8}); err != nil {
			return nil, err
		}
	}

	return b.Bytes(), nil
}

// Unmarshal a PEM encoded Certificate Block and optionally a PEM encoded private key
// block if one is available. Unmarshal is designed to perform the inverse functionality
// of Marhshal - e.g. Unmarshal will load Marshaled data. However, Unmarshal may work on
// generic PEM encoded chains - e.g. a trust chain with multiple certificates will use
// the first certificate and ignore all others. An error is returned if multiple private
// key blocks are detected in the PEM data.
func (c *Certificate) Unmarshal(data []byte) (err error) {
	// Setup block parsing to unmarshal PEM encoded data.
	chain := make([]*x509.Certificate, 0, 1)
	keys := make([]interface{}, 0)

	var block *pem.Block
	for {
		block, data = pem.Decode(data)
		if block == nil {
			break
		}

		switch block.Type {
		case trust.BlockCertificate:
			var cert *x509.Certificate
			if cert, err = x509.ParseCertificate(block.Bytes); err != nil {
				return fmt.Errorf("could not parse certificate: %s", err)
			}
			chain = append(chain, cert)
		case trust.BlockPrivateKey, trust.BlockECPrivateKey, trust.BlockRSAPrivateKey:
			var key interface{}
			if key, err = trust.ParsePrivateKey(block); err != nil {
				return fmt.Errorf("could not parse block %s: %s", block.Type, err)
			}
			keys = append(keys, key)
		}
	}

	// Validate data parsing
	if len(chain) == 0 {
		return ErrNoCertificate
	}

	if len(keys) > 1 {
		return ErrMultipleKeys
	}

	// We expect the leaf certificate to be the first certificate in the chain
	c.certs = chain[0]
	if len(keys) == 1 {
		c.privateKey = keys[0]
	}
	return nil
}
