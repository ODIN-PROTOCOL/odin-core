package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"

	"github.com/ODIN-PROTOCOL/odin-core/testing/testapp"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

func defaultVotes() []abci.VoteInfo {
	return []abci.VoteInfo{{
		Validator: abci.Validator{
			Address: testapp.Validators[0].PubKey.Address(),
			Power:   70,
		},
	}, {
		Validator: abci.Validator{
			Address: testapp.Validators[1].PubKey.Address(),
			Power:   20,
		},
	}, {
		Validator: abci.Validator{
			Address: testapp.Validators[2].PubKey.Address(),
			Power:   10,
		},
	}}
}

func SetupFeeCollector(app *testapp.TestingApp, ctx sdk.Context, k keeper.Keeper) sdk.ModuleAccountI {
	// Set collected fee to 1000000loki and 70% oracle reward proportion.
	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, Coins1000000loki)
	if err != nil {
		panic(err)
	}
	err = app.BankKeeper.SendCoinsFromModuleToModule(
		ctx,
		minttypes.ModuleName,
		authtypes.FeeCollectorName,
		Coins1000000loki,
	)
	if err != nil {
		panic(err)
	}
	app.AccountKeeper.SetAccount(ctx, feeCollector)

	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	params.OracleRewardPercentage = 70
	err = k.SetParams(ctx, params)
	if err != nil {
		panic(err)
	}

	return feeCollector
}

func TestAllocateTokenNoActiveValidators(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false, false, true)
	feeCollector := SetupFeeCollector(app, ctx, k)

	require.Equal(t, Coins1000000loki, app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()))
	// No active oracle validators so nothing should happen.
	err := k.AllocateTokens(ctx, defaultVotes())
	require.NoError(t, err)

	distAccount := app.AccountKeeper.GetModuleAccount(ctx, distrtypes.ModuleName)
	require.Equal(t, Coins1000000loki, app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()))
	require.Empty(t, app.BankKeeper.GetAllBalances(ctx, distAccount.GetAddress()))
}

func TestAllocateTokensOneActive(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false, false, true)
	feeCollector := SetupFeeCollector(app, ctx, k)

	require.Equal(t, Coins1000000loki, app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()))
	// From 70% of fee, 2% should go to community pool, the rest goes to the only active validator.
	err := k.Activate(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	err = k.AllocateTokens(ctx, defaultVotes())
	require.NoError(t, err)

	distAccount := app.AccountKeeper.GetModuleAccount(ctx, distrtypes.ModuleName)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 300000)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 700000)),
		app.BankKeeper.GetAllBalances(ctx, distAccount.GetAddress()),
	)

	communityPool, err := distrkeeper.NewQuerier(app.DistrKeeper).CommunityPool(
		ctx,
		&distrtypes.QueryCommunityPoolRequest{},
	)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(14000)}},
		communityPool.Pool,
	)

	validatorOutstandingRewards, err := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Empty(t, validatorOutstandingRewards)

	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(686000)}},
		validatorOutstandingRewards.Rewards,
	)

	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[2].ValAddress)
	require.NoError(t, err)
	require.Empty(t, validatorOutstandingRewards)
}

func TestAllocateTokensAllActive(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(true, false, true)
	feeCollector := SetupFeeCollector(app, ctx, k)

	require.Equal(t, Coins1000000loki, app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()))
	// From 70% of fee, 2% should go to community pool, the rest get split to validators.
	err := k.AllocateTokens(ctx, defaultVotes())
	require.NoError(t, err)

	distAccount := app.AccountKeeper.GetModuleAccount(ctx, distrtypes.ModuleName)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 300000)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 700000)),
		app.BankKeeper.GetAllBalances(ctx, distAccount.GetAddress()),
	)

	communityPool, err := distrkeeper.NewQuerier(app.DistrKeeper).CommunityPool(
		ctx,
		&distrtypes.QueryCommunityPoolRequest{},
	)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(14000)}},
		communityPool.Pool,
	)

	validatorOutstandingRewards, err := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(480200)}},
		validatorOutstandingRewards.Rewards,
	)

	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(137200)}},
		validatorOutstandingRewards.Rewards,
	)

	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[2].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(68600)}},
		validatorOutstandingRewards.Rewards,
	)
}

