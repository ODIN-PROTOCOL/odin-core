package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgWithdrawCoinsToAccFromTreasury = "withdraw_coins_from_treasury"
	TypeMsgMintCoins                      = "mint_coins"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgWithdrawCoinsToAccFromTreasury{}
	_ sdk.Msg = &MsgMintCoins{}
)

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	if err := m.Params.Validate(); err != nil {
		return err
	}

	return nil
}

// NewMsgWithdrawCoinsToAccFromTreasury returns a new MsgWithdrawCoinsToAccFromTreasury
func NewMsgWithdrawCoinsToAccFromTreasury(
	amt sdk.Coins,
	receiver sdk.AccAddress,
	sender sdk.AccAddress,
) MsgWithdrawCoinsToAccFromTreasury {
	return MsgWithdrawCoinsToAccFromTreasury{
		Amount:   amt,
		Receiver: receiver.String(),
		Sender:   sender.String(),
	}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgWithdrawCoinsToAccFromTreasury) ValidateBasic() error {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return err
	}
	if err := sdk.VerifyAddressFormat(sender); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if err := sdk.VerifyAddressFormat(receiver); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "receiver: %s", msg.Receiver)
	}
	if !msg.Amount.IsValid() {
		return errors.Wrapf(sdkerrors.ErrInvalidCoins, "amount: %s", msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return errors.Wrapf(ErrInvalidWithdrawalAmount, "amount: %s", msg.Amount.String())
	}

	return nil
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgWithdrawCoinsToAccFromTreasury) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgMintCoins returns a new MsgMintCoins
func NewMsgMintCoins(
	amt sdk.Coins,
	sender sdk.AccAddress,
) MsgMintCoins {
	return MsgMintCoins{
		Amount: amt,
		Sender: sender.String(),
	}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgMintCoins) ValidateBasic() error {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	if err := sdk.VerifyAddressFormat(sender); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if !msg.Amount.IsValid() {
		return errors.Wrapf(sdkerrors.ErrInvalidCoins, "amount: %s", msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return errors.Wrapf(ErrInvalidWithdrawalAmount, "amount: %s", msg.Amount.String())
	}

	return nil
}
