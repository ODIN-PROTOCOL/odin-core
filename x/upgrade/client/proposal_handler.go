package client

import (
	govclient "github.com/GeoDB-Limited/odin-core/x/gov/client"
	upgraderest "github.com/GeoDB-Limited/odin-core/x/upgrade/client/rest"
	upgradecli "github.com/cosmos/cosmos-sdk/x/upgrade/client/cli"
)

var ProposalHandler = govclient.NewProposalHandler(upgradecli.NewCmdSubmitUpgradeProposal, upgraderest.ProposalRESTHandler)
var CancelProposalHandler = govclient.NewProposalHandler(upgradecli.NewCmdSubmitCancelUpgradeProposal, upgraderest.ProposalCancelRESTHandler)
