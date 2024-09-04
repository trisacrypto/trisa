package iso3166_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/iso3166"
)

func TestValidate(t *testing.T) {
	t.Run("All", func(t *testing.T) {
		tests := []struct {
			alpha     string
			err       error
			errString string
		}{
			{
				"GB", nil, "",
			},
			{
				"TUN", nil, "",
			},
			{
				"TUND", iso3166.ErrInvalidAlpha, "",
			},
			{
				"T", iso3166.ErrInvalidAlpha, "",
			},
			{
				"gb", iso3166.ErrInvalidAlpha, "",
			},
			{
				"tun", iso3166.ErrInvalidAlpha, "",
			},
			{
				"tu1", iso3166.ErrInvalidAlpha, "",
			},
			{
				"FuN", iso3166.ErrInvalidAlpha, "",
			},
			{
				"12", iso3166.ErrInvalidAlpha, "",
			},
			{
				"2B", iso3166.ErrInvalidAlpha, "",
			},
			{
				"$A", iso3166.ErrInvalidAlpha, "",
			},
			{
				"A$A", iso3166.ErrInvalidAlpha, "",
			},
			{
				"XX", nil, "iso3166: \"XX\" is not a recognized country code",
			},
			{
				"XYZ", nil, "iso3166: \"XYZ\" is not a recognized country code",
			},
		}

		for i, tc := range tests {
			err := iso3166.Validate(tc.alpha)
			switch {
			case tc.err != nil:
				require.ErrorIs(t, err, tc.err, "test case %d failed", i)
			case tc.errString != "":
				require.EqualError(t, err, tc.errString, "test case %d failed", i)
			default:
				require.NoError(t, err, "test case %d failed", i)
			}
		}
	})

	t.Run("Alpha2", func(t *testing.T) {
		tests := []struct {
			alpha     string
			err       error
			errString string
		}{
			{
				"GB", nil, "",
			},
			{
				"TUN", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"TUND", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"T", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"gb", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"tun", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"tu1", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"FuN", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"12", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"2B", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"$A", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"A$A", iso3166.ErrInvalidAlpha2, "",
			},
			{
				"XX", nil, "iso3166: \"XX\" is not a recognized country code",
			},
			{
				"XYZ", iso3166.ErrInvalidAlpha2, "",
			},
		}

		for i, tc := range tests {
			err := iso3166.ValidateAlpha2(tc.alpha)
			switch {
			case tc.err != nil:
				require.ErrorIs(t, err, tc.err, "test case %d failed", i)
			case tc.errString != "":
				require.EqualError(t, err, tc.errString, "test case %d failed", i)
			default:
				require.NoError(t, err, "test case %d failed", i)
			}
		}
	})

	t.Run("Alpha3", func(t *testing.T) {
		tests := []struct {
			alpha     string
			err       error
			errString string
		}{
			{
				"GB", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"TUN", nil, "",
			},
			{
				"TUND", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"T", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"gb", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"tun", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"tu1", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"FuN", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"12", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"2B", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"$A", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"A$A", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"XX", iso3166.ErrInvalidAlpha3, "",
			},
			{
				"XYZ", nil, "iso3166: \"XYZ\" is not a recognized country code",
			},
		}

		for i, tc := range tests {
			err := iso3166.ValidateAlpha3(tc.alpha)
			switch {
			case tc.err != nil:
				require.ErrorIs(t, err, tc.err, "test case %d failed", i)
			case tc.errString != "":
				require.EqualError(t, err, tc.errString, "test case %d failed", i)
			default:
				require.NoError(t, err, "test case %d failed", i)
			}
		}
	})
}

func TestValidCountries(t *testing.T) {
	for _, code := range iso3166.List() {
		require.NoError(t, iso3166.Validate(code.Alpha2), "alpha validation failed for %q", code.Alpha2)
		require.NoError(t, iso3166.Validate(code.Alpha3), "alpha validation failed for %q", code.Alpha3)
		require.NoError(t, iso3166.ValidateAlpha2(code.Alpha2), "alpha2 validation failed for %q", code.Alpha2)
		require.NoError(t, iso3166.ValidateAlpha3(code.Alpha3), "alpha3 validation failed for %q", code.Alpha3)
	}
}
