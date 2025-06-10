package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/titantkx/titan/x/pointer/keeper"
	"github.com/titantkx/titan/x/pointer/types"
)

func SimulateMsgAddERC20NativePointer(
	ak types.AccountKeeper, //nolint:revive
	bk types.BankKeeper, //nolint:revive
	k keeper.Keeper, //nolint:revive
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string, //nolint:revive
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgAddERC20NativePointer{
			Authority: simAccount.Address.String(),
		}

		// TODO: Handling the AddERC20NativePointer simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "AddERC20NativePointer simulation not implemented"), nil, nil
	}
}
