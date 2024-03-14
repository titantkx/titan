package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/titantkx/titan/x/validatorreward/keeper"
	"github.com/titantkx/titan/x/validatorreward/types"
)

func SimulateMsgFundRewardPool(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgFundRewardPool{
			Depositor: simAccount.Address.String(),
		}

		// TODO: Handling the FundRewardPool simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "FundRewardPool simulation not implemented"), nil, nil
	}
}
