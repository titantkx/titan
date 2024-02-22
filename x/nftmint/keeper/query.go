package keeper

import (
	"github.com/tokenize-titan/titan/x/nftmint/types"
)

var _ types.QueryServer = Keeper{}
