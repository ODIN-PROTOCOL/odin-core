package keeper

import (
	"fmt"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis new mint genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data *minttypes.GenesisState) {
	keeper.SetMinter(ctx, data.Minter)
	keeper.SetParams(ctx, data.Params)

	moduleAcc := keeper.GetMintAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", minttypes.ModuleName))
	}

	balances := keeper.bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsZero() {
		addr, err := sdk.AccAddressFromBech32(data.ModuleCoinsAccount)
		if err != nil {
			panic(err)
		}

		if err := keeper.bankKeeper.SendCoins(ctx, addr, moduleAcc.GetAddress(), data.MintPool.TreasuryPool); err != nil {
			panic(err)
		}

		keeper.SetMintModuleCoinsAccount(ctx, addr)
		keeper.authKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	keeper.SetMintPool(ctx, data.MintPool)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) *minttypes.GenesisState {
	minter := keeper.GetMinter(ctx)
	params := keeper.GetParams(ctx)
	mintPool := keeper.GetMintPool(ctx)
	mintModuleCoinsAccount := keeper.GetMintModuleCoinsAccount(ctx)
	return minttypes.NewGenesisState(minter, params, mintPool, mintModuleCoinsAccount)
}
