package cmd_test

import (
	"fmt"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd/feemarket"
	"github.com/titantkx/titan/testutil/cmd/gov"
	"github.com/titantkx/titan/testutil/cmd/keys"
	"github.com/titantkx/titan/testutil/cmd/nft"
	"github.com/titantkx/titan/testutil/sample"

	"github.com/titantkx/titan/testutil/cmd/pointer"
	"github.com/titantkx/titan/testutil/cmd/staking"
	"github.com/titantkx/titan/testutil/cmd/tokenfactory"
	"github.com/titantkx/titan/utils"
)

type Deposit struct {
	From   string
	Amount string
}

type Vote struct {
	From   string
	Option string
}

func MustCreateVoter(t testing.TB, balance string, stakeAmount string) string {
	val := MustGetValidator(t)
	del := MustCreateAccount(t, balance).Address

	staking.MustDelegate(t, val, stakeAmount, del)
	t.Cleanup(func() {
		staking.MustUnbond(t, val, stakeAmount, del)
	})

	return del
}

func TestSubmitProposals(t *testing.T) {
	govParams := gov.MustGetParams(t)

	require.Equal(t, "0.334", govParams.Quorum.String())
	require.Equal(t, "0.5", govParams.Threshold.String())
	require.Equal(t, "0.334", govParams.VetoThreshold.String())
	require.Equal(t, "250000000000000000000"+utils.BaseDenom, govParams.MinDeposit.String())

	voter1 := keys.MustShowAddress(t, "val1") // Will represent voter3, voter4, voter5 if they do not vote
	voter2 := keys.MustShowAddress(t, "val2")
	voter3 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)
	voter4 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)
	voter5 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)

	originalFeeMarketParams := feemarket.MustGetParams(t)

	tests := []struct {
		Name           string
		Proposer       string
		Proposal       gov.ProposalMsg
		Deposits       []Deposit
		Votes          []Vote
		ExpectedStatus string
		CallbackFunc   func(string, []interface{}) // CallbackFunc will be called after the proposal is finished
	}{
		// PROPOSAL_STATUS_PASSED
		{
			"TestSubmitTextProposalAllYesPassed",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalAllYesPassed",
				Summary:  "TestSubmitTextProposalAllYesPassed",
				Metadata: "TestSubmitTextProposalAllYesPassed",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
				{voter3, gov.VOTE_OPTION_YES},
				{voter4, gov.VOTE_OPTION_YES},
				{voter5, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			nil,
		},
		{
			"TestSubmitTextProposalDepositLaterAllYesPassed",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalDepositLaterAllYesPassed",
				Summary:  "TestSubmitTextProposalDepositLaterAllYesPassed",
				Metadata: "TestSubmitTextProposalDepositLaterAllYesPassed",
				Deposit:  "150" + utils.DisplayDenom,
			},
			[]Deposit{
				{voter2, "50" + utils.DisplayDenom},
				{voter3, "50" + utils.DisplayDenom},
			},
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
				{voter3, gov.VOTE_OPTION_YES},
				{voter4, gov.VOTE_OPTION_YES},
				{voter5, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			nil,
		},
		{
			"TestSubmitTextProposalTwoYesOneNoPassed",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalTwoYesOneNoPassed",
				Summary:  "TestSubmitTextProposalTwoYesOneNoPassed",
				Metadata: "TestSubmitTextProposalTwoYesOneNoPassed",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter2, gov.VOTE_OPTION_YES},
				{voter3, gov.VOTE_OPTION_YES},
				{voter4, gov.VOTE_OPTION_NO},
			},
			gov.PROPOSAL_STATUS_PASSED,
			nil,
		},
		{
			"TestSubmitTextProposalOneYesTwoAbstainPassed",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalOneYesTwoAbstainPassed",
				Summary:  "TestSubmitTextProposalOneYesTwoAbstainPassed",
				Metadata: "TestSubmitTextProposalOneYesTwoAbstainPassed",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter2, gov.VOTE_OPTION_YES},
				{voter3, gov.VOTE_OPTION_ABSTAIN},
				{voter4, gov.VOTE_OPTION_ABSTAIN},
			},
			gov.PROPOSAL_STATUS_PASSED,
			nil,
		},
		{
			"TestSubmitTextProposalTwoYesPassed",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalTwoYesPassed",
				Summary:  "TestSubmitTextProposalTwoYesPassed",
				Metadata: "TestSubmitTextProposalTwoYesPassed",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter2, gov.VOTE_OPTION_YES},
				{voter3, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			nil,
		},
		{
			"TestSubmitUpdateParamsProposalSetMinGasPriceToZeroPassed",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitUpdateParamsProposalSetMinGasPriceToZeroPassed",
				Summary:  "TestSubmitUpdateParamsProposalSetMinGasPriceToZeroPassed",
				Metadata: "TestSubmitUpdateParamsProposalSetMinGasPriceToZeroPassed",
				Deposit:  "250" + utils.DisplayDenom,
				Messages: []any{
					gov.MsgUpdateParams{
						Type:      "/ethermint.feemarket.v1.MsgUpdateParams",
						Authority: "titan10d07y265gmmuvt4z0w9aw880jnsr700jste397",
						Params: feemarket.Params{
							BaseFee:                  originalFeeMarketParams.BaseFee,
							BaseFeeChangeDenominator: originalFeeMarketParams.BaseFeeChangeDenominator,
							ElasticityMultiplier:     originalFeeMarketParams.ElasticityMultiplier,
							EnableHeight:             originalFeeMarketParams.EnableHeight,
							MinGasMultiplier:         originalFeeMarketParams.MinGasMultiplier,
							MinGasPrice:              testutil.MakeFloat(0),
							NoBaseFee:                originalFeeMarketParams.NoBaseFee,
						},
					},
				},
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			func(_ string, _ []interface{}) {
				minGasPrice := feemarket.MustGetParams(t).MinGasPrice
				require.Equal(t, "0", minGasPrice.String())
			},
		},
		{
			"TestSubmitUpdateParamsProposalSetMinGasPriceToOriginalPassed",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitUpdateParamsProposalSetMinGasPriceToOriginalPassed",
				Summary:  "TestSubmitUpdateParamsProposalSetMinGasPriceToOriginalPassed",
				Metadata: "TestSubmitUpdateParamsProposalSetMinGasPriceToOriginalPassed",
				Deposit:  "250" + utils.DisplayDenom,
				Messages: []any{
					gov.MsgUpdateParams{
						Type:      "/ethermint.feemarket.v1.MsgUpdateParams",
						Authority: "titan10d07y265gmmuvt4z0w9aw880jnsr700jste397",
						Params:    originalFeeMarketParams,
					},
				},
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			func(_ string, _ []interface{}) {
				minGasPrice := feemarket.MustGetParams(t).MinGasPrice
				require.Equal(t, originalFeeMarketParams.MinGasPrice.String(), minGasPrice.String())
			},
		},
		{
			"TestSubmitNftCreateClassProposalPassed",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitNftCreateClassProposalPassed",
				Summary:  "TestSubmitNftCreateClassProposalPassed",
				Metadata: "TestSubmitNftCreateClassProposalPassed",
				Deposit:  "250" + utils.DisplayDenom,
				Messages: []any{
					gov.MsgNftCreateClass{
						Type:        "/titan.nftmint.MsgCreateClass",
						Creator:     "titan10d07y265gmmuvt4z0w9aw880jnsr700jste397",
						Name:        sample.Word(),
						Symbol:      sample.Word(),
						Description: sample.Paragraph(),
						Uri:         sample.URL(),
						UriHash:     sample.Hash(),
						Data:        sample.JSON(),
					},
				},
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			func(_ string, msgs []interface{}) {
				msg, ok := msgs[0].(gov.MsgNftCreateClass)
				require.True(t, ok)
				latestClass := nft.MustGetLatestClass(t)

				require.Equal(t, msg.Name, latestClass.Name)
			},
		},
		// PROPOSAL_STATUS_REJECTED
		{
			"TestSubmitTextProposalOneYesTwoNoRejected",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalOneYesTwoNoRejected",
				Summary:  "TestSubmitTextProposalOneYesTwoNoRejected",
				Metadata: "TestSubmitTextProposalOneYesTwoNoRejected",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter2, gov.VOTE_OPTION_YES},
				{voter3, gov.VOTE_OPTION_NO},
				{voter4, gov.VOTE_OPTION_NO},
			},
			gov.PROPOSAL_STATUS_REJECTED,
			nil,
		},
		{
			"TestSubmitTextProposalOneYesOneNoOneAbstainRejected",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalOneYesOneNoOneAbstainRejected",
				Summary:  "TestSubmitTextProposalOneYesOneNoOneAbstainRejected",
				Metadata: "TestSubmitTextProposalOneYesOneNoOneAbstainRejected",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter2, gov.VOTE_OPTION_YES},
				{voter3, gov.VOTE_OPTION_NO},
				{voter4, gov.VOTE_OPTION_ABSTAIN},
			},
			gov.PROPOSAL_STATUS_REJECTED,
			nil,
		},
		{
			"TestSubmitTextProposalOneYesRejected",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalOneYesRejected",
				Summary:  "TestSubmitTextProposalOneYesRejected",
				Metadata: "TestSubmitTextProposalOneYesRejected",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter2, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_REJECTED,
			nil,
		},
		{
			"TestSubmitTextProposalThreeYesTwoVetoRejected",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalThreeYesTwoVetoRejected",
				Summary:  "TestSubmitTextProposalThreeYesTwoVetoRejected",
				Metadata: "TestSubmitTextProposalThreeYesTwoVetoRejected",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
				{voter3, gov.VOTE_OPTION_YES},
				{voter4, gov.VOTE_OPTION_NO_WITH_VETO},
				{voter5, gov.VOTE_OPTION_NO_WITH_VETO},
			},
			gov.PROPOSAL_STATUS_REJECTED,
			nil,
		},
		// PROPOSAL_STATUS_DEPOSIT_FAILED
		{
			"TestSubmitTextProposalNotEnoughDepositFailed",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalNotEnoughDepositFailed",
				Summary:  "TestSubmitTextProposalNotEnoughDepositFailed",
				Metadata: "TestSubmitTextProposalNotEnoughDepositFailed",
				Deposit:  "100" + utils.DisplayDenom,
			},
			nil,
			nil,
			gov.PROPOSAL_STATUS_DEPOSIT_FAILED,
			nil,
		},
		// PROPOSAL_STATUS_FAILED
		{
			"TestSubmitValidatorRewardSetRateProposalPassed",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitValidatorRewardSetRateProposalPassed",
				Summary:  "TestSubmitValidatorRewardSetRateProposalPassed",
				Metadata: "TestSubmitValidatorRewardSetRateProposalPassed",
				Deposit:  "250" + utils.DisplayDenom,
				Messages: []any{
					gov.MsgValidatorRewardSetRate{
						Type:      "/titan.validatorreward.MsgSetRate",
						Authority: "titan10d07y265gmmuvt4z0w9aw880jnsr700jste397",
						Rate:      sdk.NewDecWithPrec(1, 1).String(),
					},
				},
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_FAILED, // because Authority (gov module) is not allowed to set rate
			nil,
		},
		{
			"TestSubmitCreateErc20NativePointerForBaseDenom",
			voter1,
			gov.ProposalMsg{
				Title:    "TestSubmitCreateErc20NativePointerForBaseDenom",
				Summary:  "TestSubmitCreateErc20NativePointerForBaseDenom",
				Metadata: "TestSubmitCreateErc20NativePointerForBaseDenom",
				Deposit:  "250" + utils.DisplayDenom,
				Messages: []any{
					gov.MsgAddERC20NativePointer{
						Type:      "/titan.pointer.MsgAddERC20NativePointer",
						Authority: "titan10d07y265gmmuvt4z0w9aw880jnsr700jste397",
						Token:     "atkx",
						Name:      "tkx",
						Symbol:    "TKX",
						Decimals:  18,
					},
				},
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_FAILED, // because can not create pointer for base denom
			nil,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			testSubmitProposal(
				t,
				test.Proposer,
				test.Proposal,
				test.Deposits,
				test.Votes,
				test.ExpectedStatus,
				test.CallbackFunc,
			)
		})
	}
}

