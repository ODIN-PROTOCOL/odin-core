package rest

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/gorilla/mux"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/testutil/rest"

// 	coinswaptypes "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
// )

// func getParamsHandler(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}
// 		res, height, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", coinswaptypes.QuerierRoute, coinswaptypes.QueryParams))
// 		if rest.CheckInternalServerError(w, err) {
// 			return
// 		}

// 		clientCtx = clientCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func getRateHandler(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}

// 		vars := mux.Vars(r)
// 		from := vars["from"]
// 		if from == "" {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid value, from: %s", from))
// 			return
// 		}

// 		to := vars["to"]
// 		if to == "" {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid value, to: %s", to))
// 			return
// 		}

// 		params := coinswaptypes.QueryRateRequest{
// 			From: from,
// 			To:   to,
// 		}

// 		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
// 		if rest.CheckBadRequestError(w, err) {
// 			return
// 		}

// 		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", coinswaptypes.QuerierRoute, coinswaptypes.QueryRate), bz)
// 		if rest.CheckInternalServerError(w, err) {
// 			return
// 		}

// 		clientCtx = clientCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }
