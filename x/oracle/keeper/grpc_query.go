package oraclekeeper

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ oracletypes.QueryServer = Querier{}

// Counts queries the number of data sources, oracle scripts, and requests.
func (k Querier) Counts(c context.Context, req *oracletypes.QueryCountsRequest) (*oracletypes.QueryCountsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &oracletypes.QueryCountsResponse{
			DataSourceCount:   k.GetDataSourceCount(ctx),
			OracleScriptCount: k.GetOracleScriptCount(ctx),
			RequestCount:      k.GetRequestCount(ctx),
		},
		nil
}

// Data queries the data source or oracle script script for given file hash.
func (k Querier) Data(c context.Context, req *oracletypes.QueryDataRequest) (*oracletypes.QueryDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	data, err := k.fileCache.GetFile(req.DataHash)
	if err != nil {
		return nil, err
	}
	return &oracletypes.QueryDataResponse{Data: data}, nil
}

// DataSource queries data source info for given data source id.
func (k Querier) DataSource(c context.Context, req *oracletypes.QueryDataSourceRequest) (*oracletypes.QueryDataSourceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	ds, err := k.GetDataSource(ctx, oracletypes.DataSourceID(req.DataSourceId))
	if err != nil {
		return nil, err
	}
	return &oracletypes.QueryDataSourceResponse{DataSource: &ds}, nil
}

// DataSources queries data sources
func (k Querier) DataSources(c context.Context, req *oracletypes.QueryDataSourcesRequest) (*oracletypes.QueryDataSourcesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	dataSources, pageRes, err := k.GetPaginatedDataSources(ctx, req.Pagination.Limit, req.Pagination.Offset)
	if err != nil {
		return nil, err
	}
	return &oracletypes.QueryDataSourcesResponse{DataSources: dataSources, Pagination: pageRes}, nil
}

// OracleScript queries oracle script info for given oracle script id.
func (k Querier) OracleScript(c context.Context, req *oracletypes.QueryOracleScriptRequest) (*oracletypes.QueryOracleScriptResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	os, err := k.GetOracleScript(ctx, oracletypes.OracleScriptID(req.OracleScriptId))
	if err != nil {
		return nil, err
	}
	return &oracletypes.QueryOracleScriptResponse{OracleScript: &os}, nil
}

// OracleScripts queries all oracle scripts with pagination.
func (k Querier) OracleScripts(c context.Context, req *oracletypes.QueryOracleScriptsRequest) (*oracletypes.QueryOracleScriptsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracleScripts, pageRes, err := k.GetPaginatedOracleScripts(ctx, req.Pagination.Limit, req.Pagination.Offset)
	if err != nil {
		return nil, err
	}
	return &oracletypes.QueryOracleScriptsResponse{OracleScripts: oracleScripts, Pagination: pageRes}, nil
}

// Request queries request info for given request id.
func (k Querier) Request(c context.Context, req *oracletypes.QueryRequestRequest) (*oracletypes.QueryRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	rid := oracletypes.RequestID(req.RequestId)

	request, err := k.GetRequest(ctx, rid)
	if err != nil {
		lastExpired := k.GetRequestLastExpired(ctx)
		if rid > lastExpired {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to get request from chain: request id (%d) > latest expired request id (%d)", rid, lastExpired))
		}
		result := k.MustGetResult(ctx, rid)
		return &oracletypes.QueryRequestResponse{Request: nil, Reports: nil, Result: &result}, nil
	}

	result, err := k.GetResult(ctx, rid)
	if err != nil {
		return nil, err
	}

	reports := k.GetRequestReports(ctx, rid)

	return &oracletypes.QueryRequestResponse{Request: &request, Result: &result, Reports: reports}, nil
}

// Requests queries all requests with pagination.
func (k Querier) Requests(c context.Context, req *oracletypes.QueryRequestsRequest) (*oracletypes.QueryRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	requests, pageRes, err := k.GetPaginatedRequests(ctx, req.Pagination.Limit, req.Pagination.Offset, req.Pagination.Reverse)
	if err != nil {
		return nil, err
	}
	return &oracletypes.QueryRequestsResponse{Requests: requests, Pagination: pageRes}, nil
}

// RequestReports queries all reports by the giver request id with pagination.
func (k Querier) RequestReports(c context.Context, req *oracletypes.QueryRequestReportsRequest) (*oracletypes.QueryRequestReportsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	reports, pageRes, err := k.GetPaginatedRequestReports(
		ctx,
		oracletypes.RequestID(req.RequestId),
		req.Pagination.Limit,
		req.Pagination.Offset,
	)
	if err != nil {
		return nil, err
	}
	return &oracletypes.QueryRequestReportsResponse{Reports: reports, Pagination: pageRes}, nil
}

// Validator queries oracle info of validator for given validator
// address.
func (k Querier) Validator(c context.Context, req *oracletypes.QueryValidatorRequest) (*oracletypes.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	validatorStatus := k.GetValidatorStatus(ctx, val)
	return &oracletypes.QueryValidatorResponse{Status: &validatorStatus}, nil
}

