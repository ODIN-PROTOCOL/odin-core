package keeper_test

import (
	bandapp "github.com/GeoDB-Limited/odin-core/app"
	"github.com/GeoDB-Limited/odin-core/x/common/testapp"
	"testing"

	"github.com/GeoDB-Limited/odin-core/x/gov/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *bandapp.BandApp
	ctx         sdk.Context
	queryClient types.QueryClient
	addrs       []sdk.AccAddress
}

func (suite *KeeperTestSuite) SetupTest() {
	app, ctx, _ := testapp.CreateTestInput(false, false)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.GovKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = queryClient
	suite.addrs = []sdk.AccAddress{testapp.TestUser1.Address, testapp.TestUser2.Address}
}

func TestIncrementProposalNumber(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(true)

	tp := TestProposal
	_, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposal6, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)

	require.Equal(t, uint64(6), proposal6.ProposalId)
}

func TestProposalQueues(t *testing.T) {
	app, ctx, _ := testapp.CreateTestInput(true)

	// create test proposals
	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)

	inactiveIterator := app.GovKeeper.InactiveProposalQueueIterator(ctx, proposal.DepositEndTime)
	require.True(t, inactiveIterator.Valid())

	proposalID := types.GetProposalIDFromBytes(inactiveIterator.Value())
	require.Equal(t, proposalID, proposal.ProposalId)
	inactiveIterator.Close()

	app.GovKeeper.ActivateVotingPeriod(ctx, proposal)

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposal.ProposalId)
	require.True(t, ok)

	activeIterator := app.GovKeeper.ActiveProposalQueueIterator(ctx, proposal.VotingEndTime)
	require.True(t, activeIterator.Valid())

	proposalID, _ = types.SplitActiveProposalQueueKey(activeIterator.Key())
	require.Equal(t, proposalID, proposal.ProposalId)

	activeIterator.Close()
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
