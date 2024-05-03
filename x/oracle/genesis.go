package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// InitGenesis performs genesis initialization for the oracle module.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data *types.GenesisState) {
	err := k.SetParams(ctx, data.Params)
	if err != nil {
		panic(err)
	}

	err = k.SetDataSourceCount(ctx, 0)
	if err != nil {
		panic(err)
	}

	err = k.SetOracleScriptCount(ctx, 0)
	if err != nil {
		panic(err)
	}

	err = k.SetRequestCount(ctx, 0)
	if err != nil {
		panic(err)
	}

	err = k.SetRequestLastExpired(ctx, 0)
	if err != nil {
		panic(err)
	}

	err = k.SetPendingResolveList(ctx, []types.RequestID{})
	if err != nil {
		panic(err)
	}

	err = k.SetAccumulatedPaymentsForData(ctx, types.AccumulatedPaymentsForData{AccumulatedAmount: sdk.NewCoins()})
	if err != nil {
		panic(err)
	}

	err = k.SetRollingSeed(ctx, make([]byte, types.RollingSeedSizeInBytes))
	if err != nil {
		panic(err)
	}

	for _, dataSource := range data.DataSources {
		_, err = k.AddDataSource(ctx, dataSource)
		if err != nil {
			panic(err)
		}
	}
	for _, oracleScript := range data.OracleScripts {
		_, err = k.AddOracleScript(ctx, oracleScript)
		if err != nil {
			panic(err)
		}
	}

	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !k.IsBound(ctx, types.PortID) {
		// oracle module binds to the oracle port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, types.PortID)
		if err != nil {
			panic(fmt.Sprintf("could not claim port capability: %v", err))
		}
	}

	moduleAcc := k.AuthKeeper.GetModuleAccount(ctx, types.ModuleName)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	k.AuthKeeper.SetModuleAccount(ctx, moduleAcc)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (*types.GenesisState, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	dataSources, err := k.GetAllDataSources(ctx)
	if err != nil {
		return nil, err
	}

	oracleScipts, err := k.GetAllOracleScripts(ctx)
	if err != nil {
		return nil, err
	}

	return &types.GenesisState{
		Params:        params,
		DataSources:   dataSources,
		OracleScripts: oracleScipts,
	}, nil
}
