package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	coinswaptypes "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
)

type Keeper struct {
	storeKey     storetypes.StoreKey
	cdc          codec.BinaryCodec
	paramstore   paramstypes.Subspace
	bankKeeper   coinswaptypes.BankKeeper
	distrKeeper  coinswaptypes.DistrKeeper
	oracleKeeper coinswaptypes.OracleKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	subspace paramstypes.Subspace,
	bk coinswaptypes.BankKeeper,
	dk coinswaptypes.DistrKeeper,
	ok coinswaptypes.OracleKeeper,
) Keeper {
	if !subspace.HasKeyTable() {
		subspace = subspace.WithKeyTable(coinswaptypes.ParamKeyTable())
	}
	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		paramstore:   subspace,
		bankKeeper:   bk,
		distrKeeper:  dk,
		oracleKeeper: ok,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", coinswaptypes.ModuleName))
}

// GetDecParam returns the parameter as specified by key as sdk.Dec
func (k Keeper) GetDecParam(ctx sdk.Context, key []byte) (res sdk.Dec) {
	k.paramstore.Get(ctx, key, &res)
	return res
}

// SetParams saves the given key-value parameter to the store.
func (k Keeper) SetParams(ctx sdk.Context, value coinswaptypes.Params) {
	k.paramstore.SetParamSet(ctx, &value)
}

// GetParams returns all current parameters as a types.Params instance.
func (k Keeper) GetParams(ctx sdk.Context) (params coinswaptypes.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetExchanges(ctx sdk.Context, value []coinswaptypes.Exchange) {
	k.paramstore.Set(ctx, coinswaptypes.KeyExchangeRates, value)
}

func (k Keeper) GetExchanges(ctx sdk.Context) (res []coinswaptypes.Exchange) {
	k.paramstore.Get(ctx, coinswaptypes.KeyExchangeRates, &res)
	return res
}

func (k Keeper) SetInitialRate(ctx sdk.Context, value sdk.Dec) {
	bz := k.cdc.MustMarshalLengthPrefixed(&gogotypes.StringValue{Value: value.String()})
	ctx.KVStore(k.storeKey).Set(coinswaptypes.InitialRateStoreKey, bz)
}

func (k Keeper) GetInitialRate(ctx sdk.Context) (rate sdk.Dec) {
	bz := ctx.KVStore(k.storeKey).Get(coinswaptypes.InitialRateStoreKey)
	var rawRate gogotypes.StringValue
	k.cdc.MustUnmarshalLengthPrefixed(bz, &rawRate)
	rate = sdk.MustNewDecFromStr(rawRate.Value)
	return rate
}
