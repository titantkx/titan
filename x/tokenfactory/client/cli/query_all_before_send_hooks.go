package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func CmdQueryAllBeforeSendHooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-before-send-hooks",
		Short: "Returns a list of all before send hooks registered",
		Args:  cobra.NoArgs,
		//nolint:revive	// keep `args` for clear meaning
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllBeforeSendHooksAddressesRequest{}

			res, err := queryClient.AllBeforeSendHooksAddresses(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
