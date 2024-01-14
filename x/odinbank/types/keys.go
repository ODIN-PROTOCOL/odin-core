package types

const (
	// ModuleName defines the module name
	ModuleName = "odinbank"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_odinbank"

    
)

var (
	ParamsKey = []byte("p_odinbank")
)


func KeyPrefix(p string) []byte {
    return []byte(p)
}
