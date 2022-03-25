/*
Package envelope replaces the handler package to provide utilities for encrypting and
decrypting trisa.SecureEnvelopes as well as sealing and unsealing them. SecureEnvelopes
are the unit of transfer in a TRISA transaction and are used to securely exchange
sensitive compliance information. Security is provided through two forms of cryptography:
symmetric cryptography to encrypt and sign the TRISA payload and asymmetric cryptography
to seal the keys and secrets of the envelopes so that only the recipient can open it.

SecureEnvelopes have a lot of terminology, the first paragraph was loaded with it! Some
of these terms are defined below in relation to SecureEnvelopes:

- Symmetric cryptography: encryption that requires both the sender and receiver to have
the same secret keys. The encryption and digital signature of the payload are symmetric,
and the encryption key and HMAC secret are stored on the envelope to ensure that both
counterparties have the secrets required to decrypt the payload.

- Asymmetric cryptography: also referred to as public key cryptography, this type of
cryptography relies on two keys: a public and a private key. When data is encrypted with
the public key, only the private key can be used to decrypt the data. In the case of
SecureEnvelopes, the encryption key and HMAC secret are encrypted using the public key
of the recipient.

- Encrypt/Decrypt: in relation to a SecureEnvelope, this refers to the symmetric
cryptography performed on the payload; these terms are used in contrast to Seal/Unseal.

- Sign/Verify: in relation to a SecureEnvelope, this refers to the digital signature or
HMAC of the encrypted payload. A digital signature ensures that the cryptographic
contents of the payload have not been tampered with and provides the counterparty
non-repudiation to affirm that the payload was received and not tampered with.

- Seal/Unseal: in relation to a SecureEnvelope, this refers to the asymmetric
cryptography performed on the encryption key and hmac secret; these terms are used in
contrast to Encrypt/Decrypt.

This package provides a wrapper, Envelope that is used to create and open
SecureEnvelopes. The envelope workflow is as follows; a new envelope is in the clear,
it is then encrypted and is referred to as "unsealed", then sealed with the public key
of the receipient. When opening an envelope with a private key, the envelope becomes
unsealed, then is decrypted in the clear.

For more details about how to work with envelopes, see the example code.
*/
package envelope

import (
	"time"

	"github.com/google/uuid"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto"
)

// Envelope is a wrapper for a trisa.SecureEnvelope that adds cryptographic
// functionality to the protocol buffer payload. An envelope can be in one of three
// states: clear, unsealed, and sealed -- referring to the cryptographic status of the
// wrapped secure envelope. For example, a clear envelope can have its payload directly
// read, but if the envelope is unsealed, then it must be decrypted before the payload
// can be parsed into specific data structures. Similarly, an unsealed envelope can be
// sealed in preparation for sending to a recipient, or remain unsealed for secure long
// term data storage.
type Envelope struct {
	msg     *api.SecureEnvelope
	payload *api.Payload
	crypto  crypto.Crypto
	seal    crypto.Cipher
}

// New creates a new Envelope from scratch with the associated payload and options.
// New should be used to initialize a secure envelope to send to a recipient, usually
// before the first transfer of an information exchange, Open is used thereafter.
func New(payload *api.Payload) (_ *Envelope, err error) {
	// TODO: handle options
	return &Envelope{
		msg: &api.SecureEnvelope{
			Id:                  uuid.NewString(),
			Payload:             nil,
			EncryptionKey:       nil,
			EncryptionAlgorithm: "",
			Hmac:                nil,
			HmacSecret:          nil,
			HmacAlgorithm:       "",
			Error:               nil,
			Timestamp:           time.Now().Format(time.RFC3339Nano),
			Sealed:              false,
			PublicKeySignature:  "",
		},
		payload: payload,
	}, nil
}

// Open initializes an Envelope from an incoming secure envelope without modifying the
// original envelope. The Envelope can then be inspected and managed using cryptographic
// and accessor functions.
func Open(msg *api.SecureEnvelope) (_ *Envelope, err error) {
	// TODO: handle options
	return &Envelope{
		msg: msg,
	}, nil
}

// Reject returns a new secure envelope that contains a rejection error but no payload.
// The original envelope is not modified, the secure envelope is cloned.
func (e *Envelope) Reject(err *api.Error) (*Envelope, error) {
	return &Envelope{
		msg: &api.SecureEnvelope{
			Id:                  e.msg.Id,
			Payload:             nil,
			EncryptionKey:       nil,
			EncryptionAlgorithm: "",
			Hmac:                nil,
			HmacSecret:          nil,
			HmacAlgorithm:       "",
			Error:               err,
			Timestamp:           time.Now().Format(time.RFC3339Nano),
			Sealed:              false,
			PublicKeySignature:  "",
		},
	}, nil
}

