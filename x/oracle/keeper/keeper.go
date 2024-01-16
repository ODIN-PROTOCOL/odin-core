package keeper

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	owasm "github.com/odin-protocol/go-owasm/api"

	"github.com/ODIN-PROTOCOL/odin-core/pkg/filecache"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

const (
	RollingSeedSizeInBytes = 32
)

type Keeper struct {
	storeKey         storetypes.StoreKey
	cdc              codec.BinaryCodec
	fileCache        filecache.Cache
	feeCollectorName string
	owasmVM          *owasm.Vm

	AuthKeeper    types.AccountKeeper
	BankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	distrKeeper   types.DistrKeeper
	authzKeeper   types.AuthzKeeper
	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper creates a new oracle Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	fileDir string,
	feeCollectorName string,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distrKeeper types.DistrKeeper,
	authzKeeper types.AuthzKeeper,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	scopeKeeper capabilitykeeper.ScopedKeeper,
	owasmVM *owasm.Vm,
	authority string,
) Keeper {
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		fileCache:        filecache.New(fileDir),
		feeCollectorName: feeCollectorName,
		owasmVM:          owasmVM,
		AuthKeeper:       authKeeper,
		BankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		distrKeeper:      distrKeeper,
		authzKeeper:      authzKeeper,
		channelKeeper:    channelKeeper,
		portKeeper:       portKeeper,
		scopedKeeper:     scopeKeeper,
		authority:        authority,
	}
}

// GetAuthority returns the x/oracle module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetRollingSeed sets the rolling seed value to be provided value.
func (k Keeper) SetRollingSeed(ctx sdk.Context, rollingSeed []byte) {
	ctx.KVStore(k.storeKey).Set(types.RollingSeedStoreKey, rollingSeed)
}

// GetRollingSeed returns the current rolling seed value.
func (k Keeper) GetRollingSeed(ctx sdk.Context) []byte {
	return ctx.KVStore(k.storeKey).Get(types.RollingSeedStoreKey)
}

// SetRequestCount sets the number of request count to the given value. Useful for genesis state.
func (k Keeper) SetRequestCount(ctx sdk.Context, count uint64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	ctx.KVStore(k.storeKey).Set(types.RequestCountStoreKey, bz)
}

// GetRequestCount returns the current number of all requests ever exist.
func (k Keeper) GetRequestCount(ctx sdk.Context) uint64 {
	bz := ctx.KVStore(k.storeKey).Get(types.RequestCountStoreKey)
	return binary.BigEndian.Uint64(bz)
}

// SetRequestLastExpired sets the ID of the last expired request.
func (k Keeper) SetRequestLastExpired(ctx sdk.Context, id types.RequestID) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(id))
	ctx.KVStore(k.storeKey).Set(types.RequestLastExpiredStoreKey, bz)
}

// GetRequestLastExpired returns the ID of the last expired request.
func (k Keeper) GetRequestLastExpired(ctx sdk.Context) types.RequestID {
	bz := ctx.KVStore(k.storeKey).Get(types.RequestLastExpiredStoreKey)
	return types.RequestID(binary.BigEndian.Uint64(bz))
}

// GetNextRequestID increments and returns the current number of requests.
func (k Keeper) GetNextRequestID(ctx sdk.Context) types.RequestID {
	requestNumber := k.GetRequestCount(ctx)
	k.SetRequestCount(ctx, requestNumber+1)
	return types.RequestID(requestNumber + 1)
}

// SetDataSourceCount sets the number of data source count to the given value.
func (k Keeper) SetDataSourceCount(ctx sdk.Context, count uint64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	ctx.KVStore(k.storeKey).Set(types.DataSourceCountStoreKey, bz)
}

// GetDataSourceCount returns the current number of all data sources ever exist.
func (k Keeper) GetDataSourceCount(ctx sdk.Context) uint64 {
	bz := ctx.KVStore(k.storeKey).Get(types.DataSourceCountStoreKey)
	return binary.BigEndian.Uint64(bz)
}

