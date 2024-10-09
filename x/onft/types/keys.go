package types

const (
	ModuleName = "onft"

	StoreKey = ModuleName

	NFTClassPrefix = "onft-"
)

var (
	NFTClassIDStoreKeyPrefix  = []byte{0x01}
	ClassOwnersStoreKeyPrefix = []byte{0x02}
)
