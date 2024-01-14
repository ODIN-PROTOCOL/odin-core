package keeper

import (
	"github.com/ODIN-PROTOCOL/odin-core/x/test/types"
)

var _ types.QueryServer = Keeper{}
