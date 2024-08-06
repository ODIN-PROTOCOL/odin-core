package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	addresscodec "cosmossdk.io/core/address"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	owasm "github.com/ODIN-PROTOCOL/wasmvm/v2"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	"github.com/ODIN-PROTOCOL/odin-core/pkg/filecache"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

type Keeper struct {
	cdc              codec.BinaryCodec
	fileCache        filecache.Cache
	feeCollectorName string
	owasmVM          *owasm.Vm

	AuthKeeper            types.AccountKeeper
	BankKeeper            types.BankKeeper
	stakingKeeper         types.StakingKeeper
	distrKeeper           types.DistrKeeper
	authzKeeper           types.AuthzKeeper
	channelKeeper         types.ChannelKeeper
	portKeeper            types.PortKeeper
	scopedKeeper          capabilitykeeper.ScopedKeeper
	validatorAddressCodec addresscodec.Codec
	addressCodec          addresscodec.Codec

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	// The (unexposed) keys used to access the stores from the Context.
	storeService corestoretypes.KVStoreService

	Schema                          collections.Schema
	Params                          collections.Item[types.Params]
	DataSources                     collections.Map[uint64, types.DataSource]
	OracleScripts                   collections.Map[uint64, types.OracleScript]
	Requests                        collections.Map[uint64, types.Request]
	PendingResolveList              collections.Item[types.PendingResolveList]
	Reports                         collections.Map[collections.Pair[uint64, []byte], types.Report]
	Results                         collections.Map[uint64, types.Result]
	ValidatorStatuses               collections.Map[[]byte, types.ValidatorStatus]
	RequestID                       collections.Sequence
	DataSourceID                    collections.Sequence
	OracleScriptID                  collections.Sequence
	RollingSeed                     collections.Item[[]byte]
	RequestLastExpired              collections.Item[uint64]
	DataProviderAccumulatedRewards  collections.Map[[]byte, types.DataProviderAccumulatedReward]
	AccumulatedDataProvidersRewards collections.Item[types.DataProvidersAccumulatedRewards]
	AccumulatedPaymentsForData      collections.Item[types.AccumulatedPaymentsForData]
}

