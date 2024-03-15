package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	ethermintcryptocodec "github.com/titantkx/ethermint/crypto/codec"
	ethermint "github.com/titantkx/ethermint/types"

	"github.com/titantkx/titan/app/params"
)

// makeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func makeEncodingConfig() params.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	return params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := makeEncodingConfig()

	// Register legacyAmino codec
	//

	sdk.RegisterLegacyAminoCodec(encodingConfig.Amino)
	// cryptocodec.RegisterCrypto(encodingConfig.Amino)
	codec.RegisterEvidences(encodingConfig.Amino)
	// Ethermint
	ethermintcryptocodec.RegisterCrypto(encodingConfig.Amino)

	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)

	// Register interfaces
	//

	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	// Ethermint
	ethermintcryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ethermint.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}
