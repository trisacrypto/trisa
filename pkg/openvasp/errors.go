package openvasp

import "errors"

var (
	ErrNilEnvelope   = errors.New("envelope is nil")
	ErrUnknownState  = errors.New("envelope is in an unknown state")
	ErrInvalidState  = errors.New("envelope is invalid")
	ErrEnvelopeError = errors.New("envelope has an error")
)
