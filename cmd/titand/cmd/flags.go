package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	ethermintsrvflags "github.com/tokenize-titan/ethermint/server/flags"
)

func UpdateFlags(cmd *cobra.Command) (*cobra.Command, error) {
	overrides := map[string]ethermintsrvflags.FlagOverride{
		flags.FlagFees: {
			Value: "",
			Usage: fmt.Sprintf("Fees to pay along with transaction; eg: 10%s", BaseDenom),
		},
		flags.FlagGasPrices: {
			Value: "",
			Usage: fmt.Sprintf("Gas prices to determine the transaction fee (e.g. 10%s)", BaseDenom),
		},
		flags.FlagKeyringBackend: {
			Usage: "Select keyring's backend (os|file|kwallet|pass|test|memory) (default in client.toml)",
		},
		genutilcli.FlagDefaultBondDenom: {
			Value: BaseDenom,
			Usage: fmt.Sprintf("genesis file default denomination, if left blank default value is '%s'", BaseDenom),
		},
	}

	return ethermintsrvflags.OverrideFlags(cmd, overrides)
}
