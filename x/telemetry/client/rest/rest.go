package rest

import (
	"fmt"
	commonrest "github.com/ODIN-PROTOCOL/odin-core/x/common/client/rest"
	telemetrytypes "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	rtr.HandleFunc(
		fmt.Sprintf("/%s/%s", telemetrytypes.ModuleName, telemetrytypes.QueryTopBalances),
		getTopBalancesHandler(clientCtx),
	).Methods("GET")

	rtr.HandleFunc(
		fmt.Sprintf("/%s/%s", telemetrytypes.ModuleName, telemetrytypes.QueryExtendedValidators),
		getExtendedValidatorsHandler(clientCtx),
	).Methods("GET")

	/*rtr.HandleFunc(
		fmt.Sprintf("/%s/%s", telemetrytypes.ModuleName, telemetrytypes.QueryAvgBlockSize),
		getAvgBlockSizeHandler(clientCtx),
	).Methods("GET")

	rtr.HandleFunc(
		fmt.Sprintf("/%s/%s", telemetrytypes.ModuleName, telemetrytypes.QueryAvgBlockTime),
		getAvgBlockTimeHandler(clientCtx),
	).Methods("GET")

	rtr.HandleFunc(
		fmt.Sprintf("/%s/%s", telemetrytypes.ModuleName, telemetrytypes.QueryAvgTxFee),
		getAvgTxFeeHandler(clientCtx),
	).Methods("GET")

	rtr.HandleFunc(
		fmt.Sprintf("/%s/%s", telemetrytypes.ModuleName, telemetrytypes.QueryTxVolume),
		getTxVolumeHandler(clientCtx),
	).Methods("GET")*/

	rtr.HandleFunc(
		fmt.Sprintf("/%s/%s", telemetrytypes.ModuleName, telemetrytypes.QueryTopValidators),
		getTopValidatorsHandler(clientCtx),
	).Methods("GET")

	/*rtr.HandleFunc(
		fmt.Sprintf("/%s/%s", telemetrytypes.ModuleName, telemetrytypes.QueryValidatorBlocks),
		getValidatorBlocksHandler(clientCtx),
	).Methods("GET")*/

	rtr.HandleFunc(
		fmt.Sprintf("/%s/%s/{consensus_address}", telemetrytypes.ModuleName, telemetrytypes.QueryValidatorByConsAddress),
		getValidatorByConsAddressHandler(clientCtx),
	).Methods("GET")
}

func getTopBalancesHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		paginationParams, ok := commonrest.CheckPaginationParams(w, r)
		if !ok {
			return
		}
		bin := clientCtx.LegacyAmino.MustMarshalJSON(paginationParams)

		query := r.URL.Query()

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf(
			"custom/%s/%s/%s",
			telemetrytypes.QuerierRoute,
			telemetrytypes.QueryTopBalances,
			query.Get(telemetrytypes.DenomTag),
		), bin)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getExtendedValidatorsHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		paginationParams, ok := commonrest.CheckPaginationParams(w, r)
		if !ok {
			return
		}
		bin := clientCtx.LegacyAmino.MustMarshalJSON(paginationParams)

		vars := mux.Vars(r)

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf(
			"custom/%s/%s/%s",
			telemetrytypes.QuerierRoute,
			telemetrytypes.QueryExtendedValidators,
			vars[telemetrytypes.StatusTag],
		), bin)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getAvgBlockSizeHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		var request telemetrytypes.QueryAvgBlockSizeRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &request) {
			return
		}
		bin := clientCtx.LegacyAmino.MustMarshalJSON(request)

		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", telemetrytypes.QuerierRoute, telemetrytypes.QueryAvgBlockSize),
			bin,
		)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getAvgBlockTimeHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		var request telemetrytypes.QueryAvgBlockTimeRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &request) {
			return
		}
		bin := clientCtx.LegacyAmino.MustMarshalJSON(request)

		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", telemetrytypes.QuerierRoute, telemetrytypes.QueryAvgBlockTime),
			bin,
		)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getAvgTxFeeHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		var request telemetrytypes.QueryAvgTxFeeRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &request) {
			return
		}
		bin := clientCtx.LegacyAmino.MustMarshalJSON(request)

		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", telemetrytypes.QuerierRoute, telemetrytypes.QueryAvgTxFee),
			bin,
		)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getTxVolumeHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		var request telemetrytypes.QueryTxVolumeRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &request) {
			return
		}
		bin := clientCtx.LegacyAmino.MustMarshalJSON(request)

		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", telemetrytypes.QuerierRoute, telemetrytypes.QueryTxVolume),
			bin,
		)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getTopValidatorsHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		var request telemetrytypes.QueryTopValidatorsRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &request) {
			return
		}
		bin := clientCtx.LegacyAmino.MustMarshalJSON(request)

		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", telemetrytypes.QuerierRoute, telemetrytypes.QueryTopValidators),
			bin,
		)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getValidatorBlocksHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		var request telemetrytypes.QueryValidatorBlocksRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &request) {
			return
		}
		bin := clientCtx.LegacyAmino.MustMarshalJSON(request)

		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", telemetrytypes.QuerierRoute, telemetrytypes.QueryValidatorBlocks),
			bin,
		)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getValidatorByConsAddressHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		validatorConsAddr := vars["consensus_address"]

		res, height, err := clientCtx.Query(fmt.Sprintf(
			"custom/%s/%s/%s",
			telemetrytypes.QuerierRoute,
			telemetrytypes.QueryValidatorByConsAddress,
			validatorConsAddr,
		))

		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