// IsReporter queries grant of account on this validator
func (k Querier) IsReporter(c context.Context, req *oracletypes.QueryIsReporterRequest) (*oracletypes.QueryIsReporterResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	rep, err := sdk.AccAddressFromBech32(req.ReporterAddress)
	if err != nil {
		return nil, err
	}
	return &oracletypes.QueryIsReporterResponse{IsReporter: k.Keeper.IsReporter(ctx, val, rep)}, nil
}

// Reporters queries all reporters of a given validator address.
func (k Querier) Reporters(c context.Context, req *oracletypes.QueryReportersRequest) (*oracletypes.QueryReportersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	reps := k.GetReporters(ctx, val)
	reporters := make([]string, len(reps))
	for idx, rep := range reps {
		reporters[idx] = rep.String()
	}
	return &oracletypes.QueryReportersResponse{Reporter: reporters}, nil
}

// ActiveValidators queries all active oracle validators.
func (k Querier) ActiveValidators(c context.Context, req *oracletypes.QueryActiveValidatorsRequest) (*oracletypes.QueryActiveValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	var vals []oracletypes.QueryActiveValidatorResult
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx,
		func(idx int64, val stakingtypes.ValidatorI) (stop bool) {
			if k.GetValidatorStatus(ctx, val.GetOperator()).IsActive {
				vals = append(vals, oracletypes.QueryActiveValidatorResult{
					Address: val.GetOperator(),
					Power:   val.GetTokens().Uint64(),
				})
			}
			return false
		})
	return &oracletypes.QueryActiveValidatorsResponse{Count: int64(len(vals))}, nil
}

// Params queries the oracle parameters.
func (k Querier) Params(c context.Context, req *oracletypes.QueryParamsRequest) (*oracletypes.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &oracletypes.QueryParamsResponse{Params: params}, nil
}

// TODO: drop or change
// RequestSearch queries the latest request that matches the given input.
func (k Querier) RequestSearch(c context.Context, req *oracletypes.QueryRequestSearchRequest) (*oracletypes.QueryRequestSearchResponse, error) {
	// TODO: revisit, maybe find another way
	//var clientCtx client.Context
	//rawClientCtx := c.Value(client.ClientContextKey)
	//if rawClientCtx != nil {
	//	clientCtx = *rawClientCtx.(*client.Context)
	//} else {
	//	// SHOULD NEVER HIT
	//	panic("client ctx is empty")
	//}
	//clientCtx := client.Context{}
	//
	//resp, _, err := oracleclientcommon.QuerySearchLatestRequest(oracletypes.QuerierRoute, clientCtx, req)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if resp == nil {
	//	return &oracletypes.QueryRequestSearchResponse{}, nil
	//}

	return nil, nil
}

// TODO:
// RequestPrice queries the latest price on standard price reference oracle script.
func (k Querier) RequestPrice(c context.Context, req *oracletypes.QueryRequestPriceRequest) (*oracletypes.QueryRequestPriceResponse, error) {
	return &oracletypes.QueryRequestPriceResponse{}, nil
}

func (k Querier) DataProvidersPool(c context.Context, req *oracletypes.QueryDataProvidersPoolRequest) (*oracletypes.QueryDataProvidersPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &oracletypes.QueryDataProvidersPoolResponse{
		Pool: k.GetOraclePool(ctx).DataProvidersPool,
	}, nil
}

// DataProviderReward returns current reward per byte for data providers
func (k Querier) DataProviderReward(
	c context.Context, _ *oracletypes.QueryDataProviderRewardRequest,
) (*oracletypes.QueryDataProviderRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	accumulatedRewards := k.GetAccumulatedDataProvidersRewards(ctx)
	return &oracletypes.QueryDataProviderRewardResponse{RewardPerByte: accumulatedRewards.CurrentRewardPerByte}, nil
}

func (k Querier) PendingRequests(c context.Context, req *oracletypes.QueryPendingRequestsRequest) (*oracletypes.QueryPendingRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	valAddress, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unable to parse given validator address: %v", err))
	}

	lastExpired := k.GetRequestLastExpired(ctx)
	requestCount := k.GetRequestCount(ctx)

	var pendingIDs []int64
	for id := lastExpired + 1; int64(id) <= requestCount; id++ {
		oracleReq := k.MustGetRequest(ctx, id)

		// If all validators reported on this request, then skip it.
		reports := k.GetRequestReports(ctx, id)
		if len(reports) == len(oracleReq.RequestedValidators) {
			continue
		}

		// Skip if validator hasn't been assigned or has been reported.
		// If the validator isn't in requested validators set, then skip it.
		isInValidatorSet := false
		for _, v := range oracleReq.RequestedValidators {
			val, err := sdk.ValAddressFromBech32(v)
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("unable to parse validator address in requested validators %v: %v", v, err))
			}
			if valAddress.Equals(val) {
				isInValidatorSet = true
				break
			}
		}
		if !isInValidatorSet {
			continue
		}

		// If the validator has reported, then skip it.
		reported := false
		for _, r := range reports {
			val, err := sdk.ValAddressFromBech32(r.Validator)
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("unable to parse validator address in requested validators %v: %v", r.Validator, err))
			}
			if valAddress.Equals(val) {
				reported = true
				break
			}
		}
		if reported {
			continue
		}

		pendingIDs = append(pendingIDs, int64(id))
	}

	return &oracletypes.QueryPendingRequestsResponse{RequestIDs: pendingIDs}, nil
}

