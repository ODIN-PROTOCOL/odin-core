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
			res, err := msgServer.RequestData(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgReportData:
			res, err := msgServer.ReportData(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCreateDataSource:
			res, err := msgServer.CreateDataSource(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgEditDataSource:
			res, err := msgServer.EditDataSource(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCreateOracleScript:
			res, err := msgServer.CreateOracleScript(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgEditOracleScript:
			res, err := msgServer.EditOracleScript(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgActivate:
			res, err := msgServer.Activate(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateParams:
			res, err := msgServer.UpdateParams(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, errors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
