package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Counts queries the number of data sources, oracle scripts, and requests.
func (k Querier) Counts(c context.Context, req *types.QueryCountsRequest) (*types.QueryCountsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryCountsResponse{
			DataSourceCount:   k.GetDataSourceCount(ctx),
			OracleScriptCount: k.GetOracleScriptCount(ctx),
			RequestCount:      k.GetRequestCount(ctx)},
		nil
}

// Data queries the data source or oracle script script for given file hash.
func (k Querier) Data(c context.Context, req *types.QueryDataRequest) (*types.QueryDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	data, err := k.fileCache.GetFile(req.DataHash)
	if err != nil {
		return nil, err
	}
	return &types.QueryDataResponse{Data: data}, nil
}

// DataSource queries data source info for given data source id.
func (k Querier) DataSource(
	c context.Context,
	req *types.QueryDataSourceRequest,
) (*types.QueryDataSourceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	ds, err := k.GetDataSource(ctx, types.DataSourceID(req.DataSourceId))
	if err != nil {
		return nil, err
	}
	return &types.QueryDataSourceResponse{DataSource: &ds}, nil
}

// DataSources queries data sources
func (k Querier) DataSources(
	c context.Context,
	req *types.QueryDataSourcesRequest,
) (*types.QueryDataSourcesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	dataSources, pageRes, err := k.GetPaginatedDataSources(ctx, req.Pagination.Limit, req.Pagination.Offset)
	if err != nil {
		return nil, err
	}
	return &types.QueryDataSourcesResponse{DataSources: dataSources, Pagination: pageRes}, nil
}

// OracleScript queries oracle script info for given oracle script id.
func (k Querier) OracleScript(
	c context.Context,
	req *types.QueryOracleScriptRequest,
) (*types.QueryOracleScriptResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	os, err := k.GetOracleScript(ctx, types.OracleScriptID(req.OracleScriptId))
	if err != nil {
		return nil, err
	}
	return &types.QueryOracleScriptResponse{OracleScript: &os}, nil
}

// OracleScripts queries all oracle scripts with pagination.
func (k Querier) OracleScripts(c context.Context, req *types.QueryOracleScriptsRequest) (*types.QueryOracleScriptsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracleScripts, pageRes, err := k.GetPaginatedOracleScripts(ctx, req.Pagination.Limit, req.Pagination.Offset)
	if err != nil {
		return nil, err
	}
	return &types.QueryOracleScriptsResponse{OracleScripts: oracleScripts, Pagination: pageRes}, nil
}

// Request queries request info for given request id.
func (k Querier) Request(c context.Context, req *types.QueryRequestRequest) (*types.QueryRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	rid := types.RequestID(req.RequestId)

	request, err := k.GetRequest(ctx, rid)
	if err != nil {
		lastExpired := k.GetRequestLastExpired(ctx)
		if rid > lastExpired {
			return nil, status.Error(
				codes.NotFound,
				fmt.Sprintf(
					"unable to get request from chain: request id (%d) > latest expired request id (%d)",
					rid,
					lastExpired,
				),
			)
		}
		result := k.MustGetResult(ctx, rid)
		return &types.QueryRequestResponse{Request: nil, Reports: nil, Result: &result}, nil
	}

	reports := k.GetReports(ctx, rid)
	if !k.HasResult(ctx, rid) {
		return &types.QueryRequestResponse{Request: &request, Reports: reports, Result: nil}, nil
	}

	result := k.MustGetResult(ctx, rid)
	return &types.QueryRequestResponse{Request: &request, Reports: reports, Result: &result}, nil
}

// Requests queries all requests with pagination.
func (k Querier) Requests(c context.Context, req *types.QueryRequestsRequest) (*types.QueryRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	requests, pageRes, err := k.GetPaginatedRequests(ctx, req.Pagination.Limit, req.Pagination.Offset, req.Pagination.Reverse)
	if err != nil {
		return nil, err
	}
	return &types.QueryRequestsResponse{Requests: requests, Pagination: pageRes}, nil
}

// RequestReports queries all reports by the giver request id with pagination.
func (k Querier) RequestReports(c context.Context, req *types.QueryRequestReportsRequest) (*types.QueryRequestReportsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	reports, pageRes, err := k.GetPaginatedRequestReports(
		ctx,
		types.RequestID(req.RequestId),
		req.Pagination.Limit,
		req.Pagination.Offset,
	)
	if err != nil {
		return nil, err
	}
	return &types.QueryRequestReportsResponse{Reports: reports, Pagination: pageRes}, nil
}

