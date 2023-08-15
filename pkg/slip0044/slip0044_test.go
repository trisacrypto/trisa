package slip0044_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/slip0044"
)

func TestParseCoinType(t *testing.T) {
	// Test that commonly used coin types are parsed correctly
	testCases := []struct {
		coin          string
		coinType      uint32
		pathComponent uint32
	}{
		{"BTC", 0, 0x80000000},
		{"Bitcoin", 0, 0x80000000},
		{"bitcoin", 0, 0x80000000},
		{"LTC", 2, 0x80000002},
		{"Litecoin", 2, 0x80000002},
		{"litecoin", 2, 0x80000002},
		{"DOGE", 3, 0x80000003},
		{"Dogecoin", 3, 0x80000003},
		{"dogecoin", 3, 0x80000003},
		{"DASH", 5, 0x80000005},
		{"Dash", 5, 0x80000005},
		{"dash", 5, 0x80000005},
		{"ETH", 60, 0x8000003c},
		{"Ether", 60, 0x8000003c},
		{"ether", 60, 0x8000003c},
		{"ETC", 61, 0x8000003d},
		{"Ether Classic", 61, 0x8000003d},
	}

	for _, test := range testCases {
		coinType, err := slip0044.ParseCoinType(test.coin)
		require.NoError(t, err, "failed to parse coin type for test case %s", test.coin)
		require.Equal(t, test.coinType, uint32(coinType), "coin type mismatch")
		require.Equal(t, test.pathComponent, coinType.PathComponent(), "path component mismatch")
	}

	_, err := slip0044.ParseCoinType("unknown")
	require.ErrorIs(t, slip0044.ErrUnknownCoin, err, "expected unknown coin error")
}
