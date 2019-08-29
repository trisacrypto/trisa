package aesgcm

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/trisacrypto/trisa/pkg/trisa/crypto"
)

// return vals
// * cipher text
// * cipher secret (encrypted secret using pub key of receiver)
// * hmac signature
// * hmac secret (encrypted secret using pub key of receiver)
func Encrypt(plainText []byte) ([]byte, []byte, []byte, []byte, error) {

	key, err := crypto.GenRandom(32)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	nonce, err := crypto.GenRandom(12)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	cipherText := aesgcm.Seal(nil, nonce, plainText, nil)

	sig, err := createHMAC(key, cipherText)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	cipherText = append(cipherText, nonce...)

	return cipherText, key, sig, key, nil
}

// Decrypt validates mac and returns decoded data.
func Decrypt(cipherText, sig, key []byte) ([]byte, error) {

	if len(cipherText) == 0 {
		return nil, errors.New("empty cipher text")
	}

	data := cipherText[:len(cipherText)-12]
	nonce := cipherText[len(cipherText)-12:]

	if err := validateHMAC(key, data, sig); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plainText, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func createHMAC(key, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	hm := hmac.New(sha256.New, key)
	hm.Write(data)
	return hm.Sum(nil), nil
}

func validateHMAC(key, data, sig []byte) error {
	hm := hmac.New(sha256.New, key)
	hm.Write(data)

	if !bytes.Equal(sig, hm.Sum(nil)) {
		return fmt.Errorf("hmac mismatch")
	}

	return nil
}
