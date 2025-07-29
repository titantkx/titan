package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/titantkx/titan/x/pointer/types"
)

var _ = strconv.Itoa(0)

func CmdAddERC20NativePointer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-erc-20-native-pointer [token] [name] [symbol] [decimals]",
		Short: "Broadcast message AddERC20NativePointer",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argToken := args[0]
			argName := args[1]
			argSymbol := args[2]
			argDecimals, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddERC20NativePointer(
				clientCtx.GetFromAddress().String(),
				argToken,
				argName,
				argSymbol,
				argDecimals,
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
