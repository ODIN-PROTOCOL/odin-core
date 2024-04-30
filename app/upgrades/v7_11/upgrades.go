package v7_11

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
	"unsafe"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	sdkerrors "cosmossdk.io/errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ODIN-PROTOCOL/odin-core/app/keepers"
	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
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
const OdinMainnet3ValPubKey = "f7pqqa+1Rkl+5j13R6iBnnKAR7bhNrOV8Cc0RfpSzjs=" // Prod
// const OdinMainnet3ValPubKey = "f7pqqa+1Rkl+5j13R6iBnnKAR7bhNrOV8Cc0RfpSzjs=" // Test

const DefiantLabPubKey = "Aw22yXnDmYKzQ1CeHh6A+PD1043vsbSBH5FmuAWIlkS7" // Prod
// const DefiantLabPubKey = "A8gI+6AHMv9Tg37JyrxSP16hUH76Umr4krXfIEqOQJMo" // Test

func deletePacketCommitment(ctx sdk.Context, storeKey storetypes.StoreKey, portID, channelID string, sequence uint64) {
	//storeKey := sdk.NewKVStoreKey(ibcexported.StoreKey)
	store := ctx.KVStore(storeKey)
	store.Delete(host.PacketCommitmentKey(portID, channelID, sequence))
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func FlushIBCPackets(ctx sdk.Context, keepers *keepers.AppKeepers) {
	// Get the IBC module's keeper
	ibcKeeper := keepers.IBCKeeper
	cdc := ibcKeeper.Codec()

	transferModule, found := ibcKeeper.PortKeeper.Router.GetRoute("transfer")
	if !found {
		log.Printf("transfer module not found")
		return
	}

	ibcStoreKey := GetUnexportedField(reflect.ValueOf(&ibcKeeper.ChannelKeeper).Elem().FieldByName("storeKey"))

	stuckedPackets := GetStuckedPackets()

	for _, packetCommitment := range ibcKeeper.ChannelKeeper.GetAllPacketCommitmentsAtChannel(ctx, "transfer", "channel-3") { //.GetAllPacketReceipts(ctx) {
		log.Printf("Packet commitment for channel %v, %v, %v", packetCommitment.ChannelId, packetCommitment.PortId, packetCommitment.Sequence)
		if stuckedPacket, ok := stuckedPackets[packetCommitment.Sequence]; ok {
			packetData := types.NewFungibleTokenPacketData(
				stuckedPacket.Denom, stuckedPacket.Amount, stuckedPacket.Sender, stuckedPacket.Receiver, stuckedPacket.Memo,
			)

			timeoutHeight := clienttypes.Height{
				RevisionNumber: stuckedPacket.RevisionNumber,
				RevisionHeight: stuckedPacket.RevisionHeight,
			}

			//commitment := ibcKeeper.ChannelKeeper.GetPacketCommitment(ctx, packetCommitment.PortId, packetCommitment.ChannelId, packetCommitment.Sequence)
			packet := channeltypes.NewPacket(packetData.GetBytes(), packetCommitment.Sequence, packetCommitment.PortId, packetCommitment.ChannelId,
				packetCommitment.PortId, "channel-258", timeoutHeight, stuckedPacket.TimeoutTimestamp)
			commitment := channeltypes.CommitPacket(cdc, packet)
			if bytes.Equal(commitment, packetCommitment.Data) {
				log.Printf("refund packet with sequence %d", packetCommitment.Sequence)
				err := transferModule.OnTimeoutPacket(ctx, packet, sdk.AccAddress{})
				if err != nil {
					log.Printf("cannot timeout packet, sequence: %d", packetCommitment.Sequence)
					panic(err)
				}
			}
		}

		deletePacketCommitment(ctx, ibcStoreKey.(storetypes.StoreKey), packetCommitment.PortId, packetCommitment.ChannelId, packetCommitment.Sequence)
	}

	ibcKeeper.ChannelKeeper.SetNextSequenceSend(ctx, "transfer", "channel-3", 6404)
}

func getBalance(
	ctx context.Context,
	ak keeper.AccountKeeper,
	bk bankkeeper.Keeper,
	addr sdk.AccAddress,
) (sdk.Coins, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)

	// Get all delegator delegations for address
	account := ak.GetAccount(ctx, addr)
	vestingAccount, ok := account.(*vestingtypes.BaseVestingAccount)
	if !ok {
		return bk.GetAllBalances(ctx, addr), nil
	} else {
		//If the account is a vesting account, create a copy of the account
		//and vest all coins with the current block header time
		newVestingAcc := vestingtypes.NewContinuousVestingAccountRaw(vestingAccount, goCtx.BlockHeader().Time.Unix())
		ak.SetAccount(ctx, newVestingAcc)
		return newVestingAcc.GetVestedCoins(goCtx.BlockTime()), nil
	}
}

