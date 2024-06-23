package cosmos

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"

	"cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/conformance"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/relayer"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	haltHeightDelta    = uint64(10) // will propose upgrade this many blocks in the future
	blocksAfterUpgrade = uint64(10)
	votingPeriod       = "10s"
	maxDepositPeriod   = "10s"
	odinChainID        = "odin-mainnet-freya"
)

func TestOdinFlushIBC(t *testing.T) {
	CosmosChainFlushIBCTest(t, "odin", "v0.7.12", "odinprotocol/core", "v0.8.1", "v0.8.1")
}

func CosmosChainFlushIBCTest(t *testing.T, chainName, initialVersion, upgradeContainerRepo, upgradeVersion, upgradeName string) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	// SDK v45 params for Juno genesis
	shortVoteGenesis := []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", votingPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", maxDepositPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.denom", "loki"),
		cosmos.NewGenesisKV("app_state.gov.params.quorum", "0.001"),
		cosmos.NewGenesisKV("app_state.gov.params.threshold", "0.001"),
		cosmos.NewGenesisKV("app_state.mint.params.mint_air", true),
	}

	one := int(1)

	chains := interchaintest.CreateChainsWithChainSpecs(t, []*interchaintest.ChainSpec{
		{
			Name:      chainName,
			ChainName: chainName,
			Version:   initialVersion,
			ChainConfig: ibc.ChainConfig{
				ModifyGenesis: cosmos.ModifyGenesis(shortVoteGenesis),
				Type:          "cosmos",
				Name:          "odin",
				ChainID:       odinChainID,
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
			NumValidators: &one,
			NumFullNodes:  &one,
		},
		{
			Name:          "osmosis",
			ChainName:     "osmosis",
			Version:       "v22.0.3",
			NumValidators: &one,
			NumFullNodes:  &one,
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

	err := r.StartRelayer(ctx, rep.RelayerExecReporter(t), path)
	require.NoError(t, err)

	ibcTestUsers := interchaintest.GetAndFundTestUsers(t, ctx, "test-ibc-user", math.NewInt(10_000_000_000), chain, counterpartyChain)

	IbcTransfersTest(t, ctx, r, rep, chain, counterpartyChain, ibcTestUsers[0], ibcTestUsers[1])
	conformance.TestChainPair(t, ctx, client, network, chain, counterpartyChain, rf, rep, r, path)

	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	t.Log("Sending upgrade proposal")

	haltHeight := height + int64(haltHeightDelta)

	proposal := cosmos.SoftwareUpgradeProposal{
		Deposit:     "1000000000" + chain.Config().Denom,
		Title:       "Chain Upgrade",
		Name:        upgradeName,
		Description: "First chain software upgrade",
		Height:      haltHeight,
		Info:        "{}",
	}

	upgradeTx, err := UpgradeProposal(chain, ctx, interchaintest.FaucetAccountKeyName, proposal)
	require.NoError(t, err)

	propId, err := strconv.ParseUint(upgradeTx.ProposalID, 10, 64)
	require.NoError(t, err, "failed to convert proposal ID to uint64")

	err = chain.VoteOnProposalAllValidators(ctx, upgradeTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit vote")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+int64(haltHeightDelta)+1, propId, govv1beta1.StatusPassed)
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

	t.Log("Upgrade successful")

	IbcTransfersTest(t, ctx, r, rep, chain, counterpartyChain, ibcTestUsers[0], ibcTestUsers[1])
	conformance.TestChainPair(t, ctx, client, network, chain, counterpartyChain, rf, rep, r, path)

}

func UpgradeProposal(chain *cosmos.CosmosChain, ctx context.Context, keyName string, prop cosmos.SoftwareUpgradeProposal) (tx cosmos.TxProposal, err error) {
	command := []string{
		"gov", "submit-legacy-proposal",
		"software-upgrade", prop.Name,
		"--upgrade-height", strconv.FormatInt(prop.Height, 10),
		"--title", prop.Title,
		"--description", prop.Description,
		"--deposit", prop.Deposit,
		"--no-validate",
	}

	if prop.Info != "" {
		command = append(command, "--upgrade-info", prop.Info)
	}

	fullNode := chain.GetNode()
	txHash, err := fullNode.ExecTx(ctx, keyName, command...)
	if err != nil {
		return cosmos.TxProposal{}, err
	}

	txResp, err := chain.GetTransaction(txHash)
	if err != nil {
		return cosmos.TxProposal{}, fmt.Errorf("failed to get transaction %s: %w", txHash, err)
	}
	tx.Height = txResp.Height
	tx.TxHash = txHash
	// In cosmos, user is charged for entire gas requested, not the actual gas used.
	tx.GasSpent = txResp.GasWanted
	events := txResp.Events

	tx.DepositAmount, _ = AttributeValue(events, "proposal_deposit", "amount")

	evtSubmitProp := "submit_proposal"
	tx.ProposalID, _ = AttributeValue(events, evtSubmitProp, "proposal_id")
	tx.ProposalType, _ = AttributeValue(events, evtSubmitProp, "proposal_type")

	return tx, nil
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

func IbcTransfersTest(t *testing.T, ctx context.Context, r ibc.Relayer, rep *testreporter.Reporter, srcChain, dstChain ibc.Chain, srcUser, dstUser ibc.Wallet) {
	channels, err := r.GetChannels(ctx, rep.RelayerExecReporter(t), odinChainID)
	require.NoError(t, err)

	srcTxs := make([]ibc.Tx, len(channels))
	dstTxs := make([]ibc.Tx, len(channels))

	srcDenoms := make([]string, len(channels))
	dstDenoms := make([]string, len(channels))
	srcToDstBalances := make([]math.Int, len(channels))
	dstToSrcBalances := make([]math.Int, len(channels))

	srcInitialBalance, err := srcChain.GetBalance(ctx, srcUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(srcChain.Config().Bech32Prefix), "loki")
	require.NoError(t, err)

	dstInitialBalance, err := dstChain.GetBalance(ctx, dstUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(dstChain.Config().Bech32Prefix), "uosmo")
	require.NoError(t, err)

	testCoinSrcToDst := ibc.WalletAmount{
		Address: srcUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(dstChain.Config().Bech32Prefix),
		Denom:   srcChain.Config().Denom,
		Amount:  math.NewInt(1_000_000),
	}

	testCoinDstToSrc := ibc.WalletAmount{
		Address: dstUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(srcChain.Config().Bech32Prefix),
		Denom:   dstChain.Config().Denom,
		Amount:  math.NewInt(1_000_000),
	}

	// START odin->osmosis
	for i, channel := range channels {
		srcDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom(channel.Counterparty.PortID, channel.Counterparty.ChannelID, "loki"))
		dstIbcDenom := srcDenomTrace.IBCDenom()
		srcDenoms[i] = dstIbcDenom

		srcToDstBalances[i], err = dstChain.GetBalance(ctx, testCoinSrcToDst.Address, dstIbcDenom)
		require.NoError(t, err)

		srcChannelID := channel.ChannelID
		srcTxs[i], err = srcChain.SendIBCTransfer(ctx, srcChannelID, srcUser.KeyName(), testCoinSrcToDst, ibc.TransferOptions{Timeout: nil})
		require.NoError(t, err)

		err = testutil.WaitForBlocks(ctx, 5, srcChain)
		require.NoError(t, err)
	}

	spend := math.NewInt(0)

	for _, srcTx := range srcTxs {
		require.NoError(t, srcTx.Validate(), "source ibc transfer tx is invalid")
		fee := srcChain.GetGasFeesInNativeDenom(srcTx.GasSpent)
		spend = spend.Add(math.NewInt(fee)).Add(math.NewInt(1_000_000))
	}

	balance, err := srcChain.GetBalance(ctx, srcUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(srcChain.Config().Bech32Prefix), "loki")
	require.NoError(t, err)

	require.Equal(t, balance, srcInitialBalance.Sub(spend))

	for i, _ := range channels {
		balance, err := dstChain.GetBalance(ctx, testCoinSrcToDst.Address, srcDenoms[i])
		require.NoError(t, err)
		require.Equal(t, balance, srcToDstBalances[i].Add(math.NewInt(1_000_000)))
	}
	// END odin->osmosis

	// START osmosis->odin
	for i, channel := range channels {
		dstDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom(channel.PortID, channel.ChannelID, "uosmo"))
		srcIbcDenom := dstDenomTrace.IBCDenom()
		dstDenoms[i] = srcIbcDenom

		dstToSrcBalances[i], err = srcChain.GetBalance(ctx, testCoinDstToSrc.Address, srcIbcDenom)
		require.NoError(t, err)

		srcChannelID := channel.Counterparty.ChannelID
		dstTxs[i], err = dstChain.SendIBCTransfer(ctx, srcChannelID, dstUser.KeyName(), testCoinDstToSrc, ibc.TransferOptions{Timeout: nil})
		require.NoError(t, err)

		err = testutil.WaitForBlocks(ctx, 1, dstChain)
		require.NoError(t, err)
	}

	spend = math.NewInt(0)

	for _, dstTx := range dstTxs {
		require.NoError(t, dstTx.Validate(), "source ibc transfer tx is invalid")
		fee := dstChain.GetGasFeesInNativeDenom(dstTx.GasSpent)
		spend = spend.Add(math.NewInt(fee)).Add(math.NewInt(1_000_000))
	}

	balance, err = dstChain.GetBalance(ctx, dstUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(dstChain.Config().Bech32Prefix), "uosmo")
	require.NoError(t, err)

	require.Equal(t, balance, dstInitialBalance.Sub(spend))

	for i, _ := range channels {
		balance, err := srcChain.GetBalance(ctx, testCoinDstToSrc.Address, dstDenoms[i])
		require.NoError(t, err)
		require.Equal(t, balance, dstToSrcBalances[i].Add(math.NewInt(1_000_000)))
	}
}
