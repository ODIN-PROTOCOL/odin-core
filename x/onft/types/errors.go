package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrClassNotFound       = sdkerrors.New(ModuleName, 1, "class not found")
	ErrSenderNotAuthorized = sdkerrors.Register(ModuleName, 2, "sender not authorized")
	ErrEmptyClassID        = sdkerrors.Register(ModuleName, 3, "empty class id")
	ErrEmptyNFTID          = sdkerrors.Register(ModuleName, 4, "empty nft id")
)
