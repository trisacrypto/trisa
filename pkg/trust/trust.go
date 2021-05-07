/*
Package trust provides support for handling PEM-encoded certificate chains. Its primary
functionality is a trust Provider that manages a single chain (either public or public
and private). A ProviderPool manages only public chains and is used to facilitate mTLS
exchanges. This package also includes certificate serialization management and helpers,
including compression utilities.
*/
package trust

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"software.sslmate.com/src/go-pkcs12"
)

// PEM Block types
const (
	BlockPublickKey         = "PUBLIC KEY"
	BlockPrivateKey         = "PRIVATE KEY"
	BlockRSAPublicKey       = "RSA PUBLIC KEY"
	BlockRSAPrivateKey      = "RSA PRIVATE KEY"
	BlockECPrivateKey       = "EC PRIVATE KEY"
	BlockCertificate        = "CERTIFICATE"
	BlockCertificateRequest = "CERTIFICATE REQUEST"
)

// Provider wraps a PEM-encoded certificate chain, which can optionally include private
// keys. Providers with keys (private providers) are used to instantiate mTLS servers,
// while public Providers are used in ProviderPools to facilitate mTLS clients.
type Provider struct {
	chain tls.Certificate
	key   interface{}
}

// New creates a Provider from PEM encoded blocks.
func New(chain []byte) (p *Provider, err error) {
	p = &Provider{}
	if err = p.Decode(chain); err != nil {
		return nil, err
	}
	return p, nil
}

// Decrypt pfxData from a PKCS12 encoded file using the specified password. The data
// must contain at least one certificate and only one private key. The first certificate
// in the data is assumed to be the leaf certificate and subsequent certificates are the
// CA chain. Certificates are appended to the provider leaf first, then CA chain.
func Decrypt(pfxData []byte, password string) (p *Provider, err error) {
	p = &Provider{}

	var crt *x509.Certificate
	var ca []*x509.Certificate
	if p.key, crt, ca, err = pkcs12.DecodeChain(pfxData, password); err != nil {
		return nil, err
	}

	p.chain.Certificate = append(p.chain.Certificate, crt.Raw)
	for _, c := range ca {
		p.chain.Certificate = append(p.chain.Certificate, c.Raw)
	}
	return p, nil
}

// Encrypt Provider as pfxData containing a private key, the end-entity certificate as
// the leaf certificate, and the rest of the chain as CA certificates. The private key
// is encrypted with the provided password, however, because of the weak encryption
// primitives used by PKCS#12, it is strongly recommended that a hardcoded password
// (such as pkcs12.DefaultPassword) is specified and the pfxData is protected using
// other means. See the software.sslmate.com go-pkcs12 package for more details.
//
// This method will return an error if the private key is nil, if there are no
// certificates stored on the provider, or if the certs cannot be parsed as x509 certs.
func (p *Provider) Encrypt(password string) (pfxData []byte, err error) {
	if len(p.chain.Certificate) == 0 {
		return nil, ErrNoCertificates
	}

	if p.key == nil {
		return nil, ErrKeyRequired
	}

	// Assume certificate is the first element in the chain
	var crt *x509.Certificate
	if crt, err = x509.ParseCertificate(p.chain.Certificate[0]); err != nil {
		return nil, fmt.Errorf("could not parse leaf certificate: %s", err)
	}

	var ca []*x509.Certificate
	if len(p.chain.Certificate) > 1 {
		ca = make([]*x509.Certificate, len(p.chain.Certificate)-1)
		for i, cacrt := range p.chain.Certificate[1:] {
			if ca[i], err = x509.ParseCertificate(cacrt); err != nil {
				return nil, fmt.Errorf("could not parse certificate %d: %s", i+1, err)
			}
		}
	}

	return pkcs12.Encode(rand.Reader, p.key, crt, ca, password)
}

// Decode PEM blocks and adds them to the provider. Certificates are appended to the
// Provider chain and Private Keys are Unmarshalled from PKCS8. All other block types
// return an error and stop processing the block or chain. Only the private key is
// verified for correctness, certificates are unvalidated.
func (p *Provider) Decode(in []byte) (err error) {
	var block *pem.Block
	for {
		block, in = pem.Decode(in)
		if block == nil {
			break
		}

		switch block.Type {
		case BlockCertificate:
			p.chain.Certificate = append(p.chain.Certificate, block.Bytes)
		case BlockPrivateKey, BlockECPrivateKey, BlockRSAPrivateKey:
			if p.key, err = parsePrivateKey(block); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unhandled block type %q", block.Type)
		}
	}
	return nil
}

