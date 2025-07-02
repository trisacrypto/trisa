package rsaoeap

import (
	gocrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"

	"github.com/trisacrypto/trisa/pkg/trisa/crypto"
	"github.com/trisacrypto/trisa/pkg/trisa/keys/signature"
)

const (
	CipherAlgorithm = "RSA-OAEP-SHA512"
	SignerAlgorithm = "RSA-PSS-SHA512"
)

// RSA implements the crypto.Cipher interface using RSA public/private key algorithm
// as specified in PKCS #1. Messages are encrypted with the public key and can only be
// decrypted using the private key. RSA objects must have a public key but the private
// key is only required for decryption.
type RSA struct {
	pub  *rsa.PublicKey
	priv *rsa.PrivateKey
}

// New creates an RSA Crypto handler with the specified key pair. If the handler
// is only being used for encryption or signature verification, simply pass the
// public key: New(pub *rsa.PublicKey); If the cipher is being used for
// decryption or signing, then pass the private key: New(key *rsa.PrivateKey).
func New(key interface{}) (_ *RSA, err error) {
	switch t := key.(type) {
	case *rsa.PublicKey:
		return &RSA{pub: t, priv: nil}, nil
	case *rsa.PrivateKey:
		return &RSA{pub: &t.PublicKey, priv: t}, nil
	default:
		return nil, fmt.Errorf("could not create RSA cipher from %T", t)
	}
}

// ############################################################################
// Cipher interface
// ############################################################################

// Encrypt the message using the public key.
func (c *RSA) Encrypt(plaintext []byte) (ciphertext []byte, err error) {
	hash := sha512.New()
	ciphertext, err = rsa.EncryptOAEP(hash, rand.Reader, c.pub, plaintext, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// Decrypt the message using the private key.
func (c *RSA) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	if c.priv == nil {
		return nil, crypto.ErrPrivateKeyRequired
	}

	hash := sha512.New()
	plaintext, err = rsa.DecryptOAEP(hash, rand.Reader, c.priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// EncryptionAlgorithm returns the name of the algorithm for adding to the Transaction.
func (c *RSA) EncryptionAlgorithm() string {
	return CipherAlgorithm
}

// ############################################################################
// KeyIdentifier interface
// ############################################################################

// PublicKeySignature implements KeyIdentifier by computing a base64 encoded SHA-256
// hash of the public key serialized as a PKIX public key without PEM encoding. This is
// a prototype method of computing the public key signature and may not match other
// external signature computation methods.
func (c *RSA) PublicKeySignature() (_ string, err error) {
	return signature.New(c.pub)
}

// ############################################################################
// Signer interface
// ############################################################################

// Sign the specified data using the private key, returning the signature. This
// signature is non-deterministic.
func (c *RSA) Sign(data []byte) (signature []byte, err error) {
	digest := sha512.Sum512(data)
	return rsa.SignPSS(rand.Reader, c.priv, gocrypto.SHA512, digest[:], nil)
}

// Verify the signature on the specified data using the public key. A valid
// signature is indicated by returning a nil error.
func (c *RSA) Verify(data, signature []byte) error {
	digest := sha512.Sum512(data)
	return rsa.VerifyPSS(c.pub, gocrypto.SHA512, digest[:], signature, nil)
}

// SignatureAlgorithm returns the name of the signing algorithm.
func (c *RSA) SignatureAlgorithm() string {
	return SignerAlgorithm
}
