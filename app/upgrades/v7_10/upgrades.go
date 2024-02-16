package v7_10

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	sdkerrors "cosmossdk.io/errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ODIN-PROTOCOL/odin-core/app/keepers"
	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const DefiantLabOldAccAddress = "odin1dnmz4yzv73lr3lmauuaa0wpwn8zm8s20fyv396"

const DefiantLabAccAddress = "odin16dcwlyrwx8duucsr363zqqsf2prc5gv52uv6zk" // Prod
//const DefiantLabAccAddress = "odin1t6hn2c9hrc33fa5slh9wtvv4ew2qhygl0rmc4q"

// const OdinMainnet3NewAccAddress = "odin1hgdq6yekx3hpz5mhph660el664pc02a4npxdas" // Test
const OdinMainnet3OldAccAddress = "odin1s0p07h5n4v2nqh0jr2gprq5cphv2mgs9twppcx"

const OdinMainnet3NewAccAddress = "odin1hgdq6yekx3hpz5mhph660el664pc02a4npxdas"

// PubKeys
const OdinMainnet3ValPubKey = "FQf4cxaS5XNv+mFEi6dtDQDOLUWVWfEyh8SqljsJz1s=" // Prod
// const OdinMainnet3ValPubKey = "f7pqqa+1Rkl+5j13R6iBnnKAR7bhNrOV8Cc0RfpSzjs=" // Test

const DefiantLabPubKey = "Aw22yXnDmYKzQ1CeHh6A+PD1043vsbSBH5FmuAWIlkS7" // Prod
// const DefiantLabPubKey = "A8gI+6AHMv9Tg37JyrxSP16hUH76Umr4krXfIEqOQJMo" // Test

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

func CreateNewAccount(ctx sdk.Context, authKeeper keeper.AccountKeeper, address sdk.AccAddress, secpPubKey []byte) error {
	// Check if the account already exists
	account := authKeeper.GetAccount(ctx, address)
	if account != nil {
		return fmt.Errorf("account %s already exists", address.String())
	}

	// Optionally, set any initial values or parameters for the new account
	// For example, you might want to set an initial balance using the bank module
	var pubkey cryptotypes.PubKey = &secp256k1.PubKey{Key: secpPubKey}

	// Create a new account with the address and public key
	newAccount := authKeeper.NewAccountWithAddress(ctx, address)
	if newAccount == nil {
		return fmt.Errorf("failed to create new account for address %s", address.String())
	}

	newAccount.SetPubKey(pubkey)
	log.Printf("New account created %v: %v", address.String(), pubkey.String())

	// Save the new account to the state
	authKeeper.SetAccount(ctx, newAccount)
	return nil
}

func getDelegations(
	ctx sdk.Context,
	stakingKeeper stakingkeeper.Keeper,
	delegatorAddr sdk.AccAddress,
) []stakingtypes.Delegation {
	delegations := stakingKeeper.GetAllDelegatorDelegations(ctx, delegatorAddr)
	return delegations
}

func InitializeValidatorSigningInfo(ctx sdk.Context, slashingKeeper slashingkeeper.Keeper, stakingKeeper stakingkeeper.Keeper, consAddr sdk.ConsAddress) error {
	// Check if signing info already exists to avoid overwriting it
	_, found := slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		startHeight := ctx.BlockHeight()
		signingInfo := slashingtypes.NewValidatorSigningInfo(consAddr, startHeight, 0, time.Unix(0, 0), false, 0)
		slashingKeeper.SetValidatorSigningInfo(ctx, consAddr, signingInfo)
	}
	return nil
}

