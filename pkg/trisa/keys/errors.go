package keys

import "errors"

var (
	ErrNoPrivateKey = errors.New("no private unsealing key available")
)
