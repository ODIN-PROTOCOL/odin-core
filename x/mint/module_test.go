package mint_test

import (
	"testing"
	// "cosmossdk.io/log"
	// minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	// abcitypes "github.com/cometbft/cometbft/abci/types"
	// tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	// authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// "github.com/stretchr/testify/require"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {

	// app := NewSimApp("ODINCHAIN", log.NewNopLogger())
	// ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// app.InitChain(
	// 	abcitypes.RequestInitChain{
	// 		AppStateBytes: []byte("{}"),
	// 		ChainId:       "test-chain-id",
	// 	},
	// )

	// acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(minttypes.ModuleName))
	// require.NotNil(t, acc)
}
