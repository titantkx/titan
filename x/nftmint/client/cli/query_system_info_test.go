package cli_test

import (
	"fmt"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"

	"github.com/tokenize-titan/titan/testutil/network"
	"github.com/tokenize-titan/titan/testutil/nullify"
	"github.com/tokenize-titan/titan/x/nftmint/client/cli"
	"github.com/tokenize-titan/titan/x/nftmint/types"
)

func networkWithSystemInfoObjects(t *testing.T) (*network.Network, types.SystemInfo) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	systemInfo := types.SystemInfo{}
	nullify.Fill(&systemInfo)
	state.SystemInfo = systemInfo
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.SystemInfo
}

func TestShowSystemInfo(t *testing.T) {
	net, obj := networkWithSystemInfoObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc string
		args []string
		err  error
		obj  types.SystemInfo
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowSystemInfo(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QuerySystemInfoResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.SystemInfo)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.SystemInfo),
				)
			}
		})
	}
}
