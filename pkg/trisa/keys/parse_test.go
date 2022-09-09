package keys_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/trisa/keys"
	"github.com/trisacrypto/trisa/pkg/trust"
)

func TestParsePublicKeys(t *testing.T) {
	var (
		err       error
		pkRSA     *rsa.PublicKey
		pkECDSA   *ecdsa.PublicKey
		pkED25519 ed25519.PublicKey
	)

	pkRSA, err = generateRSAFixture()
	require.NoError(t, err, "could not generate RSA fixture")

	pkECDSA, err = generateECDSAFixture()
	require.NoError(t, err, "could not generate ECDSA fixture")

	pkED25519, err = generateED25519Fixture()
	require.NoError(t, err, "could not generate ED25519 fixture")

	pubkeys := []interface{}{pkRSA, pkECDSA, pkED25519}
	for _, pkey := range pubkeys {
		pkix, err := x509.MarshalPKIXPublicKey(pkey)
		require.NoError(t, err, "could not marshal PKIX public key")

		p2key, err := keys.ParseKeyExchangeData(pkix)
		require.NoError(t, err, "could not parse PKIX public key of type %T", pkey)
		require.Equal(t, pkey, p2key, "parsed public key does not match original")

		pemd, err := pemEncodeFixture(pkey, trust.BlockPublicKey)
		require.NoError(t, err, "could not pem encode public key")
		p3key, err := keys.ParseKeyExchangeData(pemd)
		require.NoError(t, err, "could not parse PEM encoded public key of type %T", pkey)
		require.Equal(t, pkey, p3key, "parsed public key does not match original")
	}

	// Test PCKS1 Encoding
	pkcs1 := x509.MarshalPKCS1PublicKey(pkRSA)

	// Cannot unmarshal PKCS1 encoded keys
	_, err = keys.ParseKeyExchangeData(pkcs1)
	require.ErrorIs(t, err, keys.ErrUnparsableKeyExchange)

	//  But can marshal PEM encoded RSA PKCS1 keys (not other key types)
	pempcks1, err := pemEncodeData(pkcs1, trust.BlockRSAPublicKey)
	require.NoError(t, err, "could not PEM encode PKCS1 data")
	key, err := keys.ParseKeyExchangeData(pempcks1)
	require.NoError(t, err, "could not parse PEM encoded PKCS1 data")
	require.Equal(t, pkRSA, key, "parsing PEM encoded PKCS1 data returned unexpected key")

	// Should be able to parse PKIX data even if the block header is RSA PUBLIC KEY
	pkix, err := x509.MarshalPKIXPublicKey(pkRSA)
	require.NoError(t, err, "could not marshal PKIX RSA keys for PEM block test")

	pempkix, err := pemEncodeData(pkix, trust.BlockRSAPublicKey)
	require.NoError(t, err, "could not PEM encode PKIX data")
	key, err = keys.ParseKeyExchangeData(pempkix)
	require.NoError(t, err, "could not parse PEM encoded PKIX data")
	require.Equal(t, pkRSA, key, "parsing PEM encoded PKIX data returned unexpected key")
}

func TestParseCertificates(t *testing.T) {
	// Loading a trust chain that has a private key attached should return only the
	// public key of the leaf certificate.
	data, err := os.ReadFile("testdata/certs.pem")
	require.NoError(t, err, "could not read certs.pem fixture")

	key, err := keys.ParseKeyExchangeData(data)
	require.NoError(t, err, "could not parse PEM trust chain with private key data")
	require.IsType(t, &rsa.PublicKey{}, key, "expected RSA public key type")

	// Test parsing raw single certificate data
	block, _ := pem.Decode(data)
	key, err = keys.ParseKeyExchangeData(block.Bytes)
	require.NoError(t, err, "could not parse raw certificate data")
	require.IsType(t, &rsa.PublicKey{}, key, "expected RSA public key type")
}

func TestParseErrors(t *testing.T) {
	// Test Bad PEM Encoded Data
	data, err := pemEncodeData([]byte("this is not valid public key material"), trust.BlockCertificate)
	require.NoError(t, err, "could not create bad certificate PEM fixture")
	_, err = keys.ParseKeyExchangeData(data)
	require.ErrorIs(t, err, keys.ErrUnparsableKeyExchange, "no error on bad certificate PEM fixture")

	data, err = pemEncodeData([]byte("this is not valid public key material"), trust.BlockPublicKey)
	require.NoError(t, err, "could not create bad public key PEM fixture")
	_, err = keys.ParseKeyExchangeData(data)
	require.ErrorIs(t, err, keys.ErrUnparsableKeyExchange, "no error on bad public key PEM fixture")

	data, err = pemEncodeData([]byte("this is not valid public key material"), trust.BlockRSAPublicKey)
	require.NoError(t, err, "could not create bad rsa public key PEM fixture")
	_, err = keys.ParseKeyExchangeData(data)
	require.ErrorIs(t, err, keys.ErrUnparsableKeyExchange, "no error on bad rsa public key PEM fixture")

	// Test Mixed Certs and Keys
	var b bytes.Buffer
	f, err := os.Open("testdata/certs.pem")
	require.NoError(t, err, "could not open certs fixture")
	_, err = io.Copy(&b, f)
	require.NoError(t, err, "could not copy certs and private key to buffer")

	key, err := generateRSAFixture()
	require.NoError(t, err, "could not generate rsa fixture")

	data, err = pemEncodeFixture(key, trust.BlockPublicKey)
	require.NoError(t, err, "could not pem encode rsa fixture")
	_, err = b.Write(data)
	require.NoError(t, err, "could not write pem encoded rsa fixture to mixed pem blocks data")

	_, err = keys.ParseKeyExchangeData(b.Bytes())
	require.ErrorIs(t, err, keys.ErrUnparsableKeyExchange, "no error on mixed certs and keys")

	// Test too many keys - reusing data from above
	var b2 bytes.Buffer
	b2.Write(data)
	b2.WriteString("\n")
	b2.Write(data)

	_, err = keys.ParseKeyExchangeData(b2.Bytes())
	require.ErrorIs(t, err, keys.ErrUnparsableKeyExchange, "no error on multiple public keys")
}

func generateRSAFixture() (_ *rsa.PublicKey, err error) {
	var key *rsa.PrivateKey
	if key, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return nil, fmt.Errorf("could not create RSA key fixture: %s", err)
	}
	return &key.PublicKey, nil
}

func generateECDSAFixture() (_ *ecdsa.PublicKey, err error) {
	var key *ecdsa.PrivateKey
	if key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader); err != nil {
		return nil, fmt.Errorf("could not create ECDSA key fixture: %s", err)
	}
	return &key.PublicKey, nil
}

func generateED25519Fixture() (key ed25519.PublicKey, err error) {
	if key, _, err = ed25519.GenerateKey(rand.Reader); err != nil {
		return nil, fmt.Errorf("could not create ED25519 key fixture: %s", err)
	}
	return key, nil
}

func pemEncodeFixture(pub interface{}, block string) (_ []byte, err error) {
	var pkix []byte
	if pkix, err = x509.MarshalPKIXPublicKey(pub); err != nil {
		return nil, err
	}
	return pemEncodeData(pkix, block)
}

func pemEncodeData(data []byte, block string) (_ []byte, err error) {
	var b bytes.Buffer
	if err = pem.Encode(&b, &pem.Block{Type: block, Bytes: data}); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
