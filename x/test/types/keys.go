package types

const (
	// ModuleName defines the module name
	ModuleName = "test"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_test"

)

var (
	ParamsKey = []byte("p_test")
)



func KeyPrefix(p string) []byte {
    return []byte(p)
}
