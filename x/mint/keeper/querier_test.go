package keeper_test

import (
	"fmt"
	"testing"

	"github.com/GeoDB-Limited/odin-core/x/common/testapp"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/stretchr/testify/require"

	mintkeeper "github.com/GeoDB-Limited/odin-core/x/mint/keeper"
	minttypes "github.com/GeoDB-Limited/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewQuerier(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(true)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := mintkeeper.NewQuerier(app.MintKeeper, legacyQuerierCdc.LegacyAmino)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(ctx, []string{minttypes.QueryParams}, query)
	require.NoError(t, err)

	_, err = querier(ctx, []string{minttypes.QueryInflation}, query)
	require.NoError(t, err)

	_, err = querier(ctx, []string{minttypes.QueryAnnualProvisions}, query)
	require.NoError(t, err)

	_, err = querier(ctx, []string{minttypes.QueryIntegrationAddresses, "bsc"}, query)
	require.Error(t, err, "integration address not supported")

	_, err = querier(ctx, []string{"foo"}, query)
	require.Error(t, err)
}

func TestQueryParams(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(true)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := mintkeeper.NewQuerier(app.MintKeeper, legacyQuerierCdc.LegacyAmino)

	var params minttypes.Params

	res, sdkErr := querier(ctx, []string{minttypes.QueryParams}, abci.RequestQuery{})
	require.NoError(t, sdkErr)

	err := app.LegacyAmino().UnmarshalJSON(res, &params)
	require.NoError(t, err)

	require.Equal(t, app.MintKeeper.GetParams(ctx), params)
}

func TestQueryInflation(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(true)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := mintkeeper.NewQuerier(app.MintKeeper, legacyQuerierCdc.LegacyAmino)

	var inflation sdk.Dec

	res, sdkErr := querier(ctx, []string{minttypes.QueryInflation}, abci.RequestQuery{})
	require.NoError(t, sdkErr)

	err := app.LegacyAmino().UnmarshalJSON(res, &inflation)
	require.NoError(t, err)

	require.Equal(t, app.MintKeeper.GetMinter(ctx).Inflation, inflation)
}

func TestQueryAnnualProvisions(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(true)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := mintkeeper.NewQuerier(app.MintKeeper, legacyQuerierCdc.LegacyAmino)

	var annualProvisions sdk.Dec

	res, sdkErr := querier(ctx, []string{minttypes.QueryAnnualProvisions}, abci.RequestQuery{})
	require.NoError(t, sdkErr)

	err := app.LegacyAmino().UnmarshalJSON(res, &annualProvisions)
	require.NoError(t, err)

	require.Equal(t, app.MintKeeper.GetMinter(ctx).AnnualProvisions, annualProvisions)
}

func TestQueryIntegrationAddresses(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(true)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := mintkeeper.NewQuerier(app.MintKeeper, legacyQuerierCdc.LegacyAmino)

	_, sdkErr := querier(ctx, []string{minttypes.QueryIntegrationAddresses, "eth"}, abci.RequestQuery{})
	require.Error(t, sdkErr, "integration address not supported")
}

func TestQueryTreasuryPool(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(true)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := mintkeeper.NewQuerier(app.MintKeeper, legacyQuerierCdc.LegacyAmino)
	req := require.New(t)

	var treasuryPool sdk.Coins
	res, sdkErr := querier(ctx, []string{minttypes.QueryTreasuryPool}, abci.RequestQuery{})
	req.NoError(sdkErr)
	fmt.Printf("\n!! TreasuryPool: %s!!\n\n", treasuryPool)

	err := app.LegacyAmino().UnmarshalJSON(res, &treasuryPool)
	req.NoError(err)

	req.Equal(app.MintKeeper.GetMintPool(ctx).TreasuryPool, treasuryPool)
}

func TestQueryCommunityPool(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(true)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := mintkeeper.NewQuerier(app.MintKeeper, legacyQuerierCdc.LegacyAmino)
	req := require.New(t)

	var communityPool sdk.DecCoins
	res, sdkErr := querier(ctx, []string{minttypes.QueryCommunityPool}, abci.RequestQuery{})
	req.NoError(sdkErr)
	fmt.Printf("res: %s\n", res)

	err := app.LegacyAmino().UnmarshalJSON(res, &communityPool)
	req.NoError(err)

	req.Equal(app.DistrKeeper.GetFeePool(ctx), communityPool)
}

// func TestQueryTotalSupply(t *testing.T) {
// 	app, ctx, _ := testapp.CreateTestInput(true)
// 	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
// 	querier := mintkeeper.NewQuerier(app.MintKeeper, legacyQuerierCdc.LegacyAmino)
// 	req := require.New(t)

// 	var totalSupply sdk.Dec
// 	res, sdkErr := querier(ctx, []string{minttypes.QueryTotalSupply}, abci.RequestQuery{})
// 	req.NoError(sdkErr)
// 	fmt.Printf("res: %s\n", res)

// 	err := app.LegacyAmino().UnmarshalJSON(res, &totalSupply)
// 	req.NoError(err)

// 	req.Equal(app.DistrKeeper.GetFeePool(ctx), totalSupply)
// }
