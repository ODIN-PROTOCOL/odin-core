package oraclekeeper

import (
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type feeCollector struct {
	distrKeeper  oracletypes.DistrKeeper
	oracleKeeper Keeper
	payer        sdk.AccAddress
	collected    sdk.Coins
	limit        sdk.Coins
}

func (coll *feeCollector) Collect(ctx sdk.Context, coins sdk.Coins) error {
	coll.collected = coll.collected.Add(coins...)

	// If found any collected coin that exceed limit then return error
	for _, c := range coll.collected {
		limitAmt := coll.limit.AmountOf(c.Denom)
		if c.Amount.GT(limitAmt) {
			return sdkerrors.Wrapf(types.ErrNotEnoughFee, "require: %s, max: %s%s", c.String(), limitAmt.String(), c.Denom)
		}
	}

	// Actual send coins
	err := coll.distrKeeper.FundCommunityPool(ctx, coins, coll.payer)
	if err == nil {
		accumulatedPaymentsForData := coll.oracleKeeper.GetAccumulatedPaymentsForData(ctx)
		accumulatedPaymentsForData.AccumulatedAmount = accumulatedPaymentsForData.AccumulatedAmount.Add(coins...)

		coll.oracleKeeper.SetAccumulatedPaymentsForData(ctx, accumulatedPaymentsForData)
	}

	return err
}

func (coll *feeCollector) Collected() sdk.Coins {
	return coll.collected
}

func newFeeCollector(distrKeeper oracletypes.DistrKeeper, oracleKeeper Keeper, feeLimit sdk.Coins, payer sdk.AccAddress) FeeCollector {
	return &feeCollector{
		distrKeeper:  distrKeeper,
		oracleKeeper: oracleKeeper,
		payer:        payer,
		collected:    sdk.NewCoins(),
		limit:        feeLimit,
	}
}
