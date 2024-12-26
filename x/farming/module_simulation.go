package farming

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/titantkx/titan/testutil/sample"
	farmingsimulation "github.com/titantkx/titan/x/farming/simulation"
	"github.com/titantkx/titan/x/farming/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = farmingsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	//nolint:gosec // this is not credentials
	opWeightMsgAddReward = "op_weight_msg_add_reward"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddReward int = 100

	//nolint:gosec // this is not credentials
	opWeightMsgStake = "op_weight_msg_stake"
	// TODO: Determine the simulation weight value
	defaultWeightMsgStake int = 100

	//nolint:gosec // this is not credentials
	opWeightMsgUnstake = "op_weight_msg_unstake"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUnstake int = 100

	//nolint:gosec // this is not credentials
	opWeightMsgHarvest = "op_weight_msg_harvest"
	// TODO: Determine the simulation weight value
	defaultWeightMsgHarvest int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	addRewardGas := simtypes.RandIntBetween(simState.Rand, 0, 5_000_000)

	farmingGenesis := types.GenesisState{
		Params: types.Params{
			//nolint:gosec // G115
			AddRewardGas: uint64(addRewardGas),
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&farmingGenesis)
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

	var weightMsgAddReward int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddReward, &weightMsgAddReward, nil,
		func(_ *rand.Rand) {
			weightMsgAddReward = defaultWeightMsgAddReward
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddReward,
		farmingsimulation.SimulateMsgAddReward(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgStake int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgStake, &weightMsgStake, nil,
		func(_ *rand.Rand) {
			weightMsgStake = defaultWeightMsgStake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgStake,
		farmingsimulation.SimulateMsgStake(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUnstake int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUnstake, &weightMsgUnstake, nil,
		func(_ *rand.Rand) {
			weightMsgUnstake = defaultWeightMsgUnstake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUnstake,
		farmingsimulation.SimulateMsgUnstake(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgHarvest int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgHarvest, &weightMsgHarvest, nil,
		func(_ *rand.Rand) {
			weightMsgHarvest = defaultWeightMsgHarvest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgHarvest,
		farmingsimulation.SimulateMsgHarvest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddReward,
			defaultWeightMsgAddReward,
			func(_ *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
				farmingsimulation.SimulateMsgAddReward(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgStake,
			defaultWeightMsgStake,
			func(_ *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
				farmingsimulation.SimulateMsgStake(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUnstake,
			defaultWeightMsgUnstake,
			func(_ *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
				farmingsimulation.SimulateMsgUnstake(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgHarvest,
			defaultWeightMsgHarvest,
			func(_ *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
				farmingsimulation.SimulateMsgHarvest(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
