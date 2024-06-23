package keeper

import (
	"context"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
)

var _ minttypes.QueryServer = queryServer{}

func NewQueryServerImpl(k Keeper) minttypes.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// Params returns params of the mint module.
func (q queryServer) Params(ctx context.Context, _ *minttypes.QueryParamsRequest) (*minttypes.QueryParamsResponse, error) {
	params, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &minttypes.QueryParamsResponse{Params: params}, nil
}

// Inflation returns minter.Inflation of the mint module.
func (q queryServer) Inflation(
	ctx context.Context,
	_ *minttypes.QueryInflationRequest,
) (*minttypes.QueryInflationResponse, error) {
	minter, err := q.k.Minter.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &minttypes.QueryInflationResponse{Inflation: minter.Inflation}, nil
}

// AnnualProvisions returns minter.AnnualProvisions of the mint module.
func (q queryServer) AnnualProvisions(
	ctx context.Context,
	_ *minttypes.QueryAnnualProvisionsRequest,
) (*minttypes.QueryAnnualProvisionsResponse, error) {
	minter, err := q.k.Minter.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &minttypes.QueryAnnualProvisionsResponse{AnnualProvisions: minter.AnnualProvisions}, nil
}

// TreasuryPool returns current treasury pool
func (q queryServer) TreasuryPool(
	ctx context.Context,
	_ *minttypes.QueryTreasuryPoolRequest,
) (*minttypes.QueryTreasuryPoolResponse, error) {
	mintPool, err := q.k.MintPool.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &minttypes.QueryTreasuryPoolResponse{TreasuryPool: mintPool.TreasuryPool}, nil
}

// CurrentMintVolume returns current mint volume
func (q queryServer) CurrentMintVolume(
	ctx context.Context,
	_ *minttypes.QueryCurrentMintVolumeRequest,
) (*minttypes.QueryCurrentMintVolumeResponse, error) {
	minter, err := q.k.Minter.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &minttypes.QueryCurrentMintVolumeResponse{CurrentMintVolume: minter.CurrentMintVolume}, nil
}
