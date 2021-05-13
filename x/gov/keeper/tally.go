package keeper

import (
	govtypes "github.com/GeoDB-Limited/odin-core/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (keeper Keeper) Tally(ctx sdk.Context, proposal govtypes.Proposal) (passes bool, burnDeposits bool, tallyResults govtypes.TallyResult) {
	results := make(map[govtypes.VoteOption]sdk.Dec)
	results[govtypes.OptionYes] = sdk.ZeroDec()
	results[govtypes.OptionAbstain] = sdk.ZeroDec()
	results[govtypes.OptionNo] = sdk.ZeroDec()
	results[govtypes.OptionNoWithVeto] = sdk.ZeroDec()

	totalVotingPower := sdk.ZeroDec()

	keeper.IterateVotes(ctx, proposal.ProposalId, func(vote govtypes.Vote) bool {
		// if validator, just record it in the map
		voter, err := sdk.AccAddressFromBech32(vote.Voter)

		if err != nil {
			panic(err)
		}

		votingPower := keeper.bankKeeper.SpendableCoins(ctx, voter).AmountOf(keeper.stakingKeeper.BondDenom(ctx))

		results[vote.Option] = results[vote.Option].Add(votingPower.ToDec())
		totalVotingPower = totalVotingPower.Add(votingPower.ToDec())

		keeper.deleteVote(ctx, vote.ProposalId, voter)
		return false
	})

	tallyParams := keeper.GetTallyParams(ctx)
	tallyResults = govtypes.NewTallyResultFromMap(results)

	// If there is not enough quorum of votes, the proposal fails
	bondDenom := keeper.stakingKeeper.BondDenom(ctx)
	// total active supply = total supply - treasury pool (treasury pool does not take part in governance)
	totalActiveSupply := keeper.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(bondDenom).ToDec().Sub(keeper.mintKeeper.GetMintPool(ctx).TreasuryPool.AmountOf(bondDenom).ToDec())
	percentVoting := totalVotingPower.Quo(totalActiveSupply)
	if percentVoting.LT(tallyParams.Quorum) {
		return false, true, tallyResults
	}

	// If no one votes (everyone abstains), proposal fails
	if totalVotingPower.Sub(results[govtypes.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false, tallyResults
	}

	// If more than 1/3 of voters veto, proposal fails
	if results[govtypes.OptionNoWithVeto].Quo(totalVotingPower).GT(tallyParams.VetoThreshold) {
		return false, true, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[govtypes.OptionYes].Quo(totalVotingPower.Sub(results[govtypes.OptionAbstain])).GT(tallyParams.Threshold) {
		return true, false, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails
	return false, false, tallyResults
}
