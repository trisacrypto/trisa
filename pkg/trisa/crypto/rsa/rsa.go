package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

func Encrypt(data []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	cipherText, err := rsa.EncryptOAEP(hash, rand.Reader, pub, data, nil)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

func Decrypt(cipherText []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, priv, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
