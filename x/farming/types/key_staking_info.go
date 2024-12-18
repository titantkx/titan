package types

import (
	"encoding/binary"
	"fmt"
)

var _ binary.ByteOrder

const (
	// StakingInfoKeyPrefix is the prefix to retrieve all StakingInfo
	StakingInfoKeyPrefix = "StakingInfo/value/"
)

// StakingInfoKey returns the store key to retrieve a StakingInfo from the index fields
func StakingInfoKey(token string, staker string) []byte {
	if staker != "" {
		return []byte(fmt.Sprintf("%s/%s/", token, staker))
	}
	return []byte(token + "/")
}
