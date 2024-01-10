package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	keepertest "github.com/tokenize-titan/titan/testutil/keeper"
	"github.com/tokenize-titan/titan/utils"
)

func TestKeeper_LastDistributeTime(t *testing.T) {
	utils.InitSDKConfig()

	// return time.Time{} if the last distribute time is not set
	{
		keeper, ctx := keepertest.ValidatorrewardKeeper(t)
		value := keeper.GetLastDistributeTime(ctx)
		require.Equal(t, time.Time{}, value)
	}

	// Set the last distribute time in the store
	{
		keeper, ctx := keepertest.ValidatorrewardKeeper(t)
		// Set the last distribute time in the store
		lastDistributeTime := time.Now()
		keeper.SetLastDistributeTime(ctx, lastDistributeTime)

		// Retrieve the last distribute time
		value := keeper.GetLastDistributeTime(ctx)

		// Verify the retrieved value matches the set value
		require.Equal(t, lastDistributeTime.UTC(), value.UTC())
	}
}
