package ivms101_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestIdentityPayloadSerialization(t *testing.T) {
	in, err := os.ReadFile("testdata/identity_payload.json")
	require.NoError(t, err, "unable to read identity payload JSON fixture")

	// Unmarshal IVMS101 Identity Payload
	identity := &ivms101.IdentityPayload{}
	err = json.Unmarshal(in, identity)
	require.NoError(t, err, "unable to unmarshal identity payload JSON fixture")

	// Marshal IVMS101 Identity Payload
	out, err := json.Marshal(identity)
	require.NoError(t, err, "could not marshal identity payload to JSON")

	fmt.Println(string(out))
	diffOpts := jsondiff.DefaultConsoleOptions()
	res, _ := jsondiff.Compare(in, out, &diffOpts)
	require.Equal(t, res, jsondiff.FullMatch, "marshalled json differs from original")
}

func TestIdentityPayloadSerializationFromPB(t *testing.T) {
	// This test loads the identity payload from a protojson serialized fixture then
	// marshals and unmarshals the data as an inverse to TestIdentityPayloadSerialization
	// NOTE: we could use a raw protocol buffer here, but protojson makes it easier to
	// read and manage the JSON fixture for the future.
	pbdata, err := os.ReadFile("testdata/identity_payload.pb.json")
	require.NoError(t, err, "unable to read identity payload PB JSON fixture")

	jsonpb := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	// Load Protocol Buffer Fixture
	identity := &ivms101.IdentityPayload{}
	err = jsonpb.Unmarshal(pbdata, identity)
	require.NoError(t, err, "unable to unmarshal identity payload PB JSON fixture")

	// Marshal IVMS101 Identity Payload
	data, err := json.Marshal(identity)
	require.NoError(t, err, "could not marshal identity payload to JSON")

	// Unmarshal IVMS101 Identity Payload
	odentity := &ivms101.IdentityPayload{}
	err = json.Unmarshal(data, odentity)
	require.NoError(t, err, "unable to unmarshal identity payload from JSON")

	require.True(t, proto.Equal(identity, odentity), "serialized identity payload does not match original")
}
