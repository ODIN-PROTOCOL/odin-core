package cosmos

import (
	"context"
	"testing"

	"cosmossdk.io/math"
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
	CosmosChainWasmIBCTest(t, "odin", "v0.9.1", "odinprotocol/core", "v0.9.1", "v0.9.0")
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

	/*

		height, err := chain.Height(ctx)
		require.NoError(t, err, "error fetching height before submit upgrade proposal")

		t.Log("Sending upgrade proposal")

		haltHeight := height + int64(haltHeightDelta)

		govAddr, err := chain.AuthQueryModuleAddress(ctx, "gov")
		require.NoError(t, err)

		prop, err := chain.BuildProposal([]cosmos.ProtoMessage{&upgradetypes.MsgSoftwareUpgrade{
			Authority: govAddr,
			Plan: upgradetypes.Plan{
				Name:   "v0.9.0",
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

	*/

	err := r.GeneratePath(ctx, rep.RelayerExecReporter(t), chain.Config().ChainID, counterpartyChain.Config().ChainID, ibcPath)
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

}
