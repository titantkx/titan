package types_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/x/nftmint/types"
)

func TestMetadata(t *testing.T) {
	registry := cdctypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	types.RegisterInterfaces(registry)

	v := types.MustNewAnyWithMetadata("example metadata")
	data := types.MustGetMetadataFromAny(cdc, v)

	require.Equal(t, "example metadata", data)
}
