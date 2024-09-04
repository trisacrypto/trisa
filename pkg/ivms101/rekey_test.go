package ivms101_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
