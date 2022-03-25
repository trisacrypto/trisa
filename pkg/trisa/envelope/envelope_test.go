package envelope_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/envelope"
)

func TestEnvelopeAccessors(t *testing.T) {
	// Actual value for timestamp testing
	ats := time.Now()

	// Create a secure envelope with an error
	in := &api.SecureEnvelope{
		Id:        uuid.NewString(),
		Error:     &api.Error{Code: api.ComplianceCheckFail, Message: "afraid of the dark"},
		Timestamp: ats.Format(time.RFC3339Nano),
	}

	env, err := envelope.Open(in)
	require.NoError(t, err, "could not open envelope")

	require.Equal(t, in, env.Proto(), "proto did not return the embedded envelope")
	require.Equal(t, in.Error, env.Error(), "proto did not return the embedded error")
	require.Nil(t, env.Payload(), "payload should be nil for an error-only envelope")
	require.Nil(t, env.Crypto(), "crypto should be nil for an error-only envelope")
	require.Nil(t, env.Seal(), "seal should be nil for an error-only envelope")

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
}
