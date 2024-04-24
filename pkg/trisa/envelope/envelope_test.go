package envelope_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	generic "github.com/trisacrypto/trisa/pkg/trisa/data/generic/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/envelope"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func ExampleSeal() {
	// Create compliance payload to send to counterparty. Use key exchange or GDS to
	// fetch the public sealing key of the recipient. See the testdata fixtures for
	// example data. Note: we're loading an RSA private key and extracting its public
	// key for example and testing purposes.
	payload, _ := loadPayloadFixture("testdata/payload.json")
	key, _ := loadPrivateKey("testdata/sealing_key.pem")

	// Seal the payload: encrypting and digitally signing the marshaled protocol buffers
	// with a randomly generated encryption key and HMAC secret, then encrypting those
	// secrets with the public key of the recipient.
	msg, reject, err := envelope.Seal(payload, envelope.WithRSAPublicKey(&key.PublicKey))

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

func Example_create() {
	// Create compliance payload to send to counterparty. Use key exchange or GDS to
	// fetch the public sealing key of the recipient. See the testdata fixtures for
	// example data. Note: we're loading an RSA private key and extracting its public
	// key for example and testing purposes.
	payload, _ := loadPayloadFixture("testdata/payload.json")
	key, _ := loadPrivateKey("testdata/sealing_key.pem")

	// Envelopes transition through the following states: clear --> unsealed --> sealed.
	// First create a new envelope in the clear state with the public key of the
	// recipient that will eventually be used to seal the envelope.
	env, _ := envelope.New(payload, envelope.WithRSAPublicKey(&key.PublicKey))

	// Marshal the payload, generate random encryption and hmac secrets, and encrypt
	// the payload, creating a new envelope in the unsealed state.
	env, reject, err := env.Encrypt()

	// Two types of errors are returned from Encrypt and Seal
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

	// Seal the envelope by encrypting the encryption key and hmac secret on the secure
	// envelope with the public key of the recipient passed in at the first step.
	// Handle the reject and err errors as above.
	env, reject, err = env.Seal()

	// Fetch the secure envelope and send it.
	msg := env.Proto()
	log.Printf("sending secure envelope with id %s", msg.Id)
}

func ExampleOpen() {
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
	payload, reject, err := envelope.Open(msg, envelope.WithRSAPrivateKey(key))

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

func Example_parse() {
	// Receive a sealed secure envelope from the counterparty. Ensure you have the
	// private key paired with the public key identified by the public key signature on
	// the secure envelope in order to unseal and decrypt the payload. See testdata
	// fixtures for example data. Note: we're loading an RSA private key used in other
	// examples for demonstration and testing purposes.
	msg, _ := loadEnvelopeFixture("testdata/sealed_envelope.json")
	key, _ := loadPrivateKey("testdata/sealing_key.pem")

	// Envelopes transition through the following states: sealed --> unsealed --> clear.
	// First wrap the incoming envelope in the sealed state.
	env, _ := envelope.Wrap(msg)

	// Unseal the envelope using the private key loaded above; this decrypts the
	// encryption key and hmac secret using asymmetric encryption and returns a new
	// unsealed envelope.
	env, reject, err := env.Unseal(envelope.WithRSAPrivateKey(key))

	// Two types of errors are returned from Unseal and Decrypt
	if err != nil {
		if reject != nil {
			// If both err and reject are non-nil, then a TRISA protocol error occurred
			// and the rejection error can be sent back to the originator if you're
			// unsealing the envelope in response to a transfer request.
			out, _ := env.Reject(reject)
			log.Printf("sending TRISA rejection for envelope %s: %s", out.ID(), reject)
		} else {
			// Otherwise log the error and handle with user-specific code
			log.Fatal(err)
		}
	}

	// Decrypt the envelope using the unsealed secrets, verify the HMAC signature, then
	// unmarshal and verify the payload into new envelope in the clear state.
	// Handle the reject and err errors as above.
	env, reject, err = env.Decrypt()

	// Handle the payload with your interal compliance processing mechanism.
	payload, _ := env.Payload()
	log.Printf("received payload sent at %s", payload.SentAt)
}

// Test the creation of an envelope from scratch, moving it through each state.
func TestSendEnvelopeWorkflow(t *testing.T) {
	payload, err := loadPayloadFixture("testdata/payload.json")
	require.NoError(t, err, "could not load payload")

	key, err := loadPrivateKey("testdata/sealing_key.pem")
	require.NoError(t, err, "could not load sealing key")

	env, err := envelope.New(payload, envelope.WithRSAPublicKey(&key.PublicKey))
	require.NoError(t, err, "could not create envelope with no payload and no options")
	require.Equal(t, envelope.Clear, env.State(), "expected clear state not %q", env.State())

	eenv, reject, err := env.Encrypt()
	require.NoError(t, err, "could not encrypt envelope")
	require.Nil(t, reject, "expected no API error returned from encryption")
	require.NotSame(t, env, eenv, "Encrypt should return a clone of the original envelope")
	require.Equal(t, envelope.Unsealed, eenv.State(), "expected unsealed state not %q", eenv.State())

	senv, reject, err := eenv.Seal()
	require.NoError(t, err, "could not seal envelope")
	require.Nil(t, reject, "expected no API error returned from sealing")
	require.NotSame(t, eenv, senv, "Seal should return a clone of the original envelope")
	require.Equal(t, envelope.Sealed, senv.State(), "expected sealed state not %q", senv.State())

	// Fetch the message and check that it is ready to send
	msg := senv.Proto()
	require.NotEmpty(t, msg.Id, "message is missing an envelope ID")
	require.NotEmpty(t, msg.Payload, "message is missing encrypted payload")
	require.NotEmpty(t, msg.EncryptionKey, "message is missing encryption key")
	require.Equal(t, "AES256-GCM", msg.EncryptionAlgorithm, "unexpected encryption algorithm")
	require.NotEmpty(t, msg.Hmac, "message is missing HMAC digital signature")
	require.NotEmpty(t, msg.HmacSecret, "message is missing HMAC secret")
	require.Equal(t, "HMAC-SHA256", msg.HmacAlgorithm, "unexpected signature algorithm")
	require.Nil(t, msg.Error, "unexpected error on message")
	require.NotEmpty(t, msg.Timestamp, "message is missing timestamp")
	require.True(t, msg.Sealed, "message is not marked as sealed")
	require.NotEmpty(t, msg.PublicKeySignature, "message is missing public key signature")
	require.Equal(t, "SHA256:QhEspinUU51gK0dQGqLa56BA6xyRy5/7sN5/6GlaLZw", msg.PublicKeySignature, "unexpected public key signature")
}

// Test the handling of a secure envelope fixture through to creating a response.
func TestRecvEnvelopeWorkflow(t *testing.T) {
	msg, err := loadEnvelopeFixture("testdata/sealed_envelope.json")
	require.NoError(t, err, "could not load envelope")

	key, err := loadPrivateKey("testdata/sealing_key.pem")
	require.NoError(t, err, "could not load sealing key")

	// Wrap the envelope ensuring it's in the sealed state
	senv, err := envelope.Wrap(msg, envelope.WithRSAPrivateKey(key))
	require.NoError(t, err, "could not wrap the envelope")
	require.NoError(t, senv.ValidateMessage(), "secure envelope fixture is invalid")
	require.Equal(t, envelope.Sealed, senv.State(), "expected sealed state not %q", senv.State())

	// Unseal the envelope
	eenv, reject, err := senv.Unseal()
	require.NoError(t, err, "could not unseal the envelope")
	require.Nil(t, reject, "a rejection error was unexpectedly returned")
	require.NotSame(t, senv, eenv, "Unseal should return a clone of the original envelope")
	require.Equal(t, envelope.Unsealed, eenv.State(), "expected unsealed state not %q", eenv.State())

	// Decrypt the envelope
	env, reject, err := eenv.Decrypt()
	require.NoError(t, err, "could not decrypt envelope")
	require.Nil(t, reject, "a rejection error was unexpectedly returned")
	require.NotSame(t, eenv, env, "Decrypt should return a clone of the original envelope")
	require.Equal(t, envelope.Clear, env.State(), "expected clear state not %q", eenv.State())
	require.NotNil(t, env.Crypto(), "decrypted envelopes should maintain crytpo context")
	require.NotNil(t, env.Sealer(), "decrypted envelopes should maintain sealer context")

	// Get the payload from the envelope
	payload, err := env.Payload()
	require.NoError(t, err, "could not fetch decrypted payload from envelope")
	require.NotNil(t, payload, "nil payload unexpectedly returned")

	// Load the payload fixture for verification
	expectedPayload, err := loadPayloadFixture("testdata/pending_payload.json")
	require.NoError(t, err, "could not load payload fixture")
	require.True(t, proto.Equal(payload, expectedPayload), "decrypted payload did not match payload fixture, did fixture change?")

	// Update the payload with received at and reseal the envelope
	// TODO: does this modify the payload of the original message?
	payload.ReceivedAt = time.Now().Format(time.RFC3339)

	oenv, err := envelope.New(payload, envelope.FromEnvelope(env), envelope.WithRSAPublicKey(&key.PublicKey))
	require.NoError(t, err, "could not create new envelope from original envelope")

	eoenv, reject, err := oenv.Encrypt()
	require.NoError(t, err, "could not encrypt envelope")
	require.Nil(t, reject, "a rejection error was unexpectedly returned")
	require.NotSame(t, oenv, eoenv, "envelope unexpectedly not cloned")

	soenv, reject, err := eoenv.Seal()
	require.NoError(t, err, "could not encrypt envelope")
	require.Nil(t, reject, "a rejection error was unexpectedly returned")
	require.NotSame(t, eoenv, soenv, "envelope unexpectedly not cloned")

	out := soenv.Proto()

	// NOTE: cannot use proto.Equal since the timestamp at least will have changed
	require.Equal(t, msg.Id, out.Id, "mismatched envelope ID")
	require.NotEmpty(t, out.Payload, "missing updated, encrypted payload")
	require.NotEmpty(t, out.EncryptionKey, "sealed envelope encryption key missing")
	require.Equal(t, msg.EncryptionAlgorithm, out.EncryptionAlgorithm, "mismatched envelope encryption algorithm")
	require.NotEmpty(t, out.Hmac, "missing updated HMAC signature")
	require.NotEmpty(t, out.HmacSecret, "sealed envelope hmac secret missing")
	require.Equal(t, msg.HmacAlgorithm, out.HmacAlgorithm, "mismatched envelope HMAC algorithm")
	require.Equal(t, msg.Error, out.Error, "unexpected error on envelopes")
	require.NotEmpty(t, out.Timestamp, "no timestamp on outgoing envelope")
	require.True(t, out.Sealed, "out is not marked as sealed")
	require.Equal(t, msg.PublicKeySignature, out.PublicKeySignature, "public key signature mismatch")
}

func TestOneLiners(t *testing.T) {
	payload, err := loadPayloadFixture("testdata/pending_payload.json")
	require.NoError(t, err, "could not load pending payload")

	key, err := loadPrivateKey("testdata/sealing_key.pem")
	require.NoError(t, err, "could not load sealing key")

	// Create an envelope from the payload and the key
	msg, reject, err := envelope.Seal(payload, envelope.WithRSAPublicKey(&key.PublicKey))
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
	decryptedPayload, reject, err := envelope.Open(in, envelope.WithRSAPrivateKey(key))
	require.NoError(t, err, "could not open envelope")
	require.Nil(t, reject, "unexpected rejection error")
	require.True(t, proto.Equal(payload, decryptedPayload), "payloads do not match")
}

func TestEnvelopeAccessors(t *testing.T) {
	// Actual value for timestamp testing
	ats := time.Now()

	// Create a secure envelope with an error
	in := &api.SecureEnvelope{
		Id:        uuid.NewString(),
		Error:     &api.Error{Code: api.ComplianceCheckFail, Message: "afraid of the dark"},
		Timestamp: ats.Format(time.RFC3339Nano),
	}

	env, err := envelope.Wrap(in)
	require.NoError(t, err, "could not wrap envelope")

	require.Equal(t, in.Id, env.ID(), "did not return correct envelope ID")
	require.Equal(t, in, env.Proto(), "proto did not return the embedded envelope")
	require.Equal(t, in.Error, env.Error(), "did not return the embedded error")
	require.Nil(t, env.Crypto(), "crypto should be nil for an error-only envelope")
	require.Nil(t, env.Sealer(), "seal should be nil for an error-only envelope")

	payload, err := env.Payload()
	require.EqualError(t, err, `envelope is in state "error": payload may be invalid`)
	require.Nil(t, payload, "payload should be nil for an error-only envelope")

	ts, err := env.Timestamp()
	require.NoError(t, err, "should have been able to parse RFC3339Nano timestamp")
	require.True(t, ts.Equal(ats), "should have returned now")

	// Test parsing RFC3339 timestamp
	in.Timestamp = ats.Format(time.RFC3339)
	ts, err = env.Timestamp()
	require.NoError(t, err, "should have been able to parse RFC3339 timestamp")
	require.Less(t, ats.Sub(ts), 1*time.Second, "should have returned now without nanosecond precision")

	// Test parsing empty string timestamp
	in.Timestamp = ""
	_, err = env.Timestamp()
	require.EqualError(t, err, "trisa rejection [BAD_REQUEST]: missing ordering timestamp on secure envelope")

	// Test parsing a bad timestamp string
	in.Timestamp = "2022-15-45T38:32:12Z"
	_, err = env.Timestamp()
	require.EqualError(t, err, "trisa rejection [BAD_REQUEST]: could not parse ordering timestamp on secure envelope as RFC3339 timestamp")

	// Create an envelope with a payload
	payload, err = loadPayloadFixture("testdata/payload.json")
	require.NoError(t, err, "could not load payload fixture")

	// Create a new envelope with complete options
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err, "could not generate RSA key")
	env, err = envelope.New(payload, envelope.WithEnvelopeID(in.Id), envelope.WithTimestamp(ts), envelope.WithAESGCM(nil, nil), envelope.WithRSAPublicKey(&key.PublicKey))
	require.NoError(t, err, "could not create new envelope from payload")

	require.Equal(t, in.Id, env.ID(), "ID did not return correct envelope ID")
	require.NotNil(t, env.Proto(), "proto did not return a new secure envelope")
	require.Nil(t, env.Error(), "expected no error to be on envelope")
	require.NotNil(t, env.Crypto(), "crypto should not be nil")
	require.NotNil(t, env.Sealer(), "seal should not be nil")

	actualPayload, err := env.Payload()
	require.NoError(t, err, "error should have been returned")
	require.Equal(t, payload, actualPayload, "payload should match the one instantiated")

	actualTS, err := env.Timestamp()
	require.NoError(t, err, "could not fetch timestamp")
	require.True(t, ts.Equal(actualTS), "timestamp did not match expected timestamp")
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

const (
	expectedEnvelopeId = "2b3b4c95-0a78-4f2a-a9fa-041970f97144"
)

var (
	loadpb = protojson.UnmarshalOptions{
		AllowPartial:   false,
		DiscardUnknown: false,
	}
	dumppb = protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}
)

// Helper method to load a secure envelope fixture, generating the fixtures from the
// payloads if they have not yet been generated.
func loadEnvelopeFixture(path string) (msg *api.SecureEnvelope, err error) {
	msg = &api.SecureEnvelope{}
	if err = loadFixture(path, msg, true); err != nil {
		return nil, err
	}
	return msg, nil
}

// Helper method to load a payload fixture, generating it if it hasn't been yet
func loadPayloadFixture(path string) (payload *api.Payload, err error) {
	payload = &api.Payload{}
	if err = loadFixture(path, payload, true); err != nil {
		return nil, err
	}
	return payload, nil
}

// Helper method to load a fixture from JSON
func loadFixture(path string, m proto.Message, check bool) (err error) {
	// Check if the path exists, if it doesn't attempt to generate the fixture.
	if check {
		if _, err = os.Stat(path); os.IsNotExist(err) {
			if err = generateFixtures(); err != nil {
				return err
			}
		}
	}

	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return err
	}

	if err = loadpb.Unmarshal(data, m); err != nil {
		return err
	}
	return nil
}

