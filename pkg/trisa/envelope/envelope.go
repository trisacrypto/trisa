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
An envelope's payload is referred to as "clear" before encryption and after decryption.

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
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto/aesgcm"
	"google.golang.org/protobuf/proto"
)

// Envelope is a wrapper for a trisa.SecureEnvelope that adds cryptographic
// functionality to the protocol buffer payload. An envelope can be in one of three
// states: clear, unsealed, and sealed -- referring to the cryptographic status of the
// wrapped secure envelope. For example, a clear envelope can have its payload directly
// read, but if an envelope is unsealed, then it must be decrypted before the payload
// can be parsed into specific data structures. Similarly, an unsealed envelope can be
// sealed in preparation for sending to a recipient, or remain unsealed for secure long
// term data storage.
type Envelope struct {
	msg     *api.SecureEnvelope
	payload *api.Payload
	crypto  crypto.Crypto
	seal    crypto.Cipher
	parent  *Envelope
}

//===========================================================================
// Envelope Constructors
//===========================================================================

// New creates a new Envelope from scratch with the associated payload and options.
// New should be used to initialize a secure envelope to send to a recipient, usually
// before the first transfer of an information exchange, Open is used thereafter.
func New(payload *api.Payload, opts ...Option) (env *Envelope, err error) {
	// Create a new empty secure envelope
	env = &Envelope{
		msg: &api.SecureEnvelope{
			Id:                  uuid.NewString(),
			Payload:             nil,
			EncryptionKey:       nil,
			EncryptionAlgorithm: "",
			Hmac:                nil,
			HmacSecret:          nil,
			HmacAlgorithm:       "",
			Error:               nil,
			Timestamp:           time.Now().UTC().Format(time.RFC3339Nano),
			Sealed:              false,
			PublicKeySignature:  "",
			TransferState:       api.TransferStateUnspecified,
		},
		payload: payload,
	}

	// Apply the options
	for _, opt := range opts {
		if err = opt(env); err != nil {
			return nil, err
		}
	}
	return env, nil
}

// Wrap initializes an Envelope from an incoming secure envelope without modifying the
// original envelope. The Envelope can then be inspected and managed using cryptographic
// and accessor functions.
func Wrap(msg *api.SecureEnvelope, opts ...Option) (env *Envelope, err error) {
	if msg == nil {
		return nil, ErrNoMessage
	}

	// Wrap the secure envelope with the envelope handler.
	env = &Envelope{
		msg: msg,
	}

	// Apply the options
	for _, opt := range opts {
		if err = opt(env); err != nil {
			return nil, err
		}
	}
	return env, nil
}

// Wrap error initializes an Envelope from a TRISA error to prepare and validate a
// rejection response without going directly to the SecureEnvelope.
func WrapError(reject *api.Error, opts ...Option) (env *Envelope, err error) {
	if env, err = New(nil, opts...); err != nil {
		return nil, err
	}

	if env, err = env.Reject(reject); err != nil {
		return nil, err
	}

	if err = env.ValidateMessage(); err != nil {
		return nil, err
	}

	return env, nil
}

// Create an envelope with an encrypted and sealed payload using the public signing key
// of the TRISA peer (supplied via the WithSealingKey or WithRSAPublicKey options).
// The returned envelope has a parent chain that contains the encryption transformations
// at each step so that you can validate that the payload has been constructed correctly.
func Seal(payload *api.Payload, opts ...Option) (env *Envelope, reject *api.Error, err error) {
	if env, err = New(payload, opts...); err != nil {
		return nil, nil, err
	}

	if env, reject, err = env.Encrypt(); err != nil {
		if reject != nil {
			env, _ = env.Reject(reject)
		}
		return env, reject, err
	}

	if env, reject, err = env.Seal(); err != nil {
		if reject != nil {
			env, _ = env.Reject(reject)
		}
		return env, reject, err
	}

	return env, nil, nil
}

