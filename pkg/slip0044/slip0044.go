package slip0044

import (
	"encoding/json"
	"strings"
)

// Return the path component that corresponds to this coin type.
func (t CoinType) PathComponent() uint32 {
	return 0x80000000 | uint32(t)
}

// Return the coin symbol, which is preferred over the native String() method.
func (t CoinType) Symbol() string {
	// Because the symbol is the first entry in the enum map, it is marked as the name
	// in CoinType_name -- to ensure this happens, the slip0044.jsonl file must be
	// ordered according to the correct symbol and the generate.py file in the proto
	// folder must output the symbol as the first entry. This also depends on the way
	// that protobuf generates enum.pb.go marking any duplicate values after the first
	// value in the file.
	return CoinType_name[int32(t)]
}

// Parse a coin from its symbol or name into the corresponding coin type.
func ParseCoinType(coin string) (CoinType, error) {
	// Normalize the input to an enum-like value.
	coin = strings.Replace(strings.ToUpper(strings.TrimSpace(coin)), " ", "_", -1)

	if ct, ok := CoinType_value[coin]; ok {
		return CoinType(ct), nil
	}

	return 0, ErrUnknownCoin
}

func (t CoinType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Symbol())
}

func (t *CoinType) UnmarshalJSON(data []byte) (err error) {
	var symbol string
	if err = json.Unmarshal(data, &symbol); err != nil {
		return err
	}

	if *t, err = ParseCoinType(symbol); err != nil {
		return err
	}
	return nil
}
