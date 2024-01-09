package types

const (
	// ModuleName defines the module name
	ModuleName = "validatorreward"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_validatorreward"

	// ValidatorRewardCollectorName the root string for the validator reward collector account address
	ValidatorRewardCollectorName = "validator_reward_collector"
)

var ParamsKey = []byte("Params")

func KeyPrefix(p string) []byte {
	return []byte(p)
}
