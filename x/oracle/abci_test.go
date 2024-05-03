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
		Hash: fromHex("0100000000000000000000000000000000000000000000000000000000000000"),
	})
	require.NoError(t, err)
	_, err = app.Commit()
	require.NoError(t, err)

	rollingSeed, err = k.GetRollingSeed(ctx)
	require.NoError(t, err)
	require.Equal(
		t,
		fromHex("0000000000000000000000000000000000000000000000000000000000000001"),
		rollingSeed,
	)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Hash: fromHex("0200000000000000000000000000000000000000000000000000000000000000"),
	})
	require.NoError(t, err)
	_, err = app.Commit()
	require.NoError(t, err)

	rollingSeed, err = k.GetRollingSeed(ctx)
	require.NoError(t, err)
	require.Equal(
		t,
		fromHex("0000000000000000000000000000000000000000000000000000000000000102"),
		rollingSeed,
	)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Hash: fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	})
	require.NoError(t, err)
	_, err = app.Commit()
	require.NoError(t, err)

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
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("loki", 100)))
	require.NoError(t, err)

	mintParams.MintAir = false
	err = app.MintKeeper.SetParams(ctx, mintParams)
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToModule(
		ctx,
		minttypes.ModuleName,
		authtypes.FeeCollectorName,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 100)),
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
		sdk.NewCoins(sdk.NewInt64Coin("loki", 100)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	// If there are no validators active, Calling begin block should be no-op.
	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Hash:              fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		DecidedLastCommit: abci.CommitInfo{Votes: votes},
	})
	require.NoError(t, err)
	_, err = app.Commit()
	require.NoError(t, err)

	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 100)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	// 1 validator active, begin block should take 70% of the fee. 2% of that goes to comm pool.
	err = k.Activate(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Hash:              fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		DecidedLastCommit: abci.CommitInfo{Votes: votes},
	})
	require.NoError(t, err)
	_, err = app.Commit()
	require.NoError(t, err)

	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 30)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 70)),
		app.BankKeeper.GetAllBalances(ctx, distModule.GetAddress()),
	)

	// 100*70%*2% = 1.4loki
	communityPool, err := distrkeeper.NewQuerier(app.DistrKeeper).CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(14, 1)}},
		communityPool.Pool,
	)

	// 0loki
	validatorOutstandingRewards, err := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Empty(t, validatorOutstandingRewards)

	// 100*70%*98% = 68.6loki
	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(686, 1)}},
		validatorOutstandingRewards.Rewards,
	)

	// 2 validators active now. 70% of the remaining fee pool will be split 3 ways (comm pool + val1 + val2).
	err = k.Activate(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Hash:              fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		DecidedLastCommit: abci.CommitInfo{Votes: votes},
	})
	require.NoError(t, err)
	_, err = app.Commit()
	require.NoError(t, err)

	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 9)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 91)),
		app.BankKeeper.GetAllBalances(ctx, distModule.GetAddress()),
	)

	// 1.4loki + 30*70%*2% = 1.82loki
	communityPool, err = distrkeeper.NewQuerier(app.DistrKeeper).CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(182, 2)}},
		communityPool.Pool,
	)

	// 30*70%*98%*70% = 14.406loki
	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(14406, 3)}},
		validatorOutstandingRewards.Rewards,
	)

	// 68.6loki + 30*70%*98%*30% = 74.774loki
	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(74774, 3)}},
		validatorOutstandingRewards.Rewards,
	)

	// 1 validator becomes in active, and will not get reward this time.
	err = k.MissReport(ctx, testapp.Validators[1].ValAddress, testapp.ParseTime(100))
	require.NoError(t, err)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Hash:              fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		DecidedLastCommit: abci.CommitInfo{Votes: votes},
	})
	require.NoError(t, err)
	_, err = app.Commit()
	require.NoError(t, err)

	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 3)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 97)),
		app.BankKeeper.GetAllBalances(ctx, distModule.GetAddress()),
	)

	// 1.82loki + 6*2% = 1.82loki
	communityPool, err = distrkeeper.NewQuerier(app.DistrKeeper).CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(194, 2)}},
		communityPool.Pool,
	)

	// 14.406loki + 6*98% = 20.286loki
	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(20286, 3)}},
		validatorOutstandingRewards.Rewards,
	)

	// 74.774loki
	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(74774, 3)}},
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

	mintParams, err := app.MintKeeper.GetParams(ctx)
	require.NoError(t, err)
	mintParams.MintAir = true
	err = app.MintKeeper.SetParams(ctx, mintParams)
	require.NoError(t, err)

	// Set collected fee to 100loki + 70% oracle reward proportion + disable minting inflation.
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("loki", 50)))
	require.NoError(t, err)

	mintParams.MintAir = false
	err = app.MintKeeper.SetParams(ctx, mintParams)
	require.NoError(t, err)

	err = app.BankKeeper.SendCoinsFromModuleToModule(
		ctx,
		minttypes.ModuleName,
		authtypes.FeeCollectorName,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 50)),
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

	// Set block proposer to Validators[1], who will receive 5% bonus.
	err = app.DistrKeeper.SetPreviousProposerConsAddr(ctx, testapp.Validators[1].Address.Bytes())
	require.NoError(t, err)

	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 50)),
		app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()),
	)
	// Only Validators[0] active. After we call begin block:
	//   35loki = 70% go to oracle pool
	//     0.7loki (2%) go to community pool
	//     34.3loki go to Validators[0] (active)
	//   15loki = 30% go to distr pool
	//     0.3loki (2%) go to community pool
	//     2.25loki (15%) go to Validators[1] (proposer)
	//     12.45loki split among voters
	//        8.715loki (70%) go to Validators[0]
	//        3.735loki (30%) go to Validators[1]
	// In summary
	//   Community pool: 0.7 + 0.3 = 1
	//   Validators[0]: 34.3 + 8.715 = 43.015
	//   Validators[1]: 2.25 + 3.735 = 5.985
	err = k.Activate(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Hash:              fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		DecidedLastCommit: abci.CommitInfo{Votes: votes},
	})
	require.NoError(t, err)
	_, err = app.Commit()
	require.NoError(t, err)

	require.Equal(t, sdk.Coins{}, app.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress()))
	require.Equal(
		t,
		sdk.NewCoins(sdk.NewInt64Coin("loki", 50)),
		app.BankKeeper.GetAllBalances(ctx, distModule.GetAddress()),
	)

	communityPool, err := distrkeeper.NewQuerier(app.DistrKeeper).CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDec(1)}},
		communityPool.Pool,
	)

	validatorOutstandingRewards, err := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(44590, 3)}},
		validatorOutstandingRewards.Rewards,
	)

	validatorOutstandingRewards, err = app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.Equal(
		t,
		sdk.DecCoins{{Denom: "loki", Amount: math.LegacyNewDecWithPrec(4410, 3)}},
		validatorOutstandingRewards.Rewards,
	)
}
