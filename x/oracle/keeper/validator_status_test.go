package oraclekeeper_test

import (
	"github.com/ODIN-PROTOCOL/odin-core/x/common/testapp"
	abci "github.com/cometbft/cometbft/abci/types"
)

func defaultVotes() []abci.VoteInfo {
	return []abci.VoteInfo{{
		Validator: abci.Validator{
			Address: testapp.Validators[0].PubKey.Address(),
			Power:   70,
		},
		SignedLastBlock: true,
	}, {
		Validator: abci.Validator{
			Address: testapp.Validators[1].PubKey.Address(),
			Power:   20,
		},
		SignedLastBlock: true,
	}, {
		Validator: abci.Validator{
			Address: testapp.Validators[2].PubKey.Address(),
			Power:   10,
		},
		SignedLastBlock: true,
	}}
}

// func TestAllocateTokenNoActiveValidators(t *testing.T) {
// 	app, ctx, k := testapp.CreateTestInput(false)
// 	// Set collected fee to 1000000odin and 70% oracle reward proportion.
// 	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
// 	feeCollector.SetCoins(Coins1000000loki)
// 	app.AccountKeeper.SetAccount(ctx, feeCollector)
// 	k.SetParam(ctx, types.KeyOracleRewardPercentage, 70)
// 	require.Equal(t, Coins1000000loki, app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
// 	// No active oracle validators so nothing should happen.
// 	k.AllocateTokens(ctx, defaultVotes())
// 	require.Equal(t, Coins1000000loki, app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
// 	require.Equal(t, sdk.Coins(nil), app.SupplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins())
// }

// func TestAllocateTokensOneActive(t *testing.T) {
// 	app, ctx, k := testapp.CreateTestInput(false)
// 	// Set collected fee to 1000000odin + 70% oracle reward proportion.
// 	feeCollector := app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
// 	feeCollector.SetCoins(Coins1000000loki)
// 	app.AccountKeeper.SetAccount(ctx, feeCollector)
// 	k.SetParam(ctx, types.KeyOracleRewardPercentage, 70)
// 	require.Equal(t, Coins1000000loki, app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
// 	// From 70% of fee, 2% should go to community pool, the rest goes to the only active validator.
// 	k.Activate(ctx, testapp.Validators[1].ValAddress)
// 	k.AllocateTokens(ctx, defaultVotes())
// 	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("odin", 300000)), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
// 	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("odin", 700000)), app.SupplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins())
// 	require.Equal(t, sdk.DecCoins{{Denom: "odin", Amount: sdk.NewDec(14000)}}, app.DistrKeeper.GetFeePool(ctx).CommunityPool)
// 	require.Equal(t, sdk.DecCoins(nil), app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress))
// 	require.Equal(t, sdk.DecCoins{{Denom: "odin", Amount: sdk.NewDec(686000)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress))
// 	require.Equal(t, sdk.DecCoins(nil), app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[2].ValAddress))
// }

// func TestAllocateTokensAllActive(t *testing.T) {
// 	app, ctx, k := testapp.CreateTestInput(true)
// 	// Set collected fee to 1000000odin + 70% oracle reward proportion.
// 	feeCollector := app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
// 	feeCollector.SetCoins(Coins1000000loki)
// 	app.AccountKeeper.SetAccount(ctx, feeCollector)
// 	k.SetParam(ctx, types.KeyOracleRewardPercentage, 70)
// 	require.Equal(t, Coins1000000loki, app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
// 	// From 70% of fee, 2% should go to community pool, the rest get split to validators.
// 	k.AllocateTokens(ctx, defaultVotes())
// 	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("odin", 300000)), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
// 	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("odin", 700000)), app.SupplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins())
// 	require.Equal(t, sdk.DecCoins{{Denom: "odin", Amount: sdk.NewDec(14000)}}, app.DistrKeeper.GetFeePool(ctx).CommunityPool)
// 	require.Equal(t, sdk.DecCoins{{Denom: "odin", Amount: sdk.NewDec(480200)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress))
// 	require.Equal(t, sdk.DecCoins{{Denom: "odin", Amount: sdk.NewDec(137200)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress))
// 	require.Equal(t, sdk.DecCoins{{Denom: "odin", Amount: sdk.NewDec(68600)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[2].ValAddress))
// }

// func TestGetDefaultValidatorStatus(t *testing.T) {
// 	_, ctx, k := testapp.CreateTestInput(false)
// 	vs := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
// 	require.Equal(t, types.NewValidatorStatus(false, time.Time{}), vs)
// }

