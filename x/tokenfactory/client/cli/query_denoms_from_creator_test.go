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

	"github.com/titantkx/titan/testutil/nullify"
	"github.com/titantkx/titan/testutil/sample"
	"github.com/titantkx/titan/x/tokenfactory/client/cli"
	"github.com/titantkx/titan/x/tokenfactory/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestQueryDenomsFromCreator(t *testing.T) {
	net, objs := networkWithFactoryDenomObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc      string
		idCreator string

		args []string
		err  error
		obj  []string
	}{
		{
			desc:      "found",
			idCreator: mustGetDenomCreator(t, objs[0].Denom),

			args: common,
			obj:  []string{objs[0].Denom},
		},
		{
			desc:      "not found",
			idCreator: sample.AccAddress().String(),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idCreator,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryDenomsFromCreator(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryDenomsFromCreatorResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Denoms)

				// Fix for G601: Implicit memory aliasing in for loop
				obj := tc.obj
				require.Equal(t,
					nullify.Fill(&obj),
					nullify.Fill(&resp.Denoms),
				)
			}
		})
	}
}

func mustGetDenomCreator(t *testing.T, denom string) string {
	t.Helper()
	creator, _, err := types.DeconstructDenom(denom)
	require.NoError(t, err)
	require.NotEmpty(t, creator)
	return creator
}
