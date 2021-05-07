package aesgcm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto/aesgcm"
)

func TestAESGCM(t *testing.T) {
	plaintext := []byte("theeaglefliesatmidnight")

	// Generate a key
	cipher, err := aesgcm.New(nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, cipher.EncryptionKey())
	require.NotEmpty(t, cipher.HMACSecret())

	ciphertext, err := cipher.Encrypt(plaintext)
	require.NoError(t, err)

	signature, err := cipher.Sign(ciphertext)
	require.NoError(t, err)

	// Decode using a new cipher
	var decoder crypto.Crypto
	decoder, err = aesgcm.New(cipher.EncryptionKey(), cipher.HMACSecret())
	require.NoError(t, err)

	err = decoder.Verify(ciphertext, signature)
	require.NoError(t, err)

	decoded, err := decoder.Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decoded)
}
