package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "odin/testutil/keeper"
	"odin/x/odin/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.OdinKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
