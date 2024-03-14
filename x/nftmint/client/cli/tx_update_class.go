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

func CmdUpdateClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-class [id] [uri] [uri-hash] --class-name [name] --class-symbol [symbol] --class-description [description] --class-data [data] --from [owner] [flags]",
		Short: "Update a class",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId := args[0]
			argUri := args[1]

			var argUriHash string
			if len(args) == 3 {
				argUriHash = args[2]
			}

			f := cmd.Flags()

			argName, err := f.GetString("class-name")
			if err != nil {
				return err
			}
			argSymbol, err := f.GetString("class-symbol")
			if err != nil {
				return err
			}
			argDescription, err := f.GetString("class-description")
			if err != nil {
				return err
			}
			argData, err := f.GetString("class-data")
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateClass(
				clientCtx.GetFromAddress().String(),
				argId,
				argName,
				argSymbol,
				argDescription,
				argUri,
				argUriHash,
				argData,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	f := cmd.Flags()
	f.String("class-name", "", "Class's name")
	f.String("class-symbol", "", "Class's symbol")
	f.String("class-description", "", "Class's description")
	f.String("class-data", "", "Class's data")

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
