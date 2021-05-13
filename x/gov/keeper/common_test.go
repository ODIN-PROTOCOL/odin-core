package keeper_test

import (
	"github.com/GeoDB-Limited/odin-core/x/common/testapp"
	"testing"

	"github.com/GeoDB-Limited/odin-core/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	TestProposal = types.NewTextProposal("Test", "description")
)

func createValidators(t *testing.T, accounts ...testapp.Account) ([]sdk.AccAddress, []sdk.ValAddress) {
	t.Helper()
	addrs := make([]sdk.AccAddress, 0, len(accounts))
	valAddrs := make([]sdk.ValAddress, 0, len(accounts))
	for _, acc := range accounts {
		addrs = append(addrs, acc.Address)
		valAddrs = append(valAddrs, acc.ValAddress)
	}
	return addrs, valAddrs
}
