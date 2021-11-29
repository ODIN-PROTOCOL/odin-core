package testapp

import (
	odinapp "github.com/GeoDB-Limited/odin-core/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/rand"
)

const (
	Seed = 42
)

var (
	RAND *rand.Rand
)

func init() {
	RAND = rand.New(rand.NewSource(Seed))
	odinapp.SetBech32AddressPrefixesAndBip44CoinType(sdk.GetConfig())
}
