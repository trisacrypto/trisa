package envelope

import "errors"

var (
	ErrNoMessage                = errors.New("invalid envelope: no wrapped message")
	ErrNoEnvelopeId             = errors.New("invalid envelope: no envelope id")
	ErrInvalidEnvelopeId        = errors.New("invalid envelope: envelope id is not a uuid")
	ErrNoTimestamp              = errors.New("invalid envelope: no ordering timestamp")
	ErrInvalidTimestamp         = errors.New("invalid envelope: could not parse ordering timestamp in RFC3339 format")
	ErrNoMessageData            = errors.New("invalid envelope: must contain either error or payload")
	ErrNoEncryptionInfo         = errors.New("invalid envelope: missing encryption key or algorithm")
	ErrNoHMACInfo               = errors.New("invalid envelope: missing hmac signature, secret, or algorithm")
	ErrNoPayload                = errors.New("invalid payload: payload has not been decrypted or is missing")
	ErrNoIdentityPayload        = errors.New("invalid payload: payload does not contain identity data")
	ErrNoTransactionPayload     = errors.New("invalid payload: payload does not contain transaction data")
	ErrNoSentAtPayload          = errors.New("invalid payload: sent at timestamp is missing")
	ErrInvalidSentAtPayload     = errors.New("invalid payload: could not parse sent at timestamp in RFC3339 format")
	ErrInvalidReceivedatPayload = errors.New("invalid payload: could not parse received at timestamp in RFC3339 format")
	ErrNoError                  = errors.New("invalid rejection: missing expected rejection error")
	ErrMissingErrorCode         = errors.New("invalid rejection: missing error code")
	ErrMissingErrorMessage      = errors.New("invalid rejection: missing error message")
	ErrInvalidErrorCode         = errors.New("invalid rejection: unknown trisa error code")
	ErrCannotEncrypt            = errors.New("cannot encrypt envelope: no cryptographic handler available")
	ErrCannotSeal               = errors.New("cannot seal envelope: no public key cryptographic handler available")
	ErrCannotUnseal             = errors.New("cannot unseal envelope: no private key cryptographic handler available")
	ErrCannotVerify             = errors.New("cannot verify hmac: no cryptographic handler available")
)
