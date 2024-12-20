package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/titantkx/titan/x/farming/keeper"
	"github.com/titantkx/titan/x/farming/types"
)

func SimulateMsgUnstake(_ types.AccountKeeper, _ types.BankKeeper, _ keeper.Keeper) simtypes.Operation {
	//nolint:revive	// keep `chainID` for clear meaning
	return func(r *rand.Rand, _ *baseapp.BaseApp, _ sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgUnstake{
			Sender: simAccount.Address.String(),
		}

		// TODO: Handling the Unstake simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Unstake simulation not implemented"), nil, nil
	}
}
