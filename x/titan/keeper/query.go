package keeper

import (
	"github.com/titanlab/titan/x/titan/types"
)

var _ types.QueryServer = Keeper{}
