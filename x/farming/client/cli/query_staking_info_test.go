//nolint:dupl
package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/titantkx/titan/testutil/network"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/farming/client/cli"
	"github.com/titantkx/titan/x/farming/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithStakingInfoObjects(t *testing.T, n int) (*network.Network, []types.StakingInfo) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		stakingInfo := types.StakingInfo{
			Token:  strconv.Itoa(i),
			Staker: strconv.Itoa(i),
			Amount: math.NewInt(int64(i)),
		}
		nullify.Fill(&stakingInfo)
		state.StakingInfoList = append(state.StakingInfoList, stakingInfo)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.StakingInfoList
}

func TestShowStakingInfo(t *testing.T) {
	net, objs := networkWithStakingInfoObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc     string
		idToken  string
		idStaker string

		args []string
		err  error
		obj  types.StakingInfo
	}{
		{
			desc:     "found",
			idToken:  objs[0].Token,
			idStaker: objs[0].Staker,

			args: common,
			obj:  objs[0],
		},
		{
			desc:     "not found",
			idToken:  strconv.Itoa(100000),
			idStaker: strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idToken,
				tc.idStaker,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowStakingInfo(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryStakingInfoResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.StakingInfo)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.StakingInfo),
				)
			}
		})
	}
}

func TestListStakingInfo(t *testing.T) {
	net, objs := networkWithStakingInfoObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListStakingInfo(), args)
			require.NoError(t, err)
			var resp types.QueryStakingInfoAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.StakingInfo), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.StakingInfo),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			//nolint:gosec // G115
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListStakingInfo(), args)
			require.NoError(t, err)
			var resp types.QueryStakingInfoAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.StakingInfo), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.StakingInfo),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListStakingInfo(), args)
		require.NoError(t, err)
		var resp types.QueryStakingInfoAllResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		//nolint:gosec // G115
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.StakingInfo),
		)
	})
}