package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

var _ = strconv.Itoa(0)

func CmdSetAuthority() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-authority [new-authority]",
		Short: "Broadcast message setAuthority",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argNewAuthority := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetAuthority(
				clientCtx.GetFromAddress().String(),
				argNewAuthority,
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
