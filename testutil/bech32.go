package testutil

import (
	"testing"

	"github.com/cosmos/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func MustAccountAddressToValidatorAddress(t testing.TB, accountAddr string) string {
	config := sdk.GetConfig()

	s, b, err := bech32.DecodeNoLimit(accountAddr)
	require.NoError(t, err)
	require.NotEmpty(t, b)
	require.Equal(t, config.GetBech32AccountAddrPrefix(), s)

	valAddr, err := bech32.Encode(config.GetBech32ValidatorAddrPrefix(), b)
	require.NoError(t, err)
	require.NotEmpty(t, valAddr)

	return valAddr
}
