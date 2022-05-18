package keys_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/keys"
	"github.com/trisacrypto/trisa/pkg/trust"
)

func TestExchangePKIX(t *testing.T) {
	// Test an Exchange Key loaded from a PKIX marshaled key
	// This kind of key is the default currently in-use in TRISA but will be deprecated
	// eventually for PEM-encoded keys. This test ensures backwards compatability.
	msg := &api.SigningKey{
		Version:            42,
		Signature:          nil,
		SignatureAlgorithm: "NONE",
		PublicKeyAlgorithm: x509.RSA.String(),
		NotBefore:          "2022-05-18T08:00:00Z",
		NotAfter:           "2032-05-19T07:59:59Z",
		Revoked:            false,
		Data:               nil,
	}

	rsaKeys, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "could not generate RSA keys")

	msg.Data, err = x509.MarshalPKIXPublicKey(&rsaKeys.PublicKey)
	require.NoError(t, err, "could not marshal PKIX RSA public key fixture")

	key, err := keys.FromSigningKey(msg)
	require.NoError(t, err, "could not create exchange key from protocol buffer")
	require.False(t, key.IsPrivate(), "key should be public only")

	sealingKey, err := key.SealingKey()
	require.NoError(t, err, "could not fetch sealing key")
	require.Equal(t, &rsaKeys.PublicKey, sealingKey, "sealing key did not match generated RSA public key")

	pb, err := key.Proto()
	require.NoError(t, err, "could not create SigningKey protocol buffer")
	require.Equal(t, msg, pb, "api SigningKey messages do not match")

	_, err = key.UnsealingKey()
	require.ErrorIs(t, err, keys.ErrNoPrivateKey, "unexpectedly did not return on error on unsealing key")

	require.Equal(t, "RSA", key.PublicKeyAlgorithm())

	pks, err := key.PublicKeySignature()
	require.NoError(t, err, "could not create public key signature")
	require.NotEmpty(t, pks, "no public key signature returned")
	require.Len(t, pks, 50, "public key signature returned unexpected length")

	data, err := key.Marshal()
	require.NoError(t, err, "could not marshal exchange key")
	require.NotEmpty(t, data, "marshal returned empty data")

	key2 := &keys.Exchange{}
	err = key2.Unmarshal(data)
	require.NoError(t, err, "could not unmarshal exchange key")
	pks2, err := key2.PublicKeySignature()
	require.NoError(t, err, "could not create public key signature from unmarshaled exchange key")
	require.Equal(t, pks, pks2, "public key signatures after serialization do not match")
}

func TestExchangePEM(t *testing.T) {
	// Test an Exchange Key loaded from a PEM encoded key
	// This kind of key is the future recommended method for key serialization and the
	// Exchange Key struct should support handling it.
	msg := &api.SigningKey{
		Version:            42,
		Signature:          nil,
		SignatureAlgorithm: "NONE",
		PublicKeyAlgorithm: x509.RSA.String(),
		NotBefore:          "2022-05-18T08:00:00Z",
		NotAfter:           "2032-05-19T07:59:59Z",
		Revoked:            false,
		Data:               nil,
	}

	rsaKeys, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "could not generate RSA keys")

	pkdata, err := x509.MarshalPKIXPublicKey(&rsaKeys.PublicKey)
	require.NoError(t, err, "could not marshal PKIX RSA public key fixture")
	msg.Data = pem.EncodeToMemory(&pem.Block{Type: trust.BlockPublicKey, Bytes: pkdata})

	key, err := keys.FromSigningKey(msg)
	require.NoError(t, err, "could not create exchange key from protocol buffer")
	require.False(t, key.IsPrivate(), "key should be public only")

	sealingKey, err := key.SealingKey()
	require.NoError(t, err, "could not fetch sealing key")
	require.Equal(t, &rsaKeys.PublicKey, sealingKey, "sealing key did not match generated RSA public key")

	pb, err := key.Proto()
	require.NoError(t, err, "could not create SigningKey protocol buffer")
	require.Equal(t, msg, pb, "api SigningKey messages do not match")

	_, err = key.UnsealingKey()
	require.ErrorIs(t, err, keys.ErrNoPrivateKey, "unexpectedly did not return on error on unsealing key")

	require.Equal(t, "RSA", key.PublicKeyAlgorithm())

	pks, err := key.PublicKeySignature()
	require.NoError(t, err, "could not create public key signature")
	require.NotEmpty(t, pks, "no public key signature returned")
	require.Len(t, pks, 50, "public key signature returned unexpected length")

	data, err := key.Marshal()
	require.NoError(t, err, "could not marshal exchange key")
	require.NotEmpty(t, data, "marshal returned empty data")

	key2 := &keys.Exchange{}
	err = key2.Unmarshal(data)
	require.NoError(t, err, "could not unmarshal exchange key")
	pks2, err := key2.PublicKeySignature()
	require.NoError(t, err, "could not create public key signature from unmarshaled exchange key")
	require.Equal(t, pks, pks2, "public key signatures after serialization do not match")
}

func TestExchangeBadKey(t *testing.T) {
	// Should not be able to create an Exchange key with bad Data
	msg := &api.SigningKey{
		Version:            42,
		Signature:          nil,
		SignatureAlgorithm: "NONE",
		PublicKeyAlgorithm: x509.RSA.String(),
		NotBefore:          "2022-05-18T08:00:00Z",
		NotAfter:           "2032-05-19T07:59:59Z",
		Revoked:            false,
		Data:               nil,
	}

	key, err := keys.FromSigningKey(msg)
	require.ErrorIs(t, err, keys.ErrNoKeyData, "expected no key data when msg.Data is nil")
	require.Nil(t, key, "no key object should be returned on error")

	msg.Data = []byte("this is not a valid key data")
	key, err = keys.FromSigningKey(msg)
	require.ErrorIs(t, err, keys.ErrUnparsableKeyExchange, "expected unparsable error when msg.Data is invalid")
	require.Nil(t, key, "no key object should be returned on error")
}
