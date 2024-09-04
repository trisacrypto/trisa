package ivms101_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
)

func TestRekeying(t *testing.T) {
	ivms101.AllowRekeying()
	t.Cleanup(ivms101.DisallowRekeying)

	expected, err := os.ReadFile("testdata/identity_payload.json")
	require.NoError(t, err, "could not load testdata/identity_payload.json")

	t.Run("NoModification", func(t *testing.T) {
		f, err := os.Open("testdata/identity_payload.json")
		require.NoError(t, err, "could not load testdata/identity_payload.json")
		defer f.Close()

		payload := &ivms101.IdentityPayload{}
		err = json.NewDecoder(f).Decode(payload)
		require.NoError(t, err, "could not decode payload with rekeying")

		compat, err := json.Marshal(payload)
		require.NoError(t, err, "could not marshal payload")

		require.JSONEq(t, string(expected), string(compat), "json not equal to expected")
	})

	t.Run("Modified", func(t *testing.T) {
		f, err := os.Open("testdata/identity_payload_alt.json")
		require.NoError(t, err, "could not load testdata/identity_payload_alt.json")
		defer f.Close()

		payload := &ivms101.IdentityPayload{}
		err = json.NewDecoder(f).Decode(payload)
		require.NoError(t, err, "could not decode payload with rekeying")

		compat, err := json.Marshal(payload)
		require.NoError(t, err, "could not marshal payload")

		require.JSONEq(t, string(expected), string(compat), "json not equal to expected")
	})

	t.Run("Protobuf", func(t *testing.T) {
		f, err := os.Open("testdata/identity_payload.pb.json")
		require.NoError(t, err, "could not load testdata/identity_payload.pb.json")
		defer f.Close()

		payload := &ivms101.IdentityPayload{}
		err = json.NewDecoder(f).Decode(payload)
		require.NoError(t, err, "could not decode payload with rekeying")
	})
}

func BenchmarkRekeying(b *testing.B) {
	// Rekey1 uses a map[string]interface{} for decoding and rekeying.
	rekey1 := func(b *testing.B, data []byte) {
		var middle map[string]interface{}
		err := json.Unmarshal(data, &middle)
		assert.NoError(b, err)

		cmpt, err := json.Marshal(middle)
		assert.NoError(b, err)
		assert.Len(b, cmpt, 2733)
	}

	// Rekey2 uses a map[string]json.RawMessage for decoding and rekeying.
	// Rekey2 is the clear winner in terms of performance.
	reykey2 := func(b *testing.B, data []byte) {
		var middle map[string]json.RawMessage
		err := json.Unmarshal(data, &middle)
		assert.NoError(b, err)

		cmpt, err := json.Marshal(middle)
		assert.NoError(b, err)
		assert.Len(b, cmpt, 2733)
	}

	data, err := os.ReadFile("testdata/identity_payload.json")
	assert.NoError(b, err)

	b.Run("Rekey1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rekey1(b, data)
		}
	})

	b.Run("Rekey2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reykey2(b, data)
		}
	})
}
