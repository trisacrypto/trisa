package trust_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/trust"
)

func TestPEMPrivateKey(t *testing.T) {
	// Handling RSA keys (primary usage)
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	data, err := trust.PEMEncodePrivateKey(key)
	require.NoError(t, err)

	keyb, err := trust.PEMDecodePrivateKey(data)
	require.NoError(t, err)
	require.Equal(t, key, keyb)

	// Handling RSA PRIVATE KEY block type
	var b bytes.Buffer
	pkcs1 := x509.MarshalPKCS1PrivateKey(key)
	err = pem.Encode(&b, &pem.Block{Type: trust.BlockRSAPrivateKey, Bytes: pkcs1})
	require.NoError(t, err)

	keyc, err := trust.PEMDecodePrivateKey(b.Bytes())
	require.NoError(t, err)
	require.Equal(t, key, keyc)

	// Handling EC PRIVATE KEY block type
	eckey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	ec, err := x509.MarshalECPrivateKey(eckey)
	require.NoError(t, err)

	var c bytes.Buffer
	err = pem.Encode(&c, &pem.Block{Type: trust.BlockECPrivateKey, Bytes: ec})
	require.NoError(t, err)

	keyd, err := trust.PEMDecodePrivateKey(c.Bytes())
	require.NoError(t, err)
	require.Equal(t, eckey, keyd)
}

func TestPEMCertificate(t *testing.T) {
	crt, err := cert()
	require.NoError(t, err)

	data, err := trust.PEMEncodeCertificate(crt)
	require.NoError(t, err)

	crtb, err := trust.PEMDecodeCertificate(data)
	require.NoError(t, err)

	require.Equal(t, crt, crtb)
}

func cert() (*x509.Certificate, error) {
	tpl := &x509.Certificate{
		SerialNumber: big.NewInt(42),
		Subject: pkix.Name{
			CommonName:   "TestNet",
			Organization: []string{"Test"},
			Country:      []string{"XX"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 5, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	pub := &key.PublicKey
	signed, err := x509.CreateCertificate(rand.Reader, tpl, tpl, pub, key)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(signed)
}

func TestPEMCertificateSigningRequest(t *testing.T) {
	req, err := csr()
	require.NoError(t, err)

	data, err := trust.PEMEncodeCSR(req)
	require.NoError(t, err)

	reqb, err := trust.PEMDecodeCSR(data)
	require.NoError(t, err)

	require.Equal(t, req, reqb)
}

func csr() (*x509.CertificateRequest, error) {
	tpl := &x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{"Test"},
			Country:      []string{"XX"},
		},
		SignatureAlgorithm: x509.SHA512WithRSA,
	}

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	der, err := x509.CreateCertificateRequest(rand.Reader, tpl, key)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificateRequest(der)
}