func testSubmitProposal(
	t *testing.T,
	proposer string,
	proposalMsg gov.ProposalMsg,
	deposits []Deposit,
	votes []Vote,
	expectedStatus string,
	callbackFunc func(string, []interface{}),
) {
	proposalId := gov.MustSubmitProposal(t, proposer, proposalMsg)

	var wg1 sync.WaitGroup
	for i := range deposits {
		deposit := deposits[i]
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			gov.MustDeposit(t, deposit.From, proposalId, deposit.Amount)
		}()
	}
	wg1.Wait()

	if expectedStatus == gov.PROPOSAL_STATUS_DEPOSIT_FAILED {
		gov.MustNotPassDepositPeriod(t, proposalId)
		return
	}

	proposal := gov.MustQueryPassDepositPeriodProposal(t, proposalId)

	require.Equal(t, gov.PROPOSAL_STATUS_VOTING_PERIOD, proposal.Status)

	var wg2 sync.WaitGroup
	for i := range votes {
		vote := votes[i]
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			gov.MustVote(t, vote.From, proposalId, vote.Option)
		}()
	}
	wg2.Wait()

	if callbackFunc == nil {
		t.Parallel() // Should run in parallel from here if there is no callback function
	}

	proposal = gov.MustQueryPassVotingPeriodProposal(t, proposalId)

	require.Equal(t, expectedStatus, proposal.Status)

	if callbackFunc != nil {
		callbackFunc(proposalId, proposalMsg.Messages)
	}
}