func (k Querier) PendingRequests(
	c context.Context,
	req *types.QueryPendingRequestsRequest,
) (*types.QueryPendingRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	valAddress, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("unable to parse given validator address: %v", err),
		)
	}

	lastExpired := k.GetRequestLastExpired(ctx)
	requestCount := k.GetRequestCount(ctx)

	var pendingIDs []uint64
	for id := lastExpired + 1; uint64(id) <= requestCount; id++ {
		oracleReq := k.MustGetRequest(ctx, id)

		// If all validators reported on this request, then skip it.
		reports := k.GetReports(ctx, id)
		if len(reports) == len(oracleReq.RequestedValidators) {
			continue
		}

		// Skip if validator hasn't been assigned or has been reported.
		// If the validator isn't in requested validators set, then skip it.
		isInValidatorSet := false
		for _, v := range oracleReq.RequestedValidators {
			val, err := sdk.ValAddressFromBech32(v)
			if err != nil {
				return nil, status.Error(
					codes.Internal,
					fmt.Sprintf("unable to parse validator address in requested validators %v: %v", v, err),
				)
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
				return nil, status.Error(
					codes.Internal,
					fmt.Sprintf("unable to parse validator address in requested validators %v: %v", r.Validator, err),
				)
			}
			if valAddress.Equals(val) {
				reported = true
				break
			}
		}
		if reported {
			continue
		}

		pendingIDs = append(pendingIDs, uint64(id))
	}

	return &types.QueryPendingRequestsResponse{RequestIDs: pendingIDs}, nil
}

// Validator queries oracle info of validator for given validator
// address.
func (k Querier) Validator(
	c context.Context,
	req *types.QueryValidatorRequest,
) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	status := k.GetValidatorStatus(ctx, val)
	if err != nil {
		return nil, err
	}
	return &types.QueryValidatorResponse{Status: &status}, nil
}

// IsReporter queries grant of account on this validator
func (k Querier) IsReporter(
	c context.Context,
	req *types.QueryIsReporterRequest,
) (*types.QueryIsReporterResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	rep, err := sdk.AccAddressFromBech32(req.ReporterAddress)
	if err != nil {
		return nil, err
	}
	return &types.QueryIsReporterResponse{IsReporter: k.Keeper.IsReporter(ctx, val, rep)}, nil
}

// Reporters queries 100 gratees of a given validator address and filter for reporter.
func (k Querier) Reporters(
	c context.Context,
	req *types.QueryReportersRequest,
) (*types.QueryReportersResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	val, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	granter := sdk.AccAddress(val).String()
	granterGrantsRequest := &authz.QueryGranterGrantsRequest{
		Granter: granter,
	}
	granterGrantsRes, err := k.authzKeeper.GranterGrants(c, granterGrantsRequest)
	if err != nil {
		return nil, err
	}
	reporters := make([]string, 0)
	for _, rep := range granterGrantsRes.Grants {
		if rep.Authorization.GetCachedValue().(authz.Authorization).MsgTypeURL() == sdk.MsgTypeURL(
			&types.MsgReportData{},
		) &&
			rep.Expiration.After(ctx.BlockTime()) {
			reporters = append(reporters, rep.Grantee)
		}
	}
	return &types.QueryReportersResponse{Reporter: reporters}, nil
}

// ActiveValidators queries all active oracle validators.
func (k Querier) ActiveValidators(
	c context.Context,
	req *types.QueryActiveValidatorsRequest,
) (*types.QueryActiveValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	result := types.QueryActiveValidatorsResponse{}
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx,
		func(idx int64, val stakingtypes.ValidatorI) (stop bool) {
			if k.GetValidatorStatus(ctx, val.GetOperator()).IsActive {
				result.Validators = append(result.Validators, &types.ActiveValidator{
					Address: val.GetOperator().String(),
					Power:   val.GetTokens().Uint64(),
				})
			}
			return false
		})
	return &result, nil
}

// Params queries the oracle parameters.
func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// RequestSearch queries the latest request that match the given input.
func (k Querier) RequestSearch(
	c context.Context,
	req *types.QueryRequestSearchRequest,
) (*types.QueryRequestSearchResponse, error) {
	return nil, status.Error(codes.Unimplemented, "This feature can be taken from extra/rest branch")
}

