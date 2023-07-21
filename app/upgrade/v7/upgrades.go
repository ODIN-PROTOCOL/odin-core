package v7

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	mintkeeper "github.com/ODIN-PROTOCOL/odin-core/x/mint/keeper"
)

func getBalance(
	ctx sdk.Context,
	sk stakingkeeper.Keeper,
	ak keeper.AccountKeeper,
	bk bankkeeper.Keeper,
	addr sdk.AccAddress,
) (sdk.Coins, error) {

	//Get all delegator delegations for address
	delegations := getDelegations(ctx, sk, addr)

	for _, delegation := range delegations {
		undelegate(ctx, sk, delegation, bk)
	}

	account := ak.GetAccount(ctx, addr)

	vestingAccount, ok := account.(*vestingtypes.BaseVestingAccount)
	if !ok {
		return bk.GetAllBalances(ctx, addr), nil
	} else {
		//If the account is a vesting account, create a copy of the account
		//and vest all coins with the current block header time
		newVestingAcc := vestingtypes.NewContinuousVestingAccountRaw(vestingAccount, ctx.BlockHeader().Time.Unix())
		ak.SetAccount(ctx, newVestingAcc)
		return newVestingAcc.GetVestedCoins(ctx.BlockTime()), nil
	}
}

func getDelegations(
	ctx sdk.Context,
	stakingKeeper stakingkeeper.Keeper,
	delegatorAddr sdk.AccAddress,
) []stakingtypes.Delegation {

	delegations := stakingKeeper.GetAllDelegatorDelegations(ctx, delegatorAddr)

	return delegations
}

func undelegate(
	ctx sdk.Context,
	stakingkeeper stakingkeeper.Keeper,
	delegation stakingtypes.Delegation,
	bankkeeper bankkeeper.SendKeeper,
) {
	//Get delegator and validator addresses from the delegation
	delegatorAddr := addrToAccAddr(delegation.DelegatorAddress)
	validatorAddr := addrToValAddr(delegation.ValidatorAddress)

	//Get the amount of coins to undelegate from the delegation
	amount := sdk.NewCoin("stake", delegation.Shares.RoundInt())

	completionTime, err := stakingkeeper.Undelegate(ctx, delegatorAddr, validatorAddr, amount.Amount.ToDec())
	if err != nil {
		panic(fmt.Sprintf("Could not undelegate: %s, %s", delegation.DelegatorAddress, delegation.ValidatorAddress))
	}

	ctx.Logger().Info("Undelegation will be completet at: ", completionTime)

}

func getAddresses() [][]sdk.AccAddress {
	var addresses [][]string

	addresses = append(addresses, []string{"odin13jgurnd72ftvy79tjtl3s8e6xvsuupwreft506", "odin18u8kaxjn3pdums92exr0j6pjmt77ne7c2tkwth"})
	addresses = append(addresses, []string{"odin1lqk7gsq40dskpcukavkwuh7t73cnh9tjfgqxkp", "odin1sz7e7e90v7jsph93h46wz07sevvun2lwhcu4p0"})
	addresses = append(addresses, []string{"odin1qerev5feaft35fp2n3ept7fdsheud52z0gzne0", "odin1n27lkvt7t5cadsqn2320tpc0hsa2v8gxa35qr6"})
	addresses = append(addresses, []string{"odin14p9vgtynfy394hmz0tcrrs78e4whj5z6kmp09v", "odin1wf2uw3jd2cj7k73qyd2hz6sqgmkxvvm2heygac"})
	addresses = append(addresses, []string{"odin122qmr2s3583msah5fk5jwc7557v30kn0pj50sd", "odin17x2h8z50r5qkzlkqd56j6kwhun7ef0jy7kyktf"})
	addresses = append(addresses, []string{"odin1vmfekljy9haqgyzsfvyn85xa8xdd8wcum33cf4", "odin1jrtgrk3u808f7snxm63n5s8nqr6c2prfa3ttn5"})
	addresses = append(addresses, []string{"odin1vu0kf5ztrscm9fvc0gg3nqdfj7rfr6sc2u9r0v", "odin1mc0tkfnkfjpal62fkpfn2g0j5s3j69vat4dw4p"})
	addresses = append(addresses, []string{"odin1r532h9h54eaylxg49ll6vlr8epcl687s2x9ln7", "odin1npu6zn9mrztwkq2pr27vd3zzq485hqcg9umna8"})
	addresses = append(addresses, []string{"odin1nvzkd37yqhw9pn9eljf3mvrneug3hr0r3xyc62", "odin1ye36q77r23mnlps96wr5ckgy83x5emhazx7m3f"})
	addresses = append(addresses, []string{"odin16hxrt4scaly02caaskhe0984rzl0fj490c853y", "odin19mse0ejcpqc2wtyqgyj3rnc6qn54fn987g2pq9"})
	addresses = append(addresses, []string{"odin1n2apxttn8f3uzrrzrpgkrxascwdhcwh3uyejxe", "odin1hqpcmxrrfyqz7r895vgjl32u9yvy4vzezupk2s"})
	addresses = append(addresses, []string{"odin1sm56gxdlwd32dzcps6uzprudvawuqg6mv83eq9", "odin1mt98jffmfaa3g25yp6q00l7lk6xrfqshsjjkdv"})
	addresses = append(addresses, []string{"odin1n8gwpl4s75qhhvtlmlp5y5946acjnlutv4krud", "odin18925q6zxy7s7x7jcuru2zavypp5ll7acs2g923"})
	addresses = append(addresses, []string{"odin1qe8v4nx9l06x4m3cy3wk7al6h4ww7js5rgpgp0", "odin14mhwrtxluyzlvy98gpcld2m9qtk6qahgx3ncnj"})
	addresses = append(addresses, []string{"odin1vu0kf5ztrscm9fvc0gg3nqdfj7rfr6sc2u9r0v", "odin1c43q5hyrs5jsu303z8h44wv36f7l0fgx46m402"})

	var accAddresses [][]sdk.AccAddress

	for _, addrs := range addresses {
		var accaddrs []sdk.AccAddress
		for _, addr := range addrs {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				panic(fmt.Sprintf("account address is not valid bech32: %s", accAddr))
			}
			accaddrs = append(accaddrs, accAddr)
		}
		accAddresses = append(accAddresses, accaddrs)
	}
	return accAddresses
}

