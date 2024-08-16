package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group tokenfactory queries under a subcommand
	cmd := &cobra.Command{
		Use:                        "tokenfactory",
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdQueryDenomAuthorityMetadata())
	cmd.AddCommand(CmdQueryDenomsFromCreator())
	cmd.AddCommand(CmdQueryAllBeforeSendHooks())
	cmd.AddCommand(CmdQueryDenomBeforeSendHook())

	return cmd
}
