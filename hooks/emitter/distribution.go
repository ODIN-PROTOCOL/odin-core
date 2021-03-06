package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/ODIN-PROTOCOL/odin-core/hooks/common"
)

func (h *Hook) getCurrentRewardAndCurrentRatio(ctx sdk.Context, addr sdk.ValAddress) (string, string) {
	currentReward := "0"
	currentRatio := "0"

	reward := h.distrKeeper.GetValidatorCurrentRewards(ctx, addr)
	latestReward := h.distrKeeper.GetValidatorHistoricalRewards(ctx, addr, reward.Period-1)

	if !reward.Rewards.IsZero() {
		currentReward = reward.Rewards[0].Amount.String()
	}
	if !latestReward.CumulativeRewardRatio.IsZero() {
		currentRatio = latestReward.CumulativeRewardRatio[0].Amount.String()
	}

	return currentReward, currentRatio
}

func (h *Hook) emitUpdateValidatorRewardAndAccumulatedCommission(ctx sdk.Context, addr sdk.ValAddress) {
	currentReward, currentRatio := h.getCurrentRewardAndCurrentRatio(ctx, addr)
	accCommission, _ := h.distrKeeper.GetValidatorAccumulatedCommission(ctx, addr).Commission.TruncateDecimal()
	h.Write("UPDATE_VALIDATOR", common.JsDict{
		"operator_address":       addr.String(),
		"current_reward":         currentReward,
		"current_ratio":          currentRatio,
		"accumulated_commission": accCommission.String(),
	})
}

func (h *Hook) emitUpdateValidatorReward(ctx sdk.Context, addr sdk.ValAddress) {
	currentReward, currentRatio := h.getCurrentRewardAndCurrentRatio(ctx, addr)
	h.Write("UPDATE_VALIDATOR", common.JsDict{
		"operator_address": addr.String(),
		"current_reward":   currentReward,
		"current_ratio":    currentRatio,
	})
}

// handleMsgWithdrawDelegatorReward implements emitter handler for MsgWithdrawDelegatorReward.
func (h *Hook) handleMsgWithdrawDelegatorReward(
	ctx sdk.Context, msg *types.MsgWithdrawDelegatorReward, evMap common.EvMap, extra common.JsDict,
) {
	valAddr, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	delAddr, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	withdrawAddr := h.distrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)
	h.AddAccountsInTx(withdrawAddr.String())
	h.emitUpdateValidatorReward(ctx, valAddr)
	h.emitDelegationAfterWithdrawReward(ctx, valAddr, withdrawAddr)
	val, _ := h.stakingKeeper.GetValidator(ctx, valAddr)
	extra["moniker"] = val.Description.Moniker
	extra["identity"] = val.Description.Identity
	extra["reward_amount"] = evMap[types.EventTypeWithdrawRewards+"."+sdk.AttributeKeyAmount][0]
}

// handleMsgSetWithdrawAddress implements emitter handler for MsgSetWithdrawAddress.
func (h *Hook) handleMsgSetWithdrawAddress(msg *types.MsgSetWithdrawAddress, extra common.JsDict) {
	h.AddAccountsInTx(msg.WithdrawAddress)
	extra["delegator_address"] = msg.DelegatorAddress
	extra["withdraw_address"] = msg.WithdrawAddress
}

// handleMsgWithdrawValidatorCommission implements emitter handler for MsgWithdrawValidatorCommission.
func (h *Hook) handleMsgWithdrawValidatorCommission(
	ctx sdk.Context, msg *types.MsgWithdrawValidatorCommission, evMap common.EvMap, extra common.JsDict,
) {
	delAddr, _ := sdk.AccAddressFromBech32(msg.ValidatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	withdrawAddr := h.distrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)
	h.AddAccountsInTx(withdrawAddr.String())
	h.emitUpdateValidatorRewardAndAccumulatedCommission(ctx, valAddr)
	val, _ := h.stakingKeeper.GetValidator(ctx, valAddr)
	extra["moniker"] = val.Description.Moniker
	extra["identity"] = val.Description.Identity
	extra["commission_amount"] = evMap[types.EventTypeWithdrawCommission+"."+sdk.AttributeKeyAmount][0]
}
