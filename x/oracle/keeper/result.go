package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

const (
	packetExpireTime = int64(10 * time.Minute)
)

// HasResult checks if the result of this request ID exists in the storage.
func (k Keeper) HasResult(ctx context.Context, id types.RequestID) (bool, error) {
	return k.Results.Has(ctx, uint64(id))
}

// SetResult sets result to the store.
func (k Keeper) SetResult(ctx context.Context, reqID types.RequestID, result types.Result) error {
	return k.Results.Set(ctx, uint64(reqID), result)
}

// GetResult returns the result for the given request ID or error if not exists.
func (k Keeper) GetResult(ctx context.Context, id types.RequestID) (types.Result, error) {
	return k.Results.Get(ctx, uint64(id))
}

// MustGetResult returns the result for the given request ID. Panics on error.
func (k Keeper) MustGetResult(ctx sdk.Context, id types.RequestID) types.Result {
	result, err := k.GetResult(ctx, id)
	if err != nil {
		panic(err)
	}
	return result
}

// ResolveSuccess resolves the given request as success with the given result.
func (k Keeper) ResolveSuccess(ctx sdk.Context, id types.RequestID, result []byte, gasUsed uint64) error {
	err := k.SaveResult(ctx, id, types.RESOLVE_STATUS_SUCCESS, result)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", types.RESOLVE_STATUS_SUCCESS)),
		sdk.NewAttribute(types.AttributeKeyResult, hex.EncodeToString(result)),
		sdk.NewAttribute(types.AttributeKeyGasUsed, fmt.Sprintf("%d", gasUsed)),
	))

	return nil
}

// ResolveFailure resolves the given request as failure with the given reason.
func (k Keeper) ResolveFailure(ctx sdk.Context, id types.RequestID, reason string) error {
	err := k.SaveResult(ctx, id, types.RESOLVE_STATUS_FAILURE, []byte{})
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", types.RESOLVE_STATUS_FAILURE)),
		sdk.NewAttribute(types.AttributeKeyReason, reason),
	))

	return nil
}

// ResolveExpired resolves the given request as expired.
func (k Keeper) ResolveExpired(ctx context.Context, id types.RequestID) error {
	err := k.SaveResult(ctx, id, types.RESOLVE_STATUS_EXPIRED, []byte{})
	if err != nil {
		return err
	}

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", types.RESOLVE_STATUS_EXPIRED)),
	))

	return nil
}

// SaveResult saves the result packets for the request with the given resolve status and result.
func (k Keeper) SaveResult(
	ctx context.Context, id types.RequestID, status types.ResolveStatus, result []byte,
) error {
	goCtx := sdk.UnwrapSDKContext(ctx)
	r := k.MustGetRequest(ctx, id)
	reportCount, err := k.GetReportCount(ctx, id)
	if err != nil {
		return err
	}

	err = k.SetResult(ctx, id, types.NewResult(
		r.ClientID,                         // ClientID
		r.OracleScriptID,                   // OracleScriptID
		r.Calldata,                         // Calldata
		uint64(len(r.RequestedValidators)), // AskCount
		r.MinCount,                         // MinCount
		id,                                 // RequestID
		reportCount,                        // AnsCount
		r.RequestTime,                      // RequestTime
		goCtx.BlockTime().Unix(),           // ResolveTime
		status,                             // ResolveStatus
		result,                             // Result
	))
	if err != nil {
		return err
	}

	if r.IBCChannel != nil {
		sourceChannel := r.IBCChannel.ChannelId
		sourcePort := r.IBCChannel.PortId

		channelCap, ok := k.scopedKeeper.GetCapability(goCtx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
		if !ok {
			goCtx.EventManager().EmitEvent(sdk.NewEvent(
				types.EventTypeSendPacketFail,
				sdk.NewAttribute(types.AttributeKeyReason, "Module does not own channel capability"),
			))
			// TODO: create error?
			return nil
		}

		packetData := types.NewOracleResponsePacketData(
			r.ClientID, id, reportCount, r.RequestTime, goCtx.BlockTime().Unix(), status, result,
		)

		if _, err := k.channelKeeper.SendPacket(
			goCtx,
			channelCap,
			sourcePort,
			sourceChannel,
			clienttypes.NewHeight(0, 0),
			uint64(goCtx.BlockTime().UnixNano()+packetExpireTime),
			packetData.GetBytes(),
		); err != nil {
			goCtx.EventManager().EmitEvent(sdk.NewEvent(
				types.EventTypeSendPacketFail,
				sdk.NewAttribute(types.AttributeKeyReason, fmt.Sprintf("Unable to send packet: %s", err)),
			))
		}
	}

	return nil
}
