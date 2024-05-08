package oracle_test

import (
	"encoding/hex"
	"testing"

	"cosmossdk.io/math"
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"

	"github.com/ODIN-PROTOCOL/odin-core/testing/testapp"
)

func fromHex(hexStr string) []byte {
	res, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return res
}

func TestRollingSeedCorrect(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false)
	// Initially rolling seed should be all zeros.
	rollingSeed, err := k.GetRollingSeed(ctx)
	require.NoError(t, err)
	require.Equal(
		t,
		fromHex("0000000000000000000000000000000000000000000000000000000000000000"),
		rollingSeed,
	)

	// Every begin block, the rolling seed should get updated.
	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: app.LastBlockHeight() + 1,
		Hash:   fromHex("0100000000000000000000000000000000000000000000000000000000000000"),
	})
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(2)
	rollingSeed, err = app.OracleKeeper.GetRollingSeed(ctx)
	require.NoError(t, err)
	require.Equal(
		t,
		fromHex("0000000000000000000000000000000000000000000000000000000000000001"),
		rollingSeed,
	)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: app.LastBlockHeight() + 1,
		Hash:   fromHex("0200000000000000000000000000000000000000000000000000000000000000"),
	})
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(2)
	rollingSeed, err = app.OracleKeeper.GetRollingSeed(ctx)
	require.NoError(t, err)
	require.Equal(
		t,
		fromHex("0000000000000000000000000000000000000000000000000000000000000102"),
		rollingSeed,
	)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: app.LastBlockHeight() + 1,
		Hash:   fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	})
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(3)
	rollingSeed, err = k.GetRollingSeed(ctx)
	require.NoError(t, err)
	require.Equal(
		t,
		fromHex("00000000000000000000000000000000000000000000000000000000000102ff"),
		rollingSeed,
	)
}

func TestAllocateTokensCalledOnBeginBlock(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false)
	votes := []abci.VoteInfo{{
		Validator: abci.Validator{Address: testapp.Validators[0].PubKey.Address(), Power: 70},
	}, {
		Validator: abci.Validator{Address: testapp.Validators[1].PubKey.Address(), Power: 30},
	}}

	mintParams, err := app.MintKeeper.GetParams(ctx)
	require.NoError(t, err)
	mintParams.MintAir = true
	err = app.MintKeeper.SetParams(ctx, mintParams)
	require.NoError(t, err)

	// Set collected fee to 100loki + 70% oracle reward proportion + disable minting inflation.
	// NOTE: we intentionally keep ctx.BlockHeight = 0, so distr's AllocateTokens doesn't get called.
	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	feeCollectorStartBalance := app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	amt := sdk.NewCoins(sdk.NewInt64Coin("loki", 10000)).Sub(feeCollectorStartBalance...)
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, amt)
	require.NoError(t, err)

	mintParams.MintAir = false
	err = app.MintKeeper.SetParams(ctx, mintParams)
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToModule(
		ctx,
		minttypes.ModuleName,
		authtypes.FeeCollectorName,
		amt,
	)
	require.NoError(t, err)
	distModule := app.AccountKeeper.GetModuleAccount(ctx, distrtypes.ModuleName)

	app.AccountKeeper.SetAccount(ctx, feeCollector)
	mintParams, err = app.MintKeeper.GetParams(ctx)
	require.NoError(t, err)
	mintParams.InflationMin = math.LegacyZeroDec()
	mintParams.InflationMax = math.LegacyZeroDec()
	err = app.MintKeeper.SetParams(ctx, mintParams)
	require.NoError(t, err)
	params, err := k.GetParams(ctx)
	require.NoError(t, err)
	params.OracleRewardPercentage = 70
	err = k.SetParams(ctx, params)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 10000)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	// If there are no validators active, Calling begin block should be no-op.
	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:            app.LastBlockHeight() + 1,
		Hash:              fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		DecidedLastCommit: abci.CommitInfo{Votes: votes},
	})
	require.NoError(t, err)

	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 10000)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	// 1 validator active, begin block should take 70% of the fee. 2% of that goes to comm pool.
	err = k.Activate(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:            app.LastBlockHeight() + 1,
		Hash:              fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		DecidedLastCommit: abci.CommitInfo{Votes: votes},
	})
	require.NoError(t, err)

	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 3000)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 7000)),
		app.BankKeeper.GetAllBalances(ctx, distModule.GetAddress()),
	)

	// 10000*70%*2% = 140loki
	communityPool, err := distrkeeper.NewQuerier(app.DistrKeeper).CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(140)}},
		communityPool.Pool,
	)

	// 0loki
	validatorOutstandingRewards, err := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Empty(t, validatorOutstandingRewards)

	// 10000*70%*98% = 6860loki
	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(6860)}},
		validatorOutstandingRewards.Rewards,
	)

	// 2 validators active now. 70% of the remaining fee pool will be split 3 ways (comm pool + val1 + val2).
	err = k.Activate(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:            app.LastBlockHeight() + 1,
		Hash:              fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		DecidedLastCommit: abci.CommitInfo{Votes: votes},
	})
	require.NoError(t, err)

	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 900)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 9100)),
		app.BankKeeper.GetAllBalances(ctx, distModule.GetAddress()),
	)

	// 140loki + 3000*70%*2% = 182loki
	communityPool, err = distrkeeper.NewQuerier(app.DistrKeeper).CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(182)}},
		communityPool.Pool,
	)

	// 3000*70%*98%*70% = 1440.6loki
	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(14406, 1)}},
		validatorOutstandingRewards.Rewards,
	)

	// 6860loki + 3000*70%*98%*30% = 7477.4loki
	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(74774, 1)}},
		validatorOutstandingRewards.Rewards,
	)

	// 1 validator becomes in active, and will not get reward this time.
	err = k.MissReport(ctx, testapp.Validators[1].ValAddress, testapp.ParseTime(100))
	require.NoError(t, err)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:            app.LastBlockHeight() + 1,
		Hash:              fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		DecidedLastCommit: abci.CommitInfo{Votes: votes},
	})
	require.NoError(t, err)

	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 270)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 9730)),
		app.BankKeeper.GetAllBalances(ctx, distModule.GetAddress()),
	)

	// 182loki + 900*70%*2% = 194,6loki
	communityPool, err = distrkeeper.NewQuerier(app.DistrKeeper).CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(194)}}, // due to FundCommunityPool stupid logic
		communityPool.Pool,
	)

	// 1440.6loki + 900*70%*98% = 2058loki
	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(2058)}},
		validatorOutstandingRewards.Rewards,
	)

	// 7477.4loki
	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(74774, 1)}},
		validatorOutstandingRewards.Rewards,
	)
}

