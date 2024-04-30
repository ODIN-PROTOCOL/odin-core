package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// HasRequest checks if the request of this ID exists in the storage.
func (k Keeper) HasRequest(ctx context.Context, id types.RequestID) (bool, error) {
	return k.Requests.Has(ctx, uint64(id))
}

// GetRequest returns the request struct for the given ID or error if not exists.
func (k Keeper) GetRequest(ctx context.Context, id types.RequestID) (types.Request, error) {
	return k.Requests.Get(ctx, uint64(id))
}

// MustGetRequest returns the request struct for the given ID. Panics error if not exists.
func (k Keeper) MustGetRequest(ctx context.Context, id types.RequestID) types.Request {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		panic(err)
	}
	return request
}

// SetRequest saves the given data request to the store without performing any validation.
func (k Keeper) SetRequest(ctx context.Context, id types.RequestID, request types.Request) error {
	return k.Requests.Set(ctx, uint64(id), request)
}

// DeleteRequest removes the given data request from the store.
func (k Keeper) DeleteRequest(ctx context.Context, id types.RequestID) error {
	return k.Requests.Remove(ctx, uint64(id))
}

// AddRequest attempts to create and save a new request.
func (k Keeper) AddRequest(ctx context.Context, req types.Request) (types.RequestID, error) {
	id, err := k.GetNextRequestID(ctx)
	if err != nil {
		return 0, err
	}

	err = k.SetRequest(ctx, id, req)
	return id, err
}

// ProcessExpiredRequests resolves all expired requests and deactivates missed validators.
func (k Keeper) ProcessExpiredRequests(ctx context.Context) error {
	goCtx := sdk.UnwrapSDKContext(ctx)

	currentReqID, err := k.GetRequestLastExpired(ctx)
	if err != nil {
		return err
	}
	currentReqID += 1

	requestCount, err := k.GetRequestCount(ctx)
	if err != nil {
		return err
	}

	lastReqID := types.RequestID(requestCount)

	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	expirationBlockCount := int64(params.ExpirationBlockCount)
	// Loop through all data requests in chronological order. If a request reaches its
	// expiration height, we will deactivate validators that didn't report data on the
	// request. We also resolve requests to status EXPIRED if they are not yet resolved.
	for ; currentReqID <= lastReqID; currentReqID++ {
		req := k.MustGetRequest(ctx, currentReqID)

		// This request is not yet expired, so there's nothing to do here. Ditto for
		// all other requests that come after this. Thus we can just break the loop.
		if req.RequestHeight+expirationBlockCount > goCtx.BlockHeight() {
			break
		}

		hasResult, err := k.HasResult(ctx, currentReqID)
		if err != nil {
			return err
		}

		// If the request still does not have result, we resolve it as EXPIRED.
		if !hasResult {
			err = k.ResolveExpired(ctx, currentReqID)
			if err != nil {
				return err
			}
		}

		// Deactivate all validators that do not report to this request.
		for _, val := range req.RequestedValidators {
			v, _ := sdk.ValAddressFromBech32(val)

			hasReport, err := k.HasReport(ctx, currentReqID, v)
			if err != nil {
				return err
			}

			if !hasReport {
				err = k.MissReport(ctx, v, time.Unix(req.RequestTime, 0))
				if err != nil {
					return err
				}
			}
		}

		// Cleanup request and reports
		err = k.DeleteRequest(ctx, currentReqID)
		if err != nil {
			return err
		}

		err = k.DeleteReports(ctx, currentReqID)
		if err != nil {
			return err
		}

		// Set last expired request ID to be this current request.
		err = k.SetRequestLastExpired(ctx, currentReqID)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddPendingRequest adds the request to the pending list. DO NOT add same request more than once.
func (k Keeper) AddPendingRequest(ctx context.Context, id types.RequestID) error {
	pendingList, err := k.GetPendingResolveList(ctx)
	if err != nil {
		return err
	}

	pendingList = append(pendingList, id)
	return k.SetPendingResolveList(ctx, pendingList)
}

// SetPendingResolveList saves the list of pending request that will be resolved at end block.
func (k Keeper) SetPendingResolveList(ctx context.Context, ids []types.RequestID) error {
	intVs := make([]uint64, len(ids))
	for idx, id := range ids {
		intVs[idx] = uint64(id)
	}

	value := types.PendingResolveList{RequestIds: intVs}

	return k.PendingResolveList.Set(ctx, value)
}

// GetPendingResolveList returns the list of pending requests to be executed during EndBlock.
func (k Keeper) GetPendingResolveList(ctx context.Context) (ids []types.RequestID, err error) {
	pendingResolveList, err := k.PendingResolveList.Get(ctx)
	if err != nil {
		return nil, err
	}

	for _, rid := range pendingResolveList.RequestIds {
		ids = append(ids, types.RequestID(rid))
	}
	return ids, nil
}

// GetPaginatedRequests returns all requests with pagination
func (k Keeper) GetPaginatedRequests(
	ctx context.Context,
	limit, offset uint64, reverse bool,
) ([]types.RequestResult, *query.PageResponse, error) {
	pagination := &query.PageRequest{
		Limit:   limit,
		Offset:  offset,
		Reverse: reverse,
	}

	requests, pageRes, err := query.CollectionPaginate(ctx, k.Results, pagination, func(key uint64, result types.Result) (types.RequestResult, error) {

		request, err := k.GetRequest(ctx, result.RequestID)
		if err != nil {
			lastExpired, err := k.GetRequestLastExpired(ctx)
			if err != nil {
				return types.RequestResult{}, err
			}

			if result.RequestID > lastExpired {
				return types.RequestResult{}, status.Error(codes.NotFound, fmt.Sprintf("unable to get request from chain: request id (%d) > latest expired request id (%d)", result.RequestID, lastExpired))
			}
		}

		reports, err := k.GetReports(ctx, result.RequestID)
		if err != nil {
			return types.RequestResult{}, err
		}

		requestResult := types.RequestResult{
			Request: &request,
			Result:  &result,
			Reports: reports,
		}

		return requestResult, nil
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to paginate requests")
	}

	return requests, pageRes, nil
}
