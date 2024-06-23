package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
)

var _ minttypes.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the mint MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) minttypes.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

var _ minttypes.MsgServer = msgServer{}

func (ms msgServer) WithdrawCoinsToAccFromTreasury(
	ctx context.Context,
	msg *minttypes.MsgWithdrawCoinsToAccFromTreasury,
) (*minttypes.MsgWithdrawCoinsToAccFromTreasuryResponse, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)

	allowed, err := ms.IsEligibleAccount(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errors.Wrapf(minttypes.ErrAccountIsNotEligible, "account: %s", msg.Sender)
	}

	exceeded, err := ms.LimitExceeded(ctx, msg.Amount)
	if err != nil {
		return nil, err
	}
	if exceeded {
		return nil, errors.Wrapf(minttypes.ErrExceedsWithdrawalLimitPerTime, "amount: %s", msg.Amount.String())
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse receiver address %s", msg.Receiver)
	}

	if err := ms.WithdrawCoinsFromTreasury(ctx, receiver, msg.Amount); err != nil {
		return nil, errors.Wrapf(err, "failed to mint %s coins to account %s", msg.Amount, msg.Receiver)
	}

	goCtx.EventManager().EmitEvent(sdk.NewEvent(
		minttypes.EventTypeWithdrawal,
		sdk.NewAttribute(minttypes.AttributeKeyWithdrawalAmount, msg.Amount.String()),
		sdk.NewAttribute(minttypes.AttributeKeyReceiver, msg.Receiver),
		sdk.NewAttribute(minttypes.AttributeKeySender, msg.Sender),
	))

	return &minttypes.MsgWithdrawCoinsToAccFromTreasuryResponse{}, nil
}

func (ms msgServer) MintCoins(
	ctx context.Context,
	msg *minttypes.MsgMintCoins,
) (*minttypes.MsgMintCoinsResponse, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)

	allowed, err := ms.IsAllowedMintDenom(ctx, msg.Amount[0])
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errors.Wrapf(minttypes.ErrInvalidMintDenom, "denom: %s", msg.Amount.GetDenomByIndex(0))
	}

	allowed, err = ms.IsAllowedMinter(ctx, msg.Sender)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errors.Wrapf(minttypes.ErrAccountIsNotEligible, "account: %s", msg.Sender)
	}

	exceeded, err := ms.MintVolumeExceeded(ctx, msg.Amount)
	if err != nil {
		return nil, err
	}
	if exceeded {
		return nil, errors.Wrapf(minttypes.ErrMintVolumeExceedsLimit, "volume: %s", msg.Amount.String())
	}

	if err := ms.MintNewCoins(ctx, msg.Amount); err != nil {
		return nil, errors.Wrapf(err, "failed to mint %s new coins", msg.Amount)
	}

	goCtx.EventManager().EmitEvent(sdk.NewEvent(
		minttypes.EventTypeMinting,
		sdk.NewAttribute(minttypes.AttributeKeyMintingVolume, msg.Amount.String()),
		sdk.NewAttribute(minttypes.AttributeKeySender, msg.Sender),
	))

	return &minttypes.MsgMintCoinsResponse{}, nil
}

// UpdateParams updates the params.
func (ms msgServer) UpdateParams(ctx context.Context, msg *minttypes.MsgUpdateParams) (*minttypes.MsgUpdateParamsResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	if err := ms.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &minttypes.MsgUpdateParamsResponse{}, nil
}
