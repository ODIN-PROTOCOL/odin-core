package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/gov module sentinel errors
var (
	ErrUnknownProposal         = sdkerrors.Register(ModuleName, 132, "unknown proposal")
	ErrInactiveProposal        = sdkerrors.Register(ModuleName, 133, "inactive proposal")
	ErrAlreadyActiveProposal   = sdkerrors.Register(ModuleName, 134, "proposal already active")
	ErrInvalidProposalContent  = sdkerrors.Register(ModuleName, 135, "invalid proposal content")
	ErrInvalidProposalType     = sdkerrors.Register(ModuleName, 136, "invalid proposal type")
	ErrInvalidVote             = sdkerrors.Register(ModuleName, 137, "invalid vote option")
	ErrInvalidGenesis          = sdkerrors.Register(ModuleName, 138, "invalid genesis state")
	ErrNoProposalHandlerExists = sdkerrors.Register(ModuleName, 139, "no handler exists for proposal type")
)
