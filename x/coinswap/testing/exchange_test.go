package testing

import (
	"testing"

	swaptypes "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
	"github.com/ODIN-PROTOCOL/odin-core/x/common/testapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

const (
	minigeo     = "minigeo"
	loki        = "loki"
	initialRate = 10
)

func TestKeeper_ExchangeDenom(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(false, true)

	app.CoinswapKeeper.SetInitialRate(ctx, sdk.NewDec(initialRate))
	app.CoinswapKeeper.SetParams(ctx, swaptypes.DefaultParams())

	err := app.CoinswapKeeper.ExchangeDenom(ctx, minigeo, loki, sdk.NewInt64Coin(minigeo, 10), testapp.Alice.Address)

	assert.NoError(t, err, "exchange denom failed")
}
