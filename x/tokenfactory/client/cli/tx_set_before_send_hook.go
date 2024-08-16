package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

var _ = strconv.Itoa(0)

func CmdSetBeforeSendHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-beforesend-hook [denom] [cosmwasm-address] --from [admin] [flags]",
		Short: "Set a cosmwasm contract to be the before send hook for a factory-created denom (must have admin authority to do so)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]

			argCosmWasmAddress := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetBeforeSendHook(
				clientCtx.GetFromAddress().String(),
				argDenom,
				argCosmWasmAddress,
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