// func TestGetSetValidatorStatus(t *testing.T) {
// 	_, ctx, k := testapp.CreateTestInput(false)
// 	now := time.Now().UTC()
// 	// After setting status of the 1st validator, we should be able to get it back.
// 	k.SetValidatorStatus(ctx, testapp.Validators[0].ValAddress, types.NewValidatorStatus(true, now))
// 	vs := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
// 	require.Equal(t, types.NewValidatorStatus(true, now), vs)
// 	vs = k.GetValidatorStatus(ctx, testapp.Validators[1].ValAddress)
// 	require.Equal(t, types.NewValidatorStatus(false, time.Time{}), vs)
// }

// func TestActivateValidatorOK(t *testing.T) {
// 	_, ctx, k := testapp.CreateTestInput(false)
// 	now := time.Now().UTC()
// 	ctx = ctx.WithBlockTime(now)
// 	err := k.Activate(ctx, testapp.Validators[0].ValAddress)
// 	require.NoError(t, err)
// 	vs := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
// 	require.Equal(t, types.NewValidatorStatus(true, now), vs)
// 	vs = k.GetValidatorStatus(ctx, testapp.Validators[1].ValAddress)
// 	require.Equal(t, types.NewValidatorStatus(false, time.Time{}), vs)
// }

// func TestFailActivateAlreadyActive(t *testing.T) {
// 	_, ctx, k := testapp.CreateTestInput(false)
// 	now := time.Now().UTC()
// 	ctx = ctx.WithBlockTime(now)
// 	err := k.Activate(ctx, testapp.Validators[0].ValAddress)
// 	require.NoError(t, err)
// 	err = k.Activate(ctx, testapp.Validators[0].ValAddress)
// 	require.Error(t, err)
// }

// func TestFailActivateTooSoon(t *testing.T) {
// 	_, ctx, k := testapp.CreateTestInput(false)
// 	now := time.Now().UTC()
// 	// Set validator to be inactive just now.
// 	k.SetValidatorStatus(ctx, testapp.Validators[0].ValAddress, types.NewValidatorStatus(false, now))
// 	// You can't activate until it's been at least InactivePenaltyDuration nanosec.
// 	penaltyDuration := k.GetParam(ctx, types.KeyInactivePenaltyDuration)
// 	require.Error(t, k.Activate(ctx.WithBlockTime(now), testapp.Validators[0].ValAddress))
// 	require.Error(t, k.Activate(ctx.WithBlockTime(now.Add(time.Duration(penaltyDuration/2))), testapp.Validators[0].ValAddress))
// 	// So far there must be no changes to the validator's status.
// 	vs := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
// 	require.Equal(t, types.NewValidatorStatus(false, now), vs)
// 	// Now the time has come.
// 	require.NoError(t, k.Activate(ctx.WithBlockTime(now.Add(time.Duration(penaltyDuration))), testapp.Validators[0].ValAddress))
// 	vs = k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
// 	require.Equal(t, types.NewValidatorStatus(true, now.Add(time.Duration(penaltyDuration))), vs)
// }

// func TestMissReportSuccess(t *testing.T) {
// 	_, ctx, k := testapp.CreateTestInput(false)
// 	now := time.Now().UTC()
// 	next := now.Add(time.Duration(10))
// 	k.SetValidatorStatus(ctx, testapp.Validators[0].ValAddress, types.NewValidatorStatus(true, now))
// 	k.MissReport(ctx.WithBlockTime(next), testapp.Validators[0].ValAddress, next)
// 	vs := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
// 	require.Equal(t, types.NewValidatorStatus(false, next), vs)
// }

// func TestMissReportTooSoonNoop(t *testing.T) {
// 	_, ctx, k := testapp.CreateTestInput(false)
// 	prev := time.Now().UTC()
// 	now := prev.Add(time.Duration(10))
// 	k.SetValidatorStatus(ctx, testapp.Validators[0].ValAddress, types.NewValidatorStatus(true, now))
// 	k.MissReport(ctx.WithBlockTime(prev), testapp.Validators[0].ValAddress, prev)
// 	vs := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
// 	require.Equal(t, types.NewValidatorStatus(true, now), vs)
// }

// func TestMissReportAlreadyInactiveNoop(t *testing.T) {
// 	_, ctx, k := testapp.CreateTestInput(false)
// 	now := time.Now().UTC()
// 	next := now.Add(time.Duration(10))
// 	k.SetValidatorStatus(ctx, testapp.Validators[0].ValAddress, types.NewValidatorStatus(false, now))
// 	k.MissReport(ctx.WithBlockTime(next), testapp.Validators[0].ValAddress, next)
// 	vs := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
// 	require.Equal(t, types.NewValidatorStatus(false, now), vs)
// }
