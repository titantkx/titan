package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type MetadataI interface {
	GetData() string
}

func MustNewAnyWithMetadata(data string) *cdctypes.Any {
	if data == "" {
		return nil
	}
	v, err := cdctypes.NewAnyWithValue(&Metadata{Data: data})
	if err != nil {
		panic(err)
	}
	return v
}

func MustGetMetadataFromAny(cdc codec.Codec, v *cdctypes.Any) string {
	if v == nil {
		return ""
	}
	var data MetadataI
	if err := cdc.UnpackAny(v, &data); err != nil {
		panic(err)
	}
	return data.GetData()
}
