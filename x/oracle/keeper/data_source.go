package keeper

import (
	"bytes"
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// HasDataSource checks if the data source of this ID exists in the storage.
func (k Keeper) HasDataSource(ctx sdk.Context, id types.DataSourceID) (bool, error) {
	return k.DataSources.Has(ctx, uint64(id))
}

// GetDataSource returns the data source struct for the given ID or error if not exists.
func (k Keeper) GetDataSource(ctx context.Context, id types.DataSourceID) (types.DataSource, error) {
	return k.DataSources.Get(ctx, uint64(id))
}

// MustGetDataSource returns the data source struct for the given ID. Panic if not exists.
func (k Keeper) MustGetDataSource(ctx sdk.Context, id types.DataSourceID) types.DataSource {
	dataSource, err := k.GetDataSource(ctx, id)
	if err != nil {
		panic(err)
	}
	return dataSource
}

// SetDataSource saves the given data source to the storage without performing validation.
func (k Keeper) SetDataSource(ctx sdk.Context, id types.DataSourceID, dataSource types.DataSource) error {
	return k.DataSources.Set(ctx, uint64(id), dataSource)
}

// AddDataSource adds the given data source to the storage.
func (k Keeper) AddDataSource(ctx sdk.Context, dataSource types.DataSource) (types.DataSourceID, error) {
	id, err := k.GetNextDataSourceID(ctx)
	if err != nil {
		return 0, err
	}

	err = k.SetDataSource(ctx, id, dataSource)
	return id, err
}

// MustEditDataSource edits the given data source by id and flushes it to the storage.
func (k Keeper) MustEditDataSource(ctx sdk.Context, id types.DataSourceID, new types.DataSource) {
	dataSource := k.MustGetDataSource(ctx, id)
	dataSource.Owner = new.Owner
	dataSource.Name = modify(dataSource.Name, new.Name)
	dataSource.Description = modify(dataSource.Description, new.Description)
	dataSource.Filename = modify(dataSource.Filename, new.Filename)
	dataSource.Treasury = new.Treasury
	dataSource.Fee = new.Fee

	err := k.SetDataSource(ctx, id, dataSource)
	if err != nil {
		panic(err)
	}
}

// GetAllDataSources returns the list of all data sources in the store, or nil if there is none.
func (k Keeper) GetAllDataSources(ctx sdk.Context) (dataSources []types.DataSource, err error) {
	err = k.DataSources.Walk(ctx, nil, func(_ uint64, dataSource types.DataSource) (stop bool, err error) {
		dataSources = append(dataSources, dataSource)
		return false, err
	})
	return dataSources, err
}

// GetPaginatedDataSources returns the list of all data sources in the store with pagination
func (k Keeper) GetPaginatedDataSources(
	ctx sdk.Context,
	limit, offset uint64,
) ([]types.DataSource, *query.PageResponse, error) {
	pagination := &query.PageRequest{
		Limit:  limit,
		Offset: offset,
	}

	dataSources, pageRes, err := query.CollectionPaginate(ctx, k.DataSources, pagination, func(key uint64, dataSource types.DataSource) (types.DataSource, error) {
		return dataSource, nil
	})

	return dataSources, pageRes, err
}

// AddExecutableFile saves the given executable file to a file to filecahe storage and returns
// its sha256sum reference name. Returns do-not-modify symbol if the input is do-not-modify.
func (k Keeper) AddExecutableFile(file []byte) string {
	if bytes.Equal(file, types.DoNotModifyBytes) {
		return types.DoNotModify
	}
	return k.fileCache.AddFile(file)
}
