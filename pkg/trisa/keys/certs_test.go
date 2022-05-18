package keys_test

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/keys"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestFromCertificate(t *testing.T) {
	// Test creating Certificate based keys from a fixture
	sz, err := trust.NewSerializer(false)
	require.NoError(t, err, "could not create serializer to load fixture")

	provider, err := sz.ReadFile("testdata/certs.pem")
	require.NoError(t, err, "could not read test fixture")

	key, err := keys.FromProvider(provider)
	require.NoError(t, err, "could not create keys from provider")
	require.True(t, key.IsPrivate(), "keys expected to contain private unsealing key")

	certs, err := provider.GetLeafCertificate()
	require.NoError(t, err, "could not get certificate from provider")
	require.Contains(t, certs.DNSNames, "keys.trisa.dev", "wrong keys extracted from provider")

	key, err = keys.FromCertificate(certs)
	require.NoError(t, err, "could not create keys from certificate")
	require.False(t, key.IsPrivate(), "keys expected to not contain private unsealing key")
	_, err = key.UnsealingKey()
	require.ErrorIs(t, err, keys.ErrNoPrivateKey, "if certs are not private, should not return private key")

	key, err = keys.FromX509KeyPair(certs, provider.GetKey())
	require.NoError(t, err, "could not create keys from x509 key pair")
	require.True(t, key.IsPrivate(), "keys expected to contain private unsealing key")

	key, err = keys.FromX509KeyPair(certs, nil)
	require.NoError(t, err, "could not create keys from x509 key pair wiht nil key")
	require.False(t, key.IsPrivate(), "keys expected to not contain private unsealing key")
	_, err = key.UnsealingKey()
	require.ErrorIs(t, err, keys.ErrNoPrivateKey, "if certs are not private, should not return private key")
}

func TestCertificateWithPrivateKey(t *testing.T) {
	cert, err := loadCertificateFixture()
	require.NoError(t, err, "could not load certificate fixture with private keys")
	require.True(t, cert.IsPrivate(), "expected cert to have private keys for test")

	sealingKey, err := cert.SealingKey()
	require.NoError(t, err, "could not extract public sealing key from certs")
	require.IsType(t, &rsa.PublicKey{}, sealingKey, "expected an public key for the sealing key")

	unsealingKey, err := cert.UnsealingKey()
	require.NoError(t, err, "could not extract private unsealing key from certs")
	require.IsType(t, &rsa.PrivateKey{}, unsealingKey, "expected a private key for the sealing key")

	require.Equal(t, "RSA", cert.PublicKeyAlgorithm(), "unexpected public key algorithm")
	pks, err := cert.PublicKeySignature()
	require.NoError(t, err, "could not compute public key signature from certs")
	require.Equal(t, "SHA256:ToLuV/KNXWxQWd7/b3DK/hCgbN1Zg1PQECIHWr3F8KM", pks, "unexpected public key signature")

	// Test Proto method of certificate produces a valid api.SigningKey
	// NOTE correct parsing of data is in another test
	msg, err := cert.Proto()
	require.NoError(t, err, "could not create an exchange key with certs")
	require.Equal(t, int64(3), msg.Version, "incorrect signing key version")
	require.NotEmpty(t, msg.Signature, "incorrect signing key signature")
	require.Equal(t, "SHA256-RSA", msg.SignatureAlgorithm, "incorrect signing key signature algorithm")
	require.Equal(t, "RSA", msg.PublicKeyAlgorithm, "incorrect signing key public key algorithm")
	require.Equal(t, "2022-05-15T18:34:12Z", msg.NotBefore, "incorrect signing key not before")
	require.Equal(t, "2022-11-15T19:34:12Z", msg.NotAfter, "incorrect signing key not after")
	require.NotEmpty(t, msg.Data, "incorrect signing key data")

	// Test marshaling the certificate with the private keys
	data, err := cert.Marshal()
	require.NoError(t, err, "could not marshal certificate")
	require.NotEmpty(t, data, "no marshaled data returned")

	// Test unmarshaling the certificate from the marshaled data
	other := &keys.Certificate{}
	require.False(t, other.IsPrivate(), "empty certificate should not private")

	err = other.Unmarshal(data)
	require.NoError(t, err, "could not marshal certificate")
	require.True(t, other.IsPrivate(), "unmarshaled certificate should be private")

	opks, err := other.PublicKeySignature()
	require.NoError(t, err, "could not compute public key signature from unmarshaled certificate")
	require.Equal(t, pks, opks, "unmarshaled certificate does not match original certificate")
}

