package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	groupcodec "github.com/cosmos/cosmos-sdk/x/group/codec"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSetRate{}, "validatorreward/SetRate", nil)
	cdc.RegisterConcrete(&MsgSetAuthority{}, "validatorreward/SetAuthority", nil)
	cdc.RegisterConcrete(&MsgFundRewardPool{}, "validatorreward/FundRewardPool", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetRate{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetAuthority{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFundRewardPool{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	// AminoCdc is a amino codec created to support amino JSON compatible msgs.
	AminoCdc = codec.NewAminoCodec(Amino)
)

func init() {
	RegisterCodec(Amino)
	cryptocodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)

	// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterCodec(authzcodec.Amino)
	RegisterCodec(govcodec.Amino)
	RegisterCodec(groupcodec.Amino)
}
