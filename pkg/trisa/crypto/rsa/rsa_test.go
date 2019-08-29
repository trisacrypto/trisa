package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	plainText := []byte("encryptthisstuffplease")

	cipherText, err := Encrypt(plainText, &privkey.PublicKey)
	assert.NoError(t, err)

	newPlain, err := Decrypt(cipherText, privkey)
	assert.NoError(t, err)

	assert.Equal(t, plainText, newPlain)
}
