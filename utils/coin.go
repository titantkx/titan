package utils

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Convert sdk.Coins to wasmvmtypes.Coins
func CWCoinsFromSDKCoins(in sdk.Coins) wasmvmtypes.Coins {
	var cwCoins wasmvmtypes.Coins
	for _, coin := range in {
		cwCoins = append(cwCoins, CWCoinFromSDKCoin(coin))
	}
	return cwCoins
}

// Convert sdk.Coin to wasmvmtypes.Coin
func CWCoinFromSDKCoin(in sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom:  in.GetDenom(),
		Amount: in.Amount.String(),
	}
}
