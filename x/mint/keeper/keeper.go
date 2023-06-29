package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc              codec.BinaryCodec
	storeKey         sdk.StoreKey
	paramSpace       paramtypes.Subspace
	stakingKeeper    minttypes.StakingKeeper
	authKeeper       minttypes.AccountKeeper
	bankKeeper       minttypes.BankKeeper
	feeCollectorName string
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	sk minttypes.StakingKeeper, ak minttypes.AccountKeeper, bk minttypes.BankKeeper,
	feeCollectorName string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(minttypes.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(minttypes.ParamKeyTable())
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace,
		stakingKeeper:    sk,
		bankKeeper:       bk,
		authKeeper:       ak,
		feeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+minttypes.ModuleName)
}

// get the minter
func (k Keeper) GetMinter(ctx sdk.Context) (minter minttypes.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(minttypes.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshal(b, &minter)
	return
}

// set the minter
func (k Keeper) SetMinter(ctx sdk.Context, minter minttypes.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(minttypes.MinterKey, b)
}

// get the module coins account
func (k Keeper) GetMintModuleCoinsAccount(ctx sdk.Context) (account sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(minttypes.MintModuleCoinsAccountKey)
	if b == nil {
		return nil
	}

	return sdk.AccAddress(b)
}

// set the module coins account
func (k Keeper) SetMintModuleCoinsAccount(ctx sdk.Context, account sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(minttypes.MintModuleCoinsAccountKey, account)
}

// GetMintPool returns the mint pool info
func (k Keeper) GetMintPool(ctx sdk.Context) (mintPool minttypes.MintPool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(minttypes.MintPoolStoreKey)
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}

	k.cdc.MustUnmarshal(b, &mintPool)
	return
}

// SetMintPool sets mint pool to the store
func (k Keeper) SetMintPool(ctx sdk.Context, mintPool minttypes.MintPool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&mintPool)
	store.Set(minttypes.MintPoolStoreKey, b)
}

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params minttypes.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params minttypes.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetMintAccount returns the mint ModuleAccount
func (k Keeper) GetMintAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.authKeeper.GetModuleAccount(ctx, minttypes.ModuleName)
}

// SetMintAccount sets the module account
func (k Keeper) SetMintAccount(ctx sdk.Context, moduleAcc authtypes.ModuleAccountI) {
	k.authKeeper.SetModuleAccount(ctx, moduleAcc)
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx sdk.Context) sdk.Int {
	return k.stakingKeeper.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	return k.stakingKeeper.BondedRatio(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, minttypes.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, k.feeCollectorName, fees)
}

// LimitExceeded checks if withdrawal amount exceeds the limit
func (k Keeper) LimitExceeded(ctx sdk.Context, amt sdk.Coins) bool {
	moduleParams := k.GetParams(ctx)
	return amt.IsAnyGT(moduleParams.MaxWithdrawalPerTime)
}

// IsEligibleAccount checks if addr exists in the eligible to withdraw account pool
func (k Keeper) IsEligibleAccount(ctx sdk.Context, addr string) bool {
	params := k.GetParams(ctx)

	for _, item := range params.EligibleAccountsPool {
		if item == addr {
			return true
		}
	}

	return false
}

// WithdrawCoinsFromTreasury transfers coins from treasury pool to receiver account
func (k Keeper) WithdrawCoinsFromTreasury(ctx sdk.Context, receiver sdk.AccAddress, amount sdk.Coins) error {
	mintPool := k.GetMintPool(ctx)

	if amount.IsAllGT(mintPool.TreasuryPool) {
		return sdkerrors.Wrapf(
			minttypes.ErrWithdrawalAmountExceedsModuleBalance,
			"withdrawal amount: %s exceeds %s module balance",
			amount.String(),
			minttypes.ModuleName,
		)
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, receiver, amount); err != nil {
		return sdkerrors.Wrapf(
			err,
			"failed to withdraw %s from %s module account",
			amount.String(),
			minttypes.ModuleName,
		)
	}

	mintPool.TreasuryPool = mintPool.TreasuryPool.Sub(amount)
	k.SetMintPool(ctx, mintPool)

	return nil
}

// IsAllowedMintDenom checks if denom exists in the allowed mint denoms list
func (k Keeper) IsAllowedMintDenom(ctx sdk.Context, coin sdk.Coin) bool {
	params := k.GetParams(ctx)
	denom := coin.Denom

	for i := range params.AllowedMintDenoms {
		if denom == params.AllowedMintDenoms[i].TokenUnitDenom {
			return true
		}
	}

	return false
}

// IsAllowedMinter checks if address exists in the allowed minters list
func (k Keeper) IsAllowedMinter(ctx sdk.Context, addr string) bool {
	params := k.GetParams(ctx)

	for i := range params.AllowedMinter {
		if addr == params.AllowedMinter[i] {
			return true
		}
	}

	return false
}

// MintVolumeExceeded checks if minting volume exceeds the limit
func (k Keeper) MintVolumeExceeded(ctx sdk.Context, amt sdk.Coins) bool {
	moduleParams := k.GetParams(ctx)
	minter := k.GetMinter(ctx)
	amt = amt.Add(minter.CurrentMintVolume...)
	return amt.IsAnyGT(moduleParams.MaxAllowedMintVolume)
}

// MintNewCoins issue new coins
func (k Keeper) MintNewCoins(ctx sdk.Context, amount sdk.Coins) error {
	mintPool := k.GetMintPool(ctx)
	minter := k.GetMinter(ctx)

	if err := k.bankKeeper.MintCoins(ctx, minttypes.ModuleName, amount); err != nil {
		return sdkerrors.Wrapf(
			err,
			"failed to mint %s new coins",
			amount.String(),
		)
	}

	mintPool.TreasuryPool = mintPool.TreasuryPool.Add(amount...)
	k.SetMintPool(ctx, mintPool)

	minter.CurrentMintVolume = minter.CurrentMintVolume.Add(amount...)
	k.SetMinter(ctx, minter)

	return nil
}
