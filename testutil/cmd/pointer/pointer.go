package pointer

import (
	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
)

type Erc20Native struct {
	TokenDenom string `json:"token_denom"`
	Erc20Addr  string `json:"erc20_addr"`
}

func MustGetERC20Pointer(t testutil.TestingT, tokenDenom string) string {
	var v struct {
		Erc20Native Erc20Native `json:"erc20_native"`
	}
	cmd.MustQuery(t, &v, "pointer", "show-erc-20-native", tokenDenom)
	require.Equal(t, tokenDenom, v.Erc20Native.TokenDenom)
	return v.Erc20Native.Erc20Addr
}
