package iso3166

import (
	"errors"
	"fmt"
	"regexp"
)

// Validation Errors
var (
	ErrInvalidAlpha  = errors.New("iso3166: country alpha codes must be 2 or 3 uppercase characters")
	ErrInvalidAlpha2 = errors.New("iso3166: country alpha-2 codes must be 2 uppercase characters")
	ErrInvalidAlpha3 = errors.New("iso3166: country alpha-2 codes must be 3 uppercase characters")
)

// Regular expressions
var (
	reAlpha  = regexp.MustCompile(`^[A-Z]{2,3}$`)
	reAlpha2 = regexp.MustCompile(`^[A-Z]{2}$`)
	reAlpha3 = regexp.MustCompile(`^[A-Z]{3}$`)
)

func NotACountry(alpha string) error {
	return fmt.Errorf("iso3166: %q is not a recognized country code", alpha)
}

func Validate(alpha string) error {
	if !reAlpha.MatchString(alpha) {
		return ErrInvalidAlpha
	}

	switch len(alpha) {
	case 2:
		if _, ok := alpha2[alpha]; !ok {
			return NotACountry(alpha)
		}
	case 3:
		if _, ok := alpha3[alpha]; !ok {
			return NotACountry(alpha)
		}
	default:
		return ErrInvalidAlpha
	}

	return nil
}

func ValidateAlpha2(alpha string) error {
	if !reAlpha2.MatchString(alpha) {
		return ErrInvalidAlpha2
	}

	if _, ok := alpha2[alpha]; !ok {
		return NotACountry(alpha)
	}

	return nil
}

func ValidateAlpha3(alpha string) error {
	if !reAlpha3.MatchString(alpha) {
		return ErrInvalidAlpha3
	}

	if _, ok := alpha3[alpha]; !ok {
		return NotACountry(alpha)
	}

	return nil
}
