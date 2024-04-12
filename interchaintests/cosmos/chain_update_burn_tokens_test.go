package cosmos

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades/v7_12"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	interchaintest "github.com/odin-protocol/interchaintest/v7"
	"github.com/odin-protocol/interchaintest/v7/chain/cosmos"
	"github.com/odin-protocol/interchaintest/v7/ibc"
	"github.com/odin-protocol/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
)

const (
	haltHeightDelta    = uint64(10) // will propose upgrade this many blocks in the future
	blocksAfterUpgrade = uint64(10)
	votingPeriod       = "20s"
	maxDepositPeriod   = "10s"
)

func TestOdinUpgradeBurnTokens(t *testing.T) {
	CosmosChainUpgradeBurnTokensTest(t, "odin", "v0.7.11", "odinprotocol/core", "v0.7.12", "v0.7.12")
}

func CosmosChainUpgradeBurnTokensTest(t *testing.T, chainName, initialVersion, upgradeContainerRepo, upgradeVersion string, upgradeName string) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	ctx := context.Background()

	// SDK v45 params for Juno genesis
	shortVoteGenesis := []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", votingPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", maxDepositPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.denom", "loki"),
		cosmos.NewGenesisKV("app_state.gov.params.quorum", "0.001"),
		cosmos.NewGenesisKV("app_state.gov.params.threshold", "0.001"),
	}

	chain := interchaintest.CreateChainsWithChainSpecs(t, []*interchaintest.ChainSpec{
		{
			Name:      chainName,
			ChainName: chainName,
			Version:   initialVersion,
			ChainConfig: ibc.ChainConfig{
				ModifyGenesis: cosmos.ModifyGenesis(shortVoteGenesis),
				Type:          "cosmos",
				Name:          "odin",
				ChainID:       "odin-mainnet-freya",
				Images: []ibc.DockerImage{
					{
						Repository: upgradeContainerRepo, // FOR LOCAL IMAGE USE: Docker Image Name
						Version:    initialVersion,       // FOR LOCAL IMAGE USE: Docker Image Tag
						UidGid:     "1025:1025",
					},
				},
				Bin:            "odind",
				Bech32Prefix:   "odin",
				Denom:          "loki",
				GasPrices:      "0.00loki",
				GasAdjustment:  1.3,
				TrustingPeriod: "508h",
				NoHostMount:    false,
			},
		},
	})[0].(*cosmos.CosmosChain)

	client, network := interchaintest.DockerSetup(t)
	//t.Cleanup(func() {
	//	_ = client.Close()
	//})

	err := chain.Initialize(ctx, t.Name(), client, network)
	require.NoError(t, err)

	faucetWallet, err := chain.BuildWallet(ctx, interchaintest.FaucetAccountKeyName, "")
	require.NoError(t, err)

	randomWallet1, err := chain.BuildWallet(ctx, "random-wallet-1", "")
	require.NoError(t, err)

	randomWallet2, err := chain.BuildWallet(ctx, "random-wallet-2", "")
	require.NoError(t, err)

	additionalWallets := []ibc.WalletAmount{
		{
			Address: faucetWallet.FormattedAddress(),
			Denom:   chain.Config().Denom,
			Amount:  math.NewInt(1_000_000_000_000),
		},
		{
			Address: randomWallet1.FormattedAddress(),
			Denom:   "udoki",
			Amount:  math.NewInt(1_200_000_000),
		},
		{
			Address: randomWallet2.FormattedAddress(),
			Denom:   "umyrk",
			Amount:  math.NewInt(1_200_000_000),
		},
	}

	t.Log("Starting chain")

	err = chain.Start(t.Name(), ctx, additionalWallets...)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 5, chain)
	require.NoError(t, err)

	t.Log("Funding accounts")

	err = chain.SendFunds(ctx, interchaintest.FaucetAccountKeyName, ibc.WalletAmount{
		Address: randomWallet1.FormattedAddress(),
		Denom:   chain.Config().Denom,
		Amount:  math.NewInt(1_000_000),
	})
	require.NoError(t, err)

	err = chain.SendFunds(ctx, interchaintest.FaucetAccountKeyName, ibc.WalletAmount{
		Address: randomWallet2.FormattedAddress(),
		Denom:   chain.Config().Denom,
		Amount:  math.NewInt(1_000_000),
	})
	require.NoError(t, err)

	err = chain.SendFunds(ctx, "random-wallet-1", ibc.WalletAmount{
		Address: v7_12.AddressWithTokensToBurn,
		Denom:   "udoki",
		Amount:  math.NewInt(200_000_000),
	})
	require.NoError(t, err)

	err = chain.SendFunds(ctx, "random-wallet-2", ibc.WalletAmount{
		Address: v7_12.AddressWithTokensToBurn,
		Denom:   "umyrk",
		Amount:  math.NewInt(200_000_000),
	})
	require.NoError(t, err)

	balance, err := chain.GetBalance(ctx, v7_12.AddressWithTokensToBurn, "udoki")
	require.NoError(t, err)
	require.Equal(t, balance, math.NewInt(200_000_000))

	balance, err = chain.GetBalance(ctx, v7_12.AddressWithTokensToBurn, "umyrk")
	require.NoError(t, err)
	require.Equal(t, balance, math.NewInt(200_000_000))

	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	t.Log("Sending upgrade proposal")

	haltHeight := height + haltHeightDelta

	proposal := cosmos.SoftwareUpgradeProposal{
		Deposit:     "1000000000" + chain.Config().Denom,
		Title:       "Chain Upgrade",
		Name:        upgradeName,
		Description: "First chain software upgrade",
		Height:      haltHeight,
		Info:        "{}",
	}

	proposalTx, err := chain.UpgradeLegacyProposal(ctx, interchaintest.FaucetAccountKeyName, proposal)
	require.NoError(t, err)

	err = chain.VoteOnProposalAllValidators(ctx, proposalTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit vote")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+haltHeightDelta+1, proposalTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height before upgrade")

	timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	// this should timeout due to chain halt at upgrade height.
	_ = testutil.WaitForBlocks(timeoutCtx, int(haltHeight-height)+1, chain)

	// bring down nodes to prepare for upgrade
	err = chain.StopAllNodes(ctx)
	require.NoError(t, err, "error stopping node(s)")

	t.Log("Upgrading chain")

	// upgrade version on all nodes
	chain.UpgradeVersion(ctx, client, upgradeContainerRepo, upgradeVersion)

	// start all nodes back up.
	// validators reach consensus on first block after upgrade height
	// and chain block production resumes.
	err = chain.StartAllNodes(ctx)
	require.NoError(t, err, "error starting upgraded node(s)")
	t.Log("Successful start")
	timeoutCtx, timeoutCtxCancel = context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	err = testutil.WaitForBlocks(timeoutCtx, int(blocksAfterUpgrade), chain)
	require.NoError(t, err, "chain did not produce blocks after upgrade")

	balance, err = chain.GetBalance(ctx, v7_12.AddressWithTokensToBurn, "udoki")
	require.NoError(t, err)
	require.Equal(t, balance, math.NewInt(0))

	balance, err = chain.GetBalance(ctx, v7_12.AddressWithTokensToBurn, "umyrk")
	require.NoError(t, err)
	require.Equal(t, balance, math.NewInt(0))

	govModule, err := chain.GetModuleAddress(ctx, govtypes.ModuleName)
	require.NoError(t, err)

	balance, err = chain.GetBalance(ctx, govModule, "udoki")
	require.NoError(t, err)
	require.Equal(t, balance, math.NewInt(0))

	balance, err = chain.GetBalance(ctx, govModule, "umyrk")
	require.NoError(t, err)
	require.Equal(t, balance, math.NewInt(0))

	balance, err = chain.GetBalance(ctx, randomWallet1.FormattedAddress(), "udoki")
	require.NoError(t, err)
	require.Equal(t, balance, math.NewInt(1_000_000_000))

	balance, err = chain.GetBalance(ctx, randomWallet2.FormattedAddress(), "umyrk")
	require.NoError(t, err)
	require.Equal(t, balance, math.NewInt(1_000_000_000))
}
