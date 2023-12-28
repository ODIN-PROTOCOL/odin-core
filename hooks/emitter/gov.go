package emitter

import (
	"github.com/ODIN-PROTOCOL/odin-core/hooks/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

var (
	EventTypeInactiveProposal = types.EventTypeInactiveProposal
	EventTypeActiveProposal   = types.EventTypeActiveProposal
	StatusInactive            = 6
)

func (h *Hook) emitGovModule(ctx sdk.Context) {
	h.govKeeper.IterateProposals(ctx, func(proposal govv1.Proposal) (stop bool) {
		h.emitNewProposal(proposal, nil)
		return false
	})
	h.govKeeper.IterateAllDeposits(ctx, func(deposit govv1.Deposit) (stop bool) {
		h.Write("SET_DEPOSIT", common.JsDict{
			"proposal_id": deposit.ProposalId,
			"depositor":   deposit.Depositor,
			"amount":      deposit.Amount,
			"tx_hash":     nil,
		})
		return false
	})
	h.govKeeper.IterateAllVotes(ctx, func(vote govv1.Vote) (stop bool) {
		var answers []common.JsDict
		for _, voteOption := range vote.Options {
			answers = append(answers, common.JsDict{
				"answer": voteOption.Option,
				"weight": voteOption.Weight,
			})
		}

		h.Write("SET_VOTE", common.JsDict{
			"proposal_id": vote.ProposalId,
			"voter":       vote.Voter,
			"answers":     answers,
			"tx_hash":     nil,
		})
		return false
	})
}

func (h *Hook) emitNewProposal(proposal govv1.Proposal, proposer sdk.AccAddress) {
	h.Write("NEW_PROPOSAL", common.JsDict{
		"id":               proposal.Id,
		"proposer":         proposer,
		"title":            proposal.Title,
		"description":      proposal.Summary,
		"status":           int(proposal.Status),
		"submit_time":      proposal.SubmitTime.UnixNano(),
		"deposit_end_time": proposal.DepositEndTime.UnixNano(),
		"total_deposit":    proposal.TotalDeposit[len(proposal.TotalDeposit)-1].String(),
		"voting_time":      proposal.VotingStartTime.UnixNano(),
		"voting_end_time":  proposal.VotingEndTime.UnixNano(),
	})
}

func (h *Hook) emitSetDeposit(ctx sdk.Context, txHash []byte, id uint64, depositor sdk.AccAddress) {
	deposit, _ := h.govKeeper.GetDeposit(ctx, id, depositor)
	h.Write("SET_DEPOSIT", common.JsDict{
		"proposal_id": id,
		"depositor":   depositor,
		"amount":      deposit.Amount[len(deposit.Amount)-1].String(),
		"tx_hash":     txHash,
	})
}

func (h *Hook) emitUpdateProposalAfterDeposit(ctx sdk.Context, id uint64) {
	proposal, _ := h.govKeeper.GetProposal(ctx, id)
	h.Write("UPDATE_PROPOSAL", common.JsDict{
		"id":              id,
		"status":          int(proposal.Status),
		"total_deposit":   proposal.TotalDeposit[len(proposal.TotalDeposit)-1].String(),
		"voting_time":     proposal.VotingStartTime.UnixNano(),
		"voting_end_time": proposal.VotingEndTime.UnixNano(),
	})
}

// handleMsgSubmitProposal implements emitter handler for MsgSubmitProposal.
func (app *Hook) handleMsgSubmitProposal(
	ctx sdk.Context, txHash []byte, msg *govv1.MsgSubmitProposal, evMap common.EvMap,
) {
	proposalId := uint64(common.Atoi(evMap[types.EventTypeSubmitProposal+"."+types.AttributeKeyProposalID][0]))
	proposal, _ := app.govKeeper.GetProposal(ctx, proposalId)
	app.Write("NEW_PROPOSAL", common.JsDict{
		"id":               proposal.Id,
		"proposer":         msg.Proposer,
		"title":            proposal.Title,
		"description":      proposal.Summary,
		"status":           int(proposal.Status),
		"submit_time":      proposal.SubmitTime.UnixNano(),
		"deposit_end_time": proposal.DepositEndTime.UnixNano(),
		"total_deposit":    proposal.TotalDeposit[len(proposal.TotalDeposit)-1].String(),
		"voting_time":      proposal.VotingStartTime.UnixNano(),
		"voting_end_time":  proposal.VotingEndTime.UnixNano(),
	})
	proposer, _ := sdk.AccAddressFromBech32(msg.Proposer)
	app.emitSetDeposit(ctx, txHash, proposalId, proposer)
}

// handleMsgDeposit implements emitter handler for MsgDeposit.
func (h *Hook) handleMsgDeposit(
	ctx sdk.Context, txHash []byte, msg *govv1.MsgDeposit,
) {
	depositor, _ := sdk.AccAddressFromBech32(msg.Depositor)
	h.emitSetDeposit(ctx, txHash, msg.ProposalId, depositor)
	h.emitUpdateProposalAfterDeposit(ctx, msg.ProposalId)
}

// handleMsgVote implements emitter handler for MsgVote.
func (h *Hook) handleMsgVote(
	txHash []byte, msg *govv1.MsgVote,
) {
	h.Write("SET_VOTE", common.JsDict{
		"proposal_id": msg.ProposalId,
		"voter":       msg.Voter,
		"answer":      int(msg.Option),
		"tx_hash":     txHash,
	})
}

func (h *Hook) handleEventInactiveProposal(evMap common.EvMap) {
	h.Write("UPDATE_PROPOSAL", common.JsDict{
		"id":     common.Atoi(evMap[types.EventTypeInactiveProposal+"."+types.AttributeKeyProposalID][0]),
		"status": StatusInactive,
	})
}

func (h *Hook) handleEventTypeActiveProposal(ctx sdk.Context, evMap common.EvMap) {
	id := uint64(common.Atoi(evMap[types.EventTypeActiveProposal+"."+types.AttributeKeyProposalID][0]))
	proposal, _ := h.govKeeper.GetProposal(ctx, id)
	h.Write("UPDATE_PROPOSAL", common.JsDict{
		"id":     id,
		"status": int(proposal.Status),
	})
}
