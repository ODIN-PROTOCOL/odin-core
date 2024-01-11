package odin_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "odin/testutil/keeper"
	"odin/testutil/nullify"
	"odin/x/odin/module"
	"odin/x/odin/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.OdinKeeper(t)
	odin.InitGenesis(ctx, k, genesisState)
	got := odin.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
