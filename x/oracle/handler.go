package oracle

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

type Handler = func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error)

// NewHandler creates the msg handler of this module, as required by Cosmos-SDK standard.
func NewHandler(k keeper.Keeper) Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		msgServer := keeper.NewMsgServerImpl(k)
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgRequestData:
			res, err := msgServer.RequestData(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgReportData:
			res, err := msgServer.ReportData(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCreateDataSource:
			res, err := msgServer.CreateDataSource(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgEditDataSource:
			res, err := msgServer.EditDataSource(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCreateOracleScript:
			res, err := msgServer.CreateOracleScript(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgEditOracleScript:
			res, err := msgServer.EditOracleScript(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgActivate:
			res, err := msgServer.Activate(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateParams:
			res, err := msgServer.UpdateParams(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, errors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
