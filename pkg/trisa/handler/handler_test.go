package handler_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	protocol "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/handler"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// TestEnvelope tests that the Envelope Seal and Open operations work correctly.
func TestEnvelope(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	publicKey := &privateKey.PublicKey

	// Create a new Envelope with a payload
	payload := &protocol.Payload{
		Identity: &anypb.Any{},
	}
	identity := &ivms101.LegalPerson{
		Name: &ivms101.LegalPersonName{
			NameIdentifiers: []*ivms101.LegalPersonNameId{
				{
					LegalPersonName:               "John Doe",
					LegalPersonNameIdentifierType: ivms101.LegalPersonLegal,
				},
			},
		},
	}
	anypb.MarshalFrom(payload.Identity, identity, proto.MarshalOptions{})
	envelope := handler.New("", payload, nil)
	require.NotNil(t, envelope)
	require.NotEmpty(t, envelope.ID)
	require.NotNil(t, envelope.Payload)
	require.NotNil(t, envelope.Cipher)

	// Fail to seal the envelope with an unsupported key type
	_, err = envelope.Seal(nil)
	require.Error(t, err)
	_, err = envelope.Seal(privateKey)
	require.Error(t, err)

	// Seal the envelope with an RSA key
	secure, err := envelope.Seal(publicKey)
	require.NoError(t, err)
	require.NotNil(t, secure)
	require.Equal(t, envelope.ID, secure.Id)
	require.Equal(t, envelope.Cipher.EncryptionAlgorithm(), secure.EncryptionAlgorithm)
	require.Equal(t, envelope.Cipher.SignatureAlgorithm(), secure.HmacAlgorithm)
	require.NotEmpty(t, secure.Payload)
	require.NotEmpty(t, secure.Hmac)
	require.NotEmpty(t, secure.EncryptionKey)
	require.NotEmpty(t, secure.HmacSecret)

	// Fail to open a nil envelope
	_, err = handler.Open(nil, privateKey)
	require.Error(t, err)

	// Fail to open an envelope with an invalid encryption algorithm
	secure.EncryptionAlgorithm = "invalid"
	_, err = handler.Open(secure, privateKey)
	require.Error(t, err)

	// Fail to open an envelope with an invalid hmac algorithm
	secure.EncryptionAlgorithm = envelope.Cipher.EncryptionAlgorithm()
	secure.HmacAlgorithm = "invalid"
	_, err = handler.Open(secure, privateKey)
	require.Error(t, err)

	// Fail to open an envelope with an unsupported key type
	secure.EncryptionAlgorithm = envelope.Cipher.EncryptionAlgorithm()
	secure.HmacAlgorithm = envelope.Cipher.SignatureAlgorithm()
	_, err = handler.Open(secure, nil)
	require.Error(t, err)
	_, err = handler.Open(secure, publicKey)
	require.Error(t, err)

	// Fail to open the envelope using the wrong key
	wrongKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	_, err = handler.Open(secure, wrongKey)
	require.Error(t, err)

	// Successfully opening an envelope
	opened, err := handler.Open(secure, privateKey)
	require.NoError(t, err)
	require.NotNil(t, opened)
	require.Equal(t, envelope.ID, opened.ID)
	require.Equal(t, envelope.Cipher, opened.Cipher)
	require.True(t, proto.Equal(envelope.Payload, opened.Payload), "unexpected payload in opened envelope")
}
