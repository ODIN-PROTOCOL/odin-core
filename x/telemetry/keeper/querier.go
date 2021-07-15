package keeper

import (
	commontypes "github.com/GeoDB-Limited/odin-core/x/common/types"
	telemetrytypes "github.com/GeoDB-Limited/odin-core/x/telemetry/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper, cdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case telemetrytypes.QueryTopBalances:
			return queryTopBalances(ctx, path[1:], keeper, cdc, req)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown telemetry query endpoint")
		}
	}
}

func queryTopBalances(ctx sdk.Context, path []string, k Keeper, cdc *codec.LegacyAmino, req abci.RequestQuery) ([]byte, error) {
	if len(path) > 1 {
		return nil, sdkerrors.ErrInvalidRequest
	}
	var params commontypes.QueryPaginationParams
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return commontypes.QueryOK(cdc, k.GetPaginatedBalances(ctx, path[1], params.Limit, params.Offset))
}
