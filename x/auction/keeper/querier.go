package keeper

import (
	abci "github.com/github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	auctiontypes "github.com/ODIN-PROTOCOL/odin-core/x/auction/types"
	commontypes "github.com/ODIN-PROTOCOL/odin-core/x/common/types"
)

func NewQuerier(keeper Keeper, cdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case auctiontypes.QueryParams:
			return queryParameters(ctx, keeper, cdc)
		case auctiontypes.QueryAuctionStatus:
			return queryAuctionStatus(ctx, keeper, cdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown auction query endpoint")
		}
	}
}

func queryParameters(ctx sdk.Context, k Keeper, cdc *codec.LegacyAmino) ([]byte, error) {
	return commontypes.QueryOK(cdc, k.GetParams(ctx))
}

func queryAuctionStatus(ctx sdk.Context, k Keeper, cdc *codec.LegacyAmino) ([]byte, error) {
	return commontypes.QueryOK(cdc, k.GetAuctionStatus(ctx))
}
