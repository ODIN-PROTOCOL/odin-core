package v7_12

import (
	"cosmossdk.io/errors"
	"github.com/ODIN-PROTOCOL/odin-core/app/keepers"
	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

var AddressWithTokensToBurn = sdk.MustAccAddressFromBech32("odin1y6lz8fy3krg377kht8yugjg5uunn84nf4ux8d6")

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	am upgrades.AppManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		bankKeeper := keepers.BankKeeper
		coinsToBurn := bankKeeper.GetAllBalances(ctx, AddressWithTokensToBurn)
		err := bankKeeper.SendCoinsFromAccountToModule(ctx, AddressWithTokensToBurn, banktypes.ModuleName, coinsToBurn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to send tokens to bank module")
		}

		err = bankKeeper.Keeper.BurnCoins(ctx, banktypes.ModuleName, coinsToBurn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to burn tokens")
		}

		return vm, nil
	}
}

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v0.7.12",
	CreateUpgradeHandler: CreateUpgradeHandler,
}
