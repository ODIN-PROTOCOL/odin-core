package cosmos

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/odin-protocol/interchaintest/v7"
	"github.com/odin-protocol/interchaintest/v7/chain/cosmos"
	"github.com/odin-protocol/interchaintest/v7/conformance"
	"github.com/odin-protocol/interchaintest/v7/ibc"
	"github.com/odin-protocol/interchaintest/v7/relayer"
	"github.com/odin-protocol/interchaintest/v7/relayer/rly"
	"github.com/odin-protocol/interchaintest/v7/testreporter"
	"github.com/odin-protocol/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"golang.org/x/sync/errgroup"
)

const (
	haltHeightDelta    = uint64(20) // will propose upgrade this many blocks in the future
	blocksAfterUpgrade = uint64(10)
	votingPeriod       = "20s"
	maxDepositPeriod   = "10s"
	odinChainID        = "odin-mainnet-freya"
)

func TestOdinFlushIBC(t *testing.T) {
	CosmosChainFlushIBCTest(t, "odin", "v0.6.2", "odinprotocol/core", "v0.7.9", "v7.11", "v0.7.11")
}

func CosmosChainFlushIBCTest(t *testing.T, chainName, initialVersion, upgradeContainerRepo, hardforkVersion, upgradeVersion, upgradeName string) {
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
	}
	_ = shortVoteGenesis

	chains := interchaintest.CreateChainsWithChainSpecs(t, []*interchaintest.ChainSpec{
		{
			Name:      chainName,
			ChainName: chainName,
			Version:   initialVersion,
			ChainConfig: ibc.ChainConfig{
				//ModifyGenesis: cosmos.ModifyGenesis(shortVoteGenesis),
				Type:    "cosmos",
				Name:    "odin",
				ChainID: odinChainID,
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

	//conformance.TestChainPair(t, ctx, client, network, chain, counterpartyChain, rf, rep, r, path)

	// stop relayer
	err := r.StopRelayer(ctx, rep.RelayerExecReporter(t))
	require.NoError(t, err)

	// get channels between test chains
	channels, err := r.GetChannels(ctx, rep.RelayerExecReporter(t), odinChainID)
	require.NoError(t, err)

	// fund users before ibc transfers
	ibcTestUsers := interchaintest.GetAndFundTestUsers(t, ctx, "test-ibc-user", math.NewInt(10_000_000_000), chain, counterpartyChain)

	odinIbcTestUser := ibcTestUsers[0]
	//counterpartyIbcTestUser := ibcTestUsers[1]

	srcTxs := make([]ibc.Tx, len(channels))

	testCoinSrcToDst := ibc.WalletAmount{
		Address: odinIbcTestUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(counterpartyChain.Config().Bech32Prefix),
		Denom:   chain.Config().Denom,
		Amount:  math.NewInt(1_000_000),
	}

	// send ibc transfers with stopped relayer to simulate stacked packages
	for i, channel := range channels {
		srcChannelID := channel.ChannelID
		srcTxs[i], err = chain.SendIBCTransfer(ctx, srcChannelID, odinIbcTestUser.KeyName(), testCoinSrcToDst, ibc.TransferOptions{Timeout: nil})
		require.NoError(t, err)

		err = testutil.WaitForBlocks(ctx, 1, chain)
		require.NoError(t, err)
	}

	for _, srcTx := range srcTxs {
		require.NoError(t, srcTx.Validate(), "source ibc transfer tx is invalid")
	}

	balance, err := chain.GetBalance(ctx, odinIbcTestUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(chain.Config().Bech32Prefix), "loki")
	require.NoError(t, err)

	t.Log("balance after first ibc transfer attempt", balance)

	height, err := chain.Height(ctx)
	require.NoError(t, err)

	for _, node := range chain.Nodes() {
		err = node.StopContainer(ctx)
		require.NoError(t, err)
	}

	stateExported, err := chain.ExportState(ctx, int64(height))
	require.NoError(t, err)

	stateExported = strings.ReplaceAll(stateExported, "\"max_validators\":100,", "\"max_validators\":100,\"min_commission_rate\": \"0.000000000000000000\",")

	stateExported = strings.ReplaceAll(
		stateExported,
		"{\"deposit_params\":{\"max_deposit_period\":\"172800s\",\"min_deposit\":[{\"amount\":\"1000000000\",\"denom\":\"loki\"}]},\"deposits\":[],\"proposals\":[],\"starting_proposal_id\":\"1\",\"tally_params\":{\"quorum\":\"0.334000000000000000\",\"threshold\":\"0.500000000000000000\",\"veto_threshold\":\"0.334000000000000000\"},\"votes\":[],\"voting_params\":{\"voting_period\":\"172800s\"}},",
		"{\"starting_proposal_id\":\"1\",\"deposits\":[],\"votes\":[],\"proposals\":[],\"deposit_params\": null,\"voting_params\": null,\"tally_params\": null,\"params\": {\"min_deposit\": [{\"denom\": \"loki\",\"amount\": \"1000000000\"}],\"max_deposit_period\": \"20s\",\"voting_period\": \"20s\",\"quorum\": \"0.334000000000000000\",\"threshold\": \"0.500000000000000000\",\"veto_threshold\": \"0.334000000000000000\",\"min_initial_deposit_ratio\": \"0.000000000000000000\",\"burn_vote_quorum\": false,\"burn_proposal_deposit_prevote\": false,\"burn_vote_veto\": true}},",
	)

	stateExported = strings.ReplaceAll(
		stateExported,
		"{\"data_sources\":[],\"module_coins_account\":\"\",\"oracle_pool\":{\"data_providers_pool\":[]},\"oracle_scripts\":[],\"params\":{\"base_owasm_gas\":\"150000\",\"data_provider_reward_per_byte\":[{\"amount\":\"1000000\",\"denom\":\"loki\"},{\"amount\":\"1000000\",\"denom\":\"minigeo\"}],\"data_provider_reward_threshold\":{\"amount\":[{\"amount\":\"200000000000\",\"denom\":\"loki\"},{\"amount\":\"200000000000\",\"denom\":\"minigeo\"}],\"blocks\":\"28820\"},\"data_requester_fee_denoms\":[\"loki\",\"minigeo\"],\"expiration_block_count\":\"100\",\"inactive_penalty_duration\":\"600000000000\",\"max_ask_count\":\"16\",\"max_calldata_size\":\"1024\",\"max_data_size\":\"1024\",\"max_raw_request_count\":\"12\",\"oracle_reward_percentage\":\"70\",\"per_validator_request_gas\":\"30000\",\"reward_decreasing_fraction\":\"0.050000000000000000\",\"sampling_try_count\":\"3\"}},",
		"{\"params\":{\"max_raw_request_count\":\"16\",\"max_ask_count\":\"16\",\"max_calldata_size\":\"256\",\"max_report_data_size\":\"512\",\"expiration_block_count\": \"100\",\"base_owasm_gas\": \"50000\",\"per_validator_request_gas\":\"0\",\"sampling_try_count\": \"3\",\"oracle_reward_percentage\":\"70\",\"inactive_penalty_duration\": \"600000000000\",\"ibc_request_enabled\":true,\"data_provider_reward_per_byte\":[],\"data_provider_reward_threshold\":{\"amount\":[{\"denom\":\"loki\",\"amount\":\"200000000000\"},{\"denom\":\"minigeo\",\"amount\":\"200000000000\"}],\"blocks\":\"28820\"},\"reward_decreasing_fraction\":\"0\",\"data_requester_fee_denoms\":[]},\"data_sources\":[],\"oracle_scripts\":[],\"oracle_pool\":{\"data_providers_pool\":[]},\"module_coins_account\":\"odin1lqf6hm3nfunmhppmjhgrme9jp9d8vle90hjy5m\"},",
	)

	stateExported = strings.ReplaceAll(
		stateExported,
		"\"timestamp\":\"2024",
		"\"timestamp\":\"2023",
	)

	for _, node := range chain.Nodes() {
		//err = node.StopContainer(ctx)
		//require.NoError(t, err)

		err = node.RemoveContainer(ctx)
		require.NoError(t, err)

		_, _, err = node.ExecBin(ctx, "tendermint", "unsafe-reset-all")
		require.NoError(t, err)

		err = node.OverwriteGenesisFile(ctx, []byte(stateExported))
		require.NoError(t, err)

		// TODO: change to chain.UpgradeVersion
		node.Image = ibc.DockerImage{
			Repository: upgradeContainerRepo,
			Version:    hardforkVersion,
			UidGid:     "1026:1026",
		}

		err = node.CreateNodeContainer(ctx)
		require.NoError(t, err)
	}

	eg, egCtx := errgroup.WithContext(ctx)
	for _, node := range chain.Nodes() {
		node := node
		eg.Go(func() error {
			err = node.StartContainer(egCtx)
			//require.NoError(t, err)
			return err
		})
	}

	err = eg.Wait()
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 5, chain)
	require.NoError(t, err)

	// txs to make keys migration, ignore errors due to unavailability of unmarshalling result
	chain.SendFunds(ctx, "faucet", ibc.WalletAmount{
		Address: odinIbcTestUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(chain.Config().Bech32Prefix),
		Denom:   "loki",
		Amount:  math.NewInt(1),
	})

	chain.SendFunds(ctx, odinIbcTestUser.KeyName(), ibc.WalletAmount{
		Address: "",
		Denom:   "loki",
		Amount:  math.NewInt(1),
	})

	for _, validator := range chain.Validators {
		validator.SendFunds(ctx, "validator", ibc.WalletAmount{
			Address: "",
			Denom:   "loki",
			Amount:  math.NewInt(1),
		})
	}

	// TODO: use it later
	//send ibc transfers with stopped relayer to simulate stacked packages
	//for _, channel := range channels {
	//	srcChannelID := channel.ChannelID
	//	srcTx, err := chain.SendIBCTransfer(ctx, srcChannelID, odinIbcTestUser.KeyName(), testCoinSrcToDst, ibc.TransferOptions{Timeout: nil})
	//	require.NoError(t, err)
	//
	//	srcTxs = append(srcTxs, srcTx)
	//
	//	err = testutil.WaitForBlocks(ctx, 1, chain)
	//	require.NoError(t, err)
	//}
	//
	//for _, srcTx := range srcTxs {
	//	require.NoError(t, srcTx.Validate(), "source ibc transfer tx is invalid")
	//}

	balance, err = chain.GetBalance(ctx, odinIbcTestUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(chain.Config().Bech32Prefix), "loki")
	require.NoError(t, err)

	t.Log("balance after second ibc transfer attempt", balance)

	height, err = chain.Height(ctx)
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

	//command := []string{
	//	"gov", "submit-legacy-proposal",
	//	"software-upgrade", proposal.Name,
	//	"--upgrade-height", strconv.FormatUint(proposal.Height, 10),
	//	"--title", proposal.Title,
	//	"--description", proposal.Description,
	//	"--deposit", proposal.Deposit,
	//	"--no-validate",
	//}
	//
	//if proposal.Info != "" {
	//	command = append(command, "--upgrade-info", proposal.Info)
	//}

	var userFunds = math.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "test-ibc-user1", userFunds, chain)
	chainUser := users[0]

	t.Log("Halt height:", haltHeight)

	proposalTx, err := chain.UpgradeLegacyProposal(ctx, chainUser.KeyName(), proposal)
	require.NoError(t, err, "error submitting software upgrade proposal tx")

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

	balance, err = chain.GetBalance(ctx, odinIbcTestUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(chain.Config().Bech32Prefix), "loki")
	require.NoError(t, err)

	t.Log("balance after chain update", balance)

	//emptyStr := "\"\""
	//
	//err = r.UpdatePath(ctx, rep.RelayerExecReporter(t), path, ibc.PathUpdateOptions{
	//	SrcClientID: &emptyStr,
	//	SrcConnID:   &emptyStr,
	//})
	//require.NoError(t, err)
	//

	createClientCmd := []string{
		"rly", "tx", "client", odinChainID, counterpartyChain.Config().ChainID, path, "--client-tp", "0",
		"--home", r.(*rly.CosmosRelayer).HomeDir(), "--override",
	}

	createClientRes := r.Exec(ctx, rep.RelayerExecReporter(t), createClientCmd, nil)
	require.NoError(t, createClientRes.Err)

	//err = r.CreateClient(ctx, rep.RelayerExecReporter(t), odinChainID, counterpartyChain.Config().ChainID, path, ibc.CreateClientOptions{TrustingPeriod: "0"})
	//require.NoError(t, err)
	//
	//r.StartRelayer(ctx, rep.RelayerExecReporter(t), path)

	conn := "connection-0"
	oldClient := "07-tendermint-0"

	err = r.UpdatePath(ctx, rep.RelayerExecReporter(t), path, ibc.PathUpdateOptions{
		SrcClientID: &oldClient,
		SrcConnID:   &conn,
	})

	//proposalContent := clienttypes.NewClientUpdateProposal("Update client", "Update client", "07-tendermint-0", "07-tendermint-1")

	command := []string{
		"gov", "submit-legacy-proposal",
		"update-client", "07-tendermint-0", "07-tendermint-1",
		"--title", "Update client",
		"--description", "Update client",
		"--deposit", proposal.Deposit,
	}

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	proposalTxHash, err := chain.FullNodes[0].ExecTx(ctx, chainUser.KeyName(), command...)
	require.NoError(t, err)

	txResp, err := chain.GetTransaction(proposalTxHash)
	require.NoError(t, err, "failed to receive tx")

	events := txResp.Events
	evtSubmitProp := "submit_proposal"
	proposalID, _ := AttributeValue(events, evtSubmitProp, "proposal_id")

	err = chain.VoteOnProposalAllValidators(ctx, proposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit vote")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+haltHeightDelta, proposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	err = r.StartRelayer(ctx, rep.RelayerExecReporter(t), path)
	require.NoError(t, err)

	conformance.TestChainPair(t, ctx, client, network, chain, counterpartyChain, rf, rep, r, path)

	//err = r.CreateClient(ctx, rep.RelayerExecReporter(t), odinChainID, counterpartyChain.Config().ChainID, path, ibc.CreateClientOptions{TrustingPeriod: "0"})
	//require.NoError(t, err)

	//proposalContent := clienttypes.NewClientUpdateProposal("Update client", "Update client", "07-tendermint-0", "07-tendermint-1")

	//command := []string{
	//	"gov", "submit-legacy-proposal",
	//	"update-client", "07-tendermint-0", "07-tendermint-1",
	//	"--title", "Update client",
	//	"--description", "Update client",
	//	"--deposit", proposal.Deposit,
	//}
	//
	//proposalTxHash, err := chain.FullNodes[0].ExecTx(ctx, odinIbcTestUser.KeyName(), command...)
	//require.NoError(t, err)
	//
	//txResp, err := chain.GetTransaction(proposalTxHash)
	//require.NoError(t, err, "failed to receive tx")
	//
	//events := txResp.Events
	//evtSubmitProp := "submit_proposal"
	//proposalID, _ := AttributeValue(events, evtSubmitProp, "proposal_id")
	//
	//height, err = chain.Height(ctx)
	//require.NoError(t, err, "error fetching height before submit upgrade proposal")
	//
	//_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+haltHeightDelta, proposalID, cosmos.ProposalStatusPassed)
	//require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")
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
