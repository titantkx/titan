package keeper_test

import (
	"reflect"
	"testing"

	testkeeper "github.com/titantkx/titan/testutil/keeper"
	"github.com/titantkx/titan/utils"
	"github.com/titantkx/titan/x/validatorreward/types"
)

func TestGetParams(t *testing.T) {
	utils.InitSDKConfig()

	k, ctx := testkeeper.ValidatorrewardKeeper(t)
	params := types.DefaultParams()

	if !reflect.DeepEqual(params, k.GetParams(ctx)) {
		t.Errorf("GetParams() = %v, want %v", k.GetParams(ctx), params)
	}
}