func CreateNewAccount(ctx context.Context, authKeeper keeper.AccountKeeper, address sdk.AccAddress, secpPubKey []byte) error {
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

	err := newAccount.SetPubKey(pubkey)
	if err != nil {
		return err
	}
	log.Printf("New account created %v: %v", address.String(), pubkey.String())

	// Save the new account to the state
	authKeeper.SetAccount(ctx, newAccount)
	return nil
}

func getDelegations(
	ctx context.Context,
	stakingKeeper stakingkeeper.Keeper,
	delegatorAddr sdk.AccAddress,
) ([]stakingtypes.Delegation, error) {
	return stakingKeeper.GetAllDelegatorDelegations(ctx, delegatorAddr)
}

func InitializeValidatorSigningInfo(ctx context.Context, slashingKeeper slashingkeeper.Keeper, consAddr sdk.ConsAddress) error {
	goCtx := sdk.UnwrapSDKContext(ctx)
	// Check if signing info already exists to avoid overwriting it
	_, err := slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	if err == slashingtypes.ErrNoSigningInfoFound {
		startHeight := goCtx.BlockHeight()
		signingInfo := slashingtypes.NewValidatorSigningInfo(consAddr, startHeight, 0, time.Unix(0, 0), false, 0)
		err := slashingKeeper.SetValidatorSigningInfo(ctx, consAddr, signingInfo)
		if err != nil {
			return err
		}

		return nil
	}

	return err
}

func InitializeValidatorDistributionInfo(ctx context.Context, keepers *keepers.AppKeepers, validatorAddr sdk.ValAddress) error {
	// Initialize distribution information for the validator
	// set initial historical rewards (period 0) with reference count of 1
	err := keepers.DistrKeeper.SetValidatorHistoricalRewards(ctx, validatorAddr, 0, distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1))
	if err != nil {
		return err
	}

	// set current rewards (starting at period 1)
	err = keepers.DistrKeeper.SetValidatorCurrentRewards(ctx, validatorAddr, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 1))
	if err != nil {
		return err
	}

	// set accumulated commission
	err = keepers.DistrKeeper.SetValidatorAccumulatedCommission(ctx, validatorAddr, distrtypes.InitialValidatorAccumulatedCommission())
	if err != nil {
		return err
	}

	// set outstanding rewards
	return keepers.DistrKeeper.SetValidatorOutstandingRewards(ctx, validatorAddr, distrtypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}})
}

