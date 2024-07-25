package oracle_test

import (
	"bytes"
	gz "compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"

	"github.com/ODIN-PROTOCOL/odin-core/testing/testapp"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	wasmtypes "github.com/ODIN-PROTOCOL/wasmvm/v2/types"
)

func TestCreateDataSourceSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	dsCount, err := k.GetDataSourceCount(ctx)
	require.NoError(t, err)
	treasury := testapp.Treasury.Address
	owner := testapp.Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	executableHash := sha256.Sum256(executable)
	filename := hex.EncodeToString(executableHash[:])
	msg := types.NewMsgCreateDataSource(
		name,
		description,
		executable,
		testapp.EmptyCoins,
		treasury,
		owner,
		testapp.Alice.Address,
	)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	ds, err := k.GetDataSource(ctx, types.DataSourceID(dsCount+1))
	require.NoError(t, err)
	require.Equal(
		t,
		types.NewDataSource(testapp.Owner.Address, name, description, filename, testapp.EmptyCoins, treasury),
		ds,
	)
	event := abci.Event{
		Type: types.EventTypeCreateDataSource,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyID, Value: fmt.Sprintf("%d", dsCount+1)},
		},
	}
	require.Equal(t, event, res.Events[0])
}

func TestCreateGzippedExecutableDataSourceFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	treasury := testapp.Treasury.Address
	owner := testapp.Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(executable)
	zw.Close()
	sender := testapp.Alice.Address
	msg := types.NewMsgCreateDataSource(
		name,
		description,
		buf.Bytes()[:5],
		testapp.EmptyCoins,
		treasury,
		owner,
		sender,
	)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.ErrorIs(t, err, types.ErrUncompressionFailed)
	require.Nil(t, res)
}

func TestEditDataSourceSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	newName := "beeb"
	newDescription := "new_description"
	newExecutable := []byte("executable2")
	newExecutableHash := sha256.Sum256(newExecutable)
	newFilename := hex.EncodeToString(newExecutableHash[:])
	msg := types.NewMsgEditDataSource(
		1,
		newName,
		newDescription,
		newExecutable,
		testapp.Coins1000000loki,
		testapp.Treasury.Address,
		testapp.Alice.Address,
		testapp.Owner.Address,
	)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	ds, err := k.GetDataSource(ctx, 1)
	require.NoError(t, err)
	require.Equal(
		t,
		types.NewDataSource(
			testapp.Alice.Address,
			newName,
			newDescription,
			newFilename,
			testapp.Coins1000000loki,
			testapp.Treasury.Address,
		),
		ds,
	)
	event := abci.Event{
		Type:       types.EventTypeEditDataSource,
		Attributes: []abci.EventAttribute{{Key: types.AttributeKeyID, Value: "1"}},
	}
	require.Equal(t, event, res.Events[0])
}

func TestEditDataSourceFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	newName := "beeb"
	newDescription := "new_description"
	newExecutable := []byte("executable2")
	// Bad ID
	msg := types.NewMsgEditDataSource(
		42,
		newName,
		newDescription,
		newExecutable,
		testapp.EmptyCoins,
		testapp.Treasury.Address,
		testapp.Owner.Address,
		testapp.Owner.Address,
	)
	res, err := oracle.NewHandler(k)(ctx, msg)
	testapp.CheckErrorf(t, err, types.ErrDataSourceNotFound, "id: 42")
	require.Nil(t, res)
	// Not owner
	msg = types.NewMsgEditDataSource(
		1,
		newName,
		newDescription,
		newExecutable,
		testapp.EmptyCoins,
		testapp.Treasury.Address,
		testapp.Owner.Address,
		testapp.Bob.Address,
	)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.ErrorIs(t, err, types.ErrEditorNotAuthorized)
	require.Nil(t, res)
	// Bad Gzip
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(newExecutable)
	zw.Close()
	msg = types.NewMsgEditDataSource(
		1,
		newName,
		newDescription,
		buf.Bytes()[:5],
		testapp.EmptyCoins,
		testapp.Treasury.Address,
		testapp.Owner.Address,
		testapp.Owner.Address,
	)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.ErrorIs(t, err, types.ErrUncompressionFailed)
	require.Nil(t, res)
}

func TestCreateOracleScriptSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	osCount, err := k.GetOracleScriptCount(ctx)
	require.NoError(t, err)
	name := "os_1"
	description := "beeb"
	code := testapp.WasmExtra1
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		name,
		description,
		schema,
		url,
		code,
		testapp.Owner.Address,
		testapp.Alice.Address,
	)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	os, err := k.GetOracleScript(ctx, types.OracleScriptID(osCount+1))
	require.NoError(t, err)
	require.Equal(
		t,
		types.NewOracleScript(testapp.Owner.Address, name, description, testapp.WasmExtra1FileName, schema, url),
		os,
	)

	event := abci.Event{
		Type: types.EventTypeCreateOracleScript,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyID, Value: fmt.Sprintf("%d", osCount+1)},
		},
	}
	require.Equal(t, event, res.Events[0])
}

func TestCreateGzippedOracleScriptSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	osCount, err := k.GetOracleScriptCount(ctx)
	require.NoError(t, err)
	name := "os_1"
	description := "beeb"
	schema := "schema"
	url := "url"
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(testapp.WasmExtra1)
	zw.Close()
	msg := types.NewMsgCreateOracleScript(
		name,
		description,
		schema,
		url,
		buf.Bytes(),
		testapp.Owner.Address,
		testapp.Alice.Address,
	)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	os, err := k.GetOracleScript(ctx, types.OracleScriptID(osCount+1))
	require.NoError(t, err)
	require.Equal(
		t,
		types.NewOracleScript(testapp.Owner.Address, name, description, testapp.WasmExtra1FileName, schema, url),
		os,
	)

	event := abci.Event{
		Type: types.EventTypeCreateOracleScript,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyID, Value: fmt.Sprintf("%d", osCount+1)},
		},
	}
	require.Equal(t, event, res.Events[0])
}

func TestCreateOracleScriptFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	name := "os_1"
	description := "beeb"
	schema := "schema"
	url := "url"
	// Bad Owasm code
	msg := types.NewMsgCreateOracleScript(
		name,
		description,
		schema,
		url,
		[]byte("BAD"),
		testapp.Owner.Address,
		testapp.Alice.Address,
	)
	res, err := oracle.NewHandler(k)(ctx, msg)
	testapp.CheckErrorf(t, err, types.ErrOwasmCompilation, "caused by %s", wasmtypes.ErrValidation)
	require.Nil(t, res)
	// Bad Gzip
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(testapp.WasmExtra1)
	zw.Close()
	msg = types.NewMsgCreateOracleScript(
		name,
		description,
		schema,
		url,
		buf.Bytes()[:5],
		testapp.Owner.Address,
		testapp.Alice.Address,
	)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.ErrorIs(t, err, types.ErrUncompressionFailed)
	require.Nil(t, res)
}

func TestEditOracleScriptSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := testapp.WasmExtra2
	newSchema := "new_schema"
	newURL := "new_url"
	msg := types.NewMsgEditOracleScript(
		1,
		newName,
		newDescription,
		newSchema,
		newURL,
		newCode,
		testapp.Alice.Address,
		testapp.Owner.Address,
	)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	os, err := k.GetOracleScript(ctx, 1)
	require.NoError(t, err)
	require.Equal(
		t,
		types.NewOracleScript(
			testapp.Alice.Address,
			newName,
			newDescription,
			testapp.WasmExtra2FileName,
			newSchema,
			newURL,
		),
		os,
	)

	event := abci.Event{
		Type:       types.EventTypeEditOracleScript,
		Attributes: []abci.EventAttribute{{Key: types.AttributeKeyID, Value: "1"}},
	}
	require.Equal(t, event, res.Events[0])
}

func TestEditOracleScriptFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := testapp.WasmExtra2
	newSchema := "new_schema"
	newURL := "new_url"
	// Bad ID
	msg := types.NewMsgEditOracleScript(
		999,
		newName,
		newDescription,
		newSchema,
		newURL,
		newCode,
		testapp.Owner.Address,
		testapp.Owner.Address,
	)
	res, err := oracle.NewHandler(k)(ctx, msg)
	testapp.CheckErrorf(t, err, types.ErrOracleScriptNotFound, "id: 999")
	require.Nil(t, res)
	// Not owner
	msg = types.NewMsgEditOracleScript(
		1,
		newName,
		newDescription,
		newSchema,
		newURL,
		newCode,
		testapp.Owner.Address,
		testapp.Bob.Address,
	)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "editor not authorized")
	require.Nil(t, res)
	// Bad Owasm code
	msg = types.NewMsgEditOracleScript(
		1,
		newName,
		newDescription,
		newSchema,
		newURL,
		[]byte("BAD_CODE"),
		testapp.Owner.Address,
		testapp.Owner.Address,
	)
	res, err = oracle.NewHandler(k)(ctx, msg)
	testapp.CheckErrorf(t, err, types.ErrOwasmCompilation, "caused by %s", wasmtypes.ErrValidation)
	require.Nil(t, res)
	// Bad Gzip
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(testapp.WasmExtra2)
	zw.Close()
	msg = types.NewMsgEditOracleScript(
		1,
		newName,
		newDescription,
		newSchema,
		newURL,
		buf.Bytes()[:5],
		testapp.Owner.Address,
		testapp.Owner.Address,
	)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.ErrorIs(t, err, types.ErrUncompressionFailed)
	require.Nil(t, res)
}

