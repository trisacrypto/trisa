package crypto

import "errors"

var (
	ErrCannotSignEmpty       = errors.New("cannot sign empty data")
	ErrMissingCiphertext     = errors.New("empty cipher text")
	ErrHMACSignatureMismatch = errors.New("hmac signature mismatch")
	ErrPrivateKeyRequired    = errors.New("private key required for decryption")
)
