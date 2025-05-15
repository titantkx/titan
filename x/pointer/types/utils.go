package types

import (
	sdkerrors "cosmossdk.io/errors"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func GetErc20AddressFromString(addrStr string) (addr ethcommon.Address, err error) {
	mixCaseAddr, err := ethcommon.NewMixedcaseAddressFromString(addrStr)
	if err != nil {
		return ethcommon.Address{}, err
	}

	if !mixCaseAddr.ValidChecksum() {
		return ethcommon.Address{}, sdkerrors.Wrap(errortypes.ErrInvalidAddress, "invalid checksum for erc20 address")
	}
	addr = mixCaseAddr.Address()
	return
}
