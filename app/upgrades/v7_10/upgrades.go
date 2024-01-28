package v7_10

import (
	"encoding/base64"
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/ODIN-PROTOCOL/odin-core/app/keepers"
	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Old addresses
const DefiantLabOldAccAddress = "odin16dcwlyrwx8duucsr363zqqsf2prc5gv52uv6zk"
const OdinMainnet3OldAccAddress = "odin1s0p07h5n4v2nqh0jr2gprq5cphv2mgs9twppcx"

// New addresses
const DefiantLabNewAccAddress = "odin1dnmz4yzv73lr3lmauuaa0wpwn8zm8s20fyv396"
const OdinMainnet3NewAccAddress = "odin1hgdq6yekx3hpz5mhph660el664pc02a4npxdas"

// PubKeys
const DefiantLabsValPubKey = "+hZsfi4r1OdyIgkZBbQgCDiADkQWlzN0iQ3Szr9+Dp8="
const OdinMainnet3ValPubKey = "AgDCT69kjqdPVFSluW23HFIpqQi7bknscN3y1YA9EB57"


func getBalance(
	ctx sdk.Context,
	sk stakingkeeper.Keeper,
	ak keeper.AccountKeeper,
	bk bankkeeper.Keeper,
	addr sdk.AccAddress,
) (sdk.Coins, error) {
	// Get all delegator delegations for address
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


func createValidator(ctx sdk.Context, stakingkeeper stakingkeeper.Keeper, address string, pubKey cryptotypes.PubKey, description stakingtypes.Description, comission stakingtypes.Commission) (error){

    valAddr := sdk.ValAddress(address)
    minSelfDelegation := sdk.OneInt()

    // Create the validator
    validator, err := stakingtypes.NewValidator(valAddr, pubKey, description)
	if err != nil {
		return err
	}

    validator.MinSelfDelegation = minSelfDelegation

    // Set the validator in the store
    stakingkeeper.SetValidator(ctx, validator)
	return nil
}


func addrToValAddr(address string) sdk.ValAddress {
	valAddr, err := sdk.ValAddressFromBech32(address)
	if err != nil {
		panic(fmt.Sprintf("account address is not valid bech32: %s", valAddr))
	}
	return valAddr
}

func moveValidatorDelegations(ctx sdk.Context, k stakingkeeper.Keeper, oldValAddress sdk.ValAddress, newValAddress sdk.ValAddress) {
	for _, delegation := range k.GetValidatorDelegations(ctx, oldValAddress) {
		newDelegation := stakingtypes.Delegation{DelegatorAddress: delegation.DelegatorAddress, ValidatorAddress: newValAddress.String(), Shares: delegation.Shares}

		k.SetDelegation(ctx, newDelegation)
		k.RemoveDelegation(ctx, delegation)
	}
}


func moveDelegations(ctx sdk.Context, k stakingkeeper.Keeper, oldAddress sdk.AccAddress, newAccAddress sdk.AccAddress) {
	for _, delegation := range getDelegations(ctx, k, oldAddress) {
		newDelegation := stakingtypes.Delegation{DelegatorAddress: newAccAddress.String(), ValidatorAddress: delegation.ValidatorAddress, Shares: delegation.Shares}

		k.SetDelegation(ctx, newDelegation)
		k.RemoveDelegation(ctx, delegation)
	}
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

func getAddresses() [][]sdk.AccAddress {
	var addresses [][]string
	addresses = append(addresses, []string{DefiantLabOldAccAddress, DefiantLabNewAccAddress})
	addresses = append(addresses, []string{OdinMainnet3OldAccAddress, OdinMainnet3NewAccAddress})

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


func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	am upgrades.AppManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running v7_10 upgrade handler")

		addresses := getAddresses()

		ctx.Logger().Info("sending coins from addresses to new address")

		// Creating DefiantLabs validator
		DanPubKeyBytes, err := base64.StdEncoding.DecodeString(DefiantLabsValPubKey)
		if err != nil {
			log.Fatalf("Failed to decode base64 string: %v", err)
		}
		
		DanPubKey := ed25519.PubKey{
			Key: DanPubKeyBytes,
		}
		
		DanValDescription := stakingtypes.NewDescription("defiantlabs", "", "https://defiantlabs.net", "", "")
		DanValComission := stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.100000000000000000"), sdk.MustNewDecFromStr("0.200000000000000000"), sdk.MustNewDecFromStr("0.010000000000000000"))
		createValidator(ctx, *keepers.StakingKeeper, string(addrToValAddr(DefiantLabNewAccAddress)), &DanPubKey, DanValDescription, DanValComission)

		// Creating new Mainnet3 validator
		Odin3PubKeyBytes, err := base64.StdEncoding.DecodeString(OdinMainnet3ValPubKey)
		if err != nil {
			log.Fatalf("Failed to decode base64 string: %v", err)
		}
		
		Odin3PubKey := ed25519.PubKey{
			Key: Odin3PubKeyBytes,
		}
		
		Odin3ValDescription := stakingtypes.NewDescription("Odin Maiinet-3", "", "https://odinprotocol.io", "", "")
		Odin3ValComission := stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.050000000000000000"), sdk.MustNewDecFromStr("0.100000000000000000"), sdk.MustNewDecFromStr("0.050000000000000000"))
		createValidator(ctx, *keepers.StakingKeeper, string(addrToValAddr(OdinMainnet3NewAccAddress)), &Odin3PubKey, Odin3ValDescription, Odin3ValComission)


		for _, addrs := range addresses {

			ctx.Logger().Info(fmt.Sprintf("Sending tokens from %s to %s", addrs[0], addrs[1]))
			balance, err := getBalance(ctx, *keepers.StakingKeeper, keepers.AccountKeeper, keepers.BankKeeper, addrs[0])

			if err != nil {
				panic("Could not send coins")
			}
			
			// sending balances
			sendCoins(ctx, keepers.BankKeeper, addrs[0], addrs[1], balance)
			
			// moving delegations
			ctx.Logger().Info(fmt.Sprintf("Moving account  delegations from %s to %s", addrs[0], addrs[1]))
			moveDelegations(ctx, *keepers.StakingKeeper, addrs[0], addrs[1])

			currentValidatorAddress := addrToValAddr(addrs[0].String())
			newValidatorAddress := addrToValAddr(addrs[1].String())

			ctx.Logger().Info(fmt.Sprintf("Moving validator delegations from %s to %s", currentValidatorAddress, newValidatorAddress))
			moveValidatorDelegations(ctx, *keepers.StakingKeeper, currentValidatorAddress, newValidatorAddress)
		}

		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}
		return newVM, err
	}
}


var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v7_10",
	CreateUpgradeHandler: CreateUpgradeHandler,

}
