/*
Package crypto describes interfaces for the various encryption and hmac algorithms that
might be used to encrypt and sign transaction envelopes being passed securely in the
TRISA network. Subpackages implement specific algorithms such as aesgcm for symmetric
encryption or rsa for asymmetric encryption. Note that not all encryption mechanisms are
legal in different countries, these interfaces allow the use of different algorithms and
methodologies in the protocol without specifying what must be used.
*/
package crypto

import "crypto/rand"

// Crypto handler for TRISA transaction envelopes must be both a Cipher and a Signer.
type Crypto interface {
	Cipher
	Signer
}

// Cipher is a device that can perform encryption and decryption, This interface wraps
// different encryption algorithms that must be identified in the TRISA protocol.
type Cipher interface {
	Encrypt(plaintext []byte) (ciphertext []byte, err error)
	Decrypt(ciphertext []byte) (plaintext []byte, err error)
	EncryptionAlgorithm() string
}

// Signer creates or verifies HMAC signatures. This interface wraps multiple hmac
// algorithms that must be identified in the TRISA protocol.
type Signer interface {
	Sign(data []byte) (signature []byte, err error)
	Verify(data, signature []byte) (err error)
	SignatureAlgorithm() string
}

// Random generates a secure random sequence of bytes, this helper function is used to
// easily create keys, salts, and secrets in the crypto subpackages.
func Random(n int) (b []byte, err error) {
	b = make([]byte, n)
	if _, err = rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
