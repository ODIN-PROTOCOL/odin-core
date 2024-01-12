package keeper

import (
	"github.com/ODIN-PROTOCOL/odin-core/x/odincore/types"
)

var _ types.QueryServer = Keeper{}