// NewKeeper creates a new oracle Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService corestoretypes.KVStoreService,
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
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:                   cdc,
		fileCache:             filecache.New(fileDir),
		feeCollectorName:      feeCollectorName,
		owasmVM:               owasmVM,
		AuthKeeper:            authKeeper,
		BankKeeper:            bankKeeper,
		stakingKeeper:         stakingKeeper,
		distrKeeper:           distrKeeper,
		authzKeeper:           authzKeeper,
		channelKeeper:         channelKeeper,
		portKeeper:            portKeeper,
		scopedKeeper:          scopeKeeper,
		validatorAddressCodec: stakingKeeper.ValidatorAddressCodec(),
		authority:             authority,
		storeService:          storeService,

		Params:                          collections.NewItem(sb, types.ParamsKeyPrefix, "params", codec.CollValue[types.Params](cdc)),
		DataSources:                     collections.NewMap(sb, types.DataSourceStoreKeyPrefix, "data_sources", collections.Uint64Key, codec.CollValue[types.DataSource](cdc)),
		OracleScripts:                   collections.NewMap(sb, types.OracleScriptStoreKeyPrefix, "oracle_scripts", collections.Uint64Key, codec.CollValue[types.OracleScript](cdc)),
		Requests:                        collections.NewMap(sb, types.RequestStoreKeyPrefix, "requests", collections.Uint64Key, codec.CollValue[types.Request](cdc)),
		PendingResolveList:              collections.NewItem(sb, types.PendingResolveListStoreKey, "pending_resolve_list", codec.CollValue[types.PendingResolveList](cdc)),
		Reports:                         collections.NewMap(sb, types.ReportStoreKeyPrefix, "reports", collections.PairKeyCodec(collections.Uint64Key, collections.BytesKey), codec.CollValue[types.Report](cdc)), //collections.Map[collections.Pair[uint64, sdk.ValAddress], types.Report]{},
		Results:                         collections.NewMap(sb, types.ResultStoreKeyPrefix, "results", collections.Uint64Key, codec.CollValue[types.Result](cdc)),
		ValidatorStatuses:               collections.NewMap(sb, types.ValidatorStatusKeyPrefix, "validator_statuses", collections.BytesKey, codec.CollValue[types.ValidatorStatus](cdc)),
		RequestID:                       collections.NewSequence(sb, types.RequestCountStoreKey, "request_id"),
		DataSourceID:                    collections.NewSequence(sb, types.DataSourceCountStoreKey, "data_source_id"),
		OracleScriptID:                  collections.NewSequence(sb, types.OracleScriptCountStoreKey, "oracle_script_id"),
		RollingSeed:                     collections.NewItem(sb, types.RollingSeedStoreKey, "rolling_seed", collections.BytesValue),
		RequestLastExpired:              collections.NewItem(sb, types.RequestLastExpiredStoreKey, "request_last_expired", collections.Uint64Value),
		DataProviderAccumulatedRewards:  collections.NewMap(sb, types.DataProviderRewardsKeyPrefix, "data_provider_accumulated_rewards", collections.BytesKey, codec.CollValue[types.DataProviderAccumulatedReward](cdc)),
		AccumulatedDataProvidersRewards: collections.NewItem(sb, types.AccumulatedDataProvidersRewardsStoreKey, "accumulated_data_providers_rewards", codec.CollValue[types.DataProvidersAccumulatedRewards](cdc)),
		AccumulatedPaymentsForData:      collections.NewItem(sb, types.AccumulatedPaymentsForDataStoreKey, "accumulated_payments_for_data", codec.CollValue[types.AccumulatedPaymentsForData](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
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
func (k Keeper) SetRollingSeed(ctx context.Context, rollingSeed []byte) error {
	return k.RollingSeed.Set(ctx, rollingSeed)
}

// GetRollingSeed returns the current rolling seed value.
func (k Keeper) GetRollingSeed(ctx context.Context) ([]byte, error) {
	return k.RollingSeed.Get(ctx)
}

// SetRequestCount sets the number of request count to the given value. Useful for genesis state.
func (k Keeper) SetRequestCount(ctx context.Context, count uint64) error {
	return k.RequestID.Set(ctx, count)
}

// GetRequestCount returns the current number of all requests ever exist.
func (k Keeper) GetRequestCount(ctx context.Context) (uint64, error) {
	return k.RequestID.Peek(ctx)
}

// SetRequestLastExpired sets the ID of the last expired request.
func (k Keeper) SetRequestLastExpired(ctx context.Context, id types.RequestID) error {
	return k.RequestLastExpired.Set(ctx, uint64(id))
}

// GetRequestLastExpired returns the ID of the last expired request.
func (k Keeper) GetRequestLastExpired(ctx context.Context) (types.RequestID, error) {
	lastExpired, err := k.RequestLastExpired.Get(ctx)
	return types.RequestID(lastExpired), err
}

// GetNextRequestID increments and returns the current number of requests.
func (k Keeper) GetNextRequestID(ctx context.Context) (types.RequestID, error) {
	count, err := k.RequestID.Next(ctx)
	return types.RequestID(count + 1), err
}

// SetDataSourceCount sets the number of data source count to the given value.
func (k Keeper) SetDataSourceCount(ctx context.Context, count uint64) error {
	return k.DataSourceID.Set(ctx, count)
}

// GetDataSourceCount returns the current number of all data sources ever exist.
func (k Keeper) GetDataSourceCount(ctx context.Context) (uint64, error) {
	return k.DataSourceID.Peek(ctx)
}

// GetNextDataSourceID increments and returns the current number of data sources.
func (k Keeper) GetNextDataSourceID(ctx context.Context) (types.DataSourceID, error) {
	count, err := k.DataSourceID.Next(ctx)
	return types.DataSourceID(count + 1), err
}

// SetOracleScriptCount sets the number of oracle script count to the given value.
func (k Keeper) SetOracleScriptCount(ctx context.Context, count uint64) error {
	return k.OracleScriptID.Set(ctx, count)
}

// GetOracleScriptCount returns the current number of all oracle scripts ever exist.
func (k Keeper) GetOracleScriptCount(ctx context.Context) (uint64, error) {
	return k.OracleScriptID.Peek(ctx)
}

// GetNextOracleScriptID increments and returns the current number of oracle scripts.
func (k Keeper) GetNextOracleScriptID(ctx context.Context) (types.OracleScriptID, error) {
	count, err := k.OracleScriptID.Next(ctx)
	return types.OracleScriptID(count + 1), err
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
	capability := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, capability, host.PortPath(portID))
}

// GetPort returns the portID for the oracle module. Used in ExportGenesis
func (k Keeper) GetPort() string {
	return types.PortID
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
	capability, _ := k.authzKeeper.GetAuthorization(
		ctx,
		reporter,
		sdk.AccAddress(validator),
		sdk.MsgTypeURL(&types.MsgReportData{}),
	)
	return capability != nil
}

// GrantReporter grants the reporter to validator for testing
func (k Keeper) GrantReporter(ctx sdk.Context, validator sdk.ValAddress, reporter sdk.AccAddress) error {
	expiration := ctx.BlockTime().Add(10 * time.Minute)
	return k.authzKeeper.SaveGrant(ctx, reporter, sdk.AccAddress(validator),
		authz.NewGenericAuthorization(sdk.MsgTypeURL(&types.MsgReportData{})), &expiration,
	)
}

// RevokeReporter revokes grant from the reporter for testing
func (k Keeper) RevokeReporter(ctx context.Context, validator sdk.ValAddress, reporter sdk.AccAddress) error {
	return k.authzKeeper.DeleteGrant(ctx, reporter, sdk.AccAddress(validator), sdk.MsgTypeURL(&types.MsgReportData{}))
}

func (k Keeper) SetAccumulatedDataProvidersRewards(ctx context.Context, reward types.DataProvidersAccumulatedRewards) error {
	return k.AccumulatedDataProvidersRewards.Set(ctx, reward)
}

func (k Keeper) GetAccumulatedDataProvidersRewards(ctx context.Context) (reward types.DataProvidersAccumulatedRewards, err error) {
	return k.AccumulatedDataProvidersRewards.Get(ctx)
}

func (k Keeper) SetAccumulatedPaymentsForData(ctx context.Context, payments types.AccumulatedPaymentsForData) error {
	return k.AccumulatedPaymentsForData.Set(ctx, payments)
}

func (k Keeper) GetAccumulatedPaymentsForData(ctx context.Context) (payments types.AccumulatedPaymentsForData, err error) {
	return k.AccumulatedPaymentsForData.Get(ctx)
}