// Open a secure envelope using the private key that is paired with the public key used
// to seal the envelope (provided using the WithUnsealingKey or WithRSAPrivateKey
// options). The returned envelope has a partent chain that contains the decryption
// transformations at each step so tha tyou can validate that the payload has been
// constructed correctly.
func Open(msg *api.SecureEnvelope, opts ...Option) (env *Envelope, reject *api.Error, err error) {
	if env, err = Wrap(msg, opts...); err != nil {
		return nil, nil, err
	}

	if env, reject, err = env.Unseal(); err != nil {
		if reject != nil {
			env, _ = env.Reject(reject)
		}
		return env, reject, err
	}

	if env, reject, err = env.Decrypt(); err != nil {
		if reject != nil {
			env, _ = env.Reject(reject)
		}
		return env, reject, err
	}

	return env, nil, nil
}

//===========================================================================
// Envelope State Transitions
//===========================================================================

// Reject returns a new secure envelope that contains a TRISA rejection error but no
// payload. The original envelope is not modified, the secure envelope is cloned.
func (e *Envelope) Reject(reject *api.Error, opts ...Option) (env *Envelope, err error) {
	env = &Envelope{
		msg: &api.SecureEnvelope{
			Id:                  e.msg.Id,
			Payload:             nil,
			EncryptionKey:       nil,
			EncryptionAlgorithm: "",
			Hmac:                nil,
			HmacSecret:          nil,
			HmacAlgorithm:       "",
			Error:               reject,
			Timestamp:           time.Now().UTC().Format(time.RFC3339Nano),
			Sealed:              false,
			PublicKeySignature:  "",
			TransferState:       api.TransferStateUnspecified,
		},
		parent: e,
	}

	if reject.Retry {
		env.msg.TransferState = api.TransferRepair
	} else {
		env.msg.TransferState = api.TransferRejected
	}

	// Apply the options
	for _, opt := range opts {
		if err = opt(env); err != nil {
			return nil, err
		}
	}
	return env, nil
}

// Update the envelope with a new payload maintaining the original crypto method. This
// is useful to prepare a response to the user, for example updating the ReceivedAt
// timestamp in the payload then re-encrypting the secure envelope to send back to the
// originator. Most often, this method is also used with the WithSealingKey option so
// that the envelope workflow for sealing an envelope can be applied completely.
// The original envelope is not modified, the secure envelope is cloned.
func (e *Envelope) Update(payload *api.Payload, transferState api.TransferState, opts ...Option) (env *Envelope, err error) {
	state := e.State()
	if state != Clear && state != ClearError {
		return nil, fmt.Errorf("cannot update envelope from %q state", state)
	}

	// Clone the envelope
	env = &Envelope{
		msg: &api.SecureEnvelope{
			Id:                  e.msg.Id,
			Payload:             nil,
			EncryptionKey:       e.msg.EncryptionKey,
			EncryptionAlgorithm: e.msg.EncryptionAlgorithm,
			Hmac:                nil,
			HmacSecret:          e.msg.HmacSecret,
			HmacAlgorithm:       e.msg.HmacAlgorithm,
			Error:               e.msg.Error,
			Timestamp:           time.Now().UTC().Format(time.RFC3339Nano),
			Sealed:              false,
			PublicKeySignature:  "",
			TransferState:       transferState,
		},
		crypto:  e.crypto,
		seal:    e.seal,
		payload: payload,
		parent:  e,
	}

	// Apply the options
	for _, opt := range opts {
		if err = opt(env); err != nil {
			return nil, err
		}
	}
	return env, nil
}

