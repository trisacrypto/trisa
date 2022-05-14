/*
Package keys provides interfaces and handlers for managing public/private key pairs that
are used for sealing and unsealing secure envelopes. This package is not intended for
use with symmetric keys that are used to encrypt payloads.
*/
package keys

import api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"

// Key provides a generic interface to either a private key pair or to a public key that
// has been shared in a TRISA key-exchange. The primary use of this top level interface
// is serializing and deserializing keys with the marshaler interface and creating a
// unified mechanism to manage keys on disk.
type Key interface {
	PublicKey
	PrivateKey
	KeyMarshaler

	// Indicates if the Key contains a private key. If this method returns false, then
	// the UnsealingKey() method should always return an error.
	IsPrivate() bool
}

type PublicKey interface {
	KeyIdentifier

	// Return the key object that can be used to seal an envelope, typically an *rsa.PublicKey
	SealingKey() (interface{}, error)

	// Return the protocol buffer exchange key object to send to the counterparty
	Proto() (*api.SigningKey, error)
}

type PrivateKey interface {
	// Return the key object that can be used to unseal an envelope, typically an *rsa.PrivateKey
	UnsealingKey() (interface{}, error)
}

type KeyIdentifier interface {
	PublicKeyAlgorithm() string          // The sealing public key algorithm to identify the key type
	PublicKeySignature() (string, error) // An identifier of the public key for key management
}

type KeyMarshaler interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}
