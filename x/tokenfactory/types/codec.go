package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateDenom{}, "titan/tokenfactory/create-denom")
	legacy.RegisterAminoMsg(cdc, &MsgMint{}, "titan/tokenfactory/mint")
	legacy.RegisterAminoMsg(cdc, &MsgBurn{}, "titan/tokenfactory/burn")
	legacy.RegisterAminoMsg(cdc, &MsgChangeAdmin{}, "titan/tokenfactory/change-admin")
	legacy.RegisterAminoMsg(cdc, &MsgSetDenomMetadata{}, "titan/tokenfactory/set-denom-metadata")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "titan/tokenfactory/update-params")
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateDenom{},
		&MsgMint{},
		&MsgBurn{},
		&MsgChangeAdmin{},
		&MsgSetDenomMetadata{},
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
