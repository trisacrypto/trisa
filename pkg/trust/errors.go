package trust

import "errors"

// Standard errors for error type checking
var (
	ErrDecodePrivateKey  = errors.New("could not decode PEM private key")
	ErrDecodeCertificate = errors.New("could not decode PEM certificate")
	ErrDecodeCSR         = errors.New("could not decode PEM certificate request")
	ErrNoCertificates    = errors.New("provider does not contain any certificates")
	ErrKeyRequired       = errors.New("private key required")
	ErrZipEmpty          = errors.New("zip archive contains no providers")
	ErrZipTooMany        = errors.New("multiple providers in zip, is this a provider pool?")
)
