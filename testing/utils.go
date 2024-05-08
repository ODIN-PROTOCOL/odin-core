package ibctesting

import (
	"fmt"
	"math/rand"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/stretchr/testify/require"
)

// ApplyValSetChanges takes in cmttypes.ValidatorSet and []abci.ValidatorUpdate and will return a new cmttypes.ValidatorSet which has the
// provided validator updates applied to the provided validator set.
func ApplyValSetChanges(t *testing.T, valSet *cmttypes.ValidatorSet, valUpdates []abci.ValidatorUpdate) *cmttypes.ValidatorSet {
	updates, err := cmttypes.PB2TM.ValidatorUpdates(valUpdates)
	require.NoError(t, err)

	// must copy since validator set will mutate with UpdateWithChangeSet
	newVals := valSet.Copy()
	err = newVals.UpdateWithChangeSet(updates)
	require.NoError(t, err)

	return newVals
}

// VoteAndCheckProposalStatus votes on a gov proposal, checks if the proposal has passed, and returns an error if it has not with the failure reason.
func VoteAndCheckProposalStatus(endpoint *Endpoint, proposalID uint64) error {
	// vote on proposal
	ctx := endpoint.Chain.GetContext()
	require.NoError(endpoint.Chain.T, endpoint.Chain.App.GovKeeper.AddVote(
		ctx,
		proposalID,
		endpoint.Chain.SenderAccount.GetAddress(),
		govtypesv1.NewNonSplitVoteOption(govtypesv1.OptionYes),
		"",
	))

	// fast forward the chain context to end the voting period
	params, err := endpoint.Chain.App.GovKeeper.Params.Get(ctx)
	require.NoError(endpoint.Chain.T, err)

	endpoint.Chain.Coordinator.IncrementTimeBy(*params.VotingPeriod + *params.MaxDepositPeriod)
	endpoint.Chain.NextBlock()

	// check if proposal passed or failed on msg execution
	// we need to grab the context again since the previous context is no longer valid as the chain header time has been incremented
	p, err := endpoint.Chain.App.GovKeeper.Proposals.Get(endpoint.Chain.GetContext(), proposalID)
	require.NoError(endpoint.Chain.T, err)
	if p.Status != govtypesv1.StatusPassed {
		return fmt.Errorf("proposal failed: %s", p.FailedReason)
	}
	return nil
}

// GenerateString generates a random string of the given length in bytes
func GenerateString(length uint) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}

// UnmarshalMsgResponses parse out msg responses from a transaction result
func UnmarshalMsgResponses(cdc codec.Codec, data []byte, msgs ...codec.ProtoMarshaler) error {
	var txMsgData sdk.TxMsgData
	if err := cdc.Unmarshal(data, &txMsgData); err != nil {
		return err
	}

	if len(msgs) != len(txMsgData.MsgResponses) {
		return fmt.Errorf("expected %d message responses but got %d", len(msgs), len(txMsgData.MsgResponses))
	}

	for i, msg := range msgs {
		if err := cdc.Unmarshal(txMsgData.MsgResponses[i].Value, msg); err != nil {
			return err
		}
	}

	return nil
}
