package simulation_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/simapp"
	"github.com/ODIN-PROTOCOL/odin-core/x/mint/simulation"
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"
)

func TestDecodeStore(t *testing.T) {
	cdc := simapp.MakeTestEncodingConfig()
	dec := simulation.NewDecodeStore(cdc.Marshaler)

	minter := minttypes.NewMinter(sdk.OneDec(), sdk.NewDec(15), sdk.NewCoins(sdk.NewCoin("minigeo", sdk.NewInt(100000000))))

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: minttypes.MinterKey, Value: cdc.Marshaler.MustMarshal(&minter)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}
	testCases := []struct {
		name        string
		expectedLog string
	}{
		{"Minter", fmt.Sprintf("%v\n%v", minter, minter)},
		{"other", ""},
	}

	for i, tc := range testCases {
		i, tt := i, tc
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(testCases) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
