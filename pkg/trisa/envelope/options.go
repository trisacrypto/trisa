package envelope

import (
	"crypto/rsa"

	"github.com/trisacrypto/trisa/pkg/trisa/crypto"
)

type Option func(e *Envelope) error

func WithCrypto(crypto crypto.Crypto) Option {
	return func(e *Envelope) error {
		e.crypto = crypto
		return nil
	}
}

func WithAESGCM(encrptionKey []byte, hmacSecret []byte) Option {
	return func(e *Envelope) error {
		return nil
	}
}

func WithSeal(seal crypto.Cipher) Option {
	return func(e *Envelope) error {
		e.seal = seal
		return nil
	}
}

func WithSealingKey(key interface{}) Option {
	return func(e *Envelope) error {
		return nil
	}
}

func WithUnsealingKey(key interface{}) Option {
	return func(e *Envelope) error {
		return nil
	}
}

func WithRSAPublicKey(key *rsa.PublicKey) Option {
	return WithSealingKey(key)
}

func WithRSAPrivateKey(key *rsa.PrivateKey) Option {
	return WithUnsealingKey(key)
}
