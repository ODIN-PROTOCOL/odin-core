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
		
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}
		return newVM, err
	}
}
