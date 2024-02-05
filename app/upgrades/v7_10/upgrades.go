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
	distrbutionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
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
const OdinMainnet3ValPubKey = "FQf4cxaS5XNv+mFEi6dtDQDOLUWVWfEyh8SqljsJz1s="


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


func createValidator(ctx sdk.Context, stakingkeeper stakingkeeper.Keeper, address string, pubKey cryptotypes.PubKey, description stakingtypes.Description, comission stakingtypes.Commission) (stakingtypes.Validator,  error){

    valAddr := sdk.ValAddress(address)
    minSelfDelegation := sdk.OneInt()

    // Create the validator
    validator, err := stakingtypes.NewValidator(valAddr, pubKey, description)
	if err != nil {
		log.Printf("Error when creating a validator %v: %s", valAddr, err)
		return stakingtypes.Validator{}, err
	}

    validator.MinSelfDelegation = minSelfDelegation

    // Set the validator in the store
    stakingkeeper.SetValidator(ctx, validator)
	return validator, nil
}


func withdrawRewardsAndCommission(ctx sdk.Context, sk stakingkeeper.Keeper, dk distrbutionkeeper.Keeper, oldValAddress sdk.ValAddress, newValAddress sdk.ValAddress) {
	oldValAccAddress := sdk.AccAddress(oldValAddress)
	newValAccAddress := sdk.AccAddress(newValAddress)

	// withdrawing all rewards, self-delegation rewards mapped to new account
	for _, delegation := range sk.GetValidatorDelegations(ctx, oldValAddress) {
		withdrawAddress := dk.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(delegation.DelegatorAddress))
		delegatorAddress := delegation.GetDelegatorAddr()
		
		// we suppose that old Odin accounts are unavailable, so we're routing rewards to new addresses and proceeding wit hwithdraws
		if withdrawAddress.String() == oldValAccAddress.String() {
			log.Printf("Found delegation which withdrawal address is the old one: %v. Setting withdrawal address to new account: %v", oldValAccAddress.String(), newValAccAddress.String())
			dk.SetDelegatorWithdrawAddr(ctx, delegatorAddress, newValAccAddress)
		}
		
		log.Printf("Withdrawing reward for %v delegator address from %v", delegatorAddress.String(), oldValAddress.String())
		dk.WithdrawDelegationRewards(ctx, delegatorAddress, oldValAddress)
	}

	// Comission
	// explicitly setting validator withdrawal address, in case it has no self-delegation in the loop above
	dk.SetDelegatorWithdrawAddr(ctx, oldValAccAddress, newValAccAddress)
	dk.WithdrawValidatorCommission(ctx, oldValAddress)
}

func addrToValAddr(address string) (sdk.ValAddress, error) {
	bytes, err := sdk.GetFromBech32(address, "odin")
	if err != nil {
		log.Printf("account address %s is not valid bech32: %s", address, err)
		return nil, err
	}

	valAddr := sdk.ValAddress(bytes)
	return valAddr, nil
}

func moveValidatorDelegations(ctx sdk.Context, k stakingkeeper.Keeper, oldValAddress sdk.ValAddress, newValAddress sdk.ValAddress) {
	for _, delegation := range k.GetValidatorDelegations(ctx, oldValAddress) {
		newDelegation := stakingtypes.Delegation{DelegatorAddress: delegation.DelegatorAddress, ValidatorAddress: newValAddress.String(), Shares: delegation.Shares}
		log.Printf("Moving validator delegation from %v to %v", oldValAddress,  newDelegation.ValidatorAddress)
		
		k.SetDelegation(ctx, newDelegation)
		k.RemoveDelegation(ctx, delegation)
	}
}

func moveDelegations(ctx sdk.Context, k stakingkeeper.Keeper, oldAddress sdk.AccAddress, newAccAddress sdk.AccAddress) {
	for _, delegation := range getDelegations(ctx, k, oldAddress) {
		newDelegation := stakingtypes.Delegation{DelegatorAddress: newAccAddress.String(), ValidatorAddress: delegation.ValidatorAddress, Shares: delegation.Shares}
		log.Printf("Moving validator delegation from %v to %v", delegation.ValidatorAddress, newDelegation.ValidatorAddress)
		
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
) (error) {
	//send coins to new address
	err := bankkeeper.SendCoins(ctx, fromAddr, toAddr, coins)
	if err != nil {
		log.Printf("Could not send coins from: %s, to: %s, error: %s", fromAddr, toAddr, err)
		return err
	}
	return nil
}

func getAddresses() ([][]sdk.AccAddress, error) {
	var addresses [][]string
	addresses = append(addresses, []string{DefiantLabOldAccAddress, DefiantLabNewAccAddress})
	addresses = append(addresses, []string{OdinMainnet3OldAccAddress, OdinMainnet3NewAccAddress})

	var accAddresses [][]sdk.AccAddress

	for _, addrs := range addresses {
		var accaddrs []sdk.AccAddress
		for _, addr := range addrs {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				log.Printf("account address is not valid bech32: %s: %s", accAddr, err)
				return nil, err
			}
			accaddrs = append(accaddrs, accAddr)
		}
		accAddresses = append(accAddresses, accaddrs)
	}
	return accAddresses, nil
}


