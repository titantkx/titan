package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// FarmKeyPrefix is the prefix to retrieve all Farm
	FarmKeyPrefix = "Farm/value/"
)

// FarmKey returns the store key to retrieve a Farm from the index fields
func FarmKey(token string) []byte {
	var key []byte

	tokenBytes := []byte(token)
	key = append(key, tokenBytes...)
	key = append(key, []byte("/")...)

	return key
}
