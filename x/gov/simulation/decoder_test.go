package simulation_test

import (
	"encoding/binary"
	"fmt"
	app "github.com/GeoDB-Limited/odin-core/app"
	"github.com/GeoDB-Limited/odin-core/x/gov/simulation"
	govtypes "github.com/GeoDB-Limited/odin-core/x/gov/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
)

func TestDecodeStore(t *testing.T) {
	cdc, _ := app.MakeCodecs()
	dec := simulation.NewDecodeStore(cdc)

	endTime := time.Now().UTC()

	content := govtypes.ContentFromProposalType("test", "test", govtypes.ProposalTypeText)
	proposal, err := govtypes.NewProposal(content, 1, endTime, endTime.Add(24*time.Hour))
	require.NoError(t, err)

	proposalIDBz := make([]byte, 8)
	binary.LittleEndian.PutUint64(proposalIDBz, 1)
	deposit := govtypes.NewDeposit(1, delAddr1, sdk.NewCoins(sdk.NewCoin(govtypes.DefaultBondDenom, sdk.OneInt())))
	vote := govtypes.NewVote(1, delAddr1, govtypes.OptionYes)

	proposalBz, err := cdc.MarshalBinaryBare(&proposal)
	require.NoError(t, err)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: govtypes.ProposalKey(1), Value: proposalBz},
			{Key: govtypes.InactiveProposalQueueKey(1, endTime), Value: proposalIDBz},
			{Key: govtypes.DepositKey(1, delAddr1), Value: cdc.MustMarshalBinaryBare(&deposit)},
			{Key: govtypes.VoteKey(1, delAddr1), Value: cdc.MustMarshalBinaryBare(&vote)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"proposals", fmt.Sprintf("%v\n%v", proposal, proposal)},
		{"proposal IDs", "proposalIDA: 1\nProposalIDB: 1"},
		{"deposits", fmt.Sprintf("%v\n%v", deposit, deposit)},
		{"votes", fmt.Sprintf("%v\n%v", vote, vote)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