func createValidator(ctx context.Context, keeppers *keepers.AppKeepers, address string, pubKey cryptotypes.PubKey, description stakingtypes.Description, comission stakingtypes.Commission) (stakingtypes.Validator, error) {

	valAddr := sdk.ValAddress(address)
	minSelfDelegation := math.OneInt()

	// Create the validator
	validator, err := stakingtypes.NewValidator(address, pubKey, description)
	if err != nil {
		log.Printf("Error when creating a validator %v: %s", valAddr, err)
		return stakingtypes.Validator{}, err
	}

	validator.MinSelfDelegation = minSelfDelegation
	validator.Status = stakingtypes.Bonded
	validator.Tokens = math.ZeroInt()
	validator.DelegatorShares = math.LegacyZeroDec()
	validator.Commission = comission

	// Update validators in the store
	err = keeppers.StakingKeeper.SetValidator(ctx, validator)
	if err != nil {
		return stakingtypes.Validator{}, err
	}

	consAddr := sdk.ConsAddress(pubKey.Address())
	valconsAddr, err := validator.GetConsAddr()
	if err != nil {
		log.Printf("Error when converting validator consensus address to string: %s", err)
		return stakingtypes.Validator{}, err
	}

	log.Printf("Created validator %v (%v:%v)", valAddr.String(), consAddr.String(), valconsAddr)

	err = keeppers.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
	if err != nil {
		return stakingtypes.Validator{}, err
	}

	err = InitializeValidatorSigningInfo(ctx, keeppers.SlashingKeeper, consAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}

	err = InitializeValidatorDistributionInfo(ctx, keeppers, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}

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

func withdrawRewardsAndCommission(ctx context.Context, ak authkeeper.AccountKeeper, sk stakingkeeper.Keeper, dk distributionkeeper.Keeper, oldValAddress sdk.ValAddress, newValAddress sdk.ValAddress) error {
	oldValAccAddress := sdk.AccAddress(oldValAddress)
	newValAccAddress := sdk.AccAddress(newValAddress)

	delegations, err := sk.GetValidatorDelegations(ctx, oldValAddress)
	if err != nil {
		return err
	}

	// withdrawing all rewards, self-delegation rewards mapped to new account
	for _, delegation := range delegations {
		withdrawAddress, err := dk.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(delegation.DelegatorAddress))
		if err != nil {
			return err
		}

		delegatorAccAddress, err := ak.AddressCodec().StringToBytes(delegation.GetDelegatorAddr())
		if err != nil {
			return err
		}

		// we suppose that old Odin accounts are unavailable, so we're routing rewards to new addresses and proceeding wit hwithdraws
		if withdrawAddress.String() == oldValAccAddress.String() {
			log.Printf("Found delegation which withdrawal address is the old one: %v. Setting withdrawal address to new account: %v", oldValAccAddress.String(), newValAccAddress.String())
			err = dk.SetDelegatorWithdrawAddr(ctx, delegatorAccAddress, newValAccAddress)
			if err != nil {
				return err
			}
		}

		log.Printf("Withdrawing reward for %v delegator address from %v", string(delegatorAccAddress), oldValAddress.String())
		_, err = dk.WithdrawDelegationRewards(ctx, delegatorAccAddress, oldValAddress)
		if err != nil {
			return err
		}
	}

	// Comission
	// explicitly setting validator withdrawal address, in case it has no self-delegation in the loop above
	err = dk.SetDelegatorWithdrawAddr(ctx, oldValAccAddress, newValAccAddress)
	if err != nil {
		return err
	}

	_, err = dk.WithdrawValidatorCommission(ctx, oldValAddress)
	if err != nil {
		return err
	}

	return nil
}

