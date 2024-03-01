package setup

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/testutil/cmd/gov"
	"github.com/tokenize-titan/titan/testutil/cmd/status"
)

const ImageName = "titanlab/titand"

func Install(t testutil.TestingT, w io.Writer, rootDir string) {
	cwd := testutil.Getwd(t)
	testutil.Chdir(t, rootDir)
	defer testutil.Chdir(t, cwd)
	cmd.MustExecWrite(t, w, "make", "build")
	cmd.MustExecWrite(t, w, "cp", rootDir+"/build/titand", testutil.UserHomeDir(t)+"/go/bin")
}

func BuildImage(t testutil.TestingT, w io.Writer, rootDir string, tag string) {
	cmd.MustExecWrite(t, w, "docker", "build", rootDir, "-t", ImageName+":"+tag)
}

func StartChain(t testutil.TestingT, w io.Writer, dcFile string) (ready chan struct{}, done chan struct{}) {
	readyCh := make(chan struct{})
	doneCh := make(chan struct{})

	s := testutil.NewStreamer()

	go func() {
		cmd.MustExecWrite(t, s, "docker", "compose", "-f", dcFile, "up")
		s.Close()
		doneCh <- struct{}{}
	}()

	go func() {
		isRunning := false
		r := bufio.NewReader(s)
		for {
			line, isPrefix, err := r.ReadLine()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
			w.Write(line)
			if !isPrefix {
				fmt.Fprintln(w)
			}
			if !isRunning && strings.Contains(string(line), "executed block") {
				isRunning = true
				readyCh <- struct{}{}
			}
		}
	}()

	return readyCh, doneCh
}

func StartChainAndListenForUpgrade(t testutil.TestingT, w io.Writer, dcFile string, upgradeName string) (ready <-chan struct{}, upgrade <-chan struct{}, done <-chan struct{}) {
	upgradeCh := make(chan struct{})

	s := testutil.NewStreamer()

	readyCh, doneCh := StartChain(t, s, dcFile)

	go func() {
		needUpgrade := false
		r := bufio.NewReader(s)
		for {
			line, isPrefix, err := r.ReadLine()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
			w.Write(line)
			if !isPrefix {
				fmt.Fprintln(w)
			}
			if !needUpgrade && strings.Contains(string(line), `UPGRADE "`+upgradeName+`" NEEDED`) {
				needUpgrade = true
				upgradeCh <- struct{}{}
			}
		}
	}()

	return readyCh, upgradeCh, doneCh
}

func StopChain(t testutil.TestingT, w io.Writer, dcFile string) {
	cmd.MustExecWrite(t, w, "docker", "compose", "-f", dcFile, "down")
}

func UpgradeChain(t testutil.TestingT, upgradeName string, vals ...string) {
	require.NotEmpty(t, vals)

	curHeight := status.MustGetLatestBlockHeight(t)
	depositAmount := gov.MustGetParams(t).MinDeposit.String()

	proposalId := gov.MustSubmitProposal(t, vals[0], gov.ProposalMsg{
		Title:    "Upgrade titan chain",
		Summary:  "Upgrade titan chain",
		Metadata: "Upgrade titan chain",
		Deposit:  depositAmount,
		Messages: []any{
			gov.MsgSoftwareUpgrade{
				Type:      "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
				Authority: "titan10d07y265gmmuvt4z0w9aw880jnsr700jste397",
				Plan: gov.SoftwareUpgradePlan{
					Name:   upgradeName,
					Height: testutil.MakeInt(curHeight + 50),
					Info:   "https://example.com/titand-info.json",
				},
			},
		},
	})

	for _, val := range vals {
		gov.MustVote(t, val, proposalId, gov.VOTE_OPTION_YES)
	}

	proposal := gov.MustQueryPassVotingPeriodProposal(t, proposalId)

	require.Equal(t, gov.PROPOSAL_STATUS_PASSED, proposal.Status)
}
