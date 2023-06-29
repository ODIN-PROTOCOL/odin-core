package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
}

// BankKeeper defines the expected supply Keeper.
type BankKeeper interface {
	GetSupply(ctx sdk.Context, denom string) (supply sdk.Coin)
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

type DistrKeeper interface {
	GetFeePool(ctx sdk.Context) (feePool distrtypes.FeePool)
	SetFeePool(ctx sdk.Context, feePool distrtypes.FeePool)
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

type OracleKeeper interface {
	GetAccumulatedPaymentsForData(ctx sdk.Context) (payments oracletypes.AccumulatedPaymentsForData)
	SetAccumulatedPaymentsForData(ctx sdk.Context, payments oracletypes.AccumulatedPaymentsForData)
}
