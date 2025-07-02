package rsaoeap_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto/rsaoeap"
)

func TestRSA(t *testing.T) {
	// Generate a random key
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	// Test data
	plaintext := []byte("for your eyes only -- classified") // to encrypt
	data := []byte("trust, but verify.")                    // to sign
	dataWrong := []byte("trust only after verification.")   // to fail verification

	// Cipher only takes RSA keys
	_, err = rsaoeap.New("foo")
	require.ErrorContains(t, err, "could not create RSA cipher from", "unexpected or no error")

	// RSA with public key only
	pubRSA, err := rsaoeap.New(&priv.PublicKey)
	require.NoError(t, err)

	// RSA with private and public keys
	privRSA, err := rsaoeap.New(priv)
	require.NoError(t, err)

	// Encrypt using a new pubRSA with just the public key
	ciphertext, err := pubRSA.Encrypt(plaintext)
	require.NoError(t, err)

	// Decrypt using a new cipher with both public and private key
	decoded, err := privRSA.Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decoded)

	// Sign using the private key
	signature, err := privRSA.Sign(data)
	require.NoError(t, err, "could not sign message")
	require.NotNil(t, signature, "expected a signature")

	// Verify using just the public key
	err = pubRSA.Verify(data, signature)
	require.NoError(t, err, "could not verify message")

	// Error when verifying an unmatched data/signature pair
	err = pubRSA.Verify(dataWrong, signature)
	require.ErrorIs(t, err, rsa.ErrVerification, "message was mistakenly verified")
}
