package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/titantkx/titan/x/pointer/types"
)

func CmdListErc20Native() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-erc-20-native",
		Short: "list all erc20_native",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryErc20NativeAllRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Erc20NativeAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowErc20Native() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-erc-20-native [token-denom]",
		Short: "shows a erc20_native",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argTokenDenom := args[0]

			params := &types.QueryErc20NativeRequest{
				TokenDenom: argTokenDenom,
			}

			res, err := queryClient.Erc20Native(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
