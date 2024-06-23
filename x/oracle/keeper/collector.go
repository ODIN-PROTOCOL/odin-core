package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FeeCollector interface {
	Collect(sdk.Context, sdk.Coins) error
	Collected() sdk.Coins
}

type RewardCollector interface {
	Collect(sdk.Context, sdk.Coins, sdk.AccAddress) error
	CalculateReward([]byte, sdk.Coins) sdk.Coins
	Collected() sdk.Coins
}

// CollectReward subtract reward from fee pool and sends it to the data providers for reporting data
func (k Keeper) CollectReward(
	ctx sdk.Context, rawReports []oracletypes.RawReport, rawRequests []oracletypes.RawRequest,
) (sdk.Coins, error) {
	collector := newRewardCollector(k, k.BankKeeper)
	oracleParams, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	rawReportsMap := make(map[oracletypes.ExternalID]oracletypes.RawReport)
	for _, rawRep := range rawReports {
		rawReportsMap[rawRep.ExternalID] = rawRep
	}

	accumulatedDataProvidersRewards, err := k.GetAccumulatedDataProvidersRewards(ctx)
	if err != nil {
		return nil, err
	}

	accumulatedAmount := accumulatedDataProvidersRewards.AccumulatedAmount
	currentRewardPerByte := accumulatedDataProvidersRewards.CurrentRewardPerByte
	var rewPerByteInFeeDenom sdk.Coins

	for _, rawReq := range rawRequests {
		rawRep, ok := rawReportsMap[rawReq.GetExternalID()]
		if !ok {
			// this request had no report
			continue
		}

		ds := k.MustGetDataSource(ctx, rawReq.GetDataSourceID())
		dsOwnerAddr, err := sdk.AccAddressFromBech32(ds.Owner)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing data source owner address: %s", dsOwnerAddr)
		}

		for _, fee := range ds.Fee {
			rewPerByteInFeeDenom = rewPerByteInFeeDenom.Add(sdk.NewCoins(sdk.NewCoin(fee.Denom, currentRewardPerByte.AmountOf(fee.Denom)))...)
		}

		var reward sdk.Coins
		for {
			reward = collector.CalculateReward(rawRep.Data, rewPerByteInFeeDenom)
			if reward.Add(accumulatedAmount...).IsAllLT(oracleParams.DataProviderRewardThreshold.Amount) {
				break
			}

			rewPerByteInFeeDenom, _ = sdk.NewDecCoinsFromCoins(rewPerByteInFeeDenom...).MulDec(
				math.LegacyOneDec().Sub(oracleParams.RewardDecreasingFraction),
			).TruncateDecimal()
		}

		accumulatedAmount = accumulatedAmount.Add(reward...)
		err = collector.Collect(ctx, reward, dsOwnerAddr)
		if err != nil {
			return nil, err
		}
	}

	err = k.SetAccumulatedDataProvidersRewards(
		ctx,
		oracletypes.NewDataProvidersAccumulatedRewards(rewPerByteInFeeDenom, accumulatedAmount),
	)
	if err != nil {
		return nil, err
	}

	return collector.Collected(), nil
}
