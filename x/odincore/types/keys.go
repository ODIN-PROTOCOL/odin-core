package types

const (
	// ModuleName defines the module name
	ModuleName = "odincore"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_odincore"
)

var (
	ParamsKey = []byte("p_odincore")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
