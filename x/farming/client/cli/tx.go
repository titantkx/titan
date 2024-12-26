package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/titantkx/titan/x/farming/types"
)

var DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())

const (
	//nolint:unused
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	//nolint:unused
	listSeparator = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdAddReward())
	cmd.AddCommand(CmdStake())
	cmd.AddCommand(CmdUnstake())
	cmd.AddCommand(CmdHarvest())
	// this line is used by starport scaffolding # 1

	return cmd
}
