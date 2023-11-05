package testapp

import (
	"encoding/json"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	odinapp "github.com/ODIN-PROTOCOL/odin-core/app"
)

type TestAppBuilder interface {
	Build(chainID string, stateBytes []byte, params ...bool) *odinapp.OdinApp
	Codec() codec.Codec
	AddGenesis() TestAppBuilder
	UpdateModules(modulesGenesis map[string]json.RawMessage) TestAppBuilder

	GetAuthBuilder() *AuthBuilder
	GetStakingBuilder() *StakingBuilder
	GetBankBuilder() *BankBuilder
	GetOracleBuilder() *OracleBuilder

	SetAuthBuilder(*AuthBuilder)
	SetStakingBuilder(*StakingBuilder)
	SetBankBuilder(*BankBuilder)
	SetOracleBuilder(*OracleBuilder)
}

type testAppBuilder struct {
	app     *odinapp.OdinApp
	genesis odinapp.GenesisState

	*AuthBuilder
	*StakingBuilder
	*BankBuilder
	*OracleBuilder
}

func (b *testAppBuilder) SetAuthBuilder(builder *AuthBuilder) {
	b.AuthBuilder = builder
}

func (b *testAppBuilder) SetStakingBuilder(builder *StakingBuilder) {
	b.StakingBuilder = builder
}

func (b *testAppBuilder) SetBankBuilder(builder *BankBuilder) {
	b.BankBuilder = builder
}

func (b *testAppBuilder) SetOracleBuilder(builder *OracleBuilder) {
	b.OracleBuilder = builder
}

func (b *testAppBuilder) GetAuthBuilder() *AuthBuilder {
	return b.AuthBuilder
}

func (b *testAppBuilder) GetStakingBuilder() *StakingBuilder {
	return b.StakingBuilder
}

func (b *testAppBuilder) GetBankBuilder() *BankBuilder {
	return b.BankBuilder
}

func (b *testAppBuilder) GetOracleBuilder() *OracleBuilder {
	return b.OracleBuilder
}

func NewTestAppBuilder(dir string, logger log.Logger) TestAppBuilder {
	builder := testAppBuilder{}

	db := dbm.NewMemDB()
	encCdc := odinapp.MakeEncodingConfig()
	builder.app = odinapp.NewOdinApp(logger, db, nil, true, map[int64]bool{}, dir, 0, encCdc, EmptyAppOptions{}, false, 0)
	return &builder
}

func (b *testAppBuilder) Codec() codec.Codec {
	return b.app.AppCodec()
}

func (b *testAppBuilder) Build(chainID string, stateBytes []byte, params ...bool) *odinapp.OdinApp {
	stateBytesNew := stateBytes
	if stateBytes == nil {
		stateBytesNew, _ = json.MarshalIndent(b.genesis, "", " ")
	}
	// Initialize the sim blockchain. We are ready for testing!
	b.app.InitChain(abci.RequestInitChain{
		ChainId:       chainID,
		Validators:    []abci.ValidatorUpdate{},
		AppStateBytes: stateBytesNew,
	})
	return b.app
}

func (b *testAppBuilder) AddGenesis() TestAppBuilder {
	b.genesis = odinapp.NewDefaultGenesisState()
	return b
}

func (b *testAppBuilder) UpdateModules(modulesGenesis map[string]json.RawMessage) TestAppBuilder {
	for k, v := range modulesGenesis {
		if v != nil {
			b.genesis[k] = v
		}
	}
	return b
}
