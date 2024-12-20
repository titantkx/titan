package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// RewardKeyPrefix is the prefix to retrieve all Reward
	RewardKeyPrefix = "Reward/value/"
)

// RewardKey returns the store key to retrieve a Reward from the index fields
func RewardKey(farmer string) []byte {
	var key []byte

	farmerBytes := []byte(farmer)
	key = append(key, farmerBytes...)
	key = append(key, []byte("/")...)

	return key
}
