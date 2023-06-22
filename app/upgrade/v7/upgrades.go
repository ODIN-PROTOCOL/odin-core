package v7

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	mintkeeper "github.com/ODIN-PROTOCOL/odin-core/x/mint/keeper"
)

const newAddress = "odin17zhnwfs7rh78kz628l2mxjt0u6456rznjxyu6f"

func getBalance(ctx sdk.Context, bankkeeper bankkeeper.Keeper, addr sdk.AccAddress) sdk.Coins {
	return bankkeeper.GetAllBalances(ctx, addr)
}

func getAddresses() []sdk.AccAddress {
	var addresses [15]string

	addresses[0] = "odin13jgurnd72ftvy79tjtl3s8e6xvsuupwreft506"
	addresses[1] = "odin1lqk7gsq40dskpcukavkwuh7t73cnh9tjfgqxkp"
	addresses[2] = "odin1qerev5feaft35fp2n3ept7fdsheud52z0gzne0"
	addresses[3] = "odin14p9vgtynfy394hmz0tcrrs78e4whj5z6kmp09v"
	addresses[4] = "odin122qmr2s3583msah5fk5jwc7557v30kn0pj50sd"
	addresses[5] = "odin1vmfekljy9haqgyzsfvyn85xa8xdd8wcum33cf4"
	addresses[6] = "odin1vu0kf5ztrscm9fvc0gg3nqdfj7rfr6sc2u9r0v"
	addresses[7] = "odin1r532h9h54eaylxg49ll6vlr8epcl687s2x9ln7"
	addresses[8] = "odin1nvzkd37yqhw9pn9eljf3mvrneug3hr0r3xyc62"
	addresses[9] = "odin16hxrt4scaly02caaskhe0984rzl0fj490c853y"
	addresses[10] = "odin1n2apxttn8f3uzrrzrpgkrxascwdhcwh3uyejxe"
	addresses[11] = "odin1sm56gxdlwd32dzcps6uzprudvawuqg6mv83eq9"
	addresses[12] = "odin1n8gwpl4s75qhhvtlmlp5y5946acjnlutv4krud"
	addresses[13] = "odin1qe8v4nx9l06x4m3cy3wk7al6h4ww7js5rgpgp0"
	addresses[14] = "odin1vu0kf5ztrscm9fvc0gg3nqdfj7rfr6sc2u9r0v"

	var accAddresses []sdk.AccAddress

	for _, addr := range addresses {
		accAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			panic(fmt.Sprintf("account address is not valid bech32: %s", accAddr))
		}
		accAddresses = append(accAddresses, accAddr)

	}
	return accAddresses
}

func sumBalances(
	ctx sdk.Context,
	bankkeeper bankkeeper.Keeper,
	addresses []sdk.AccAddress,
) sdk.Coins {
	totalCoins := sdk.NewCoins()
	for _, addr := range addresses {
		balance := getBalance(ctx, bankkeeper, addr)

		for _, coin := range balance {
			totalCoins = totalCoins.Add(coin)
		}

	}
	return totalCoins
}

func mintAndSendCoins(
	ctx sdk.Context,
	bankkeeper bankkeeper.Keeper,
	mintkeeper mintkeeper.Keeper,
	addr sdk.AccAddress,
	coins sdk.Coins,
) error {
	// mint coins
	err := mintkeeper.MintCoins(ctx, coins)
	if err != nil {
		return err
	}

	// send coins to new address
	err = bankkeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	if err != nil {
		return err
	}

	return nil
}

func burnCoins(
	ctx sdk.Context,
	bankkeeper bankkeeper.Keeper,
	addresses []sdk.AccAddress,
) error {
	for _, addr := range addresses {
		balance := getBalance(ctx, bankkeeper, addr)

		for _, coin := range balance {
			if !coin.Amount.IsZero() {
				err := bankkeeper.BurnCoins(ctx, minttypes.ModuleName, balance)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CreateUpgradeHandler(mm module.Manager, configurator module.Configurator, bankkeeper bankkeeper.Keeper, mintkeeper mintkeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		ctx.Logger().Info("getting all addresses")
		addresses := getAddresses()

		ctx.Logger().Info("getting all account balances")
		totalCoins := sumBalances(ctx, bankkeeper, addresses)

		// Getting accound address of new address
		newAddr, err := sdk.AccAddressFromBech32(newAddress)
		if err != nil {
			panic(fmt.Sprintf("account address is not valid bech32: %s", newAddr))
		}

		ctx.Logger().Info("minting new coins and sending them to new address")
		mintAndSendCoins(ctx, bankkeeper, mintkeeper, newAddr, totalCoins)

		ctx.Logger().Info("Burning coins from old addresses")
		burnCoins(ctx, bankkeeper, addresses)

		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return newVM, err
		}
		return newVM, err
	}
}
