package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// InitGenesis performs genesis initialization for the oracle module.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data *types.GenesisState) {
	k.SetParams(ctx, data.Params)
	k.SetDataSourceCount(ctx, 0)
	k.SetOracleScriptCount(ctx, 0)
	k.SetRequestCount(ctx, 0)
	k.SetRequestLastExpired(ctx, 0)
	k.SetRollingSeed(ctx, make([]byte, types.RollingSeedSizeInBytes))
	for _, dataSource := range data.DataSources {
		_ = k.AddDataSource(ctx, dataSource)
	}
	for _, oracleScript := range data.OracleScripts {
		_ = k.AddOracleScript(ctx, oracleScript)
	}

	k.SetPort(ctx, types.PortID)
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

	balances := k.BankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsZero() {
		if err := k.BankKeeper.SendCoins(ctx, sdk.AccAddress(data.ModuleCoinsAccount), moduleAcc.GetAddress(), data.OraclePool.DataProvidersPool); err != nil {
			panic(err)
		}

		k.AuthKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	k.SetOraclePool(ctx, data.OraclePool)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:             k.GetParams(ctx),
		DataSources:        k.GetAllDataSources(ctx),
		OracleScripts:      k.GetAllOracleScripts(ctx),
		OraclePool:         k.GetOraclePool(ctx),
		ModuleCoinsAccount: k.GetOracleModuleCoinsAccount(ctx).String(),
	}
}
