package coinswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	coinswapkeeper "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/keeper"
	coinswaptypes "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
)

func NewHandler(k coinswapkeeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		msgServer := coinswapkeeper.NewMsgServerImpl(k)
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *coinswaptypes.MsgExchange:
			res, err := msgServer.Exchange(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", coinswaptypes.ModuleName, msg)
		}
	}
}
