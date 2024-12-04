package gov

import (
	"context"
	"encoding/json"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
	txcmd "github.com/titantkx/titan/testutil/cmd/tx"
)

const (
	//nolint:revive
	PROPOSAL_STATUS_PASSED = "PROPOSAL_STATUS_PASSED"
	//nolint:revive
	PROPOSAL_STATUS_REJECTED = "PROPOSAL_STATUS_REJECTED"
	//nolint:revive
	PROPOSAL_STATUS_FAILED = "PROPOSAL_STATUS_FAILED"
	//nolint:revive
	PROPOSAL_STATUS_DEPOSIT_PERIOD = "PROPOSAL_STATUS_DEPOSIT_PERIOD"
	//nolint:revive
	PROPOSAL_STATUS_DEPOSIT_FAILED = "PROPOSAL_STATUS_DEPOSIT_FAILED"
	//nolint:revive
	PROPOSAL_STATUS_VOTING_PERIOD = "PROPOSAL_STATUS_VOTING_PERIOD"
)

const (
	//nolint:revive
	VOTE_OPTION_YES = "Yes"
	//nolint:revive
	VOTE_OPTION_NO = "No"
	//nolint:revive
	VOTE_OPTION_ABSTAIN = "Abstain"
	//nolint:revive
	VOTE_OPTION_NO_WITH_VETO = "NoWithVeto"
)

type Params struct {
	MinDeposit                 testutil.Coins    `json:"min_deposit"`
	MaxDepositPeriod           testutil.Duration `json:"max_deposit_period"`
	VotingPeriod               testutil.Duration `json:"voting_period"`
	Quorum                     testutil.Float    `json:"quorum"`
	Threshold                  testutil.Float    `json:"threshold"`
	VetoThreshold              testutil.Float    `json:"veto_threshold"`
	MinInitialDepositRatio     testutil.Float    `json:"min_initial_deposit_ratio"`
	BurnVoteQuorum             bool              `json:"burn_vote_quorum"`
	BurnProposalDepositPrevote bool              `json:"burn_proposal_deposit_prevote"`
	BurnVoteVeto               bool              `json:"burn_vote_veto"`
}

func MustGetParams(t testutil.TestingT) Params {
	var v struct {
		Params Params `json:"params"`
	}
	cmd.MustQuery(t, &v, "gov", "params")
	return v.Params
}

type Proposal struct {
	Id               string           `json:"id"`
	Status           string           `json:"status"`
	FinalTallyResult FinalTallyResult `json:"final_tally_result"`
	SubmitTime       time.Time        `json:"submit_time"`
	DepositEndTime   time.Time        `json:"deposit_end_time"`
	TotalDeposit     testutil.Coins   `json:"total_deposit"`
	VotingStartTime  time.Time        `json:"voting_start_time"`
	VotingEndTime    time.Time        `json:"voting_end_time"`
	Metadata         string           `json:"metadata"`
	Title            string           `json:"title"`
	Summary          string           `json:"summary"`
	Proposer         string           `json:"proposer"`
}

type FinalTallyResult struct {
	YesCount        testutil.Int `json:"yes_count"`
	AbstainCount    testutil.Int `json:"abstain_count"`
	NoCount         testutil.Int `json:"no_count"`
	NoWithVetoCount testutil.Int `json:"no_with_veto_count"`
}

func GetProposal(proposalId string) (*Proposal, error) {
	var proposal Proposal
	err := cmd.Query(&proposal, "gov", "proposal", proposalId)
	if err != nil {
		return nil, err
	}
	return &proposal, nil
}

func MustGetProposal(t testutil.TestingT, proposalId string) Proposal {
	var proposal Proposal
	cmd.MustQuery(t, &proposal, "gov", "proposal", proposalId)
	require.Equal(t, proposalId, proposal.Id)
	return proposal
}

func MustNotPassDepositPeriod(t testutil.TestingT, proposalId string) {
	for {
		proposal, err := GetProposal(proposalId)
		if err != nil {
			require.ErrorContains(t, err, "NotFound")
			return
		}
		require.NotNil(t, proposal)
		require.Equal(t, PROPOSAL_STATUS_DEPOSIT_PERIOD, proposal.Status)
		if proposal.DepositEndTime.After(time.Now()) {
			time.Sleep(time.Until(proposal.DepositEndTime) + 1*time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func MustQueryPassDepositPeriodProposal(t testutil.TestingT, proposalId string) Proposal {
	for {
		proposal := MustGetProposal(t, proposalId)
		if proposal.Status != PROPOSAL_STATUS_DEPOSIT_PERIOD {
			return proposal
		}
		if proposal.DepositEndTime.After(time.Now()) {
			time.Sleep(time.Until(proposal.DepositEndTime) + 1*time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func MustQueryPassVotingPeriodProposal(t testutil.TestingT, proposalId string) Proposal {
	for {
		proposal := MustGetProposal(t, proposalId)
		if proposal.Status != PROPOSAL_STATUS_DEPOSIT_PERIOD && proposal.Status != PROPOSAL_STATUS_VOTING_PERIOD {
			return proposal
		}
		if proposal.VotingEndTime.After(time.Now()) {
			time.Sleep(time.Until(proposal.VotingEndTime) + 1*time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

type ProposalMsg struct {
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Metadata string `json:"metadata"`
	Deposit  string `json:"deposit"`
	Messages []any  `json:"messages,omitempty"`
}

type MsgUpdateParams struct {
	Type      string `json:"@type"`
	Authority string `json:"authority"`
	Params    any    `json:"params"`
}

type MsgSoftwareUpgrade struct {
	Type      string              `json:"@type"`
	Authority string              `json:"authority"`
	Plan      SoftwareUpgradePlan `json:"plan"`
}

type SoftwareUpgradePlan struct {
	Name   string       `json:"name"`
	Height testutil.Int `json:"height"`
	Info   string       `json:"info"`
}

type MsgNftCreateClass struct {
	Type        string `json:"@type"`
	Creator     string `json:"creator,omitempty"`
	Name        string `json:"name,omitempty"`
	Symbol      string `json:"symbol,omitempty"`
	Description string `json:"description,omitempty"`
	Uri         string `json:"uri,omitempty"`
	UriHash     string `json:"uri_hash,omitempty"`
	Data        string `json:"data,omitempty"`
}

type MsgValidatorRewardSetRate struct {
	Type      string `json:"@type"`
	Authority string `json:"authority,omitempty"`
	Rate      string `json:"rate"`
}

func MustSubmitProposal(t testutil.TestingT, from string, proposal ProposalMsg) string {
	file := testutil.MustCreateTemp(t, "proposal_*.json")
	err := json.NewEncoder(file).Encode(proposal)

	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "gov", "submit-proposal", file.Name(), "--from="+from)

	proposalId := tx.MustGetEventAttributeValue(t, "submit_proposal", "proposal_id")

	return proposalId
}

func MustDeposit(t testutil.TestingT, from string, proposalId string, amount string) {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()
	txcmd.MustExecTx(t, ctx, "gov", "deposit", proposalId, amount, "--from="+from)
}

func MustVote(t testutil.TestingT, from string, proposalId string, option string) {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()
	txcmd.MustExecTx(t, ctx, "gov", "vote", proposalId, option, "--from="+from)
}
