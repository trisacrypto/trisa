/*
Helper functions for encoding and decoding (LNURLs) which are used to specify which
VASP controls a specific virtual asset address. OpenVASP and TRP both use LNURLs to
facilitate travel rule transfers and the TRISA protocol recommends its use.

https://www.21analytics.ch/blog/lnurl-for-fatf-travel-rule-software-solution/

The bech32 implementation of the bech32 format specified in BIP 173 was ported from
https://github.com/fiatjaf/go-lnurl per their MIT license and the test cases from the
BIP were ported from the https://github.com/btcsuite/btcd repository per their ISC
license.
*/
package lnurl

import (
	"errors"
	"strings"
)

// Encode a plain-text https URL into a bech32-encoded uppercased lnurl string.
func Encode(url string) (lnurl string, err error) {
	var converted []byte
	if converted, err = convertBits([]byte(url), 8, 5, true); err != nil {
		return "", err
	}

	if lnurl, err = encode("lnurl", converted); err != nil {
		return "", err
	}
	return strings.ToUpper(lnurl), nil
}

// Decode a bech32 encoded lnurl string and returns a plain-text https URL.
func Decode(lnurl string) (url string, err error) {
	lnurl = strings.ToLower(lnurl)

	if !strings.HasPrefix(lnurl, "lnurl1") {
		return "", ErrUnhandledScheme
	}

	// bech32
	tag, data, err := decode(lnurl)
	if err != nil {
		return "", err
	}

	if tag != "lnurl" {
		return "", errors.New("tag is not 'lnurl', but '" + tag + "'")
	}

	converted, err := convertBits(data, 5, 8, false)
	if err != nil {
		return "", err
	}

	return string(converted), nil
}
