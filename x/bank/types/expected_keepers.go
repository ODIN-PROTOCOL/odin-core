package types

import (
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	//minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// AccountKeeper defines the account contract that must be fulfilled when
// creating a x/bank keeper.
type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
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
