package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter Minter, params Params, mintPool MintPool, mintModuleCoinsAccount sdk.AccAddress) *GenesisState {
	return &GenesisState{
		Minter:             minter,
		Params:             params,
		MintPool:           mintPool,
		ModuleCoinsAccount: mintModuleCoinsAccount.String(),
	}
}


// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Minter:   DefaultInitialMinter(),
		Params:   DefaultParams(),
		MintPool: InitialMintPool(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	return ValidateMinter(gs.Minter)
}
