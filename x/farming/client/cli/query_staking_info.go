package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/titantkx/titan/x/farming/types"
)

func CmdListStakingInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-staking-info [token]",
		Short: "List all staking info of a token",
		Args:  cobra.RangeArgs(0, 1),
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

			var argToken string
			if len(args) > 0 {
				argToken = args[0]
			}

			params := &types.QueryStakingInfoAllRequest{
				Token:      argToken,
				Pagination: pageReq,
			}

			res, err := queryClient.StakingInfoAll(cmd.Context(), params)
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

func CmdShowStakingInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-staking-info [token] [staker]",
		Short: "Show staking info of a staker for a token",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argToken := args[0]
			argStaker := args[1]

			params := &types.QueryStakingInfoRequest{
				Token:  argToken,
				Staker: argStaker,
			}

			res, err := queryClient.StakingInfo(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
