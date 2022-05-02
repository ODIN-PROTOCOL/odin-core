package keeper

import (
	commontypes "github.com/ODIN-PROTOCOL/odin-core/x/common/types"
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, _ abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case minttypes.QueryParams:
			return queryParams(ctx, k, legacyQuerierCdc)

		case minttypes.QueryInflation:
			return queryInflation(ctx, k, legacyQuerierCdc)

		case minttypes.QueryAnnualProvisions:
			return queryAnnualProvisions(ctx, k, legacyQuerierCdc)

		case minttypes.QueryIntegrationAddresses:
			return queryIntegrationAddresses(ctx, path[1:], k, legacyQuerierCdc)

		case minttypes.QueryTreasuryPool:
			return queryTreasuryPool(ctx, k, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryInflation(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	minter := k.GetMinter(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, minter.Inflation)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryAnnualProvisions(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	minter := k.GetMinter(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, minter.AnnualProvisions)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryIntegrationAddresses(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "integration addresses not specified")
	}

	integrationAddress, ok := k.GetParams(ctx).IntegrationAddresses[path[0]]
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotSupported, "integration address not supported")
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, integrationAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryTreasuryPool(ctx sdk.Context, k Keeper, cdc *codec.LegacyAmino) ([]byte, error) {
	return commontypes.QueryOK(cdc, k.GetMintPool(ctx).TreasuryPool)
}
