package yoda

import (
	"context"
	"fmt"
	"strings"
	"time"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	app "github.com/ODIN-PROTOCOL/odin-core/app"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// Proto codec for encoding/decoding proto message
var cdc = app.MakeEncodingConfig().Marshaler

func signAndBroadcast(
	c *Context, key keyring.Info, msgs []sdk.Msg, gasLimit uint64, memo string,
) (string, error) {
	clientCtx := client.Context{
		Client:            c.client,
		TxConfig:          app.MakeEncodingConfig().TxConfig,
		BroadcastMode:     "async",
		InterfaceRegistry: app.MakeEncodingConfig().InterfaceRegistry,
	}
	acc, err := queryAccount(clientCtx, key)
	if err != nil {
		return "", fmt.Errorf("unable to get account: %w", err)
	}

	txf := tx.Factory{}.
		WithAccountNumber(acc.GetAccountNumber()).
		WithSequence(acc.GetSequence()).
		WithTxConfig(app.MakeEncodingConfig().TxConfig).
		WithGas(gasLimit).WithGasAdjustment(1).
		WithChainID(cfg.ChainID).
		WithMemo(memo).
		WithGasPrices(c.gasPrices).
		WithKeybase(kb).
		WithAccountRetriever(clientCtx.AccountRetriever)

	txb, err := tx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return "", err
	}

	err = tx.Sign(txf, key.GetName(), txb, true)
	if err != nil {
		return "", err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return "", err
	}

	// broadcast to a Tendermint node
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return "", err
	}
	// out, err := txBldr.WithKeybase(keybase).BuildAndSign(key.GetName(), ckeys.DefaultKeyPass, msgs)
	// if err != nil {
	// 	return "", fmt.Errorf("Failed to build tx with error: %s", err.Error())
	// }
	return res.TxHash, nil
}

