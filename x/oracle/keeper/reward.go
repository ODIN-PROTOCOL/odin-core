package keeper

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

func (k Keeper) SetDataProviderAccumulatedReward(ctx sdk.Context, acc sdk.AccAddress, reward sdk.Coins) error {
	hasReward, err := k.HasDataProviderReward(ctx, acc)
	if err != nil {
		return err
	}

	if !hasReward {
		return k.DataProviderAccumulatedRewards.Set(ctx, acc.Bytes(), oracletypes.NewDataProviderAccumulatedReward(acc, reward))
	}

	oldReward, err := k.GetDataProviderAccumulatedReward(ctx, acc.Bytes())
	if err != nil {
		return err
	}

	newReward := oldReward.Add(reward...)
	return k.DataProviderAccumulatedRewards.Set(ctx, acc.Bytes(), oracletypes.NewDataProviderAccumulatedReward(acc, newReward))
}

func (k Keeper) ClearDataProviderAccumulatedReward(ctx sdk.Context, acc sdk.AccAddress) error {
	return k.DataProviderAccumulatedRewards.Remove(ctx, acc.Bytes())
}

func (k Keeper) GetDataProviderAccumulatedReward(ctx sdk.Context, acc sdk.AccAddress) (sdk.Coins, error) {
	rewards, err := k.DataProviderAccumulatedRewards.Get(ctx, acc.Bytes())
	if err != nil {
		return nil, err
	}

	return rewards.DataProviderReward, nil
}

func (k Keeper) HasDataProviderReward(ctx sdk.Context, acc sdk.AccAddress) (bool, error) {
	return k.DataProviderAccumulatedRewards.Has(ctx, acc.Bytes())
}

// AllocateRewardsToDataProviders sends rewards from fee pool to data providers, that have given data for the passed request
func (k Keeper) AllocateRewardsToDataProviders(ctx sdk.Context, rid oracletypes.RequestID) error {
	request, err := k.GetRequest(ctx, rid)
	if err != nil {
		return err
	}

	for _, rawReq := range request.RawRequests {
		ds := k.MustGetDataSource(ctx, rawReq.GetDataSourceID())

		ownerAccAddr, err := sdk.AccAddressFromBech32(ds.Owner)
		if err != nil {
			return err
		}

		hasReward, err := k.HasDataProviderReward(ctx, ownerAccAddr)
		if err != nil {
			return err
		}

		if !hasReward {
			continue
		}

		reward, err := k.GetDataProviderAccumulatedReward(ctx, ownerAccAddr)
		if err != nil {
			return err
		}

		err = k.distrKeeper.DistributeFromFeePool(ctx, reward, ownerAccAddr)
		if err != nil {
			if err == distrtypes.ErrBadDistribution {
				k.Logger(ctx).Error("oracle pool does not have enough coins to reward data providers")
				continue
			}

			return err
		}

		// we are sure to have paid the reward to the provider, we can remove him now
		err = k.ClearDataProviderAccumulatedReward(ctx, ownerAccAddr)
		if err != nil {
			return err
		}
	}

	return nil
}
