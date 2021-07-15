package rest

import (
	"fmt"
	commonrest "github.com/GeoDB-Limited/odin-core/x/common/client/rest"
	telemetrytypes "github.com/GeoDB-Limited/odin-core/x/telemetry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	rtr.HandleFunc(fmt.Sprintf("%s/top_balances", telemetrytypes.ModuleName), getTopBalancesHandler(clientCtx)).Methods("GET")
}

func getTopBalancesHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		paginationParams, ok := commonrest.СheckPaginationParams(w, r)
		if !ok {
			return
		}
		bin := clientCtx.LegacyAmino.MustMarshalJSON(paginationParams)
		res, height, err := clientCtx.QueryWithData(fmt.Sprintf(
			"custom/%s/%s",
			telemetrytypes.QuerierRoute,
			telemetrytypes.QueryTopBalances,
		), bin)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
