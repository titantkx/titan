package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/titantkx/titan/x/farming/keeper"
	"github.com/titantkx/titan/x/farming/types"
)

func SimulateMsgStake(ak types.AccountKeeper, bk types.BankKeeper, _ keeper.Keeper) simtypes.Operation {
	//nolint:revive	// keep `chainID` for clear meaning
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		sender, _ := simtypes.RandomAcc(r, accs)
		spendable := bk.SpendableCoins(ctx, sender.Address)

		if spendable.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStake, "no balance to stake"), nil, nil
		}

		msg := &types.MsgStake{
			Sender: sender.Address.String(),
			Amount: simtypes.RandSubsetCoins(r, spendable),
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
