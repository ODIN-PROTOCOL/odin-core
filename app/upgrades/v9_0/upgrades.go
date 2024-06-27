package v9_0

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/ODIN-PROTOCOL/odin-core/app/keepers"
	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"
	wasmlctypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
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

		goCtx.Logger().Info("Upgrade v9 complete")
		return vm, nil
	}
}

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v0.9.0",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{wasmtypes.StoreKey, wasmlctypes.StoreKey, icqtypes.StoreKey},
	},
}
