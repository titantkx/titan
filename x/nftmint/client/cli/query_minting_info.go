package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/tokenize-titan/titan/x/nftmint/types"
)

func CmdListMintingInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-minting-info",
		Short: "List all minting info",
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

			params := &types.QueryMintingInfosRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.MintingInfos(cmd.Context(), params)
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

func CmdShowMintingInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-minting-info [class-id]",
		Short: "Show minting info for a given class",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argClassId := args[0]

			params := &types.QueryMintingInfoRequest{
				ClassId: argClassId,
			}

			res, err := queryClient.MintingInfo(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
