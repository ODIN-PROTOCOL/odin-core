package mint_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/runtime"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := runtime.AppI.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.InitChain(
		abcitypes.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(minttypes.ModuleName))
	require.NotNil(t, acc)
}