// Encode Provider in PCKS12 PEM format for serialization. Certificates are written to
// the array first. If the private key exists, it is written as the last PEM block.
func (p *Provider) Encode() (_ []byte, err error) {
	var b bytes.Buffer
	var block []byte

	for i, asn1Data := range p.chain.Certificate {
		var crt *x509.Certificate
		if crt, err = x509.ParseCertificate(asn1Data); err != nil {
			return nil, fmt.Errorf("could not parse certificate %d: %s", i, err)
		}

		if block, err = PEMEncodeCertificate(crt); err != nil {
			return nil, fmt.Errorf("could not encode certificate %d: %s", i, err)
		}

		b.Write(block)
	}

	if p.key != nil {
		if block, err = PEMEncodePrivateKey(p.key); err != nil {
			return nil, fmt.Errorf("could not encode private key: %s", err)
		}

		b.Write(block)
	}

	return b.Bytes(), nil
}

// GetCertPool returns the x509.CertPool certificate set representing the root,
// intermediate, and leaf certificates of the Provider. This pool is provider-specific
// and does not include system certificates.
func (p *Provider) GetCertPool() (_ *x509.CertPool, err error) {
	pool := x509.NewCertPool()
	for _, c := range p.chain.Certificate {
		var x509Cert *x509.Certificate
		if x509Cert, err = x509.ParseCertificate(c); err != nil {
			return nil, err
		}
		pool.AddCert(x509Cert)
	}
	return pool, nil
}

// GetKeyPair returns a tls.Certificate parsed from the PEM encoded data maintained by
// the provider. This method uses tls.X509KeyPair to ensure that the public/private key
// pair are suitable for use with an HTTP Server.
func (p *Provider) GetKeyPair() (_ tls.Certificate, err error) {
	if p.key == nil {
		return tls.Certificate{}, ErrKeyRequired
	}

	var block []byte
	var certs bytes.Buffer
	for i, asn1Data := range p.chain.Certificate {
		var crt *x509.Certificate
		if crt, err = x509.ParseCertificate(asn1Data); err != nil {
			return tls.Certificate{}, fmt.Errorf("could not parse certificate %d: %s", i, err)
		}

		if block, err = PEMEncodeCertificate(crt); err != nil {
			return tls.Certificate{}, fmt.Errorf("could not encode certificate %d: %s", i, err)
		}

		certs.Write(block)
	}

	var key []byte
	if key, err = PEMEncodePrivateKey(p.key); err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(certs.Bytes(), key)
}

// GetLeafCertificate returns the parsed x509 leaf certificate if it exists, returning
// an error if there are no certificates or if there is a parse error.
func (p *Provider) GetLeafCertificate() (*x509.Certificate, error) {
	if p.chain.Leaf != nil {
		return p.chain.Leaf, nil
	}

	if len(p.chain.Certificate) == 0 {
		return nil, ErrNoCertificates
	}
	return x509.ParseCertificate(p.chain.Certificate[0])
}

// GetKey returns the private key, or nil if this is a public provider.
func (p *Provider) GetKey() interface{} {
	return p.key
}

// GetRSAKeys returns a fully constructed RSA PrivateKey that includes the public key
// material property. This method errors if the key is not an RSA key or does not exist.
func (p *Provider) GetRSAKeys() (key *rsa.PrivateKey, err error) {
	if p.key == nil {
		return nil, ErrKeyRequired
	}

	var ok bool
	if key, ok = p.key.(*rsa.PrivateKey); !ok {
		return nil, fmt.Errorf("private key is not RSA but is %T", p.key)
	}

	var cert *x509.Certificate
	if cert, err = p.GetLeafCertificate(); err != nil {
		return nil, err
	}

	key.PublicKey = *cert.PublicKey.(*rsa.PublicKey)
	return key, nil
}

// IsPrivate returns true if the Provider contains a non-nil key.
func (p *Provider) IsPrivate() bool {
	return p.key != nil
}

// Public returns a Provider without the key. If the Provider is already public, then
// the pointer to the same Provider is returned (does not clone).
func (p *Provider) Public() *Provider {
	if p.key == nil {
		return p
	}
	return &Provider{chain: p.chain, key: nil}
}

// String returns the common name of the Provider from the leaf certificate.
func (p *Provider) String() string {
	cert, err := p.GetLeafCertificate()
	if err != nil {
		return ""
	}
	return cert.Subject.CommonName
}
