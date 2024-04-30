package keeper

import (
	"cosmossdk.io/errors"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type feeCollector struct {
	distrKeeper  types.DistrKeeper
	oracleKeeper Keeper
	payer        sdk.AccAddress
	collected    sdk.Coins
	limit        sdk.Coins
}

func newFeeCollector(distrKeeper types.DistrKeeper, oracleKeeper Keeper, feeLimit sdk.Coins, payer sdk.AccAddress) FeeCollector {
	return &feeCollector{
		distrKeeper:  distrKeeper,
		oracleKeeper: oracleKeeper,
		payer:        payer,
		collected:    sdk.NewCoins(),
		limit:        feeLimit,
	}
}

func (coll *feeCollector) Collect(ctx sdk.Context, coins sdk.Coins) error {
	coll.collected = coll.collected.Add(coins...)

	// If found any collected coin that exceed limit then return error
	for _, c := range coll.collected {
		limitAmt := coll.limit.AmountOf(c.Denom)
		if c.Amount.GT(limitAmt) {
			return errors.Wrapf(
				types.ErrNotEnoughFee,
				"require: %s, max: %s%s",
				c.String(),
				limitAmt.String(),
				c.Denom,
			)
		}
	}

	// Actual send coins
	err := coll.distrKeeper.FundCommunityPool(ctx, coins, coll.payer)
	if err == nil {
		accumulatedPaymentsForData, err := coll.oracleKeeper.GetAccumulatedPaymentsForData(ctx)
		if err != nil {
			return err
		}

		accumulatedPaymentsForData.AccumulatedAmount = accumulatedPaymentsForData.AccumulatedAmount.Add(coins...)

		err = coll.oracleKeeper.SetAccumulatedPaymentsForData(ctx, accumulatedPaymentsForData)
		if err != nil {
			return err
		}
	}

	return err
}

func (coll *feeCollector) Collected() sdk.Coins {
	return coll.collected
}
