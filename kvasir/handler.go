package kvasir

import (
	"encoding/hex"
	"strconv"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

type processingResult struct {
	result  []byte
	version string
}

func handleTransaction(c *Context, l *Logger, tx abci.TxResult) {
	l.Debug(":eyes: Inspecting incoming transaction: %X", tmhash.Sum(tx.Tx))
	if tx.Result.Code != 0 {
		l.Debug(":alien: Skipping transaction with non-zero code: %d", tx.Result.Code)
		return
	}

	go handleRequestLog(c, l, tx.Result.Events)
}

func handleRequestLog(c *Context, l *Logger, log []abci.Event) {
	idStrs := GetEventValues(log, "wasm-wasm_request", "id")
	contractStrs := GetEventValues(log, "wasm-wasm_request", "contract")

	for i, idStr := range idStrs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			l.Error(":cold_sweat: Failed to convert %s to integer with error: %s", c, idStr, err.Error())
			return
		}

		// TODO: write better solution
		contract := contractStrs[i]

		// If id is in pending requests list, then skip it.
		if c.pendingRequests[RequestKey{
			ContractAddress: contract,
			RequestID:       id,
		}] {
			l.Debug(":eyes: Request is in pending list, then skip")
			return
		}

		request, err := queryPendingRequest(c, l, id, contract)
		if err != nil {
			l.Error(":cold_sweat: Failed to query pending request with error: %s", c, err.Error())
			continue
		}

		go handleRequest(c, l, contract, request)
	}
}

func handleRequest(c *Context, l *Logger, contract string, request Request) {
	l = l.With("contract", contract, "rid", request.RequestID)

	hasMe := false
	for _, val := range request.ChosenValidators {
		if val == c.validator.String() {
			hasMe = true
			break
		}
	}
	if !hasMe {
		l.Debug(":next_track_button: Skip request not related to this validator")
		return
	}

	l.Info(":delivery_truck: Processing request")

	keyIndex := c.nextKeyIndex()
	key := c.keys[keyIndex]

	rawRequest := rawRequest{
		contract:  contract,
		requestID: request.RequestID,
		calldata:  request.Metadata,
	}

	// process raw requests
	result, err := handleRawRequest(c, l, rawRequest, key)
	if err != nil {
		l.Debug("Cannot process request: %s", err.Error())
		return
	}

	c.pendingMsgs <- ReportMsgWithKey{
		contractAddress: contract,
		requestID:       request.RequestID,
		//msg:             types.NewMsgReportData(types.RequestID(id), reports, c.validator),
		execVersion: result.version,
		keyIndex:    keyIndex,
		result:      result.result,
	}
}

func handleRawRequest(
	c *Context,
	l *Logger,
	req rawRequest,
	key *keyring.Record,
) (*processingResult, error) {
	c.updateHandlingGauge(1)
	defer c.updateHandlingGauge(-1)

	vmsg := NewRequestVerification(cfg.ChainID, c.validator, req.requestID, req.contract)
	sig, pubkey, err := kb.Sign(key.Name, vmsg.GetSignBytes(), signing.SignMode_SIGN_MODE_DIRECT)
	if err != nil {
		l.Error(":skull: Failed to sign verify message: %s", c, err.Error())
		return nil, err
	}

	result, err := c.executor.Exec(req.calldata, map[string]interface{}{
		"ODIN_CHAIN_ID":   vmsg.ChainID,
		"ODIN_CONTRACT":   vmsg.Contract,
		"ODIN_VALIDATOR":  vmsg.Validator,
		"ODIN_REQUEST_ID": strconv.Itoa(int(vmsg.RequestID)),
		"ODIN_REPORTER":   hex.EncodeToString(pubkey.Bytes()),
		"ODIN_SIGNATURE":  sig,
	})

	if err != nil {
		l.Error(":skull: Failed to execute data source script: %s", c, err.Error())
		return nil, err
	} else {
		l.Debug(
			":sparkles: Query data done with calldata: %q, result: %q, exitCode: %d",
			req.calldata, result.Output, result.Code,
		)
		return &processingResult{
			result:  result.Output,
			version: result.Version,
		}, err
	}
}
