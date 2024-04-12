package v7_12

import (
	"cosmossdk.io/errors"
	"github.com/ODIN-PROTOCOL/odin-core/app/keepers"
	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const AddressWithTokensToBurn = "odin1y6lz8fy3krg377kht8yugjg5uunn84nf4ux8d6"

func CreateUpgradeHandler(
	_ *module.Manager,
	_ module.Configurator,
	_ upgrades.AppManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		accAddressWithTokensToBurn := sdk.MustAccAddressFromBech32(AddressWithTokensToBurn)
		bankKeeper := keepers.BankKeeper

		dokiToBurn := bankKeeper.GetBalance(ctx, accAddressWithTokensToBurn, "udoki")
		myrkToBurn := bankKeeper.GetBalance(ctx, accAddressWithTokensToBurn, "umyrk")
		coinsToBurn := sdk.NewCoins(dokiToBurn, myrkToBurn)

		err := bankKeeper.SendCoinsFromAccountToModule(ctx, accAddressWithTokensToBurn, govtypes.ModuleName, coinsToBurn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to send tokens to bank module")
		}

		err = bankKeeper.Keeper.BurnCoins(ctx, govtypes.ModuleName, coinsToBurn)
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
