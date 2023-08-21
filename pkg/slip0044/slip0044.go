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
	return symbolsByType[t]
}

// Parse a coin from its symbol or name into the corresponding coin type.
func ParseCoinType(coin string) (CoinType, error) {
	// First try to parse as a symbol.
	if coinType, ok := typesBySymbol[strings.ToUpper(coin)]; ok {
		return coinType, nil
	}

	// Otherwise try to parse as a name.
	if coinType, ok := typesByName[strings.ToLower(coin)]; ok {
		return coinType, nil
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
