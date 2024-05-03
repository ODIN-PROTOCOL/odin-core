package keeper

import (
	"bytes"
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// HasOracleScript checks if the oracle script of this ID exists in the storage.
func (k Keeper) HasOracleScript(ctx context.Context, id types.OracleScriptID) (bool, error) {
	return k.OracleScripts.Has(ctx, uint64(id))
}

// GetOracleScript returns the oracle script struct for the given ID or error if not exists.
func (k Keeper) GetOracleScript(ctx context.Context, id types.OracleScriptID) (types.OracleScript, error) {
	oracleScript, err := k.OracleScripts.Get(ctx, uint64(id))
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.OracleScript{}, sdkerrors.Wrapf(types.ErrOracleScriptNotFound, "id: %d", id)
		}

		return oracleScript, err
	}

	return oracleScript, nil
}

// MustGetOracleScript returns the oracle script struct for the given ID. Panic if not exists.
func (k Keeper) MustGetOracleScript(ctx context.Context, id types.OracleScriptID) types.OracleScript {
	oracleScript, err := k.GetOracleScript(ctx, id)
	if err != nil {
		panic(err)
	}
	return oracleScript
}

// SetOracleScript saves the given oracle script to the storage without performing validation.
func (k Keeper) SetOracleScript(ctx context.Context, id types.OracleScriptID, oracleScript types.OracleScript) error {
	return k.OracleScripts.Set(ctx, uint64(id), oracleScript)
}

// AddOracleScript adds the given oracle script to the storage.
func (k Keeper) AddOracleScript(ctx context.Context, oracleScript types.OracleScript) (types.OracleScriptID, error) {
	id, err := k.GetNextOracleScriptID(ctx)
	if err != nil {
		return 0, err
	}
	err = k.SetOracleScript(ctx, id, oracleScript)
	return id, err
}

// MustEditOracleScript edits the given oracle script by id and flushes it to the storage. Panic if not exists.
func (k Keeper) MustEditOracleScript(ctx context.Context, id types.OracleScriptID, new types.OracleScript) {
	oracleScript := k.MustGetOracleScript(ctx, id)
	oracleScript.Owner = new.Owner
	oracleScript.Name = modify(oracleScript.Name, new.Name)
	oracleScript.Description = modify(oracleScript.Description, new.Description)
	oracleScript.Filename = modify(oracleScript.Filename, new.Filename)
	oracleScript.Schema = modify(oracleScript.Schema, new.Schema)
	oracleScript.SourceCodeURL = modify(oracleScript.SourceCodeURL, new.SourceCodeURL)

	err := k.SetOracleScript(ctx, id, oracleScript)
	if err != nil {
		panic(err)
	}
}

// GetAllOracleScripts returns the list of all oracle scripts in the store, or nil if there is none.
func (k Keeper) GetAllOracleScripts(ctx context.Context) (oracleScripts []types.OracleScript, err error) {
	err = k.OracleScripts.Walk(ctx, nil, func(_ uint64, oracleScript types.OracleScript) (stop bool, err error) {
		oracleScripts = append(oracleScripts, oracleScript)
		return false, err
	})
	return oracleScripts, err
}

// GetPaginatedOracleScripts returns oracle scripts with pagination.
func (k Keeper) GetPaginatedOracleScripts(
	ctx context.Context,
	limit, offset uint64,
) ([]types.OracleScript, *query.PageResponse, error) {
	pagination := &query.PageRequest{
		Limit:  limit,
		Offset: offset,
	}

	oracleScripts, pageRes, err := query.CollectionPaginate(ctx, k.OracleScripts, pagination, func(key uint64, oracleScript types.OracleScript) (types.OracleScript, error) {
		return oracleScript, nil
	})
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "failed to paginate oracle scripts")
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
