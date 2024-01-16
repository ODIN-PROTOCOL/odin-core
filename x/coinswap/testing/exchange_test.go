package testing

import (
	"testing"

	"github.com/ODIN-PROTOCOL/odin-core/testing/testapp"
	swaptypes "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

const (
	minigeo     = "minigeo"
	loki        = "loki"
	initialRate = 10
)

func TestKeeper_ExchangeDenom(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(false, true, true)

	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin("loki", sdk.NewInt(1))))
	assert.NoError(t, err)

	// add tokens for exchange
	err = app.DistrKeeper.FundCommunityPool(ctx, sdk.NewCoins(sdk.NewCoin("loki", sdk.NewInt(1))), app.AccountKeeper.GetModuleAddress(minttypes.ModuleName))
	assert.NoError(t, err)

	app.CoinswapKeeper.SetInitialRate(ctx, sdk.NewDec(initialRate))
	app.CoinswapKeeper.SetParams(ctx, swaptypes.DefaultParams())

	err = app.CoinswapKeeper.ExchangeDenom(ctx, minigeo, loki, sdk.NewInt64Coin(minigeo, 10), testapp.Alice.Address)

	assert.NoError(t, err, "exchange denom failed")
}
