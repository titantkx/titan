package utils

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "github.com/tokenize-titan/ethermint/cmd/config"
)

const (
	AccountAddressPrefix = "titan"

	MainnetChainID = "titan-18888"

	TestnetChainID = "titan-18889"
	// DisplayDenom defines the denomination displayed to users in client applications.
	DisplayDenom = "tkx"
)

var (
	// BaseDenom defines to the default denomination used in titan (staking, governance, etc.)
	BaseDenom = fmt.Sprintf("a%s", DisplayDenom)
	// utkx
	MicroDenom = fmt.Sprintf("u%s", DisplayDenom)
	// mtkx
	MilliDenom = fmt.Sprintf("m%s", DisplayDenom)
)

func InitSDKConfig() {
	sdk.DefaultBondDenom = BaseDenom
	sdk.DefaultPowerReduction = sdk.NewIntFromUint64(1e18)

	// Set prefixes
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"

	if sdk.GetConfig().GetBech32AccountAddrPrefix() != AccountAddressPrefix {
		// Set and seal config
		config := sdk.GetConfig()
		config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
		config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
		config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)

		// Ethermint config coin type to 60
		cmdcfg.SetBip44CoinType(config)

		config.Seal()
	}
}

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if _, registed := sdk.GetDenomUnit(DisplayDenom); !registed {
		if err := sdk.RegisterDenom(DisplayDenom, sdk.OneDec()); err != nil {
			panic(err)
		}
	}

	if _, registed := sdk.GetDenomUnit(BaseDenom); !registed {
		if err := sdk.RegisterDenom(BaseDenom, sdk.NewDecWithPrec(1, 18)); err != nil {
			panic(err)
		}
	}

	if _, registed := sdk.GetDenomUnit(MicroDenom); !registed {
		if err := sdk.RegisterDenom(MicroDenom, sdk.NewDecWithPrec(1, 6)); err != nil {
			panic(err)
		}
	}

	if _, registed := sdk.GetDenomUnit(MilliDenom); !registed {
		if err := sdk.RegisterDenom(MilliDenom, sdk.NewDecWithPrec(1, 3)); err != nil {
			panic(err)
		}
	}
}