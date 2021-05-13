package rest

import (
	govrest "github.com/GeoDB-Limited/odin-core/x/gov/client/rest"
	govtypes "github.com/GeoDB-Limited/odin-core/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	upgraderest "github.com/cosmos/cosmos-sdk/x/upgrade/client/rest"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"net/http"
	"time"
)

func ProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "upgrade",
		Handler:  newPostPlanHandler(clientCtx),
	}
}

func ProposalCancelRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "upgrade",
		Handler:  newCancelPlanHandler(clientCtx),
	}
}

func newPostPlanHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req upgraderest.PlanRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		var t time.Time
		if req.UpgradeTime != "" {
			t, err = time.Parse(time.RFC3339, req.UpgradeTime)
			if rest.CheckBadRequestError(w, err) {
				return
			}
		}

		plan := types.Plan{Name: req.UpgradeName, Time: t, Height: req.UpgradeHeight, Info: req.UpgradeInfo}
		content := types.NewSoftwareUpgradeProposal(req.Title, req.Description, plan)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func newCancelPlanHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req upgraderest.CancelRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewCancelSoftwareUpgradeProposal(req.Title, req.Description)

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
