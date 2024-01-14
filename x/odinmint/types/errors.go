package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/odinmint module sentinel errors
var (
	ErrInvalidMintDenom                     = sdkerrors.Register(ModuleName, 121, "The given mint denom is invalid")
	ErrAccountIsNotEligible                 = sdkerrors.Register(ModuleName, 122, "The given account is not eligible to mint")
	ErrInvalidWithdrawalAmount              = sdkerrors.Register(ModuleName, 123, "The given withdrawal amount is invalid")
	ErrExceedsWithdrawalLimitPerTime        = sdkerrors.Register(ModuleName, 124, "The given amount exceeds the withdrawal limit per time")
	ErrWithdrawalAmountExceedsModuleBalance = sdkerrors.Register(ModuleName, 125, "The given amount to withdraw exceeds module balance")
	ErrMintVolumeExceedsLimit               = sdkerrors.Register(ModuleName, 126, "The given volume to mint exceeds allowed mint volume")
)
