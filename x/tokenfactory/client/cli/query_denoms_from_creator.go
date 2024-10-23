package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func CmdQueryDenomsFromCreator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denoms-from-creator [creator]",
		Short: "Returns a list of all tokens created by a specific creator address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argCreator := args[0]

			params := &types.QueryDenomsFromCreatorRequest{
				Creator: argCreator,
			}

			res, err := queryClient.DenomsFromCreator(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
