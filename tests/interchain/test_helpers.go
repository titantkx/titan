package interchain

import (
	"context"

	sdktestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	interchaintest "github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"

	titanapp "github.com/titantkx/titan/app"
	titanutils "github.com/titantkx/titan/utils"
)

var (
	IBCRelayerImage   = "ghcr.io/cosmos/relayer"
	IBCRelayerVersion = "main"

	titanImageInfo = []ibc.DockerImage{
		{
			Repository: "docker.io/titantkx/titand",
			Version:    "local",
			UidGid:     "1025:1025",
		},
	}

	overrideGenesisKV = []cosmos.GenesisKV{
		// {
		// 	Key:   "app_state.gov.params.min_deposit.0.denom",
		// 	Value: titanutils.BondDenom,
		// },
		// {
		// 	Key:   "app_state.feepay.params.enable_feepay",
		// 	Value: false,
		// },
		{
			Key:   "app_state.evm.params.evm_denom",
			Value: titanutils.BondDenom,
		},
	}
)

func titanEncoding() *sdktestutil.TestEncodingConfig {
	encodingConfig := titanapp.MakeEncodingConfig()

	return &sdktestutil.TestEncodingConfig{
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		Codec:             encodingConfig.Marshaler,
		TxConfig:          encodingConfig.TxConfig,
		Amino:             encodingConfig.Amino,
	}
}

func TitanChainSpec(
	_ context.Context,
	chainID string,
	nv, nf int,
) *interchaintest.ChainSpec {
	decimals := int64(18)
	return &interchaintest.ChainSpec{
		NumValidators: &nv,
		NumFullNodes:  &nf,
		ChainConfig: ibc.ChainConfig{
			Type:           "cosmos",
			Name:           "titan",
			ChainID:        chainID,
			Bin:            "titand",
			Denom:          "atkx",
			Bech32Prefix:   "titan",
			CoinType:       "60",
			GasPrices:      "100000000000atkx",
			GasAdjustment:  2,
			TrustingPeriod: "168h0m0s",
			NoHostMount:    false,
			Images:         titanImageInfo,
			EncodingConfig: titanEncoding(),
			CoinDecimals:   &decimals,
			ModifyGenesis:  cosmos.ModifyGenesis(overrideGenesisKV),
			UseGasUsed:     true,
		},
	}
}