func TestAllocateTokensWithDistrAllocateTokens(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false)
	ctx = ctx.WithBlockHeight(10) // Set block height to ensure distr's AllocateTokens gets called.
	votes := []abci.VoteInfo{{
		Validator: abci.Validator{Address: testapp.Validators[0].PubKey.Address(), Power: 70},
	}, {
		Validator: abci.Validator{Address: testapp.Validators[1].PubKey.Address(), Power: 30},
	}}

	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	distModule := app.AccountKeeper.GetModuleAccount(ctx, distrtypes.ModuleName)

	feeCollectorStartBalance := app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	amt := sdk.NewCoins(sdk.NewInt64Coin("loki", 1000)).Sub(feeCollectorStartBalance...)

	mintParams, err := app.MintKeeper.GetParams(ctx)
	require.NoError(t, err)
	mintParams.MintAir = true
	err = app.MintKeeper.SetParams(ctx, mintParams)
	require.NoError(t, err)

	// Set collected fee to 1000loki + 70% oracle reward proportion + disable minting inflation.
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, amt)
	require.NoError(t, err)

	mintParams.MintAir = false
	err = app.MintKeeper.SetParams(ctx, mintParams)
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToModule(
		ctx,
		minttypes.ModuleName,
		authtypes.FeeCollectorName,
		amt,
	)
	require.NoError(t, err)

	app.AccountKeeper.SetAccount(ctx, feeCollector)

	mintParams, err = app.MintKeeper.GetParams(ctx)
	require.NoError(t, err)
	mintParams.InflationMin = math.LegacyZeroDec()
	mintParams.InflationMax = math.LegacyZeroDec()
	err = app.MintKeeper.SetParams(ctx, mintParams)
	require.NoError(t, err)

	params, err := k.GetParams(ctx)
	require.NoError(t, err)
	params.OracleRewardPercentage = 70
	err = k.SetParams(ctx, params)
	require.NoError(t, err)

	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 1000)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	// Only Validators[0] active. After we call begin block:
	//   700loki = 70% go to oracle pool
	//     14loki (2%) go to community pool
	//     686loki go to Validators[0] (active)
	//   300loki = 30% go to distr pool
	//     6loki (2%) go to community pool
	//     294loki split among voters
	//        205.8loki (70%) go to Validators[0]
	//        88.2loki (30%) go to Validators[1]
	// In summary
	//   Community pool: 20 + 6 = 26
	//   Validators[0]: 686 + 205.8 = 891.8
	//   Validators[1]: 88.2
	err = k.Activate(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(2).
		WithVoteInfos(votes).
		WithHeaderHash(fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))
	_, err = app.BeginBlocker(ctx)
	require.NoError(t, err)

	//_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
	//	Height:            app.LastBlockHeight() + 1,
	//	Hash:              fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	//	DecidedLastCommit: abci.CommitInfo{Votes: votes},
	//})
	require.NoError(t, err)

	require.Equal(t, sdk.Coins{}, app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()))
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 1000)),
		app.BankKeeper.GetAllBalances(ctx, distModule.GetAddress()),
	)

	communityPool, err := distrkeeper.NewQuerier(app.DistrKeeper).CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(20)}},
		communityPool.Pool,
	)

	validatorOutstandingRewards, err := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(8918, 1)}},
		validatorOutstandingRewards.Rewards,
	)

	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(882, 1)}},
		validatorOutstandingRewards.Rewards,
	)
}
