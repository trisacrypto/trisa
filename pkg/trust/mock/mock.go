package mock

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"sync"
	"time"

	"github.com/trisacrypto/trisa/pkg/trust"
	"software.sslmate.com/src/go-pkcs12"
)

// TestCA variables
var (
	initCAonce     sync.Once
	rootCA         tls.Certificate
	intermediateCA tls.Certificate
	icaPrivKey     *rsa.PrivateKey
)

// Create a chain with a leaf node, an intermediate, and root ca + private key.
func Chain() (data []byte, err error) {
	initCAonce.Do(initCA)

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(44),
		Subject: pkix.Name{
			CommonName:   "Test",
			Organization: []string{"Test Net"},
			Country:      []string{"XX"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 7),
		SubjectKeyId: []byte{1, 2, 3, 4, 5, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 4096)
	pub := &priv.PublicKey

	// Sign the certificate
	var ca *x509.Certificate
	if ca, err = x509.ParseCertificate(intermediateCA.Certificate[0]); err != nil {
		return nil, err
	}

	var signed []byte
	if signed, err = x509.CreateCertificate(rand.Reader, tmpl, ca, pub, icaPrivKey); err != nil {
		return nil, err
	}

	var cert *x509.Certificate
	if cert, err = x509.ParseCertificate(signed); err != nil {
		return nil, err
	}

	var rca *x509.Certificate
	if rca, err = x509.ParseCertificate(rootCA.Certificate[0]); err != nil {
		return nil, err
	}

	return pkcs12.Encode(rand.Reader, priv, cert, []*x509.Certificate{ca, rca}, pkcs12.DefaultPassword)
}

func initCA() {
	// Root CA
	rootCAtmpl := &x509.Certificate{
		SerialNumber: big.NewInt(42),
		Subject: pkix.Name{
			CommonName:   "Test Root CA",
			Organization: []string{"Test Root"},
			Country:      []string{"XX"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	data, err := x509.CreateCertificate(rand.Reader, rootCAtmpl, rootCAtmpl, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		panic(err)
	}

	if rootCA, err = parseCertificate(data, caPrivKey); err != nil {
		panic(err)
	}

	// Intermediate CA
	interCAtmpl := &x509.Certificate{
		SerialNumber: big.NewInt(43),
		Subject: pkix.Name{
			CommonName:   "Test Intermediate CA",
			Organization: []string{"Test Intermediate"},
			Country:      []string{"XX"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	ca, err := x509.ParseCertificate(rootCA.Certificate[0])
	if err != nil {
		panic(err)
	}

	icaPrivKey, _ = rsa.GenerateKey(rand.Reader, 4096)
	if data, err = x509.CreateCertificate(rand.Reader, interCAtmpl, ca, &icaPrivKey.PublicKey, caPrivKey); err != nil {
		panic(err)
	}

	if intermediateCA, err = parseCertificate(data, icaPrivKey); err != nil {
		panic(err)
	}

}

func parseCertificate(data []byte, priv *rsa.PrivateKey) (tls.Certificate, error) {
	crt, err := x509.ParseCertificate(data)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEMBlock, err := trust.PEMEncodeCertificate(crt)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyPEMBlock, err := trust.PEMEncodePrivateKey(priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(certPEMBlock, keyPEMBlock)
}
