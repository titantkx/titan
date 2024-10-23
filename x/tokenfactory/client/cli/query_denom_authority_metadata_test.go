//nolint:dupl
package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/titantkx/titan/testutil/network"
	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/x/tokenfactory/client/cli"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithFactoryDenomObjects(t *testing.T, n int) (*network.Network, []types.GenesisDenom) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		factoryDenom := types.GenesisDenom{
			Denom: fmt.Sprintf("factory/%s/%d", sample.AccAddress(), i),
			AuthorityMetadata: types.DenomAuthorityMetadata{
				Admin: sample.AccAddress().String(),
			},
		}
		state.FactoryDenoms = append(state.FactoryDenoms, factoryDenom)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.FactoryDenoms
}

func TestQueryDenomAuthorityMetadata(t *testing.T) {
	net, objs := networkWithFactoryDenomObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc    string
		idDenom string

		args []string
		err  error
		obj  types.DenomAuthorityMetadata
	}{
		{
			desc:    "found",
			idDenom: objs[0].Denom,

			args: common,
			obj:  objs[0].AuthorityMetadata,
		},
		{
			desc:    "not found",
			idDenom: fmt.Sprintf("factory/%s/%d", sample.AccAddress(), 100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idDenom,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryDenomAuthorityMetadata(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryDenomAuthorityMetadataResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.AuthorityMetadata)

				// Fix for G601: Implicit memory aliasing in for loop
				obj := tc.obj
				require.Equal(t,
					nullify.Fill(&obj),
					nullify.Fill(&resp.AuthorityMetadata),
				)
			}
		})
	}
}
