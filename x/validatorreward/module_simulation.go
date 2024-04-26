package validatorreward

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/titantkx/titan/testutil/sample"
	validatorrewardsimulation "github.com/titantkx/titan/x/validatorreward/simulation"
	"github.com/titantkx/titan/x/validatorreward/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = validatorrewardsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	//nolint:gosec // this is not credentials
	opWeightMsgSetRate = "op_weight_msg_set_rate"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetRate int = 100

	//nolint:gosec // this is not credentials
	opWeightMsgSetAuthority = "op_weight_msg_set_authority"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetAuthority int = 100

	//nolint:gosec // this is not credentials
	opWeightMsgFundRewardPool = "op_weight_msg_fund_reward_pool"
	// TODO: Determine the simulation weight value
	defaultWeightMsgFundRewardPool int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	validatorrewardGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&validatorrewardGenesis)
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

	var weightMsgSetRate int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSetRate, &weightMsgSetRate, nil,
		func(_ *rand.Rand) {
			weightMsgSetRate = defaultWeightMsgSetRate
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetRate,
		validatorrewardsimulation.SimulateMsgSetRate(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSetAuthority int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSetAuthority, &weightMsgSetAuthority, nil,
		func(_ *rand.Rand) {
			weightMsgSetAuthority = defaultWeightMsgSetAuthority
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetAuthority,
		validatorrewardsimulation.SimulateMsgSetAuthority(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgFundRewardPool int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgFundRewardPool, &weightMsgFundRewardPool, nil,
		func(_ *rand.Rand) {
			weightMsgFundRewardPool = defaultWeightMsgFundRewardPool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgFundRewardPool,
		validatorrewardsimulation.SimulateMsgFundRewardPool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgSetRate,
			defaultWeightMsgSetRate,
			func(_ *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
				validatorrewardsimulation.SimulateMsgSetRate(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgSetAuthority,
			defaultWeightMsgSetAuthority,
			func(_ *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
				validatorrewardsimulation.SimulateMsgSetAuthority(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgFundRewardPool,
			defaultWeightMsgFundRewardPool,
			func(_ *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
				validatorrewardsimulation.SimulateMsgFundRewardPool(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
