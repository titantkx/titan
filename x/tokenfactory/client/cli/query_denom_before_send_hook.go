package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

func CmdQueryDenomBeforeSendHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denom-before-send-hook [denom]",
		Short: "Get the before send hook for a specific denom",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argDenom := args[0]

			params := &types.QueryBeforeSendHookAddressRequest{
				Denom: argDenom,
			}

			res, err := queryClient.BeforeSendHookAddress(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
