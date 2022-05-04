package oraclekeeper

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ODIN-PROTOCOL/odin-core/pkg/bandrng"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// 1 cosmos gas is equal to 7 owasm gas
const gasConversionFactor = 7

func convertToOwasmGas(cosmos uint64) uint32 {
	return uint32(cosmos * gasConversionFactor)
}

// GetRandomValidators returns a pseudorandom subset of active validators. Each validator has
// chance of getting selected directly proportional to the amount of voting power it has.
func (k Keeper) GetRandomValidators(ctx sdk.Context, size int, id int64) ([]sdk.ValAddress, error) {
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
	rng, err := bandrng.NewRng(k.GetRollingSeed(ctx), sdk.Uint64ToBigEndian(uint64(id)), []byte(ctx.ChainID()))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrBadDrbgInitialization, err.Error())
	}
	tryCount := int(k.GetParamUint64(ctx, types.KeySamplingTryCount))
	chosenValIndexes := bandrng.ChooseSomeMaxWeight(rng, valPowers, size, tryCount)
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
	ibcSource *types.IBCSource,
) (types.RequestID, error) {
	askCount := r.GetAskCount()
	if askCount > k.GetParamUint64(ctx, types.KeyMaxAskCount) {
		return 0, sdkerrors.Wrapf(types.ErrInvalidAskCount, "got: %d, max: %d", askCount, k.GetParamUint64(ctx, types.KeyMaxAskCount))
	}

	// Consume gas for data requests.
	ctx.GasMeter().ConsumeGas(askCount*k.GetParamUint64(ctx, types.KeyPerValidatorRequestGas), "PER_VALIDATOR_REQUEST_FEE")

	// Get a random validator set to perform this request.
	validators, err := k.GetRandomValidators(ctx, int(askCount), k.GetRequestCount(ctx)+1)
	if err != nil {
		return 0, err
	}

	// Create a request object. Note that RawRequestIDs will be populated after preparation is done.
	req := types.NewRequest(
		r.GetOracleScriptID(), r.GetCalldata(), validators, r.GetMinCount(),
		ctx.BlockHeight(), ctx.BlockTime(), r.GetClientID(), nil, ibcSource, r.GetExecuteGas(),
	)

	// Create an execution environment and call Owasm prepare function.
	env := types.NewPrepareEnv(req, int64(k.GetParamUint64(ctx, types.KeyMaxDataSize)), int64(k.GetParamUint64(ctx, types.KeyMaxRawRequestCount)))
	script, err := k.GetOracleScript(ctx, req.OracleScriptID)
	if err != nil {
		return 0, err
	}

	// Consume fee and execute owasm code
	ctx.GasMeter().ConsumeGas(k.GetParamUint64(ctx, types.KeyBaseOwasmGas), "BASE_OWASM_FEE")
	ctx.GasMeter().ConsumeGas(r.GetPrepareGas(), "OWASM_PREPARE_FEE")
	code := k.GetFile(script.Filename)
	maxDataSize := k.GetParamUint64(ctx, types.KeyMaxDataSize)
	output, err := k.owasmVM.Prepare(code, convertToOwasmGas(r.GetPrepareGas()), int64(maxDataSize), env)
	if err != nil {
		return 0, sdkerrors.Wrapf(types.ErrBadWasmExecution, err.Error())
	}

	// Preparation complete! It's time to collect raw request ids.
	req.RawRequests = env.GetRawRequests()
	// TODO compare Oracle fee implementation
	//fee := k.GetDataRequesterBasicFeeParam(ctx)

	//err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, feePayer, types.ModuleName, sdk.NewCoins(fee.Value()))
	//if err != nil {
	//	return 0, sdkerrors.Wrap(err, "sending coins from account to module")
	//}
	//
	//oraclePool := k.GetOraclePool(ctx)
	//oraclePool.DataProvidersPool = oraclePool.DataProvidersPool.Add(sdk.NewDecCoinFromCoin(fee.Value()))
	//k.SetOraclePool(ctx, oraclePool)

	// Preparation complete! Nothing can go wrong now (naive). It's time to collect raw request ids.
	if len(req.RawRequests) == 0 {
		return 0, types.ErrEmptyRawRequests
	}
	// TODO now fees are sent to the data source 'Treasury' which is, what we want is to send it to Data Providers Pool
	// TODO rework this and remove
	// Collect ds fee
	if _, err := k.CollectFee(ctx, feePayer, r.GetFeeLimit(), askCount, req.RawRequests); err != nil {
		return 0, err
	}
	// We now have everything we need to the request, so let's add it to the store.
	rid := k.AddRequest(ctx, req)

	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeRequest)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", rid)),
		sdk.NewAttribute(types.AttributeKeyClientID, req.ClientID),
		sdk.NewAttribute(types.AttributeKeyOracleScriptID, fmt.Sprintf("%d", req.OracleScriptID)),
		sdk.NewAttribute(types.AttributeKeyCalldata, hex.EncodeToString(req.Calldata)),
		sdk.NewAttribute(types.AttributeKeyAskCount, fmt.Sprintf("%d", askCount)),
		sdk.NewAttribute(types.AttributeKeyMinCount, fmt.Sprintf("%d", req.MinCount)),
		sdk.NewAttribute(types.AttributeKeyGasUsed, fmt.Sprintf("%d", output.GasUsed)),
	)
	for _, val := range req.RequestedValidators {
		event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyValidator, val))
	}
	ctx.EventManager().EmitEvent(event)

	// Subtract execute fee
	ctx.GasMeter().ConsumeGas(k.GetParamUint64(ctx, types.KeyBaseOwasmGas), "BASE_OWASM_FEE")
	ctx.GasMeter().ConsumeGas(r.GetExecuteGas(), "OWASM_EXECUTE_FEE")

	// Emit an event for each of the raw data requests
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
		))
	}
	return rid, nil
}

// ResolveRequest resolves the given request and saves the result to the store. The function
// assumes that the given request is in a resolvable state with sufficient reporters.
func (k Keeper) ResolveRequest(ctx sdk.Context, reqID types.RequestID) {
	req := k.MustGetRequest(ctx, reqID)
	env := types.NewExecuteEnv(req, k.GetRequestReports(ctx, reqID), ctx.BlockTime())
	script := k.MustGetOracleScript(ctx, req.OracleScriptID)
	code := k.GetFile(script.Filename)
	maxDataSize := k.GetParamUint64(ctx, types.KeyMaxDataSize)
	output, err := k.owasmVM.Execute(code, convertToOwasmGas(req.GetExecuteGas()), int64(maxDataSize), env)
	if err != nil {
		k.ResolveFailure(ctx, reqID, err.Error())
		// TODO: send response to IBC module on fail request
	} else if env.Retdata == nil {
		k.ResolveFailure(ctx, reqID, "no return data")
		// TODO: send response to IBC module on fail request
	} else {
		k.ResolveSuccess(ctx, reqID, env.Retdata, output.GasUsed)
	}
}
