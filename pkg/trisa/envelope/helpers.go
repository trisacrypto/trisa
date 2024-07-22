package envelope

import (
	"time"

	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto/aesgcm"
)

//===========================================================================
// Quick One-Line Helper Functions
//===========================================================================

// SealPayload an envelope using the public signing key of the TRISA peer (must be
// supplied via the WithSealingKey or WithRSAPublicKey options). A secure envelope is
// created by marshaling the payload, encrypting it, then sealing the envelope by
// encrypting the encryption key and hmac secret with the public key of the recipient.
// This method returns two types of errors: a rejection error can be returned to the
// sender to indicate that the TRISA protocol failed, otherwise an error is returned for
// the user to handle. This method is a convenience one-liner, for more control of the
// sealing process or to manage intermediate steps, use the Envelope wrapper directly.
func SealPayload(payload *api.Payload, opts ...Option) (_ *api.SecureEnvelope, reject *api.Error, err error) {
	var env *Envelope
	if env, err = New(payload, opts...); err != nil {
		return nil, nil, err
	}

	// Validate the payload before encrypting
	if err = env.ValidatePayload(); err != nil {
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

	if reject, err = env.encrypt(payload); reject != nil || err != nil {
		if reject != nil {
			msg, _ := env.Reject(reject)
			return msg.Proto(), reject, err
		}
		return nil, nil, err
	}

	if reject, err = env.sealEnvelope(); reject != nil || err != nil {
		if reject != nil {
			msg, _ := env.Reject(reject)
			return msg.Proto(), reject, err
		}
		return nil, nil, err
	}

	return env.Proto(), nil, nil
}

// OpenPayload a secure envelope using the private key that is paired with the public
// key that was used to seal the envelope (must be supplied via the WithUnsealingKey or
// WithRSAPrivateKey options). This method decrypts the encryption key and hmac secret,
// decrypts and verifies the payload HMAC signature, then unmarshals the payload and
// verifies its contents. This method returns two types of errors: a rejection error
// that can be returned to the sender to indicate that the TRISA protocol failed,
// otherwise an error is returned for the user to handle. This method is a convenience
// one-liner, for more control of the open envelope process or to manage intermediate
// steps, use the Envelope wrapper directly.
func OpenPayload(msg *api.SecureEnvelope, opts ...Option) (payload *api.Payload, reject *api.Error, err error) {
	var env *Envelope
	if env, err = Wrap(msg, opts...); err != nil {
		return nil, nil, err
	}

	// A rejection here would be related to a sealing key failure
	if reject, err = env.unsealEnvelope(); reject != nil || err != nil {
		return nil, reject, err
	}

	// A rejection here is related to the decryption, verification, and parsing the payload
	if reject, err = env.decrypt(); reject != nil || err != nil {
		return nil, reject, err
	}

	if payload, err = env.Payload(); err != nil {
		return nil, nil, err
	}
	return payload, nil, nil
}

// Reject returns a new rejection error to send to the counterparty
func Reject(reject *api.Error, opts ...Option) (_ *api.SecureEnvelope, err error) {
	var env *Envelope
	if env, err = New(nil, opts...); err != nil {
		return nil, err
	}

	// Add the error to the envelope and validate
	env.msg.Error = reject

	// Determine the transfer state from the rejection
	if reject.Retry {
		env.msg.TransferState = api.TransferRepair
	} else {
		env.msg.TransferState = api.TransferRejected
	}

	// Validate the message and the error
	if err = env.ValidateMessage(); err != nil {
		return nil, err
	}

	return env.Proto(), nil
}

// Check returns any error on the specified envelope as well as a bool that indicates
// if the envelope is in an error state (even if the envelope contains a payload).
func Check(msg *api.SecureEnvelope) (_ *api.Error, iserr bool) {
	env := &Envelope{msg: msg}
	return env.Error(), env.IsError()
}

// Status returns the state the secure envelope is currently in.
func Status(msg *api.SecureEnvelope) State {
	env := &Envelope{msg: msg}
	return env.State()
}

// Timestamp returns the parsed timestamp from the secure envelope.
func Timestamp(msg *api.SecureEnvelope) (time.Time, error) {
	env := &Envelope{msg: msg}
	return env.Timestamp()
}

// Validate is a one-liner for Wrap(msg).ValidateMessage() and can be used to ensure
// that a secure envelope has been correctly initialized and can be processed.
func Validate(msg *api.SecureEnvelope) (err error) {
	var env *Envelope
	if env, err = Wrap(msg); err != nil {
		return err
	}
	return env.ValidateMessage()
}
