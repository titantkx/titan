package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/titantkx/titan/x/nftmint/types"
)

var _ = strconv.Itoa(0)

func CmdTransferClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-class [class-id] [receiver] --from [sender] [flags]",
		Short: "Transfer a class to new owner",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argClassId := args[0]
			argReceiver := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferClass(
				clientCtx.GetFromAddress().String(),
				argClassId,
				argReceiver,
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
