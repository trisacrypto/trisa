package trp_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

func TestTravelAddress(t *testing.T) {
	testCases := []struct {
		address  string
		expected string
		err      error
	}{
		{
			"lnurl1dp68gurn8ghj7cn9dejkv6trd9shy7fwvdhk6tm5wfcr7arpvu7hgunpwejkcun4d3jkjmn3w45hy7gmsy37e",
			"https://beneficiary.com/trp?tag=travelruleinquiry",
			nil,
		},
		{
			"LNURL1DP68GURN8GHJ7MMSV4H8VCTNWQH8GETNWSKKUET59E5K7TE3XGEN7ARPVU7KJMN3W45HY7GF5KZ53",
			"https://openvasp.test-net.io/123?tag=inquiry",
			nil,
		},
		{
			"taGw1e4cjuujwyHBV51aspbLBUbcDhebW7ss8iF6dLEn19WfyQZt6HaUTStSE2YcadkFG",
			"https://beneficiary.com/x/12345?color=red&foo=bar&t=i",
			nil,
		},
		{
			"https://beneficiary.com/x/12345?t=i",
			"https://beneficiary.com/x/12345?t=i",
			nil,
		},
		{
			"foo", "",
			ErrUnknownTravelAddress,
		},
	}

	for i, tc := range testCases {
		info := &Info{Address: tc.address, APIVersion: "3.1.0", RequestIdentifier: "704c548a-70af-480c-af83-6fb7803df85c"}
		actual, err := info.URL()
		if tc.err != nil {
			require.Error(t, err, "expected error on test case %d", i)
			require.ErrorIs(t, err, tc.err, "unexpected error on test case %d", i)
		} else {
			require.NoError(t, err, "expected no error on test case %d", i)
			require.Equal(t, tc.expected, actual.String(), "unexpected mismatch on test case %d", i)
		}
	}
}
