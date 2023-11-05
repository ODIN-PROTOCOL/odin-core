// package rest

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	// "github.com/cosmos/cosmos-sdk/testutil/rest"

// 	auctiontypes "github.com/ODIN-PROTOCOL/odin-core/x/auction/types"
// )

// func getParamsHandler(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}
// 		res, height, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", auctiontypes.QuerierRoute, auctiontypes.QueryParams))
// 		if rest.CheckInternalServerError(w, err) {
// 			return
// 		}

// 		clientCtx = clientCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func getAuctionStatusHandler(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}
// 		res, height, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", auctiontypes.QuerierRoute, auctiontypes.QueryAuctionStatus))
// 		if rest.CheckInternalServerError(w, err) {
// 			return
// 		}

// 		clientCtx = clientCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }
