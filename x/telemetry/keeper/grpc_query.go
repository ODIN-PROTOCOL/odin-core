package keeper

import (
	"context"
	telemetrytypes "github.com/GeoDB-Limited/odin-core/x/telemetry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ telemetrytypes.QueryServer = Keeper{}

func (k Keeper) TopBalances(c context.Context, request *telemetrytypes.QueryTopBalancesRequest) (*telemetrytypes.QueryTopBalancesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &telemetrytypes.QueryTopBalancesResponse{
		Balances: k.GetPaginatedBalances(ctx, request.GetDenom(), request.Pagination.GetLimit(), request.Pagination.GetOffset()),
	}, nil
}
