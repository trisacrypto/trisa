package openvasp

import (
	"errors"
	"net/http"
)

var (
	ErrEmptyConfirmation     = errors.New("must specify either txid or canceled in confirmation")
	ErrAmbiguousConfirmation = errors.New("cannot specify both txid and canceled in confirmation")
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
