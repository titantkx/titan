package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	etherminttests "github.com/titantkx/ethermint/tests"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/titantkx/titan/testutil/network"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/pointer/client/cli"
	"github.com/titantkx/titan/x/pointer/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithErc20NativeObjects(t *testing.T, n int) (*network.Network, []types.Erc20Native) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		erc20Native := types.Erc20Native{
			TokenDenom: strconv.Itoa(i),
			Erc20Addr:  etherminttests.GenerateAddress().String(),
		}
		nullify.Fill(&erc20Native)
		state.Erc20NativeList = append(state.Erc20NativeList, erc20Native)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.Erc20NativeList
}

func TestShowErc20Native(t *testing.T) {
	net, objs := networkWithErc20NativeObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc         string
		idTokenDenom string

		args []string
		err  error
		obj  types.Erc20Native
	}{
		{
			desc:         "found",
			idTokenDenom: objs[0].TokenDenom,

			args: common,
			obj:  objs[0],
		},
		{
			desc:         "not found",
			idTokenDenom: strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idTokenDenom,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowErc20Native(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryErc20NativeResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Erc20Native)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.Erc20Native),
				)
			}
		})
	}
}

func TestListErc20Native(t *testing.T) {
	net, objs := networkWithErc20NativeObjects(t, 5)

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
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListErc20Native(), args)
			require.NoError(t, err)
			var resp types.QueryErc20NativeAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Erc20Native), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Erc20Native),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListErc20Native(), args)
			require.NoError(t, err)
			var resp types.QueryErc20NativeAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Erc20Native), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Erc20Native),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListErc20Native(), args)
		require.NoError(t, err)
		var resp types.QueryErc20NativeAllResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total)) //nolint:gosec
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.Erc20Native),
		)
	})
}
