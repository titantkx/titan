package cmd_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/tokenize-titan/titan/testutil/cmd/nft"
	"github.com/tokenize-titan/titan/testutil/sample"
	"github.com/tokenize-titan/titan/utils"
)

func MustCreateClass(t testing.TB, creator string) nft.Class {
	return nft.MustCreateClass(t, sample.URL(), sample.Hash(), sample.Word(), sample.Word(), sample.Paragraph(), sample.JSON(), creator)
}

func MustUpdateClass(t testing.TB, classId, owner string) {
	nft.MustUpdateClass(t, classId, sample.URL(), sample.Hash(), sample.Word(), sample.Word(), sample.Paragraph(), sample.JSON(), owner)
}

func MustErrUpdateClass(t testing.TB, errExpr, classId, owner string) {
	nft.MustErrUpdateClass(t, errExpr, classId, sample.URL(), sample.Hash(), sample.Word(), sample.Word(), sample.Paragraph(), sample.JSON(), owner)
}

func MustMint(t testing.TB, receiver, classId, minter string) nft.NFT {
	return nft.MustMint(t, receiver, classId, sample.URL(), sample.Hash(), sample.JSON(), minter)
}

func MustErrMint(t testing.TB, errExpr, receiver, classId, minter string) {
	nft.MustErrMint(t, errExpr, receiver, classId, sample.URL(), sample.Hash(), sample.JSON(), minter)
}

func TestCreateClass(t *testing.T) {
	t.Parallel()

	creator := MustCreateAccount(t, "1"+utils.DisplayDenom).Address

	MustCreateClass(t, creator)
}

func TestUpdateClass(t *testing.T) {
	t.Parallel()

	creator := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	class := MustCreateClass(t, creator)

	MustUpdateClass(t, class.Id, creator)
}

func TestUpdateClassNotFound(t *testing.T) {
	t.Parallel()

	updater := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	nonExistentClassId := strconv.FormatUint(math.MaxUint64, 10)

	MustErrUpdateClass(t, "class not found", nonExistentClassId, updater)
}

func TestUpdateClassUnauthorized(t *testing.T) {
	t.Parallel()

	creator := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	class := MustCreateClass(t, creator)
	unauthorizedUpdater := MustCreateAccount(t, "1"+utils.DisplayDenom).Address

	MustErrUpdateClass(t, "unauthorized", class.Id, unauthorizedUpdater)
}

func TestTransferClass(t *testing.T) {
	t.Parallel()

	creator := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	receiver := MustAddKey(t).Address
	class := MustCreateClass(t, creator)

	nft.MustTransferClass(t, class.Id, receiver, creator)
}

func TestTransferClassNotFound(t *testing.T) {
	t.Parallel()

	sender := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	receiver := MustAddKey(t).Address
	nonExistentClassId := strconv.FormatUint(math.MaxUint64, 10)

	nft.MustErrTransferClass(t, "class not found", nonExistentClassId, receiver, sender)
}

func TestTransferClassUnauthorized(t *testing.T) {
	t.Parallel()

	creator := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	receiver := MustAddKey(t).Address
	class := MustCreateClass(t, creator)
	unauthorizedSender := MustCreateAccount(t, "1"+utils.DisplayDenom).Address

	nft.MustErrTransferClass(t, "unauthorized", class.Id, receiver, unauthorizedSender)
}

func TestMint(t *testing.T) {
	t.Parallel()

	creator := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	receiver := MustAddKey(t).Address
	class := MustCreateClass(t, creator)

	MustMint(t, receiver, class.Id, creator)
}

func TestMintClassNotFound(t *testing.T) {
	t.Parallel()

	minter := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	receiver := MustAddKey(t).Address
	nonExistentClassId := strconv.FormatUint(math.MaxUint64, 10)

	MustErrMint(t, "class not found", receiver, nonExistentClassId, minter)
}

func TestMintUnauthorized(t *testing.T) {
	t.Parallel()

	creator := MustCreateAccount(t, "1"+utils.DisplayDenom).Address
	receiver := MustAddKey(t).Address
	class := MustCreateClass(t, creator)
	unauthorizedMinter := MustCreateAccount(t, "1"+utils.DisplayDenom).Address

	MustErrMint(t, "unauthorized", receiver, class.Id, unauthorizedMinter)
}
