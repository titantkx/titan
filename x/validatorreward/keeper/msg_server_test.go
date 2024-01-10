package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/tokenize-titan/titan/testutil/keeper"
	"github.com/tokenize-titan/titan/x/validatorreward/keeper"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, sdk.Context, *keeper.Keeper) {
	k, ctx := keepertest.ValidatorrewardKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx, k
}

func TestMsgServer(t *testing.T) {
	ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
