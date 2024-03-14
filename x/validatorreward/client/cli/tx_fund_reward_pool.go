package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/titantkx/titan/x/validatorreward/types"
)

var _ = strconv.Itoa(0)

func CmdFundRewardPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund-reward-pool [amount]",
		Short: "Broadcast message FundRewardPool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAmount := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinsNormalized(argAmount)
			if err != nil {
				return err
			}

			msg := types.NewMsgFundRewardPool(
				clientCtx.GetFromAddress(),
				amount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