func TestRequestDataSuccess(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(true)
	ctx = ctx.WithBlockHeight(124).WithBlockTime(testapp.ParseTime(1581589790))
	msg := types.NewMsgRequestData(
		1,
		[]byte("beeb"),
		2,
		2,
		"CID",
		testapp.Coins100000000loki,
		testapp.TestDefaultPrepareGas,
		testapp.TestDefaultExecuteGas,
		testapp.FeePayer.Address,
	)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, types.NewRequest(
		1,
		[]byte("beeb"),
		[]sdk.ValAddress{testapp.Validators[2].ValAddress, testapp.Validators[0].ValAddress},
		2,
		124,
		testapp.ParseTime(1581589790),
		"CID",
		[]types.RawRequest{
			types.NewRawRequest(1, 1, []byte("beeb")),
			types.NewRawRequest(2, 2, []byte("beeb")),
			types.NewRawRequest(3, 3, []byte("beeb")),
		},
		nil,
		testapp.TestDefaultExecuteGas,
	), k.MustGetRequest(ctx, 1))
	event := abci.Event{
		Type: authtypes.EventTypeCoinSpent,
		Attributes: []abci.EventAttribute{
			{Key: authtypes.AttributeKeySpender, Value: testapp.FeePayer.Address.String()},
			{Key: sdk.AttributeKeyAmount, Value: "2000000loki"},
		},
	}
	require.Equal(t, event, res.Events[0])
	require.Equal(t, event, res.Events[4])
	require.Equal(t, event, res.Events[8])
	event = abci.Event{
		Type: authtypes.EventTypeCoinReceived,
		Attributes: []abci.EventAttribute{
			{Key: authtypes.AttributeKeyReceiver, Value: app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName).String()},
			{Key: sdk.AttributeKeyAmount, Value: "2000000loki"},
		},
	}
	require.Equal(t, event, res.Events[1])
	require.Equal(t, event, res.Events[5])
	require.Equal(t, event, res.Events[9])
	event = abci.Event{
		Type: authtypes.EventTypeTransfer,
		Attributes: []abci.EventAttribute{
			{Key: authtypes.AttributeKeyRecipient, Value: app.AccountKeeper.GetModuleAddress(distrtypes.ModuleName).String()},
			{Key: authtypes.AttributeKeySender, Value: testapp.FeePayer.Address.String()},
			{Key: sdk.AttributeKeyAmount, Value: "2000000loki"},
		},
	}
	require.Equal(t, event, res.Events[2])
	require.Equal(t, event, res.Events[6])
	require.Equal(t, event, res.Events[10])
	event = abci.Event{
		Type: sdk.EventTypeMessage,
		Attributes: []abci.EventAttribute{
			{Key: authtypes.AttributeKeySender, Value: testapp.FeePayer.Address.String()},
		},
	}
	require.Equal(t, event, res.Events[3])
	require.Equal(t, event, res.Events[7])
	require.Equal(t, event, res.Events[11])

	event = abci.Event{
		Type: types.EventTypeRequest,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyID, Value: "1"},
			{Key: types.AttributeKeyClientID, Value: "CID"},
			{Key: types.AttributeKeyOracleScriptID, Value: "1"},
			{Key: types.AttributeKeyCalldata, Value: "62656562"}, // "beeb" in hex
			{Key: types.AttributeKeyAskCount, Value: "2"},
			{Key: types.AttributeKeyMinCount, Value: "2"},
			{Key: types.AttributeKeyGasUsed, Value: "5294700000"},
			{Key: types.AttributeKeyTotalFees, Value: "6000000loki"},
			{Key: types.AttributeKeyValidator, Value: testapp.Validators[2].ValAddress.String()},
			{Key: types.AttributeKeyValidator, Value: testapp.Validators[0].ValAddress.String()},
		},
	}
	require.Equal(t, event, res.Events[12])
	event = abci.Event{
		Type: types.EventTypeRawRequest,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyDataSourceID, Value: "1"},
			{Key: types.AttributeKeyDataSourceHash, Value: testapp.DataSources[1].Filename},
			{Key: types.AttributeKeyExternalID, Value: "1"},
			{Key: types.AttributeKeyCalldata, Value: "beeb"},
			{Key: types.AttributeKeyFee, Value: "1000000loki"},
		},
	}
	require.Equal(t, event, res.Events[13])
	event = abci.Event{
		Type: types.EventTypeRawRequest,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyDataSourceID, Value: "2"},
			{Key: types.AttributeKeyDataSourceHash, Value: testapp.DataSources[2].Filename},
			{Key: types.AttributeKeyExternalID, Value: "2"},
			{Key: types.AttributeKeyCalldata, Value: "beeb"},
			{Key: types.AttributeKeyFee, Value: "1000000loki"},
		},
	}
	require.Equal(t, event, res.Events[14])
	event = abci.Event{
		Type: types.EventTypeRawRequest,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyDataSourceID, Value: "3"},
			{Key: types.AttributeKeyDataSourceHash, Value: testapp.DataSources[3].Filename},
			{Key: types.AttributeKeyExternalID, Value: "3"},
			{Key: types.AttributeKeyCalldata, Value: "beeb"},
			{Key: types.AttributeKeyFee, Value: "1000000loki"},
		},
	}
	require.Equal(t, event, res.Events[15])
}

func TestRequestDataFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	// No active oracle validators
	res, err := oracle.NewHandler(
		k,
	)(
		ctx,
		types.NewMsgRequestData(
			1,
			[]byte("beeb"),
			2,
			2,
			"CID",
			testapp.Coins100000000loki,
			testapp.TestDefaultPrepareGas,
			testapp.TestDefaultExecuteGas,
			testapp.FeePayer.Address,
		),
	)
	testapp.CheckErrorf(t, err, types.ErrInsufficientValidators, "0 < 2")
	require.Nil(t, res)
	err = k.Activate(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	err = k.Activate(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	// Too large calldata
	res, err = oracle.NewHandler(
		k,
	)(
		ctx,
		types.NewMsgRequestData(
			1,
			[]byte(strings.Repeat("beeb", 2000)),
			2,
			2,
			"CID",
			testapp.Coins100000000loki,
			testapp.TestDefaultPrepareGas,
			testapp.TestDefaultExecuteGas,
			testapp.FeePayer.Address,
		),
	)
	testapp.CheckErrorf(t, err, types.ErrTooLargeCalldata, "got: 8000, max: 512")
	require.Nil(t, res)
	// Too high ask count
	res, err = oracle.NewHandler(
		k,
	)(
		ctx,
		types.NewMsgRequestData(
			1,
			[]byte("beeb"),
			3,
			2,
			"CID",
			testapp.Coins100000000loki,
			testapp.TestDefaultPrepareGas,
			testapp.TestDefaultExecuteGas,
			testapp.FeePayer.Address,
		),
	)
	testapp.CheckErrorf(t, err, types.ErrInsufficientValidators, "2 < 3")
	require.Nil(t, res)
	// Bad oracle script ID
	res, err = oracle.NewHandler(
		k,
	)(
		ctx,
		types.NewMsgRequestData(
			999,
			[]byte("beeb"),
			2,
			2,
			"CID",
			testapp.Coins100000000loki,
			testapp.TestDefaultPrepareGas,
			testapp.TestDefaultExecuteGas,
			testapp.FeePayer.Address,
		),
	)
	testapp.CheckErrorf(t, err, types.ErrOracleScriptNotFound, "id: 999")
	require.Nil(t, res)
	// Pay not enough fee
	res, err = oracle.NewHandler(
		k,
	)(
		ctx,
		types.NewMsgRequestData(
			1,
			[]byte("beeb"),
			2,
			2,
			"CID",
			testapp.EmptyCoins,
			testapp.TestDefaultPrepareGas,
			testapp.TestDefaultExecuteGas,
			testapp.FeePayer.Address,
		),
	)
	testapp.CheckErrorf(t, err, types.ErrNotEnoughFee, "require: 2000000loki, max: 0loki")
	require.Nil(t, res)
}

func TestReportSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Set up a mock request asking 3 validators with min count 2.
	err := k.SetRequest(ctx, 42, types.NewRequest(
		1,
		[]byte("beeb"),
		[]sdk.ValAddress{
			testapp.Validators[2].ValAddress,
			testapp.Validators[1].ValAddress,
			testapp.Validators[0].ValAddress,
		},
		2,
		124,
		testapp.ParseTime(1581589790),
		"CID",
		[]types.RawRequest{
			types.NewRawRequest(1, 1, []byte("beeb")),
			types.NewRawRequest(2, 2, []byte("beeb")),
		},
		nil,
		0,
	))
	require.NoError(t, err)
	// Common raw reports for everyone.
	reports := []types.RawReport{types.NewRawReport(1, 0, []byte("data1")), types.NewRawReport(2, 0, []byte("data2"))}
	// Validators[0] reports data.
	res, err := oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Validators[0].ValAddress))
	require.NoError(t, err)

	pendingResolveList, err := k.GetPendingResolveList(ctx)
	require.NoError(t, err)
	require.Equal(t, []types.RequestID(nil), pendingResolveList)
	event := abci.Event{
		Type: types.EventTypeReport,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyID, Value: "42"},
			{Key: types.AttributeKeyValidator, Value: testapp.Validators[0].ValAddress.String()},
		},
	}
	require.Equal(t, event, res.Events[0])
	// Validators[1] reports data. Now the request should move to pending resolve.
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Validators[1].ValAddress))
	require.NoError(t, err)

	pendingResolveList, err = k.GetPendingResolveList(ctx)
	require.NoError(t, err)
	require.Equal(t, []types.RequestID{42}, pendingResolveList)
	event = abci.Event{
		Type: types.EventTypeReport,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyID, Value: "42"},
			{Key: types.AttributeKeyValidator, Value: testapp.Validators[1].ValAddress.String()},
		},
	}
	require.Equal(t, event, res.Events[0])
	// Even if we resolve the request, Validators[2] should still be able to report.
	err = k.SetPendingResolveList(ctx, []types.RequestID{})
	require.NoError(t, err)
	err = k.ResolveSuccess(ctx, 42, []byte("RESOLVE_RESULT!"), 1234)
	require.NoError(t, err)

	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Validators[2].ValAddress))
	require.NoError(t, err)
	event = abci.Event{
		Type: types.EventTypeReport,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyID, Value: "42"},
			{Key: types.AttributeKeyValidator, Value: testapp.Validators[2].ValAddress.String()},
		},
	}
	require.Equal(t, event, res.Events[0])
	// Check the reports of this request. We should see 3 reports, with report from Validators[2] comes after resolve.
	finalReport, err := k.GetReports(ctx, 42)
	require.NoError(t, err)
	require.Contains(t, finalReport, types.NewReport(testapp.Validators[0].ValAddress, true, reports))
	require.Contains(t, finalReport, types.NewReport(testapp.Validators[1].ValAddress, true, reports))
	require.Contains(t, finalReport, types.NewReport(testapp.Validators[2].ValAddress, false, reports))
}

func TestReportFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Set up a mock request asking 3 validators with min count 2.
	err := k.SetRequest(ctx, 42, types.NewRequest(
		1,
		[]byte("beeb"),
		[]sdk.ValAddress{
			testapp.Validators[2].ValAddress,
			testapp.Validators[1].ValAddress,
			testapp.Validators[0].ValAddress,
		},
		2,
		124,
		testapp.ParseTime(1581589790),
		"CID",
		[]types.RawRequest{
			types.NewRawRequest(1, 1, []byte("beeb")),
			types.NewRawRequest(2, 2, []byte("beeb")),
		},
		nil,
		0,
	))
	require.NoError(t, err)

	// Common raw reports for everyone.
	reports := []types.RawReport{types.NewRawReport(1, 0, []byte("data1")), types.NewRawReport(2, 0, []byte("data2"))}
	// Bad ID
	res, err := oracle.NewHandler(k)(ctx, types.NewMsgReportData(999, reports, testapp.Validators[0].ValAddress))
	testapp.CheckErrorf(t, err, types.ErrRequestNotFound, "id: 999")
	require.Nil(t, res)
	// Not-asked validator
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Alice.ValAddress))
	testapp.CheckErrorf(
		t,
		err,
		types.ErrValidatorNotRequested,
		"reqID: 42, val: %s",
		testapp.Alice.ValAddress.String(),
	)
	require.Nil(t, res)
	// Too large report data size
	res, err = oracle.NewHandler(
		k,
	)(
		ctx,
		types.NewMsgReportData(
			42,
			[]types.RawReport{
				types.NewRawReport(1, 0, []byte("data1")),
				types.NewRawReport(2, 0, []byte(strings.Repeat("data2", 2000))),
			},
			testapp.Validators[0].ValAddress,
		),
	)
	testapp.CheckErrorf(t, err, types.ErrTooLargeRawReportData, "got: 10000, max: 512")
	require.Nil(t, res)
	// Not having all raw reports
	res, err = oracle.NewHandler(
		k,
	)(
		ctx,
		types.NewMsgReportData(
			42,
			[]types.RawReport{types.NewRawReport(1, 0, []byte("data1"))},
			testapp.Validators[0].ValAddress,
		),
	)
	require.ErrorIs(t, err, types.ErrInvalidReportSize)
	require.Nil(t, res)
	// Incorrect external IDs
	res, err = oracle.NewHandler(
		k,
	)(
		ctx,
		types.NewMsgReportData(
			42,
			[]types.RawReport{types.NewRawReport(1, 0, []byte("data1")), types.NewRawReport(42, 0, []byte("data2"))},
			testapp.Validators[0].ValAddress,
		),
	)
	testapp.CheckErrorf(t, err, types.ErrRawRequestNotFound, "reqID: 42, extID: 42")
	require.Nil(t, res)
	// Request already expired
	err = k.SetRequestLastExpired(ctx, 42)
	require.NoError(t, err)
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Validators[0].ValAddress))
	require.ErrorIs(t, err, types.ErrRequestAlreadyExpired)
	require.Nil(t, res)
}

func TestActivateSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	ctx = ctx.WithBlockTime(testapp.ParseTime(1000000))
	validatorStatus, err := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(t,
		types.NewValidatorStatus(false, time.Time{}),
		validatorStatus,
	)
	msg := types.NewMsgActivate(testapp.Validators[0].ValAddress)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)

	validatorStatus, err = k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.Equal(t,
		types.NewValidatorStatus(true, testapp.ParseTime(1000000)),
		validatorStatus,
	)
	event := abci.Event{
		Type: types.EventTypeActivate,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyValidator, Value: testapp.Validators[0].ValAddress.String()},
		},
	}
	require.Equal(t, event, res.Events[0])
}

func TestActivateFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	msg := types.NewMsgActivate(testapp.Validators[0].ValAddress)
	// Already active.
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.ErrorIs(t, err, types.ErrValidatorAlreadyActive)
	require.Nil(t, res)
	// Too soon to activate.
	ctx = ctx.WithBlockTime(testapp.ParseTime(100000))
	err = k.MissReport(ctx, testapp.Validators[0].ValAddress, testapp.ParseTime(99999))
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(testapp.ParseTime(100001))
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.ErrorIs(t, err, types.ErrTooSoonToActivate)
	require.Nil(t, res)
	// OK
	ctx = ctx.WithBlockTime(testapp.ParseTime(200000))
	_, err = oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
}

func TestUpdateParamsSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	expectedParams := types.Params{
		MaxRawRequestCount:       1,
		MaxAskCount:              10,
		MaxCalldataSize:          256,
		MaxReportDataSize:        512,
		ExpirationBlockCount:     30,
		BaseOwasmGas:             50000,
		PerValidatorRequestGas:   3000,
		SamplingTryCount:         3,
		OracleRewardPercentage:   50,
		InactivePenaltyDuration:  1000,
		IBCRequestEnabled:        true,
		RewardDecreasingFraction: math.LegacyNewDec(10),
	}
	msg := types.NewMsgUpdateParams(k.GetAuthority(), expectedParams)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)

	params, err := k.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedParams, params)
	event := abci.Event{
		Type: types.EventTypeUpdateParams,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyParams, Value: expectedParams.String()},
		},
	}
	require.Equal(t, event, res.Events[0])

	expectedParams = types.Params{
		MaxRawRequestCount:       2,
		MaxAskCount:              20,
		MaxCalldataSize:          512,
		MaxReportDataSize:        256,
		ExpirationBlockCount:     40,
		BaseOwasmGas:             0,
		PerValidatorRequestGas:   0,
		SamplingTryCount:         5,
		OracleRewardPercentage:   0,
		InactivePenaltyDuration:  0,
		IBCRequestEnabled:        false,
		RewardDecreasingFraction: math.LegacyNewDec(10),
	}
	msg = types.NewMsgUpdateParams(k.GetAuthority(), expectedParams)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)

	params, err = k.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedParams, params)
	event = abci.Event{
		Type: types.EventTypeUpdateParams,
		Attributes: []abci.EventAttribute{
			{Key: types.AttributeKeyParams, Value: expectedParams.String()},
		},
	}
	require.Equal(t, event, res.Events[0])
}

func TestUpdateParamsFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	expectedParams := types.Params{
		MaxRawRequestCount:      1,
		MaxAskCount:             10,
		MaxCalldataSize:         256,
		MaxReportDataSize:       512,
		ExpirationBlockCount:    30,
		BaseOwasmGas:            50000,
		PerValidatorRequestGas:  3000,
		SamplingTryCount:        3,
		OracleRewardPercentage:  50,
		InactivePenaltyDuration: 1000,
		IBCRequestEnabled:       true,
	}
	msg := types.NewMsgUpdateParams("foo", expectedParams)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.ErrorContains(t, err, "invalid authority")
	require.Nil(t, res)

	expectedParams = types.Params{
		MaxRawRequestCount:      0,
		MaxAskCount:             10,
		MaxCalldataSize:         256,
		MaxReportDataSize:       512,
		ExpirationBlockCount:    30,
		BaseOwasmGas:            50000,
		PerValidatorRequestGas:  3000,
		SamplingTryCount:        3,
		OracleRewardPercentage:  50,
		InactivePenaltyDuration: 1000,
		IBCRequestEnabled:       true,
	}
	msg = types.NewMsgUpdateParams(k.GetAuthority(), expectedParams)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.ErrorContains(t, err, "max raw request count must be positive")
	require.Nil(t, res)
}
