package keeper

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ODIN-PROTOCOL/odin-core/pkg/odinrng"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// 1 cosmos gas is equal to 20000000 owasm gas
const gasConversionFactor = 20_000_000

func ConvertToOwasmGas(cosmos uint64) uint64 {
	return uint64(cosmos * gasConversionFactor)
}

// GetSpanSize return maximum value between MaxReportDataSize and MaxCallDataSize
func (k Keeper) GetSpanSize(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	if params.MaxReportDataSize > params.MaxCalldataSize {
		return params.MaxReportDataSize
	}
	return params.MaxCalldataSize
}

// GetRandomValidators returns a pseudorandom subset of active validators. Each validator has
// chance of getting selected directly proportional to the amount of voting power it has.
func (k Keeper) GetRandomValidators(ctx sdk.Context, size int, id uint64) ([]sdk.ValAddress, error) {
	valOperators := []sdk.ValAddress{}
	valPowers := []uint64{}
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx,
		func(idx int64, val stakingtypes.ValidatorI) (stop bool) {
			if k.GetValidatorStatus(ctx, val.GetOperator()).IsActive {
				valOperators = append(valOperators, val.GetOperator())
				valPowers = append(valPowers, val.GetTokens().Uint64())
			}
			return false
		})
	if len(valOperators) < size {
		return nil, sdkerrors.Wrapf(
			types.ErrInsufficientValidators, "%d < %d", len(valOperators), size)
	}
	rng, err := odinrng.NewRng(k.GetRollingSeed(ctx), sdk.Uint64ToBigEndian(id), []byte(ctx.ChainID()))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrBadDrbgInitialization, err.Error())
	}
	tryCount := int(k.GetParams(ctx).SamplingTryCount)
	chosenValIndexes := odinrng.ChooseSomeMaxWeight(rng, valPowers, size, tryCount)
	validators := make([]sdk.ValAddress, size)
	for i, idx := range chosenValIndexes {
		validators[i] = valOperators[idx]
	}
	return validators, nil
}