// Encrypt the envelope by marshaling the payload, encrypting, and digitally signing it.
// If the original envelope does not have a crypto method (either by the user supplying
// one via options or from an incoming secure envelope) then a new AESGCM crypto is
// created with a random encryption key and hmac secret. Encrypt returns a new unsealed
// envelope that maintains the original crypto and cipher mechanisms. Two types of
// errors may be returned: a rejection error is intended to be returned to the sender
// due to some protocol-specific issue, otherwise an error is returned if the user has
// made some error that requires rehandling of the envelope.
// The original envelope is not modified, the secure envelope is cloned.
func (e *Envelope) Encrypt(opts ...Option) (env *Envelope, reject *api.Error, err error) {
	// The envelope may only be encrypted from the Clear state (e.g. it has a payload)
	state := e.State()
	if state != Clear && state != ClearError {
		return nil, nil, fmt.Errorf("cannot encrypt envelope from %q state", state)
	}

	// Clone the envelope
	env = &Envelope{
		msg: &api.SecureEnvelope{
			Id:                  e.msg.Id,
			Payload:             nil,
			EncryptionKey:       nil,
			EncryptionAlgorithm: "",
			Hmac:                nil,
			HmacSecret:          nil,
			HmacAlgorithm:       "",
			Error:               e.msg.Error,
			Timestamp:           time.Now().UTC().Format(time.RFC3339Nano),
			Sealed:              false,
			PublicKeySignature:  "",
			TransferState:       e.msg.TransferState,
		},
		crypto: e.crypto,
		seal:   e.seal,
		parent: e,
	}

	// Apply the options
	for _, opt := range opts {
		if err = opt(env); err != nil {
			return nil, nil, err
		}
	}

	// Validate the payload before encrypting
	if err = e.ValidatePayload(); err != nil {
		return nil, nil, err
	}

	// Create a new AES-GCM crypto handler if one is not supplied on the envelope. This
	// generates a random encryption key and hmac secret on a per-envelope basis,
	// helping to prevent statistical cryptographic attacks.
	if env.crypto == nil {
		if env.crypto, err = aesgcm.New(nil, nil); err != nil {
			return nil, nil, err
		}
	}

	// Encrypt the payload and fill in the details on the envelope.
	if reject, err = env.encrypt(e.payload); err != nil {
		return nil, reject, err
	}
	return env, reject, nil
}

// Internal encrypt method to update an envelope directly and for short-circuit methods.
func (e *Envelope) encrypt(payload *api.Payload) (_ *api.Error, err error) {
	// TODO: should we inspect the envelope for encryption metadata to create the cipher?
	if e.crypto == nil {
		return nil, ErrCannotEncrypt
	}

	var cleartext []byte
	if cleartext, err = proto.Marshal(payload); err != nil {
		return nil, fmt.Errorf("could not marshal payload: %s", err)
	}

	if e.msg.Payload, err = e.crypto.Encrypt(cleartext); err != nil {
		return nil, fmt.Errorf("could not encrypt payload data: %s", err)
	}

	if e.msg.Hmac, err = e.crypto.Sign(e.msg.Payload); err != nil {
		return nil, fmt.Errorf("could not sign payload data: %s", err)
	}

	// Populate metadata on envelope and reset public key signature since its not present
	e.msg.EncryptionKey = e.crypto.EncryptionKey()
	e.msg.HmacSecret = e.crypto.HMACSecret()
	e.msg.EncryptionAlgorithm = e.crypto.EncryptionAlgorithm()
	e.msg.HmacAlgorithm = e.crypto.SignatureAlgorithm()
	e.msg.Sealed = false
	e.msg.PublicKeySignature = ""

	// Validate the message before returning
	if err = e.ValidateMessage(); err != nil {
		return nil, err
	}
	return nil, nil
}

