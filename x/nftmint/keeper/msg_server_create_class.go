package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/titantkx/titan/x/nftmint/types"
)

func (k msgServer) CreateClass(goCtx context.Context, msg *types.MsgCreateClass) (*types.MsgCreateClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	systemInfo, ok := k.GetSystemInfo(ctx)
	if !ok {
		panic("SystemInfo not found")
	}

	classId := strconv.FormatUint(systemInfo.NextClassId, 10)
	classData := types.MustNewAnyWithMetadata(msg.Data)

	class := nft.Class{
		Id:          classId,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Description: msg.Description,
		Uri:         msg.Uri,
		UriHash:     msg.UriHash,
		Data:        classData,
	}

	if err := k.Keeper.nftKeeper.SaveClass(ctx, class); err != nil {
		return nil, types.WrapInternalError(err)
	}

	mintingInfo := types.MintingInfo{
		ClassId:     classId,
		Owner:       msg.Creator,
		NextTokenId: types.DefaultIndex,
	}

	k.Keeper.SetMintingInfo(ctx, mintingInfo)

	systemInfo.NextClassId++
	k.Keeper.SetSystemInfo(ctx, systemInfo)

	ctx.EventManager().EmitTypedEvent(&types.EventCreateClass{
		Id:    classId,
		Owner: msg.Creator,
	})

	return &types.MsgCreateClassResponse{Id: classId}, nil
}
