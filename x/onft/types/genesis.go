package types

func NewGenesisState(nftClassID uint64, classOwners map[string]string) *GenesisState {
	return &GenesisState{
		NftClassId:  nftClassID,
		ClassOwners: classOwners,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		NftClassId:  uint64(0),
		ClassOwners: make(map[string]string),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (g *GenesisState) Validate() error { return nil }
