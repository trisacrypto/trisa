package slip0044

import "strings"

// Return the path component that corresponds to this coin type.
func (t CoinType) PathComponent() uint32 {
	return 0x80000000 | uint32(t)
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
