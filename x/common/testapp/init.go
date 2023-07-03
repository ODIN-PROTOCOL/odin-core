package testapp

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"

	odinapp "github.com/ODIN-PROTOCOL/odin-core/app"
)

const (
	Seed = 42
)

var RAND *rand.Rand

func init() {
	RAND = rand.New(rand.NewSource(Seed))
	odinapp.SetBech32AddressPrefixesAndBip44CoinType(sdk.GetConfig())
}
