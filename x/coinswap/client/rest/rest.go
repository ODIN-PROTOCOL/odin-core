package rest

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"

	coinswaptypes "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
)

func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	rtr.HandleFunc(fmt.Sprintf("/%s/params", coinswaptypes.ModuleName), getParamsHandler(clientCtx)).Methods("GET")
	rtr.HandleFunc(fmt.Sprintf("/%s/rate/{from}/{to}", coinswaptypes.ModuleName), getRateHandler(clientCtx)).Methods("GET")
}