func queryAccount(clientCtx client.Context, key keyring.Info) (client.Account, error) {
	accountRetriever := authtypes.AccountRetriever{}
	acc, err := accountRetriever.GetAccount(clientCtx, key.GetAddress())
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func SubmitReport(c *Context, l *Logger, keyIndex int64, reports []ReportMsgWithKey) {
	// Return key and update pending metric when done with SubmitReport whether successfully or not.
	defer func() {
		c.freeKeys <- keyIndex
	}()
	defer c.updatePendingGauge(int64(-len(reports)))

	// Summarize execute version
	versionMap := make(map[string]bool)
	msgs := make([]sdk.Msg, len(reports))
	ids := make([]oracletypes.RequestID, len(reports))
	feeEstimations := make([]FeeEstimationData, len(reports))

	for i, report := range reports {
		if err := report.msg.ValidateBasic(); err != nil {
			l.Error(":exploding_head: Failed to validate basic with error: %s", c, err.Error())
			return
		}
		msgs[i] = report.msg
		ids[i] = report.msg.RequestID
		feeEstimations[i] = report.feeEstimationData
		for _, exec := range report.execVersion {
			versionMap[exec] = true
		}
	}
	l = l.With("rids", ids)

	versions := make([]string, 0, len(versionMap))
	for exec := range versionMap {
		versions = append(versions, exec)
	}
	memo := fmt.Sprintf("yoda:%s/exec:%s", version.Version, strings.Join(versions, ","))
	key := c.keys[keyIndex]
	// cliCtx := sdkCtx.CLIContext{Client: c.client, TrustNode: true, Codec: cdc}
	clientCtx := client.Context{
		Client:            c.client,
		TxConfig:          app.MakeEncodingConfig().TxConfig,
		InterfaceRegistry: app.MakeEncodingConfig().InterfaceRegistry,
	}

	gasLimit := estimateGas(c, l, msgs, feeEstimations)
	// We want to resend transaction only if tx returns Out of gas error.
	for sendAttempt := uint64(1); sendAttempt <= c.maxTry; sendAttempt++ {
		var txHash string
		l.Info(":e-mail: Sending report transaction attempt: (%d/%d)", sendAttempt, c.maxTry)
		for broadcastTry := uint64(1); broadcastTry <= c.maxTry; broadcastTry++ {
			l.Info(":writing_hand: Try to sign and broadcast report transaction(%d/%d)", broadcastTry, c.maxTry)
			hash, err := signAndBroadcast(c, key, msgs, gasLimit, memo)
			if err != nil {
				// Use info level because this error can happen and retry process can solve this error.
				l.Info(":warning: %s", err.Error())
				time.Sleep(c.rpcPollInterval)
				continue
			}
			// Transaction passed CheckTx process and wait to include in block.
			txHash = hash
			break
		}
		if txHash == "" {
			l.Error(":exploding_head: Cannot try to broadcast more than %d try", c, c.maxTry)
			return
		}
		txFound := false
	FindTx:
		for start := time.Now(); time.Since(start) < c.broadcastTimeout; {
			time.Sleep(c.rpcPollInterval)
			txRes, err := authtx.QueryTx(clientCtx, txHash)
			if err != nil {
				l.Debug(":warning: Failed to query tx with error: %s", err.Error())
				continue
			}

			if txRes.Code == 0 {
				l.Info(":smiling_face_with_sunglasses: Successfully broadcast tx with hash: %s", txHash)
				c.updateSubmittedCount(int64(len(reports)))
				return
			}
			if txRes.Codespace == sdkerrors.RootCodespace &&
				txRes.Code == sdkerrors.ErrOutOfGas.ABCICode() {
				// Increase gas limit and try to broadcast again
				gasLimit = gasLimit * 130 / 100
				l.Info(":fuel_pump: Tx(%s) is out of gas and will be rebroadcasted with %d gas", txHash, gasLimit)
				txFound = true
				break FindTx
			} else {
				l.Error(":exploding_head: Tx returned nonzero code %d with log %s, tx hash: %s", c, txRes.Code, txRes.RawLog, txRes.TxHash)
				return
			}
		}
		if !txFound {
			l.Error(":question_mark: Cannot get transaction response from hash: %s transaction might be included in the next few blocks or check your node's health.", c, txHash)
			return
		}
	}
	l.Error(":anxious_face_with_sweat: Cannot send reports with adjusted gas: %d", c, gasLimit)
}

// GetExecutable fetches data source executable using the provided client.
func GetExecutable(c *Context, l *Logger, hash string) ([]byte, error) {
	resValue, err := c.fileCache.GetFile(hash)
	if err != nil {
		l.Debug(":magnifying_glass_tilted_left: Fetching data source hash: %s from bandchain querier", hash)
		res, err := abciQuery(c, l, fmt.Sprintf("custom/%s/%s/%s", oracletypes.StoreKey, oracletypes.QueryData, hash), nil)
		if err != nil {
			l.Error(":exploding_head: Failed to get data source with error: %s", c, err.Error())
			return nil, err
		}
		resValue = res.Response.GetValue()
		c.fileCache.AddFile(resValue)
	} else {
		l.Debug(":card_file_box: Found data source hash: %s in cache file", hash)
	}

	l.Debug(":balloon: Received data source hash: %s content: %q", hash, resValue[:32])
	return resValue, nil
}

// GetDataSourceHash fetches data source hash by id
func GetDataSourceHash(c *Context, l *Logger, id oracletypes.DataSourceID) (string, error) {
	if hash, ok := c.dataSourceCache.Load(id); ok {
		return hash.(string), nil
	}

	res, err := abciQuery(c, l, fmt.Sprintf("/store/%s/key", oracletypes.StoreKey), oracletypes.DataSourceStoreKey(id))
	if err != nil {
		l.Error(":skull: Failed to get data source with error: %s", c, err.Error())
		return "", err
	}

	var d oracletypes.DataSource
	cdc.MustUnmarshal(res.Response.Value, &d)

	hash, _ := c.dataSourceCache.LoadOrStore(id, d.Filename)

	return hash.(string), nil
}

// GetRequest fetches request by id
func GetRequest(c *Context, l *Logger, id oracletypes.RequestID) (oracletypes.Request, error) {
	res, err := abciQuery(c, l, fmt.Sprintf("/store/%s/key", oracletypes.StoreKey), oracletypes.RequestStoreKey(id))
	if err != nil {
		l.Error(":skull: Failed to get request with error: %s", c, err.Error())
		return oracletypes.Request{}, err
	}

	var r oracletypes.Request
	cdc.MustUnmarshal(res.Response.Value, &r)

	return r, nil
}

// abciQuery will try to query data from BandChain node maxTry time before give up and return error
func abciQuery(c *Context, l *Logger, path string, data []byte) (*ctypes.ResultABCIQuery, error) {
	var lastErr error
	for try := 0; try < int(c.maxTry); try++ {
		res, err := c.client.ABCIQuery(context.Background(), path, data)
		if err != nil {
			l.Debug(":skull: Failed to query on %s request with error: %s", path, err.Error())
			lastErr = err
			time.Sleep(c.rpcPollInterval)
			continue
		}
		return res, nil
	}
	return nil, lastErr
}
