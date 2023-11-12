package helpers

// import (
// 	"net/http"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	// "github.com/cosmos/cosmos-sdk/testutil/rest"
// 	// govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
// )

// func EmptyRestHandler(client.Context) govrest.ProposalRESTHandler {
// 	return govrest.ProposalRESTHandler{
// 		SubRoute: "unsupported-ibc-client",
// 		Handler: func(w http.ResponseWriter, r *http.Request) {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, "Legacy REST Routes are not supported for IBC proposals")
// 		},
// 	}
// }
