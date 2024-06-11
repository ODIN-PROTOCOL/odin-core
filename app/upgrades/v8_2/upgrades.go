package v8_2

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	circuittypes "cosmossdk.io/x/circuit/types"
	"cosmossdk.io/x/nft"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/ODIN-PROTOCOL/odin-core/app/keepers"
	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.AppManager,
	_ *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		goCtx := sdk.UnwrapSDKContext(ctx)
		goCtx.Logger().Info("Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		goCtx.Logger().Info("Upgrade v17 complete")
		return vm, nil
	}
}

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v0.8.2",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{nft.StoreKey, circuittypes.StoreKey},
	},
}
