package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/log"

	errortypes "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/ODIN-PROTOCOL/odin-core/x/bank/types"
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

	distrKeeper   types.DistributionKeeper
	mintKeeper    types.MintKeeper
	accountKeeper types.AccountKeeper
}

// NewWrappedBankKeeperBurnToCommunityPool creates a new instance of WrappedBankKeeper
// with its distrKeeper and accountKeeper members set to nil.
func NewWrappedBankKeeperBurnToCommunityPool(bk bankkeeper.Keeper, acc types.AccountKeeper) WrappedBankKeeper {
	return WrappedBankKeeper{bk, nil, nil, acc}
}

// SetDistrKeeper sets distr module keeper for this WrappedBankKeeper instance.
func (k *WrappedBankKeeper) SetDistrKeeper(dk types.DistributionKeeper) {
	k.distrKeeper = dk
}

func (k *WrappedBankKeeper) SetMintKeeper(mintKeeper types.MintKeeper) {
	k.mintKeeper = mintKeeper
}

// Logger returns a module-specific logger.
func (k WrappedBankKeeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/wrappedbank")
}

// BurnCoins moves the specified amount of coins from the given module name to
// the community pool. The total bank of the coins will not change.
func (k WrappedBankKeeper) BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	// If distrKeeper is not set OR we want to burn coins in distr itself, we will
	// just use the original BurnCoins function.

	if k.distrKeeper == nil || moduleName == distrtypes.ModuleName {
		return k.Keeper.BurnCoins(ctx, moduleName, amt)
	}

	// Create the account if it doesn't yet exist.
	acc := k.accountKeeper.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(errortypes.Wrapf(
			sdkerrors.ErrUnknownAddress,
			"module account %s does not exist", moduleName,
		))
	}

	if !acc.HasPermission(authtypes.Burner) {
		panic(errortypes.Wrapf(
			sdkerrors.ErrUnauthorized,
			"module account %s does not have permissions to burn tokens",
			moduleName,
		))
	}

	// Instead of burning coins, we send them to the community pool.
	err := k.distrKeeper.FundCommunityPool(ctx, amt, acc.GetAddress())
	if err != nil {
		return err
	}

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf(
		"sent %s from %s module account to community pool", amt.String(), moduleName,
	))
	return nil
}

// MintCoins does not create any new coins, just gets them from the community pull
func (k WrappedBankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	if k.distrKeeper == nil || moduleName == distrtypes.ModuleName {
		return k.Keeper.MintCoins(ctx, moduleName, amt)
	}

	mintParams, err := k.mintKeeper.GetParams(ctx)
	if err != nil {
		return errortypes.Wrap(err, "failed to get mint module params")
	}

	if mintParams.MintAir {
		return k.Keeper.MintCoins(ctx, moduleName, amt)
	}
	acc := k.accountKeeper.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(errortypes.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(authtypes.Minter) {
		panic(errortypes.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to mint tokens", moduleName))
	}

	logger := k.Logger(ctx)
	err = k.SendCoinsFromModuleToModule(ctx, distrtypes.ModuleName, moduleName, amt)
	if err != nil {
		err = errortypes.Wrap(err, fmt.Sprintf("failed to mint %s from %s module account", amt.String(), moduleName))
		logger.Error(err.Error())
		return err
	}
	logger.Info(fmt.Sprintf("minted %s from %s module account", amt.String(), moduleName))

	return nil
}
