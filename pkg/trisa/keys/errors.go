package keys

import "errors"

var (
	ErrNoPrivateKey          = errors.New("no private unsealing key available")
	ErrNoCertificate         = errors.New("no certificates found in PEM encoded data")
	ErrMultipleKeys          = errors.New("too many private keys found in PEM encoded data")
	ErrUnparsableKeyExchange = errors.New("could not parse key exchange data with known key serialization methods")
	ErrNoPublicKey           = errors.New("no public keys found in PEM encoded data")
	ErrTooManyBlocks         = errors.New("too many public key blocks found in PEM encoded data")
)
