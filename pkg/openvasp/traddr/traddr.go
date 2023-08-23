/*
Helper functions for encoding and decoding Travel Addresses, which are used to specify
which VASP controls a specific virtual asset address. OpenVASP and TRP have replaced
LNURLs in recent versions to more easily facilitate travel rule transfers and the TRISA
protocol recommends the use of a travel address over the use of an LNURL.

https://gitlab.com/OpenVASP/travel-rule-protocol/-/blob/master/core/specification.md?ref_type=heads#travel-address

The base58 encoding implementation was ported from the https://github.com/btcsuite/btcd
repository per their ISC license.
*/
package traddr

import "strings"

const scheme = "ta"

func Encode(url string) (_ string, err error) {
	return scheme + checkEncode([]byte(url)), nil
}

func Decode(traddr string) (_ string, err error) {
	if !strings.HasPrefix(traddr, scheme) {
		return "", ErrUnhandledScheme
	}

	var url []byte
	if url, err = checkDecode(strings.TrimPrefix(traddr, scheme)); err != nil {
		return "", err
	}
	return string(url), nil
}
