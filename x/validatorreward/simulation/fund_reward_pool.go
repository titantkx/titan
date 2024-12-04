package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/titantkx/titan/x/validatorreward/keeper"
	"github.com/titantkx/titan/x/validatorreward/types"
)

func SimulateMsgFundRewardPool(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	_ keeper.Keeper,
) simtypes.Operation {
	//nolint:revive	// keep `chainID` for clear meaning
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		depositor, _ := simtypes.RandomAcc(r, accs)
		spendable := bk.SpendableCoins(ctx, depositor.Address)
		depositAmount := simtypes.RandSubsetCoins(r, spendable)

		// if coins slice is empty, we cannot create valid types.MsgFundRewardPool
		if len(depositAmount) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgFundRewardPool, "empty coins slice"), nil, nil
		}

		msg := &types.MsgFundRewardPool{
			Depositor: depositor.Address.String(),
			Amount:    depositAmount,
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      depositor,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: depositAmount,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
