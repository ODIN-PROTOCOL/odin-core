package simulation

import (
	"math/rand"

	"cosmossdk.io/math"
	"github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgUpdateParams int = 100

	OpWeightMsgUpdateParams = "op_weight_msg_update_params" //nolint:gosec
)

// ProposalMsgs defines the module weighted proposals' contents
func ProposalMsgs() []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			OpWeightMsgUpdateParams,
			DefaultWeightMsgUpdateParams,
			SimulateMsgUpdateParams,
		),
	}
}

// SimulateMsgUpdateParams returns a random MsgUpdateParams
func SimulateMsgUpdateParams(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	// use the default gov module account address as authority
	var authority sdk.AccAddress = address.Module("gov")

	params := types.DefaultParams()
	params.BlocksPerYear = uint64(simtypes.RandIntBetween(r, 1, 1000000))
	params.GoalBonded = math.LegacyNewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 100)), 2)
	params.InflationMin = math.LegacyNewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 50)), 2)
	params.InflationMax = math.LegacyNewDecWithPrec(int64(simtypes.RandIntBetween(r, 50, 100)), 2)
	params.InflationRateChange = math.LegacyNewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 100)), 2)
	params.MintDenom = simtypes.RandStringOfLength(r, 10)

	return &types.MsgUpdateParams{
		Authority: authority.String(),
		Params:    params,
	}
}
