package oraclekeeper

import (
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// modify returns new value if it is not `DoNotModify`. Returns old value otherwise
func modify(oldVal string, newVal string) string {
	if newVal == types.DoNotModify {
		return oldVal
	}
	return newVal
}

func modifyFee(oldVal, newVal sdk.Coins) sdk.Coins {
	if newVal == nil {
		return oldVal
	}
	return newVal
}
