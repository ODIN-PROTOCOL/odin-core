package types

import (
	distrtypes "x/distribution/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DistrKeeper defines the expected distribution keeper.
type DistrKeeper interface {
	GetValidatorHistoricalRewards(ctx sdk.Context, val sdk.ValAddress, period uint64) (rewards distrtypes.ValidatorHistoricalRewards)
}
