package v7_10_test

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	interchaintest "github.com/odin-protocol/interchaintest/v7"
	"github.com/odin-protocol/interchaintest/v7/chain/cosmos"
	"github.com/odin-protocol/interchaintest/v7/ibc"
	"github.com/odin-protocol/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
)

const (
	haltHeightDelta    = uint64(10) // will propose upgrade this many blocks in the future
	blocksAfterUpgrade = uint64(10)
	votingPeriod       = "10s"
	maxDepositPeriod   = "10s"
)

func TestOdinUpgrade(t *testing.T) {
	CosmosChainUpgradeTest(t, "odin", "v0.7.9", "gcr.io/odinprotocol/core", "v0.7.10", "v7_10")
}

func CosmosChainUpgradeTest(t *testing.T, chainName, initialVersion, upgradeContainerRepo, upgradeVersion string, upgradeName string) {
	// Performs upgrade test

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
						Repository: upgradeContainerRepo, 
						Version:    initialVersion,
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
				GenesisPath: "./genesis.json",
			},
		},
	})

	client, _ := interchaintest.DockerSetup(t)

	chain := chains[0].(*cosmos.CosmosChain)
	ctx := context.Background()

	log.Printf("new chain validators: %v", chain.NumValidators)

	var userFunds = sdkmath.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chain)
	chainUser := users[0]

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

	//upgradeTx, err := chain.UpgradeProposal(ctx, chainUser.KeyName(), proposal)
	//require.NoError(t, err, "error submitting software upgrade proposal tx")

	//for i := 0; i < 5; i++ {
	//	res, _, err := chain.Validators[0].ExecQuery(ctx, "gov", "proposals")
	//	fmt.Println(res)
	//	proposal, err := chain.QueryProposal(ctx, proposalID)
	//	fmt.Println(proposal)
	//	fmt.Println(err)
	//	time.Sleep(5 * time.Second)
	//}

	err = chain.VoteOnProposalAllValidators(ctx, proposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

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
}

//
//func StartWithGenesisFile(
//	ctx context.Context,
//	testName string,
//	client *client.Client,
//	network string,
//	genesisFilePath string,
//	c *cosmos.CosmosChain,
//) error {
//	genBz, err := os.ReadFile(genesisFilePath)
//	if err != nil {
//		return fmt.Errorf("failed to read genesis file: %w", err)
//	}
//
//	chainCfg := c.Config()
//
//	var genesisFile cosmos.GenesisFile
//	if err := json.Unmarshal(genBz, &genesisFile); err != nil {
//		return err
//	}
//
//	genesisValidators := genesisFile.Validators
//	totalPower := int64(0)
//
//	validatorsWithPower := make([]cosmos.ValidatorWithIntPower, 0)
//
//	for _, genesisValidator := range genesisValidators {
//		power, err := strconv.ParseInt(genesisValidator.Power, 10, 64)
//		if err != nil {
//			return err
//		}
//		totalPower += power
//		validatorsWithPower = append(validatorsWithPower, cosmos.ValidatorWithIntPower{
//			Address:      genesisValidator.Address,
//			Power:        power,
//			PubKeyBase64: genesisValidator.PubKey.Value,
//		})
//	}
//
//	sort.Slice(validatorsWithPower, func(i, j int) bool {
//		return validatorsWithPower[i].Power > validatorsWithPower[j].Power
//	})
//
//	var eg errgroup.Group
//	var mu sync.Mutex
//	genBzReplace := func(find, replace []byte) {
//		mu.Lock()
//		defer mu.Unlock()
//		genBz = bytes.ReplaceAll(genBz, find, replace)
//	}
//
//	twoThirdsConsensus := int64(math.Ceil(float64(totalPower) * 2 / 3))
//	totalConsensus := int64(0)
//
//	var activeVals []cosmos.ValidatorWithIntPower
//	for _, validator := range validatorsWithPower {
//		activeVals = append(activeVals, validator)
//
//		totalConsensus += validator.Power
//
//		if totalConsensus > twoThirdsConsensus {
//			break
//		}
//	}
//
//	c.numValidators = len(activeVals)
//
//	if err := c.initializeChainNodes(ctx, testName, client, network); err != nil {
//		return err
//	}
//
//	if err := c.prepNodes(ctx, true, nil, types.Coin{}); err != nil {
//		return err
//	}
//
//	if c.cfg.PreGenesis != nil {
//		err := c.cfg.PreGenesis(chainCfg)
//		if err != nil {
//			return err
//		}
//	}
//
//	for i, validator := range activeVals {
//		v := c.Validators[i]
//		validator := validator
//		eg.Go(func() error {
//			testNodePubKeyJsonBytes, err := v.ReadFile(ctx, "config/priv_validator_key.json")
//			if err != nil {
//				return fmt.Errorf("failed to read priv_validator_key.json: %w", err)
//			}
//
//			var testNodePrivValFile PrivValidatorKeyFile
//			if err := json.Unmarshal(testNodePubKeyJsonBytes, &testNodePrivValFile); err != nil {
//				return fmt.Errorf("failed to unmarshal priv_validator_key.json: %w", err)
//			}
//
//			// modify genesis file overwriting validators address with the one generated for this test node
//			genBzReplace([]byte(validator.Address), []byte(testNodePrivValFile.Address))
//
//			// modify genesis file overwriting validators base64 pub_key.value with the one generated for this test node
//			genBzReplace([]byte(validator.PubKeyBase64), []byte(testNodePrivValFile.PubKey.Value))
//
//			existingValAddressBytes, err := hex.DecodeString(validator.Address)
//			if err != nil {
//				return err
//			}
//
//			testNodeAddressBytes, err := hex.DecodeString(testNodePrivValFile.Address)
//			if err != nil {
//				return err
//			}
//
//			valConsPrefix := fmt.Sprintf("%svalcons", chainCfg.Bech32Prefix)
//
//			existingValBech32ValConsAddress, err := bech32.ConvertAndEncode(valConsPrefix, existingValAddressBytes)
//			if err != nil {
//				return err
//			}
//
//			testNodeBech32ValConsAddress, err := bech32.ConvertAndEncode(valConsPrefix, testNodeAddressBytes)
//			if err != nil {
//				return err
//			}
//
//			genBzReplace([]byte(existingValBech32ValConsAddress), []byte(testNodeBech32ValConsAddress))
//
//			return nil
//		})
//	}
//
//	if err := eg.Wait(); err != nil {
//		return err
//	}
//
//	return c.startWithFinalGenesis(ctx, genBz)
//}

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