func TestGetDefaultValidatorStatus(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	vs, err := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(t, types.NewValidatorStatus(false, time.Time{}), vs)
}

func TestGetSetValidatorStatus(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	// After setting status of the 1st validator, we should be able to get it back.
	err := k.SetValidatorStatus(ctx, testapp.Validators[0].ValAddress, types.NewValidatorStatus(true, now))
	require.NoError(t, err)

	vs, err := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(t, types.NewValidatorStatus(true, now), vs)

	vs, err = k.GetValidatorStatus(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(t, types.NewValidatorStatus(false, time.Time{}), vs)
}

func TestActivateValidatorOK(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	ctx = ctx.WithBlockTime(now)
	err := k.Activate(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)

	vs, err := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(t, types.NewValidatorStatus(true, now), vs)

	vs, err = k.GetValidatorStatus(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(t, types.NewValidatorStatus(false, time.Time{}), vs)
}

func TestFailActivateAlreadyActive(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	ctx = ctx.WithBlockTime(now)
	err := k.Activate(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	err = k.Activate(ctx, testapp.Validators[0].ValAddress)
	require.ErrorIs(t, err, types.ErrValidatorAlreadyActive)
}

func TestFailActivateTooSoon(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()

	// Set validator to be inactive just now.
	err := k.SetValidatorStatus(ctx, testapp.Validators[0].ValAddress, types.NewValidatorStatus(false, now))
	require.NoError(t, err)

	// You can't activate until it's been at least InactivePenaltyDuration nanosec.
	params, err := k.GetParams(ctx)
	require.NoError(t, err)
	penaltyDuration := params.InactivePenaltyDuration
	require.ErrorIs(
		t,
		k.Activate(ctx.WithBlockTime(now), testapp.Validators[0].ValAddress),
		types.ErrTooSoonToActivate,
	)
	require.ErrorIs(
		t,
		k.Activate(ctx.WithBlockTime(now.Add(time.Duration(penaltyDuration/2))), testapp.Validators[0].ValAddress),
		types.ErrTooSoonToActivate,
	)

	// So far there must be no changes to the validator's status.
	vs, err := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(t, types.NewValidatorStatus(false, now), vs)

	// Now the time has come.
	require.NoError(
		t,
		k.Activate(ctx.WithBlockTime(now.Add(time.Duration(penaltyDuration))), testapp.Validators[0].ValAddress),
	)
	vs, err = k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(t, types.NewValidatorStatus(true, now.Add(time.Duration(penaltyDuration))), vs)
}

func TestMissReportSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	next := now.Add(time.Duration(10))
	err := k.SetValidatorStatus(ctx, testapp.Validators[0].ValAddress, types.NewValidatorStatus(true, now))
	require.NoError(t, err)

	err = k.MissReport(ctx.WithBlockTime(next), testapp.Validators[0].ValAddress, next)
	require.NoError(t, err)

	vs, err := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(t, types.NewValidatorStatus(false, next), vs)
}

func TestMissReportTooSoonNoop(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	prev := time.Now().UTC()
	now := prev.Add(time.Duration(10))

	err := k.SetValidatorStatus(ctx, testapp.Validators[0].ValAddress, types.NewValidatorStatus(true, now))
	require.NoError(t, err)

	err = k.MissReport(ctx.WithBlockTime(prev), testapp.Validators[0].ValAddress, prev)
	require.NoError(t, err)

	vs, err := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(t, types.NewValidatorStatus(true, now), vs)
}

func TestMissReportAlreadyInactiveNoop(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	next := now.Add(time.Duration(10))

	err := k.SetValidatorStatus(ctx, testapp.Validators[0].ValAddress, types.NewValidatorStatus(false, now))
	require.NoError(t, err)

	err = k.MissReport(ctx.WithBlockTime(next), testapp.Validators[0].ValAddress, next)
	require.NoError(t, err)

	vs, err := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(t, types.NewValidatorStatus(false, now), vs)
}
