package keeper

import (
	"context"
	telemetrytypes "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ telemetrytypes.QueryServer = Keeper{}

func (k Keeper) ValidatorByConsAddr(c context.Context, request *telemetrytypes.QueryValidatorByConsAddrRequest) (*telemetrytypes.QueryValidatorByConsAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	consAddr, err := sdk.ConsAddressFromHex(request.ConsensusAddress)
	if err != nil {
		return nil, err
	}
	val, ok := k.stakingQuerier.Keeper.GetValidatorByConsAddr(ctx, consAddr)
	if !ok {
		return nil, sdkerrors.ErrNotFound
	}

	return &telemetrytypes.QueryValidatorByConsAddrResponse{Validator: val}, nil
}

func (k Keeper) TopBalances(
	c context.Context,
	request *telemetrytypes.QueryTopBalancesRequest,
) (*telemetrytypes.QueryTopBalancesResponse, error) {

	ctx := sdk.UnwrapSDKContext(c)
	balances, total := k.GetPaginatedBalances(ctx, request.GetDenom(), request.GetDesc(), request.Pagination)
	// TODO: optimize or remove
	//txsCount, err := k.GetAccountTxsCount(BalancesToAccounts(balances)...)
	//if err != nil {
	//	return nil, sdkerrors.Wrap(err, "failed to get accounts txs")
	//}
	return &telemetrytypes.QueryTopBalancesResponse{
		//TransactionsCount: txsCount,
		Balances: balances,
		Pagination: &query.PageResponse{
			Total: total,
		},
	}, nil
}

func (k Keeper) AvgBlockSize(
	_ context.Context,
	request *telemetrytypes.QueryAvgBlockSizeRequest,
) (*telemetrytypes.QueryAvgBlockSizeResponse, error) {

	blockSizePerDay, err := k.GetAvgBlockSizePerDay(request.GetStartDate(), request.GetEndDate())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get average block size per day")
	}

	return &telemetrytypes.QueryAvgBlockSizeResponse{
		AvgBlockSizePerDay: blockSizePerDay,
	}, nil
}

func (k Keeper) AvgBlockTime(
	_ context.Context,
	request *telemetrytypes.QueryAvgBlockTimeRequest,
) (*telemetrytypes.QueryAvgBlockTimeResponse, error) {

	blockTimePerDay, err := k.GetAvgBlockTimePerDay(request.GetStartDate(), request.GetEndDate())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get average block time per day")
	}

	return &telemetrytypes.QueryAvgBlockTimeResponse{
		AvgBlockTimePerDay: blockTimePerDay,
	}, nil
}

func (k Keeper) AvgTxFee(
	c context.Context,
	request *telemetrytypes.QueryAvgTxFeeRequest,
) (*telemetrytypes.QueryAvgTxFeeResponse, error) {

	avgTxFee, err := k.GetAvgTxFeePerDay(request.GetStartDate(), request.GetEndDate())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get average tx fee per day")
	}

	return &telemetrytypes.QueryAvgTxFeeResponse{
		AvgTxFeePerDay: avgTxFee,
	}, nil
}

func (k Keeper) TxVolume(
	c context.Context,
	request *telemetrytypes.QueryTxVolumeRequest,
) (*telemetrytypes.QueryTxVolumeResponse, error) {

	txVolume, err := k.GetTxVolumePerDay(request.GetStartDate(), request.GetEndDate())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get tx volume")
	}

	return &telemetrytypes.QueryTxVolumeResponse{
		TxVolumePerDay: txVolume,
	}, nil
}

func (k Keeper) ExtendedValidators(
	c context.Context,
	request *telemetrytypes.QueryExtendedValidatorsRequest,
) (*telemetrytypes.QueryExtendedValidatorsResponse, error) {

	ctx := sdk.UnwrapSDKContext(c)
	validatorsResp, err := k.stakingQuerier.Validators(c, ExtendedValidatorsRequestToValidatorsRequest(request))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get validators")
	}
	accounts, err := ValidatorsToAccounts(validatorsResp.GetValidators())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get validators accounts addresses")
	}
	extendedValidatorsResp := ValidatorsResponseToExtendedValidatorsResponse(validatorsResp)
	extendedValidatorsResp.Balances = k.GetBalances(ctx, accounts...)
	return extendedValidatorsResp, nil
}

func (k Keeper) ValidatorBlocks(
	c context.Context,
	request *telemetrytypes.QueryValidatorBlocksRequest,
) (*telemetrytypes.QueryValidatorBlocksResponse, error) {

	if request.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	address, err := sdk.ValAddressFromBech32(request.ValidatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	blocks, total, err := k.GetValidatorBlocks(
		ctx,
		address,
		request.GetDesc(),
		request.GetPagination(),
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get validator blocks")
	}

	return &telemetrytypes.QueryValidatorBlocksResponse{
		Blocks: blocks,
		Pagination: &query.PageResponse{
			Total: total,
		},
	}, nil
}

func (k Keeper) TopValidators(
	c context.Context,
	request *telemetrytypes.QueryTopValidatorsRequest,
) (*telemetrytypes.QueryTopValidatorsResponse, error) {

	ctx := sdk.UnwrapSDKContext(c)
	topValidators, total, err := k.GetTopValidatorsByBlocks(
		ctx,
		request.GetStartDate(),
		request.GetEndDate(),
		request.GetDesc(),
		request.GetPagination(),
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get top validators by blocks")
	}

	return &telemetrytypes.QueryTopValidatorsResponse{
		TopValidators: topValidators,
		Pagination: &query.PageResponse{
			Total: total,
		},
	}, nil
}
