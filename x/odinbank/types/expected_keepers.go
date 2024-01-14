package types

import (
	"context"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/odinmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) sdk.ModuleAccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}


// DistributionKeeper defines the distribution contract that must be fulfilled when
// creating a x/bank keeper.
type DistributionKeeper interface {
	GetFeePool(sdk.Context) disttypes.FeePool
	SetFeePool(sdk.Context, disttypes.FeePool)
}

type MintKeeper interface {
	GetParams(sdk.Context) minttypes.Params
}
