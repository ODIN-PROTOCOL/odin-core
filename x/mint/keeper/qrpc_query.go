package keeper

import (
	"context"

	minttypes "github.com/GeoDB-Limited/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ minttypes.QueryServer = Keeper{}

// Params returns params of the mint module.
func (k Keeper) Params(c context.Context, _ *minttypes.QueryParamsRequest) (*minttypes.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &minttypes.QueryParamsResponse{Params: params}, nil
}

// Inflation returns minter.Inflation of the mint module.
func (k Keeper) Inflation(
	c context.Context,
	_ *minttypes.QueryInflationRequest,
) (*minttypes.QueryInflationResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minter := k.GetMinter(ctx)

	return &minttypes.QueryInflationResponse{Inflation: minter.Inflation}, nil
}

// AnnualProvisions returns minter.AnnualProvisions of the mint module.
func (k Keeper) AnnualProvisions(
	c context.Context,
	_ *minttypes.QueryAnnualProvisionsRequest,
) (*minttypes.QueryAnnualProvisionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	minter := k.GetMinter(ctx)

	return &minttypes.QueryAnnualProvisionsResponse{AnnualProvisions: minter.AnnualProvisions}, nil
}

// IntegrationAddress returns ethereum integration address
func (k Keeper) IntegrationAddress(
	c context.Context,
	req *minttypes.QueryIntegrationAddressRequest,
) (*minttypes.QueryIntegrationAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	addresses := k.GetParams(ctx).IntegrationAddresses

	return &minttypes.QueryIntegrationAddressResponse{IntegrationAddress: addresses[req.NetworkName]}, nil
}

// TreasuryPool returns current treasury pool
func (k Keeper) TreasuryPool(
	c context.Context,
	_ *minttypes.QueryTreasuryPoolRequest,
) (*minttypes.QueryTreasuryPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	mintPool := k.GetMintPool(ctx)

	return &minttypes.QueryTreasuryPoolResponse{TreasuryPool: mintPool.TreasuryPool}, nil
}

func (k Keeper) OdinInfo(c context.Context, request *minttypes.QueryOdinInfoRequest) (*minttypes.QueryOdinInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	mintPool := k.GetMintPool(ctx)

	totalActiveSupply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(bondDenom).ToDec().Sub(k.mintKeeper.GetMintPool(ctx).TreasuryPool.AmountOf(bondDenom).ToDec())

	// communityPool := app.DistrKeeper.GetFeePool(ctx).CommunityPool
	feePool := k.distrKeeper.GetFeePool(ctx)

	// balances
	validatorsResp, err := k.stakingQuerier.Validators(c, OdinInfoRequestToValidatorsRequest(request))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get validators")
	}
	accounts, err := ValidatorsToAccounts(validatorsResp.GetValidators())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get validators accounts addresses")
	}
	// extendedValidatorsResp := ValidatorsResponseToExtendedValidatorsResponse(validatorsResp)
	// extendedValidatorsResp.Balances = k.GetBalances(ctx, accounts...)
	balances := k.GetBalances(ctx, accounts...)

	return &minttypes.QueryOdinInfoResponse{
		// Total:
		Balances:          balances,
		CommunityPool:     feePool.CommunityPool,
		TreasuryPool:      mintPool.TreasuryPool,
		DataProvidersPool: mintPool.DataProvidersPool,
	}, nil
}
