package keys

import (
	"crypto/x509"
	"fmt"
	"time"

	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/keys/signature"
	"github.com/trisacrypto/trisa/pkg/trust"
)

func FromCertificate(cert *x509.Certificate) (Key, error) {
	return &Certificate{certs: cert}, nil
}

func FromProvider(certs *trust.Provider) (_ Key, err error) {
	key := &Certificate{
		privateKey: certs.GetKey(),
	}

	if key.certs, err = certs.GetLeafCertificate(); err != nil {
		return nil, err
	}
	return key, nil
}

type Certificate struct {
	certs      *x509.Certificate
	privateKey interface{}
}

func (c *Certificate) IsPrivate() bool {
	return c.privateKey != nil
}

func (c *Certificate) SealingKey() (interface{}, error) {
	return c.certs.PublicKey, nil
}

func (c *Certificate) Proto() (msg *api.SigningKey, err error) {
	msg = &api.SigningKey{
		Version:            int64(c.certs.Version),
		Signature:          c.certs.Signature,
		SignatureAlgorithm: c.certs.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: c.PublicKeyAlgorithm(),
		NotBefore:          c.certs.NotBefore.Format(time.RFC3339),
		NotAfter:           c.certs.NotAfter.Format(time.RFC3339),
	}

	if msg.Data, err = x509.MarshalPKIXPublicKey(c.certs.PublicKey); err != nil {
		return nil, fmt.Errorf("could not marshal pkix public key: %s", err)
	}

	return msg, nil
}

func (c *Certificate) UnsealingKey() (interface{}, error) {
	if c.IsPrivate() {
		return c.privateKey, nil
	}
	return nil, ErrNoPrivateKey
}

func (c *Certificate) PublicKeyAlgorithm() string {
	return c.certs.PublicKeyAlgorithm.String()
}

func (c *Certificate) PublicKeySignature() (string, error) {
	// TODO: better signature handler for public keys
	return signature.New(c.certs.PublicKey)
}

func (c *Certificate) Marshal() ([]byte, error) {
	return nil, nil
}

func (c *Certificate) Unmarshal(data []byte) error {
	return nil
}
