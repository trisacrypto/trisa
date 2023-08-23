package lnurl

import (
	"errors"
	"fmt"
)

var (
	ErrUnhandledScheme        = errors.New("unhandled lnurl scheme")
	ErrMixedCase              = errors.New("string is not all lowercase or all uppercase")
	ErrInvalidBitGroups       = errors.New("only bit groups between 1 and 8 allowed")
	ErrInvalidIncompleteGroup = errors.New("invalid incomplete group")
)

// ErrNonCharsetChar is returned when a character outside of the specific
// bech32 charset is used in the string.
type ErrNonCharsetChar rune

func (e ErrNonCharsetChar) Error() string {
	return fmt.Sprintf("invalid character not part of charset: %v", int(e))
}

// ErrInvalidDataByte is returned when a byte outside the range required for
// conversion into a string was found.
type ErrInvalidDataByte byte

func (e ErrInvalidDataByte) Error() string {
	return fmt.Sprintf("invalid data byte: %v", byte(e))
}

// ErrInvalidChecksum is returned when the extracted checksum of the string
// is different than what was expected.
type ErrInvalidChecksum struct {
	Expected string
	Actual   string
}

func (e ErrInvalidChecksum) Error() string {
	return fmt.Sprintf("expected %v, got %v", e.Expected, e.Actual)
}

// ErrInvalidCharacter is returned when the bech32 string has a character
// outside the range of the supported charset.
type ErrInvalidCharacter rune

func (e ErrInvalidCharacter) Error() string {
	return fmt.Sprintf("invalid character in string: '%c'", rune(e))
}

// ErrInvalidSeparatorIndex is returned when the separator character '1' is
// in an invalid position in the bech32 string.
type ErrInvalidSeparatorIndex int

func (e ErrInvalidSeparatorIndex) Error() string {
	return fmt.Sprintf("invalid separator index %d", int(e))
}