func addrToAccAddr(address string) sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(fmt.Sprintf("account address is not valid bech32: %s", accAddr))
	}
	return accAddr
}

func addrToValAddr(address string) sdk.ValAddress {
	valAddr, err := sdk.ValAddressFromBech32(address)
	if err != nil {
		panic(fmt.Sprintf("account address is not valid bech32: %s", valAddr))
	}
	return valAddr
}

// func sumBalances(
// 	ctx sdk.Context,
// 	stakingkeeper stakingkeeper.Keeper,
// 	accountkeeper keeper.AccountKeeper,
// 	bankkeeper bankkeeper.Keeper,
// 	addresses []sdk.AccAddress,
// ) sdk.Coins {
// 	totalCoins := sdk.NewCoins()
// 	for _, addr := range addresses {
// 		balance, _ := getBalance(ctx, stakingkeeper, accountkeeper, bankkeeper, addr)

// 		for _, coin := range balance {
// 			totalCoins = totalCoins.Add(coin)
// 		}

// 	}
// 	return totalCoins
// }

func sendCoins(
	ctx sdk.Context,
	bankkeeper bankkeeper.Keeper,
	fromAddr sdk.AccAddress,
	toAddr sdk.AccAddress,
	coins sdk.Coins,
) {

	//send coins to new address
	err := bankkeeper.SendCoins(ctx, fromAddr, toAddr, coins)

	if err != nil {
		panic(fmt.Sprintf("Could not send coins from: %s, to: %s, error: %s", fromAddr, toAddr, err))
	}
}

// func burnCoins(
// 	ctx sdk.Context,
// 	accountkeeper keeper.AccountKeeper,
// 	bankkeeper bankkeeper.Keeper,
// 	address sdk.AccAddress,
// 	coins sdk.Coins,
// ) error {

// 	//acc := accountkeeper.GetAccount(ctx, address)

// 	if err := bankkeeper.BurnCoins(ctx, "module_name", coins); err != nil {
// 		panic(fmt.Sprintf("Failed to burn coins on account: %s", address))
// 		return err
// 	}
// 	return nil
// }

func CreateUpgradeHandler(
	mm module.Manager,
	configurator module.Configurator,
	stakingkeeper stakingkeeper.Keeper,
	accountkeeper keeper.AccountKeeper,
	bankkeeper bankkeeper.Keeper,
	mintkeeper mintkeeper.Keeper,
) upgradetypes.UpgradeHandler {

	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		ctx.Logger().Info("getting all addresses")
		addresses := getAddresses()

		// ctx.Logger().Info("getting all account balances")
		// totalCoins := sumBalances(ctx, stakingkeeper, accountkeeper, bankkeeper, addresses)

		ctx.Logger().Info("sending coins from addresses to new address")
		for _, addrs := range addresses {
			balance, err := getBalance(ctx, stakingkeeper, accountkeeper, bankkeeper, addrs[0])
			if err != nil {
				panic(fmt.Sprintf("Could not send coins"))
			}
			sendCoins(ctx, bankkeeper, addrs[0], addrs[1], balance)
		}

		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}
		return newVM, err
	}
}