func InitializeValidatorDistributionInfo(ctx sdk.Context, keepers *keepers.AppKeepers, validatorAddr sdk.ValAddress) {
	// Initialize distribution information for the validator
	// set initial historical rewards (period 0) with reference count of 1
	keepers.DistrKeeper.SetValidatorHistoricalRewards(ctx, validatorAddr, 0, distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1))

	// set current rewards (starting at period 1)
	keepers.DistrKeeper.SetValidatorCurrentRewards(ctx, validatorAddr, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 1))

	// set accumulated commission
	keepers.DistrKeeper.SetValidatorAccumulatedCommission(ctx, validatorAddr, distrtypes.InitialValidatorAccumulatedCommission())

	// set outstanding rewards
	keepers.DistrKeeper.SetValidatorOutstandingRewards(ctx, validatorAddr, distrtypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}})
}

func createValidator(ctx sdk.Context, keeppers *keepers.AppKeepers, address string, pubKey cryptotypes.PubKey, description stakingtypes.Description, comission stakingtypes.Commission) (stakingtypes.Validator, error) {

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
	keeppers.StakingKeeper.SetValidator(ctx, validator)

	consAddr := sdk.ConsAddress(pubKey.Address())
	valconsAddr, err := validator.GetConsAddr()
	if err != nil {
		log.Printf("Error when converting validator consensus address to string: %s", err)
		return stakingtypes.Validator{}, err
	}

	log.Printf("Created validator %v (%v:%v)", valAddr.String(), consAddr.String(), valconsAddr)

	keeppers.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
	InitializeValidatorSigningInfo(ctx, keeppers.SlashingKeeper, *keeppers.StakingKeeper, consAddr)
	InitializeValidatorDistributionInfo(ctx, keeppers, valAddr)

	err = keeppers.StakingKeeper.Hooks().AfterValidatorCreated(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}

	err = keeppers.DistrKeeper.Hooks().AfterValidatorCreated(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}

	return validator, nil
}

