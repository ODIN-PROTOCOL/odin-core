package testapp

import (
	"math/rand"

	odinapp "github.com/ODIN-PROTOCOL/odin-core/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
