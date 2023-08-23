package traddr

import "errors"

var (
	// ErrChecksum indicates that the checksum of a check-encoded string does not
	// verify against the checksum.
	ErrChecksum = errors.New("checksum error")

	// ErrInvalidFormat indicates that the check-encoded string has an invalid format.
	ErrInvalidFormat = errors.New("invalid format: version and/or checksum bytes missing")

	// ErrUnhandledScheme indicates that the travel address was not prefixed with ta.
	ErrUnhandledScheme = errors.New("unhandled travel address scheme")

	// ErrMissingQuery indicates that the travel address does not have the query string t=i
	ErrMissingQueryString = errors.New("missing query string")

	// ErrURIScheme indicates that there is a protocol scheme, e.g. https:// in the url
	ErrURIScheme = errors.New("travel address should not contain protocol scheme")

	// ErrInvalidTLD indicates that the URL has no top-level domain
	ErrInvalidTLD = errors.New("host missing top level domain")
)
