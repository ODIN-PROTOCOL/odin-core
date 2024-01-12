package types

const (
	// ModuleName
	ModuleName = "mint"

	// ModuleVersion defines the current module version
	ModuleVersion = 1

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// RouterKey is the message route for mint
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the minting querier
	QueryParams               = "parameters"
	QueryInflation            = "inflation"
	QueryAnnualProvisions     = "annual_provisions"
	QueryIntegrationAddresses = "integration_addresses"
	QueryTreasuryPool         = "treasury_pool"
	QueryCurrentMintVolume    = "current_mint_volume"
)

var (
	// GlobalStoreKeyPrefix is used as prefix for the store keys
	GlobalStoreKeyPrefix = []byte{0x00}
	// MinterKey is used for the keeper store
	MinterKey = append(GlobalStoreKeyPrefix, []byte("Minter")...)
	ParamsKey = []byte{0x01}
	// MintPoolStoreKey is the key for global mint pool state
	MintPoolStoreKey          = append(GlobalStoreKeyPrefix, []byte("MintPool")...)
	MintModuleCoinsAccountKey = append(GlobalStoreKeyPrefix, []byte("MintModuleCoinsAccount")...)
)
