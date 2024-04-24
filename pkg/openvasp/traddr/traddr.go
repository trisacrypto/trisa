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

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	tld "github.com/bombsimon/tld-validator"
)

const scheme = "ta"

// Make converts a a URI into a travel address string. If the URI contains a scheme such
// as https:// it is stripped because the https protocol is assumed. If the URI does
// not contain the t=i query parameter, it is added to the travel address.
func Make(uri string) (_ string, err error) {
	var u *URL
	if u, err = Parse(uri); err != nil {
		return "", err
	}

	u.Scheme = ""

	params := u.Query()
	params.Set("t", "i")
	u.RawQuery = params.Encode()

	return Encode(strings.TrimPrefix(u.String(), "//"))
}

// Encode a URI into a travel address string. The encode method is strict, which means
// that the URI must contain the t=i query parameter, a valid TLD, and must not have
// a URI scheme such as https:// which is implied by the travel rule address.
func Encode(uri string) (_ string, err error) {
	// Parse the url to check the scheme, query string, and TLD
	var u *URL
	if u, err = Parse(uri); err != nil {
		return "", err
	}

	if u.Scheme != "" {
		return "", ErrURIScheme
	}

	if tag := u.Query().Get("t"); tag != "i" {
		return "", ErrMissingQueryString
	}

	if err = u.ValidTLD(); err != nil {
		return "", err
	}

	return scheme + checkEncode([]byte(uri)), nil
}

// Decode a travel address string into a raw URI. The decode method is strict, which
// means that the travel address must contain the t=i query parameter, a valid TLD, and
// must not have a URI scheme such as https:// which is implied by the address.
func Decode(traddr string) (_ string, err error) {
	if !strings.HasPrefix(traddr, scheme) {
		return "", ErrUnhandledScheme
	}

	var url []byte
	if url, err = checkDecode(strings.TrimPrefix(traddr, scheme)); err != nil {
		return "", err
	}

	var u *URL
	if u, err = Parse(string(url)); err != nil {
		return "", err
	}

	if u.Scheme != "" {
		return "", ErrURIScheme
	}

	if tag := u.Query().Get("t"); tag != "i" {
		return "", ErrMissingQueryString
	}

	if err = u.ValidTLD(); err != nil {
		return "", err
	}

	return string(url), nil
}

// DecodeURL returns a fully formed URL including the https:// URI scheme.
func DecodeURL(traddr string) (_ string, err error) {
	var rawURL string
	if rawURL, err = Decode(traddr); err != nil {
		return "", err
	}

	var u *URL
	if u, err = Parse(rawURL); err != nil {
		return "", err
	}

	u.Scheme = "https"
	return u.String(), nil
}

// URL extends the net/url package with methods for validing the hostname and TLD.
type URL struct {
	url.URL
}

// Parse a raw url into a URL data structure for validation.
func Parse(rawURL string) (_ *URL, err error) {
	var u *url.URL
	if u, err = url.Parse(rawURL); err != nil {
		return nil, err
	}

	// If the scheme and host is missing, try parsing with an empty scheme
	if (u.Scheme == "" && u.Host == "") || u.Hostname() == "" {
		if u, err = url.Parse("//" + rawURL); err != nil {
			return nil, err
		}
	}

	return &URL{*u}, nil
}

func (u *URL) ValidTLD() error {
	// localhost is a valid TLD
	hostname := u.Hostname()
	if hostname == "localhost" {
		return nil
	}

	// Check if the hostname is an IP address
	if net.ParseIP(hostname) != nil {
		return nil
	}

	// Otherwise validate the TLD
	if !tld.FromDomainName(hostname).IsValid() {
		return fmt.Errorf("%q is an %w", hostname, ErrInvalidTLD)
	}
	return nil
}
