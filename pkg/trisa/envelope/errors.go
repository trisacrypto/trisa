package envelope

import "errors"

var (
	ErrNoMessage                = errors.New("invalid envelope: no wrapped message")
	ErrNoEnvelopeId             = errors.New("invalid envelope: no envelope id")
	ErrNoTimestamp              = errors.New("invalid envelope: no ordering timestamp")
	ErrNoMessageData            = errors.New("invalid envelope: must contain either error or payload")
	ErrNoEncryptionInfo         = errors.New("invalid envelope: missing encryption key or algorithm")
	ErrNoHMACInfo               = errors.New("invalid envelope: missing hmac signature, secret, or algorithm")
	ErrNoPayload                = errors.New("invalid payload: payload has not been decrypted")
	ErrNoIdentityPayload        = errors.New("invalid payload: payload does not contain identity data")
	ErrNoTransactionPayload     = errors.New("invalid payload: payload does not contain transaction data")
	ErrNoSentAtPayload          = errors.New("invalid payload: sent at timestamp is missing")
	ErrInvalidSentAtPayload     = errors.New("invalid payload: could not parse sent at timestamp in RFC3339 format")
	ErrInvalidReceivedatPayload = errors.New("invalid payload: could not parse received at timestamp in RFC3339 format")
	ErrCannotEncrypt            = errors.New("cannot encrypt envelope: no cryptographic handler available")
	ErrCannotSeal               = errors.New("cannot seal envelope: no public key cryptographic handler available")
)