// GetNextDataSourceID increments and returns the current number of data sources.
func (k Keeper) GetNextDataSourceID(ctx sdk.Context) types.DataSourceID {
	dataSourceCount := k.GetDataSourceCount(ctx)
	k.SetDataSourceCount(ctx, dataSourceCount+1)
	return types.DataSourceID(dataSourceCount + 1)
}

// SetOracleScriptCount sets the number of oracle script count to the given value.
func (k Keeper) SetOracleScriptCount(ctx sdk.Context, count uint64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	ctx.KVStore(k.storeKey).Set(types.OracleScriptCountStoreKey, bz)
}

// GetOracleScriptCount returns the current number of all oracle scripts ever exist.
func (k Keeper) GetOracleScriptCount(ctx sdk.Context) uint64 {
	bz := ctx.KVStore(k.storeKey).Get(types.OracleScriptCountStoreKey)
	return binary.BigEndian.Uint64(bz)
}

// GetNextOracleScriptID increments and returns the current number of oracle scripts.
func (k Keeper) GetNextOracleScriptID(ctx sdk.Context) types.OracleScriptID {
	oracleScriptCount := k.GetOracleScriptCount(ctx)
	k.SetOracleScriptCount(ctx, oracleScriptCount+1)
	return types.OracleScriptID(oracleScriptCount + 1)
}

// GetFile loads the file from the file storage. Panics if the file does not exist.
func (k Keeper) GetFile(name string) []byte {
	return k.fileCache.MustGetFile(name)
}

// IsBound checks if the oracle module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the oracle module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the oracle module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the oracle module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// IsReporter checks if the validator granted to the reporter
func (k Keeper) IsReporter(ctx sdk.Context, validator sdk.ValAddress, reporter sdk.AccAddress) bool {
	cap, _ := k.authzKeeper.GetAuthorization(
		ctx,
		reporter,
		sdk.AccAddress(validator),
		sdk.MsgTypeURL(&types.MsgReportData{}),
	)
	return cap != nil
}

// GrantReporter grants the reporter to validator for testing
func (k Keeper) GrantReporter(ctx sdk.Context, validator sdk.ValAddress, reporter sdk.AccAddress) error {
	expiration := ctx.BlockTime().Add(10 * time.Minute)
	return k.authzKeeper.SaveGrant(ctx, reporter, sdk.AccAddress(validator),
		authz.NewGenericAuthorization(sdk.MsgTypeURL(&types.MsgReportData{})), &expiration,
	)
}

// RevokeReporter revokes grant from the reporter for testing
func (k Keeper) RevokeReporter(ctx sdk.Context, validator sdk.ValAddress, reporter sdk.AccAddress) error {
	return k.authzKeeper.DeleteGrant(ctx, reporter, sdk.AccAddress(validator), sdk.MsgTypeURL(&types.MsgReportData{}))
}

func (k Keeper) SetAccumulatedDataProvidersRewards(ctx sdk.Context, reward types.DataProvidersAccumulatedRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&reward)
	store.Set(types.AccumulatedDataProvidersRewardsStoreKey, b)
}

func (k Keeper) GetAccumulatedDataProvidersRewards(ctx sdk.Context) (reward types.DataProvidersAccumulatedRewards) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AccumulatedDataProvidersRewardsStoreKey)
	k.cdc.MustUnmarshal(bz, &reward)
	return
}

func (k Keeper) SetAccumulatedPaymentsForData(ctx sdk.Context, payments types.AccumulatedPaymentsForData) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&payments)
	store.Set(types.AccumulatedPaymentsForDataStoreKey, b)
}

func (k Keeper) GetAccumulatedPaymentsForData(ctx sdk.Context) (payments types.AccumulatedPaymentsForData) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AccumulatedPaymentsForDataStoreKey)
	k.cdc.MustUnmarshal(bz, &payments)
	return
}

// get the module coins account
func (k Keeper) GetOracleModuleCoinsAccount(ctx sdk.Context) (account sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.OracleModuleCoinsAccountKey)
	if b == nil {
		return nil
	}

	return sdk.AccAddress(b)
}

// set the module coins account
func (k Keeper) SetOracleModuleCoinsAccount(ctx sdk.Context, account sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(types.OracleModuleCoinsAccountKey, account)
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
