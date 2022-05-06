package oraclekeeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	host "github.com/cosmos/ibc-go/v2/modules/core/24-host"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ODIN-PROTOCOL/odin-core/pkg/filecache"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	owasm "github.com/bandprotocol/go-owasm/api"
)

type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              codec.BinaryCodec
	fileCache        filecache.Cache
	feeCollectorName string
	paramstore       paramtypes.Subspace
	owasmVM          *owasm.Vm

	authKeeper    oracletypes.AccountKeeper
	bankKeeper    oracletypes.BankKeeper
	distrKeeper   oracletypes.DistrKeeper
	stakingKeeper oracletypes.StakingKeeper
	channelKeeper oracletypes.ChannelKeeper
	portKeeper    oracletypes.PortKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new oracle Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	ps paramtypes.Subspace,
	fileDir string,
	feeCollectorName string,
	authKeeper oracletypes.AccountKeeper,
	bankKeeper oracletypes.BankKeeper,
	stakingKeeper oracletypes.StakingKeeper,
	distrKeeper oracletypes.DistrKeeper,
	channelKeeper oracletypes.ChannelKeeper,
	portKeeper oracletypes.PortKeeper,
	scopeKeeper capabilitykeeper.ScopedKeeper,
	owasmVM *owasm.Vm,
) Keeper {
	if addr := authKeeper.GetModuleAddress(oracletypes.ModuleName); addr == nil {
		panic("the oracle module account has not been set")
	}

	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(oracletypes.ParamKeyTable())
	}
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		fileCache:        filecache.New(fileDir),
		feeCollectorName: feeCollectorName,
		paramstore:       ps,
		owasmVM:          owasmVM,
		authKeeper:       authKeeper,
		bankKeeper:       bankKeeper,
		distrKeeper:      distrKeeper,
		stakingKeeper:    stakingKeeper,
		channelKeeper:    channelKeeper,
		portKeeper:       portKeeper,
		scopedKeeper:     scopeKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", oracletypes.ModuleName))
}

// GetParamUint64 returns the parameter as specified by key as an uint64.
func (k Keeper) GetParamUint64(ctx sdk.Context, key []byte) (res uint64) {
	k.paramstore.Get(ctx, key, &res)
	return res
}

func (k Keeper) SetParamUint64(ctx sdk.Context, key []byte, value uint64) {
	k.paramstore.Set(ctx, key, value)
}

// SetParams saves the given key-value parameter to the store.
func (k Keeper) SetParams(ctx sdk.Context, value oracletypes.Params) {
	k.paramstore.SetParamSet(ctx, &value)
}

// GetParams returns all current parameters as a oracletypes.Params instance.
func (k Keeper) GetParams(ctx sdk.Context) (params oracletypes.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetDataProviderRewardPerByteParam(ctx sdk.Context, value sdk.Coins) {
	k.paramstore.Set(ctx, oracletypes.KeyDataProviderRewardPerByte, value)
}

func (k Keeper) GetDataProviderRewardPerByteParam(ctx sdk.Context) (res sdk.Coins) {
	k.paramstore.Get(ctx, oracletypes.KeyDataProviderRewardPerByte, &res)
	return res
}

func (k Keeper) SetDataProviderRewardThresholdParam(ctx sdk.Context, value oracletypes.RewardThreshold) {
	k.paramstore.Set(ctx, oracletypes.KeyDataProviderRewardThreshold, value)
}

func (k Keeper) GetDataProviderRewardThresholdParam(ctx sdk.Context) (res oracletypes.RewardThreshold) {
	k.paramstore.Get(ctx, oracletypes.KeyDataProviderRewardThreshold, &res)
	return res
}

func (k Keeper) SetRewardDecreasingFractionParam(ctx sdk.Context, value sdk.Dec) {
	k.paramstore.Set(ctx, oracletypes.KeyRewardDecreasingFraction, value)
}

func (k Keeper) GetRewardDecreasingFractionParam(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, oracletypes.KeyRewardDecreasingFraction, &res)
	return res
}

func (k Keeper) SetDataRequesterFeeDenomsParam(ctx sdk.Context, value []string) {
	k.paramstore.Set(ctx, oracletypes.KeyDataRequesterFeeDenoms, value)
}

func (k Keeper) GetDataRequesterFeeDenomsParam(ctx sdk.Context) (res []string) {
	k.paramstore.Get(ctx, oracletypes.KeyDataRequesterFeeDenoms, &res)
	return res
}

// SetRollingSeed sets the rolling seed value to be provided value.
func (k Keeper) SetRollingSeed(ctx sdk.Context, rollingSeed []byte) {
	ctx.KVStore(k.storeKey).Set(oracletypes.RollingSeedStoreKey, rollingSeed)
}

// GetRollingSeed returns the current rolling seed value.
func (k Keeper) GetRollingSeed(ctx sdk.Context) []byte {
	return ctx.KVStore(k.storeKey).Get(oracletypes.RollingSeedStoreKey)
}

// SetRequestCount sets the number of request count to the given value. Useful for genesis state.
func (k Keeper) SetRequestCount(ctx sdk.Context, count int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(oracletypes.RequestCountStoreKey, k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: count}))
}

// GetRequestCount returns the current number of all requests ever exist.
func (k Keeper) GetRequestCount(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(oracletypes.RequestCountStoreKey)
	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

// SetRequestLastExpired sets the ID of the last expired request.
func (k Keeper) SetRequestLastExpired(ctx sdk.Context, id oracletypes.RequestID) {
	store := ctx.KVStore(k.storeKey)
	store.Set(oracletypes.RequestLastExpiredStoreKey, k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: int64(id)}))
}

