package rest

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/gorilla/mux"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/types/rest"

// 	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
// )

// func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
// 	r.HandleFunc(
// 		fmt.Sprintf("%s/%s", minttypes.LegacyRoute, minttypes.QueryParams),
// 		queryParamsHandlerFn(clientCtx),
// 	).Methods("GET")

// 	r.HandleFunc(
// 		fmt.Sprintf("%s/%s", minttypes.LegacyRoute, minttypes.QueryInflation),
// 		queryInflationHandlerFn(clientCtx),
// 	).Methods("GET")

// 	r.HandleFunc(
// 		fmt.Sprintf("%s/%s", minttypes.LegacyRoute, minttypes.QueryAnnualProvisions),
// 		queryAnnualProvisionsHandlerFn(clientCtx),
// 	).Methods("GET")

// 	r.HandleFunc(
// 		fmt.Sprintf("%s/%s/{%s}", minttypes.LegacyRoute, minttypes.QueryIntegrationAddresses, networkNameTag),
// 		queryIntegrationAddressHandlerFn(clientCtx),
// 	).Methods("GET")

// 	r.HandleFunc(
// 		fmt.Sprintf("%s/%s", minttypes.LegacyRoute, minttypes.QueryTreasuryPool),
// 		queryTreasuryPoolHandlerFn(clientCtx),
// 	).Methods("GET")

// 	r.HandleFunc(
// 		fmt.Sprintf("%s/%s", minttypes.LegacyRoute, minttypes.QueryCurrentMintVolume),
// 		queryCurrentMintVolumeHandlerFn(clientCtx),
// 	).Methods("GET")
// }

// func queryParamsHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		route := fmt.Sprintf("custom/%s/%s", minttypes.QuerierRoute, minttypes.QueryParams)

// 		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}

// 		res, height, err := clientCtx.QueryWithData(route, nil)
// 		if rest.CheckInternalServerError(w, err) {
// 			return
// 		}

// 		clientCtx = clientCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryInflationHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		route := fmt.Sprintf("custom/%s/%s", minttypes.QuerierRoute, minttypes.QueryInflation)

// 		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}

// 		res, height, err := clientCtx.QueryWithData(route, nil)
// 		if rest.CheckInternalServerError(w, err) {
// 			return
// 		}

// 		clientCtx = clientCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryAnnualProvisionsHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		route := fmt.Sprintf("custom/%s/%s", minttypes.QuerierRoute, minttypes.QueryAnnualProvisions)

// 		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}

// 		res, height, err := clientCtx.QueryWithData(route, nil)
// 		if rest.CheckInternalServerError(w, err) {
// 			return
// 		}

// 		clientCtx = clientCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, clientCtx, res)
// 	}
// }

// func queryIntegrationAddressHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}

// 		vars := mux.Vars(r)

// 		res, height, err := clientCtx.Query(fmt.Sprintf(
// 			"custom/%s/%s/%s",
// 			minttypes.QuerierRoute,
// 			minttypes.QueryIntegrationAddresses,
// 			vars[networkNameTag],
// 		))
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		cliCtx = cliCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, cliCtx, res)
// 	}
// }

// func queryTreasuryPoolHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		route := fmt.Sprintf("custom/%s/%s", minttypes.QuerierRoute, minttypes.QueryTreasuryPool)

// 		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}

// 		res, height, err := cliCtx.Query(route)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		cliCtx = cliCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, cliCtx, res)
// 	}
// }

// func queryCurrentMintVolumeHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		route := fmt.Sprintf("custom/%s/%s", minttypes.QuerierRoute, minttypes.QueryCurrentMintVolume)

// 		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
// 		if !ok {
// 			return
// 		}

// 		res, height, err := cliCtx.Query(route)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		cliCtx = cliCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, cliCtx, res)
// 	}
// }
