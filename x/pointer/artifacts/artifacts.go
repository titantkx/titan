package artifacts

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/titantkx/titan/x/pointer/artifacts/erc20native"
)

func GetParsedABI(typ string) *abi.ABI {
	switch typ {
	case "erc20native":
		return erc20native.GetParsedABI()
	default:
		panic(fmt.Sprintf("unknown artifact type %s", typ))
	}
}

func GetBin(typ string) []byte {
	switch typ {
	case "erc20native":
		return erc20native.GetBin()
	default:
		panic(fmt.Sprintf("unknown artifact type %s", typ))
	}
}
