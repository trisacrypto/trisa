package keys

import "errors"

var (
	ErrNoPrivateKey  = errors.New("no private unsealing key available")
	ErrNoCertificate = errors.New("no certificates found in PEM encoded data")
	ErrMultipleKeys  = errors.New("too many private keys found in PEM encoded data")
)
