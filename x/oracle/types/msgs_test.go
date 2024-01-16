package types

import (
	"strings"
	"testing"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	GoodTestAddr    = sdk.AccAddress(make([]byte, 20))
	EmptyAddr       = sdk.AccAddress([]byte(""))
	GoodTestValAddr = sdk.ValAddress(make([]byte, 20))
	EmptyValAddr    = sdk.ValAddress([]byte(""))

	MsgPk            = secp256k1.GenPrivKey().PubKey()
	GoodTestAddr2    = sdk.AccAddress(MsgPk.Address())
	GoodTestValAddr2 = sdk.ValAddress(MsgPk.Address())

	GoodCoins = sdk.NewCoins()
	BadCoins  = []sdk.Coin{{Denom: "loki", Amount: sdk.NewInt(-1)}}
	FeeCoins  = sdk.NewCoins(sdk.NewCoin("loki", sdk.NewInt(1000)))
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
}

func TestMsgType(t *testing.T) {
	require.Equal(t, "create_data_source", MsgCreateDataSource{}.Type())
	require.Equal(t, "edit_data_source", MsgEditDataSource{}.Type())
	require.Equal(t, "create_oracle_script", MsgCreateOracleScript{}.Type())
	require.Equal(t, "edit_oracle_script", MsgEditOracleScript{}.Type())
	require.Equal(t, "request", MsgRequestData{}.Type())
	require.Equal(t, "report", MsgReportData{}.Type())
	require.Equal(t, "activate", MsgActivate{}.Type())
}

func TestMsgGetSigners(t *testing.T) {
	signerAcc := sdk.AccAddress([]byte("01234567890123456789"))
	signerVal := sdk.ValAddress([]byte("01234567890123456789"))
	anotherAcc := sdk.AccAddress([]byte("98765432109876543210"))
	anotherVal := sdk.ValAddress([]byte("98765432109876543210"))
	treasuryAcc := sdk.AccAddress([]byte("treasury"))
	signers := []sdk.AccAddress{signerAcc}
	emptyCoins := sdk.NewCoins()
	require.Equal(
		t,
		signers,
		NewMsgCreateDataSource(
			"name",
			"desc",
			[]byte("exec"),
			emptyCoins,
			treasuryAcc,
			anotherAcc,
			signerAcc,
		).GetSigners(),
	)
	require.Equal(
		t,
		signers,
		NewMsgEditDataSource(
			1,
			"name",
			"desc",
			[]byte("exec"),
			emptyCoins,
			treasuryAcc,
			anotherAcc,
			signerAcc,
		).GetSigners(),
	)
	require.Equal(
		t,
		signers,
		NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte("code"), anotherAcc, signerAcc).GetSigners(),
	)
	require.Equal(
		t,
		signers,
		NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte("code"), anotherAcc, signerAcc).GetSigners(),
	)
	require.Equal(
		t,
		signers,
		NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", emptyCoins, 1, 1, signerAcc).GetSigners(),
	)
	require.Equal(
		t,
		[]sdk.AccAddress{anotherAcc},
		NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, anotherVal).GetSigners(),
	)
	require.Equal(t, signers, NewMsgActivate(signerVal).GetSigners())
}

