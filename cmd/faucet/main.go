package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	odin "github.com/ODIN-PROTOCOL/odin-core/app"
)

const (
	DefaultKeyringBackend = "test"
	DefaultHomeEnv        = "$HOME/.faucet"
)

// Global instance.
var faucet Faucet

type Faucet struct {
	config  Config
	keybase keyring.Keyring
}

func main() {
	appConfig := sdk.GetConfig()
	odin.SetBech32AddressPrefixesAndBip44CoinTypeAndSeal(appConfig)

	rootCmd := &cobra.Command{
		Use:   "faucet",
		Short: "Faucet server for developers' network",
	}
	rootCmd.AddCommand(
		runCmd(),
		configCmd(),
		KeysCmd(),
	)

	clientCtx := client.GetClientContextFromCmd(rootCmd)
	cdc := clientCtx.Codec

	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		home, err := rootCmd.PersistentFlags().GetString(flags.FlagHome)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to parse home directory")
		}
		keyringBackend, err := rootCmd.Flags().GetString(flags.FlagKeyringBackend)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to parse keyring backend")
		}
		if err := os.MkdirAll(home, os.ModePerm); err != nil {
			return sdkerrors.Wrap(err, "failed to create a directory")
		}
		faucet.keybase, err = keyring.New(sdk.KeyringServiceName(), keyringBackend, home, nil, cdc)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to create a new keyring")
		}
		return initConfig(home)
	}
	rootCmd.PersistentFlags().String(flags.FlagHome, os.ExpandEnv(DefaultHomeEnv), "home directory")
	rootCmd.PersistentFlags().String(flags.FlagKeyringBackend, DefaultKeyringBackend, "keyring backend")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
