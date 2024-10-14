package v10_0

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/ODIN-PROTOCOL/odin-core/app/keepers"
	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades"
	onfttypes "github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
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

		goCtx.Logger().Info("Upgrade v10 complete")
		return vm, nil
	}
}

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v0.10.0",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{onfttypes.StoreKey},
	},
}
