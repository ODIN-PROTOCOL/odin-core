package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
)

// InitGenesis new mint genesis
func (k Keeper) InitGenesis(ctx sdk.Context, ak minttypes.AccountKeeper, data *minttypes.GenesisState) {
	if err := k.SetMinter(ctx, data.Minter); err != nil {
		panic(err)
	}

	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}

	moduleAcc := k.authKeeper.GetModuleAccount(ctx, minttypes.ModuleName)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", minttypes.ModuleName))
	}

	balances := k.bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsAllLTE(data.MintPool.TreasuryPool) {
		diff := data.MintPool.TreasuryPool.Sub(balances...)

		if err := k.MintCoins(ctx, diff); err != nil {
			panic(err)
		}

		k.authKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	if err := k.SetMintPool(ctx, data.MintPool); err != nil {
		panic(err)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *minttypes.GenesisState {
	minter, err := k.GetMinter(ctx)
	if err != nil {
		panic(err)
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	mintPool, err := k.GetMintPool(ctx)
	if err != nil {
		panic(err)
	}

	return minttypes.NewGenesisState(minter, params, mintPool)
}