// RequestVerification verifies oracle request for validation before executing data sources
func (k Querier) RequestVerification(
	c context.Context,
	req *oracletypes.QueryRequestVerificationRequest,
) (*oracletypes.QueryRequestVerificationResponse, error) {
	// Request should not be empty
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// Provided chain ID should match current chain ID
	if ctx.ChainID() != req.ChainId {
		return nil, status.Error(
			codes.FailedPrecondition,
			fmt.Sprintf(
				"provided chain ID does not match the validator's chain ID; expected %s, got %s",
				ctx.ChainID(),
				req.ChainId,
			),
		)
	}

	// Provided validator's address should be valid
	validator, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("unable to parse validator address: %s", err.Error()),
		)
	}

	// Provided signature should be valid, which means this query request should be signed by the provided reporter
	pk, err := hex.DecodeString(req.Reporter)
	if err != nil || len(pk) != secp256k1.PubKeySize {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unable to get reporter's public key"))
	}
	reporterPubKey := secp256k1.PubKey(pk[:])

	requestVerificationContent := oracletypes.NewRequestVerification(
		req.ChainId, validator,
		oracletypes.RequestID(req.RequestId),
		oracletypes.ExternalID(req.ExternalId),
	)
	signByte := requestVerificationContent.GetSignBytes()
	if !reporterPubKey.VerifySignature(signByte, req.Signature) {
		return nil, status.Error(codes.Unauthenticated, "invalid reporter's signature")
	}

	// Provided reporter should be authorized by the provided validator
	reporter := sdk.AccAddress(reporterPubKey.Address().Bytes())
	if !k.Keeper.IsReporter(ctx, validator, reporter) {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("%s is not an authorized reporter of %s", reporter, req.Validator))
	}

	// Provided request should exist on chain
	request, err := k.GetRequest(ctx, oracletypes.RequestID(req.RequestId))
	if err != nil {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("unable to get request from chain: %s", err.Error()),
		)
	}

	// Provided validator should be assigned to response to the request
	isValidatorAssigned := false
	for _, requestedValidator := range request.RequestedValidators {
		v, _ := sdk.ValAddressFromBech32(requestedValidator)
		if validator.Equals(v) {
			isValidatorAssigned = true
			break
		}
	}
	if !isValidatorAssigned {
		return nil, status.Error(
			codes.PermissionDenied,
			fmt.Sprintf("%s is not assigned for request ID %d", validator, req.RequestId),
		)
	}

	// Provided external ID should be required by the request determined by oracle script
	var dataSourceID *oracletypes.DataSourceID
	for _, rawRequest := range request.RawRequests {
		if rawRequest.ExternalID == oracletypes.ExternalID(req.ExternalId) {
			dataSourceID = &rawRequest.DataSourceID
			break
		}
	}
	if dataSourceID == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf(
				"no data source required by the request %d found which relates to the external data source with ID %d.",
				req.RequestId,
				req.ExternalId,
			),
		)
	}

	// Provided validator should not have reported data for the request
	reports := k.GetRequestReports(ctx, oracletypes.RequestID(req.RequestId))
	isValidatorReported := false
	for _, report := range reports {
		reportVal, _ := sdk.ValAddressFromBech32(report.Validator)
		if reportVal.Equals(validator) {
			isValidatorReported = true
			break
		}
	}
	if isValidatorReported {
		return nil, status.Error(
			codes.AlreadyExists,
			fmt.Sprintf("validator %s already submitted data report for this request", validator),
		)
	}

	params := k.GetParams(ctx)

	// The request should not be expired
	if request.RequestHeight+int64(params.ExpirationBlockCount) < ctx.BlockHeader().Height {
		return nil, status.Error(
			codes.DeadlineExceeded,
			fmt.Sprintf("Request with ID %d is already expired", req.RequestId),
		)
	}

	return &oracletypes.QueryRequestVerificationResponse{
		ChainId:      req.ChainId,
		Validator:    req.Validator,
		RequestId:    req.RequestId,
		ExternalId:   req.ExternalId,
		DataSourceId: uint64(*dataSourceID),
	}, nil
}

// DataProviderAccumulatedReward queries reward of a given data provider address.
func (k Querier) DataProviderAccumulatedReward(c context.Context, req *oracletypes.QueryDataProviderAccumulatedRewardRequest) (*oracletypes.QueryDataProviderAccumulatedRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.DataProviderAddress)
	if err != nil {
		return nil, err
	}
	accumulatedReward := k.GetDataProviderAccumulatedReward(ctx, addr)
	return &oracletypes.QueryDataProviderAccumulatedRewardResponse{AccumulatedReward: accumulatedReward}, nil
}
