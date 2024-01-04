package rest

// import (
// 	"fmt"

// 	"github.com/gorilla/mux"

// 	"github.com/cosmos/cosmos-sdk/client"

// 	auctiontypes "github.com/ODIN-PROTOCOL/odin-core/x/auction/types"
// )

// func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
// 	rtr.HandleFunc(fmt.Sprintf("/%s/params", auctiontypes.ModuleName), getParamsHandler(clientCtx)).Methods("GET")
// 	rtr.HandleFunc(fmt.Sprintf("/%s/status", auctiontypes.ModuleName), getAuctionStatusHandler(clientCtx)).Methods("GET")
// }
