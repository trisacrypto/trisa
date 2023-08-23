package openvasp

import (
	"errors"
	"net/http"
)

var (
	ErrNilEnvelope           = errors.New("envelope is nil")
	ErrUnknownState          = errors.New("envelope is in an unknown state")
	ErrInvalidState          = errors.New("envelope is invalid")
	ErrEnvelopeError         = errors.New("envelope has an error")
	ErrEmptyConfirmation     = errors.New("must specify either txid or canceled in confirmation")
	ErrAmbiguousConfirmation = errors.New("cannot specify both txid and canceled in confirmation")
	ErrUnknownTravelAddress  = errors.New("could not identify travel address scheme")
)

type StatusError struct {
	Code    int    // HTTP Status Code to return
	Message string // Message to return, if empty, default status message is used
}

func (e *StatusError) Error() string {
	if e.Message == "" {
		return http.StatusText(e.Code)
	}
	return e.Message
}