// GetRequestLastExpired returns the ID of the last expired request.
func (k Keeper) GetRequestLastExpired(ctx sdk.Context) oracletypes.RequestID {
	bz := ctx.KVStore(k.storeKey).Get(oracletypes.RequestLastExpiredStoreKey)
	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return oracletypes.RequestID(intV.GetValue())
}

// GetNextRequestID increments and returns the current number of requests.
func (k Keeper) GetNextRequestID(ctx sdk.Context) oracletypes.RequestID {
	requestNumber := k.GetRequestCount(ctx)
	bz := k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: requestNumber + 1})
	ctx.KVStore(k.storeKey).Set(oracletypes.RequestCountStoreKey, bz)
	return oracletypes.RequestID(requestNumber + 1)
}

// SetDataSourceCount sets the number of data source count to the given value.
func (k Keeper) SetDataSourceCount(ctx sdk.Context, count int64) {
	bz := k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: count})
	ctx.KVStore(k.storeKey).Set(oracletypes.DataSourceCountStoreKey, bz)
}

// GetDataSourceCount returns the current number of all data sources ever exist.
func (k Keeper) GetDataSourceCount(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(oracletypes.DataSourceCountStoreKey)
	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

// GetNextDataSourceID increments and returns the current number of data sources.
func (k Keeper) GetNextDataSourceID(ctx sdk.Context) oracletypes.DataSourceID {
	dataSourceCount := k.GetDataSourceCount(ctx)
	k.SetDataSourceCount(ctx, dataSourceCount+1)
	return oracletypes.DataSourceID(dataSourceCount + 1)
}

// SetOracleScriptCount sets the number of oracle script count to the given value.
func (k Keeper) SetOracleScriptCount(ctx sdk.Context, count int64) {
	bz := k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: count})
	ctx.KVStore(k.storeKey).Set(oracletypes.OracleScriptCountStoreKey, bz)
}

// GetOracleScriptCount returns the current number of all oracle scripts ever exist.
func (k Keeper) GetOracleScriptCount(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(oracletypes.OracleScriptCountStoreKey)
	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

// GetNextOracleScriptID increments and returns the current number of oracle scripts.
func (k Keeper) GetNextOracleScriptID(ctx sdk.Context) oracletypes.OracleScriptID {
	oracleScriptCount := k.GetOracleScriptCount(ctx)
	k.SetOracleScriptCount(ctx, oracleScriptCount+1)
	return oracletypes.OracleScriptID(oracleScriptCount + 1)
}

// GetFile loads the file from the file storage. Panics if the file does not exist.
func (k Keeper) GetFile(name string) []byte {
	return k.fileCache.MustGetFile(name)
}

// IsBound checks if the transfer module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	capability := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, capability, host.PortPath(portID))
}

// GetPort returns the portID for the transfer module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(oracletypes.PortKey))
}

// SetPort sets the portID for the transfer module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(oracletypes.PortKey, []byte(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

func (k Keeper) SetAccumulatedDataProvidersRewards(ctx sdk.Context, reward oracletypes.DataProvidersAccumulatedRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&reward)
	store.Set(oracletypes.AccumulatedDataProvidersRewardsStoreKey, b)
}

func (k Keeper) GetAccumulatedDataProvidersRewards(ctx sdk.Context) (reward oracletypes.DataProvidersAccumulatedRewards) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(oracletypes.AccumulatedDataProvidersRewardsStoreKey)
	k.cdc.MustUnmarshal(bz, &reward)
	return
}

func (k Keeper) SetAccumulatedPaymentsForData(ctx sdk.Context, payments oracletypes.AccumulatedPaymentsForData) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&payments)
	store.Set(oracletypes.AccumulatedPaymentsForDataStoreKey, b)
}

func (k Keeper) GetAccumulatedPaymentsForData(ctx sdk.Context) (payments oracletypes.AccumulatedPaymentsForData) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(oracletypes.AccumulatedPaymentsForDataStoreKey)
	k.cdc.MustUnmarshal(bz, &payments)
	return
}

// get the module coins account
func (k Keeper) GetOracleModuleCoinsAccount(ctx sdk.Context) (account sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(oracletypes.OracleModuleCoinsAccountKey)
	if b == nil {
		return nil
	}

	return sdk.AccAddress(b)
}

// set the module coins account
func (k Keeper) SetOracleModuleCoinsAccount(ctx sdk.Context, account sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(oracletypes.OracleModuleCoinsAccountKey, account)
}

// IsAllowedFeeDenom checks if fee denom exists in data requester fee denoms list
func (k Keeper) IsAllowedFeeDenom(ctx sdk.Context, feeDenom string) bool {
	params := k.GetParams(ctx)

	for _, paramsFeeDenom := range params.DataRequesterFeeDenoms {
		if feeDenom == paramsFeeDenom {
			return true
		}
	}

	return false
}

// IsInDataProviderRewardPerByte checks if fee denom exists in data provider reward per byte list
func (k Keeper) IsInDataProviderRewardPerByte(ctx sdk.Context, feeDenom string) bool {
	params := k.GetParams(ctx)

	for _, paramsRewardPerByte := range params.DataProviderRewardPerByte {
		if feeDenom == paramsRewardPerByte.Denom {
			return true
		}
	}

	return false
}

// IsInDataProviderRewardThreshold checks if fee denom exists in data provider reward threshold list
func (k Keeper) IsInDataProviderRewardThreshold(ctx sdk.Context, feeDenom string) bool {
	params := k.GetParams(ctx)

	for _, paramsRewardThreshold := range params.DataProviderRewardThreshold.Amount {
		if feeDenom == paramsRewardThreshold.Denom {
			return true
		}
	}

	return false
}
