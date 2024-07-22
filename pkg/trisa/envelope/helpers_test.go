package envelope_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/envelope"
	"google.golang.org/protobuf/proto"
)

func ExampleSealPayload() {
	// Create compliance payload to send to counterparty. Use key exchange or GDS to
	// fetch the public sealing key of the recipient. See the testdata fixtures for
	// example data. Note: we're loading an RSA private key and extracting its public
	// key for example and testing purposes.
	payload, _ := loadPayloadFixture("testdata/payload.json")
	key, _ := loadPrivateKey("testdata/sealing_key.pem")

	// Seal the payload: encrypting and digitally signing the marshaled protocol buffers
	// with a randomly generated encryption key and HMAC secret, then encrypting those
	// secrets with the public key of the recipient.
	msg, reject, err := envelope.SealPayload(payload, envelope.WithRSAPublicKey(&key.PublicKey))

	// Two types errors may be returned from envelope.Seal
	if err != nil {
		if reject != nil {
			// If both err and reject are non-nil, then a TRISA protocol error occurred
			// and the rejection error can be sent back to the originator if you're
			// sealing the envelope in response to a transfer request
			log.Println(reject.String())
		} else {
			// Otherwise log the error and handle with user-specific code
			log.Fatal(err)
		}
	}

	// Otherwise send the secure envelope to the recipient
	log.Printf("sending secure envelope with id %s", msg.Id)
}

func ExampleOpenPayload() {
	// Receive a sealed secure envelope from the counterparty. Ensure you have the
	// private key paired with the public key identified by the public key signature on
	// the secure envelope in order to unseal and decrypt the payload. See testdata
	// fixtures for example data. Note: we're loading an RSA private key used in other
	// examples for demonstration and testing purposes.
	msg, _ := loadEnvelopeFixture("testdata/sealed_envelope.json")
	key, _ := loadPrivateKey("testdata/sealing_key.pem")

	// Open the secure envelope, unsealing the encryption key and hmac secret with the
	// supplied private key, then decrypting, verifying, and unmarshaling the payload
	// using those secrets.
	payload, reject, err := envelope.OpenPayload(msg, envelope.WithRSAPrivateKey(key))

	// Two types errors may be returned from envelope.Open
	if err != nil {
		if reject != nil {
			// If both err and reject are non-nil, then a TRISA protocol error occurred
			// and the rejection error can be sent back to the originator if you're
			// opening the envelope in response to a transfer request
			out, _ := envelope.Reject(reject, envelope.WithEnvelopeID(msg.Id))
			log.Printf("sending TRISA rejection for envelope %s: %s", out.Id, reject)
		} else {
			// Otherwise log the error and handle with user-specific code
			log.Fatal(err)
		}
	}

	// Handle the payload with your interal compliance processing mechanism.
	log.Printf("received payload sent at %s", payload.SentAt)
}

func TestOneLiners(t *testing.T) {
	payload, err := loadPayloadFixture("testdata/pending_payload.json")
	require.NoError(t, err, "could not load pending payload")

	key, err := loadPrivateKey("testdata/sealing_key.pem")
	require.NoError(t, err, "could not load sealing key")

	// Create an envelope from the payload and the key
	msg, reject, err := envelope.SealPayload(payload, envelope.WithRSAPublicKey(&key.PublicKey))
	require.NoError(t, err, "could not seal envelope")
	require.Nil(t, reject, "unexpected rejection error")

	// Ensure the msg is valid
	require.NotEmpty(t, msg.Id, "no envelope id on the message")
	require.NotEmpty(t, msg.Payload, "no payload on the message")
	require.NotEmpty(t, msg.EncryptionKey, "no encryption key on the message")
	require.NotEmpty(t, msg.EncryptionAlgorithm, "no encryption algorithm on the message")
	require.NotEmpty(t, msg.Hmac, "no hmac signature on the message")
	require.NotEmpty(t, msg.HmacSecret, "no hmac secret on the message")
	require.NotEmpty(t, msg.HmacAlgorithm, "no hmac algorithm on the message")
	require.Empty(t, msg.Error, "unexpected error on the message")
	require.NotEmpty(t, msg.Timestamp, "no timestamp on the message")
	require.True(t, msg.Sealed, "message not marked as sealed")
	require.NotEmpty(t, msg.PublicKeySignature, "no public key signature on the message")

	// Serialize and Deserialize the message
	data, err := proto.Marshal(msg)
	require.NoError(t, err, "could not marshal outgoing message")

	in := &api.SecureEnvelope{}
	require.NoError(t, proto.Unmarshal(data, in), "could not unmarshal incoming message")

	// Open the envelope with the private key
	decryptedPayload, reject, err := envelope.OpenPayload(in, envelope.WithRSAPrivateKey(key))
	require.NoError(t, err, "could not open envelope")
	require.Nil(t, reject, "unexpected rejection error")
	require.True(t, proto.Equal(payload, decryptedPayload), "payloads do not match")
}

