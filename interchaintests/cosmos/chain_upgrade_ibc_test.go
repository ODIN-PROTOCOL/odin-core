package cosmos_test

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	interchaintest "github.com/odin-protocol/interchaintest/v7"
	"github.com/odin-protocol/interchaintest/v7/chain/cosmos"
	//	"github.com/odin-protocol/interchaintest/v7/conformance"
	"github.com/odin-protocol/interchaintest/v7/ibc"
	"github.com/odin-protocol/interchaintest/v7/relayer"
	"github.com/odin-protocol/interchaintest/v7/testreporter"
	"github.com/odin-protocol/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	haltHeightDelta    = uint64(10) // will propose upgrade this many blocks in the future
	blocksAfterUpgrade = uint64(10)
	votingPeriod       = "20s"
	maxDepositPeriod   = "10s"
)

func TestOdinUpgradeIBC(t *testing.T) {
	CosmosChainUpgradeIBCTest(t, "odin", "v0.7.9", "gcr.io/odinprotocol/core", "v0.7.10", "v0.7.10")
}

func CosmosChainUpgradeIBCTest(t *testing.T, chainName, initialVersion, upgradeContainerRepo, upgradeVersion string, upgradeName string) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	err := DownloadGenesis()
	require.NoError(t, err)

	// SDK v45 params for Juno genesis
	shortVoteGenesis := []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", votingPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", maxDepositPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.denom", "loki"),
		cosmos.NewGenesisKV("app_state.gov.params.quorum", "0.001"),
		cosmos.NewGenesisKV("app_state.gov.params.threshold", "0.001"),
	}

	chains := interchaintest.CreateChainsWithChainSpecs(t, []*interchaintest.ChainSpec{
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
				//SkipGenTx:      true,
				GenesisPath: "genesis.json",
			},
		},
		{
			Name:      "osmosis",
			ChainName: "osmosis",
			Version:   "v22.0.3",
		},
	})

	client, network := interchaintest.DockerSetup(t)

	chain, counterpartyChain := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	const (
		path        = "ibc-upgrade-test-path"
		relayerName = "relayer"
	)

	// Get a relayer instance
	rf := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
		relayer.StartupFlags("-b", "100"),
	)

	r := rf.Build(t, client, network)

	ic := interchaintest.NewInterchain().
		AddChain(chain).
		AddChain(counterpartyChain).
		AddRelayer(r, relayerName).
		AddLink(interchaintest.InterchainLink{
			Chain1:  chain,
			Chain2:  counterpartyChain,
			Relayer: r,
			Path:    path,
		})

	ctx := context.Background()

	rep := testreporter.NewNopReporter()

	require.NoError(t, ic.Build(ctx, rep.RelayerExecReporter(t), interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation: false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	pubkey, _, err := chain.Validators[0].ExecBin(ctx, "tendermint", "show-validator")
	require.NoError(t, err)

	commandCreateVal := []string{
		"staking", "create-validator",
		"--amount", "5000000000000loki",
		"--commission-max-change-rate", "0.05",
		"--commission-max-rate", "0.10",
		"--commission-rate", "0.05",
		"--min-self-delegation", "1",
		"--pubkey", string(pubkey),
		//"--from", "validator",
	}

	_, err = chain.Validators[0].ExecTx(ctx, "validator", commandCreateVal...)
	require.NoError(t, err)

	var userFunds = sdkmath.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chain)
	chainUser := users[0]

	// test IBC conformance before chain upgrade
	//conformance.TestChainPair(t, ctx, client, network, chain, counterpartyChain, rf, rep, r, path)

	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	haltHeight := height + haltHeightDelta

	proposal := cosmos.SoftwareUpgradeProposal{
		Deposit:     "1000000000" + chain.Config().Denom, // greater than min deposit
		Title:       "Chain Upgrade 1",
		Name:        upgradeName,
		Description: "First chain software upgrade",
		Height:      haltHeight,
		Info:        "{}",
	}

	command := []string{
		"gov", "submit-legacy-proposal",
		"software-upgrade", proposal.Name,
		"--upgrade-height", strconv.FormatUint(proposal.Height, 10),
		"--title", proposal.Title,
		"--description", proposal.Description,
		"--deposit", proposal.Deposit,
		"--no-validate",
	}

	if proposal.Info != "" {
		command = append(command, "--upgrade-info", proposal.Info)
	}

	proposalTxHash, err := chain.FullNodes[0].ExecTx(ctx, chainUser.KeyName(), command...)
	require.NoError(t, err, "error submitting software upgrade proposal tx")

	txResp, err := chain.GetTransaction(proposalTxHash)
	require.NoError(t, err, "failed to receive tx")

	events := txResp.Events
	evtSubmitProp := "submit_proposal"
	proposalID, _ := AttributeValue(events, evtSubmitProp, "proposal_id")

	err = chain.Validators[0].VoteOnProposal(ctx, "validator", proposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit vote")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+haltHeightDelta, proposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height before upgrade")

	timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	// this should timeout due to chain halt at upgrade height.
	_ = testutil.WaitForBlocks(timeoutCtx, int(haltHeight-height)+1, chain)

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height after chain should have halted")

	// make sure that chain is halted
	require.Equal(t, haltHeight, height, "height is not equal to halt height")

	// bring down nodes to prepare for upgrade
	err = chain.StopAllNodes(ctx)
	require.NoError(t, err, "error stopping node(s)")

	// upgrade version on all nodes
	chain.UpgradeVersion(ctx, client, upgradeContainerRepo, upgradeVersion)

	// start all nodes back up.
	// validators reach consensus on first block after upgrade height
	// and chain block production resumes.
	err = chain.StartAllNodes(ctx)
	require.NoError(t, err, "error starting upgraded node(s)")

	timeoutCtx, timeoutCtxCancel = context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	err = testutil.WaitForBlocks(timeoutCtx, int(blocksAfterUpgrade), chain)
	require.NoError(t, err, "chain did not produce blocks after upgrade")

	// test IBC conformance after chain upgrade on same path
	//conformance.TestChainPair(t, ctx, client, network, chain, counterpartyChain, rf, rep, r, path)
}

func AttributeValue(events []abcitypes.Event, eventType, attrKey string) (string, bool) {
	for _, event := range events {
		if event.Type != eventType {
			continue
		}
		for _, attr := range event.Attributes {
			if attr.Key == attrKey {
				return attr.Value, true
			}

			// tendermint < v0.37-alpha returns base64 encoded strings in events.
			key, err := base64.StdEncoding.DecodeString(attr.Key)
			if err != nil {
				continue
			}
			if string(key) == attrKey {
				value, err := base64.StdEncoding.DecodeString(attr.Value)
				if err != nil {
					continue
				}
				return string(value), true
			}
		}
	}
	return "", false
}

func DownloadGenesis() error {

	// Get the data
	resp, err := http.Get("https://storage.googleapis.com/odin-mainnet-freya/genesis-4.1.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create("genesis.json")
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
