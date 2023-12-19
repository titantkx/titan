package cmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "github.com/tokenize-titan/ethermint/cmd/config"

	"github.com/tokenize-titan/titan/app"
)

const (
	MainnetChainID = "titan-18888"

	TestnetChainID = "titan-18889"
	// DisplayDenom defines the denomination displayed to users in client applications.
	DisplayDenom = "tkx"
	// BaseDenom defines to the default denomination used in titan (staking, governance, etc.)
	BaseDenom = "utkx"
	// BaseDenomUnit defines the base denomination unit for Titan.
	// 1 tkx = 1x10^{BaseDenomUnit} utkx
	BaseDenomUnit = 18
)

func InitSDKConfig() {
	// Set prefixes
	accountPubKeyPrefix := app.AccountAddressPrefix + "pub"
	validatorAddressPrefix := app.AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := app.AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := app.AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := app.AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)

	// Ethermint config coin type to 60
	cmdcfg.SetBip44CoinType(config)

	config.Seal()
}

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(DisplayDenom, sdk.OneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(BaseDenom, sdk.NewDecWithPrec(1, BaseDenomUnit)); err != nil {
		panic(err)
	}
}
