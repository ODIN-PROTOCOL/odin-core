package gov_test

import (
	"github.com/GeoDB-Limited/odin-core/x/common/testapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/GeoDB-Limited/odin-core/x/gov/types"
	"github.com/stretchr/testify/require"
)

func TestEqualProposals(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(false, true)
	addrs := []sdk.AccAddress{testapp.TestUser1.Address, testapp.TestUser2.Address}

	SortAddresses(addrs)

	// Submit two proposals
	proposal := TestProposal
	proposal1, err := app.GovKeeper.SubmitProposal(ctx, proposal)
	require.NoError(t, err)

	proposal2, err := app.GovKeeper.SubmitProposal(ctx, proposal)
	require.NoError(t, err)

	// They are similar but their IDs should be different
	require.NotEqual(t, proposal1, proposal2)
	require.NotEqual(t, proposal1, proposal2)

	// Now create two genesis blocks
	state1 := types.GenesisState{Proposals: []types.Proposal{proposal1}}
	state2 := types.GenesisState{Proposals: []types.Proposal{proposal2}}
	require.NotEqual(t, state1, state2)
	require.False(t, state1.Equal(state2))

	// Now make proposals identical by setting both IDs to 55
	proposal1.ProposalId = 55
	proposal2.ProposalId = 55
	require.Equal(t, proposal1, proposal1)
	require.Equal(t, proposal1, proposal2)

	// Reassign proposals into state
	state1.Proposals[0] = proposal1
	state2.Proposals[0] = proposal2

	// State should be identical now..
	require.Equal(t, state1, state2)
	require.True(t, state1.Equal(state2))
}
