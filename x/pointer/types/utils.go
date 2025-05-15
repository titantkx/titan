package types

import (
	sdkerrors "cosmossdk.io/errors"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func GetErc20AddressFromString(addrStr string) (*ethcommon.Address, error) {
	mixCaseAddr, err := ethcommon.NewMixedcaseAddressFromString(addrStr)
	if err != nil {
		return nil, err
	}

	if !mixCaseAddr.ValidChecksum() {
		return nil, sdkerrors.Wrap(errortypes.ErrInvalidAddress, "invalid checksum for erc20 address")
	}
	addr := mixCaseAddr.Address()
	return &addr, nil
}
