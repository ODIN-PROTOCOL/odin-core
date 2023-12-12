package cmd

import (
	"io"
	"os"

	dbm "github.com/cometbft/cometbft-db"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	odin "github.com/ODIN-PROTOCOL/odin-core/app"
	"github.com/ODIN-PROTOCOL/odin-core/app/params"
	"github.com/ODIN-PROTOCOL/odin-core/hooks/emitter"
	"github.com/ODIN-PROTOCOL/odin-core/hooks/request"
	tmcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

const (
	flagWithEmitter           = "with-emitter"
	flagDisableFeelessReports = "disable-feeless-reports"
	flagEnableFastSync        = "enable-fast-sync"
	flagWithPricer            = "with-pricer"
	flagWithRequestSearch     = "with-request-search"
	flagWithOwasmCacheSize    = "oracle-script-cache-size"
	flagEnableApi             = "api.enable"
)

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := odin.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(odin.DefaultNodeHome).WithViper("ODIN")

	srvCfg := serverconfig.DefaultConfig()
	cfg := tmcfg.DefaultConfig()

	rootCmd := &cobra.Command{
		Use:   "odind",
		Short: "Odin Consumer App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd, serverconfig.DefaultConfigTemplate, srvCfg, cfg)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	rootCmd.AddCommand(
		InitCmd(odin.NewDefaultGenesisState(), odin.DefaultNodeHome),
		// genesisCommand(
		// 	encodingConfig,
		// 	AddGenesisDataSourceCmd(odin.DefaultNodeHome),
		// 	AddGenesisOracleScriptCmd(odin.DefaultNodeHome),
		// ),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, odin.DefaultNodeHome, nil),
		genutilcli.GenTxCmd(odin.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, odin.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(odin.ModuleBasics),
		AddGenesisAccountCmd(odin.DefaultNodeHome),
		AddGenesisDataSourceCmd(odin.DefaultNodeHome),
		AddGenesisOracleScriptCmd(odin.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
	)

	server.AddCommands(rootCmd, odin.DefaultNodeHome, newApp, createSimappAndExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(odin.DefaultNodeHome),
	)

	rootCmd.PersistentFlags().String(flagWithRequestSearch, "", "[Experimental] Enable mode to save request in sql database")
	rootCmd.PersistentFlags().String(flagWithEmitter, "", "[Experimental] Enable mode with emitter")
	rootCmd.PersistentFlags().Uint32(flagWithOwasmCacheSize, 100, "[Experimental] Number of oracle scripts to cache")
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	odin.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	odin.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// newApp is an AppCreator
func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	// snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	// snapshotDB, err := dbm.NewGoLevelDB("metadata", snapshotDir)
	// if err != nil {
	// 	panic(err)
	// }

	// snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	// if err != nil {
	// 	panic(err)
	// }

	odinApp := odin.NewOdinApp(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		odin.MakeEncodingConfig(), // Ideally, we would reuse the one created by NewRootCmd.
		appOpts,
		cast.ToBool(appOpts.Get(flagDisableFeelessReports)),
		cast.ToUint32(appOpts.Get(flagWithOwasmCacheSize)),
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),

		// baseapp.SetSnapshotStore(snapshotStore),
		// baseapp.SetSnapshotInterval(cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval))),
		// baseapp.SetSnapshotKeepRecent(cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent))),
	)
	connStr, _ := appOpts.Get(flagWithRequestSearch).(string)
	if connStr != "" {
		odinApp.AddHook(request.NewHook(
			odinApp.AppCodec(), odinApp.OracleKeeper, connStr))
	}

	connStr, _ = appOpts.Get(flagWithEmitter).(string)
	if connStr != "" {
		odinApp.AddHook(
			emitter.NewHook(odinApp.AppCodec(), odinApp.LegacyAmino(), odin.MakeEncodingConfig(), odinApp.AccountKeeper, odinApp.BankKeeper,
				odinApp.StakingKeeper, odinApp.MintKeeper, odinApp.DistrKeeper, odinApp.GovKeeper,
				odinApp.OracleKeeper, connStr, false))
	}

	return odinApp
}

func createSimappAndExport(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string,
	appOpts servertypes.AppOptions, modulesToExport []string,
) (servertypes.ExportedApp, error) {
	encCfg := odin.MakeEncodingConfig() // Ideally, we would reuse the one created by NewRootCmd.
	encCfg.Marshaler = codec.NewProtoCodec(encCfg.InterfaceRegistry)
	var odinConsumerApp *odin.OdinApp
	if height != -1 {
		odinConsumerApp = odin.NewOdinApp(logger, db, traceStore, false, map[int64]bool{}, "", uint(1), encCfg, appOpts, false, cast.ToUint32(appOpts.Get(flagWithOwasmCacheSize)))

		if err := odinConsumerApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		odinConsumerApp = odin.NewOdinApp(logger, db, traceStore, true, map[int64]bool{}, "", uint(1), encCfg, appOpts, false, cast.ToUint32(appOpts.Get(flagWithOwasmCacheSize)))
	}

	return odinConsumerApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
