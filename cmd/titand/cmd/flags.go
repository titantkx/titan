package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"cosmossdk.io/tools/rosetta"
	"github.com/cosmos/cosmos-sdk/client/flags"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	ethermintsrvflags "github.com/tokenize-titan/ethermint/server/flags"
	"github.com/tokenize-titan/titan/utils"
)

func UpdateFlags(cmd *cobra.Command) (*cobra.Command, error) {
	overrides := map[string]ethermintsrvflags.FlagOverride{
		flags.FlagFees: {
			Usage: fmt.Sprintf("Fees to pay along with transaction; eg: 10%s", utils.BaseDenom),
		},
		flags.FlagGasPrices: {
			Usage: fmt.Sprintf("Gas prices to determine the transaction fee (e.g. 10%s)", utils.BaseDenom),
		},
		flags.FlagKeyringBackend: {
			Usage: "Select keyring's backend (os|file|kwallet|pass|test|memory) (default in client.toml)",
		},
		genutilcli.FlagDefaultBondDenom: {
			Value: utils.BaseDenom,
			Usage: fmt.Sprintf("genesis file default denomination, if left blank default value is '%s'", utils.BaseDenom),
		},
		rosetta.FlagDenomToSuggest: {
			Value: utils.BaseDenom,
		},
		rosetta.FlagPricesToSuggest: {
			Value: fmt.Sprintf("10%s", utils.BaseDenom),
		},
	}

	return ethermintsrvflags.OverrideFlags(cmd, overrides)
}
