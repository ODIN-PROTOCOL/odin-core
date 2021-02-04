package bank

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// WrappedBankKeeper encapsulates the underlying bank keeper and overrides
// its BurnCoins function to send the coins to the community pool instead of
// just destroying them.
//
// Note that distrKeeper keeps the reference to the distr module keeper.
// Due to the circular dependency between bank-distr, distrKeeper
// cannot be initialized when the struct is created. Rather, SetDistrKeeper and SetAccountKeeper
// are expected to be called to set `distrKeeper` and `accountKeeper` respectively.
type WrappedBankKeeper struct {
	bankkeeper.Keeper

	distrKeeper   *distrkeeper.Keeper
	accountKeeper banktypes.AccountKeeper
}

// NewWrappedBankKeeperBurnToCommunityPool creates a new instance of WrappedBankKeeper
// with its distrKeeper and accountKeeper members set to nil.
func NewWrappedBankKeeperBurnToCommunityPool(bk bankkeeper.Keeper) WrappedBankKeeper {
	return WrappedBankKeeper{bk, nil, nil}
}

// SetDistrKeeper sets distr module keeper for this WrappedBankKeeper instance.
func (k *WrappedBankKeeper) SetDistrKeeper(dk *distrkeeper.Keeper) {
	k.distrKeeper = dk
}

// SetAccountKeeper sets account module keeper for this WrappedBankKeeper instance.
func (k *WrappedBankKeeper) SetAccountKeeper(ak banktypes.AccountKeeper) {
	k.accountKeeper = ak
}

// Logger returns a module-specific logger.
func (k WrappedBankKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprint("x/wrappedbank"))
}

// BurnCoins moves the specified amount of coins from the given module name to
// the community pool. The total bank of the coins will not change.
func (k WrappedBankKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	// If distrKeeper is not set OR we want to burn coins in distr itself, we will
	// just use the original BurnCoins function.
	if k.distrKeeper == nil || moduleName == distrtypes.ModuleName {
		return k.Keeper.BurnCoins(ctx, moduleName, amt)
	}

	// Create the account if it doesn't yet exist.
	acc := k.accountKeeper.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(sdkerrors.Wrapf(
			sdkerrors.ErrUnknownAddress,
			"module account %s does not exist", moduleName,
		))
	}

	if !acc.HasPermission(authtypes.Burner) {
		panic(sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"module account %s does not have permissions to burn tokens",
			moduleName,
		))
	}

	// Instead of burning coins, we send them to the community pool.
	k.SendCoinsFromModuleToModule(ctx, moduleName, distrtypes.ModuleName, amt)
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(amt...)...)
	k.distrKeeper.SetFeePool(ctx, feePool)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf(
		"sent %s from %s module account to community pool", amt.String(), moduleName,
	))
	return nil
}