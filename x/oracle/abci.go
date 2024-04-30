package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// handleBeginBlock re-calculates and saves the rolling seed value based on block hashes.
func handleBeginBlock(ctx sdk.Context, k keeper.Keeper) error {
	// Update rolling seed used for pseudorandom oracle provider selection.
	rollingSeed, err := k.GetRollingSeed(ctx)
	if err != nil {
		return err
	}

	err = k.SetRollingSeed(ctx, append(rollingSeed[1:], ctx.HeaderHash()[0]))
	if err != nil {
		return err
	}

	// Reward a portion of block rewards (inflation + tx fee) to active oracle validators.
	err = k.AllocateTokens(ctx, ctx.VoteInfos())
	if err != nil {
		return err
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	// Reset the price to the original price if a new boundary period has begun
	rewardThresholdBlocks := params.DataProviderRewardThreshold.Blocks
	previousBlock := uint64(ctx.BlockHeight() - 1)
	if previousBlock%rewardThresholdBlocks == 0 {
		initialReward := params.DataProviderRewardPerByte
		err = k.SetAccumulatedDataProvidersRewards(
			ctx,
			types.NewDataProvidersAccumulatedRewards(initialReward, sdk.NewCoins()),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// handleEndBlock cleans up the state during end block. See comment in the implementation!
func handleEndBlock(ctx sdk.Context, k keeper.Keeper) error {
	pendingResolveList, err := k.GetPendingResolveList(ctx)
	if err != nil {
		return err
	}

	// Loops through all requests in the resolvable list to resolve all of them!
	for _, reqID := range pendingResolveList {
		err = k.ResolveRequest(ctx, reqID)
		if err != nil {
			return err
		}

		err = k.AllocateRewardsToDataProviders(ctx, reqID)
		if err != nil {
			return err
		}
	}

	// Once all the requests are resolved, we can clear the list.
	err = k.SetPendingResolveList(ctx, []types.RequestID{})
	if err != nil {
		return err
	}

	// Lastly, we clean up data requests that are supposed to be expired.
	err = k.ProcessExpiredRequests(ctx)
	if err != nil {
		return err
	}
	// NOTE: We can remove old requests from state to optimize space, using `k.DeleteRequest`
	// and `k.DeleteReports`. We don't do that now as it is premature optimization at this state.

	return nil
}
