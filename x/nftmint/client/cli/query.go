package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/titantkx/titan/x/nftmint/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(_ string) *cobra.Command {
	// Group nftmint queries under a subcommand
	cmd := &cobra.Command{
		Use:                        "nft-mint",
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListMintingInfo())
	cmd.AddCommand(CmdShowMintingInfo())
	cmd.AddCommand(CmdShowSystemInfo())
	// this line is used by starport scaffolding # 1

	return cmd
}
