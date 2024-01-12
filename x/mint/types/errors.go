package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidMintDenom                     = errors.Register(ModuleName, 121, "The given mint denom is invalid")
	ErrAccountIsNotEligible                 = errors.Register(ModuleName, 122, "The given account is not eligible to mint")
	ErrInvalidWithdrawalAmount              = errors.Register(ModuleName, 123, "The given withdrawal amount is invalid")
	ErrExceedsWithdrawalLimitPerTime        = errors.Register(ModuleName, 124, "The given amount exceeds the withdrawal limit per time")
	ErrWithdrawalAmountExceedsModuleBalance = errors.Register(ModuleName, 125, "The given amount to withdraw exceeds module balance")
	ErrMintVolumeExceedsLimit               = errors.Register(ModuleName, 126, "The given volume to mint exceeds allowed mint volume")
)