func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	am upgrades.AppManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running v7_10 upgrade handler")

		addresses, err := getAddresses()
		if err != nil {
			log.Printf("Error when retrieving addresses from getAddresses: %s", err)
			return nil, err
		}

		ctx.Logger().Info("sending coins from addresses to new address")
		
		// Creating DefiantLabs validator
		DanOldValAddress, err := addrToValAddr(DefiantLabOldAccAddress)
		if err != nil {
			return nil, err
		}
		DanOldVal, found := keepers.StakingKeeper.GetValidator(ctx, DanOldValAddress)
		if !found {
			log.Printf("Validator with %v has not been found", DanOldValAddress)
			return nil, err
		}

		DanPubKeyBytes, err := base64.StdEncoding.DecodeString(DefiantLabsValPubKey)
		if err != nil {
			log.Printf("Failed to decode base64 string: %v", err)
			return nil, err
		}
		
		DanPubKey := ed25519.PubKey{
			Key: DanPubKeyBytes,
		}
		
		DanValAddr, err := addrToValAddr(DefiantLabNewAccAddress)
		if err != nil {
			return nil, err
		}

		DanNewVal, err := createValidator(ctx, *keepers.StakingKeeper, string(DanValAddr), &DanPubKey, DanOldVal.Description, DanOldVal.Commission)
		if err != nil {
			return nil, err
		}

		ctx.Logger().Info(fmt.Sprintf("Moving validator delegations from %s to %s", DanOldValAddress, DanValAddr))
		moveValidatorDelegations(ctx, *keepers.StakingKeeper, DanOldValAddress, DanValAddr)
		DanNewVal.UpdateStatus(stakingtypes.Bonded)

		// Creating new Mainnet3 validator
		Odin3OldValAddress, err := addrToValAddr(OdinMainnet3OldAccAddress)
		if err != nil {
			return nil, err
		}
		Odin3OldVal, found := keepers.StakingKeeper.GetValidator(ctx, Odin3OldValAddress)
		if !found {
			log.Printf("Validator with %v has not been found", Odin3OldValAddress)
			return nil, err
		}

		Odin3PubKeyBytes, err := base64.StdEncoding.DecodeString(OdinMainnet3ValPubKey)
		if err != nil {
			log.Printf("Error whend decoding public key from string %v", err)
			return nil, err
		}
		
		Odin3PubKey := ed25519.PubKey{
			Key: Odin3PubKeyBytes,
		}
		Odin3ValAddr, err := addrToValAddr(OdinMainnet3NewAccAddress)
		if err != nil {
			return nil, err
		}

		Odin3Val, err := createValidator(ctx, *keepers.StakingKeeper, string(Odin3ValAddr), &Odin3PubKey, Odin3OldVal.Description, Odin3OldVal.Commission)
		if err != nil {
			return nil, err
		}

		ctx.Logger().Info(fmt.Sprintf("Moving validator delegations from %s to %s",  Odin3OldValAddress, Odin3ValAddr))
		moveValidatorDelegations(ctx, *keepers.StakingKeeper, Odin3OldValAddress, Odin3ValAddr)
		Odin3Val.UpdateStatus(stakingtypes.Bonded)
		
		// rewards and comission
		withdrawRewardsAndCommission(ctx, *keepers.StakingKeeper, keepers.DistrKeeper, DanOldValAddress, DanValAddr)
		withdrawRewardsAndCommission(ctx, *keepers.StakingKeeper, keepers.DistrKeeper, Odin3OldValAddress, Odin3ValAddr)

		for _, addrs := range addresses {

			ctx.Logger().Info(fmt.Sprintf("Sending tokens from %s to %s", addrs[0], addrs[1]))
			balance, err := getBalance(ctx, *keepers.StakingKeeper, keepers.AccountKeeper, keepers.BankKeeper, addrs[0])

			if err != nil {
				log.Printf("Error when retrieving balance for address %s: %s",  addrs[0], err)
				return nil, err
			}
			
			// sending balances
			sendCoins(ctx, keepers.BankKeeper, addrs[0], addrs[1], balance)
						
			// moving delegations
			ctx.Logger().Info(fmt.Sprintf("Moving account  delegations from %s to %s", addrs[0], addrs[1]))
			moveDelegations(ctx, *keepers.StakingKeeper, addrs[0], addrs[1])
		}
		
		DanOldVal.UpdateStatus(stakingtypes.Unbonded)
		Odin3OldVal.UpdateStatus(stakingtypes.Unbonded)

		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			log.Printf("Error when running migrations: %s", err)
			return nil, err
		}
		return newVM, err
	}
}


var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v0.7.10",
	CreateUpgradeHandler: CreateUpgradeHandler,
}
