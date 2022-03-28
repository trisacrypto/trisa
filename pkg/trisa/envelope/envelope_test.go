package envelope_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
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

func Example_create() {
	payload, _ := loadPayloadFixture("testdata/payload.json")
	fmt.Println(payload)
}

// Test the creation of an envelope from scratch, moving it through each state.
func TestSendEnvelopeWorkflow(t *testing.T) {
	payload, err := loadPayloadFixture("testdata/payload.json")
	require.NoError(t, err, "could not load payload")

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err, "could not generate sealing key")

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
	require.EqualError(t, err, `envelope is in state "unsealed-error": payload may be invalid`)
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
	if data, err = ioutil.ReadFile(path); err != nil {
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
	return nil
}

func dumpFixture(path string, m proto.Message) (err error) {
	var data []byte
	if data, err = dumppb.Marshal(m); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}
