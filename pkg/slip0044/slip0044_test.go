package slip0044_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/slip0044"
)

func TestParseCoinType(t *testing.T) {
	// Test that commonly used coin types are parsed correctly
	testCases := []struct {
		coin          string
		symbol        string
		coinType      uint32
		pathComponent uint32
	}{
		{"BTC", "BTC", 0, 0x80000000},
		{"Bitcoin", "BTC", 0, 0x80000000},
		{"bitcoin", "BTC", 0, 0x80000000},
		{"LTC", "LTC", 2, 0x80000002},
		{"Litecoin", "LTC", 2, 0x80000002},
		{"litecoin", "LTC", 2, 0x80000002},
		{"DOGE", "DOGE", 3, 0x80000003},
		{"Dogecoin", "DOGE", 3, 0x80000003},
		{"dogecoin", "DOGE", 3, 0x80000003},
		{"DASH", "DASH", 5, 0x80000005},
		{"Dash", "DASH", 5, 0x80000005},
		{"dash", "DASH", 5, 0x80000005},
		{"ETH", "ETH", 60, 0x8000003c},
		{"Ether", "ETH", 60, 0x8000003c},
		{"ether", "ETH", 60, 0x8000003c},
		{"ETC", "ETC", 61, 0x8000003d},
		{"Ether Classic", "ETC", 61, 0x8000003d},
	}

	for _, test := range testCases {
		coinType, err := slip0044.ParseCoinType(test.coin)
		require.NoError(t, err, "failed to parse coin type for test case %s", test.coin)
		require.Equal(t, test.coinType, uint32(coinType), "coin type mismatch")
		require.Equal(t, test.pathComponent, coinType.PathComponent(), "path component mismatch")
		require.Equal(t, test.symbol, coinType.Symbol(), "symbol mismatch")
	}

	_, err := slip0044.ParseCoinType("unknown")
	require.ErrorIs(t, slip0044.ErrUnknownCoin, err, "expected unknown coin error")
}

func TestJSON(t *testing.T) {
	type Example struct {
		Asset        slip0044.CoinType   `json:"asset"`
		Alternatives []slip0044.CoinType `json:"alternatives"`
	}

	data, err := os.ReadFile("testdata/example.json")
	require.NoError(t, err, "could not load testdata/example.json")

	var example *Example
	err = json.Unmarshal(data, &example)
	require.NoError(t, err, "could not unmarshal example json")

	cmpt, err := json.Marshal(&example)
	require.NoError(t, err, "could not marshal json data")
	require.JSONEq(t, string(data), string(cmpt), "marshaled json data did not match example")

}
