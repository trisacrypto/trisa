package aesgcm

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/trisacrypto/trisa/pkg/trisa/crypto"
)

const (
	Algorithm          = "AES-GCM"
	SignatureAlgorithm = "HMAC-SHA256"
	Algorithm256       = "AES256-GCM"
	Algorithm192       = "AES192-GCM"
	Algorithm128       = "AES128-GCM"
)

// AESGCM implements the crypto.Crypto interface using the AES-GCM algorithm for
// symmetric-key encryption. This algorithm is widely adopted for it's performance and
// throughput rates for state-of-the-art high-speed communication on inexpensive
// hardware (Wikipedia). This implementation generates a 32 byte random encryption key
// when initialized, if one not specified by default. Users should create a new AESGCM
// to encrypt and sign different messages with different keys.
type AESGCM struct {
	key    []byte // the symmetric encryption key
	secret []byte // the HMAC secret used to calculate the signature
}

// New creates an AESGCM Crypto handler, generating an encryption key if it is nil or
// zero length. If hmac_secret isn't specified, the encryption key is used. The key and
// secret should be 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
func New(encryptionKey, hmacSecret []byte) (_ *AESGCM, err error) {
	if len(encryptionKey) == 0 {
		if encryptionKey, err = crypto.Random(32); err != nil {
			return nil, fmt.Errorf("could not generate encryption key: %w", err)
		}
	}

	if len(hmacSecret) == 0 {
		hmacSecret = encryptionKey
	}

	return &AESGCM{key: encryptionKey, secret: hmacSecret}, nil
}

// Encrypt a message using the struct key, appending a 12 byte random nonce to the end
// of the ciphertext message.
func (c *AESGCM) Encrypt(plaintext []byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := crypto.Random(12)
	if err != nil {
		return nil, err
	}

	ciphertext = aesgcm.Seal(nil, nonce, plaintext, nil)
	ciphertext = append(ciphertext, nonce...)
	return ciphertext, nil
}

// Decrypt a message using the struct key, extracting the nonce from the end.
func (c *AESGCM) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	if len(ciphertext) < 12 {
		return nil, crypto.ErrMissingCiphertext
	}

	data := ciphertext[:len(ciphertext)-12]
	nonce := ciphertext[len(ciphertext)-12:]

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err = aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt ciphertext: %w", err)
	}
	return plaintext, nil
}

// EncryptionAlgorithm returns the name of the algorithm for adding to the Transaction.
func (c *AESGCM) EncryptionAlgorithm() string {
	switch len(c.key) {
	case 32:
		return Algorithm256
	case 24:
		return Algorithm192
	case 16:
		return Algorithm128
	default:
		return Algorithm
	}
}

// Sign the specified data (ususally the ciphertext) using the struct secret.
func (c *AESGCM) Sign(data []byte) (signature []byte, err error) {
	if len(data) == 0 {
		return nil, crypto.ErrCannotSignEmpty
	}

	hm := hmac.New(sha256.New, c.secret)
	hm.Write(data)
	return hm.Sum(nil), nil
}

// Verify the signature on the specified data using the struct secret.
func (c *AESGCM) Verify(data, signature []byte) (err error) {
	hm := hmac.New(sha256.New, c.secret)
	hm.Write(data)

	if !bytes.Equal(signature, hm.Sum(nil)) {
		return crypto.ErrHMACSignatureMismatch
	}

	return nil
}

// SignatureAlgorithm returns the name of the hmac_algorithm for adding to the Transaction.
func (c *AESGCM) SignatureAlgorithm() string {
	return SignatureAlgorithm
}

// EncryptionKey is a read-only getter.
func (c *AESGCM) EncryptionKey() []byte {
	return c.key
}

// HMACSecret is a read-only getter.
func (c *AESGCM) HMACSecret() []byte {
	return c.secret
}
