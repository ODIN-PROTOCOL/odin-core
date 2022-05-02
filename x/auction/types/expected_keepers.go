package types

import (
	coinswaptypes "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OracleKeeper defines the expected oracle Keeper.
type OracleKeeper interface {
	GetAccumulatedPaymentsForData(ctx sdk.Context) (payments oracletypes.AccumulatedPaymentsForData)
}

// CoinswapKeeper defines the expected coinswap Keeper.
type CoinswapKeeper interface {
	ExchangeDenom(
		ctx sdk.Context,
		fromDenom, toDenom string,
		amt sdk.Coin,
		requester sdk.AccAddress,
		additionalExchanges ...coinswaptypes.Exchange,
	) error
}
