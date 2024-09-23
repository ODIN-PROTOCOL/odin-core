package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrClassNotFound       = sdkerrors.New(ModuleName, 1, "class not found")
	ErrSenderNotAuthorized = sdkerrors.Register(ModuleName, 2, "sender not authorized")
)
