package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/titantkx/titan/x/tokenfactory/keeper"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func SimulateMsgChangeAdmin(
	_ types.AccountKeeper,
	_ types.BankKeeper,
	_ keeper.Keeper,
) simtypes.Operation {
	//nolint:revive	// keep `chainID` for clear meaning
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgChangeAdmin{
			Sender: simAccount.Address.String(),
		}

		// TODO: Handling the ChangeAdmin simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "ChangeAdmin simulation not implemented"), nil, nil
	}
}
