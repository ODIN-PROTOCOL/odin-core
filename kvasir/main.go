package kvasir

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	app "github.com/ODIN-PROTOCOL/odin-core/app"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagValidator        = "validator"
	flagLogLevel         = "log-level"
	flagExecutor         = "executor"
	flagBroadcastTimeout = "broadcast-timeout"
	flagRPCPollInterval  = "rpc-poll-interval"
	flagMaxTry           = "max-try"
	flagMaxReport        = "max-report"
)

type Config struct {
	ChainID           string   `mapstructure:"chain-id"`  // ChainID of the target chain
	NodeURI           string   `mapstructure:"node"`      // Remote RPC URI of OdinChain node to connect to
	Validator         string   `mapstructure:"validator"` // The validator address that I'm responsible for
	ValidatorAccAddr  string   `mapstructure:"validator-acc"`
	GasPrices         string   `mapstructure:"gas-prices"`          // Gas prices of the transaction
	LogLevel          string   `mapstructure:"log-level"`           // Log level of the logger
	Executor          string   `mapstructure:"executor"`            // Executor name and URL (example: "Executor name:URL")
	BroadcastTimeout  string   `mapstructure:"broadcast-timeout"`   // The time that Yoda will wait for tx commit
	RPCPollInterval   string   `mapstructure:"rpc-poll-interval"`   // The duration of rpc poll interval
	MaxTry            uint64   `mapstructure:"max-try"`             // The maximum number of tries to submit a report transaction
	MaxReport         uint64   `mapstructure:"max-report"`          // The maximum number of reports in one transaction
	MetricsListenAddr string   `mapstructure:"metrics-listen-addr"` // Address to listen on for prometheus metrics
	Contracts         []string `mapstructure:"contracts"`
	IPFS              string   `mapstructure:"ipfs"`
	GRPC              string   `mapstructure:"grpc"`
}

// Global instances.
var (
	cfg               Config
	kb                keyring.Keyring
	DefaultKvasirHome string
)

func initConfig(c *Context, cmd *cobra.Command) error {
	viper.SetConfigFile(path.Join(c.home, "config.yaml"))
	_ = viper.ReadInConfig() // If we fail to read config file, we'll just rely on cmd flags.
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	return nil
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultKvasirHome = filepath.Join(userHomeDir, ".kvasir")
}

func Main() {
	appConfig := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinTypeAndSeal(appConfig)

	ctx := &Context{}
	rootCmd := &cobra.Command{
		Use:   "yoda",
		Short: "OdinChain oracle daemon to subscribe and response to oracle requests",
	}

	rootCmd.AddCommand(
		configCmd(),
		keysCmd(ctx),
		runCmd(ctx),
		version.NewVersionCommand(),
	)
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		home, err := rootCmd.PersistentFlags().GetString(flags.FlagHome)
		if err != nil {
			return err
		}
		ctx.home = home
		if err := os.MkdirAll(home, os.ModePerm); err != nil {
			return err
		}
		kb, err = keyring.New("odin", keyring.BackendTest, home, nil, cdc)
		if err != nil {
			return err
		}
		return initConfig(ctx, rootCmd)
	}
	rootCmd.PersistentFlags().String(flags.FlagHome, DefaultKvasirHome, "home directory")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
