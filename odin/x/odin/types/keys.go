package types

const (
	// ModuleName defines the module name
	ModuleName = "odin"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_odin"
)

var (
	ParamsKey = []byte("p_odin")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