func withdrawRewardsAndCommission(ctx sdk.Context, sk stakingkeeper.Keeper, dk distributionkeeper.Keeper, oldValAddress sdk.ValAddress, newValAddress sdk.ValAddress) {
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

func selfDelegate(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, bankKeeper bankkeeper.Keeper, delegatorAddr sdk.AccAddress, validator stakingtypes.Validator, amount sdk.Coin) error {
	// Ensure the delegator (validator account) has enough balance for the delegation
	if !bankKeeper.HasBalance(ctx, delegatorAddr, amount) {
		return sdkerrors.Wrapf(errortypes.ErrInsufficientFunds, "not enough balance to self-delegate to validator: %s", validator.OperatorAddress)
	}

	// Send coins from the delegator's account to the module account (staking module account) as part of delegation
	err := bankKeeper.SendCoinsFromAccountToModule(ctx, delegatorAddr, stakingtypes.NotBondedPoolName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	// Delegate tokens to the validator
	_, err = stakingKeeper.Delegate(ctx, delegatorAddr, amount.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return err
	}
	return nil
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

func moveValidatorDelegations(ctx sdk.Context, k stakingkeeper.Keeper, d distributionkeeper.Keeper, b bankkeeper.Keeper, oldVal stakingtypes.Validator, newVal stakingtypes.Validator, selfDelegationTokens math.Int) error {
	// Retrieving self-delegation address
	validatorDelegatorAddr := sdk.AccAddress(oldVal.GetOperator())
	newValidatorDelegatorAddr := sdk.AccAddress(newVal.GetOperator())

	log.Printf("Validator withdraw address: %v", validatorDelegatorAddr)

	minSelfDelegation := math.NewInt(oldVal.MinSelfDelegation.Int64())
	minDelegationShares, err := oldVal.SharesFromTokens(minSelfDelegation)

	if err != nil {
		log.Printf("Error when converting min self-delegation to shares: %s", err)
		return err
	}

	totalSharesToMove := oldVal.DelegatorShares.Sub(minDelegationShares)
	tokensToMove := oldVal.TokensFromShares(totalSharesToMove)

	_, found := k.GetDelegation(ctx, validatorDelegatorAddr, oldVal.GetOperator())
	if !found {
		log.Printf("%s self delegation not found, self delegating 100 Odin", validatorDelegatorAddr)
		selfDelegate(ctx, k, b, newValidatorDelegatorAddr, newVal, sdk.NewCoin("loki", selfDelegationTokens))
	}

	k.RemoveValidatorTokensAndShares(ctx, oldVal, totalSharesToMove)
	k.AddValidatorTokensAndShares(ctx, newVal, tokensToMove.TruncateInt().Add(selfDelegationTokens)) // Adding new self-delegation

	for _, delegation := range k.GetValidatorDelegations(ctx, oldVal.GetOperator()) {
		log.Printf("Moving validator delegation from %v to %v", delegation.DelegatorAddress, newVal.OperatorAddress)

		withdrawAddress := d.GetDelegatorWithdrawAddr(ctx, delegation.GetDelegatorAddr())
		log.Printf("Delegator withdraw address: %v", validatorDelegatorAddr)

		var newDelegationAmt math.LegacyDec
		var newDelegatorAddress string

		// Processing self-delergation, keeping old validator's self-delegation min amount ot make it survive the upgrade
		log.Printf("Withdraw address and validator delegator %s VS %s", withdrawAddress.String(), validatorDelegatorAddr.String())
		if withdrawAddress.String() == validatorDelegatorAddr.String() {
			newDelegationAmt = delegation.Shares.Sub(minDelegationShares)

			// Create a new delegation to the new validator
			oldDelegationReplacement := stakingtypes.Delegation{
				DelegatorAddress: delegation.DelegatorAddress,
				ValidatorAddress: newVal.OperatorAddress,
				Shares:           minDelegationShares,
			}

			newDelegatorAddress = newValidatorDelegatorAddr.String()

			log.Printf("New delegator address for self delegation: %s", newDelegatorAddress)
			log.Printf("Old delegation amount: %s", minDelegationShares)

			err := k.Hooks().BeforeDelegationCreated(ctx, delegation.GetDelegatorAddr(), oldVal.GetOperator())
			if err != nil {
				log.Printf("Error when running hook after adding delegation %v to %v", delegation.GetDelegatorAddr(), oldVal.GetOperator())
				return err
			}

			// Creating old validator's new self-delegation
			k.SetDelegation(ctx, oldDelegationReplacement)

			err = d.Hooks().AfterDelegationModified(ctx, delegation.GetDelegatorAddr(), oldVal.GetOperator())
			if err != nil {
				log.Printf("Error when running hook after adding delegation %v to %v", delegation.GetDelegatorAddr(), oldVal.GetOperator())
				return err
			}

		} else {
			newDelegationAmt = delegation.Shares
			newDelegatorAddress = delegation.DelegatorAddress
		}

		// Remove the delegation to the old validator
		err = d.Hooks().BeforeDelegationRemoved(ctx, delegation.GetDelegatorAddr(), newVal.GetOperator())
		if err != nil {
			log.Printf("Error when running hook before adding delegation %v to %v", delegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}

		k.RemoveDelegation(ctx, delegation)

		newDelegation := stakingtypes.Delegation{
			DelegatorAddress: newDelegatorAddress,
			ValidatorAddress: newVal.OperatorAddress,
			Shares:           newDelegationAmt,
		}

		err := k.Hooks().BeforeDelegationCreated(ctx, newDelegation.GetDelegatorAddr(), newVal.GetOperator())
		if err != nil {
			log.Printf("Error when running hook before adding delegation %v to %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}

		err = d.Hooks().BeforeDelegationCreated(ctx, newDelegation.GetDelegatorAddr(), newVal.GetOperator())
		if err != nil {
			log.Printf("Error when running hook before adding delegation %v to %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}
		log.Printf("New delegation amount: %s", newDelegationAmt)

		k.SetDelegation(ctx, newDelegation)

		err = k.Hooks().AfterDelegationModified(ctx, newDelegation.GetDelegatorAddr(), newVal.GetOperator())
		if err != nil {
			log.Printf("Error when running hook after adding delegation %v to %v. %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator(), err)
			return err
		}
		err = d.Hooks().AfterDelegationModified(ctx, newDelegation.GetDelegatorAddr(), newVal.GetOperator())
		if err != nil {
			log.Printf("Error when running hook after adding delegation %v to %v: %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator(), err)
			return err
		}

		// Double checking if the delegation has been added
		found := d.HasDelegatorStartingInfo(ctx, newVal.GetOperator(), newDelegation.GetDelegatorAddr())
		if !found {
			return fmt.Errorf("delegator starting info not found")
		}
	}

	return nil
}

func moveDelegations(ctx sdk.Context, keepers *keepers.AppKeepers, oldAddress sdk.AccAddress, newVal stakingtypes.Validator) error {
	for _, delegation := range getDelegations(ctx, *keepers.StakingKeeper, oldAddress) {
		log.Printf("Moving delegation from %v to %v", delegation.DelegatorAddress, newVal.OperatorAddress)
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

		err = keepers.DistrKeeper.Hooks().BeforeDelegationCreated(ctx, newDelegation.GetDelegatorAddr(), newVal.GetOperator())
		if err != nil {
			log.Printf("Error when running hook before adding delegation %v to %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}

		keepers.StakingKeeper.SetDelegation(ctx, newDelegation)
		err = keepers.DistrKeeper.Hooks().AfterDelegationModified(ctx, newDelegation.GetDelegatorAddr(), newVal.GetOperator())
		if err != nil {
			log.Printf("Error when running hook after addig delegation %v to %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}
	}
	return nil
}

func moveSelfDelegation(ctx sdk.Context, keepers *keepers.AppKeepers, oldDelegatorAddress sdk.AccAddress, newDelegatorAddress sdk.AccAddress, validatorAddr sdk.ValAddress) error {
	stakingKeeper := keepers.StakingKeeper

	// Get the delegation from the old validator
	delegation, found := stakingKeeper.GetDelegation(ctx, oldDelegatorAddress, validatorAddr)
	if !found {
		log.Printf("self delegation not found: %s", oldDelegatorAddress)
		return fmt.Errorf("self delegation not found")
	}
	amount := delegation.Shares

	stakingKeeper.RemoveDelegation(ctx, delegation)

	// Create a new delegation to the new validator
	newDelegation := stakingtypes.Delegation{
		DelegatorAddress: newDelegatorAddress.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           amount,
	}

	err := stakingKeeper.Hooks().BeforeDelegationCreated(ctx, delegation.GetDelegatorAddr(), validatorAddr)
	if err != nil {
		log.Printf("Error when running hook before adding delegation %v to %v", delegation.GetDelegatorAddr(), validatorAddr)
		return err
	}
	// Save the new delegation
	stakingKeeper.SetDelegation(ctx, newDelegation)
	err = keepers.DistrKeeper.Hooks().AfterDelegationModified(ctx, newDelegatorAddress, validatorAddr)
	if err != nil {
		log.Printf("Error when running hook before adding delegation %v to %v", delegation.GetDelegatorAddr(), validatorAddr)
		return err
	}
	return nil
}

func fixUnbondingHeight(ctx sdk.Context, keepers *keepers.AppKeepers, validator stakingtypes.Validator) error {
	validator.UnbondingHeight = 0
	validator.UnbondingTime = time.Unix(0, 0)
	validator.Status = stakingtypes.Bonded
	keepers.StakingKeeper.SetValidator(ctx, validator)
	return nil
}

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
		log.Printf("Could not send coins from: %s, to: %s, error: %s", fromAddr, toAddr, err)
		return err
	}
	return nil
}

func fixDefiantLabs(ctx sdk.Context, keepers *keepers.AppKeepers) error {

	// Fixing self delegation
	DefiantLabsValAddress, err := addrToValAddr(DefiantLabAccAddress)
	if err != nil {
		return err
	}

	DanVal, found := keepers.StakingKeeper.GetValidator(ctx, DefiantLabsValAddress)
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

	PubKeyBytes, err := base64.StdEncoding.DecodeString(DefiantLabPubKey)
	if err != nil {
		log.Printf("Error whend decoding public key from string %v", err)
		return err
	}

	err = CreateNewAccount(ctx, keepers.AccountKeeper, DefiantLabsAcc, PubKeyBytes)
	if err != nil {
		log.Printf("Error when creating new account for %v: %s", DefiantLabsAcc, err)
	}

	keepers.DistrKeeper.SetWithdrawAddr(ctx, DefiantLabsAcc, DefiantLabsAcc)

	// sending balances
	ctx.Logger().Info(fmt.Sprintf("Sending tokens from %s to %s", DefiantLabOldAccAddress, DefiantLabAccAddress))

	balance, err := getBalance(ctx, *keepers.StakingKeeper, keepers.AccountKeeper, keepers.BankKeeper, DefiantLabsOldAcc)
	if err != nil {
		log.Printf("Error when retrieving balance for address %s: %s", DefiantLabOldAccAddress, err)
		return err
	}
	sendCoins(ctx, keepers.BankKeeper, DefiantLabsOldAcc, DefiantLabsAcc, balance)

	// Moving delegations
	moveSelfDelegation(ctx, keepers, DefiantLabsOldAcc, DefiantLabsAcc, DefiantLabsValAddress)
	moveDelegations(ctx, keepers, DefiantLabsOldAcc, DanVal)

	err = fixUnbondingHeight(ctx, keepers, DanVal)
	if err != nil {
		log.Printf("Error when updating unbonding height %s: %s", DefiantLabOldAccAddress, err)
		return err
	}
	return nil
}

func fixMainnet3(ctx sdk.Context, keepers *keepers.AppKeepers) error {

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

	Odin3ValAddr, err := addrToValAddr(OdinMainnet3NewAccAddress)
	if err != nil {
		return err
	}

	// rewards and comission
	withdrawRewardsAndCommission(ctx, *keepers.StakingKeeper, keepers.DistrKeeper, Odin3OldValAddress, Odin3ValAddr)
	ctx.Logger().Info(fmt.Sprintf("Sending tokens from %s to %s", OldMainnet3Addr, NewMainnet3Addr))
	balance, err := getBalance(ctx, *keepers.StakingKeeper, keepers.AccountKeeper, keepers.BankKeeper, OldMainnet3Addr)
	if err != nil {
		log.Printf("Error when retrieving balance for address %s: %s", OldMainnet3Addr, err)
		return err
	}

	// sending balances
	sendCoins(ctx, keepers.BankKeeper, OldMainnet3Addr, NewMainnet3Addr, balance)

	Odin3PubKeyBytes, err := base64.StdEncoding.DecodeString(OdinMainnet3ValPubKey)
	if err != nil {
		log.Printf("Error whend decoding public key from string %v", err)
		return err
	}

	Odin3PubKey := ed25519.PubKey{
		Key: Odin3PubKeyBytes,
	}

	Odin3Val, err := createValidator(ctx, keepers, string(Odin3ValAddr), &Odin3PubKey, Odin3OldVal.Description, Odin3OldVal.Commission)
	if err != nil {
		return err
	}

	ctx.Logger().Info(fmt.Sprintf("Moving validator delegations from %s to %s", Odin3OldValAddress, Odin3ValAddr))

	err = moveValidatorDelegations(ctx, *keepers.StakingKeeper, keepers.DistrKeeper, keepers.BankKeeper, Odin3OldVal, Odin3Val, math.NewInt(100000000))
	if err != nil {
		return err
	}

	// moveSelfDelegation(ctx, keepers, OldMainnet3Addr, NewMainnet3Addr, Odin3ValAddr)
	moveDelegations(ctx, keepers, sdk.AccAddress(OdinMainnet3OldAccAddress), Odin3Val)

	Odin3Val.UpdateStatus(stakingtypes.Bonded)

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
