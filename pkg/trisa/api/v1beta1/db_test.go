package api_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"google.golang.org/protobuf/proto"
)

func TestSecureEnvelopeDatabase(t *testing.T) {
	ip := &api.SecureEnvelope{
		Id:                  uuid.NewString(),
		Payload:             []byte("payload"),
		EncryptionKey:       []byte("supersecretkey"),
		EncryptionAlgorithm: "NO-CIPHER",
		Hmac:                []byte("digitalsignature"),
		HmacSecret:          []byte("whispersecret"),
		HmacAlgorithm:       "NO-HMAC",
		Timestamp:           "2024-04-04T18:38:58-05:00",
		Sealed:              false,
		PublicKeySignature:  "notakey",
	}

	value, err := ip.Value()
	require.NoError(t, err, "could not fetch value for secure envelope")

	cp := &api.SecureEnvelope{}
	err = cp.Scan(value)
	require.NoError(t, err, "could not scan value into secure envelope")
	require.True(t, proto.Equal(ip, cp), "loaded value not equal to copied value")
}
