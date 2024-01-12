package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	commontypes "github.com/ODIN-PROTOCOL/odin-core/x/common/types"
	telemetrytypes "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/types"
)

func NewQuerier(keeper Keeper, cdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case telemetrytypes.QueryTopBalances:
			return queryTopBalances(ctx, path[1:], keeper, cdc, req)
		case telemetrytypes.QueryExtendedValidators:
			return queryExtendedValidators(ctx, path[1:], keeper, cdc, req)
		/*case telemetrytypes.QueryAvgBlockSize:
			return queryAvgBlockSize(ctx, path[1:], keeper, cdc, req)
		case telemetrytypes.QueryAvgBlockTime:
			return queryAvgBlockTime(ctx, path[1:], keeper, cdc, req)
		case telemetrytypes.QueryAvgTxFee:
			return queryAvgTxFee(ctx, path[1:], keeper, cdc, req)
		case telemetrytypes.QueryTxVolume:
			return queryTxVolume(ctx, path[1:], keeper, cdc, req)*/
		case telemetrytypes.QueryValidatorBlocks:
			return queryValidatorBlocks(ctx, path[1:], keeper, cdc, req)
		case telemetrytypes.QueryTopValidators:
			return queryTopValidators(ctx, path[1:], keeper, cdc, req)
		case telemetrytypes.QueryValidatorByConsAddress:
			return queryValidatorByConsAddr(ctx, path[1:], keeper, cdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown telemetry query endpoint")
		}
	}
}

func queryTopBalances(
	ctx sdk.Context,
	path []string,
	k Keeper,
	cdc *codec.LegacyAmino,
	req abci.RequestQuery,
) ([]byte, error) {
	if len(path) > 1 {
		return nil, sdkerrors.ErrInvalidRequest
	}
	var params commontypes.QueryPaginationParams
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	balances, total := k.GetPaginatedBalances(ctx, path[0], params.Desc, &query.PageRequest{
		Offset: params.Offset,
		Limit:  params.Limit,
	})
	return commontypes.QueryOK(cdc, telemetrytypes.QueryTopBalancesResponse{
		Balances: balances,
		Pagination: &query.PageResponse{
			Total: total,
		},
	})
}

func queryExtendedValidators(
	ctx sdk.Context,
	path []string,
	k Keeper,
	cdc *codec.LegacyAmino,
	req abci.RequestQuery,
) ([]byte, error) {
	if len(path) > 1 {
		return nil, sdkerrors.ErrInvalidRequest
	}
	var params commontypes.QueryPaginationParams
	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal query pagination params")
	}

	total := 0
	if params.GetCountTotal() {
		total = len(k.stakingQuerier.GetValidators(ctx, k.stakingQuerier.MaxValidators(ctx)))
	}
	validators, err := k.ExtendedValidators(sdk.WrapSDKContext(ctx), &telemetrytypes.QueryExtendedValidatorsRequest{
		Status: path[0],
		Pagination: &query.PageRequest{
			Offset: params.Offset,
			Limit:  params.Limit,
		},
	})
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to query extended validators")
	}

	validators.Pagination.Total = uint64(total)
	return commontypes.QueryOK(cdc, validators)
}

func queryAvgBlockSize(
	_ sdk.Context,
	_ []string,
	k Keeper,
	cdc *codec.LegacyAmino,
	req abci.RequestQuery,
) ([]byte, error) {
	var request telemetrytypes.QueryAvgBlockSizeRequest
	if err := cdc.UnmarshalJSON(req.Data, &request); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	blockSizePerDay, err := k.GetAvgBlockSizePerDay(request.GetStartDate(), request.GetEndDate())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get average block size per day")
	}
	return commontypes.QueryOK(cdc, telemetrytypes.QueryAvgBlockSizeResponse{
		AvgBlockSizePerDay: blockSizePerDay,
	})
}

