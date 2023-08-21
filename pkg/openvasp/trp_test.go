package openvasp_test

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/openvasp"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/envelope"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestEnvelopeExtension(t *testing.T) {
	// Load a sealed envelope fixture
	env, err := loadEnvelopeFixture("testdata/sealed_envelope.json")
	require.NoError(t, err, "could not load unsealed envelope fixture")

	// Convert the sealed envelope to a SealedTRISAEnvelope extension
	ext, err := openvasp.EnvelopeToExtension(env)
	require.NoError(t, err, "could not convert envelope to TRP extension")
	sealed, ok := ext.(*openvasp.SealedTRISAEnvelope)
	require.True(t, ok, "was not converted to a sealed envelope")
	actual := &api.SecureEnvelope{}
	require.NoError(t, protojson.Unmarshal([]byte(sealed.Envelope), actual), "could not unmarshal envelope from JSON string")
	require.Equal(t, env.ID(), actual.Id, "envelope ID does not match")

	// Load an unsealed envelope fixture
	env, err = loadEnvelopeFixture("testdata/unsealed_envelope.json")
	require.NoError(t, err, "could not load unsealed envelope fixture")

	// Convert the unsealed envelope to a UnsealedTRISAEnvelope extension
	ext, err = openvasp.EnvelopeToExtension(env)
	require.NoError(t, err, "could not convert envelope to TRP extension")
	unsealed, ok := ext.(*openvasp.UnsealedTRISAEnvelope)
	require.True(t, ok, "was not converted to an unsealed envelope")
	require.Equal(t, env.ID(), unsealed.Id, "envelope ID does not match")
	proto := env.Proto()
	require.Equal(t, proto.Payload, unsealed.Payload, "payload does not match")
	require.Equal(t, proto.EncryptionKey, unsealed.EncryptionKey, "encryption key does not match")
	require.Equal(t, proto.EncryptionAlgorithm, unsealed.EncryptionAlgorithm, "encryption algorithm does not match")
	require.Equal(t, proto.Hmac, unsealed.HMAC, "hmac does not match")
	require.Equal(t, proto.HmacSecret, unsealed.HMACSecret, "hmac secret does not match")
	require.Equal(t, proto.HmacAlgorithm, unsealed.HMACAlgorithm, "hmac algorithm does not match")

	// Create a clear envelope from the payload fixture
	payload, err := loadPayloadFixture("testdata/payload.json")
	require.NoError(t, err, "could not load payload")
	//key, err := loadPrivateKey("testdata/sealing_key.pem")
	//require.NoError(t, err, "could not load private key")
	env, err = envelope.New(payload) //, envelope.WithRSAPublicKey(&key.PublicKey))
	require.NoError(t, err, "could not create clear envelope")

	// Convert the clear envelope to an UnsealedTRISAEnvelope extension
	ext, err = openvasp.EnvelopeToExtension(env)
	require.NoError(t, err, "could not convert envelope to TRP extension")
	unsealed, ok = ext.(*openvasp.UnsealedTRISAEnvelope)
	require.True(t, ok, "was not converted to an unsealed envelope")
	require.Equal(t, env.ID(), unsealed.Id, "envelope ID does not match")
	proto = env.Proto()
	require.Equal(t, proto.Payload, unsealed.Payload, "payload does not match")
	require.Equal(t, proto.EncryptionKey, unsealed.EncryptionKey, "encryption key does not match")
	require.Equal(t, proto.EncryptionAlgorithm, unsealed.EncryptionAlgorithm, "encryption algorithm does not match")
	require.Equal(t, proto.Hmac, unsealed.HMAC, "hmac does not match")
	require.Equal(t, proto.HmacSecret, unsealed.HMACSecret, "hmac secret does not match")
	require.Equal(t, proto.HmacAlgorithm, unsealed.HMACAlgorithm, "hmac algorithm does not match")
}

func loadEnvelopeFixture(path string) (env *envelope.Envelope, err error) {
	msg := &api.SecureEnvelope{}
	if err = loadFixture(path, msg); err != nil {
		return nil, err
	}
	return envelope.Wrap(msg)
}

func loadPayloadFixture(path string) (payload *api.Payload, err error) {
	payload = &api.Payload{}
	if err = loadFixture(path, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func loadPrivateKey(path string) (key *rsa.PrivateKey, err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("could not decode PEM data")
	}

	var keyt interface{}
	if keyt, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return nil, err
	}

	return keyt.(*rsa.PrivateKey), nil
}

// Helper method to load a fixture from JSON
func loadFixture(path string, m proto.Message) (err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return err
	}

	opts := protojson.UnmarshalOptions{
		AllowPartial:   false,
		DiscardUnknown: false,
	}
	if err = opts.Unmarshal(data, m); err != nil {
		return err
	}
	return nil
}