// RequestPrice queries the latest price on standard price reference oracle
// script.
func (k Querier) RequestPrice(
	c context.Context,
	req *types.QueryRequestPriceRequest,
) (*types.QueryRequestPriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "This feature can be taken from extra/rest branch")
}

// RequestVerification verifies oracle request for validation before executing data sources
func (k Querier) RequestVerification(
	c context.Context,
	req *types.QueryRequestVerificationRequest,
) (*types.QueryRequestVerificationResponse, error) {
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
		return nil, status.Error(codes.InvalidArgument, "unable to get reporter's public key")
	}
	reporterPubKey := secp256k1.PubKey(pk[:])

	requestVerificationContent := types.NewRequestVerification(
		req.ChainId,
		validator,
		types.RequestID(req.RequestId),
		types.ExternalID(req.ExternalId),
		types.DataSourceID(req.DataSourceId),
	)
	signByte := requestVerificationContent.GetSignBytes()
	if !reporterPubKey.VerifySignature(signByte, req.Signature) {
		return nil, status.Error(codes.Unauthenticated, "invalid reporter's signature")
	}

	// Provided reporter should be authorized by the provided validator
	reporter := sdk.AccAddress(reporterPubKey.Address().Bytes())
	if !k.Keeper.IsReporter(ctx, validator, reporter) {
		return nil, status.Error(
			codes.PermissionDenied,
			fmt.Sprintf("%s is not an authorized reporter of %s", reporter, req.Validator),
		)
	}

	// Provided request should exist on chain
	request, err := k.GetRequest(ctx, types.RequestID(req.RequestId))
	if err != nil {
		// return uncertain result if request id is in range of max delay
		if req.RequestId-k.GetRequestCount(ctx) > 0 && req.RequestId-k.GetRequestCount(ctx) <= req.MaxDelay {
			return &types.QueryRequestVerificationResponse{
				ChainId:      req.ChainId,
				Validator:    req.Validator,
				RequestId:    req.RequestId,
				ExternalId:   req.ExternalId,
				DataSourceId: req.DataSourceId,
				IsDelay:      true,
			}, nil
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to get request from chain: %s", err.Error()))
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
	var dataSourceID *types.DataSourceID
	for _, rawRequest := range request.RawRequests {
		if rawRequest.ExternalID == types.ExternalID(req.ExternalId) {
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
	if *dataSourceID != types.DataSourceID(req.DataSourceId) {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf(
				"data source required by the request %d which relates to the external data source with ID %d is not match with data source id provided in request.",
				req.RequestId,
				req.ExternalId,
			),
		)
	}

	// Provided validator should not have reported data for the request
	reports := k.GetReports(ctx, types.RequestID(req.RequestId))
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

	// The request should not be expired
	if request.RequestHeight+int64(k.GetParams(ctx).ExpirationBlockCount) < ctx.BlockHeader().Height {
		return nil, status.Error(
			codes.DeadlineExceeded,
			fmt.Sprintf("Request with ID %d is already expired", req.RequestId),
		)
	}

	return &types.QueryRequestVerificationResponse{
		ChainId:      req.ChainId,
		Validator:    req.Validator,
		RequestId:    req.RequestId,
		ExternalId:   req.ExternalId,
		DataSourceId: uint64(*dataSourceID),
		IsDelay:      false,
	}, nil
}

func (k Querier) DataProvidersPool(c context.Context, req *types.QueryDataProvidersPoolRequest) (*types.QueryDataProvidersPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryDataProvidersPoolResponse{
		Pool: k.GetOraclePool(ctx).DataProvidersPool,
	}, nil
}

// DataProviderReward returns current reward per byte for data providers
func (k Querier) DataProviderReward(
	c context.Context, _ *types.QueryDataProviderRewardRequest,
) (*types.QueryDataProviderRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	accumulatedRewards := k.GetAccumulatedDataProvidersRewards(ctx)
	return &types.QueryDataProviderRewardResponse{RewardPerByte: accumulatedRewards.CurrentRewardPerByte}, nil
}

// DataProviderAccumulatedReward queries reward of a given data provider address.
func (k Querier) DataProviderAccumulatedReward(c context.Context, req *types.QueryDataProviderAccumulatedRewardRequest) (*types.QueryDataProviderAccumulatedRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.DataProviderAddress)
	if err != nil {
		return nil, err
	}
	accumulatedReward := k.GetDataProviderAccumulatedReward(ctx, addr)
	return &types.QueryDataProviderAccumulatedRewardResponse{AccumulatedReward: accumulatedReward}, nil
}
