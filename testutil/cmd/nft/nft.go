package nft

import (
	"context"
	"strings"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	"github.com/titantkx/titan/testutil/cmd"
	txcmd "github.com/titantkx/titan/testutil/cmd/tx"
)

type MintingInfo struct {
	ClassId     string       `json:"class_id"`
	Owner       string       `json:"owner"`
	NextTokenId testutil.Int `json:"next_token_id"`
}

func MustGetMintingInfo(t testutil.TestingT, classId string) MintingInfo {
	var v struct {
		MintingInfo MintingInfo `json:"minting_info"`
	}
	cmd.MustQuery(t, &v, "nft-mint", "show-minting-info", classId)
	require.Equal(t, classId, v.MintingInfo.ClassId)
	return v.MintingInfo
}

type Class struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
	Description string   `json:"description"`
	Uri         string   `json:"uri"`
	UriHash     string   `json:"uri_hash"`
	Data        Metadata `json:"data"`
}

type NFT struct {
	ClassId string   `json:"class_id"`
	Id      string   `json:"id"`
	Uri     string   `json:"uri"`
	UriHash string   `json:"uri_hash"`
	Data    Metadata `json:"data"`
}

type Metadata struct {
	Type string `json:"@type"`
	Data string `json:"data"`
}

func MustGetClass(t testutil.TestingT, classId string) Class {
	var v struct {
		Class Class `json:"class"`
	}
	cmd.MustQuery(t, &v, "nft", "class", classId)
	require.Equal(t, classId, v.Class.Id)
	return v.Class
}

func MustGetLatestClass(t testutil.TestingT) Class {
	var v struct {
		Classes []Class `json:"classes"`
	}
	cmd.MustQuery(t, &v, "nft", "classes", "--reverse")
	require.NotEmpty(t, v.Classes)
	return v.Classes[0]
}

func MustGetNFT(t testutil.TestingT, classId string, tokenId string) NFT {
	var v struct {
		NFT NFT `json:"nft"`
	}
	cmd.MustQuery(t, &v, "nft", "nft", classId, tokenId)
	require.Equal(t, classId, v.NFT.ClassId)
	require.Equal(t, tokenId, v.NFT.Id)
	return v.NFT
}

func MustCreateClass(t testutil.TestingT, uri, uriHash, name, symbol, description, data, creator string) Class {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	args := []string{
		"nft-mint",
		"create-class",
		uri,
		uriHash,
		"--from=" + creator,
	}
	if name != "" {
		args = append(args, "--class-name="+name)
	}
	if symbol != "" {
		args = append(args, "--class-symbol="+symbol)
	}
	if description != "" {
		args = append(args, "--class-description="+description)
	}
	if data != "" {
		args = append(args, "--class-data="+data)
	}

	tx := txcmd.MustExecTx(t, ctx, args...)

	classId := strings.Trim(tx.MustGetEventAttributeValue(t, "titan.nftmint.EventCreateClass", "id"), "\"")
	owner := strings.Trim(tx.MustGetEventAttributeValue(t, "titan.nftmint.EventCreateClass", "owner"), "\"")

	require.NotEmpty(t, classId)
	require.Equal(t, creator, owner)

	class := MustGetClass(t, classId)

	require.Equal(t, uri, class.Uri)
	require.Equal(t, uriHash, class.UriHash)
	require.Equal(t, name, class.Name)
	require.Equal(t, symbol, class.Symbol)
	require.Equal(t, description, class.Description)
	require.Equal(t, "/titan.nftmint.Metadata", class.Data.Type)
	require.Equal(t, data, class.Data.Data)

	mintingInfo := MustGetMintingInfo(t, classId)

	require.Equal(t, creator, mintingInfo.Owner)

	return class
}

