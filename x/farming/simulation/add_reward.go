package simulation

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/titantkx/titan/x/farming/keeper"
	"github.com/titantkx/titan/x/farming/types"
)

func SimulateMsgAddReward(ak types.AccountKeeper, bk types.BankKeeper, _ keeper.Keeper) simtypes.Operation {
	//nolint:revive	// keep `chainID` for clear meaning
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		sender, _ := simtypes.RandomAcc(r, accs)
		spendable := bk.SpendableCoins(ctx, sender.Address)

		if spendable.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAddReward, "no blance to add reward"), nil, nil
		}

		startTime := ctx.BlockHeader().Time.Add(time.Duration(simtypes.RandIntBetween(r, 0, 100_0000_000)))
		endTime := startTime.Add(time.Duration(simtypes.RandIntBetween(r, 0, 100_0000_000)))

		msg := &types.MsgAddReward{
			Sender:    sender.Address.String(),
			Token:     "stake",
			Amount:    simtypes.RandSubsetCoins(r, spendable),
			EndTime:   endTime,
			StartTime: startTime,
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    sender,
			AccountKeeper: ak,
			Bankkeeper:    bk,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
