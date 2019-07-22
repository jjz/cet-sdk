package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bankx Keeper will implement the interface
type ExpectedBankxKeeper interface {
	SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error // to tranfer coins
	FreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error                   // freeze some coins when creating orders
	UnFreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error                 // unfreeze coins and then orders can be executed
}

// Asset Keeper will implement the interface
type ExpectedAssetStatusKeeper interface {
	IsTokenExists(ctx sdk.Context, denom string) bool // check whether there is a coin named "denom"
	IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool
}