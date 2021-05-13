package client

import (
	distrrest "github.com/GeoDB-Limited/odin-core/x/distribution/client/rest"
	govclient "github.com/GeoDB-Limited/odin-core/x/gov/client"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/cli"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, distrrest.ProposalRESTHandler)
)
