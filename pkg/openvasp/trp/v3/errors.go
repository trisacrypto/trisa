package trp

import (
	"errors"
	"net/http"
)

var (
	ErrNilEnvelope           = errors.New("envelope is nil")
	ErrUnknownState          = errors.New("envelope is in an unknown state")
	ErrInvalidState          = errors.New("envelope is invalid")
	ErrEnvelopeError         = errors.New("envelope has an error")
	ErrUnknownTravelAddress  = errors.New("could not identify travel address scheme")
	ErrEmptyConfirmation     = errors.New("invalid: must specify either txid or canceled in confirmation")
	ErrAmbiguousConfirmation = errors.New("invalid: cannot specify both txid and canceled in confirmation")
	ErrEmptyAddress          = errors.New("invalid: payment address is required")
	ErrEmptyCallback         = errors.New("invalid: callback URL is required")
	ErrEmptyResolution       = errors.New("invalid: must specify one of version, approved, or rejected")
	ErrAmbiguousResolution   = errors.New("invalid: may specify only one of version, approved, or rejected")
	ErrEmptyAsset            = errors.New("invalid: must specify either DTI or SLIP-0044 asset identifier")
	ErrNoAmount              = errors.New("invalid: must specify a non-zero amount")
	ErrMissingIVMS101        = errors.New("invalid: must specify IVMS101 identity payload")
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
