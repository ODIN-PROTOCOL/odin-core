package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc              codec.BinaryCodec
	storeService     storetypes.KVStoreService
	stakingKeeper    minttypes.StakingKeeper
	authKeeper       minttypes.AccountKeeper
	bankKeeper       minttypes.BankKeeper
	feeCollectorName string

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	Schema   collections.Schema
	Params   collections.Item[minttypes.Params]
	Minter   collections.Item[minttypes.Minter]
	MintPool collections.Item[minttypes.MintPool]
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	sk minttypes.StakingKeeper,
	ak minttypes.AccountKeeper,
	bk minttypes.BankKeeper,
	feeCollectorName string,
	authority string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(minttypes.ModuleName); addr == nil {
		panic(fmt.Sprintf("the x/%s module account has not been set", minttypes.ModuleName))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:              cdc,
		storeService:     storeService,
		stakingKeeper:    sk,
		bankKeeper:       bk,
		authKeeper:       ak,
		feeCollectorName: feeCollectorName,
		authority:        authority,
		Params:           collections.NewItem(sb, minttypes.ParamsKey, "params", codec.CollValue[minttypes.Params](cdc)),
		Minter:           collections.NewItem(sb, minttypes.MinterKey, "minter", codec.CollValue[minttypes.Minter](cdc)),
		MintPool:         collections.NewItem(sb, minttypes.MintPoolStoreKey, "mint_pool", codec.CollValue[minttypes.MintPool](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

// GetAuthority returns the x/mint module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+minttypes.ModuleName)
}

// get the minter
func (k Keeper) GetMinter(ctx context.Context) (minttypes.Minter, error) {
	return k.Minter.Get(ctx)
}

// set the minter
func (k Keeper) SetMinter(ctx context.Context, minter minttypes.Minter) error {
	return k.Minter.Set(ctx, minter)
}

// GetMintPool returns the mint pool info
func (k Keeper) GetMintPool(ctx context.Context) (minttypes.MintPool, error) {
	return k.MintPool.Get(ctx)
}

// SetMintPool sets mint pool to the store
func (k Keeper) SetMintPool(ctx context.Context, mintPool minttypes.MintPool) error {
	return k.MintPool.Set(ctx, mintPool)
}

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx context.Context) (minttypes.Params, error) {
	return k.Params.Get(ctx)
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx context.Context, params minttypes.Params) error {
	return k.Params.Set(ctx, params)
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx context.Context) (math.Int, error) {
	return k.stakingKeeper.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx context.Context) (math.LegacyDec, error) {
	return k.stakingKeeper.BondedRatio(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx context.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, minttypes.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx context.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, k.feeCollectorName, fees)
}

// LimitExceeded checks if withdrawal amount exceeds the limit
func (k Keeper) LimitExceeded(ctx context.Context, amt sdk.Coins) (bool, error) {
	moduleParams, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	return amt.IsAnyGT(moduleParams.MaxWithdrawalPerTime), nil
}

// IsEligibleAccount checks if addr exists in the eligible to withdraw account pool
func (k Keeper) IsEligibleAccount(ctx context.Context, addr string) (bool, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	for _, item := range params.EligibleAccountsPool {
		if item == addr {
			return true, nil
		}
	}

	return false, nil
}

// WithdrawCoinsFromTreasury transfers coins from treasury pool to receiver account
func (k Keeper) WithdrawCoinsFromTreasury(ctx context.Context, receiver sdk.AccAddress, amount sdk.Coins) error {
	mintPool, err := k.GetMintPool(ctx)
	if err != nil {
		return err
	}

	if amount.IsAllGT(mintPool.TreasuryPool) {
		return errors.Wrapf(
			minttypes.ErrWithdrawalAmountExceedsModuleBalance,
			"withdrawal amount: %s exceeds %s module balance",
			amount.String(),
			minttypes.ModuleName,
		)
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, receiver, amount); err != nil {
		return errors.Wrapf(
			err,
			"failed to withdraw %s from %s module account",
			amount.String(),
			minttypes.ModuleName,
		)
	}

	for _, coinAmount := range amount {
		mintPool.TreasuryPool = mintPool.TreasuryPool.Sub(coinAmount)
	}
	err = k.SetMintPool(ctx, mintPool)
	if err != nil {
		return err
	}

	return nil
}

// IsAllowedMintDenom checks if denom exists in the allowed mint denoms list
func (k Keeper) IsAllowedMintDenom(ctx context.Context, coin sdk.Coin) (bool, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	denom := coin.Denom

	for i := range params.AllowedMintDenoms {
		if denom == params.AllowedMintDenoms[i].TokenUnitDenom {
			return true, nil
		}
	}

	return false, nil
}

// IsAllowedMinter checks if address exists in the allowed minters list
func (k Keeper) IsAllowedMinter(ctx context.Context, addr string) (bool, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	for i := range params.AllowedMinter {
		if addr == params.AllowedMinter[i] {
			return true, nil
		}
	}

	return false, nil
}

// MintVolumeExceeded checks if minting volume exceeds the limit
func (k Keeper) MintVolumeExceeded(ctx context.Context, amt sdk.Coins) (bool, error) {
	moduleParams, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	minter, err := k.GetMinter(ctx)
	if err != nil {
		return false, err
	}

	amt = amt.Add(minter.CurrentMintVolume...)
	return amt.IsAnyGT(moduleParams.MaxAllowedMintVolume), nil
}

// MintNewCoins issue new coins
func (k Keeper) MintNewCoins(ctx context.Context, amount sdk.Coins) error {
	mintPool, err := k.GetMintPool(ctx)
	if err != nil {
		return err
	}

	minter, err := k.GetMinter(ctx)
	if err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, minttypes.ModuleName, amount); err != nil {
		return errors.Wrapf(
			err,
			"failed to mint %s new coins",
			amount.String(),
		)
	}

	mintPool.TreasuryPool = mintPool.TreasuryPool.Add(amount...)
	err = k.SetMintPool(ctx, mintPool)
	if err != nil {
		return err
	}

	minter.CurrentMintVolume = minter.CurrentMintVolume.Add(amount...)
	err = k.SetMinter(ctx, minter)
	if err != nil {
		return err
	}

	return nil
}
