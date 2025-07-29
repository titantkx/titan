package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// Erc20NativeKeyPrefix is the prefix to retrieve all Erc20Native
	Erc20NativeKeyPrefix = "Erc20Native/value/"
)

// Erc20NativeKey returns the store key to retrieve a Erc20Native from the index fields
func Erc20NativeKey(
	tokenDenom string,
) []byte {
	var key []byte

	tokenDenomBytes := []byte(tokenDenom)
	key = append(key, tokenDenomBytes...)
	key = append(key, []byte("/")...)

	return key
}
