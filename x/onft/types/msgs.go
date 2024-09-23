package types

import (
	"cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgCreateNFTClass{}
	_ sdk.Msg = &MsgTransferClassOwnership{}
	_ sdk.Msg = &MsgMintNFT{}

	_ sdk.HasValidateBasic = &MsgCreateNFTClass{}
	_ sdk.HasValidateBasic = &MsgTransferClassOwnership{}
	_ sdk.HasValidateBasic = &MsgMintNFT{}
)

func NewMsgCreateNFTClass(
	name string,
	symbol string,
	description string,
	uri string,
	uriHash string,
	data *codectypes.Any,
	sender sdk.AccAddress,
) *MsgCreateNFTClass {
	return &MsgCreateNFTClass{
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Data:        data,
		Sender:      sender.String(),
	}
}

func (msg *MsgCreateNFTClass) ValidateBasic() error {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	if err := sdk.VerifyAddressFormat(sender); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}

	return nil
}

func NewMsgTransferClassOwnership(
	classID string,
	sender sdk.AccAddress,
	newOwner string,
) *MsgTransferClassOwnership {
	return &MsgTransferClassOwnership{
		ClassId:  classID,
		Sender:   sender.String(),
		NewOwner: newOwner,
	}
}

func (msg *MsgTransferClassOwnership) ValidateBasic() error {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	if err := sdk.VerifyAddressFormat(sender); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}

	newOwner, err := sdk.AccAddressFromBech32(msg.NewOwner)
	if err != nil {
		return err
	}
	if err := sdk.VerifyAddressFormat(newOwner); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "new_owner: %s", msg.Sender)
	}

	return nil
}

func NewMsgMintNFT(
	classID string,
	uri string,
	uriHash string,
	sender sdk.AccAddress,
	receiver string,
	data *codectypes.Any,
) *MsgMintNFT {
	return &MsgMintNFT{
		ClassId:  classID,
		Uri:      uri,
		UriHash:  uriHash,
		Sender:   sender.String(),
		Receiver: receiver,
		Data:     data,
	}
}

func (msg *MsgMintNFT) ValidateBasic() error {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	if err := sdk.VerifyAddressFormat(sender); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return err
	}
	if err := sdk.VerifyAddressFormat(receiver); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "new_owner: %s", msg.Sender)
	}

	return nil
}
