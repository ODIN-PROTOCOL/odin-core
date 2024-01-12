package oraclekeeper

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/slashing/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

const (
	packetExpireTime = int64(10 * time.Minute)
)

// HasResult checks if the result of this request ID exists in the storage.
func (k Keeper) HasResult(ctx sdk.Context, id oracletypes.RequestID) bool {
	return ctx.KVStore(k.storeKey).Has(oracletypes.ResultStoreKey(id))
}

// SetResult sets result to the store.
func (k Keeper) SetResult(ctx sdk.Context, reqID oracletypes.RequestID, result oracletypes.Result) {
	ctx.KVStore(k.storeKey).Set(oracletypes.ResultStoreKey(reqID), k.cdc.MustMarshal(&result))
}

// GetResult returns the result for the given request ID or error if not exists.
func (k Keeper) GetResult(ctx sdk.Context, id oracletypes.RequestID) (oracletypes.Result, error) {
	bz := ctx.KVStore(k.storeKey).Get(oracletypes.ResultStoreKey(id))
	if bz == nil {
		return oracletypes.Result{}, sdkerrors.Wrapf(oracletypes.ErrResultNotFound, "id: %d", id)
	}
	var result oracletypes.Result
	k.cdc.MustUnmarshal(bz, &result)
	return result, nil
}

// MustGetResult returns the result for the given request ID. Panics on error.
func (k Keeper) MustGetResult(ctx sdk.Context, id oracletypes.RequestID) oracletypes.Result {
	result, err := k.GetResult(ctx, id)
	if err != nil {
		panic(err)
	}
	return result
}

// ResolveSuccess resolves the given request as success with the given result.
func (k Keeper) ResolveSuccess(ctx sdk.Context, id oracletypes.RequestID, result []byte, gasUsed uint64) {
	k.SaveResult(ctx, id, oracletypes.RESOLVE_STATUS_SUCCESS, result)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		oracletypes.EventTypeResolve,
		sdk.NewAttribute(oracletypes.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(oracletypes.AttributeKeyResolveStatus, fmt.Sprintf("%d", oracletypes.RESOLVE_STATUS_SUCCESS)),
		sdk.NewAttribute(oracletypes.AttributeKeyResult, hex.EncodeToString(result)),
		sdk.NewAttribute(oracletypes.AttributeKeyGasUsed, fmt.Sprintf("%d", gasUsed)),
	))
}

// ResolveFailure resolves the given request as failure with the given reason.
func (k Keeper) ResolveFailure(ctx sdk.Context, id oracletypes.RequestID, reason string) {
	k.SaveResult(ctx, id, oracletypes.RESOLVE_STATUS_FAILURE, []byte{})
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		oracletypes.EventTypeResolve,
		sdk.NewAttribute(oracletypes.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(oracletypes.AttributeKeyResolveStatus, fmt.Sprintf("%d", oracletypes.RESOLVE_STATUS_FAILURE)),
		sdk.NewAttribute(oracletypes.AttributeKeyReason, reason),
	))
}

// ResolveExpired resolves the given request as expired.
func (k Keeper) ResolveExpired(ctx sdk.Context, id oracletypes.RequestID) {
	k.SaveResult(ctx, id, oracletypes.RESOLVE_STATUS_EXPIRED, []byte{})
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		oracletypes.EventTypeResolve,
		sdk.NewAttribute(oracletypes.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(oracletypes.AttributeKeyResolveStatus, fmt.Sprintf("%d", oracletypes.RESOLVE_STATUS_EXPIRED)),
	))
}

// SaveResult saves the result packets for the request with the given resolve status and result.
func (k Keeper) SaveResult(
	ctx sdk.Context, id oracletypes.RequestID, status oracletypes.ResolveStatus, result []byte,
) {
	r := k.MustGetRequest(ctx, id)
	reportCount := k.GetReportCount(ctx, id)
	k.SetResult(ctx, id, oracletypes.NewResult(
		r.ClientID,                         // ClientID
		r.OracleScriptID,                   // OracleScriptID
		r.Calldata,                         // Calldata
		uint64(len(r.RequestedValidators)), // AskCount
		r.MinCount,                         // MinCount
		id,                                 // RequestID
		reportCount,                        // AnsCount
		int64(r.RequestTime),               // RequestTime
		ctx.BlockTime().Unix(),             // ResolveTime
		status,                             // ResolveStatus
		result,                             // Result
	))

	if r.IBCChannel != nil {
		sourceChannel := r.IBCChannel.ChannelId
		sourcePort := r.IBCChannel.PortId

		channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
		if !ok {
			ctx.EventManager().EmitEvent(sdk.NewEvent(
				oracletypes.EventTypeSendPacketFail,
				sdk.NewAttribute(oracletypes.AttributeKeyReason, "Module does not own channel capability"),
			))
			return
		}

		packetData := oracletypes.NewOracleResponsePacketData(
			r.ClientID, id, reportCount, int64(r.RequestTime), ctx.BlockTime().Unix(), status, result,
		)

		if _, err := k.channelKeeper.SendPacket(
			ctx,
			channelCap,
			sourcePort,
			sourceChannel,
			clienttypes.NewHeight(0, 0),
			uint64(ctx.BlockTime().UnixNano()+packetExpireTime),
			packetData.GetBytes(),
		); err != nil {
			ctx.EventManager().EmitEvent(sdk.NewEvent(
				oracletypes.EventTypeSendPacketFail,
				sdk.NewAttribute(types.AttributeKeyReason, fmt.Sprintf("Unable to send packet: %s", err)),
			))
		}
	}
}