func TestSubmitAddErc20NativePointer(t *testing.T) {
	// Create new token factory
	creator := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	tokenDenom := tokenfactory.MustCreateNewToken(t, creator, "erc20nativepointerTest")
	// Mint token factory
	tokenfactory.MustMintTokenFactory(t, creator, creator, tokenDenom, testutil.MakeInt(1000000000))

	// Create new token factory 2
	creator2 := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	tokenDenom2 := tokenfactory.MustCreateNewToken(t, creator2, "erc20nativepointerTest2")
	// Mint token factory
	tokenfactory.MustMintTokenFactory(t, creator2, creator2, tokenDenom2, testutil.MakeInt(1000000000))

	govParams := gov.MustGetParams(t)

	require.Equal(t, "0.334", govParams.Quorum.String())
	require.Equal(t, "0.5", govParams.Threshold.String())
	require.Equal(t, "0.334", govParams.VetoThreshold.String())
	require.Equal(t, "250000000000000000000"+utils.BaseDenom, govParams.MinDeposit.String())

	voter1 := keys.MustShowAddress(t, "val1") // Will represent voter3, voter4, voter5 if they do not vote
	voter2 := keys.MustShowAddress(t, "val2")

	proposal := struct {
		Name           string
		Proposer       string
		Proposal       gov.ProposalMsg
		Deposits       []Deposit
		Votes          []Vote
		ExpectedStatus string
		CallbackFunc   func(string, []interface{}) // CallbackFunc will be called after the proposal is finished
	}{
		"TestSubmitAddErc20NativePointer",
		voter1,
		gov.ProposalMsg{
			Title:    "TestSubmitAddErc20NativePointer",
			Summary:  "TestSubmitAddErc20NativePointer",
			Metadata: "TestSubmitAddErc20NativePointer",
			Deposit:  "250" + utils.DisplayDenom,
			Messages: []any{
				gov.MsgAddERC20NativePointer{
					Type:      "/titan.pointer.MsgAddERC20NativePointer",
					Authority: "titan10d07y265gmmuvt4z0w9aw880jnsr700jste397",
					Token:     tokenDenom,
					Name:      "erc20nativepointerTest",
					Symbol:    "ERC20NATIVEPOINTERTEST",
					Decimals:  6,
				},
			},
		},
		nil,
		[]Vote{
			{voter1, gov.VOTE_OPTION_YES},
			{voter2, gov.VOTE_OPTION_YES},
		},
		gov.PROPOSAL_STATUS_PASSED,
		func(s string, i []interface{}) {}, //nolint:revive
	}

	testSubmitProposal(
		t,
		proposal.Proposer,
		proposal.Proposal,
		proposal.Deposits,
		proposal.Votes,
		proposal.ExpectedStatus,
		proposal.CallbackFunc,
	)

	erc20PointerAddr := pointer.MustGetERC20Pointer(t, tokenDenom)
	fmt.Println("ERC20 Pointer Address:", erc20PointerAddr)

	proposal2 := struct {
		Name           string
		Proposer       string
		Proposal       gov.ProposalMsg
		Deposits       []Deposit
		Votes          []Vote
		ExpectedStatus string
		CallbackFunc   func(string, []interface{}) // CallbackFunc will be called after the proposal is finished
	}{
		"TestSubmitAddErc20NativePointer2",
		voter1,
		gov.ProposalMsg{
			Title:    "TestSubmitAddErc20NativePointer2",
			Summary:  "TestSubmitAddErc20NativePointer2",
			Metadata: "TestSubmitAddErc20NativePointer2",
			Deposit:  "250" + utils.DisplayDenom,
			Messages: []any{
				gov.MsgAddERC20NativePointer{
					Type:      "/titan.pointer.MsgAddERC20NativePointer",
					Authority: "titan10d07y265gmmuvt4z0w9aw880jnsr700jste397",
					Token:     tokenDenom2,
					Name:      "erc20nativepointerTest2",
					Symbol:    "ERC20NATIVEPOINTERTEST2",
					Decimals:  8,
				},
			},
		},
		nil,
		[]Vote{
			{voter1, gov.VOTE_OPTION_YES},
			{voter2, gov.VOTE_OPTION_YES},
		},
		gov.PROPOSAL_STATUS_PASSED,
		func(s string, i []interface{}) {}, //nolint:revive
	}

	testSubmitProposal(
		t,
		proposal2.Proposer,
		proposal2.Proposal,
		proposal2.Deposits,
		proposal2.Votes,
		proposal2.ExpectedStatus,
		proposal2.CallbackFunc,
	)

	erc20PointerAddr2 := pointer.MustGetERC20Pointer(t, tokenDenom2)
	fmt.Println("ERC20 Pointer Address 2:", erc20PointerAddr2)

	// Verify that the two pointers are different
	require.NotEqual(t, erc20PointerAddr, erc20PointerAddr2, "ERC20 Pointer addresses should be different")
}
