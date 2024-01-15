package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// GetOraclePool gets the oracle pool info
func (k Keeper) GetOraclePool(ctx sdk.Context) (oraclePool oracletypes.OraclePool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(oracletypes.OraclePoolStoreKey)
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}
	k.cdc.MustUnmarshalLengthPrefixed(b, &oraclePool)
	return
}

// SetOraclePool sets the oracle pool info
func (k Keeper) SetOraclePool(ctx sdk.Context, oraclePool oracletypes.OraclePool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalLengthPrefixed(&oraclePool)
	store.Set(oracletypes.OraclePoolStoreKey, b)
}