func TestMsgGetSignBytes(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount("odin", "odin"+sdk.PrefixPublic)
	sdk.GetConfig().
		SetBech32PrefixForValidator("odin"+sdk.PrefixValidator+sdk.PrefixOperator, "odin"+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	sdk.GetConfig().
		SetBech32PrefixForConsensusNode("odin"+sdk.PrefixValidator+sdk.PrefixConsensus, "odin"+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
	require.Equal(
		t,
		`{"type":"oracle/CreateDataSource","value":{"description":"desc","executable":"ZXhlYw==","fee":[{"amount":"1000","denom":"loki"}],"name":"name","owner":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth","sender":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth","treasury":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth"}}`,
		string(
			NewMsgCreateDataSource(
				"name",
				"desc",
				[]byte("exec"),
				FeeCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			).GetSignBytes(),
		),
	)
	require.Equal(
		t,
		`{"type":"oracle/EditDataSource","value":{"data_source_id":"1","description":"desc","executable":"ZXhlYw==","fee":[{"amount":"1000","denom":"loki"}],"name":"name","owner":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth","sender":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth","treasury":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth"}}`,
		string(
			NewMsgEditDataSource(
				1,
				"name",
				"desc",
				[]byte("exec"),
				FeeCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			).GetSignBytes(),
		),
	)
	require.Equal(
		t,
		`{"type":"oracle/CreateOracleScript","value":{"code":"Y29kZQ==","description":"desc","name":"name","owner":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth","schema":"schema","sender":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth","source_code_url":"url"}}`,
		string(
			NewMsgCreateOracleScript(
				"name",
				"desc",
				"schema",
				"url",
				[]byte("code"),
				GoodTestAddr,
				GoodTestAddr,
			).GetSignBytes(),
		),
	)
	require.Equal(
		t,
		`{"type":"oracle/EditOracleScript","value":{"code":"Y29kZQ==","description":"desc","name":"name","oracle_script_id":"1","owner":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth","schema":"schema","sender":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth","source_code_url":"url"}}`,
		string(
			NewMsgEditOracleScript(
				1,
				"name",
				"desc",
				"schema",
				"url",
				[]byte("code"),
				GoodTestAddr,
				GoodTestAddr,
			).GetSignBytes(),
		),
	)
	require.Equal(
		t,
		`{"type":"oracle/Request","value":{"ask_count":"10","calldata":"Y2FsbGRhdGE=","client_id":"client-id","execute_gas":"250000","fee_limit":[{"amount":"1000","denom":"loki"}],"min_count":"5","oracle_script_id":"1","prepare_gas":"50000","sender":"odin1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqhrxuth"}}`,
		string(
			NewMsgRequestData(
				1,
				[]byte("calldata"),
				10,
				5,
				"client-id",
				FeeCoins,
				50000,
				250000,
				GoodTestAddr,
			).GetSignBytes(),
		),
	)
	require.Equal(
		t,
		`{"type":"oracle/Report","value":{"raw_reports":[{"data":"ZGF0YTE=","exit_code":1,"external_id":"1"},{"data":"ZGF0YTI=","exit_code":2,"external_id":"2"}],"request_id":"1","validator":"odinvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq67jzzr"}}`,
		string(
			NewMsgReportData(
				1,
				[]RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}},
				GoodTestValAddr,
			).GetSignBytes(),
		),
	)
	require.Equal(t,
		`{"type":"oracle/Activate","value":{"validator":"odinvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq67jzzr"}}`,
		string(NewMsgActivate(GoodTestValAddr).GetSignBytes()),
	)
}

func TestMsgCreateDataSourceValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{
			true,
			NewMsgCreateDataSource(
				"name",
				"desc",
				[]byte("exec"),
				GoodCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgCreateDataSource("name", "desc", []byte("exec"), GoodCoins, EmptyAddr, GoodTestAddr, GoodTestAddr),
		},
		{
			false,
			NewMsgCreateDataSource("name", "desc", []byte("exec"), GoodCoins, GoodTestAddr, EmptyAddr, GoodTestAddr),
		},
		{
			false,
			NewMsgCreateDataSource("name", "desc", []byte("exec"), GoodCoins, GoodTestAddr, GoodTestAddr, EmptyAddr),
		},
		{
			false,
			NewMsgCreateDataSource("name", "desc", []byte("exec"), BadCoins, GoodTestAddr, GoodTestAddr, GoodTestAddr),
		},
		{
			false,
			NewMsgCreateDataSource(
				strings.Repeat("x", 200),
				"desc",
				[]byte("exec"),
				GoodCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgCreateDataSource(
				"name",
				strings.Repeat("x", 5000),
				[]byte("exec"),
				GoodCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{false, NewMsgCreateDataSource("name", "desc", []byte{}, GoodCoins, GoodTestAddr, GoodTestAddr, GoodTestAddr)},
		{
			false,
			NewMsgCreateDataSource(
				"name",
				"desc",
				[]byte(strings.Repeat("x", 20000)),
				GoodCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgCreateDataSource(
				"name",
				"desc",
				DoNotModifyBytes,
				GoodCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			),
		},
	})
}

func TestMsgEditDataSourceValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{
			true,
			NewMsgEditDataSource(
				1,
				"name",
				"desc",
				[]byte("exec"),
				GoodCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgEditDataSource(1, "name", "desc", []byte("exec"), GoodCoins, EmptyAddr, GoodTestAddr, GoodTestAddr),
		},
		{
			false,
			NewMsgEditDataSource(1, "name", "desc", []byte("exec"), GoodCoins, GoodTestAddr, EmptyAddr, GoodTestAddr),
		},
		{
			false,
			NewMsgEditDataSource(1, "name", "desc", []byte("exec"), GoodCoins, GoodTestAddr, GoodTestAddr, EmptyAddr),
		},
		{
			false,
			NewMsgEditDataSource(
				1,
				"name",
				"desc",
				[]byte("exec"),
				BadCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgEditDataSource(
				1,
				strings.Repeat("x", 200),
				"desc",
				[]byte("exec"),
				GoodCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgEditDataSource(
				1,
				"name",
				strings.Repeat("x", 5000),
				[]byte("exec"),
				GoodCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgEditDataSource(1, "name", "desc", []byte{}, GoodCoins, GoodTestAddr, GoodTestAddr, GoodTestAddr),
		},
		{
			false,
			NewMsgEditDataSource(
				1,
				"name",
				"desc",
				[]byte(strings.Repeat("x", 20000)),
				GoodCoins,
				GoodTestAddr,
				GoodTestAddr,
				GoodTestAddr,
			),
		},
	})
}

func TestMsgCreateOracleScriptValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte("code"), EmptyAddr, GoodTestAddr)},
		{
			false,
			NewMsgCreateOracleScript(
				strings.Repeat("x", 200),
				"desc",
				"schema",
				"url",
				[]byte("code"),
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgCreateOracleScript(
				"name",
				strings.Repeat("x", 5000),
				"schema",
				"url",
				[]byte("code"),
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgCreateOracleScript(
				"name",
				"desc",
				strings.Repeat("x", 1000),
				"url",
				[]byte("code"),
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgCreateOracleScript(
				"name",
				"desc",
				"schema",
				strings.Repeat("x", 200),
				[]byte("code"),
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{false, NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte{}, GoodTestAddr, GoodTestAddr)},
		{
			false,
			NewMsgCreateOracleScript(
				"name",
				"desc",
				"schema",
				"url",
				[]byte(strings.Repeat("x", 600000)),
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgCreateOracleScript("name", "desc", "schema", "url", DoNotModifyBytes, GoodTestAddr, GoodTestAddr),
		},
		{false, NewMsgCreateOracleScript("name", "desc", "schema", "url", []byte("code"), GoodTestAddr, EmptyAddr)},
	})
}

func TestMsgEditOracleScriptValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte("code"), GoodTestAddr, GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte("code"), EmptyAddr, GoodTestAddr)},
		{
			false,
			NewMsgEditOracleScript(
				1,
				strings.Repeat("x", 200),
				"desc",
				"schema",
				"url",
				[]byte("code"),
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgEditOracleScript(
				1,
				"name",
				strings.Repeat("x", 5000),
				"schema",
				"url",
				[]byte("code"),
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgEditOracleScript(
				1,
				"name",
				"desc",
				strings.Repeat("x", 1000),
				"url",
				[]byte("code"),
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{
			false,
			NewMsgEditOracleScript(
				1,
				"name",
				"desc",
				"schema",
				strings.Repeat("x", 200),
				[]byte("code"),
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{false, NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte{}, GoodTestAddr, GoodTestAddr)},
		{
			false,
			NewMsgEditOracleScript(
				1,
				"name",
				"desc",
				"schema",
				"url",
				[]byte(strings.Repeat("x", 600000)),
				GoodTestAddr,
				GoodTestAddr,
			),
		},
		{false, NewMsgEditOracleScript(1, "name", "desc", "schema", "url", []byte("code"), GoodTestAddr, EmptyAddr)},
	})
}

func TestMsgRequestDataValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodCoins, 1, 1, GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 2, 5, "client-id", GoodCoins, 1, 1, GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 0, 0, "client-id", GoodCoins, 1, 1, GoodTestAddr)},
		{
			false,
			NewMsgRequestData(1, []byte("calldata"), 10, 5, strings.Repeat("x", 300), GoodCoins, 1, 1, GoodTestAddr),
		},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodCoins, 1, 1, EmptyAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", BadCoins, 1, 1, GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodCoins, 0, 1, GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodCoins, 1, 0, GoodTestAddr)},
	})
}

func TestMsgReportDataValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, GoodTestValAddr)},
		{false, NewMsgReportData(1, []RawReport{}, GoodTestValAddr)},
		{false, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {1, 1, []byte("data2")}}, GoodTestValAddr)},
		{false, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, EmptyValAddr)},
	})
}

func TestMsgActivateValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgActivate(GoodTestValAddr)},
		{false, NewMsgActivate(EmptyValAddr)},
	})
}
