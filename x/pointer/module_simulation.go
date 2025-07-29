package pointer

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/titantkx/titan/testutil/sample"
	pointersimulation "github.com/titantkx/titan/x/pointer/simulation"
	"github.com/titantkx/titan/x/pointer/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = pointersimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgAddERC20NativePointer = "op_weight_msg_add_erc_20_native_pointer" //nolint:gosec
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddERC20NativePointer int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	pointerGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&pointerGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent { //nolint:staticcheck
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgAddERC20NativePointer int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddERC20NativePointer, &weightMsgAddERC20NativePointer, nil,
		func(_ *rand.Rand) {
			weightMsgAddERC20NativePointer = defaultWeightMsgAddERC20NativePointer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddERC20NativePointer,
		pointersimulation.SimulateMsgAddERC20NativePointer(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg { //nolint:revive
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddERC20NativePointer,
			defaultWeightMsgAddERC20NativePointer,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg { //nolint:revive
				pointersimulation.SimulateMsgAddERC20NativePointer(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
