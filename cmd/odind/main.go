package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"

	app "github.com/ODIN-PROTOCOL/odin-core/app"
	"github.com/ODIN-PROTOCOL/odin-core/cmd/odind/cmd"
)

func main() {
	app.SetBech32AddressPrefixesAndBip44CoinTypeAndSeal(sdk.GetConfig())
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
