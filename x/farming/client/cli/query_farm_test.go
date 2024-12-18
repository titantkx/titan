//nolint:dupl
package cli_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil/network"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/x/farming/client/cli"
	"github.com/titantkx/titan/x/farming/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithFarmObjects(t *testing.T, n int) (*network.Network, []types.Farm) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		farm := types.Farm{
			Token: strconv.Itoa(i),
			Rewards: []*types.FarmReward{
				{
					Sender:    sample.AccAddress().String(),
					Amount:    sdk.NewCoins(sdk.NewCoin("tkx", sdk.NewInt(1000))),
					EndTime:   time.Now().Add(1 * time.Hour),
					StartTime: time.Now(),
				},
			},
		}
		state.FarmList = append(state.FarmList, farm)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.FarmList
}

func TestShowFarm(t *testing.T) {
	net, objs := networkWithFarmObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc    string
		idToken string

		args []string
		err  error
		obj  types.Farm
	}{
		{
			desc:    "found",
			idToken: objs[0].Token,

			args: common,
			obj:  objs[0],
		},
		{
			desc:    "not found",
			idToken: strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idToken,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowFarm(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryFarmResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Farm)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.Farm),
				)
			}
		})
	}
}

func TestListFarm(t *testing.T) {
	net, objs := networkWithFarmObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			//nolint:gosec // G115
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListFarm(), args)
			require.NoError(t, err)
			var resp types.QueryFarmAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Farm), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Farm),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			//nolint:gosec // G115
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListFarm(), args)
			require.NoError(t, err)
			var resp types.QueryFarmAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Farm), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Farm),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListFarm(), args)
		require.NoError(t, err)
		var resp types.QueryFarmAllResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		//nolint:gosec // G115
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.Farm),
		)
	})
}