func MustUpdateClass(t testutil.TestingT, classId, uri, uriHash, name, symbol, description, data, owner string) {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	args := []string{
		"nft-mint",
		"update-class",
		classId,
		uri,
		uriHash,
		"--from=" + owner,
	}
	if name != "" {
		args = append(args, "--class-name="+name)
	}
	if symbol != "" {
		args = append(args, "--class-symbol="+symbol)
	}
	if description != "" {
		args = append(args, "--class-description="+description)
	}
	if data != "" {
		args = append(args, "--class-data="+data)
	}

	tx := txcmd.MustExecTx(t, ctx, args...)

	updatedClassId := strings.Trim(tx.MustGetEventAttributeValue(t, "titan.nftmint.EventUpdateClass", "id"), "\"")

	require.Equal(t, classId, updatedClassId)

	class := MustGetClass(t, classId)

	require.Equal(t, uri, class.Uri)
	require.Equal(t, uriHash, class.UriHash)
	require.Equal(t, name, class.Name)
	require.Equal(t, symbol, class.Symbol)
	require.Equal(t, description, class.Description)
	require.Equal(t, "/titan.nftmint.Metadata", class.Data.Type)
	require.Equal(t, data, class.Data.Data)
}

func MustErrUpdateClass(t testutil.TestingT, errExpr, classId, uri, uriHash, name, symbol, description, data, owner string) {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	args := []string{
		"nft-mint",
		"update-class",
		classId,
		uri,
		uriHash,
		"--from=" + owner,
	}
	if name != "" {
		args = append(args, "--class-name="+name)
	}
	if symbol != "" {
		args = append(args, "--class-symbol="+symbol)
	}
	if description != "" {
		args = append(args, "--class-description="+description)
	}
	if data != "" {
		args = append(args, "--class-data="+data)
	}

	txcmd.MustErrExecTx(t, ctx, errExpr, args...)
}

func MustTransferClass(t testutil.TestingT, classId, receiver, sender string) {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "nft-mint", "transfer-class", classId, receiver, "--from="+sender)

	transferredClassId := strings.Trim(tx.MustGetEventAttributeValue(t, "titan.nftmint.EventTransferClass", "id"), "\"")
	oldOwner := strings.Trim(tx.MustGetEventAttributeValue(t, "titan.nftmint.EventTransferClass", "old_owner"), "\"")
	newOwner := strings.Trim(tx.MustGetEventAttributeValue(t, "titan.nftmint.EventTransferClass", "new_owner"), "\"")

	require.Equal(t, classId, transferredClassId)
	require.Equal(t, sender, oldOwner)
	require.Equal(t, receiver, newOwner)

	mintingInfo := MustGetMintingInfo(t, classId)

	require.Equal(t, receiver, mintingInfo.Owner)
}

func MustErrTransferClass(t testutil.TestingT, errExpr, classId, receiver, sender string) {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	txcmd.MustErrExecTx(t, ctx, errExpr, "nft-mint", "transfer-class", classId, receiver, "--from="+sender)
}

func MustMint(t testutil.TestingT, receiver, classId, uri, uriHash, data, minter string) NFT {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "nft-mint", "mint", receiver, classId, uri, uriHash, data, "--from="+minter)

	createdClassId := strings.Trim(tx.MustGetEventAttributeValue(t, "cosmos.nft.v1beta1.EventMint", "class_id"), "\"")
	tokenId := strings.Trim(tx.MustGetEventAttributeValue(t, "cosmos.nft.v1beta1.EventMint", "id"), "\"")
	owner := strings.Trim(tx.MustGetEventAttributeValue(t, "cosmos.nft.v1beta1.EventMint", "owner"), "\"")

	require.Equal(t, classId, createdClassId)
	require.NotEmpty(t, tokenId)
	require.Equal(t, receiver, owner)

	token := MustGetNFT(t, classId, tokenId)

	require.Equal(t, uri, token.Uri)
	require.Equal(t, uriHash, token.UriHash)
	require.Equal(t, "/titan.nftmint.Metadata", token.Data.Type)
	require.Equal(t, data, token.Data.Data)

	return token
}

func MustErrMint(t testutil.TestingT, errExpr, receiver, classId, uri, uriHash, data, minter string) {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	txcmd.MustErrExecTx(t, ctx, errExpr, "nft-mint", "mint", receiver, classId, uri, uriHash, data, "--from="+minter)
}
