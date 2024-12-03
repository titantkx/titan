package setup

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
	"github.com/titantkx/titan/testutil/cmd/gov"
	"github.com/titantkx/titan/testutil/cmd/status"
)

const ImageName = "titantkx/titand"

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
		numberExecutedBlocksNeeded := 6
		for {
			line, isPrefix, err := r.ReadLine()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
			if !isRunning {
				fmt.Println(string(line))
			}
			//nolint:errcheck	// accept the error
			w.Write(line)
			if !isPrefix {
				fmt.Fprintln(w)
			}
			if !isRunning && strings.Contains(string(line), "executed block") {
				fmt.Printf("Executed block %d\n", numberExecutedBlocksNeeded)
				numberExecutedBlocksNeeded--
				if numberExecutedBlocksNeeded == 0 {
					isRunning = true
					readyCh <- struct{}{}
				}
			}
		}
	}()

	return readyCh, doneCh
}

func StartChainAndListenForUpgrade(t testutil.TestingT, w io.Writer, dcFile string, upgradeName string) (ready <-chan struct{}, upgrade <-chan struct{}, done <-chan struct{}) {
	upgradeCh := make(chan struct{})

	s := testutil.NewStreamer()

	readyCh, doneCh := StartChain(t, s, dcFile)
	needUpgrade := false

	go func() {
		r := bufio.NewReader(s)

		for {
			line, isPrefix, err := r.ReadLine()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
			//nolint:errcheck	// accept the error
			w.Write(line)
			if !isPrefix {
				fmt.Fprintln(w)
			}
			if !needUpgrade && strings.Contains(string(line), `UPGRADE "`+upgradeName+`" NEEDED`) {
				needUpgrade = true
				go func() {
					// wait for 15 seconds before sending the upgrade signal
					// to make sure all node already go to need upgrade state
					time.Sleep(15 * time.Second)
					upgradeCh <- struct{}{}
				}()
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
