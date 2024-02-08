package v7_10

import (
	"encoding/base64"
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/ODIN-PROTOCOL/odin-core/app/keepers"
	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Old addresses
const DefiantLabAccAddress = "odin16dcwlyrwx8duucsr363zqqsf2prc5gv52uv6zk"
const OdinMainnet3OldAccAddress = "odin1s0p07h5n4v2nqh0jr2gprq5cphv2mgs9twppcx"

// New addresses
const DefiantLabOldAccAddress = "odin1dnmz4yzv73lr3lmauuaa0wpwn8zm8s20fyv396"
const OdinMainnet3NewAccAddress = "odin1hgdq6yekx3hpz5mhph660el664pc02a4npxdas"

// PubKeys
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


func createValidator(ctx sdk.Context, sk stakingkeeper.Keeper, dk distributionkeeper.Keeper, address string, pubKey cryptotypes.PubKey, description stakingtypes.Description, comission stakingtypes.Commission) (stakingtypes.Validator,  error){

    valAddr := sdk.ValAddress(address)
    minSelfDelegation := sdk.OneInt()

    // Create the validator
    validator, err := stakingtypes.NewValidator(valAddr, pubKey, description)
	if err != nil {
		log.Printf("Error when creating a validator %v: %s", valAddr, err)
		return stakingtypes.Validator{}, err
	}

    validator.MinSelfDelegation = minSelfDelegation
	validator.Status = stakingtypes.Bonded
    validator.Tokens = sdk.ZeroInt()
	validator.DelegatorShares = sdk.ZeroDec()
	validator.Commission = comission

	// Update validators in the store
	sk.SetValidator(ctx, validator)

	err = sk.Hooks().AfterValidatorCreated(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}

	err = dk.Hooks().AfterValidatorCreated(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}

	return validator, nil
}


func withdrawRewardsAndCommission(ctx sdk.Context, sk stakingkeeper.Keeper, dk distributionkeeper.Keeper,  oldValAddress sdk.ValAddress, newValAddress sdk.ValAddress) {
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


func moveValidatorDelegations(ctx sdk.Context, k stakingkeeper.Keeper, d distributionkeeper.Keeper, oldVal stakingtypes.Validator, newVal stakingtypes.Validator) (error) {
	cumOldValShares := math.LegacyZeroDec()

	for _, delegation := range k.GetValidatorDelegations(ctx, oldVal.GetOperator()) {
		log.Printf("Moving validator delegation from %v to %v", delegation.DelegatorAddress,  newVal.OperatorAddress)
		// Remove the delegation to the old validator
		k.RemoveDelegation(ctx, delegation)

		// Create a new delegation to the new validator
		newDelegation := stakingtypes.Delegation{
			DelegatorAddress: delegation.DelegatorAddress,
			ValidatorAddress: newVal.OperatorAddress,
			Shares:           delegation.Shares,
		}
		
		err := k.Hooks().BeforeDelegationCreated(ctx, delegation.GetDelegatorAddr(), newVal.GetOperator()) 
		if err != nil {
			log.Printf("Error when running hook after adding delegation %v to %v", delegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}
		k.SetDelegation(ctx, newDelegation)	

		err = d.Hooks().BeforeDelegationCreated(ctx, delegation.GetDelegatorAddr(), newVal.GetOperator())
		if err != nil {
			log.Printf("Error when running hook after addig delegation %v to %v", delegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}
		cumOldValShares = cumOldValShares.Add(delegation.Shares)
	}
	
	// tokens
	tokens := oldVal.TokensFromShares(cumOldValShares) 
	
	k.AddValidatorTokensAndShares(ctx, newVal, tokens.TruncateInt())
	k.RemoveValidatorTokensAndShares(ctx, oldVal, cumOldValShares)

	k.SetValidatorByConsAddr(ctx, oldVal)
	k.SetValidatorByConsAddr(ctx, newVal)

	return nil
}

func moveDelegations(ctx sdk.Context, keepers *keepers.AppKeepers, oldAddress sdk.AccAddress, newVal stakingtypes.Validator) error {
	for _, delegation := range getDelegations(ctx, *keepers.StakingKeeper, oldAddress) {
		log.Printf("Moving delegation from %v to %v", delegation.DelegatorAddress,  newVal.OperatorAddress)
		keepers.StakingKeeper.RemoveDelegation(ctx, delegation)

		newDelegation := stakingtypes.Delegation{
			DelegatorAddress: delegation.DelegatorAddress,
			ValidatorAddress: newVal.OperatorAddress,
			Shares:           delegation.Shares,
		}
		
		err := keepers.StakingKeeper.Hooks().BeforeDelegationCreated(ctx, delegation.GetDelegatorAddr(), newVal.GetOperator()) 
		if err != nil {
			log.Printf("Error when running hook after adding delegation %v to %v", delegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}
		keepers.StakingKeeper.SetDelegation(ctx, newDelegation)	

		err = keepers.DistrKeeper.Hooks().BeforeDelegationCreated(ctx, delegation.GetDelegatorAddr(), newVal.GetOperator())
		if err != nil {
			log.Printf("Error when running hook after addig delegation %v to %v", delegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}
	}
	return nil
}


func forceCompleteUnbonding(ctx sdk.Context, keepers *keepers.AppKeepers, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	ubd, found := keepers.StakingKeeper.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoUnbondingDelegation
	}

	bondDenom := keepers.StakingKeeper.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()

	delegatorAddress, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if !entry.OnHold() {
			ubd.RemoveEntry(int64(i))
			i--
			keepers.StakingKeeper.DeleteUnbondingIndex(ctx, entry.UnbondingId)

			// track undelegation only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewCoin(bondDenom, entry.Balance)
				if err := keepers.BankKeeper.UndelegateCoinsFromModuleToAccount(
					ctx, stakingtypes.NotBondedPoolName, delegatorAddress, sdk.NewCoins(amt),
				); err != nil {
					return nil, err
				}

				balances = balances.Add(amt)
			}
		}
	}

	// set the unbonding delegation or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		keepers.StakingKeeper.RemoveUnbondingDelegation(ctx, ubd)
	} else {
		keepers.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
	}

	return balances, nil
}

func moveSelfDelegation(ctx sdk.Context, keepers  *keepers.AppKeepers, oldDelegatorAddress sdk.AccAddress, newDelegatorAddress sdk.AccAddress, validatorAddr sdk.ValAddress) error {
    stakingKeeper := keepers.StakingKeeper
	
	// Get the delegation from the old validator
    delegation, found := stakingKeeper.GetDelegation(ctx, oldDelegatorAddress, validatorAddr)
    if !found {
		log.Printf("self delegation not found: %s", oldDelegatorAddress)
        return fmt.Errorf("self delegation not found")
    }
	
	// stakingKeeper.RemoveDelegation(ctx, delegation)

	// Unbonding
	tm, err := stakingKeeper.Undelegate(ctx, oldDelegatorAddress, validatorAddr, delegation.Shares)
	// amt, err := stakingKeeper.Unbond(ctx, oldDelegatorAddress, validatorAddr, delegation.Shares) 
	log.Printf("Unbonding time from %v: %v", tm, oldDelegatorAddress)
	if err != nil{
		log.Printf("Error when unbonding %v: %v",  oldDelegatorAddress, err)
		return err
	}

	selfBond, err := forceCompleteUnbonding(ctx, keepers, oldDelegatorAddress,  validatorAddr)
	log.Printf("Unbonding finished with amount from %v: %v", selfBond, oldDelegatorAddress)

	if err != nil{
		log.Printf("Error when unbonding %v: %v",  oldDelegatorAddress, err)
		return err
	}

	balance, err := getBalance(ctx, *keepers.StakingKeeper, keepers.AccountKeeper, keepers.BankKeeper, oldDelegatorAddress)
	
	if err != nil {
		log.Printf("Error when retrieving balance for address %s: %s",  oldDelegatorAddress, err)
		return err
	}

	log.Printf("Balance after unbonding %s: %s",  oldDelegatorAddress, balance)
	log.Printf("Sending coins from %v to %v (%v)", oldDelegatorAddress, newDelegatorAddress, balance)

	err = sendCoins(ctx, keepers.BankKeeper, oldDelegatorAddress, newDelegatorAddress, balance)
	if err != nil {
		log.Printf("Error when sending coins for address %s: %s",  oldDelegatorAddress, err)
		return err
	}

    validator, found := stakingKeeper.GetValidator(ctx, validatorAddr)
    if !found {
        return fmt.Errorf("old validator not found")
    }
	
 	SelfDelegate(ctx, *keepers.StakingKeeper, keepers.BankKeeper, newDelegatorAddress, validator, selfBond)
	if err != nil {
		log.Printf("Error when self delegating to %v", newDelegatorAddress)
		return err
	}


    // // Create a new delegation to the new validator
    // newDelegation := stakingtypes.Delegation{
    //     DelegatorAddress: newDelegatorAddress.String(),
    //     ValidatorAddress: validatorAddr.String(),
    //     Shares:           amount,
    // }
	
	// err = stakingKeeper.Hooks().BeforeDelegationCreated(ctx, delegation.GetDelegatorAddr(), validatorAddr) 
	// if err != nil {
	// 	log.Printf("Error when running hook before adding delegation %v to %v", delegation.GetDelegatorAddr(), validatorAddr)
	// 	return err
	// }
    // // Save the new delegation
    // stakingKeeper.SetDelegation(ctx, newDelegation)

    // Update the old validator's and new validator's tokens and delegator shares

    // Convert shares to tokens if necessary. This example assumes `amount` is already in tokens.
    // You might need to adjust based on your actual data and the SDK's version.
    // tokens := sdk.TokensFromConsensusPower(int64(amount.TruncateInt64()), sdk.DefaultPowerReduction)

    // oldValidator.Tokens = oldValidator.Tokens.Sub(tokens)
    // newValidator.Tokens = newValidator.Tokens.Add(tokens)

    // Ensure to update the validators' power index if their tokens have changed
	
	validator.Jailed = false
    stakingKeeper.SetValidator(ctx, validator)
    // stakingKeeper.SetValidator(ctx, newValidator)

    // This operation should ideally be followed by updates to the pool's bonded tokens
    // and triggering any necessary events or hooks for state consistency

    return nil
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


func SelfDelegate(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, bankKeeper bankkeeper.Keeper, delegatorAddr sdk.AccAddress, validator stakingtypes.Validator, amount sdk.Coins) error {
    // Delegate tokens to the validator
	for _, balance := range amount {
		
		// Ensure the delegator (validator account) has enough balance for the delegation		
		if !bankKeeper.HasBalance(ctx, delegatorAddr, balance) {
			return sdkerrors.Wrapf(errortypes.ErrInsufficientFunds, "not enough balance to self-delegate to validator: %s", validator.OperatorAddress)
		}
		
		// Send coins from the delegator's account to the module account (staking module account) as part of delegation
		err := bankKeeper.SendCoinsFromAccountToModule(ctx, delegatorAddr, stakingtypes.NotBondedPoolName, amount)
		if err != nil {
			return err
		}

		_, err = stakingKeeper.Delegate(ctx, delegatorAddr, balance.Amount, stakingtypes.Unbonded, validator, true)
		if err != nil {
			return err
		}
	}
    return nil
}


func fixDefiantLabs(ctx sdk.Context, keepers *keepers.AppKeepers) (error) {
	

	// Fixing self delegation
	DefiantLabsValAddress, err := addrToValAddr(DefiantLabAccAddress)
	if err != nil {
		return err
	}

	DefiantLabsVal, found := keepers.StakingKeeper.GetValidator(ctx, DefiantLabsValAddress)
	if !found {
		log.Printf("Validator with %v has not been found", DefiantLabsValAddress)
		return err
	}

	DefiantLabsOldAcc, err := sdk.AccAddressFromBech32(DefiantLabOldAccAddress)
	if err != nil {
		log.Printf("account address is not valid bech32: %s: %s", DefiantLabAccAddress, err)
		return err
	}

	DefiantLabsAcc, err := sdk.AccAddressFromBech32(DefiantLabAccAddress)
	if err != nil {
		log.Printf("account address is not valid bech32: %s: %s", DefiantLabAccAddress, err)
		return err
	}
	
	// Moving DefiantLabs self-delegation
	err = moveSelfDelegation(ctx, keepers, DefiantLabsOldAcc, DefiantLabsAcc,  DefiantLabsValAddress)
	if err != nil {
		return err
	}

	// Donating 10K loki to DefiantLab to automate his self-delegation process
	
	// OdinMainnet3DonorAddr, err := sdk.AccAddressFromBech32(OdinMainnet3OldAccAddress)
	// if err != nil {
	// 	log.Printf("account address is not valid bech32: %s: %s", OdinMainnet3OldAccAddress, err)
	// 	return err
	// }

	// log.Printf("Donating coins for self-delegation %v", DefiantLabsValAddress)
	// sendCoins(ctx, keepers.BankKeeper, OdinMainnet3DonorAddr, DefiantLabsAcc, sdk.NewCoins(sdk.NewCoin("loki", sdk.NewInt(10000))))
	
	// keepers.DistrKeeper.SetWithdrawAddr(ctx, DefiantLabsAcc, DefiantLabsAcc)
	
	// // sending balances
	// ctx.Logger().Info(fmt.Sprintf("Sending tokens from %s to %s", DefiantLabOldAccAddress, DefiantLabAccAddress))

	// balance, err := getBalance(ctx, *keepers.StakingKeeper, keepers.AccountKeeper, keepers.BankKeeper, DefiantLabsOldAcc)
	// if err != nil {
	// 	log.Printf("Error when retrieving balance for address %s: %s",  DefiantLabOldAccAddress, err)
	// 	return err
	// }
	// sendCoins(ctx, keepers.BankKeeper, DefiantLabsOldAcc, DefiantLabsAcc, balance)

	// Moving delegations
	moveDelegations(ctx, keepers, DefiantLabsOldAcc, DefiantLabsVal)
	
	// Self delegating to new validator
	// selfDel := keepers.StakingKeeper.Delegation(ctx, DefiantLabsAcc, DefiantLabsValAddress)
	// if selfDel == nil {
	// 	err := SelfDelegate(ctx, *keepers.StakingKeeper, keepers.BankKeeper, DefiantLabsAcc, DefiantLabsVal, sdk.NewCoin("loki", math.NewInt(1000)))
	// 	if err != nil {
	// 		log.Printf("Error when self delegating to %v", DefiantLabsValAddress)
	// 		return err
	// 	}
	// }
	return nil
}

func fixMainnet3(ctx sdk.Context, keepers *keepers.AppKeepers) (error) { 

	// Showing all validator powers
	OldMainnet3Addr, err := sdk.AccAddressFromBech32(OdinMainnet3OldAccAddress)
	if err != nil {
		log.Printf("account address is not valid bech32: %s: %s", DefiantLabAccAddress, err)
		return err
	}
	
	NewMainnet3Addr, err := sdk.AccAddressFromBech32(OdinMainnet3NewAccAddress)
	if err != nil {
		log.Printf("account address is not valid bech32: %s: %s", DefiantLabAccAddress, err)
		return err
	}

	ctx.Logger().Info(fmt.Sprintf("Sending tokens from %s to %s", OldMainnet3Addr, NewMainnet3Addr))
	balance, err := getBalance(ctx, *keepers.StakingKeeper, keepers.AccountKeeper, keepers.BankKeeper, OldMainnet3Addr)
	if err != nil {
		log.Printf("Error when retrieving balance for address %s: %s",  OldMainnet3Addr, err)
		return err
	}

	// sending balances
	sendCoins(ctx, keepers.BankKeeper, OldMainnet3Addr, NewMainnet3Addr, balance)

	// Creating new Mainnet3 validator
	Odin3OldValAddress, err := addrToValAddr(OdinMainnet3OldAccAddress)
	if err != nil {
		return err
	}
	Odin3OldVal, found := keepers.StakingKeeper.GetValidator(ctx, Odin3OldValAddress)
	if !found {
		log.Printf("Validator with %v has not been found", Odin3OldValAddress)
		return err
	}

	Odin3PubKeyBytes, err := base64.StdEncoding.DecodeString(OdinMainnet3ValPubKey)
	if err != nil {
		log.Printf("Error whend decoding public key from string %v", err)
		return err
	}
	
	Odin3PubKey := ed25519.PubKey{
		Key: Odin3PubKeyBytes,
	}
	Odin3ValAddr, err := addrToValAddr(OdinMainnet3NewAccAddress)
	if err != nil {
		return err
	}

	Odin3Val, err := createValidator(ctx, *keepers.StakingKeeper, keepers.DistrKeeper, string(Odin3ValAddr), &Odin3PubKey, Odin3OldVal.Description, Odin3OldVal.Commission)
	if err != nil {
		return err
	}
	
	ctx.Logger().Info(fmt.Sprintf("Moving validator delegations from %s to %s",  Odin3OldValAddress, Odin3ValAddr))
	
	err = moveValidatorDelegations(ctx, *keepers.StakingKeeper, keepers.DistrKeeper, Odin3OldVal, Odin3Val)
	if err != nil {
		return err
	}
	moveDelegations(ctx, keepers, sdk.AccAddress(OdinMainnet3OldAccAddress), Odin3Val)

	Odin3Val.UpdateStatus(stakingtypes.Bonded)

	// rewards and comission
	withdrawRewardsAndCommission(ctx, *keepers.StakingKeeper, keepers.DistrKeeper, Odin3OldValAddress, Odin3ValAddr)
	Odin3OldVal.UpdateStatus(stakingtypes.Unbonded)

	// Showing all validator powers
	log.Printf("Validator power after update:")
	for _, validator := range keepers.StakingKeeper.GetAllValidators(ctx) {
		log.Printf("%v: %v", validator.OperatorAddress, validator.ConsensusPower(keepers.StakingKeeper.PowerReduction(ctx)))
	}

	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	am upgrades.AppManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running v7_10 upgrade handler")

		log.Printf("Validator power before update:")
		for _, validator := range keepers.StakingKeeper.GetAllValidators(ctx) {
			log.Printf("%v: %v", validator.OperatorAddress, validator.ConsensusPower(keepers.StakingKeeper.PowerReduction(ctx)))
		}

		// Fixinng Dan's validator account association 
		err := fixDefiantLabs(ctx, keepers)
		if err != nil {
			return nil, err
		}
		
		err = fixMainnet3(ctx, keepers)
		if err != nil {
			return nil, err
		}

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
