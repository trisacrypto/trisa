package ivms101_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestIdentityPayloadSerialization(t *testing.T) {
	t.Skip("not yet implemented")
	in, err := ioutil.ReadFile("testdata/identity_payload.json")
	require.NoError(t, err, "unable to read identity payload JSON fixture")

	// Unmarshal IVMS101 Identity Payload
	identity := &ivms101.IdentityPayload{}
	err = json.Unmarshal(in, identity)
	require.NoError(t, err, "unable to unmarshal identity payload JSON fixture")

	// Marshal IVMS101 Identity Payload
	out, err := json.Marshal(identity)
	require.NoError(t, err, "could not marshal identity payload to JSON")

	require.Equal(t, in, out, "marshaled and unmarshaled JSON does not match")
}

func TestIdentityPayloadSerializationFromPB(t *testing.T) {
	t.Skip("not yet implemented")
	// This test loads the identity payload from a protojson serialized fixture then
	// marshals and unmarshals the data as an inverse to TestIdentityPayloadSerialization
	// NOTE: we could use a raw protocol buffer here, but protojson makes it easier to
	// read and manage the JSON fixture for the future.
	pbdata, err := ioutil.ReadFile("testdata/identity_payload.pb.json")
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

	fmt.Println(string(data))

	// Unmarshal IVMS101 Identity Payload
	odentity := &ivms101.IdentityPayload{}
	err = json.Unmarshal(data, odentity)
	require.NoError(t, err, "unable to unmarshal identity payload from JSON")

	require.True(t, proto.Equal(identity, odentity), "serialized identity payload does not match original")
}
