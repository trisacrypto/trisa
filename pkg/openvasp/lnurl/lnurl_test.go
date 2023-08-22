package lnurl_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/openvasp/lnurl"
)

func TestLNURL(t *testing.T) {
	testCases := []struct {
		lnurl string
		wburl string
	}{
		{
			"LNURL1DP68GURN8GHJ7MMSV4H8VCTNWQH8GETNWSKKUET59E5K7TE3XGEN7ARPVU7KJMN3W45HY7GF5KZ53",
			"https://openvasp.test-net.io/123?tag=inquiry",
		},
	}

	for i, tc := range testCases {
		actual, err := lnurl.Decode(tc.lnurl)
		require.NoError(t, err, "test case %d decode failed", i)
		require.Equal(t, tc.wburl, actual, "test case %d docde equality failed", i)

		actual, err = lnurl.Encode(tc.wburl)
		require.NoError(t, err, "test case %d encode failed", i)
		require.Equal(t, tc.lnurl, actual, "test case %d encode equality failed", i)
	}
}