// The original envelope is not modified, the secure envelope is cloned.
func (e *Envelope) Decrypt(opts ...Option) (env *Envelope, reject *api.Error, err error) {
	// The envelope may only be decrypted from the Unsealed state (e.g. the encryption key is in the clear)
	state := e.State()
	if state != Unsealed && state != UnsealedError {
		return nil, nil, fmt.Errorf("cannot decrypt envelope from %q state", state)
	}

	// Clone the envelope
	env = &Envelope{
		msg: &api.SecureEnvelope{
			Id:                  e.msg.Id,
			Payload:             e.msg.Payload,
			EncryptionKey:       e.msg.EncryptionKey,
			EncryptionAlgorithm: e.msg.EncryptionAlgorithm,
			Hmac:                e.msg.Hmac,
			HmacSecret:          e.msg.HmacSecret,
			HmacAlgorithm:       e.msg.HmacAlgorithm,
			Error:               e.msg.Error,
			Timestamp:           time.Now().UTC().Format(time.RFC3339Nano),
			Sealed:              false,
			PublicKeySignature:  "",
			TransferState:       e.msg.TransferState,
		},
		crypto: e.crypto,
		seal:   e.seal,
		parent: e,
	}

	// Apply the options
	for _, opt := range opts {
		if err = opt(env); err != nil {
			return nil, nil, err
		}
	}

	// Decrypt the payload and update the details on the envelope
	if reject, err = env.decrypt(); err != nil {
		return nil, reject, err
	}
	return env, reject, nil
}

// Internal decrypt method to update an envelope directly and for short-circuit methods.
func (e *Envelope) decrypt() (_ *api.Error, err error) {
	// Validate the message contains the required data to decrypt
	if err = e.ValidateMessage(); err != nil {
		return nil, err
	}

	if e.crypto == nil {
		// Create the cipher from the data on the envelope
		// Check if the encryption algorithms are supported
		// TODO: allow more algorithms by adding composition functionality to the crypto pacakge
		if e.msg.EncryptionAlgorithm != "AES256-GCM" {
			err = fmt.Errorf("unsupported encryption algorithm %q", e.msg.EncryptionAlgorithm)
			return api.Errorf(api.UnhandledAlgorithm, err.Error()), err
		}
		if e.msg.HmacAlgorithm != "HMAC-SHA256" {
			err = fmt.Errorf("unsupported digital signature algorithm %q", e.msg.HmacAlgorithm)
			return api.Errorf(api.UnhandledAlgorithm, err.Error()), err
		}

		if e.crypto, err = aesgcm.New(e.msg.EncryptionKey, e.msg.HmacSecret); err != nil {
			return nil, fmt.Errorf("could not create AES-GCM cipher for payload decryption: %v", err)
		}
	}

	// Verify the HMAC signature of the envelope
	if err = e.crypto.Verify(e.msg.Payload, e.msg.Hmac); err != nil {
		return api.Errorf(api.InvalidSignature, "could not verify HMAC signature"), err
	}

	// Decrypt the payload data
	var data []byte
	if data, err = e.crypto.Decrypt(e.msg.Payload); err != nil {
		return api.Errorf(api.InvalidKey, "could not decrypt payload with encryption key"), err
	}

	// Parse the payload
	e.payload = &api.Payload{}
	if err = proto.Unmarshal(data, e.payload); err != nil {
		return api.Errorf(api.EnvelopeDecodeFail, "could not unmarshal payload from decrypted data"), err
	}

	// Validate the payload
	// TODO: use more specific error such as UNPARSEABLE_TRANSACTION or INCOMPLETE_IDENTITY
	if err = e.ValidatePayload(); err != nil {
		return api.Errorf(api.ValidationError, err.Error()), err
	}

	// Set the payload and the signature to nil now that the message is in clear text
	e.msg.Payload = nil
	e.msg.Hmac = nil
	return nil, nil
}

