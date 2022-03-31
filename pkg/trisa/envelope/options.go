package envelope

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/trisacrypto/trisa/pkg/trisa/crypto"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto/aesgcm"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto/rsaoeap"
)

type Option func(e *Envelope) error

func FromEnvelope(env *Envelope) Option {
	return func(e *Envelope) error {
		e.msg.Id = env.msg.Id
		e.crypto = env.crypto
		return nil
	}
}

func WithEnvelopeID(id string) Option {
	return func(e *Envelope) error {
		e.msg.Id = id
		return nil
	}
}

func WithTimestamp(ts time.Time) Option {
	return func(e *Envelope) error {
		e.msg.Timestamp = ts.Format(time.RFC3339Nano)
		return nil
	}
}

func WithCrypto(crypto crypto.Crypto) Option {
	return func(e *Envelope) error {
		e.crypto = crypto
		return nil
	}
}

func WithAESGCM(encryptionKey []byte, hmacSecret []byte) Option {
	return func(e *Envelope) (err error) {
		if e.crypto, err = aesgcm.New(encryptionKey, hmacSecret); err != nil {
			return err
		}
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
	return func(e *Envelope) (err error) {
		switch t := key.(type) {
		case *rsa.PublicKey:
			if e.seal, err = rsaoeap.New(t); err != nil {
				return err
			}
		default:
			return fmt.Errorf("could not use %T for sealing", t)
		}
		return nil
	}
}

func WithUnsealingKey(key interface{}) Option {
	return func(e *Envelope) (err error) {
		switch t := key.(type) {
		case *rsa.PrivateKey:
			if e.seal, err = rsaoeap.New(t); err != nil {
				return err
			}
		default:
			return fmt.Errorf("could not use %T for unsealing", t)
		}
		return nil
	}
}

func WithRSAPublicKey(key *rsa.PublicKey) Option {
	return WithSealingKey(key)
}

func WithRSAPrivateKey(key *rsa.PrivateKey) Option {
	return WithUnsealingKey(key)
}
