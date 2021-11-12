package iso3166_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/iso3166"
)

func TestFind(t *testing.T) {
	code, err := iso3166.Find("United States")
	require.NoError(t, err)
	require.Equal(t, "USA", code.Alpha3)

	code, err = iso3166.Find("GB")
	require.NoError(t, err)
	require.Equal(t, "United Kingdom", code.Country)

	code, err = iso3166.Find("gb")
	require.NoError(t, err)
	require.Equal(t, "United Kingdom", code.Country)

	code, err = iso3166.Find("BRA")
	require.NoError(t, err)
	require.Equal(t, "Brazil", code.Country)

	code, err = iso3166.Find("bra")
	require.NoError(t, err)
	require.Equal(t, "Brazil", code.Country)

	code, err = iso3166.Find("376")
	require.NoError(t, err)
	require.Equal(t, "Israel", code.Country)

	code, err = iso3166.Find("turks and caicos")
	require.NoError(t, err)
	require.Equal(t, "TC", code.Alpha2)

	_, err = iso3166.Find("Foo")
	require.Error(t, err)

	_, err = iso3166.Find("turk")
	require.Contains(t, err.Error(), "ambiguous, multiple countries matched")
}

func TestNormalizedSearch(t *testing.T) {
	for _, code := range iso3166.List() {
		// Unnormalized result should always return
		found, err := iso3166.Find(code.Country)
		require.NoError(t, err)
		require.Equal(t, code.Alpha3, found.Alpha3)

		// Normalized results might return an error
		found, err = iso3166.Find(strings.ToLower(code.Country))
		if err != nil {
			t.Logf("%q has an ambiguous, case-insensitive lookup", code.Country)
			require.Contains(t, err.Error(), "ambiguous, multiple countries matched")
		} else {
			require.Equal(t, code.Alpha3, found.Alpha3)
		}
	}
}

func TestList(t *testing.T) {
	codes := iso3166.List()
	require.Equal(t, 249, len(codes))
}