// PrepareRequest takes an request specification object, performs the prepare call, and saves
// the request object to store. Also emits events related to the request.
func (k Keeper) PrepareRequest(
	ctx sdk.Context,
	r types.RequestSpec,
	feePayer sdk.AccAddress,
	ibcChannel *types.IBCChannel,
) (types.RequestID, error) {
	calldataSize := len(r.GetCalldata())
	if calldataSize > int(k.GetSpanSize(ctx)) {
		return 0, types.WrapMaxError(types.ErrTooLargeCalldata, calldataSize, int(k.GetSpanSize(ctx)))
	}

	askCount := r.GetAskCount()
	params := k.GetParams(ctx)
	if askCount > params.MaxAskCount {
		return 0, types.WrapMaxError(types.ErrInvalidAskCount, int(askCount), int(params.MaxAskCount))
	}

	// Consume gas for data requests.
	ctx.GasMeter().ConsumeGas(askCount*params.PerValidatorRequestGas, "PER_VALIDATOR_REQUEST_FEE")

	// Get a random validator set to perform this request.
	validators, err := k.GetRandomValidators(ctx, int(askCount), k.GetRequestCount(ctx)+1)
	if err != nil {
		return 0, err
	}

	// Create a request object. Note that RawRequestIDs will be populated after preparation is done.
	req := types.NewRequest(
		r.GetOracleScriptID(), r.GetCalldata(), validators, r.GetMinCount(),
		ctx.BlockHeight(), ctx.BlockTime(), r.GetClientID(), nil, ibcChannel, r.GetExecuteGas(),
	)

	// Create an execution environment and call Owasm prepare function.
	env := types.NewPrepareEnv(
		req,
		int64(params.MaxCalldataSize),
		int64(params.MaxRawRequestCount),
		int64(k.GetSpanSize(ctx)),
	)
	script, err := k.GetOracleScript(ctx, req.OracleScriptID)
	if err != nil {
		return 0, err
	}

	// Consume fee and execute owasm code
	ctx.GasMeter().ConsumeGas(params.BaseOwasmGas, "BASE_OWASM_FEE")
	ctx.GasMeter().ConsumeGas(r.GetPrepareGas(), "OWASM_PREPARE_FEE")
	code := k.GetFile(script.Filename)
	output, err := k.owasmVM.Prepare(code, ConvertToOwasmGas(r.GetPrepareGas()), env)
	if err != nil {
		return 0, sdkerrors.Wrapf(types.ErrBadWasmExecution, err.Error())
	}

	// Preparation complete! It's time to collect raw request ids.
	req.RawRequests = env.GetRawRequests()
	if len(req.RawRequests) == 0 {
		return 0, types.ErrEmptyRawRequests
	}
	// Collect ds fee
	totalFees, err := k.CollectFee(ctx, feePayer, r.GetFeeLimit(), askCount, req.RawRequests)
	if err != nil {
		return 0, err
	}
	// We now have everything we need to the request, so let's add it to the store.
	id := k.AddRequest(ctx, req)

	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeRequest)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyClientID, req.ClientID),
		sdk.NewAttribute(types.AttributeKeyOracleScriptID, fmt.Sprintf("%d", req.OracleScriptID)),
		sdk.NewAttribute(types.AttributeKeyCalldata, hex.EncodeToString(req.Calldata)),
		sdk.NewAttribute(types.AttributeKeyAskCount, fmt.Sprintf("%d", askCount)),
		sdk.NewAttribute(types.AttributeKeyMinCount, fmt.Sprintf("%d", req.MinCount)),
		sdk.NewAttribute(types.AttributeKeyGasUsed, fmt.Sprintf("%d", output.GasUsed)),
		sdk.NewAttribute(types.AttributeKeyTotalFees, totalFees.String()),
	)
	for _, val := range req.RequestedValidators {
		event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyValidator, val))
	}
	ctx.EventManager().EmitEvent(event)

	// Subtract execute fee
	ctx.GasMeter().ConsumeGas(params.BaseOwasmGas, "BASE_OWASM_FEE")
	ctx.GasMeter().ConsumeGas(r.GetExecuteGas(), "OWASM_EXECUTE_FEE")

	// Emit an event for each of the raw data requests.
	for _, rawReq := range env.GetRawRequests() {
		ds, err := k.GetDataSource(ctx, rawReq.DataSourceID)
		if err != nil {
			return 0, err
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, fmt.Sprintf("%d", rawReq.DataSourceID)),
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, ds.Filename),
			sdk.NewAttribute(types.AttributeKeyExternalID, fmt.Sprintf("%d", rawReq.ExternalID)),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(rawReq.Calldata)),
			sdk.NewAttribute(types.AttributeKeyFee, ds.Fee.String()),
		))
	}
	return id, nil
}

// ResolveRequest resolves the given request and saves the result to the store. The function
// assumes that the given request is in a resolvable state with sufficient reporters.
func (k Keeper) ResolveRequest(ctx sdk.Context, reqID types.RequestID) {
	req := k.MustGetRequest(ctx, reqID)
	env := types.NewExecuteEnv(req, k.GetReports(ctx, reqID), ctx.BlockTime(), int64(k.GetSpanSize(ctx)))
	script := k.MustGetOracleScript(ctx, req.OracleScriptID)
	code := k.GetFile(script.Filename)
	output, err := k.owasmVM.Execute(code, ConvertToOwasmGas(req.GetExecuteGas()), env)

	if err != nil {
		k.ResolveFailure(ctx, reqID, err.Error())
	} else if env.Retdata == nil {
		k.ResolveFailure(ctx, reqID, "no return data")
	} else {
		k.ResolveSuccess(ctx, reqID, env.Retdata, output.GasUsed)
	}
}

// CollectFee subtract fee from fee payer and send them to treasury
func (k Keeper) CollectFee(
	ctx sdk.Context,
	payer sdk.AccAddress,
	feeLimit sdk.Coins,
	askCount uint64,
	rawRequests []types.RawRequest,
) (sdk.Coins, error) {
	collector := newFeeCollector(k.distrKeeper, k, feeLimit, payer)

	for _, r := range rawRequests {
		ds, err := k.GetDataSource(ctx, r.DataSourceID)
		if err != nil {
			return nil, err
		}

		fee := sdk.NewCoins()
		for _, c := range ds.Fee {
			c.Amount = c.Amount.Mul(sdk.NewInt(int64(askCount)))
			fee = fee.Add(c)
		}

		if err := collector.Collect(ctx, fee); err != nil {
			return nil, err
		}
	}

	return collector.Collected(), nil
}