// Seal the envelope using public key cryptography so that the envelope can only be
// decrypted by the recipient. This method encrypts the encryption key and hmac secret
// using the supplied public key, marking the secure envelope as sealed and updates the
// signature of the public key used to seal the secure envelope. Two types of errors may
// be returned from this method: a rejection error used to communicate to the sender
// that something went wrong and they should resend the envelope or an error that the
// user should handle in their own code base.
// The original envelope is not modified, the secure envelope is cloned.
func (e *Envelope) Seal(opts ...Option) (env *Envelope, reject *api.Error, err error) {
	// The envelope may only be sealed from the Unsealed state (e.g. it has been encrypted)
	state := e.State()
	if state != Unsealed && state != UnsealedError {
		return nil, nil, fmt.Errorf("cannot seal envelope from %q state", state)
	}

	// Clone the envelope
	env = &Envelope{
		msg: &api.SecureEnvelope{
			Id:                  e.msg.Id,
			Payload:             e.msg.Payload,
			EncryptionKey:       e.msg.EncryptionKey,
			EncryptionAlgorithm: e.msg.EncryptionAlgorithm,
			Hmac:                e.msg.Hmac,
			HmacSecret:          e.msg.HmacSecret,
			HmacAlgorithm:       e.msg.HmacAlgorithm,
			Error:               e.msg.Error,
			Timestamp:           time.Now().UTC().Format(time.RFC3339Nano),
			Sealed:              false,
			PublicKeySignature:  "",
			TransferState:       e.msg.TransferState,
		},
		crypto: e.crypto,
		seal:   e.seal,
		parent: e,
	}

	// Apply the options
	for _, opt := range opts {
		if err = opt(env); err != nil {
			return nil, nil, err
		}
	}

	// Seal the payload and fill in the details on the envelope.
	if reject, err = env.sealEnvelope(); err != nil {
		return nil, reject, err
	}
	return env, reject, nil
}

// Internal seal envelope method to update an envelope directly and for short-circuit methods.
func (e *Envelope) sealEnvelope() (_ *api.Error, err error) {
	if e.seal == nil {
		return nil, ErrCannotSeal
	}

	if e.msg.EncryptionKey, err = e.seal.Encrypt(e.msg.EncryptionKey); err != nil {
		return nil, fmt.Errorf("could not seal encryption key: %s", err)
	}

	if e.msg.HmacSecret, err = e.seal.Encrypt(e.msg.HmacSecret); err != nil {
		return nil, fmt.Errorf("could not seal hmac secret: %s", err)
	}

	// Message is now sealed
	e.msg.Sealed = true

	// Add public key signature if the key supports it
	if pks, ok := e.seal.(crypto.KeyIdentifier); ok {
		if e.msg.PublicKeySignature, err = pks.PublicKeySignature(); err != nil {
			return nil, fmt.Errorf("could not compute public key signature: %s", err)
		}
	}
	return nil, nil
}

// Unseal the envelope using public key cryptography so that the envelope can be opened
// and decrypted. This method requires a sealing private key and will return an error
// if one is not available. If the envelope is not able to be opened because the secure
// envelope contains an unknown or improper state a rejection error is returned to
// communicate back to the sender that the envelope could not be unsealed.
// The original envelope is not modified, the secure envelope is cloned.
func (e *Envelope) Unseal(opts ...Option) (env *Envelope, reject *api.Error, err error) {
	// The envelope may only be unsealed from the Sealed state (e.g. it is a valid incoming secure envelope)
	state := e.State()
	if state != Sealed && state != SealedError {
		return nil, nil, fmt.Errorf("cannot seal envelope from %q state", state)
	}

	// Clone the envelope
	env = &Envelope{
		msg: &api.SecureEnvelope{
			Id:                  e.msg.Id,
			Payload:             e.msg.Payload,
			EncryptionKey:       e.msg.EncryptionKey,
			EncryptionAlgorithm: e.msg.EncryptionAlgorithm,
			Hmac:                e.msg.Hmac,
			HmacSecret:          e.msg.HmacSecret,
			HmacAlgorithm:       e.msg.HmacAlgorithm,
			Error:               e.msg.Error,
			Timestamp:           time.Now().UTC().Format(time.RFC3339Nano),
			Sealed:              e.msg.Sealed,
			PublicKeySignature:  e.msg.PublicKeySignature,
			TransferState:       e.msg.TransferState,
		},
		crypto: e.crypto,
		seal:   e.seal,
		parent: e,
	}

	// Apply the options
	for _, opt := range opts {
		if err = opt(env); err != nil {
			return nil, nil, err
		}
	}

	// Seal the payload and fill in the details on the envelope.
	if reject, err = env.unsealEnvelope(); err != nil {
		return nil, reject, err
	}
	return env, reject, nil
}

