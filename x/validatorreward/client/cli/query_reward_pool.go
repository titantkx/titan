package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/titantkx/titan/x/validatorreward/types"
)

var _ = strconv.Itoa(0)

func CmdRewardPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-pool",
		Short: "Query reward-pool",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) { //nolint:revive	// keep args for clear meaning
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRewardPoolRequest{}

			res, err := queryClient.RewardPool(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
