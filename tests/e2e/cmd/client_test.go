package cmd_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/app/client"
	"github.com/titantkx/titan/testutil/cmd"
	"github.com/titantkx/titan/testutil/cmd/config"
	"github.com/titantkx/titan/utils"
)

var defaultClientCtx client.Context

var initClientCtxOnce sync.Once

func initClientContext(t testing.TB) {
	initClientCtxOnce.Do(func() {
		defaultClientCtx = MustCreateClientContext(t)
	})
}

func MustCreateClientContext(t testing.TB) client.Context {
	conf := config.MustGetConfig(t)
	clientCtx, err := client.CreateClientContext(client.Config{
		ChainID:        conf.ChainID,
		HomeDir:        cmd.HomeDir,
		KeyringDir:     cmd.HomeDir,
		KeyringBackend: conf.KeyringBackend,
		Output:         conf.Output,
		Node:           conf.Node,
		BroadcastMode:  conf.BroadcastMode,
	})
	require.NoError(t, err)
	return clientCtx
}

func MustSetFlag(t testing.TB, flags *pflag.FlagSet, name string, value string) {
	err := flags.Set(name, value)
	require.NoError(t, err)
}

func TestClient_BankSend(t *testing.T) {
	t.Parallel()

	initClientContext(t)

	sender := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	receiver := MustCreateAccount(t, "").Address
	amount := types.NewCoins(types.NewCoin(utils.BaseDenom, types.NewInt(1000)))

	msg := banktypes.NewMsgSend(
		types.MustAccAddressFromBech32(sender),
		types.MustAccAddressFromBech32(receiver),
		amount,
	)

	flags := client.NewTxFlags()
	MustSetFlag(t, flags, "from", sender)
	MustSetFlag(t, flags, "gas", "auto")
	MustSetFlag(t, flags, "gas-prices", "1000000000000atkx")
	MustSetFlag(t, flags, "gas-adjustment", "2")

	clientCtx, err := client.ReadTxFlags(defaultClientCtx, flags)
	require.NoError(t, err)

	clientCtx = clientCtx.WithDeadline(time.Now().Add(5 * time.Second))

	tx, err := client.BroadcastAndQueryTx(clientCtx, flags, msg)
	require.NoError(t, err)
	require.NotNil(t, tx)

	client := banktypes.NewQueryClient(defaultClientCtx)

	resp, err := client.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: receiver,
		Denom:   utils.BaseDenom,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, utils.BaseDenom, resp.Balance.Denom)
	require.Equal(t, int64(1000), resp.Balance.Amount.Int64())
}
