package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/odinmint/types"
)

// InitGenesis new mint genesis
func InitGenesis(ctx sdk.Context, kp Keeper, data *minttypes.GenesisState) {
	kp.SetMinter(ctx, data.Minter)
	kp.SetParams(ctx, data.Params)

	moduleAcc := kp.GetMintAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", minttypes.ModuleName))
	}

	balances := kp.bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsZero() {
		addr, err := sdk.AccAddressFromBech32(data.ModuleCoinsAccount)
		if err != nil {
			panic(err)
		}

		if err := kp.bankKeeper.SendCoins(ctx, addr, moduleAcc.GetAddress(), data.MintPool.TreasuryPool); err != nil {
			panic(err)
		}

		kp.SetMintModuleCoinsAccount(ctx, addr)
	}

	kp.SetMintPool(ctx, data.MintPool)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, kp Keeper) *minttypes.GenesisState {
	minter := kp.GetMinter(ctx)
	params := kp.GetParams(ctx)
	mintPool := kp.GetMintPool(ctx)
	
	mintModuleCoinsAccount := kp.GetMintModuleCoinsAccount(ctx)
	return minttypes.NewGenesisState(minter, params, mintPool, mintModuleCoinsAccount)
}
