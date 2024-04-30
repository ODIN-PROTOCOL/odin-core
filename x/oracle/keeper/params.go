package keeper

import (
	"context"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// SetParams sets the x/oracle module parameters.
func (k Keeper) SetParams(ctx context.Context, p types.Params) error {
	if err := p.Validate(); err != nil {
		return err
	}

	return k.Params.Set(ctx, p)
}

// GetParams returns the current x/oracle module parameters.
func (k Keeper) GetParams(ctx context.Context) (p types.Params, err error) {
	return k.Params.Get(ctx)
}
