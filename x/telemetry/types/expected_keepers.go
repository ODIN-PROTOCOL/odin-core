package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// DistrKeeper defines the expected distribution keeper.
type DistrKeeper interface {
	GetValidatorHistoricalRewards(ctx sdk.Context, val sdk.ValAddress, period uint64) (rewards distrtypes.ValidatorHistoricalRewards)
}