// Internal unseal envelope method to update envelope directly and for short-circuit methods.
func (e *Envelope) unsealEnvelope() (reject *api.Error, err error) {
	if e.seal == nil {
		return nil, ErrCannotUnseal
	}

	if e.msg.EncryptionKey, err = e.seal.Decrypt(e.msg.EncryptionKey); err != nil {
		return api.Errorf(api.InvalidKey, "could not unseal encryption key").WithRetry(), err
	}

	if e.msg.HmacSecret, err = e.seal.Decrypt(e.msg.HmacSecret); err != nil {
		return api.Errorf(api.InvalidKey, "could not unseal HMAC secret").WithRetry(), err
	}

	// Mark the envelope as unsealed and remove the public key signature
	e.msg.Sealed = false
	e.msg.PublicKeySignature = ""

	return nil, nil
}

//===========================================================================
// Envelope Accessors
//===========================================================================

// ID returns the envelope ID
func (e *Envelope) ID() string {
	return e.msg.Id
}

// UUID returns the envelope ID parsed as a uuid or an error
func (e *Envelope) UUID() (uuid.UUID, error) {
	return uuid.Parse(e.msg.Id)
}

// Proto returns the trisa.SecureEnvelope protocol buffer.
func (e *Envelope) Proto() *api.SecureEnvelope {
	return e.msg
}

// Payload returns the parsed trisa.Payload protocol buffer if available. If the
// envelope is not decrypted then an error is returned.
func (e *Envelope) Payload() (_ *api.Payload, err error) {
	state := e.State()
	if state != Clear && state != ClearError {
		err = fmt.Errorf("envelope is in state %q: payload may be invalid", state)
	}
	return e.payload, err
}

// Error returns the TRISA rejection error on the envelope if it exists
func (e *Envelope) Error() *api.Error {
	// Ensure a nil error is returned if the error is zero-valued
	if e.msg.Error != nil && e.msg.Error.IsZero() {
		return nil
	}
	return e.msg.Error
}

// TransferState returns the TRISA TransferState on the envelope
func (e *Envelope) TransferState() api.TransferState {
	return e.msg.TransferState
}

// IsError returns true if the envelope is in an error state
func (e *Envelope) IsError() bool {
	state := e.State()
	return state == Error || state == ClearError || state == UnsealedError || state == SealedError
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

	return ts.UTC(), err
}

// Crypto returns the cryptographic method used to encrypt/decrypt the payload.
func (e *Envelope) Crypto() crypto.Crypto {
	return e.crypto
}

// Sealer returns the cryptographic cipher method used to seal/unseal the envelope.
func (e *Envelope) Sealer() crypto.Cipher {
	return e.seal
}

// Returns the parent envelope for envelope sequence chain lookups.
func (e *Envelope) Parent() *Envelope {
	return e.parent
}

// Finds the public key signature by looking up the envelope tree until a non-zero
// public key signature is available. Returns empty string if none exists. This is
// useful when you receive an envelope and fully decrypt it, but need to refer back to
// the public key signature that was used in the original sealed envelope.
func (e *Envelope) FindPublicKeySignature() string {
	switch {
	case e.msg.PublicKeySignature != "":
		return e.msg.PublicKeySignature
	case e.parent != nil:
		return e.parent.FindPublicKeySignature()
	default:
		return ""
	}
}

