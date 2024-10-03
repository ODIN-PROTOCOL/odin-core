package kvasir

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"time"

	shell "github.com/ipfs/go-ipfs-api"

	app "github.com/ODIN-PROTOCOL/odin-core/app"
	wasmtypes "github.com/ODIN-PROTOCOL/wasmd/x/wasm/types"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var (
	// Proto codec for encoding/decoding proto message
	cdc = app.MakeEncodingConfig().Marshaler
)

func signAndBroadcast(
	c *Context, key *keyring.Record, msgs []sdk.Msg, gasLimit uint64, memo string,
) (string, error) {
	clientCtx := client.Context{
		Client:            c.client,
		Codec:             cdc,
		TxConfig:          app.MakeEncodingConfig().TxConfig,
		BroadcastMode:     "sync",
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

	address, err := key.GetAddress()
	if err != nil {
		return "", err
	}

	execMsg := authz.NewMsgExec(address, msgs)

	txb, err := txf.BuildUnsignedTx(&execMsg)
	if err != nil {
		return "", err
	}

	err = tx.Sign(context.Background(), txf, key.Name, txb, true)
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

	return res.TxHash, nil
}

func queryAccount(clientCtx client.Context, key *keyring.Record) (client.Account, error) {
	accountRetriever := authtypes.AccountRetriever{}

	address, err := key.GetAddress()
	if err != nil {
		return nil, err
	}

	acc, err := accountRetriever.GetAccount(clientCtx, address)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

// SubmitReport TODO: rework
func SubmitReport(c *Context, l *Logger, keyIndex int64, reports []ReportMsgWithKey) {
	// Return key and update pending metric when done with SubmitReport whether successfully or not.
	defer func() {
		c.freeKeys <- keyIndex
	}()
	defer c.updatePendingGauge(int64(-len(reports)))

	// Summarize execute version
	versionMap := make(map[string]bool)
	msgs := make([]sdk.Msg, len(reports))
	ids := make([]uint64, len(reports))

	for i, report := range reports {
		hash, err := uploadToIPFS(c, l, report.result, report)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(hash)

		contractMsg := fmt.Sprintf(
			"{\"report_data\":{\"val_address\": \"%s\", \"request_id\": %d, \"response\": \"%s\"}}",
			c.validator.String(),
			report.request.RequestID,
			hash,
		)

		msg := wasmtypes.MsgExecuteContract{
			Sender:   c.validatorAccAddr.String(),
			Contract: report.contractAddress,
			Msg:      []byte(contractMsg),
			Funds:    sdk.Coins{},
		}

		if err := msg.ValidateBasic(); err != nil {
			l.Error(":exploding_head: Failed to validate basic with error: %s", c, err.Error())
			return
		}
		msgs[i] = sdk.Msg(&msg)
		ids[i] = report.request.RequestID
		versionMap[report.execVersion] = true
	}
	l = l.With("rids", ids)

	versions := make([]string, 0, len(versionMap))
	for exec := range versionMap {
		versions = append(versions, exec)
	}
	memo := fmt.Sprintf("kvasir:%s/exec:%s", version.Version, strings.Join(versions, ","))
	key := c.keys[keyIndex]

	clientCtx := client.Context{
		Client:            c.client,
		TxConfig:          app.MakeEncodingConfig().TxConfig,
		InterfaceRegistry: app.MakeEncodingConfig().InterfaceRegistry,
	}

	gasLimit := uint64(300_000 * len(msgs))
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
				gasLimit = gasLimit * 110 / 100
				l.Info(":fuel_pump: Tx(%s) is out of gas and will be rebroadcasted with %d gas", txHash, gasLimit)
				txFound = true
				break FindTx
			} else {
				l.Error(":exploding_head: Tx returned nonzero code %d with log %s, tx hash: %s", c, txRes.Code, txRes.RawLog, txRes.TxHash)
				return
			}
		}
		if !txFound {
			l.Error(
				":question_mark: Cannot get transaction response from hash: %s transaction might be included in the next few blocks or check your node's health.",
				c,
				txHash,
			)
			return
		}
	}
	l.Error(":anxious_face_with_sweat: Cannot send reports with adjusted gas: %d", c, gasLimit)
}

// abciQuery will try to query data from OdinChain node maxTry time before give up and return error
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

func uploadToIPFS(c *Context, l *Logger, file []byte, report ReportMsgWithKey) (string, error) {
	pngFile, err := os.Open(string(file))
	if err != nil {
		return "", err
	}
	defer func() {
		pngFile.Close()
		os.Remove(string(file))
	}()

	img, err := png.Decode(pngFile)
	if err != nil {
		return "", err
	}

	// Create a new JPG file
	jpgFile, err := os.Create(string(file) + ".jpg")
	if err != nil {
		return "", err
	}
	defer func() {
		jpgFile.Close()
		os.Remove(string(file) + ".jpg")
	}()

	options := jpeg.Options{
		Quality: 80,
	}

	// Encode the image to JPEG and save it
	err = jpeg.Encode(jpgFile, img, &options)
	if err != nil {
		return "", err
	}

	// Connect to local IPFS node
	sh := shell.NewShell(c.ipfs)
	// Upload the image to IPFS
	pngFile.Seek(0, 0)
	pngCID, err := sh.Add(pngFile, shell.Pin(true))
	if err != nil {
		return "", err
	}

	jpgFile.Seek(0, 0)
	jpgCID, err := sh.Add(jpgFile, shell.Pin(true))
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"image":          fmt.Sprintf("https://ipfs.io/ipfs/%s", pngCID),
		"preview":        fmt.Sprintf("https://ipfs.io/ipfs/%s", jpgCID),
		"request_id":     report.request.RequestID,
		"request_height": report.request.RequestHeight,
		"prompt":         report.request.Metadata,
		"owner":          report.request.Sender,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	cid, err := sh.Add(bytes.NewReader(jsonData), shell.Pin(true))
	if err != nil {
		return "", err
	}

	return cid, nil
}
