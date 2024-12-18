package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/titantkx/titan/x/farming/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(_ string) *cobra.Command {
	// Group farming queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListFarm())
	cmd.AddCommand(CmdShowFarm())
	cmd.AddCommand(CmdListStakingInfo())
	cmd.AddCommand(CmdShowStakingInfo())
	cmd.AddCommand(CmdShowDistributionInfo())
	cmd.AddCommand(CmdListReward())
	cmd.AddCommand(CmdShowReward())
	// this line is used by starport scaffolding # 1

	return cmd
}
