package types_test

import (
	govtypes "github.com/GeoDB-Limited/odin-core/x/gov/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	coinsPos   = sdk.NewCoins(sdk.NewInt64Coin(govtypes.DefaultBondDenom, 1000))
	coinsZero  = sdk.NewCoins()
	coinsMulti = sdk.NewCoins(sdk.NewInt64Coin(govtypes.DefaultBondDenom, 1000), sdk.NewInt64Coin("foo", 10000))
	addrs      = []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress("test2"),
	}
)

func init() {
	coinsMulti.Sort()
}

// test ValidateBasic for MsgCreateValidator
func TestMsgSubmitProposal(t *testing.T) {
	tests := []struct {
		title, description string
		proposalType       string
		proposerAddr       sdk.AccAddress
		initialDeposit     sdk.Coins
		expectPass         bool
	}{
		{"Test Proposal", "the purpose of this proposal is to test", govtypes.ProposalTypeText, addrs[0], coinsPos, true},
		{"", "the purpose of this proposal is to test", govtypes.ProposalTypeText, addrs[0], coinsPos, false},
		{"Test Proposal", "", govtypes.ProposalTypeText, addrs[0], coinsPos, false},
		{"Test Proposal", "the purpose of this proposal is to test", govtypes.ProposalTypeText, sdk.AccAddress{}, coinsPos, false},
		{"Test Proposal", "the purpose of this proposal is to test", govtypes.ProposalTypeText, addrs[0], coinsZero, true},
		{"Test Proposal", "the purpose of this proposal is to test", govtypes.ProposalTypeText, addrs[0], coinsMulti, true},
		{strings.Repeat("#", govtypes.MaxTitleLength*2), "the purpose of this proposal is to test", govtypes.ProposalTypeText, addrs[0], coinsMulti, false},
		{"Test Proposal", strings.Repeat("#", govtypes.MaxDescriptionLength*2), govtypes.ProposalTypeText, addrs[0], coinsMulti, false},
	}

	for i, tc := range tests {
		msg, err := govtypes.NewMsgSubmitProposal(
			govtypes.ContentFromProposalType(tc.title, tc.description, tc.proposalType),
			tc.initialDeposit,
			tc.proposerAddr,
		)

		require.NoError(t, err)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgDepositGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("addr1")
	msg := govtypes.NewMsgDeposit(addr, 0, coinsPos)
	res := msg.GetSignBytes()

	expected := `{"type":"gov/MsgDeposit","value":{"amount":[{"amount":"1000","denom":"odin"}],"depositor":"cosmos1v9jxgu33kfsgr5","proposal_id":"0"}}`
	require.Equal(t, expected, string(res))
}

// test ValidateBasic for MsgDeposit
func TestMsgDeposit(t *testing.T) {
	tests := []struct {
		proposalID    uint64
		depositorAddr sdk.AccAddress
		depositAmount sdk.Coins
		expectPass    bool
	}{
		{0, addrs[0], coinsPos, true},
		{1, sdk.AccAddress{}, coinsPos, false},
		{1, addrs[0], coinsZero, true},
		{1, addrs[0], coinsMulti, true},
	}

	for i, tc := range tests {
		msg := govtypes.NewMsgDeposit(tc.depositorAddr, tc.proposalID, tc.depositAmount)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

// test ValidateBasic for MsgDeposit
func TestMsgVote(t *testing.T) {
	tests := []struct {
		proposalID uint64
		voterAddr  sdk.AccAddress
		option     govtypes.VoteOption
		expectPass bool
	}{
		{0, addrs[0], govtypes.OptionYes, true},
		{0, sdk.AccAddress{}, govtypes.OptionYes, false},
		{0, addrs[0], govtypes.OptionNo, true},
		{0, addrs[0], govtypes.OptionNoWithVeto, true},
		{0, addrs[0], govtypes.OptionAbstain, true},
		{0, addrs[0], govtypes.VoteOption(0x13), false},
	}

	for i, tc := range tests {
		msg := govtypes.NewMsgVote(tc.voterAddr, tc.proposalID, tc.option)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

// this tests that Amino JSON MsgSubmitProposal.GetSignBytes() still works with Content as Any using the ModuleCdc
func TestMsgSubmitProposal_GetSignBytes(t *testing.T) {
	msg, err := govtypes.NewMsgSubmitProposal(govtypes.NewTextProposal("test", "abcd"), sdk.NewCoins(), sdk.AccAddress{})
	require.NoError(t, err)
	var bz []byte
	require.NotPanics(t, func() {
		bz = msg.GetSignBytes()
	})
	require.Equal(t,
		`{"type":"gov/MsgSubmitProposal","value":{"content":{"type":"gov/TextProposal","value":{"description":"abcd","title":"test"}},"initial_deposit":[]}}`,
		string(bz))
}
