package keeper

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

func (k Keeper) SetDataProviderAccumulatedReward(ctx sdk.Context, acc sdk.AccAddress, reward sdk.Coins) {
	key := oracletypes.DataProviderRewardsPrefixKey(acc)
	if !k.HasDataProviderReward(ctx, acc) {
		ctx.KVStore(k.storeKey).Set(key, k.cdc.MustMarshal(oracletypes.NewDataProviderAccumulatedReward(acc, reward)))
		return
	}
	oldReward := k.GetDataProviderAccumulatedReward(ctx, acc)
	newReward := oldReward.Add(reward...)
	ctx.KVStore(k.storeKey).Set(key, k.cdc.MustMarshal(oracletypes.NewDataProviderAccumulatedReward(acc, newReward)))
}

func (k Keeper) ClearDataProviderAccumulatedReward(ctx sdk.Context, acc sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(oracletypes.DataProviderRewardsPrefixKey(acc))
}

func (k Keeper) GetDataProviderAccumulatedReward(ctx sdk.Context, acc sdk.AccAddress) sdk.Coins {
	key := oracletypes.DataProviderRewardsPrefixKey(acc)
	bz := ctx.KVStore(k.storeKey).Get(key)
	dataProviderReward := oracletypes.DataProviderAccumulatedReward{}
	k.cdc.MustUnmarshal(bz, &dataProviderReward)
	return dataProviderReward.DataProviderReward
}

func (k Keeper) HasDataProviderReward(ctx sdk.Context, acc sdk.AccAddress) bool {
	return ctx.KVStore(k.storeKey).Has(oracletypes.DataProviderRewardsPrefixKey(acc))
}

// AllocateRewardsToDataProviders sends rewards from fee pool to data providers, that have given data for the passed request
func (k Keeper) AllocateRewardsToDataProviders(ctx sdk.Context, rid oracletypes.RequestID) {
	logger := k.Logger(ctx)
	request := k.MustGetRequest(ctx, rid)

	// rewards are lying in the distribution fee pool
	feePool := k.distrKeeper.GetFeePool(ctx)

	for _, rawReq := range request.RawRequests {
		ds := k.MustGetDataSource(ctx, rawReq.GetDataSourceID())

		ownerAccAddr, err := sdk.AccAddressFromBech32(ds.Owner)
		if err != nil {
			panic(err)
		}
		if !k.HasDataProviderReward(ctx, ownerAccAddr) {
			continue
		}
		reward := k.GetDataProviderAccumulatedReward(ctx, ownerAccAddr)

		diff, hasNeg := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(reward...))
		if hasNeg {
			logger.With("lack", diff).Error("oracle pool does not have enough coins to reward data providers")
			// not return because maybe still enough coins to pay someone
			continue
		}
		feePool.CommunityPool = diff

		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, distrtypes.ModuleName, ownerAccAddr, reward)
		if err != nil {
			panic(err)
		}

		// we are sure to have paid the reward to the provider, we can remove him now
		k.ClearDataProviderAccumulatedReward(ctx, ownerAccAddr)
	}

	k.distrKeeper.SetFeePool(ctx, feePool)
}
