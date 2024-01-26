package cmd_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd/feemarket"
	"github.com/tokenize-titan/titan/testutil/cmd/gov"
	"github.com/tokenize-titan/titan/testutil/cmd/keys"
	"github.com/tokenize-titan/titan/testutil/cmd/staking"
	"github.com/tokenize-titan/titan/utils"
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
	t.Parallel()

	govParams := gov.MustGetParams(t)

	require.Equal(t, govParams.MinDeposit.String(), "2.5e+20"+utils.BaseDenom)
	require.Equal(t, govParams.Quorum.String(), "0.334")
	require.Equal(t, govParams.Threshold.String(), "0.5")
	require.Equal(t, govParams.VetoThreshold.String(), "0.334")

	proposer := keys.MustShowAddress(t, "val1") // Will represent others if they do not vote
	voter1 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)
	voter2 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)
	voter3 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)
	voter4 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)

	originalFeeMarketParams := feemarket.MustGetParams(t)

	tests := []struct {
		Name           string
		Proposer       string
		Proposal       gov.ProposalMsg
		Deposits       []Deposit
		Votes          []Vote
		ExpectedStatus string
		CallbackFunc   func() // CallbackFunc will be called after the proposal is finished
	}{
		// PROPOSAL_STATUS_PASSED
		{
			"TestSubmitTextProposalAllYesPassed",
			proposer,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalAllYesPassed",
				Summary:  "TestSubmitTextProposalAllYesPassed",
				Metadata: "TestSubmitTextProposalAllYesPassed",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{proposer, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			nil,
		},
		{
			"TestSubmitTextProposalDepositLaterAllYesPassed",
			proposer,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalDepositLaterAllYesPassed",
				Summary:  "TestSubmitTextProposalDepositLaterAllYesPassed",
				Metadata: "TestSubmitTextProposalDepositLaterAllYesPassed",
				Deposit:  "150" + utils.DisplayDenom,
			},
			[]Deposit{
				{voter1, "50" + utils.DisplayDenom},
				{voter2, "50" + utils.DisplayDenom},
			},
			[]Vote{
				{proposer, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			nil,
		},
		{
			"TestSubmitTextProposalTwoYesOneNoPassed",
			proposer,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalTwoYesOneNoPassed",
				Summary:  "TestSubmitTextProposalTwoYesOneNoPassed",
				Metadata: "TestSubmitTextProposalTwoYesOneNoPassed",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
				{voter3, gov.VOTE_OPTION_NO},
			},
			gov.PROPOSAL_STATUS_PASSED,
			nil,
		},
		{
			"TestSubmitTextProposalOneYesTwoAbstainPassed",
			proposer,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalOneYesTwoAbstainPassed",
				Summary:  "TestSubmitTextProposalOneYesTwoAbstainPassed",
				Metadata: "TestSubmitTextProposalOneYesTwoAbstainPassed",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_ABSTAIN},
				{voter3, gov.VOTE_OPTION_ABSTAIN},
			},
			gov.PROPOSAL_STATUS_PASSED,
			nil,
		},
		{
			"TestSubmitTextProposalTwoYesPassed",
			proposer,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalTwoYesPassed",
				Summary:  "TestSubmitTextProposalTwoYesPassed",
				Metadata: "TestSubmitTextProposalTwoYesPassed",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			nil,
		},
		{
			"TestSubmitUpdateParamsProposalSetMinGasPriceToZeroPassed",
			proposer,
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
				{proposer, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			func() {
				minGasPrice := feemarket.MustGetParams(t).MinGasPrice
				require.Equal(t, "0", minGasPrice.String())
			},
		},
		{
			"TestSubmitUpdateParamsProposalSetMinGasPriceToOriginalPassed",
			proposer,
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
				{proposer, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_PASSED,
			func() {
				minGasPrice := feemarket.MustGetParams(t).MinGasPrice
				require.Equal(t, originalFeeMarketParams.MinGasPrice.String(), minGasPrice.String())
			},
		},
		// PROPOSAL_STATUS_REJECTED
		{
			"TestSubmitTextProposalOneYesTwoNoRejected",
			proposer,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalOneYesTwoNoRejected",
				Summary:  "TestSubmitTextProposalOneYesTwoNoRejected",
				Metadata: "TestSubmitTextProposalOneYesTwoNoRejected",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_NO},
				{voter3, gov.VOTE_OPTION_NO},
			},
			gov.PROPOSAL_STATUS_REJECTED,
			nil,
		},
		{
			"TestSubmitTextProposalOneYesOneNoOneAbstainRejected",
			proposer,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalOneYesOneNoOneAbstainRejected",
				Summary:  "TestSubmitTextProposalOneYesOneNoOneAbstainRejected",
				Metadata: "TestSubmitTextProposalOneYesOneNoOneAbstainRejected",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_NO},
				{voter3, gov.VOTE_OPTION_ABSTAIN},
			},
			gov.PROPOSAL_STATUS_REJECTED,
			nil,
		},
		{
			"TestSubmitTextProposalOneYesRejected",
			proposer,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalOneYesRejected",
				Summary:  "TestSubmitTextProposalOneYesRejected",
				Metadata: "TestSubmitTextProposalOneYesRejected",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{voter1, gov.VOTE_OPTION_YES},
			},
			gov.PROPOSAL_STATUS_REJECTED,
			nil,
		},
		{
			"TestSubmitTextProposalThreeYesTwoVetoRejected",
			proposer,
			gov.ProposalMsg{
				Title:    "TestSubmitTextProposalThreeYesTwoVetoRejected",
				Summary:  "TestSubmitTextProposalThreeYesTwoVetoRejected",
				Metadata: "TestSubmitTextProposalThreeYesTwoVetoRejected",
				Deposit:  "250" + utils.DisplayDenom,
			},
			nil,
			[]Vote{
				{proposer, gov.VOTE_OPTION_YES},
				{voter1, gov.VOTE_OPTION_YES},
				{voter2, gov.VOTE_OPTION_YES},
				{voter3, gov.VOTE_OPTION_NO_WITH_VETO},
				{voter4, gov.VOTE_OPTION_NO_WITH_VETO},
			},
			gov.PROPOSAL_STATUS_REJECTED,
			nil,
		},
		// PROPOSAL_STATUS_DEPOSIT_FAILED
		{
			"TestSubmitTextProposalNotEnoughDepositFailed",
			proposer,
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
	callbackFunc func(),
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
		callbackFunc()
	}
}
