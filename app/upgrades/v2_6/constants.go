package v2_6

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

const UpgradeName = "v2_6"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{group.StoreKey, consensusparamtypes.StoreKey, crisistypes.StoreKey},
	},
}

var ICAAllowMessages = []string{
	sdk.MsgTypeURL(&authz.MsgExec{}),
	sdk.MsgTypeURL(&authz.MsgGrant{}),
	sdk.MsgTypeURL(&authz.MsgRevoke{}),
	sdk.MsgTypeURL(&banktypes.MsgSend{}),
	sdk.MsgTypeURL(&banktypes.MsgMultiSend{}),
	sdk.MsgTypeURL(&distrtypes.MsgSetWithdrawAddress{}),
	sdk.MsgTypeURL(&distrtypes.MsgWithdrawValidatorCommission{}),
	sdk.MsgTypeURL(&distrtypes.MsgFundCommunityPool{}),
	sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}),
	sdk.MsgTypeURL(&feegrant.MsgGrantAllowance{}),
	sdk.MsgTypeURL(&feegrant.MsgRevokeAllowance{}),
	sdk.MsgTypeURL(&govv1beta1.MsgVoteWeighted{}),
	sdk.MsgTypeURL(&govv1beta1.MsgSubmitProposal{}),
	sdk.MsgTypeURL(&govv1beta1.MsgDeposit{}),
	sdk.MsgTypeURL(&govv1beta1.MsgVote{}),
	// Change: add messages from Group module
	sdk.MsgTypeURL(&group.MsgCreateGroupPolicy{}),
	sdk.MsgTypeURL(&group.MsgCreateGroupWithPolicy{}),
	sdk.MsgTypeURL(&group.MsgCreateGroup{}),
	sdk.MsgTypeURL(&group.MsgExec{}),
	sdk.MsgTypeURL(&group.MsgLeaveGroup{}),
	sdk.MsgTypeURL(&group.MsgSubmitProposal{}),
	sdk.MsgTypeURL(&group.MsgUpdateGroupAdmin{}),
	sdk.MsgTypeURL(&group.MsgUpdateGroupMembers{}),
	sdk.MsgTypeURL(&group.MsgUpdateGroupMetadata{}),
	sdk.MsgTypeURL(&group.MsgUpdateGroupPolicyAdmin{}),
	sdk.MsgTypeURL(&group.MsgUpdateGroupPolicyDecisionPolicy{}),
	sdk.MsgTypeURL(&group.MsgUpdateGroupPolicyMetadata{}),
	sdk.MsgTypeURL(&group.MsgVote{}),
	sdk.MsgTypeURL(&group.MsgWithdrawProposal{}),
	// Change: add messages from Oracle module
	sdk.MsgTypeURL(&oracletypes.MsgActivate{}),
	sdk.MsgTypeURL(&oracletypes.MsgCreateDataSource{}),
	sdk.MsgTypeURL(&oracletypes.MsgCreateOracleScript{}),
	sdk.MsgTypeURL(&oracletypes.MsgEditDataSource{}),
	sdk.MsgTypeURL(&oracletypes.MsgEditOracleScript{}),
	sdk.MsgTypeURL(&oracletypes.MsgReportData{}),
	sdk.MsgTypeURL(&oracletypes.MsgRequestData{}),

	sdk.MsgTypeURL(&stakingtypes.MsgEditValidator{}),
	sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}),
	sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}),
	sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}),
	sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{}),
	sdk.MsgTypeURL(&vestingtypes.MsgCreateVestingAccount{}),
	sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}),
}
