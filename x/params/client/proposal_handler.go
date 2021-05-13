package client

import (
	govclient "github.com/GeoDB-Limited/odin-core/x/gov/client"
	paramsrest "github.com/GeoDB-Limited/odin-core/x/params/client/rest"
	"github.com/cosmos/cosmos-sdk/x/params/client/cli"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd, paramsrest.ProposalRESTHandler)