// Find payload returns the nearest payload by looking up the envelope tree until a
// non-nil payload is available. Returns nil if none exists. This is useful when you
// have a payload that is encrypted and sealed, and want to refer back to the original
// payload without keeping track of all of the original envelopes.
func (e *Envelope) FindPayload() *api.Payload {
	switch {
	case e.payload != nil:
		return e.payload
	case e.parent != nil:
		return e.parent.FindPayload()
	default:
		return nil
	}
}

//===========================================================================
// Envelope Validation
//===========================================================================

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
	if err := e.ValidateMessage(); err != nil {
		return Corrupted
	}

	// If there is no payload and there is an error it is in error mode
	if len(e.msg.Payload) == 0 {
		if e.msg.Error == nil || e.msg.Error.IsZero() {
			// Shouldn't happen because of ValidateMessage
			return Unknown
		}
		return Error
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

// ValidateMessage returns an error if the secure envelope does not have the required fields to send.
func (e *Envelope) ValidateMessage() error {
	if e.msg == nil {
		return ErrNoMessage
	}

	if e.msg.Id == "" {
		return ErrNoEnvelopeId
	}

	if _, err := e.UUID(); err != nil {
		return ErrInvalidEnvelopeId
	}

	if e.msg.Timestamp == "" {
		return ErrNoTimestamp
	}

	if _, err := e.Timestamp(); err != nil {
		return ErrInvalidTimestamp
	}

	// The message should have either an error or an encrypted payload
	if len(e.msg.Payload) == 0 {
		if e.msg.Error == nil || e.msg.Error.IsZero() {
			return ErrNoMessageData
		}
		return e.ValidateError()
	}

	// If there is a payload then all payload fields should be set
	if len(e.msg.EncryptionKey) == 0 || e.msg.EncryptionAlgorithm == "" {
		return ErrNoEncryptionInfo
	}

	if len(e.msg.Hmac) == 0 || len(e.msg.HmacSecret) == 0 || e.msg.HmacAlgorithm == "" {
		return ErrNoHMACInfo
	}

	if e.msg.TransferState == api.TransferRejected || e.msg.TransferState == api.TransferRepair {
		return ErrMessageWithErrorState
	}

	// Note: not validating public_key_signature or sealed fields
	return nil
}

// ValidatePayload returns an error if the payload is not ready to be encrypted.
// TODO: should we parse the types of the payload to ensure they're TRISA types?
func (e *Envelope) ValidatePayload() error {
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

// ValidateError returns an error if the error message is missing details
func (e *Envelope) ValidateError() error {
	if e.msg.Error == nil {
		return ErrNoError
	}

	if e.msg.Error.Code == api.Unhandled {
		return ErrMissingErrorCode
	}

	if _, ok := api.Error_Code_name[int32(e.msg.Error.Code)]; !ok {
		return ErrInvalidErrorCode
	}

	if e.msg.Error.Message == "" {
		return ErrMissingErrorMessage
	}

	if e.msg.TransferState != api.TransferRejected && e.msg.TransferState != api.TransferRepair {
		return ErrInvalidErrorTransferState
	}

	return nil
}

// ValidateHMAC checks if the HMAC signature is valid with respect to the HMAC secret
// and encrypted payload. This is generally used for non-repudiation purposes.
func (e *Envelope) ValidateHMAC() (valid bool, err error) {
	// The payload is required to validate the HMAC signature
	if len(e.msg.Payload) == 0 {
		return false, ErrNoPayload
	}

	// An HMAC signature is required for validating the HMAC signature!
	if len(e.msg.Hmac) == 0 {
		return false, ErrNoHMACInfo
	}

	// The cryptography mechanism must have been created, either from encryption or
	// decryption. We do not check the state in case crypto has been added by the user.
	if e.crypto == nil {
		return false, ErrCannotVerify
	}

	// Validate the HMAC signature
	if err = e.crypto.Verify(e.msg.Payload, e.msg.Hmac); err != nil {
		if errors.Is(err, crypto.ErrHMACSignatureMismatch) {
			return false, nil
		}
		return false, err
	}

	// HMAC signature is valid
	return true, nil
}
