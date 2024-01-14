package types

const (
	// ModuleName defines the module name
	ModuleName = "odinmint"

	// ModuleVersion defines the current module version
	ModuleVersion = 1

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_odinmint"
	
	// RouterKey is the message route for mint
	RouterKey = ModuleName
	
	// Query endpoints supported by the minting querier
	QueryParams               = "parameters"
	QueryInflation            = "inflation"
	QueryAnnualProvisions     = "annual_provisions"
	QueryIntegrationAddresses = "integration_addresses"
	QueryTreasuryPool         = "treasury_pool"
	QueryCurrentMintVolume    = "current_mint_volume"
)

var (
	ParamsKey = []byte("p_odinmint")
	// GlobalStoreKeyPrefix is used as prefix for the store keys
	GlobalStoreKeyPrefix = []byte{0x00}
	// MinterKey is used for the keeper store
	MinterKey = append(GlobalStoreKeyPrefix, []byte("Minter")...)
	// MintPoolStoreKey is the key for global mint pool state
	MintPoolStoreKey          = append(GlobalStoreKeyPrefix, []byte("MintPool")...)
	MintModuleCoinsAccountKey = append(GlobalStoreKeyPrefix, []byte("MintModuleCoinsAccount")...)
)


func KeyPrefix(p string) []byte {
    return []byte(p)
}
