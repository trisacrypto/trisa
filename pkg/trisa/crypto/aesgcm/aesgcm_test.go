package aesgcm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	plainText := []byte("thisbetterworks")
	cipherText, key, sig, _, err := Encrypt(plainText)

	assert.NoError(t, err)

	newPlain, err := Decrypt(cipherText, sig, key)
	assert.NoError(t, err)

	assert.Equal(t, plainText, newPlain)
}
