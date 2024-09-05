package gleif_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/gleif"
)

func TestValidAuthorities(t *testing.T) {
	f, err := os.Open("testdata/registrationAuthorities.json")
	require.NoError(t, err)
	defer f.Close()

	var authorities gleif.RegistrationAuthorities
	err = json.NewDecoder(f).Decode(&authorities)
	require.NoError(t, err)
	require.Len(t, authorities, 1037)

	for _, ra := range authorities {
		require.NoError(t, gleif.Validate(ra.Option))
	}
}

func TestInvalidAuthorities(t *testing.T) {
	tests := []struct {
		ra  string
		err error
	}{
		{
			"VA001231",
			gleif.ErrIncorrectFormat,
		},
		{
			"ra001231",
			gleif.ErrIncorrectFormat,
		},
		{
			"RA14",
			gleif.ErrIncorrectFormat,
		},
		{
			"RA101005",
			gleif.ErrNotFound,
		},
	}

	for i, tc := range tests {
		require.ErrorIs(t, gleif.Validate(tc.ra), tc.err, "test case %d failed", i)
	}
}
