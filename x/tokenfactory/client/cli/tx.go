package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/titantkx/titan/x/tokenfactory/types"
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
		Use:                        "tokenfactory",
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdBurn())
	cmd.AddCommand(CmdChangeAdmin())
	cmd.AddCommand(CmdCreateDenom())
	cmd.AddCommand(CmdMint())
	cmd.AddCommand(CmdSetDenomMetadata())
	// this line is used by starport scaffolding # 1

	return cmd
}
