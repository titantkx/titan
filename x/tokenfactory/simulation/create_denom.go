package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/titantkx/titan/x/tokenfactory/keeper"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func SimulateMsgCreateDenom(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	//nolint:revive	// keep `chainID` for clear meaning
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		creator, _ := simtypes.RandomAcc(r, accs)
		spendable := bk.SpendableCoins(ctx, creator.Address)

		fee := k.GetParams(ctx).DenomCreationFee

		if spendable.IsAllLT(fee) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateDenom, "not enough balance to pay denom creation fee"), nil, nil
		}

		subDenom := simtypes.RandStringOfLength(r, types.MaxSubdenomLength)

		msg := &types.MsgCreateDenom{
			Sender:   creator.Address.String(),
			Subdenom: subDenom,
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      creator,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: fee,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
