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

// Route implements the sdk.Msg interface.
func (msg MsgWithdrawCoinsToAccFromTreasury) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface.
func (msg MsgWithdrawCoinsToAccFromTreasury) Type() string {
	return TypeMsgWithdrawCoinsToAccFromTreasury
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

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgWithdrawCoinsToAccFromTreasury) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
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

// Route implements the sdk.Msg interface.
func (msg MsgMintCoins) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface.
func (msg MsgMintCoins) Type() string {
	return TypeMsgMintCoins
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

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgMintCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgMintCoins) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