// Proto returns the trisa.SecureEnvelope protocol buffer.
func (e *Envelope) Proto() *api.SecureEnvelope {
	return e.msg
}

// Payload returns the parsed trisa.Payload protocol buffer if available.
func (e *Envelope) Payload() *api.Payload {
	return e.payload
}

// Crypto returns the cryptographic method used to encrypt/decrypt the payload.
func (e *Envelope) Crypto() crypto.Crypto {
	return e.crypto
}

// Seal returns the cryptographic cipher method used to seal/unseal the envelope.
func (e *Envelope) Seal() crypto.Cipher {
	return e.seal
}

// Error returns the error on the envelope if it exists
func (e *Envelope) Error() error {
	// Ensure a nil error is returned if the error is zero-valued
	if e.msg.Error != nil && e.msg.Error.IsZero() {
		return nil
	}
	return e.msg.Error
}

// Timestamp returns the ordering timestamp of the secure envelope. If the timestamp is
// not on the envelope or it cannot be parsed, an error is returned.
func (e *Envelope) Timestamp() (ts time.Time, err error) {
	if e.msg.Timestamp == "" {
		return ts, &api.Error{Code: api.BadRequest, Message: "missing ordering timestamp on secure envelope"}
	}

	// First attempt: parse with nanosecond resolution
	if ts, err = time.Parse(time.RFC3339Nano, e.msg.Timestamp); err != nil {
		// Second attempt: try without nanosecond resolution
		if ts, err = time.Parse(time.RFC3339, e.msg.Timestamp); err != nil {
			return ts, &api.Error{Code: api.BadRequest, Message: "could not parse ordering timestamp on secure envelope as RFC3339 timestamp"}
		}
	}
	return ts, err
}

// State returns the current state of the envelope.
func (e *Envelope) State() State {
	// If a payload exists on the envelope, then it is in the clear
	if e.payload != nil {
		if e.msg.Error == nil || e.msg.Error.IsZero() {
			return Clear
		}
		return ClearError
	}

	// If there is no payload, there should be a secure envelope
	if e.msg == nil {
		return Unknown
	}

	// If there is a secure envelope, it should be valid
	if err := e.ValidMessage(); err != nil {
		return Corrupted
	}

	// Check if the envelope is marked as sealed
	if e.msg.Sealed {
		if e.msg.Error == nil || e.msg.Error.IsZero() {
			return Sealed
		}
		return SealedError
	}

	// Message is unsealed
	if e.msg.Error == nil || e.msg.Error.IsZero() {
		return Unsealed
	}
	return UnsealedError
}

// Returns true if the secure envelope has required fields to send.
func (e *Envelope) ValidMessage() error {
	if e.msg == nil {
		return ErrNoMessage
	}

	if e.msg.Id == "" {
		return ErrNoEnvelopeId
	}

	if e.msg.Timestamp == "" {
		return ErrNoTimestamp
	}

	// The message should have either an error or an encrypted payload
	if len(e.msg.Payload) == 0 {
		if e.msg.Error == nil || e.msg.Error.IsZero() {
			return ErrNoMessageData
		}
		return nil
	}

	// If there is a payload then all payload fields should be set
	if len(e.msg.EncryptionKey) == 0 || e.msg.EncryptionAlgorithm == "" {
		return ErrNoEncryptionInfo
	}

	if len(e.msg.Hmac) == 0 || len(e.msg.HmacSecret) == 0 || e.msg.HmacAlgorithm == "" {
		return ErrNoHMACInfo
	}

	// Note: not validating public_key_signature or sealed fields
	return nil
}

// Returns true if the secure envelope has a valid payload
func (e *Envelope) ValidPayload() error {
	if e.payload == nil {
		return ErrNoPayload
	}

	// Identity payload is required
	if e.payload.Identity == nil {
		return ErrNoIdentityPayload
	}

	// A transaction is required
	if e.payload.Transaction == nil {
		return ErrNoTransactionPayload
	}

	//  The SentAt timestamp is required and should be parseable
	if e.payload.SentAt == "" {
		return ErrNoSentAtPayload
	}

	if _, err := time.Parse(time.RFC3339, e.payload.SentAt); err != nil {
		return ErrInvalidSentAtPayload
	}

	// If the ReceivedAt timestamp is available, it should be parseable
	if e.payload.ReceivedAt != "" {
		if _, err := time.Parse(time.RFC3339, e.payload.ReceivedAt); err != nil {
			return ErrInvalidReceivedatPayload
		}
	}

	return nil
}
