package keeper_test

import (
	"reflect"
	"testing"

	testkeeper "github.com/tokenize-titan/titan/testutil/keeper"
	"github.com/tokenize-titan/titan/utils"
	"github.com/tokenize-titan/titan/x/validatorreward/types"
)

func TestGetParams(t *testing.T) {
	utils.InitSDKConfig()

	k, ctx := testkeeper.ValidatorrewardKeeper(t)
	params := types.DefaultParams()

	if !reflect.DeepEqual(params, k.GetParams(ctx)) {
		t.Errorf("GetParams() = %v, want %v", k.GetParams(ctx), params)
	}
}