func queryAvgBlockTime(
	_ sdk.Context,
	_ []string,
	k Keeper,
	cdc *codec.LegacyAmino,
	req abci.RequestQuery,
) ([]byte, error) {
	var request telemetrytypes.QueryAvgBlockTimeRequest
	if err := cdc.UnmarshalJSON(req.Data, &request); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	blockTimePerDay, err := k.GetAvgBlockTimePerDay(request.GetStartDate(), request.GetEndDate())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get average block time per day")
	}

	return commontypes.QueryOK(cdc, telemetrytypes.QueryAvgBlockTimeResponse{
		AvgBlockTimePerDay: blockTimePerDay,
	})
}

func queryAvgTxFee(_ sdk.Context, _ []string, k Keeper, cdc *codec.LegacyAmino, req abci.RequestQuery) ([]byte, error) {
	var request telemetrytypes.QueryAvgTxFeeRequest
	if err := cdc.UnmarshalJSON(req.Data, &request); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	avgTxFee, err := k.GetAvgTxFeePerDay(request.GetStartDate(), request.GetEndDate())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get average tx fee per day")
	}
	return commontypes.QueryOK(cdc, telemetrytypes.QueryAvgTxFeeResponse{
		AvgTxFeePerDay: avgTxFee,
	})
}

func queryTxVolume(_ sdk.Context, _ []string, k Keeper, cdc *codec.LegacyAmino, req abci.RequestQuery) ([]byte, error) {
	var request telemetrytypes.QueryTxVolumeRequest
	if err := cdc.UnmarshalJSON(req.Data, &request); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	txVolume, err := k.GetTxVolumePerDay(request.GetStartDate(), request.GetEndDate())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get tx volume")
	}

	return commontypes.QueryOK(cdc, telemetrytypes.QueryTxVolumeResponse{
		TxVolumePerDay: txVolume,
	})
}

func queryValidatorBlocks(
	ctx sdk.Context, _ []string, k Keeper, cdc *codec.LegacyAmino, req abci.RequestQuery,
) ([]byte, error) {
	var request telemetrytypes.QueryValidatorBlocksRequest
	if err := cdc.UnmarshalJSON(req.Data, &request); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	address, err := sdk.ValAddressFromBech32(request.ValidatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	blocks, total, err := k.GetValidatorBlocks(
		ctx,
		address,
		request.GetDesc(),
		request.GetPagination(),
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get validators blocks")
	}

	return commontypes.QueryOK(cdc, telemetrytypes.QueryValidatorBlocksResponse{
		Blocks: blocks,
		Pagination: &query.PageResponse{
			Total: total,
		},
	})
}

func queryTopValidators(
	ctx sdk.Context, _ []string, k Keeper, cdc *codec.LegacyAmino, req abci.RequestQuery,
) ([]byte, error) {
	var request telemetrytypes.QueryTopValidatorsRequest
	if err := cdc.UnmarshalJSON(req.Data, &request); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	validatorsBlocks, total, err := k.GetTopValidatorsByBlocks(
		ctx,
		request.GetStartDate(),
		request.GetEndDate(),
		request.GetDesc(),
		request.GetPagination(),
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get top validators by blocks")
	}

	return commontypes.QueryOK(cdc, telemetrytypes.QueryTopValidatorsResponse{
		TopValidators: validatorsBlocks,
		Pagination: &query.PageResponse{
			Total: total,
		},
	})
}

func queryValidatorByConsAddr(
	ctx sdk.Context,
	path []string,
	k Keeper,
	cdc *codec.LegacyAmino,
) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "consensus address not specified")
	}

	request := telemetrytypes.QueryValidatorByConsAddrRequest{
		ConsensusAddress: path[0],
	}

	validator, err := k.ValidatorByConsAddr(sdk.WrapSDKContext(ctx), &request)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to query validator by cons address")
	}

	return commontypes.QueryOK(cdc, validator.Validator)
}