func TestCertificateNoPrivateKey(t *testing.T) {
	provider, err := loadCertificateProvider("testdata/certs.pem")
	require.NoError(t, err, "could not load certificate fixture")

	leaf, err := provider.GetLeafCertificate()
	require.NoError(t, err, "could not get leaf x509 certificate from provider")
	require.Contains(t, leaf.DNSNames, "keys.trisa.dev", "wrong keys extracted from provider")

	cert, err := keys.FromCertificate(leaf)
	require.NoError(t, err, "could not create certificate fixture without private keys")
	require.False(t, cert.IsPrivate(), "expected cert not to have private keys for test")

	sealingKey, err := cert.SealingKey()
	require.NoError(t, err, "could not extract public sealing key from certs")
	require.IsType(t, &rsa.PublicKey{}, sealingKey, "expected an public key for the sealing key")

	unsealingKey, err := cert.UnsealingKey()
	require.ErrorIs(t, err, keys.ErrNoPrivateKey)
	require.Nil(t, unsealingKey)

	require.Equal(t, "RSA", cert.PublicKeyAlgorithm(), "unexpected public key algorithm")
	pks, err := cert.PublicKeySignature()
	require.NoError(t, err, "could not compute public key signature from certs")
	require.Equal(t, "SHA256:ToLuV/KNXWxQWd7/b3DK/hCgbN1Zg1PQECIHWr3F8KM", pks, "unexpected public key signature")

	// Test Proto method of certificate produces a valid api.SigningKey
	// NOTE correct parsing of data is in another test
	msg, err := cert.Proto()
	require.NoError(t, err, "could not create an exchange key with certs")
	require.Equal(t, int64(3), msg.Version, "incorrect signing key version")
	require.NotEmpty(t, msg.Signature, "incorrect signing key signature")
	require.Equal(t, "SHA256-RSA", msg.SignatureAlgorithm, "incorrect signing key signature algorithm")
	require.Equal(t, "RSA", msg.PublicKeyAlgorithm, "incorrect signing key public key algorithm")
	require.Equal(t, "2022-05-15T18:34:12Z", msg.NotBefore, "incorrect signing key not before")
	require.Equal(t, "2022-11-15T19:34:12Z", msg.NotAfter, "incorrect signing key not after")
	require.NotEmpty(t, msg.Data, "incorrect signing key data")

	// Test marshaling the certificate with the private keys
	data, err := cert.Marshal()
	require.NoError(t, err, "could not marshal certificate")
	require.NotEmpty(t, data, "no marshaled data returned")

	// Test unmarshaling the certificate from the marshaled data
	other := &keys.Certificate{}
	require.False(t, other.IsPrivate(), "empty certificate should not private")

	err = other.Unmarshal(data)
	require.NoError(t, err, "could not marshal certificate")
	require.False(t, other.IsPrivate(), "unmarshaled certificate should continue not to have private keys")

	opks, err := other.PublicKeySignature()
	require.NoError(t, err, "could not compute public key signature from unmarshaled certificate")
	require.Equal(t, pks, opks, "unmarshaled certificate does not match original certificate")

}

func TestGDSLookup(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/lookup.json")
	require.NoError(t, err, "could not load lookup reply fixture")

	rep := &gds.LookupReply{}
	err = protojson.Unmarshal(data, rep)
	require.NoError(t, err, "could not unmarshal lookup reply fixture")

	key, err := keys.FromGDSLookup(rep)
	require.NoError(t, err, "could not parse keys from gds lookup reply")
	require.False(t, key.IsPrivate(), "no private key should be returned from gds")

	pks, err := key.PublicKeySignature()
	require.NoError(t, err, "could not create public key signature")
	require.Equal(t, "SHA256:p6e+TfID7MG1l6V0QsfJpuv49t0q5sHfFY1WUNkkSgk", pks, "unexpected public key signature for alice")
}

func loadCertificateProvider(path string) (_ *trust.Provider, err error) {
	var sz *trust.Serializer
	if sz, err = trust.NewSerializer(false); err != nil {
		return nil, err
	}
	return sz.ReadFile(path)
}

func loadCertificateFixture() (cert *keys.Certificate, err error) {
	var provider *trust.Provider
	if provider, err = loadCertificateProvider("testdata/certs.pem"); err != nil {
		return nil, err
	}

	var key keys.Key
	if key, err = keys.FromProvider(provider); err != nil {
		return nil, err
	}

	var ok bool
	if cert, ok = key.(*keys.Certificate); !ok {
		return nil, fmt.Errorf("%T is not a *keys.Certificate", cert)
	}
	return cert, nil
}
