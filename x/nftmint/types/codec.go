package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateClass{}, "nftmint/CreateClass", nil)
	cdc.RegisterConcrete(&MsgMint{}, "nftmint/Mint", nil)
	cdc.RegisterConcrete(&Metadata{}, "nftmint/Metadata", nil)
	cdc.RegisterConcrete(&MsgUpdateClass{}, "nftmint/UpdateClass", nil)
	cdc.RegisterConcrete(&MsgTransferClass{}, "nftmint/TransferClass", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateClass{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMint{},
	)
	registry.RegisterImplementations((*MetadataI)(nil),
		&Metadata{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateClass{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransferClass{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