func TestReject(t *testing.T) {
	envelopeID := "63c763bc-f820-4a76-b64f-15587ec84a13"

	t.Run("Repair", func(t *testing.T) {
		err := &api.Error{
			Code:    api.Error_EXCEEDED_TRADING_VOLUME,
			Message: "our system is overloaded, please try again later",
			Retry:   true,
		}

		msg, e := envelope.Reject(err, envelope.WithEnvelopeID(envelopeID))
		require.NoError(t, e, "expected no error creating a transfer state")

		require.Equal(t, envelopeID, msg.Id, "envelope does not have correct id")
	})

	t.Run("Reject", func(t *testing.T) {
		err := &api.Error{
			Code:    api.Error_UNKNOWN_ORIGINATOR,
			Message: "the originator specified does not have an account with us",
			Retry:   false,
		}

		msg, e := envelope.Reject(err, envelope.WithEnvelopeID(envelopeID))
		require.NoError(t, e, "expected no error creating a transfer state")

		require.Equal(t, envelopeID, msg.Id, "envelope does not have correct id")
	})

	t.Run("Invalid", func(t *testing.T) {
		msg, e := envelope.Reject(&api.Error{}, envelope.WithEnvelopeID(envelopeID))
		require.ErrorIs(t, e, envelope.ErrNoMessageData)
		require.Nil(t, msg)
	})
}

func TestCheck(t *testing.T) {
	emsg, err := loadEnvelopeFixture("testdata/error_envelope.json")
	require.NoError(t, err, "could not load error envelope fixture")

	terr, iserr := envelope.Check(emsg)
	require.True(t, iserr, "expected error envelope to return iserr true")
	require.Equal(t, api.ComplianceCheckFail, terr.Code)
	require.Equal(t, "specified account has been frozen temporarily", terr.Message)
	require.False(t, terr.Retry)

	for _, path := range []string{"testdata/unsealed_envelope.json", "testdata/sealed_envelope.json"} {
		msg, err := loadEnvelopeFixture(path)
		require.NoError(t, err, "could not load %s", path)

		terr, iserr = envelope.Check(msg)
		require.False(t, iserr)
		require.Nil(t, terr)
	}
}

func TestStatus(t *testing.T) {
	testCases := []struct {
		path  string
		state envelope.State
	}{
		{"testdata/error_envelope.json", envelope.Error},
		{"testdata/unsealed_envelope.json", envelope.Unsealed},
		{"testdata/sealed_envelope.json", envelope.Sealed},
	}

	for i, tc := range testCases {
		msg, err := loadEnvelopeFixture(tc.path)
		require.NoError(t, err, "could not load fixture from %s", tc.path)

		state := envelope.Status(msg)
		require.Equal(t, tc.state, state, "test case %d expected %s got %s", i+1, tc.state, state)
	}
}

func TestTimestamp(t *testing.T) {
	testCases := []struct {
		path     string
		expected time.Time
	}{
		{"testdata/error_envelope.json", time.Time(time.Date(2022, time.January, 27, 8, 21, 43, 0, time.UTC))},
		{"testdata/unsealed_envelope.json", time.Date(2022, time.March, 29, 14, 16, 27, 453444000, time.UTC)},
		{"testdata/sealed_envelope.json", time.Date(2022, time.March, 29, 14, 16, 29, 755212000, time.UTC)},
	}

	for i, tc := range testCases {
		msg, err := loadEnvelopeFixture(tc.path)
		require.NoError(t, err, "could not load fixture from %s", tc.path)

		ts, err := envelope.Timestamp(msg)
		require.NoError(t, err, "timestamp parsing error on test case %d", i)
		require.True(t, tc.expected.Equal(ts), "timestamp mismatch on test case %d", i)
	}
}
