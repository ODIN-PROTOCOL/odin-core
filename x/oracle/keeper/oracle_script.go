package keeper

import (
	"bytes"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// HasOracleScript checks if the oracle script of this ID exists in the storage.
func (k Keeper) HasOracleScript(ctx sdk.Context, id types.OracleScriptID) bool {
	return ctx.KVStore(k.storeKey).Has(types.OracleScriptStoreKey(id))
}

// GetOracleScript returns the oracle script struct for the given ID or error if not exists.
func (k Keeper) GetOracleScript(ctx sdk.Context, id types.OracleScriptID) (types.OracleScript, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.OracleScriptStoreKey(id))
	if bz == nil {
		return types.OracleScript{}, errors.Wrapf(types.ErrOracleScriptNotFound, "id: %d", id)
	}
	var oracleScript types.OracleScript
	k.cdc.MustUnmarshal(bz, &oracleScript)
	return oracleScript, nil
}

// MustGetOracleScript returns the oracle script struct for the given ID. Panic if not exists.
func (k Keeper) MustGetOracleScript(ctx sdk.Context, id types.OracleScriptID) types.OracleScript {
	oracleScript, err := k.GetOracleScript(ctx, id)
	if err != nil {
		panic(err)
	}
	return oracleScript
}

// SetOracleScript saves the given oracle script to the storage without performing validation.
func (k Keeper) SetOracleScript(ctx sdk.Context, id types.OracleScriptID, oracleScript types.OracleScript) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OracleScriptStoreKey(id), k.cdc.MustMarshal(&oracleScript))
}

// AddOracleScript adds the given oracle script to the storage.
func (k Keeper) AddOracleScript(ctx sdk.Context, oracleScript types.OracleScript) types.OracleScriptID {
	id := k.GetNextOracleScriptID(ctx)
	k.SetOracleScript(ctx, id, oracleScript)
	return id
}

// MustEditOracleScript edits the given oracle script by id and flushes it to the storage. Panic if not exists.
func (k Keeper) MustEditOracleScript(ctx sdk.Context, id types.OracleScriptID, new types.OracleScript) {
	oracleScript := k.MustGetOracleScript(ctx, id)
	oracleScript.Owner = new.Owner
	oracleScript.Name = modify(oracleScript.Name, new.Name)
	oracleScript.Description = modify(oracleScript.Description, new.Description)
	oracleScript.Filename = modify(oracleScript.Filename, new.Filename)
	oracleScript.Schema = modify(oracleScript.Schema, new.Schema)
	oracleScript.SourceCodeURL = modify(oracleScript.SourceCodeURL, new.SourceCodeURL)
	k.SetOracleScript(ctx, id, oracleScript)
}

// GetAllOracleScripts returns the list of all oracle scripts in the store, or nil if there is none.
func (k Keeper) GetAllOracleScripts(ctx sdk.Context) (oracleScripts []types.OracleScript) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleScriptStoreKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oracleScript types.OracleScript
		k.cdc.MustUnmarshal(iterator.Value(), &oracleScript)
		oracleScripts = append(oracleScripts, oracleScript)
	}
	return oracleScripts
}

// GetPaginatedOracleScripts returns oracle scripts with pagination.
func (k Keeper) GetPaginatedOracleScripts(
	ctx sdk.Context,
	limit, offset uint64,
) ([]types.OracleScript, *query.PageResponse, error) {
	oracleScripts := make([]types.OracleScript, 0)
	oracleScriptsStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.OracleScriptStoreKeyPrefix)
	pagination := &query.PageRequest{
		Limit:  limit,
		Offset: offset,
	}

	pageRes, err := query.FilteredPaginate(
		oracleScriptsStore,
		pagination,
		func(key []byte, value []byte, accumulate bool) (bool, error) {
			var oracleScript types.OracleScript
			if err := k.cdc.Unmarshal(value, &oracleScript); err != nil {
				return false, err
			}
			if accumulate {
				oracleScripts = append(oracleScripts, oracleScript)
			}
			return true, nil
		},
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to paginate oracle scripts")
	}

	return oracleScripts, pageRes, nil
}

// AddOracleScriptFile compiles Wasm code (see go-owasm), adds the compiled file to filecache,
// and returns its sha256 reference name. Returns do-not-modify symbol if input is do-not-modify.
func (k Keeper) AddOracleScriptFile(file []byte) (string, error) {
	if bytes.Equal(file, types.DoNotModifyBytes) {
		return types.DoNotModify, nil
	}
	compiledFile, err := k.owasmVM.Compile(file, types.MaxCompiledWasmCodeSize)
	if err != nil {
		return "", sdkerrors.Wrapf(types.ErrOwasmCompilation, "caused by %s", err.Error())
	}
	return k.fileCache.AddFile(compiledFile), nil
}
