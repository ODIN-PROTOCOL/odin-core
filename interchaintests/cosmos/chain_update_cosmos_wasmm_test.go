package cosmos

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos/wasm"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

//const (
//	haltHeightDelta    = uint64(10) // will propose upgrade this many blocks in the future
//	blocksAfterUpgrade = uint64(10)
//	votingPeriod       = "10s"
//	maxDepositPeriod   = "10s"
//	odinChainID        = "odin-mainnet-freya"
//)

func TestOdinWasmIBC(t *testing.T) {
	CosmosChainWasmIBCTest(t, "odin", "v0.8.4", "odinprotocol/core", "v0.9.3", "v0.9.3")
}

func CosmosChainWasmIBCTest(t *testing.T, chainName, initialVersion, upgradeContainerRepo, upgradeVersion, upgradeName string) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

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
				EncodingConfig: wasm.WasmEncoding(),
			},
			NumValidators: &one,
			NumFullNodes:  &one,
		},
		{
			Name:          "juno",
			ChainName:     "juno",
			Version:       "latest",
			NumValidators: &one,
			NumFullNodes:  &one,
			ChainConfig: ibc.ChainConfig{
				EncodingConfig: wasm.WasmEncoding(),
				GasPrices:      "0.00ujuno",
			},
		},
	})

	chain, counterpartyChain := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Relayer Factory
	client, network := interchaintest.DockerSetup(t)
	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(
		t, client, network)

	// Prep Interchain
	const ibcPath = "wasmpath"
	ic := interchaintest.NewInterchain().
		AddChain(chain).
		AddChain(counterpartyChain).
		AddRelayer(r, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  chain,
			Chain2:  counterpartyChain,
			Relayer: r,
			Path:    ibcPath,
		})

	ctx := context.Background()

	rep := testreporter.NewNopReporter()

	require.NoError(t, ic.Build(ctx, rep.RelayerExecReporter(t), interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation: true,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	t.Log("Sending upgrade proposal")

	haltHeight := height + int64(haltHeightDelta)

	govAddr, err := chain.AuthQueryModuleAddress(ctx, "gov")
	require.NoError(t, err)

	prop, err := chain.BuildProposal([]cosmos.ProtoMessage{&upgradetypes.MsgSoftwareUpgrade{
		Authority: govAddr,
		Plan: upgradetypes.Plan{
			Name:   "v0.9.3",
			Height: haltHeight,
			Info:   "",
		},
	}}, "Chain upgrade", "sum", "", "1000000000"+chain.Config().Denom, "", false)
	require.NoError(t, err)

	upgradeTx, err := chain.SubmitProposal(ctx, interchaintest.FaucetAccountKeyName, prop)
	require.NoError(t, err)

	propId, err := strconv.ParseUint(upgradeTx.ProposalID, 10, 64)
	require.NoError(t, err, "failed to convert proposal ID to uint64")

	err = chain.VoteOnProposalAllValidators(ctx, propId, cosmos.ProposalVoteYes)
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

	err = r.GeneratePath(ctx, rep.RelayerExecReporter(t), chain.Config().ChainID, counterpartyChain.Config().ChainID, ibcPath)
	require.NoError(t, err, "error generating path")

	err = r.LinkPath(ctx, rep.RelayerExecReporter(t), ibcPath, ibc.DefaultChannelOpts(), ibc.DefaultClientOpts())
	require.NoError(t, err, "error linking path")

	// Create and Fund User Wallets
	initBal := math.NewInt(100_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", initBal, chain, counterpartyChain)
	odinUser := users[0]
	junoUser := users[1]

	err = testutil.WaitForBlocks(ctx, 2, chain, counterpartyChain)
	require.NoError(t, err)

	odinUserBalInitial, err := chain.GetBalance(ctx, odinUser.FormattedAddress(), chain.Config().Denom)
	require.NoError(t, err)
	require.True(t, odinUserBalInitial.Equal(initBal))

	junoUserBalInitial, err := counterpartyChain.GetBalance(ctx, junoUser.FormattedAddress(), counterpartyChain.Config().Denom)
	require.NoError(t, err)
	require.True(t, junoUserBalInitial.Equal(initBal))

	err = r.StartRelayer(ctx, rep.RelayerExecReporter(t), ibcPath)
	require.NoError(t, err, "error starting relayer")

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, rep.RelayerExecReporter(t))
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	// Store ibc_reflect_send.wasm contract on juno1
	juno1ContractCodeId, err := chain.StoreContract(
		ctx, odinUser.KeyName(), "sample_contracts/cw_ibc_example.wasm")
	require.NoError(t, err)

	// Instantiate the contract on juno1
	juno1ContractAddr, err := chain.InstantiateContract(
		ctx, odinUser.KeyName(), juno1ContractCodeId, "{}", true)
	require.NoError(t, err)

	// Store ibc_reflect_send.wasm on juno2
	juno2ContractCodeId, err := counterpartyChain.StoreContract(
		ctx, junoUser.KeyName(), "sample_contracts/cw_ibc_example.wasm")
	require.NoError(t, err)

	// Instantiate contract on juno2
	juno2ContractAddr, err := counterpartyChain.InstantiateContract(
		ctx, junoUser.KeyName(), juno2ContractCodeId, "{}", true)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 1, chain, counterpartyChain)
	require.NoError(t, err)

	// Query the reflect sender contract on Juno1 for it's port id
	chainContractInfo, err := chain.QueryContractInfo(ctx, juno1ContractAddr)
	require.NoError(t, err)
	chainContractPortId := chainContractInfo.ContractInfo.IbcPortID

	// Query the reflect contract on Juno2 for it's port id
	counterpartyContractInfo, err := counterpartyChain.QueryContractInfo(ctx, juno2ContractAddr)
	require.NoError(t, err)
	counterpartyContractPortId := counterpartyContractInfo.ContractInfo.IbcPortID

	// Create channel between Juno1 and Juno2
	err = r.CreateChannel(ctx, rep.RelayerExecReporter(t), ibcPath, ibc.CreateChannelOptions{
		SourcePortName: chainContractPortId,
		DestPortName:   counterpartyContractPortId,
		Order:          ibc.Unordered,
		Version:        "counter-1",
	})
	require.NoError(t, err)

	// Wait for the channel to get set up and whoami message to exchange
	err = testutil.WaitForBlocks(ctx, 10, chain, counterpartyChain)
	require.NoError(t, err)

	// Get contract channel
	chainChannelInfo, err := r.GetChannels(ctx, rep.RelayerExecReporter(t), chain.Config().ChainID)
	require.NoError(t, err)
	chainChannelID := chainChannelInfo[len(chainChannelInfo)-1].ChannelID

	// Get contract channel
	counterpartyChannelInfo, err := r.GetChannels(ctx, rep.RelayerExecReporter(t), counterpartyChain.Config().ChainID)
	require.NoError(t, err)
	counterpartyChannelID := counterpartyChannelInfo[len(counterpartyChannelInfo)-1].ChannelID

	// Prepare the query and execute messages to interact with the contracts
	queryChainCountMsg := fmt.Sprintf(`{"get_count":{"channel":"%s"}}`, chainChannelID)
	queryCounterpartyCountMsg := fmt.Sprintf(`{"get_count":{"channel":"%s"}}`, counterpartyChannelID)
	chainIncrementMsg := fmt.Sprintf(`{"increment": {"channel":"%s"}}`, chainChannelID)
	counterpartyIncrementMsg := fmt.Sprintf(`{"increment": {"channel":"%s"}}`, counterpartyChannelID)

	_, err = chain.Height(ctx)
	require.NoError(t, err)

	// Query the count of the contract on juno1- should be 0 as no packets have been sent through
	var chainInitialCountResponse CwIbcCountResponse
	err = chain.QueryContract(ctx, juno1ContractAddr, queryChainCountMsg, &chainInitialCountResponse)
	require.NoError(t, err)
	require.Equal(t, 0, chainInitialCountResponse.Data.Count)

	// Query the count of the contract on juno1- should be 0 as no packets have been sent through
	var counterpartyInitialCountResponse CwIbcCountResponse
	err = counterpartyChain.QueryContract(ctx, juno2ContractAddr, queryCounterpartyCountMsg, &counterpartyInitialCountResponse)
	require.NoError(t, err)
	require.Equal(t, 0, counterpartyInitialCountResponse.Data.Count)

	// Send packet from juno1 to juno2 and increment the juno2 contract count
	juno1Increment, err := chain.ExecuteContract(ctx, odinUser.KeyName(), juno1ContractAddr, chainIncrementMsg)
	require.NoError(t, err)
	// Check if the transaction was successful
	require.Equal(t, uint32(0), juno1Increment.Code)

	// Wait for the ibc packet to be delivered
	err = testutil.WaitForBlocks(ctx, 2, chain, counterpartyChain)
	require.NoError(t, err)

	// Query the count of the contract on juno2- should be 1 as a single packet has been sent through
	var juno2IncrementedCountResponse CwIbcCountResponse
	err = counterpartyChain.QueryContract(ctx, juno2ContractAddr, queryCounterpartyCountMsg, &juno2IncrementedCountResponse)
	require.NoError(t, err)
	require.Equal(t, 1, juno2IncrementedCountResponse.Data.Count)

	// Query the count of the contract on juno1- should still be 0 as no packets have been sent through
	var juno1PreIncrementedCountResponse CwIbcCountResponse
	err = chain.QueryContract(ctx, juno1ContractAddr, queryChainCountMsg, &juno1PreIncrementedCountResponse)
	require.NoError(t, err)
	require.Equal(t, 0, juno1PreIncrementedCountResponse.Data.Count)

	// send packet from juno2 to juno1 and increment the juno1 contract count
	juno2Increment, err := counterpartyChain.ExecuteContract(ctx, junoUser.KeyName(), juno2ContractAddr, counterpartyIncrementMsg)
	require.NoError(t, err)
	require.Equal(t, uint32(0), juno2Increment.Code)

	// Wait for the ibc packet to be delivered
	err = testutil.WaitForBlocks(ctx, 2, chain, counterpartyChain)
	require.NoError(t, err)

	// Query the count of the contract on juno1- should still be 1 as a single packet has been sent through
	var juno1IncrementedCountResponse CwIbcCountResponse
	err = chain.QueryContract(ctx, juno1ContractAddr, queryChainCountMsg, &juno1IncrementedCountResponse)
	require.NoError(t, err)
	require.Equal(t, 1, juno1IncrementedCountResponse.Data.Count)

	// Query the count of the contract on juno2- should be 1 as a single packet has now been sent through from juno1 to juno2
	var juno2PreIncrementedCountResponse CwIbcCountResponse
	err = counterpartyChain.QueryContract(ctx, juno2ContractAddr, queryCounterpartyCountMsg, &juno2PreIncrementedCountResponse)
	require.NoError(t, err)
	require.Equal(t, 1, juno2PreIncrementedCountResponse.Data.Count)
}

// cw_ibc_example response data
type CwIbcCountResponse struct {
	Data struct {
		Count int `json:"count"`
	} `json:"data"`
}
