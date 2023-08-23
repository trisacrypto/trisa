package openvasp_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	. "github.com/trisacrypto/trisa/pkg/openvasp"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	generic "github.com/trisacrypto/trisa/pkg/trisa/data/generic/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/envelope"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestTravelAddress(t *testing.T) {
	testCases := []struct {
		address  string
		expected string
		err      error
	}{
		{
			"lnurl1dp68gurn8ghj7cn9dejkv6trd9shy7fwvdhk6tm5wfcr7arpvu7hgunpwejkcun4d3jkjmn3w45hy7gmsy37e",
			"https://beneficiary.com/trp?tag=travelruleinquiry",
			nil,
		},
		{
			"LNURL1DP68GURN8GHJ7MMSV4H8VCTNWQH8GETNWSKKUET59E5K7TE3XGEN7ARPVU7KJMN3W45HY7GF5KZ53",
			"https://openvasp.test-net.io/123?tag=inquiry",
			nil,
		},
		{
			"taGw1e4cjuujwyHBV51aspbLBUbcDhebW7ss8iF6dLEn19WfyQZt6HaUTStSE2YcadkFG",
			"https://beneficiary.com/x/12345?color=red&foo=bar&t=i",
			nil,
		},
		{
			"https://beneficiary.com/x/12345?t=i",
			"https://beneficiary.com/x/12345?t=i",
			nil,
		},
		{
			"https://beneficiary.com/x/12345",
			"https://beneficiary.com/x/12345?t=i",
			nil,
		},
		{
			"http://beneficiary.com/x/12345",
			"http://beneficiary.com/x/12345?t=i",
			nil,
		},
		{
			"foo", "",
			ErrUnknownTravelAddress,
		},
	}

	for i, tc := range testCases {
		info := &TRPInfo{Address: tc.address, APIVersion: APIVersion, RequestIdentifier: "704c548a-70af-480c-af83-6fb7803df85c"}
		actual, err := info.GetURL()
		if tc.err != nil {
			require.Error(t, err, "expected error on test case %d", i)
			require.ErrorIs(t, err, tc.err, "unexpected error on test case %d", i)
		} else {
			require.NoError(t, err, "expected no error on test case %d", i)
			require.Equal(t, tc.expected, actual, "unexpected mismatch on test case %d", i)
		}
	}
}

func TestEnvelopePayload(t *testing.T) {
	// Nil envelope should return an error
	_, err := EnvelopeToPayload(nil)
	require.ErrorIs(t, err, ErrNilEnvelope, "nil envelope should return an error")

	// Load a sealed envelope fixture
	env, err := loadEnvelope("testdata/sealed_envelope.json")
	require.NoError(t, err, "could not load unsealed envelope fixture")

	// Convert the sealed envelope to a TRP payload
	payload, err := EnvelopeToPayload(env)
	require.NoError(t, err, "could not convert envelope to TRP extension")
	ext, ok := payload.Extensions[SealedTRISAExtension]
	require.True(t, ok, "payload does not contain sealed envelope extension")
	sealed, ok := ext.(*SealedTRISAEnvelope)
	require.True(t, ok, "expected sealed envelope extension")
	actual := &api.SecureEnvelope{}
	require.NoError(t, protojson.Unmarshal([]byte(sealed.Envelope), actual), "could not unmarshal sealed envelope from JSON string")
	require.Equal(t, env.ID(), actual.Id, "envelope ID does not match")

	// Load an unsealed envelope fixture
	env, err = loadEnvelope("testdata/unsealed_envelope.json")
	require.NoError(t, err, "could not load unsealed envelope fixture")

	// Convert the unsealed envelope to a TRP payload
	payload, err = EnvelopeToPayload(env)
	require.NoError(t, err, "could not convert envelope to TRP payload")
	ext, ok = payload.Extensions[UnsealedTRISAExtension]
	require.True(t, ok, "payload does not contain unsealed envelope extension")
	unsealed, ok := ext.(*UnsealedTRISAEnvelope)
	require.True(t, ok, "was not converted to an unsealed envelope")
	require.Equal(t, env.ID(), unsealed.Id, "envelope ID does not match")
	se := env.Proto()
	require.Equal(t, se.Payload, unsealed.Payload, "payload does not match")
	require.Equal(t, se.EncryptionKey, unsealed.EncryptionKey, "encryption key does not match")
	require.Equal(t, se.EncryptionAlgorithm, unsealed.EncryptionAlgorithm, "encryption algorithm does not match")
	require.Equal(t, se.Hmac, unsealed.HMAC, "hmac does not match")
	require.Equal(t, se.HmacSecret, unsealed.HMACSecret, "hmac secret does not match")
	require.Equal(t, se.HmacAlgorithm, unsealed.HMACAlgorithm, "hmac algorithm does not match")

	// Create a clear envelope with a TRISA payload
	fixture := &api.Payload{
		SentAt: time.Now().Format(time.RFC3339),
	}
	identity, err := loadIdentity("testdata/identity.json")
	require.NoError(t, err, "could not load identity payload fixture")
	transaction, err := loadTransaction("testdata/transaction.json")
	require.NoError(t, err, "could not load transaction payload fixture")
	fixture.Identity, err = anypb.New(identity)
	require.NoError(t, err, "could not marshal identity payload")
	fixture.Transaction, err = anypb.New(transaction)
	require.NoError(t, err, "could not marshal transaction payload")
	env, err = envelope.New(fixture)
	require.NoError(t, err, "could not create clear envelope")

	// Convert the clear envelope to a TRP payload
	payload, err = EnvelopeToPayload(env)
	require.NoError(t, err, "could not convert envelope to TRP payload")
	require.Equal(t, &Asset{SLIP044: 0}, payload.Asset, "asset type does not match")
	require.Equal(t, transaction.Amount, payload.Amount, "amount does not match")
	require.True(t, proto.Equal(payload.IVMS101, identity), "identity does not match")
	require.Nil(t, payload.Extensions, "payload should not contain any extensions")
}

func loadEnvelope(path string) (env *envelope.Envelope, err error) {
	msg := &api.SecureEnvelope{}
	if err = loadFixture(path, msg); err != nil {
		return nil, err
	}
	return envelope.Wrap(msg)
}

func loadIdentity(path string) (identity *ivms101.IdentityPayload, err error) {
	identity = &ivms101.IdentityPayload{}
	if err = loadFixture(path, identity); err != nil {
		return nil, err
	}
	return identity, nil
}

func loadTransaction(path string) (transaction *generic.Transaction, err error) {
	transaction = &generic.Transaction{}
	if err = loadFixture(path, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
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

func TestConfirmationValidate(t *testing.T) {
	testCases := []struct {
		confirm *Confirmation
		err     error
	}{
		{&Confirmation{}, ErrEmptyConfirmation},
		{&Confirmation{TRP: &TRPInfo{APIVersion: APIVersion}}, ErrEmptyConfirmation},
		{&Confirmation{TXID: "foo", Canceled: "bar"}, ErrAmbiguousConfirmation},
		{&Confirmation{TXID: "foo"}, nil},
		{&Confirmation{Canceled: "bar"}, nil},
	}

	for i, tc := range testCases {
		err := tc.confirm.Validate()
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "test case %d failed with mismatched error", i)
		} else {
			require.NoError(t, err, "test case %d failed: expected valid confirmation", i)
		}
	}
}
