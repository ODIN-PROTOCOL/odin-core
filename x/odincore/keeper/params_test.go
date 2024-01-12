package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/ODIN-PROTOCOL/odin-core/testutil/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/odincore/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.OdincoreKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
