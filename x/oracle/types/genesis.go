package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState instance
func NewGenesisState(params Params, dataSources []DataSource, oracleScripts []OracleScript, oracleModuleCoinsAccount sdk.AccAddress) *GenesisState {
	return &GenesisState{
		Params:             params,
		DataSources:        dataSources,
		OracleScripts:      oracleScripts,
		OraclePool:         InitialOraclePool(),
		ModuleCoinsAccount: oracleModuleCoinsAccount.String(),
	}
}

// DefaultGenesisState returns the default oracle genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		DataSources:   []DataSource{},
		OracleScripts: []OracleScript{},
		OraclePool:    InitialOraclePool(),
	}
}

// GetGenesisStateFromAppState returns oracle GenesisState given raw application genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	// for i := range g.Validators {
	// 	if err := g.Validators[i].UnpackInterfaces(c); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (g GenesisState) Validate() error { return nil }