// Helper method to generate secure envelopes from the payload fixtures
func generateFixtures() (err error) {
	// Load the components of the various payloads that will be created
	var (
		payload        *api.Payload
		pendingPayload *api.Payload
	)

	identity := &ivms101.IdentityPayload{}
	if err = loadFixture("testdata/payload/identity.json", identity, false); err != nil {
		return fmt.Errorf("could not unmarshal identity payload: %v", err)
	}

	pending := &generic.Pending{}
	if err = loadFixture("testdata/payload/pending.json", pending, false); err != nil {
		return fmt.Errorf("could not read pending payload: %v", err)
	}

	transaction := &generic.Transaction{}
	if err = loadFixture("testdata/payload/transaction.json", transaction, false); err != nil {
		return fmt.Errorf("could not read transaction payload: %v", err)
	}

	payload = &api.Payload{
		SentAt:     "2022-01-27T08:21:43Z",
		ReceivedAt: "2022-01-30T16:28:39Z",
	}
	if payload.Identity, err = anypb.New(identity); err != nil {
		return fmt.Errorf("could not create identity payload: %v", err)
	}
	if payload.Transaction, err = anypb.New(transaction); err != nil {
		return fmt.Errorf("could not create transaction payload: %v", err)
	}

	pendingPayload = &api.Payload{
		Identity: payload.Identity,
		SentAt:   payload.SentAt,
	}
	if pendingPayload.Transaction, err = anypb.New(pending); err != nil {
		return fmt.Errorf("could not create pending payload: %v", err)
	}

	// Serialize the payloads
	if err = dumpFixture("testdata/payload.json", payload); err != nil {
		return fmt.Errorf("could not marshal payload: %v", err)
	}

	if err = dumpFixture("testdata/pending_payload.json", pendingPayload); err != nil {
		return fmt.Errorf("could not marshal pending payload: %v", err)
	}

	// Create error-only envelope
	env := &api.SecureEnvelope{
		Id:        expectedEnvelopeId,
		Timestamp: "2022-01-27T08:21:43Z",
		Error: &api.Error{
			Code:    api.Error_COMPLIANCE_CHECK_FAIL,
			Message: "specified account has been frozen temporarily",
		},
	}

	if err = dumpFixture("testdata/error_envelope.json", env); err != nil {
		return fmt.Errorf("could not marshal error only envelope: %v", err)
	}

	// Create unsealed envelope
	var handler *envelope.Envelope
	if handler, err = envelope.New(payload); err != nil {
		return err
	}

	if handler, _, err = handler.Encrypt(); err != nil {
		return err
	}

	if err = dumpFixture("testdata/unsealed_envelope.json", handler.Proto()); err != nil {
		return fmt.Errorf("could not marshal unsealed envelope: %v", err)
	}

	// Create RSA keys for sealed secure envelope fixtures
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("could not generate RSA key fixture")
	}
	if err = dumpPrivateKey("testdata/sealing_key.pem", key); err != nil {
		return err
	}

	if env, _, err = envelope.Seal(pendingPayload, envelope.WithRSAPublicKey(&key.PublicKey)); err != nil {
		return err
	}
	if err = dumpFixture("testdata/sealed_envelope.json", env); err != nil {
		return fmt.Errorf("could not marshal sealed envelope: %v", err)
	}
	return nil
}

func dumpFixture(path string, m proto.Message) (err error) {
	var data []byte
	if data, err = dumppb.Marshal(m); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func dumpPrivateKey(path string, key *rsa.PrivateKey) (err error) {
	var data []byte
	if data, err = x509.MarshalPKCS8PrivateKey(key); err != nil {
		return err
	}

	block := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: data,
	})

	return os.WriteFile(path, block, 0600)
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
