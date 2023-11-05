package types

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/stretchr/testify/require"
)

var (
	GoodTestAddr    = sdk.AccAddress(make([]byte, 20))
	BadTestAddr     = sdk.AccAddress(make([]byte, 256))
	GoodTestValAddr = sdk.ValAddress(make([]byte, 20))
	BadTestValAddr  = sdk.ValAddress(make([]byte, 256))

	MsgPk            = secp256k1.GenPrivKey().PubKey()
	GoodTestAddr2    = sdk.AccAddress(MsgPk.Address())
	GoodTestValAddr2 = sdk.ValAddress(MsgPk.Address())

	GoodCoins = sdk.NewCoins()
	BadCoins  = []sdk.Coin{{Denom: "odin", Amount: sdk.NewInt(-1)}}
)

type validateTestCase struct {
	valid bool
	msg   sdk.Msg
}

func performValidateTests(t *testing.T, cases []validateTestCase) {
	for _, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestMsgRoute(t *testing.T) {
	require.Equal(t, "oracle", MsgCreateDataSource{}.Route())
	require.Equal(t, "oracle", MsgEditDataSource{}.Route())
	require.Equal(t, "oracle", MsgCreateOracleScript{}.Route())
	require.Equal(t, "oracle", MsgEditOracleScript{}.Route())
	require.Equal(t, "oracle", MsgRequestData{}.Route())
	require.Equal(t, "oracle", MsgReportData{}.Route())
	require.Equal(t, "oracle", MsgActivate{}.Route())
	require.Equal(t, "oracle", MsgAddReporter{}.Route())
	require.Equal(t, "oracle", MsgRemoveReporter{}.Route())
}

func TestMsgType(t *testing.T) {
	require.Equal(t, "create_data_source", MsgCreateDataSource{}.Type())
	require.Equal(t, "edit_data_source", MsgEditDataSource{}.Type())
	require.Equal(t, "create_oracle_script", MsgCreateOracleScript{}.Type())
	require.Equal(t, "edit_oracle_script", MsgEditOracleScript{}.Type())
	require.Equal(t, "request", MsgRequestData{}.Type())
	require.Equal(t, "report", MsgReportData{}.Type())
	require.Equal(t, "activate", MsgActivate{}.Type())
	require.Equal(t, "add_reporter", MsgAddReporter{}.Type())
	require.Equal(t, "remove_reporter", MsgRemoveReporter{}.Type())
}

func TestMsgGetSigners(t *testing.T) {
	signerAcc := sdk.AccAddress([]byte("01234567890123456789"))
	signerVal := sdk.ValAddress([]byte("01234567890123456789"))
	anotherAcc := sdk.AccAddress([]byte("98765432109876543210"))
	anotherVal := sdk.ValAddress([]byte("98765432109876543210"))
	signers := []sdk.AccAddress{signerAcc}
	emptyCoins := sdk.NewCoins()
	require.Equal(t, signers, NewMsgCreateDataSource("name", "desc", []byte("exec"), emptyCoins, anotherAcc, signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgEditDataSource(1, "name", "desc", []byte("exec"), emptyCoins, anotherAcc, signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte("code"), anotherAcc, signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte("code"), anotherAcc, signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", emptyCoins, 1, 1, signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, anotherVal, signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgActivate(signerVal).GetSigners())
	require.Equal(t, signers, NewMsgAddReporter(signerVal, anotherAcc).GetSigners())
	require.Equal(t, signers, NewMsgRemoveReporter(signerVal, anotherAcc).GetSigners())
}

// func TestMsgGetSignBytes(t *testing.T) {
// 	sdk.GetConfig().SetBech32PrefixForAccount("band", "band"+sdk.PrefixPublic)
// 	sdk.GetConfig().SetBech32PrefixForValidator("band"+sdk.PrefixValidator+sdk.PrefixOperator, "band"+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
// 	sdk.GetConfig().SetBech32PrefixForConsensusNode("band"+sdk.PrefixValidator+sdk.PrefixConsensus, "band"+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
// 	require.Equal(t,
// 		`{"type":"oracle/CreateDataSource","value":{"description":"desc","executable":"ZXhlYw==","name":"name","owner":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","sender":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4"}}`,
// 		string(NewMsgCreateDataSource("name", "desc", []byte("exec"), GoodTestAddr, GoodTestAddr).GetSignBytes()),
// 	)
// 	require.Equal(t,
// 		`{"type":"oracle/EditDataSource","value":{"data_source_id":"1","description":"desc","executable":"ZXhlYw==","name":"name","owner":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","sender":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4"}}`,
// 		string(NewMsgEditDataSource(1, "name", "desc", []byte("exec"), GoodTestAddr, GoodTestAddr).GetSignBytes()),
// 	)
// 	require.Equal(t,
// 		`{"type":"oracle/CreateOracleScript","value":{"code":"Y29kZQ==","description":"desc","name":"name","owner":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","schema":"schema","sender":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","source_code_url":"url"}}`,
// 		string(NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte("code"), GoodTestAddr, GoodTestAddr).GetSignBytes()),
// 	)
// 	require.Equal(t,
// 		`{"type":"oracle/EditOracleScript","value":{"code":"Y29kZQ==","description":"desc","name":"name","oracle_script_id":"1","owner":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","schema":"schema","sender":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","source_code_url":"url"}}`,
// 		string(NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte("code"), GoodTestAddr, GoodTestAddr).GetSignBytes()),
// 	)
// 	require.Equal(t,
// 		`{"type":"oracle/Request","value":{"ask_count":"10","calldata":"Y2FsbGRhdGE=","client_id":"client-id","min_count":"5","oracle_script_id":"1","sender":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4"}}`,
// 		string(NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodTestAddr).GetSignBytes()),
// 	)
// 	require.Equal(t,
// 		`{"type":"oracle/Report","value":{"raw_reports":[{"data":"ZGF0YTE=","exit_code":1,"external_id":"1"},{"data":"ZGF0YTI=","exit_code":2,"external_id":"2"}],"reporter":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","request_id":"1","validator":"bandvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqx6y767"}}`,
// 		string(NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, GoodTestValAddr, GoodTestAddr).GetSignBytes()),
// 	)
// 	require.Equal(t,
// 		`{"type":"oracle/Activate","value":{"validator":"bandvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqx6y767"}}`,
// 		string(NewMsgActivate(GoodTestValAddr).GetSignBytes()),
// 	)
// 	require.Equal(t,
// 		`{"type":"oracle/AddReporter","value":{"reporter":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","validator":"bandvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqx6y767"}}`,
// 		string(NewMsgAddReporter(GoodTestValAddr, GoodTestAddr).GetSignBytes()),
// 	)
// 	require.Equal(t,
// 		`{"type":"oracle/RemoveReporter","value":{"reporter":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","validator":"bandvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqx6y767"}}`,
// 		string(NewMsgRemoveReporter(GoodTestValAddr, GoodTestAddr).GetSignBytes()),
// 	)
// }

func TestMsgCreateDataSourceValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgCreateDataSource("name", "desc", []byte("exec"), GoodCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateDataSource("name", "desc", []byte("exec"), GoodCoins, BadTestAddr, GoodTestAddr)},
		{false, NewMsgCreateDataSource("name", "desc", []byte("exec"), GoodCoins, GoodTestAddr, BadTestAddr)},
		{false, NewMsgCreateDataSource("name", "desc", []byte("exec"), BadCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateDataSource(strings.Repeat("x", 200), "desc", []byte("exec"), GoodCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateDataSource("name", strings.Repeat("x", 5000), []byte("exec"), GoodCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateDataSource("name", "desc", []byte{}, GoodCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateDataSource("name", "desc", []byte(strings.Repeat("x", 20000)), GoodCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateDataSource("name", "desc", DoNotModifyBytes, GoodCoins, GoodTestAddr, GoodTestAddr)},
	})
}

func TestMsgEditDataSourceValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgEditDataSource(1, "name", "desc", []byte("exec"), GoodCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditDataSource(1, "name", "desc", []byte("exec"), GoodCoins, BadTestAddr, GoodTestAddr)},
		{false, NewMsgEditDataSource(1, "name", "desc", []byte("exec"), GoodCoins, GoodTestAddr, BadTestAddr)},
		{false, NewMsgEditDataSource(1, "name", "desc", []byte("exec"), BadCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditDataSource(1, strings.Repeat("x", 200), "desc", []byte("exec"), GoodCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditDataSource(1, "name", strings.Repeat("x", 5000), []byte("exec"), GoodCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditDataSource(1, "name", "desc", []byte{}, GoodCoins, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditDataSource(1, "name", "desc", []byte(strings.Repeat("x", 20000)), GoodCoins, GoodTestAddr, GoodTestAddr)},
	})
}

func TestMsgCreateOracleScriptValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte("code"), BadTestAddr, GoodTestAddr)},
		{false, NewMsgCreateOracleScript(strings.Repeat("x", 200), "desc", "schema", "url", []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateOracleScript("name", strings.Repeat("x", 5000), "schema", "url", []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateOracleScript("name", "desc", strings.Repeat("x", 1000), "url", []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateOracleScript("name", "desc", "schema", strings.Repeat("x", 200), []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte{}, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte(strings.Repeat("x", 600000)), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateOracleScript("name", "desc", "schema", "url", DoNotModifyBytes, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte("code"), GoodTestAddr, BadTestAddr)},
	})
}

func TestMsgEditOracleScriptValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte("code"), BadTestAddr, GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, strings.Repeat("x", 200), "desc", "schema", "url", []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, "name", strings.Repeat("x", 5000), "schema", "url", []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, "name", "desc", strings.Repeat("x", 1000), "url", []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, "name", "desc", "schema", strings.Repeat("x", 200), []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte{}, GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte(strings.Repeat("x", 600000)), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte("code"), GoodTestAddr, BadTestAddr)},
	})
}

func TestMsgRequestDataValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodCoins, 1, 1, GoodTestAddr)},
		{true, NewMsgRequestData(1, []byte(strings.Repeat("x", 2000)), 10, 5, "client-id", GoodCoins, 1, 1, GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 2, 5, "client-id", GoodCoins, 1, 1, GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 0, 0, "client-id", GoodCoins, 1, 1, GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, strings.Repeat("x", 300), GoodCoins, 1, 1, GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodCoins, 1, 1, BadTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", BadCoins, 1, 1, GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodCoins, 0, 1, GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodCoins, 1, 0, GoodTestAddr)},
	})
}

func TestMsgReportDataValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, GoodTestValAddr, GoodTestAddr)},
		{false, NewMsgReportData(1, []RawReport{}, GoodTestValAddr, GoodTestAddr)},
		{true, NewMsgReportData(1, []RawReport{{1, 1, []byte(strings.Repeat("x", 500))}, {2, 2, []byte("data2")}}, GoodTestValAddr, GoodTestAddr)},
		{false, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {1, 1, []byte("data2")}}, GoodTestValAddr, GoodTestAddr)},
		{false, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, BadTestValAddr, GoodTestAddr)},
		{false, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, GoodTestValAddr, BadTestAddr)},
	})
}

func TestMsgActivateValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgActivate(GoodTestValAddr)},
		{false, NewMsgActivate(BadTestValAddr)},
	})
}

func TestMsgAddReporterValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgAddReporter(GoodTestValAddr, GoodTestAddr2)},
		{false, NewMsgAddReporter(BadTestValAddr, GoodTestAddr)},
		{false, NewMsgAddReporter(GoodTestValAddr, BadTestAddr)},
		{false, NewMsgAddReporter(GoodTestValAddr, GoodTestAddr)},
	})
}

func TestMsgRemoveReporterValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgRemoveReporter(GoodTestValAddr, GoodTestAddr2)},
		{false, NewMsgRemoveReporter(BadTestValAddr, GoodTestAddr)},
		{false, NewMsgRemoveReporter(GoodTestValAddr, BadTestAddr)},
		{false, NewMsgRemoveReporter(GoodTestValAddr, GoodTestAddr)},
	})
}
