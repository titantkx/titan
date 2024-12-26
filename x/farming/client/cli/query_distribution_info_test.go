package cli_test

import (
	"fmt"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"

	"github.com/titantkx/titan/testutil/network"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/x/farming/client/cli"
	"github.com/titantkx/titan/x/farming/types"
)

func networkWithDistributionInfoObjects(t *testing.T) (*network.Network, types.DistributionInfo) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	distributionInfo := &types.DistributionInfo{}
	nullify.Fill(&distributionInfo)
	state.DistributionInfo = distributionInfo
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), *state.DistributionInfo
}

func TestShowDistributionInfo(t *testing.T) {
	net, obj := networkWithDistributionInfoObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc string
		args []string
		err  error
		obj  types.DistributionInfo
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowDistributionInfo(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryDistributionInfoResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.DistributionInfo)
			}
		})
	}
}
