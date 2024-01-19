package cmd_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
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
	del1 := keys.MustShowAddress(t, "val1")
	val1 := testutil.MustAccountAddressToValidatorAddress(t, del1)
	del := MustCreateAccount(t, balance).Address

	staking.MustDelegate(t, val1, stakeAmount, del)
	t.Cleanup(func() {
		staking.MustUnbond(t, val1, stakeAmount, del)
	})

	return del
}

func TestSubmitProposals(t *testing.T) {
	t.Parallel()

	params := gov.MustGetParams(t)

	require.Equal(t, params.MinDeposit.String(), "2.5e+20"+utils.BaseDenom)
	require.Equal(t, params.Quorum.String(), "0.334")
	require.Equal(t, params.Threshold.String(), "0.5")
	require.Equal(t, params.VetoThreshold.String(), "0.334")

	proposer := keys.MustShowAddress(t, "val1") // Will represent others if they do not vote
	voter1 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)
	voter2 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)
	voter3 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)
	voter4 := MustCreateVoter(t, "1000000"+utils.DisplayDenom, "100000"+utils.DisplayDenom)

	tests := []struct {
		Name           string
		Proposer       string
		Proposal       any
		Deposits       []Deposit
		Votes          []Vote
		ExpectedStatus string
	}{
		// PROPOSAL_STATUS_PASSED
		{
			"TestSubmitTextProposalAllYesPassed",
			proposer,
			gov.TextProposal{
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
		},
		{
			"TestSubmitTextProposalDepositLaterAllYesPassed",
			proposer,
			gov.TextProposal{
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
		},
		{
			"TestSubmitTextProposalTwoYesOneNoPassed",
			proposer,
			gov.TextProposal{
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
		},
		{
			"TestSubmitTextProposalOneYesTwoAbstainPassed",
			proposer,
			gov.TextProposal{
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
		},
		{
			"TestSubmitTextProposalTwoYesPassed",
			proposer,
			gov.TextProposal{
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
		},
		// PROPOSAL_STATUS_REJECTED
		{
			"TestSubmitTextProposalOneYesTwoNoRejected",
			proposer,
			gov.TextProposal{
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
		},
		{
			"TestSubmitTextProposalOneYesOneNoOneAbstainRejected",
			proposer,
			gov.TextProposal{
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
		},
		{
			"TestSubmitTextProposalOneYesRejected",
			proposer,
			gov.TextProposal{
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
		},
		{
			"TestSubmitTextProposalThreeYesTwoVetoRejected",
			proposer,
			gov.TextProposal{
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
		},
		// PROPOSAL_STATUS_DEPOSIT_FAILED
		{
			"TestSubmitTextProposalNotEnoughDepositFailed",
			proposer,
			gov.TextProposal{
				Title:    "TestSubmitTextProposalNotEnoughDepositFailed",
				Summary:  "TestSubmitTextProposalNotEnoughDepositFailed",
				Metadata: "TestSubmitTextProposalNotEnoughDepositFailed",
				Deposit:  "100" + utils.DisplayDenom,
			},
			nil,
			nil,
			gov.PROPOSAL_STATUS_DEPOSIT_FAILED,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			testSubmitProposal(t, test.Proposer, test.Proposal, test.Deposits, test.Votes, test.ExpectedStatus)
		})
	}
}

func testSubmitProposal(t *testing.T, proposer string, proposal any, deposits []Deposit, votes []Vote, expectedStatus string) {
	proposalId := gov.MustSubmitProposal(t, proposer, proposal)

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

	p := gov.MustQueryPassDepositPeriodProposal(t, proposalId)

	require.Equal(t, gov.PROPOSAL_STATUS_VOTING_PERIOD, p.Status)

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

	t.Parallel()

	p = gov.MustQueryPassVotingPeriodProposal(t, proposalId)

	require.Equal(t, expectedStatus, p.Status)
}