func selfDelegate(ctx context.Context, stakingKeeper stakingkeeper.Keeper, bankKeeper bankkeeper.Keeper, delegatorAddr sdk.AccAddress, validator stakingtypes.Validator, amount sdk.Coin) error {
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

func moveValidatorDelegations(ctx context.Context, ak authkeeper.AccountKeeper, k stakingkeeper.Keeper, d distributionkeeper.Keeper, b bankkeeper.Keeper, oldVal stakingtypes.Validator, newVal stakingtypes.Validator, selfDelegationTokens math.Int) error {
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

	oldValBz, err := k.ValidatorAddressCodec().StringToBytes(oldVal.GetOperator())
	if err != nil {
		return err
	}

	newValBz, err := k.ValidatorAddressCodec().StringToBytes(newVal.GetOperator())
	if err != nil {
		return err
	}

	selfDelegation, err := k.GetDelegation(ctx, validatorDelegatorAddr, oldValBz)
	if err == stakingtypes.ErrNoDelegation {
		log.Printf("%s self delegation not found, self delegating 100 Odin", validatorDelegatorAddr)
		err = selfDelegate(ctx, k, b, newValidatorDelegatorAddr, newVal, sdk.NewCoin("loki", selfDelegationTokens))
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		log.Printf("Self delegation found %s, %s, %s", selfDelegation.DelegatorAddress, selfDelegation.Shares.String(), selfDelegation.ValidatorAddress)
	}

	_, _, err = k.RemoveValidatorTokensAndShares(ctx, oldVal, totalSharesToMove)
	if err != nil {
		return err
	}

	_, _, err = k.AddValidatorTokensAndShares(ctx, newVal, tokensToMove.TruncateInt().Add(selfDelegationTokens)) // Adding new self-delegation
	if err != nil {
		return err
	}

	delegations, err := k.GetValidatorDelegations(ctx, oldValBz)
	if err != nil {
		return err
	}

	for _, delegation := range delegations {
		delegatorAccAddress, err := ak.AddressCodec().StringToBytes(delegation.GetDelegatorAddr())
		if err != nil {
			return err
		}

		log.Printf("Moving validator delegation from %v to %v", delegation.DelegatorAddress, newVal.OperatorAddress)

		withdrawAddress, err := d.GetDelegatorWithdrawAddr(ctx, delegatorAccAddress)
		if err != nil {
			return err
		}

		log.Printf("Delegator withdraw address: %v", validatorDelegatorAddr)

		var newDelegationAmt math.LegacyDec
		var newDelegatorAddress string

		// Processing self-delergation, keeping old validator's self-delegation min amount ot make it survive the upgrade
		log.Printf("Withdraw address and validator delegator %s VS %s", withdrawAddress.String(), validatorDelegatorAddr.String())

		if delegation.DelegatorAddress == selfDelegation.DelegatorAddress && delegation.ValidatorAddress == selfDelegation.ValidatorAddress {
			log.Printf("Processing self delegation")

			newDelegationAmt = delegation.Shares.Sub(minDelegationShares)

			// Create a new delegation to the new validator
			oldDelegationReplacement := stakingtypes.Delegation{
				DelegatorAddress: delegation.DelegatorAddress,
				ValidatorAddress: oldVal.OperatorAddress,
				Shares:           minDelegationShares,
			}

			newDelegatorAddress = newValidatorDelegatorAddr.String()

			log.Printf("New delegator address for self delegation: %s", newDelegatorAddress)
			log.Printf("Old delegation amount: %s", minDelegationShares)

			err := k.Hooks().BeforeDelegationCreated(ctx, delegatorAccAddress, oldValBz)
			if err != nil {
				log.Printf("Error when running hook after adding delegation %v to %v", delegation.GetDelegatorAddr(), oldVal.GetOperator())
				return err
			}

			// Creating old validator's new self-delegation
			err = k.SetDelegation(ctx, oldDelegationReplacement)
			if err != nil {
				return err
			}

			err = k.Hooks().AfterDelegationModified(ctx, delegatorAccAddress, oldValBz)
			if err != nil {
				log.Printf("Error when running hook after adding delegation %v to %v", delegation.GetDelegatorAddr(), newVal.GetOperator())
				return err
			}
			err = d.Hooks().AfterDelegationModified(ctx, delegatorAccAddress, oldValBz)
			if err != nil {
				log.Printf("Error when running hook after adding delegation %v to %v", delegation.GetDelegatorAddr(), newVal.GetOperator())
				return err
			}

		} else {
			newDelegationAmt = delegation.Shares
			newDelegatorAddress = delegation.DelegatorAddress
		}

		// Remove the delegation to the old validator
		err = d.Hooks().BeforeDelegationRemoved(ctx, delegatorAccAddress, newValBz)
		if err != nil {
			log.Printf("Error when running hook before adding delegation %v to %v", delegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}

		err = k.RemoveDelegation(ctx, delegation)
		if err != nil {
			return err
		}

		newDelegation := stakingtypes.Delegation{
			DelegatorAddress: newDelegatorAddress,
			ValidatorAddress: newVal.OperatorAddress,
			Shares:           newDelegationAmt,
		}

		newDelegatorAccAddress, err := ak.AddressCodec().StringToBytes(newDelegation.GetDelegatorAddr())
		if err != nil {
			return err
		}

		err = k.Hooks().BeforeDelegationCreated(ctx, newDelegatorAccAddress, newValBz)
		if err != nil {
			log.Printf("Error when running hook before adding delegation %v to %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}

		err = d.Hooks().BeforeDelegationCreated(ctx, newDelegatorAccAddress, newValBz)
		if err != nil {
			log.Printf("Error when running hook before adding delegation %v to %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}
		log.Printf("New delegation amount: %s", newDelegationAmt)

		err = k.SetDelegation(ctx, newDelegation)
		if err != nil {
			return err
		}

		err = k.Hooks().AfterDelegationModified(ctx, newDelegatorAccAddress, newValBz)
		if err != nil {
			log.Printf("Error when running hook after adding delegation %v to %v. %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator(), err)
			return err
		}
		err = d.Hooks().AfterDelegationModified(ctx, newDelegatorAccAddress, newValBz)
		if err != nil {
			log.Printf("Error when running hook after adding delegation %v to %v: %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator(), err)
			return err
		}

		// Double checking if the delegation has been added
		found, err := d.HasDelegatorStartingInfo(ctx, newValBz, newDelegatorAccAddress)
		if err != nil {
			return err
		}

		if !found {
			return fmt.Errorf("delegator starting info not found")
		}
	}

	return nil
}

func moveDelegations(ctx context.Context, keepers *keepers.AppKeepers, oldAddress sdk.AccAddress, newVal stakingtypes.Validator) error {
	newValBz, err := keepers.StakingKeeper.ValidatorAddressCodec().StringToBytes(newVal.GetOperator())
	if err != nil {
		return err
	}

	delegations, err := getDelegations(ctx, *keepers.StakingKeeper, oldAddress)
	if err != nil {
		return err
	}

	for _, delegation := range delegations {
		log.Printf("Moving delegation from %v to %v", delegation.DelegatorAddress, newVal.OperatorAddress)
		err := keepers.StakingKeeper.RemoveDelegation(ctx, delegation)
		if err != nil {
			return err
		}

		newDelegation := stakingtypes.Delegation{
			DelegatorAddress: delegation.DelegatorAddress,
			ValidatorAddress: newVal.OperatorAddress,
			Shares:           delegation.Shares,
		}

		newDelegatorAccAddress, err := keepers.AccountKeeper.AddressCodec().StringToBytes(newDelegation.GetDelegatorAddr())
		if err != nil {
			return err
		}

		err = keepers.StakingKeeper.Hooks().BeforeDelegationCreated(ctx, newDelegatorAccAddress, newValBz)
		if err != nil {
			log.Printf("Error when running hook after adding delegation %v to %v", delegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}

		err = keepers.DistrKeeper.Hooks().BeforeDelegationCreated(ctx, newDelegatorAccAddress, newValBz)
		if err != nil {
			log.Printf("Error when running hook before adding delegation %v to %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}

		err = keepers.StakingKeeper.SetDelegation(ctx, newDelegation)
		if err != nil {
			return err
		}

		err = keepers.DistrKeeper.Hooks().AfterDelegationModified(ctx, newDelegatorAccAddress, newValBz)
		if err != nil {
			log.Printf("Error when running hook after addig delegation %v to %v", newDelegation.GetDelegatorAddr(), newVal.GetOperator())
			return err
		}
	}
	return nil
}

func moveSelfDelegation(ctx context.Context, keepers *keepers.AppKeepers, oldDelegatorAddress sdk.AccAddress, newDelegatorAddress sdk.AccAddress, validatorAddr sdk.ValAddress) error {
	stakingKeeper := keepers.StakingKeeper

	// Get the delegation from the old validator
	delegation, err := stakingKeeper.GetDelegation(ctx, oldDelegatorAddress, validatorAddr)
	if err == stakingtypes.ErrNoDelegation {
		log.Printf("self delegation not found: %s", oldDelegatorAddress)
		return fmt.Errorf("self delegation not found")
	}
	if err != nil {
		return err
	}

	amount := delegation.Shares

	err = stakingKeeper.RemoveDelegation(ctx, delegation)
	if err != nil {
		return err
	}

	// Create a new delegation to the new validator
	newDelegation := stakingtypes.Delegation{
		DelegatorAddress: newDelegatorAddress.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           amount,
	}

	delegatorAccAddress, err := keepers.AccountKeeper.AddressCodec().StringToBytes(delegation.GetDelegatorAddr())
	if err != nil {
		return err
	}

	err = stakingKeeper.Hooks().BeforeDelegationCreated(ctx, delegatorAccAddress, validatorAddr)
	if err != nil {
		log.Printf("Error when running hook before adding delegation %v to %v", delegation.GetDelegatorAddr(), validatorAddr)
		return err
	}
	// Save the new delegation
	err = stakingKeeper.SetDelegation(ctx, newDelegation)
	if err != nil {
		return err
	}

	err = keepers.DistrKeeper.Hooks().AfterDelegationModified(ctx, newDelegatorAddress, validatorAddr)
	if err != nil {
		log.Printf("Error when running hook before adding delegation %v to %v", delegation.GetDelegatorAddr(), validatorAddr)
		return err
	}
	return nil
}

func fixUnbondingHeight(ctx context.Context, keepers *keepers.AppKeepers, validator stakingtypes.Validator) error {
	validator.UnbondingHeight = 0
	validator.UnbondingTime = time.Unix(0, 0)
	validator.Status = stakingtypes.Bonded
	return keepers.StakingKeeper.SetValidator(ctx, validator)
}

func sendCoins(
	ctx context.Context,
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

func fixDefiantLabs(ctx context.Context, keepers *keepers.AppKeepers) error {
	goCtx := sdk.UnwrapSDKContext(ctx)

	// Fixing self delegation
	DefiantLabsValAddress, err := addrToValAddr(DefiantLabAccAddress)
	if err != nil {
		return err
	}

	DanVal, err := keepers.StakingKeeper.GetValidator(ctx, DefiantLabsValAddress)
	if err == stakingtypes.ErrNoValidatorFound {
		log.Printf("Validator with %v has not been found", DefiantLabsValAddress)
		return err
	}
	if err != nil {
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

	err = keepers.DistrKeeper.SetWithdrawAddr(ctx, DefiantLabsAcc, DefiantLabsAcc)
	if err != nil {
		return err
	}

	// sending balances
	goCtx.Logger().Info(fmt.Sprintf("Sending tokens from %s to %s", DefiantLabOldAccAddress, DefiantLabAccAddress))

	balance, err := getBalance(ctx, keepers.AccountKeeper, keepers.BankKeeper, DefiantLabsOldAcc)
	if err != nil {
		log.Printf("Error when retrieving balance for address %s: %s", DefiantLabOldAccAddress, err)
		return err
	}

	err = sendCoins(ctx, keepers.BankKeeper, DefiantLabsOldAcc, DefiantLabsAcc, balance)
	if err != nil {
		return err
	}

	// Moving delegations
	err = moveSelfDelegation(ctx, keepers, DefiantLabsOldAcc, DefiantLabsAcc, DefiantLabsValAddress)
	if err != nil {
		return err
	}

	err = moveDelegations(ctx, keepers, DefiantLabsOldAcc, DanVal)
	if err != nil {
		return err
	}

	err = fixUnbondingHeight(ctx, keepers, DanVal)
	if err != nil {
		log.Printf("Error when updating unbonding height %s: %s", DefiantLabOldAccAddress, err)
		return err
	}
	return nil
}

func fixMainnet3(ctx context.Context, keepers *keepers.AppKeepers) error {
	goCtx := sdk.UnwrapSDKContext(ctx)

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

	goCtx.Logger().Info(fmt.Sprintf("Sending tokens from %s to %s", OldMainnet3Addr, NewMainnet3Addr))
	// Creating new Mainnet3 validator
	Odin3OldValAddress, err := addrToValAddr(OdinMainnet3OldAccAddress)
	if err != nil {
		return err
	}

	Odin3OldVal, err := keepers.StakingKeeper.GetValidator(ctx, Odin3OldValAddress)
	if err != nil {
		log.Printf("Validator with %v has not been found", Odin3OldValAddress)
		return err
	}

	Odin3ValAddr, err := addrToValAddr(OdinMainnet3NewAccAddress)
	if err != nil {
		return err
	}

	// rewards and comission
	err = withdrawRewardsAndCommission(ctx, keepers.AccountKeeper, *keepers.StakingKeeper, keepers.DistrKeeper, Odin3OldValAddress, Odin3ValAddr)
	if err != nil {
		return err
	}

	goCtx.Logger().Info(fmt.Sprintf("Sending tokens from %s to %s", OldMainnet3Addr, NewMainnet3Addr))
	balance, err := getBalance(ctx, keepers.AccountKeeper, keepers.BankKeeper, OldMainnet3Addr)
	if err != nil {
		log.Printf("Error when retrieving balance for address %s: %s", OldMainnet3Addr, err)
		return err
	}

	// sending balances
	err = sendCoins(ctx, keepers.BankKeeper, OldMainnet3Addr, NewMainnet3Addr, balance)
	if err != nil {
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

	Odin3Val, err := createValidator(ctx, keepers, string(Odin3ValAddr), &Odin3PubKey, Odin3OldVal.Description, Odin3OldVal.Commission)
	if err != nil {
		return err
	}

	goCtx.Logger().Info(fmt.Sprintf("Moving validator delegations from %s to %s", Odin3OldValAddress, Odin3ValAddr))

	err = moveValidatorDelegations(ctx, keepers.AccountKeeper, *keepers.StakingKeeper, keepers.DistrKeeper, keepers.BankKeeper, Odin3OldVal, Odin3Val, math.NewInt(100000000))
	if err != nil {
		return err
	}

	// moveSelfDelegation(ctx, keepers, OldMainnet3Addr, NewMainnet3Addr, Odin3ValAddr)
	err = moveDelegations(ctx, keepers, sdk.AccAddress(OdinMainnet3OldAccAddress), Odin3Val)
	if err != nil {
		return err
	}

	Odin3Val.UpdateStatus(stakingtypes.Bonded)

	// Showing all validator powers
	log.Printf("Validator power after update:")
	validators, err := keepers.StakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}

	for _, validator := range validators {
		log.Printf("%v: %v", validator.OperatorAddress, validator.ConsensusPower(keepers.StakingKeeper.PowerReduction(ctx)))
	}

	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.AppManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		goCtx := sdk.UnwrapSDKContext(ctx)
		goCtx.Logger().Info("running v7_10 upgrade handler")

		log.Printf("Validator power before update:")
		validators, err := keepers.StakingKeeper.GetAllValidators(ctx)
		if err != nil {
			return nil, err
		}

		for _, validator := range validators {
			log.Printf("%v: %v", validator.OperatorAddress, validator.ConsensusPower(keepers.StakingKeeper.PowerReduction(ctx)))
		}

		Odin3OldValAddress, err := addrToValAddr(OdinMainnet3OldAccAddress)
		if err != nil {
			return nil, err
		}

		oldMainmet3Val, err := keepers.StakingKeeper.GetValidator(ctx, Odin3OldValAddress)
		if err != nil {
			log.Printf("failed to find old mainnet3 validator")
			return nil, errors.New("failed to find old mainnet3 validator")
		}

		oldMainmet3Val.Jailed = true

		err = keepers.StakingKeeper.SetValidator(ctx, oldMainmet3Val)
		if err != nil {
			return nil, err
		}

		// Fixinng Dan's validator account association
		err = fixDefiantLabs(ctx, keepers)
		if err != nil {
			return nil, err
		}

		err = fixMainnet3(ctx, keepers)
		if err != nil {
			return nil, err
		}

		log.Printf("Flushing IBC packets...")
		FlushIBCPackets(goCtx, keepers)
		log.Printf("Flushing IBC packets complete")

		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			log.Printf("Error when running migrations: %s", err)
			return nil, err
		}
		return newVM, err
	}
}

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v0.7.11",
	CreateUpgradeHandler: CreateUpgradeHandler,
}
