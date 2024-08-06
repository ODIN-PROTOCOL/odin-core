package keeper

import (
	"context"
	"errors"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// valWithPower is an internal type to track validator with voting power inside of AllocateTokens.
type valWithPower struct {
	val   stakingtypes.ValidatorI
	power int64
}

// AllocateTokens allocates a portion of fee collected in the previous blocks to validators that
// that are actively performing oracle tasks. Note that this reward is also subjected to comm tax.
func (k Keeper) AllocateTokens(ctx sdk.Context, previousVotes []abci.VoteInfo) error {
	toReward := make([]valWithPower, 0)
	totalPower := int64(0)
	for _, vote := range previousVotes {
		val, err := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
		if err != nil {
			return err
		}

		valAddress, err := k.validatorAddressCodec.StringToBytes(val.GetOperator())
		status, err := k.GetValidatorStatus(ctx, valAddress)
		if err != nil {
			return err
		}

		if status.IsActive {
			toReward = append(toReward, valWithPower{val: val, power: vote.Validator.Power})
			totalPower += vote.Validator.Power
		}
	}
	if totalPower == 0 {
		// No active validators performing oracle tasks, nothing needs to be done here.
		return nil
	}

	feeCollector := k.AuthKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	distrModule := k.AuthKeeper.GetModuleAccount(ctx, distr.ModuleName)
	totalFee := sdk.NewDecCoinsFromCoins(k.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())...)

	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	// Compute the fee allocated for oracle module to distribute to active validators.
	oracleRewardRatio := math.LegacyNewDecWithPrec(int64(params.OracleRewardPercentage), 2)
	oracleRewardInt, _ := totalFee.MulDecTruncate(oracleRewardRatio).TruncateDecimal()
	// Transfer the oracle reward portion from fee collector to distr module.
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, distr.ModuleName, oracleRewardInt)
	if err != nil {
		return err
	}

	communityTax, err := k.distrKeeper.GetCommunityTax(ctx)
	if err != nil {
		return err
	}

	// Convert the transferred tokens back to DecCoins for internal distr allocations.
	oracleReward := sdk.NewDecCoinsFromCoins(oracleRewardInt...)
	remaining := oracleReward
	rewardMultiplier := math.LegacyOneDec().Sub(communityTax)
	// Allocate non-community pool tokens to active validators weighted by voting power.
	for _, each := range toReward {
		powerFraction := math.LegacyNewDec(each.power).QuoTruncate(math.LegacyNewDec(totalPower))
		reward := oracleReward.MulDecTruncate(rewardMultiplier).MulDecTruncate(powerFraction)
		err = k.distrKeeper.AllocateTokensToValidator(ctx, each.val, reward)
		if err != nil {
			return err
		}
		remaining = remaining.Sub(reward)
	}

	// Allocate the remaining coins to the community pool.
	// Recreate coins to sanitize them
	remainingNormalized := sdk.NewCoins(sdk.NormalizeCoins(remaining)...)
	if remainingNormalized.Empty() {
		return nil
	}

	return k.distrKeeper.FundCommunityPool(ctx, remainingNormalized, distrModule.GetAddress())

}

// GetValidatorStatus returns the validator status for the given validator. Note that validator
// status is default to [inactive, 0], so new validators start with inactive state.
func (k Keeper) GetValidatorStatus(ctx context.Context, val sdk.ValAddress) (types.ValidatorStatus, error) {
	validatorStatus, err := k.ValidatorStatuses.Get(ctx, val.Bytes())
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.ValidatorStatus{IsActive: false}, nil
		}
		return validatorStatus, err
	}

	return validatorStatus, nil
}

func (k Keeper) HasValidatorStatus(ctx context.Context, val sdk.ValAddress) (bool, error) {
	return k.ValidatorStatuses.Has(ctx, val.Bytes())
}

// SetValidatorStatus sets the validator status for the given validator.
func (k Keeper) SetValidatorStatus(ctx context.Context, val sdk.ValAddress, status types.ValidatorStatus) error {
	return k.ValidatorStatuses.Set(ctx, val.Bytes(), status)
}

// Activate changes the given validator's status to active. Returns error if the validator is
// already active or was deactivated recently, as specified by InactivePenaltyDuration parameter.
func (k Keeper) Activate(ctx context.Context, val sdk.ValAddress) error {
	goCtx := sdk.UnwrapSDKContext(ctx)

	hasValidator, err := k.HasValidatorStatus(ctx, val)
	if err != nil {
		return err
	}

	if hasValidator {
		status, err := k.GetValidatorStatus(ctx, val)
		if err != nil {
			return err
		}

		if status.IsActive {
			return types.ErrValidatorAlreadyActive
		}

		params, err := k.GetParams(ctx)
		if err != nil {
			return err
		}

		penaltyDuration := time.Duration(params.InactivePenaltyDuration)
		if !status.Since.IsZero() && status.Since.Add(penaltyDuration).After(goCtx.BlockHeader().Time) {
			return types.ErrTooSoonToActivate
		}
	}

	return k.SetValidatorStatus(ctx, val, types.NewValidatorStatus(true, goCtx.BlockHeader().Time))
}

// MissReport changes the given validator's status to inactive. No-op if already inactive or
// if the validator was active after the time the request happened.
func (k Keeper) MissReport(ctx context.Context, val sdk.ValAddress, requestTime time.Time) error {
	goCtx := sdk.UnwrapSDKContext(ctx)
	status, err := k.GetValidatorStatus(ctx, val)
	if err != nil {
		return err
	}

	if status.IsActive && status.Since.Before(requestTime) {
		err = k.SetValidatorStatus(ctx, val, types.NewValidatorStatus(false, goCtx.BlockHeader().Time))
		if err != nil {
			return err
		}

		goCtx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeDeactivate,
			sdk.NewAttribute(types.AttributeKeyValidator, val.String()),
		))
	}

	return nil
}
