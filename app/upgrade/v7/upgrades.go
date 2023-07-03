package v7

import (
	"fmt"

	mintkeeper "github.com/ODIN-PROTOCOL/odin-core/x/mint/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const newAddress = "odin17zhnwfs7rh78kz628l2mxjt0u6456rznjxyu6f"

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
	var addresses [15][15]string

	addresses[0][0] = "odin13jgurnd72ftvy79tjtl3s8e6xvsuupwreft506"  //continuous vesting account
	addresses[0][1] = "odin18ffxuvq37jj6dwrsrd3qs0dy2s2tg3htz6h9ed"  //new account
	addresses[1][0] = "odin1lqk7gsq40dskpcukavkwuh7t73cnh9tjfgqxkp"  //base account
	addresses[1][1] = "odin1sz7e7e90v7jsph93h46wz07sevvun2lwhcu4p0"  //new account
	addresses[2][0] = "odin1qerev5feaft35fp2n3ept7fdsheud52z0gzne0"  //base account
	addresses[2][1] = "odin1n27lkvt7t5cadsqn2320tpc0hsa2v8gxa35qr6"  //new account
	addresses[3][0] = "odin14p9vgtynfy394hmz0tcrrs78e4whj5z6kmp09v"  //continuous vesting account
	addresses[3][1] = "odin1wf2uw3jd2cj7k73qyd2hz6sqgmkxvvm2heygac"  //new account
	addresses[4][0] = "odin122qmr2s3583msah5fk5jwc7557v30kn0pj50sd"  //continuous vesting account
	addresses[4][1] = "odin17x2h8z50r5qkzlkqd56j6kwhun7ef0jy7kyktf"  //new account
	addresses[5][0] = "odin1vmfekljy9haqgyzsfvyn85xa8xdd8wcum33cf4"  //base account
	addresses[5][1] = "odin1jrtgrk3u808f7snxm63n5s8nqr6c2prfa3ttn5"  //new account
	addresses[6][0] = "odin1vu0kf5ztrscm9fvc0gg3nqdfj7rfr6sc2u9r0v"  //continuous vesting account
	addresses[6][1] = "odin1mc0tkfnkfjpal62fkpfn2g0j5s3j69vat4dw4p"  //new account
	addresses[7][0] = "odin1r532h9h54eaylxg49ll6vlr8epcl687s2x9ln7"  //continuous vesting account
	addresses[7][1] = "odin1npu6zn9mrztwkq2pr27vd3zzq485hqcg9umna8"  //new account
	addresses[8][0] = "odin1nvzkd37yqhw9pn9eljf3mvrneug3hr0r3xyc62"  //continuous vesting account
	addresses[8][1] = "odin1ye36q77r23mnlps96wr5ckgy83x5emhazx7m3f"  //new account
	addresses[9][0] = "odin16hxrt4scaly02caaskhe0984rzl0fj490c853y"  //continuous vesting account
	addresses[9][1] = "odin19mse0ejcpqc2wtyqgyj3rnc6qn54fn987g2pq9"  //new account
	addresses[10][0] = "odin1n2apxttn8f3uzrrzrpgkrxascwdhcwh3uyejxe" //continuous vesting account
	addresses[10][1] = "odin1hqpcmxrrfyqz7r895vgjl32u9yvy4vzezupk2s" //new account
	addresses[11][0] = "odin1sm56gxdlwd32dzcps6uzprudvawuqg6mv83eq9" //continuous vesting account
	addresses[11][1] = "odin1mt98jffmfaa3g25yp6q00l7lk6xrfqshsjjkdv" //new account
	addresses[12][0] = "odin1n8gwpl4s75qhhvtlmlp5y5946acjnlutv4krud" //continuous vesting account
	addresses[12][1] = "odin18925q6zxy7s7x7jcuru2zavypp5ll7acs2g923" //new account
	addresses[13][0] = "odin1qe8v4nx9l06x4m3cy3wk7al6h4ww7js5rgpgp0" //continuous vesting account
	addresses[13][1] = "odin14mhwrtxluyzlvy98gpcld2m9qtk6qahgx3ncnj" //new account
	addresses[14][0] = "odin1vu0kf5ztrscm9fvc0gg3nqdfj7rfr6sc2u9r0v" //continuous vesting account
	addresses[14][1] = "odin1c43q5hyrs5jsu303z8h44wv36f7l0fgx46m402" //new account

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
) error {

	//send coins to new address
	err := bankkeeper.SendCoins(ctx, fromAddr, toAddr, coins)
	if err != nil {
		panic(fmt.Sprintf("Could not send coins from: %s, to: %s, error: %s", fromAddr, toAddr, err))
	}

	return nil
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
				sendCoins(ctx, bankkeeper, addrs[0], addrs[1], balance)
			}

		}

		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return newVM, err
		}
		return newVM, err
	}
}
