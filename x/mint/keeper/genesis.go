package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
)

// InitGenesis new mint genesis
func (k Keeper) InitGenesis(ctx sdk.Context, ak minttypes.AccountKeeper, data *minttypes.GenesisState) {
	k.SetMinter(ctx, data.Minter)

	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}

	moduleAcc := k.GetMintAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", minttypes.ModuleName))
	}

	balances := k.bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsZero() {
		addr, err := sdk.AccAddressFromBech32(data.ModuleCoinsAccount)
		if err != nil {
			panic(err)
		}

		if err := k.bankKeeper.SendCoins(ctx, addr, moduleAcc.GetAddress(), data.MintPool.TreasuryPool); err != nil {
			panic(err)
		}

		k.SetMintModuleCoinsAccount(ctx, addr)
		k.authKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	k.SetMintPool(ctx, data.MintPool)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *minttypes.GenesisState {
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)
	mintPool := k.GetMintPool(ctx)
	mintModuleCoinsAccount := k.GetMintModuleCoinsAccount(ctx)
	return minttypes.NewGenesisState(minter, params, mintPool, mintModuleCoinsAccount)
}
