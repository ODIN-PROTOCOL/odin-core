package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/ODIN-PROTOCOL/odin-core/pkg/gzip"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) RequestData(
	goCtx context.Context,
	msg *types.MsgRequestData,
) (*types.MsgRequestDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	payer, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	_, err = k.PrepareRequest(ctx, msg, payer, nil)
	if err != nil {
		return nil, err
	}
	return &types.MsgRequestDataResponse{}, nil
}

func (k msgServer) ReportData(goCtx context.Context, msg *types.MsgReportData) (*types.MsgReportDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	maxReportDataSize := int(k.GetParams(ctx).MaxReportDataSize)
	for _, r := range msg.RawReports {
		if len(r.Data) > maxReportDataSize {
			return nil, types.WrapMaxError(types.ErrTooLargeRawReportData, len(r.Data), maxReportDataSize)
		}
	}

	validator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}

	// check request must not expire.
	if msg.RequestID <= k.GetRequestLastExpired(ctx) {
		return nil, types.ErrRequestAlreadyExpired
	}

	reportInTime := !k.HasResult(ctx, msg.RequestID)
	err = k.AddReport(ctx, msg.RequestID, validator, reportInTime, msg.RawReports)
	if err != nil {
		return nil, err
	}

	// if request has not been resolved, check if it need to resolve at the endblock
	if reportInTime {
		req := k.MustGetRequest(ctx, msg.RequestID)
		if k.GetReportCount(ctx, msg.RequestID) == req.MinCount {
			// at this moment we are sure, that all the raw reports here are validated
			// so we can distribute the reward for them in end-block
			if _, err := k.CollectReward(ctx, msg.GetRawReports(), req.RawRequests); err != nil {
				return nil, err
			}
			// At the exact moment when the number of reports is sufficient, we add the request to
			// the pending resolve list. This can happen at most one time for any request.
			k.AddPendingRequest(ctx, msg.RequestID)
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeReport,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", msg.RequestID)),
		sdk.NewAttribute(types.AttributeKeyValidator, validator.String()),
	))
	return &types.MsgReportDataResponse{}, nil
}

func (k msgServer) CreateDataSource(
	goCtx context.Context,
	msg *types.MsgCreateDataSource,
) (*types.MsgCreateDataSourceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// unzip if it's a zip file
	if gzip.IsGzipped(msg.Executable) {
		var err error
		msg.Executable, err = gzip.Uncompress(msg.Executable, types.MaxExecutableSize)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrUncompressionFailed, err.Error())
		}
	}

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	treasury, err := sdk.AccAddressFromBech32(msg.Treasury)
	if err != nil {
		return nil, err
	}

	id := k.AddDataSource(ctx, types.NewDataSource(
		owner, msg.Name, msg.Description, k.AddExecutableFile(msg.Executable), msg.Fee, treasury,
	))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreateDataSource,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
	))

	return &types.MsgCreateDataSourceResponse{}, nil
}

func (k msgServer) EditDataSource(
	goCtx context.Context,
	msg *types.MsgEditDataSource,
) (*types.MsgEditDataSourceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	dataSource, err := k.GetDataSource(ctx, msg.DataSourceID)
	if err != nil {
		return nil, err
	}

	owner, err := sdk.AccAddressFromBech32(dataSource.Owner)
	if err != nil {
		return nil, err
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// sender must be the owner of data source
	if !owner.Equals(sender) {
		return nil, types.ErrEditorNotAuthorized
	}

	treasury, err := sdk.AccAddressFromBech32(msg.Treasury)
	if err != nil {
		return nil, err
	}

	// unzip if it's a zip file
	if gzip.IsGzipped(msg.Executable) {
		msg.Executable, err = gzip.Uncompress(msg.Executable, types.MaxExecutableSize)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrUncompressionFailed, err.Error())
		}
	}

	newOwner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	// Can safely use MustEdit here, as we already checked that the data source exists above.
	k.MustEditDataSource(ctx, msg.DataSourceID, types.NewDataSource(
		newOwner, msg.Name, msg.Description, k.AddExecutableFile(msg.Executable), msg.Fee, treasury,
	))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEditDataSource,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", msg.DataSourceID)),
	))

	return &types.MsgEditDataSourceResponse{}, nil
}

func (k msgServer) CreateOracleScript(
	goCtx context.Context,
	msg *types.MsgCreateOracleScript,
) (*types.MsgCreateOracleScriptResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// unzip if it's a zip file
	if gzip.IsGzipped(msg.Code) {
		var err error
		msg.Code, err = gzip.Uncompress(msg.Code, types.MaxWasmCodeSize)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrUncompressionFailed, err.Error())
		}
	}

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	filename, err := k.AddOracleScriptFile(msg.Code)
	if err != nil {
		return nil, err
	}

	id := k.AddOracleScript(ctx, types.NewOracleScript(
		owner, msg.Name, msg.Description, filename, msg.Schema, msg.SourceCodeURL,
	))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreateOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
	))

	return &types.MsgCreateOracleScriptResponse{}, nil
}

func (k msgServer) EditOracleScript(
	goCtx context.Context,
	msg *types.MsgEditOracleScript,
) (*types.MsgEditOracleScriptResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	oracleScript, err := k.GetOracleScript(ctx, msg.OracleScriptID)
	if err != nil {
		return nil, err
	}

	owner, err := sdk.AccAddressFromBech32(oracleScript.Owner)
	if err != nil {
		return nil, err
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// sender must be the owner of oracle script
	if !owner.Equals(sender) {
		return nil, types.ErrEditorNotAuthorized
	}

	// unzip if it's a zip file
	if gzip.IsGzipped(msg.Code) {
		msg.Code, err = gzip.Uncompress(msg.Code, types.MaxWasmCodeSize)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrUncompressionFailed, err.Error())
		}
	}

	filename, err := k.AddOracleScriptFile(msg.Code)
	if err != nil {
		return nil, err
	}

	newOwner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	k.MustEditOracleScript(ctx, msg.OracleScriptID, types.NewOracleScript(
		newOwner, msg.Name, msg.Description, filename, msg.Schema, msg.SourceCodeURL,
	))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEditOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", msg.OracleScriptID)),
	))

	return &types.MsgEditOracleScriptResponse{}, nil
}

func (k msgServer) Activate(goCtx context.Context, msg *types.MsgActivate) (*types.MsgActivateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	err = k.Keeper.Activate(ctx, valAddr)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeActivate,
		sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator),
	))
	return &types.MsgActivateResponse{}, nil
}

func (k msgServer) UpdateParams(
	goCtx context.Context,
	msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			k.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUpdateParams,
		sdk.NewAttribute(types.AttributeKeyParams, msg.Params.String()),
	))

	return &types.MsgUpdateParamsResponse{}, nil
}
