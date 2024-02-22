package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// MintingInfoKeyPrefix is the prefix to retrieve all MintingInfo
	MintingInfoKeyPrefix = "MintingInfo/value/"
)

// MintingInfoKey returns the store key to retrieve a MintingInfo from the index fields
func MintingInfoKey(
	classId string,
) []byte {
	var key []byte

	classIdBytes := []byte(classId)
	key = append(key, classIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
