package handler

import (
	"crypto/rsa"

	"github.com/google/uuid"
	protocol "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	apierr "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1/errors"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto/aesgcm"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto/rsaoeap"
	"google.golang.org/protobuf/proto"
)

// Envelope wraps a SecureEnvelope containing all of the information necessary to access
// the payload data. The envelope can be edited and resealed to simplify TRISA exchanges.
type Envelope struct {
	ID      string
	Payload *protocol.Payload
	Cipher  crypto.Crypto
}

// New creates a new envelope, generating an ID if the ID is empty and creating a new
// AES-GCM cipher if the cipher is nil.
func New(id string, payload *protocol.Payload, cipher crypto.Crypto) *Envelope {
	if id == "" {
		id = uuid.New().String()
	}

	if cipher == nil {
		var err error
		if cipher, err = aesgcm.New(nil, nil); err != nil {
			panic(err)
		}
	}

	return &Envelope{ID: id, Payload: payload, Cipher: cipher}
}

// Open a secure envelope using the private signing key paired with the public key that
// was used to encrypt the symmetric payload encryption key. The open method decrypts
// the payload key, then decrypts and verifies the payload data using the algorithm
// information stored in the envelope. It returns a data structure with discovered
// cipher and decrypted Payload for access. On error returns *protocol.Error so that
// the error can be directly returned to the client.
func Open(in *protocol.SecureEnvelope, key interface{}) (*Envelope, *protocol.Error) {
	var (
		err           error
		asym          crypto.Cipher
		encryptionKey []byte
		hmacSecret    []byte
		payloadData   []byte
	)

	if in == nil {
		return nil, apierr.Errorf(apierr.BadRequest, "missing secure envelope")
	}

	// Check the algorithms to make sure they're supported
	// TODO: allow more algorithms than just AES256-GCM and HMAC-SHA256
	if in.EncryptionAlgorithm != "AES256-GCM" {
		return nil, apierr.Errorf(apierr.UnhandledAlgorithm, "%s encryption unsupported", in.EncryptionAlgorithm)
	}
	if in.HmacAlgorithm != "HMAC-SHA256" {
		return nil, apierr.Errorf(apierr.UnhandledAlgorithm, "%s hmac unsupported", in.HmacAlgorithm)
	}

	// Create the asymmetric cipher with the private key to decrypt the payload key.
	// TODO: add other asymmetric encryption algorithms
	switch t := key.(type) {
	case *rsa.PrivateKey:
		if asym, err = rsaoeap.New(t); err != nil {
			return nil, apierr.Errorf(apierr.InternalError, "could not create RSA cipher for asymmetric decryption: %s", err)
		}
	default:
		return nil, apierr.Errorf(apierr.UnhandledAlgorithm, "could not use %T for asymetric decryption", t)
	}

	// Decrypt the payload encryption key and hmac secret.
	if encryptionKey, err = asym.Decrypt(in.EncryptionKey); err != nil {
		return nil, apierr.WithRetry(apierr.Errorf(apierr.InvalidKey, "encryption key signed incorrectly: %s", err))
	}
	if hmacSecret, err = asym.Decrypt(in.HmacSecret); err != nil {
		return nil, apierr.WithRetry(apierr.Errorf(apierr.InvalidKey, "hmac secret signed incorrectly: %s", err))
	}

	// Create the envelope with the AES-GCM Cipher
	// TODO: allow multiple Cipher/Signer types
	env := &Envelope{ID: in.Id, Payload: &protocol.Payload{}}
	if env.Cipher, err = aesgcm.New(encryptionKey, hmacSecret); err != nil {
		return nil, apierr.Errorf(apierr.InternalError, "could not create AES-GCM cipher for symmetric decryption: %s", err)
	}

	// Verify the signature and decrypt the payload
	if err = env.Cipher.Verify(in.Payload, in.Hmac); err != nil {
		return nil, apierr.Errorf(apierr.InvalidSignature, "could not verify HMAC signature: %s", err)
	}

	if payloadData, err = env.Cipher.Decrypt(in.Payload); err != nil {
		return nil, apierr.Errorf(apierr.InvalidKey, "could not decrypt payload with key: %s", err)
	}

	// Parse the payload
	if err = proto.Unmarshal(payloadData, env.Payload); err != nil {
		return nil, apierr.Errorf(apierr.EnvelopeDecodeFail, "could not unmarshal payload from decrypted data: %s", err)
	}
	return env, nil
}

// Seal an envelope using the public signing key of the TRISA peer. The envelope uses
// the internal Cipher to encrypt the Payload then encrypts the keys in the Cipher with
// the public key. On error returns *protocol.Error so that the error can be directly
// returned to the client.
func (e *Envelope) Seal(key interface{}) (out *protocol.SecureEnvelope, _ *protocol.Error) {
	var (
		err         error
		asym        crypto.Cipher
		payloadData []byte
	)

	// Create the asymmetric cipher with the private key to decrypt the payload key.
	// TODO: add other asymmetric encryption algorithms
	switch t := key.(type) {
	case *rsa.PublicKey:
		if asym, err = rsaoeap.New(t); err != nil {
			return nil, apierr.Errorf(apierr.InternalError, "could not create RSA cipher for asymmetric encryption: %s", err)
		}
	default:
		return nil, apierr.Errorf(apierr.UnhandledAlgorithm, "could not use %T for asymmetric encryption", t)
	}

	if payloadData, err = proto.Marshal(e.Payload); err != nil {
		return nil, apierr.Errorf(apierr.InternalError, "could not marshal payload data: %s", err)
	}

	// Create the secure envelope with algorithm details
	out = &protocol.SecureEnvelope{
		Id:                  e.ID,
		EncryptionAlgorithm: e.Cipher.EncryptionAlgorithm(),
		HmacAlgorithm:       e.Cipher.SignatureAlgorithm(),
	}

	// Encrypt and sign the payload
	if out.Payload, err = e.Cipher.Encrypt(payloadData); err != nil {
		return nil, apierr.Errorf(apierr.InternalError, "could not encrypt payload data: %s", err)
	}

	if out.Hmac, err = e.Cipher.Sign(out.Payload); err != nil {
		return nil, apierr.Errorf(apierr.InternalError, "could not sign payload data: %s", err)
	}

	// Encrypt the payload encryption key and hmac secret
	// TODO: Update the crypto interface to allow fetching the key and secret
	keys := e.Cipher.(*aesgcm.AESGCM)
	if out.EncryptionKey, err = asym.Encrypt(keys.EncryptionKey()); err != nil {
		return nil, apierr.Errorf(apierr.InternalError, "could not encrypt payload encryption key: %s", err)
	}

	if out.HmacSecret, err = asym.Encrypt(keys.HMACSecret()); err != nil {
		return nil, apierr.Errorf(apierr.InternalError, "could not encrypt hmac secret: %s", err)
	}

	return out, nil
}

// Seal a payload using the specified symmetric key cipher and public signing key. This
// is a convienience method for users who do not want to directly Seal an Envelope.
func Seal(id string, payload *protocol.Payload, cipher crypto.Crypto, key interface{}) (*protocol.SecureEnvelope, *protocol.Error) {
	env := &Envelope{ID: id, Payload: payload, Cipher: cipher}
	return env.Seal(key)
}
