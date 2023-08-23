package lnurl_test

import (
	"strings"
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
		{
			"lnurl1dp68gurn8ghj7cn9dejkv6trd9shy7fwvdhk6tm5wfcr7arpvu7hgunpwejkcun4d3jkjmn3w45hy7gmsy37e",
			"https://beneficiary.com/trp?tag=travelruleinquiry",
		},
	}

	for i, tc := range testCases {
		actual, err := lnurl.Decode(tc.lnurl)
		require.NoError(t, err, "test case %d decode failed", i)
		require.Equal(t, tc.wburl, actual, "test case %d decode equality failed", i)

		actual, err = lnurl.Encode(tc.wburl)
		require.NoError(t, err, "test case %d encode failed", i)
		require.Equal(t, strings.ToUpper(tc.lnurl), actual, "test case %d encode equality failed", i)
	}
}

func TestLNURLErrors(t *testing.T) {
	testCases := []struct {
		input string
		err   error
	}{
		{
			"https://DP68GURN8GHJ7MMSV4H8VCTNWQH8GETNWSKKUET59E5K7TE3XGEN7ARPVU7KJMN3W45HY7GF5KZ53",
			lnurl.ErrUnhandledScheme,
		},
		{
			"lnurl1split1checkupstagehandshakeupstreamerranterredcaperred2y9e2w",
			lnurl.ErrInvalidChecksum{"6gr7g4", "2y9e2w"},
		},
		{
			"lnurl1s lit1checkupstagehandshakeupstreamerranterredcaperredp8hs2p",
			lnurl.ErrInvalidCharacter(' '),
		},
		{
			"lnurl1spl\x7Ft1checkupstagehandshakeupstreamerranterredcaperred2y9e3w",
			lnurl.ErrInvalidCharacter(127),
		},
		{
			"lnurl1split1cheo2y9e2w",
			lnurl.ErrNonCharsetChar('o'),
		},
		{
			"lnurl1split1a2y9w",
			lnurl.ErrInvalidSeparatorIndex(1),
		},
	}

	for i, tc := range testCases {
		actual, err := lnurl.Decode(tc.input)
		require.Error(t, err, "test case %d did not error", i)
		require.ErrorIs(t, err, tc.err, "test case %d error did not match", i)
		require.Empty(t, actual, "test case %d returned data", i)
	}
}
