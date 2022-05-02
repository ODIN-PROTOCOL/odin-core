package auction

import (
	auctionkeeper "github.com/ODIN-PROTOCOL/odin-core/x/auction/keeper"
	auctiontypes "github.com/ODIN-PROTOCOL/odin-core/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k auctionkeeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		msgServer := auctionkeeper.NewMsgServerImpl(k)
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *auctiontypes.MsgBuyCoins:
			res, err := msgServer.BuyCoins(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", auctiontypes.ModuleName, msg)
		}
	}
}
