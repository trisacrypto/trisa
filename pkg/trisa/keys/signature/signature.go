/*
Package signature provides a mechanism for computing public key signatures, which are
used to help identify public keys used in sealing TRISA envelopes and select the
matching private key pair when a secure envelope is received.

A PublicKeySignature takes the form of "ALG:base64data" where ALG is one of the valid
hashing algorithms used by this package and base64data contains the hash of the key.
*/
package signature

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const (
	UnknownSignatureAlgorithm Algorithm = iota
	MD5
	SHA256
	SHA512
)

// New creates a public key signature for an *rsa.PublicKey, *ecdsa.PublicKey or
// ed25519.PublicKey (or any key that can be marshalled by x509.MarshalPKIXPublicKey).
// It returns the default SHA256 signature that the package recommends for public key
// identification and matching.
func New(pub interface{}) (_ string, err error) {
	return Sign(pub, SHA256)
}

// Match determes if the public key signature matches the specified public key.
func Match(pks string, pub interface{}) bool {
	algorithm, checksum, err := Parse(pks)
	if err != nil {
		return false
	}

	var sum []byte
	if sum, err = Hash(pub, algorithm); err != nil {
		return false
	}

	return bytes.Equal(sum, checksum)
}

// Sign creates a public key signature from a public key that can bbe marshalled as an
// PKIX public key. It then takes the hash of the marshalled data using the specified
// signature algorithm and returns a string that concatenates the name of the hashing
// algorithm with the raw base64 encoded of the hash sum.
func Sign(pub interface{}, algorithm Algorithm) (_ string, err error) {
	var sum []byte
	if sum, err = Hash(pub, algorithm); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", algorithm, base64.RawStdEncoding.EncodeToString(sum)), nil
}

// Hash returns the checksum of the marshalled PKIX public key
func Hash(pub interface{}, algorithm Algorithm) (_ []byte, err error) {
	var data []byte
	if data, err = x509.MarshalPKIXPublicKey(pub); err != nil {
		return nil, fmt.Errorf("could not marshal pkix public key: %s", err)
	}

	switch algorithm {
	case SHA256:
		sum := sha256.Sum256(data)
		return sum[:], nil
	case MD5:
		sum := md5.Sum(data)
		return sum[:], nil
	case SHA512:
		sum := sha512.Sum512(data)
		return sum[:], nil
	default:
		return nil, fmt.Errorf("unhandled signature algorithm %q", algorithm)
	}
}

// Parse a public key signature into its algorithm and hash components
func Parse(pks string) (algorithm Algorithm, sum []byte, err error) {
	// Determine the signature hash algorithm
	parts := strings.Split(pks, ":")
	if len(parts) != 2 {
		return UnknownSignatureAlgorithm, nil, errors.New("could not parse the signature hash algorithm")
	}

	algorithm = ParseAlgorithm(parts[0])
	if algorithm == UnknownSignatureAlgorithm {
		return UnknownSignatureAlgorithm, nil, errors.New("unknown signature hash algorithm")
	}

	if sum, err = base64.RawStdEncoding.DecodeString(parts[1]); err != nil {
		return UnknownSignatureAlgorithm, nil, errors.New("could not base64 decode the signature hash sum")
	}

	return algorithm, sum, nil
}

type Algorithm uint8

func (s Algorithm) String() string {
	switch s {
	case SHA256:
		return "SHA256"
	case MD5:
		return "MD5"
	case SHA512:
		return "SHA512"
	default:
		return "unknown"
	}
}

func ParseAlgorithm(s string) Algorithm {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case "SHA256":
		return SHA256
	case "MD5":
		return MD5
	case "SHA512":
		return SHA512
	default:
		return UnknownSignatureAlgorithm
	}
}
