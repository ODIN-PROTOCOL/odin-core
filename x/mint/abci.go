package mint

import (
	"context"
	"time"

	"github.com/ODIN-PROTOCOL/odin-core/x/mint/keeper"
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx context.Context, k keeper.Keeper, ic minttypes.InflationCalculationFn) error {
	defer telemetry.ModuleMeasureSince(minttypes.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	goCtx := sdk.UnwrapSDKContext(ctx)
	goCtx.Logger().Info("here1")
	// fetch stored minter & params
	minter, err := k.GetMinter(ctx)
	if err != nil {
		return err
	}

	goCtx.Logger().Info("here2")
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	goCtx.Logger().Info("here3")
	// recalculate inflation rate
	totalStakingSupply, err := k.StakingTokenSupply(ctx)
	if err != nil {
		return err
	}

	goCtx.Logger().Info("here4")
	bondedRatio, err := k.BondedRatio(ctx)
	if err != nil {
		return err
	}

	goCtx.Logger().Info("here5")
	minter.Inflation = ic(ctx, minter, params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
	if err = k.SetMinter(ctx, minter); err != nil {
		return err
	}

	goCtx.Logger().Info("here6")
	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	goCtx.Logger().Info("here7")
	err = k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	goCtx.Logger().Info("here8")
	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(minttypes.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			minttypes.EventTypeMint,
			sdk.NewAttribute(minttypes.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(minttypes.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(minttypes.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)

	return nil
}
