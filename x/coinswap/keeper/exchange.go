package keeper

import (
	"strings"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	errortypes "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	coinswaptypes "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
)

// ExchangeDenom exchanges given amount
func (k Keeper) ExchangeDenom(
	ctx sdk.Context,
	fromDenom, toDenom string,
	amt sdk.Coin,
	requester sdk.AccAddress,
	additionalExchangeRates ...coinswaptypes.Exchange,
) error {
	// convert source amount to destination amount according to rate
	convertedAmt, err := k.convertToRate(ctx, fromDenom, toDenom, amt, additionalExchangeRates...)
	if err != nil {
		return errortypes.Wrap(err, "converting rate")
	}

	// first send source tokens to module
	err = k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(amt), requester)
	if err != nil {
		return errortypes.Wrapf(
			err,
			"sending coins from account: %s, to module: %s",
			requester.String(),
			distrtypes.ModuleName,
		)
	}

	toSend, remainder := convertedAmt.TruncateDecimal()
	if !remainder.IsZero() {
		k.Logger(ctx).With(
			"coins",
			remainder.String(),
		).Info("performing exchange according to limited precision some coins are lost")
	}

	feePool := k.distrKeeper.GetFeePool(ctx)

	diff, hasNeg := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(toSend))
	if hasNeg {
		k.Logger(ctx).With("lack", diff).Error("oracle pool does not have enough coins to reward data providers")
		// not return because maybe still enough coins to pay someone
		return sdkerrors.ErrInsufficientFunds
	}
	feePool.CommunityPool = diff
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, distrtypes.ModuleName, requester, sdk.NewCoins(toSend))
	if err != nil {
		return err
	}

	k.distrKeeper.SetFeePool(ctx, feePool)

	accumulatedPaymentsForData := k.oracleKeeper.GetAccumulatedPaymentsForData(ctx)

	accumulatedPaymentsForData.AccumulatedAmount, _ = accumulatedPaymentsForData.AccumulatedAmount.SafeSub(toSend)
	k.oracleKeeper.SetAccumulatedPaymentsForData(ctx, accumulatedPaymentsForData)

	return nil
}

// convertToRate returns the converted amount according to current rate
func (k Keeper) convertToRate(
	ctx sdk.Context,
	fromDenom, toDenom string,
	amt sdk.Coin,
	additionalExchangeRates ...coinswaptypes.Exchange,
) (sdk.DecCoin, error) {
	rate, err := k.GetRate(ctx, fromDenom, toDenom, additionalExchangeRates...)
	if err != nil {
		return sdk.DecCoin{}, errortypes.Wrap(err, "failed to convert to rate")
	}

	if rate.GT(sdk.NewDecFromInt(amt.Amount)) {
		return sdk.DecCoin{}, errortypes.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"current rate: %s is higher then amount provided: %s",
			rate.String(),
			amt.String(),
		)
	}

	convertedAmt := sdk.NewDecFromInt(amt.Amount).QuoRoundUp(rate)
	return sdk.NewDecCoinFromDec(toDenom, convertedAmt), nil
}

// GetRate returns the exchange rate for the given pair
func (k Keeper) GetRate(
	ctx sdk.Context,
	fromDenom, toDenom string,
	additionalExchangeRates ...coinswaptypes.Exchange,
) (sdk.Dec, error) {
	initialRate := k.GetInitialRate(ctx)
	params := k.GetParams(ctx)
	exchangeRates := append(params.ExchangeRates, additionalExchangeRates...)

	for _, ex := range exchangeRates {
		if strings.ToLower(ex.From) == strings.ToLower(fromDenom) && strings.ToLower(ex.To) == strings.ToLower(toDenom) {
			return initialRate.Mul(ex.RateMultiplier), nil
		}
	}

	return sdk.Dec{}, sdkerrors.Wrapf(
		coinswaptypes.ErrInvalidExchangeDenom,
		"failed to get rate from: %s, to: %s",
		fromDenom,
		toDenom,
	)
}
