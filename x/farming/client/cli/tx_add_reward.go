package cli

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/titantkx/titan/x/farming/types"
)

var _ = strconv.Itoa(0)

func CmdAddReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-reward [token] [rewards] [end-time] [start-time]",
		Short: "Add farming rewards for a token",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argToken := args[0]

			argAmount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			argEndTime, err := time.Parse(time.RFC3339, args[2])
			if err != nil {
				return err
			}

			var argStartTime time.Time
			if len(args) > 3 {
				argStartTime, err = time.Parse(time.RFC3339, args[3])
				if err != nil {
					return err
				}
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddReward(
				clientCtx.GetFromAddress().String(),
				argToken,
				argAmount,
				argEndTime,
				argStartTime,
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
