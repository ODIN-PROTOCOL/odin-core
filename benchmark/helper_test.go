package benchmark

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"

	types "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	owasm "github.com/odin-protocol/go-owasm/api"
	"github.com/stretchr/testify/require"

	"github.com/ODIN-PROTOCOL/odin-core/pkg/obi"
	"github.com/ODIN-PROTOCOL/odin-core/testing/testapp"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

type Account struct {
	testapp.Account
	Num uint64
	Seq uint64
}

type BenchmarkCalldata struct {
	DataSourceId uint64
	Scenario     uint64
	Value        uint64
	Text         string
}

func GetBenchmarkWasm() ([]byte, error) {
	oCode, err := ioutil.ReadFile("./testdata/benchmark-oracle-script.wasm")
	return oCode, err
}

func GenMsgRequestData(
	sender *Account,
	oracleScriptId uint64,
	dataSourceId uint64,
	scenario uint64,
	value uint64,
	stringLength int,
	prepareGas uint64,
	executeGas uint64,
) []sdk.Msg {
	msg := oracletypes.MsgRequestData{
		OracleScriptID: oracletypes.OracleScriptID(oracleScriptId),
		Calldata: obi.MustEncode(BenchmarkCalldata{
			DataSourceId: dataSourceId,
			Scenario:     scenario,
			Value:        value,
			Text:         strings.Repeat("#", stringLength),
		}),
		AskCount:   1,
		MinCount:   1,
		ClientID:   "",
		FeeLimit:   sdk.Coins{sdk.NewInt64Coin("loki", 1)},
		PrepareGas: prepareGas,
		ExecuteGas: executeGas,
		Sender:     sender.Address.String(),
	}

	return []sdk.Msg{&msg}
}

func GenMsgSend(
	sender *Account,
	receiver *Account,
) []sdk.Msg {
	msg := banktypes.MsgSend{
		FromAddress: sender.Address.String(),
		ToAddress:   receiver.Address.String(),
		Amount:      sdk.Coins{sdk.NewInt64Coin("loki", 1)},
	}

	return []sdk.Msg{&msg}
}

func GenMsgCreateOracleScript(sender *Account, code []byte) []sdk.Msg {
	msg := oracletypes.MsgCreateOracleScript{
		Name:          "test",
		Description:   "test",
		Schema:        "test",
		SourceCodeURL: "test",
		Code:          code,
		Owner:         sender.Address.String(),
		Sender:        sender.Address.String(),
	}

	return []sdk.Msg{&msg}
}

func GenMsgCreateDataSource(sender *Account, code []byte) []sdk.Msg {
	msg := oracletypes.MsgCreateDataSource{
		Name:        "test",
		Description: "test",
		Executable:  code,
		Fee:         sdk.Coins{},
		Treasury:    sender.Address.String(),
		Owner:       sender.Address.String(),
		Sender:      sender.Address.String(),
	}

	return []sdk.Msg{&msg}
}

func GenMsgActivate(account *Account) []sdk.Msg {
	msg := oracletypes.MsgActivate{
		Validator: account.ValAddress.String(),
	}

	return []sdk.Msg{&msg}
}

func GenSequenceOfTxs(
	txConfig client.TxConfig,
	msgs []sdk.Msg,
	account *Account,
	numTxs int,
) []sdk.Tx {
	txs := make([]sdk.Tx, numTxs)

	for i := 0; i < numTxs; i++ {
		txs[i], _ = testapp.GenTx(
			txConfig,
			msgs,
			sdk.Coins{sdk.NewInt64Coin("loki", 1)},
			math.MaxInt64,
			"",
			[]uint64{account.Num},
			[]uint64{account.Seq},
			account.PrivKey,
		)
		account.Seq += 1
	}

	return txs
}

type Event struct {
	Type       string
	Attributes map[string]string
}

func DecodeEvents(events []types.Event) []Event {
	evs := []Event{}
	for _, event := range events {
		attrs := make(map[string]string, 0)
		for _, attributes := range event.Attributes {
			attrs[string(attributes.Key)] = string(attributes.Value)
		}
		evs = append(evs, Event{
			Type:       event.Type,
			Attributes: attrs,
		})
	}

	return evs
}

func LogEvents(b testing.TB, events []types.Event) {
	evs := DecodeEvents(events)
	for i, ev := range evs {
		b.Logf("Event %d: %+v\n", i, ev)
	}

	if len(evs) == 0 {
		b.Logf("No Event")
	}
}

func GetFirstAttributeOfLastEventValue(events []types.Event) (int, error) {
	evt := events[len(events)-1]
	attr := evt.Attributes[0]
	value, err := strconv.Atoi(string(attr.Value))

	return value, err
}

func InitOwasmTestEnv(
	b testing.TB,
	cacheSize uint32,
	scenario uint64,
	parameter uint64,
	stringLength int,
) (*owasm.Vm, []byte, oracletypes.Request) {
	// prepare owasm vm
	owasmVM, err := owasm.NewVm(cacheSize)
	require.NoError(b, err)

	// prepare owasm code
	oCode, err := GetBenchmarkWasm()
	require.NoError(b, err)
	compiledCode, err := owasmVM.Compile(oCode, oracletypes.MaxCompiledWasmCodeSize)
	require.NoError(b, err)

	// prepare request
	req := oracletypes.NewRequest(
		1, obi.MustEncode(BenchmarkCalldata{
			DataSourceId: 1,
			Scenario:     scenario,
			Value:        parameter,
			Text:         strings.Repeat("#", stringLength),
		}), []sdk.ValAddress{[]byte{}}, 1,
		1, time.Now(), "", nil, nil, ExecuteGasLimit,
	)

	return owasmVM, compiledCode, req
}

func GetConsensusParams(maxGas int64) *tmproto.ConsensusParams {
	return &tmproto.ConsensusParams{
		Block: &tmproto.BlockParams{
			MaxBytes: 200000,
			MaxGas:   maxGas,
		},
		Evidence: &tmproto.EvidenceParams{
			MaxAgeNumBlocks: 302400,
			MaxAgeDuration:  504 * time.Hour,
		},
		Validator: &tmproto.ValidatorParams{
			PubKeyTypes: []string{
				tmtypes.ABCIPubKeyTypeSecp256k1,
			},
		},
	}
}

func ChunkSlice(slice []uint64, chunkSize int) [][]uint64 {
	var chunks [][]uint64
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func GenOracleReports() []oracletypes.Report {
	return []oracletypes.Report{
		{
			Validator:       "",
			InBeforeResolve: true,
			RawReports: []oracletypes.RawReport{
				{
					ExternalID: 0,
					ExitCode:   0,
					Data:       []byte{},
				},
			},
		},
	}
}

func GetSpanSize() uint64 {
	if oracletypes.DefaultMaxReportDataSize > oracletypes.DefaultMaxCalldataSize {
		return oracletypes.DefaultMaxReportDataSize
	}
	return oracletypes.DefaultMaxCalldataSize
}
