package nft

import (
	"encoding/base64"

	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	nftminttypes "github.com/titantkx/titan/x/nftmint/types"
)

type Class struct {
	id   string
	uri  string
	data string
}

func (class Class) GetID() string {
	return class.id
}

func (class Class) GetURI() string {
	return class.uri
}

func (class Class) GetData() string {
	return class.data
}

type NFT struct {
	classId string
	id      string
	uri     string
	data    string
}

func (nft NFT) GetClassID() string {
	return nft.classId
}

func (nft NFT) GetID() string {
	return nft.id
}

func (nft NFT) GetURI() string {
	return nft.uri
}

func (nft NFT) GetData() string {
	return nft.data
}

type Keeper struct {
	cdc codec.Codec
	nftkeeper.Keeper
}

func NewKeeper(cdc codec.Codec, k nftkeeper.Keeper) nfttransfertypes.NFTKeeper {
	return Keeper{cdc: cdc, Keeper: k}
}

func (k Keeper) CreateOrUpdateClass(ctx sdk.Context, classID, classURI, classData string) error {
	classDataRaw, err := base64.RawStdEncoding.DecodeString(classData)
	if err != nil {
		return err
	}
	class := nft.Class{
		Id:   classID,
		Uri:  classURI,
		Data: nftminttypes.MustNewAnyWithMetadata(string(classDataRaw)),
	}
	if !k.HasClass(ctx, classID) {
		return k.SaveClass(ctx, class)
	}
	return k.UpdateClass(ctx, class)
}

func (k Keeper) Mint(ctx sdk.Context, classID, tokenID, tokenURI, tokenData string, receiver sdk.AccAddress) error {
	tokenDataRaw, err := base64.RawStdEncoding.DecodeString(tokenData)
	if err != nil {
		return err
	}
	token := nft.NFT{
		ClassId: classID,
		Id:      tokenID,
		Uri:     tokenURI,
		Data:    nftminttypes.MustNewAnyWithMetadata(string(tokenDataRaw)),
	}
	return k.Keeper.Mint(ctx, token, receiver)
}

//nolint:revive	// keep tokenData for clear meaning of param
func (k Keeper) Transfer(ctx sdk.Context, classID, tokenID, tokenData string, receiver sdk.AccAddress) error {
	return k.Keeper.Transfer(ctx, classID, tokenID, receiver)
}

func (k Keeper) GetClass(ctx sdk.Context, classID string) (nfttransfertypes.Class, bool) {
	class, ok := k.Keeper.GetClass(ctx, classID)
	if !ok {
		return nil, false
	}
	classData := nftminttypes.MustGetMetadataFromAny(k.cdc, class.Data)
	return Class{
		id:   class.Id,
		uri:  class.Uri,
		data: base64.RawStdEncoding.EncodeToString([]byte(classData)),
	}, true
}

func (k Keeper) GetNFT(ctx sdk.Context, classID, tokenID string) (nfttransfertypes.NFT, bool) {
	token, ok := k.Keeper.GetNFT(ctx, classID, tokenID)
	if !ok {
		return nil, false
	}
	tokenData := nftminttypes.MustGetMetadataFromAny(k.cdc, token.Data)
	return NFT{
		classId: token.ClassId,
		id:      token.Id,
		uri:     token.Uri,
		data:    base64.RawStdEncoding.EncodeToString([]byte(tokenData)),
	}, true
}
