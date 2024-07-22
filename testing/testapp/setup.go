package testapp

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store/snapshots"
	snapshottypes "cosmossdk.io/store/snapshots/types"
	owasm "github.com/ODIN-PROTOCOL/wasmvm/v2"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"github.com/stretchr/testify/require"

	odinapp "github.com/ODIN-PROTOCOL/odin-core/app"
	"github.com/ODIN-PROTOCOL/odin-core/pkg/filecache"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// Account is a data structure to store key of test account.
type Account struct {
	PrivKey    cryptotypes.PrivKey
	PubKey     cryptotypes.PubKey
	Address    sdk.AccAddress
	ValAddress sdk.ValAddress
}

type Validator struct {
	PrivKey    crypto.PrivKey
	PubKey     crypto.PubKey
	Address    sdk.AccAddress
	ValAddress sdk.ValAddress
}

// nolint
var (
	Owner           Account
	Treasury        Account
	FeePayer        Account
	Alice           Account
	Bob             Account
	Carol           Account
	FeePoolProvider Account
	Validators      []Account
	DataSources     []types.DataSource
	OracleScripts   []types.OracleScript
	OwasmVM         *owasm.Vm
)

// nolint
var (
	EmptyCoins               = sdk.Coins(nil)
	Coins1loki               = sdk.NewCoins(sdk.NewInt64Coin("loki", 1))
	Coins10loki              = sdk.NewCoins(sdk.NewInt64Coin("loki", 10))
	Coins11loki              = sdk.NewCoins(sdk.NewInt64Coin("loki", 11))
	Coins1000000minigeo      = sdk.NewCoins(sdk.NewInt64Coin("minigeo", 1000000))
	Coin100000000minigeo     = sdk.NewInt64Coin("minigeo", 100000000)
	Coins1000000loki         = sdk.NewCoins(sdk.NewInt64Coin("loki", 1000000))
	Coin100000000loki        = sdk.NewInt64Coin("loki", 100000000)
	Coins99999999loki        = sdk.NewCoins(sdk.NewInt64Coin("loki", 99999999))
	Coins100000000loki       = sdk.NewCoins(sdk.NewInt64Coin("loki", 100000000))
	Coins10000000000loki     = sdk.NewCoins(sdk.NewInt64Coin("loki", 10000000000))
	BadCoins                 = []sdk.Coin{{Denom: "loki", Amount: math.NewInt(-1)}}
	Port1                    = "port-1"
	Port2                    = "port-2"
	Channel1                 = "channel-1"
	Channel2                 = "channel-2"
	DefaultCommunityPool     = sdk.NewCoins(Coin100000000minigeo, Coin100000000loki)
	DefaultDataProvidersPool = sdk.NewCoins(Coin100000000loki)
)

const (
	TestDefaultPrepareGas uint64 = 40000
	TestDefaultExecuteGas uint64 = 300000
)

// DefaultConsensusParams defines the default Tendermint consensus params used in TestingApp.
var DefaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		// MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeSecp256k1,
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

type TestingApp struct {
	*odinapp.OdinApp
}

func (app *TestingApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetStakingKeeper implements the TestingApp interface.
func (app *TestingApp) GetStakingKeeper() *stakingkeeper.Keeper {
	return app.StakingKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *TestingApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *TestingApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

// GetTxConfig implements the TestingApp interface.
func (app *TestingApp) GetTxConfig() client.TxConfig {
	return odinapp.MakeEncodingConfig().TxConfig
}

func init() {
	odinapp.SetBech32AddressPrefixesAndBip44CoinTypeAndSeal(sdk.GetConfig())
	r := rand.New(rand.NewSource(time.Now().Unix()))
	Owner = createArbitraryAccount(r)
	Treasury = createArbitraryAccount(r)
	FeePayer = createArbitraryAccount(r)
	Alice = createArbitraryAccount(r)
	Bob = createArbitraryAccount(r)
	Carol = createArbitraryAccount(r)
	for i := 0; i < 3; i++ {
		Validators = append(Validators, createArbitraryAccount(r))
	}

	// Sorted list of validators is needed for ibctest when signing a commit block
	sort.Slice(Validators, func(i, j int) bool {
		return Validators[i].PubKey.Address().String() < Validators[j].PubKey.Address().String()
	})

	owasmVM, err := owasm.NewOracleVm(10)
	if err != nil {
		panic(err)
	}
	OwasmVM = owasmVM
}

func createArbitraryAccount(r *rand.Rand) Account {
	privkeySeed := make([]byte, 12)
	r.Read(privkeySeed)
	privKey := secp256k1.GenPrivKeyFromSecret(privkeySeed)
	return Account{
		PrivKey:    privKey,
		PubKey:     privKey.PubKey(),
		Address:    sdk.AccAddress(privKey.PubKey().Address()),
		ValAddress: sdk.ValAddress(privKey.PubKey().Address()),
	}
}

func getGenesisDataSources(homePath, denom string) []types.DataSource {
	dir := filepath.Join(homePath, "files")
	fc := filecache.New(dir)
	DataSources = []types.DataSource{{}} // 0th index should be ignored
	for idx := 0; idx < 5; idx++ {
		idxStr := fmt.Sprintf("%d", idx+1)
		hash := fc.AddFile([]byte("code" + idxStr))
		DataSources = append(DataSources, types.NewDataSource(
			Owner.Address, "name"+idxStr, "desc"+idxStr, hash, sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(1000000))), Treasury.Address,
		))
	}
	return DataSources[1:]
}

func getGenesisOracleScripts(homePath string) []types.OracleScript {
	dir := filepath.Join(homePath, "files")
	fc := filecache.New(dir)
	OracleScripts = []types.OracleScript{{}} // 0th index should be ignored
	wasms := [][]byte{
		Wasm1, Wasm2, Wasm3, Wasm4, Wasm56(10), Wasm56(10000000), Wasm78(10), Wasm78(2000), Wasm9,
	}
	for idx := 0; idx < len(wasms); idx++ {
		idxStr := fmt.Sprintf("%d", idx+1)
		hash := fc.AddFile(compile(wasms[idx]))
		OracleScripts = append(OracleScripts, types.NewOracleScript(
			Owner.Address, "name"+idxStr, "desc"+idxStr, hash, "schema"+idxStr, "url"+idxStr,
		))
	}
	return OracleScripts[1:]
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

// NewTestApp creates instance of our app using in test.
func NewTestApp(chainID string, logger log.Logger) *TestingApp {
	// Set HomeFlag to a temp folder for simulation run.
	dir, err := os.MkdirTemp("", "odind")
	if err != nil {
		panic(err)
	}
	db := dbm.NewMemDB()
	// db, _ := dbm.NewGoLevelDB("db", dir)

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = dir
	appOptions[server.FlagInvCheckPeriod] = 0

	snapshotDir := filepath.Join(dir, "data", "snapshots")
	snapshotDB, err := dbm.NewDB("metadata", dbm.GoLevelDBBackend, snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	snapshotOptions := snapshottypes.NewSnapshotOptions(
		1000,
		2,
	)

	app := &TestingApp{
		OdinApp: odinapp.NewOdinApp(
			log.NewNopLogger(),
			db,
			nil,
			true,
			map[int64]bool{},
			appOptions,
			100,
			baseapp.SetSnapshot(snapshotStore, snapshotOptions),
			baseapp.SetChainID(chainID),
		),
	}
	genesis := odinapp.NewDefaultGenesisState()
	acc := []authtypes.GenesisAccount{
		&authtypes.BaseAccount{Address: Owner.Address.String()},
		&authtypes.BaseAccount{Address: FeePayer.Address.String()},
		&authtypes.BaseAccount{Address: Alice.Address.String()},
		&authtypes.BaseAccount{Address: Bob.Address.String()},
		&authtypes.BaseAccount{Address: Carol.Address.String()},
		&authtypes.BaseAccount{Address: Validators[0].Address.String()},
		&authtypes.BaseAccount{Address: Validators[1].Address.String()},
		&authtypes.BaseAccount{Address: Validators[2].Address.String()},
	}
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), acc)
	genesis[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(Validators))
	signingInfos := make([]slashingtypes.SigningInfo, 0, len(Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(Validators))
	bamt := []math.Int{Coins100000000loki[0].Amount, Coins1000000loki[0].Amount, Coins99999999loki[0].Amount}
	// bondAmt := sdk.NewInt(1000000)
	for idx, val := range Validators {
		tmpk, err := cryptocodec.ToCmtPubKeyInterface(val.PubKey)
		if err != nil {
			panic(err)
		}
		pk, err := cryptocodec.FromCmtPubKeyInterface(tmpk)
		if err != nil {
			panic(err)
		}
		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			panic(err)
		}
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bamt[idx],
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		}
		consAddr, err := validator.GetConsAddr()
		validatorSigningInfo := slashingtypes.NewValidatorSigningInfo(consAddr, 0, 0, time.Unix(0, 0), false, 0)
		if err != nil {
			panic(err)
		}
		validators = append(validators, validator)
		signingInfos = append(
			signingInfos,
			slashingtypes.SigningInfo{Address: sdk.ConsAddress(consAddr).String(), ValidatorSigningInfo: validatorSigningInfo},
		)
		delegations = append(
			delegations,
			stakingtypes.NewDelegation(acc[4+idx].GetAddress().String(), sdk.ValAddress(val.Address).String(), math.LegacyOneDec()),
		)
	}
	// set validators and delegations
	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = "loki"
	stakingGenesis := stakingtypes.NewGenesisState(stakingParams, validators, delegations)
	genesis[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	slashingParams := slashingtypes.DefaultParams()
	slashingGenesis := slashingtypes.NewGenesisState(slashingParams, signingInfos, nil)
	genesis[slashingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(slashingGenesis)

	// Fund seed accounts and validators with 1000000loki and 100000000loki initially.
	balances := []banktypes.Balance{
		{
			Address: Owner.Address.String(),
			Coins:   Coins1000000loki,
		},
		{Address: FeePayer.Address.String(), Coins: Coins100000000loki},
		{Address: Alice.Address.String(), Coins: Coins1000000loki.Add(Coin100000000minigeo)},
		{Address: Bob.Address.String(), Coins: Coins1000000loki},
		{Address: Carol.Address.String(), Coins: Coins1000000loki.Add(Coin100000000minigeo)},
		{Address: Validators[0].Address.String(), Coins: Coins100000000loki},
		{Address: Validators[1].Address.String(), Coins: Coins100000000loki},
		{Address: Validators[2].Address.String(), Coins: Coins100000000loki},
	}
	totalSupply := sdk.NewCoins()
	for idx := 0; idx < len(balances)-len(validators); idx++ {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(balances[idx].Coins...)
	}
	for idx := 0; idx < len(validators); idx++ {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(
			balances[idx+len(balances)-len(validators)].Coins.Add(sdk.NewCoin("loki", bamt[idx]))...)
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin("loki", math.NewInt(200999999))},
	})

	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
		[]banktypes.SendEnabled{},
	)
	genesis[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	// Add genesis data sources and oracle scripts
	oracleGenesis := types.DefaultGenesisState()
	oracleGenesis.DataSources = getGenesisDataSources(dir, "loki")
	oracleGenesis.OracleScripts = getGenesisOracleScripts(dir)
	genesis[types.ModuleName] = app.AppCodec().MustMarshalJSON(oracleGenesis)
	stateBytes, err := json.MarshalIndent(genesis, "", " ")
	if err != nil {
		panic(err)
	}

	// Initialize the sim blockchain. We are ready for testing!
	_, err = app.InitChain(&abci.RequestInitChain{
		ChainId:         chainID,
		Validators:      []abci.ValidatorUpdate{},
		ConsensusParams: DefaultConsensusParams,
		AppStateBytes:   stateBytes,
	})
	if err != nil {
		panic(err)
	}

	hash, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: app.LastBlockHeight() + 1,
		Hash:   hash,
	})
	if err != nil {
		panic(err)
	}

	return app
}

// CreateTestInput creates a new test environment for unit tests.
func CreateTestInput(params ...bool) (*TestingApp, sdk.Context, keeper.Keeper) {
	app := NewTestApp("ODINCHAIN", log.NewNopLogger())
	ctx := app.NewContextLegacy(false, tmproto.Header{Height: app.LastBlockHeight()})
	if len(params) > 0 && params[0] {
		app.OracleKeeper.Activate(ctx, Validators[0].ValAddress)
		app.OracleKeeper.Activate(ctx, Validators[1].ValAddress)
		app.OracleKeeper.Activate(ctx, Validators[2].ValAddress)
	}

	if len(params) > 1 && params[1] {
		app.DistrKeeper.FundCommunityPool(ctx, DefaultCommunityPool, FeePoolProvider.Address)
		accumulatedPaymentsForData, err := app.OracleKeeper.GetAccumulatedPaymentsForData(ctx)
		if err != nil {
			panic(err)
		}
		accumulatedPaymentsForData.AccumulatedAmount = accumulatedPaymentsForData.AccumulatedAmount.Add(DefaultDataProvidersPool...)

		err = app.OracleKeeper.SetAccumulatedPaymentsForData(ctx, accumulatedPaymentsForData)
		if err != nil {
			panic(err)
		}

		ctx = app.NewContext(false)
	}

	if len(params) > 2 && params[2] {
		mintParams, err := app.MintKeeper.GetParams(ctx)
		if err != nil {
			panic(err)
		}

		mintParams.MintAir = true
		err = app.MintKeeper.SetParams(ctx, mintParams)
		if err != nil {
			panic(err)
		}
	}

	return app, ctx, app.OracleKeeper
}

func setup(withGenesis bool, invCheckPeriod uint, chainID string) (*TestingApp, odinapp.GenesisState, string) {
	dir, err := os.MkdirTemp("", "odinibc")
	if err != nil {
		panic(err)
	}
	db := dbm.NewMemDB()

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = dir
	appOptions[server.FlagInvCheckPeriod] = 0

	snapshotDir := filepath.Join(dir, "data", "snapshots")
	snapshotDB, err := dbm.NewDB("metadata", dbm.GoLevelDBBackend, snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	snapshotOptions := snapshottypes.NewSnapshotOptions(
		1000,
		2,
	)

	app := &TestingApp{
		OdinApp: odinapp.NewOdinApp(
			log.NewNopLogger(),
			db,
			nil,
			true,
			map[int64]bool{},
			appOptions,
			0,
			baseapp.SetSnapshot(snapshotStore, snapshotOptions),
			baseapp.SetChainID(chainID),
		),
	}
	if withGenesis {
		return app, odinapp.NewDefaultGenesisState(), dir
	}
	return app, odinapp.GenesisState{}, dir
}

// SetupWithEmptyStore setup a TestingApp instance with empty DB
func SetupWithEmptyStore() *TestingApp {
	app, _, _ := setup(false, 0, "ODINCHAIN")
	return app
}

// SetupWithGenesisValSet initializes a new TestingApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the OdinChain from first genesis
// account. A Nop logger is set in TestingApp.
func SetupWithGenesisValSet(
	t *testing.T,
	valSet *tmtypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	chainID string,
	powerReduction math.Int,
	balances ...banktypes.Balance,
) *TestingApp {
	app, genesisState, dir := setup(true, 5, chainID)

	// ensure baseapp has a chain-id set before running InitChain
	baseapp.SetChainID(chainID)(app.GetBaseApp())

	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.TokensFromConsensusPower(1, powerReduction)

	for i, val := range valSet.Validators {
		pk, err := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(Validators[i].Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		}

		validators = append(validators, validator)
		delegations = append(
			delegations,
			stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(Validators[i].Address).String(), math.LegacyOneDec()),
		)
	}

	// set validators and delegations
	var stakingGenesis stakingtypes.GenesisState
	app.AppCodec().MustUnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesis)

	stakingGenesis.Params.BondDenom = "loki"
	bondDenom := stakingGenesis.Params.BondDenom

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(bondDenom, bondAmt.Mul(math.NewInt(int64(len(valSet.Validators)))))},
	})

	// set validators and delegations
	stakingGenesis = *stakingtypes.NewGenesisState(stakingGenesis.Params, validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(&stakingGenesis)

	// update total supply
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		sdk.NewCoins(), // TODO: check validity
		[]banktypes.Metadata{},
		[]banktypes.SendEnabled{},
	)
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	// Add genesis data sources and oracle scripts
	oracleGenesis := types.DefaultGenesisState()
	oracleGenesis.DataSources = getGenesisDataSources(dir, "minigeo")
	oracleGenesis.OracleScripts = getGenesisOracleScripts(dir)
	genesisState[types.ModuleName] = app.AppCodec().MustMarshalJSON(oracleGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	_, err = app.InitChain(
		&abci.RequestInitChain{
			ChainId:         chainID,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simtestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)
	require.NoError(t, err)

	return app
}

const (
	DefaultGenTxGas = 1000000
)

/*
// GenTx generates a signed mock transaction.
func GenTx(
	gen client.TxConfig,
	msgs []sdk.Msg,
	feeAmt sdk.Coins,
	gas uint64,
	chainID string,
	accNums, accSeqs []uint64,
	priv ...cryptotypes.PrivKey,
) (sdk.Tx, error) {
	sigs := make([]signing.SignatureV2, len(priv))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))

	signMode, err := authsign.APISignModeToInternal(gen.SignModeHandler().DefaultMode())
	if err != nil {
		return nil, err
	}

	// 1st round: set SignatureV2 with empty signatures, to set correct
	// signer infos.
	for i, p := range priv {
		sigs[i] = signing.SignatureV2{
			PubKey: p.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode: signMode,
			},
			Sequence: accSeqs[i],
		}
	}

	tx := gen.NewTxBuilder()
	err = tx.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	err = tx.SetSignatures(sigs...)
	if err != nil {
		return nil, err
	}
	tx.SetMemo(memo)
	tx.SetFeeAmount(feeAmt)
	tx.SetGasLimit(gas)

	// 2nd round: once all signer infos are set, every signer can sign.
	for i, p := range priv {
		signerData := authsign.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		signBytes, err := gen.SignModeHandler().GetSignBytes(signMode, signerData, tx.GetTx())
		if err != nil {
			panic(err)
		}
		sig, err := p.Sign(signBytes)
		if err != nil {
			panic(err)
		}
		sigs[i].Data.(*signing.SingleSignatureData).Signature = sig
		err = tx.SetSignatures(sigs...)
		if err != nil {
			panic(err)
		}
	}

	return tx.GetTx(), nil
}

// SignAndDeliver signs and delivers a transaction. No simulation occurs as the
// ibc testing package causes checkState and deliverState to diverge in block time.
func SignAndDeliver(
	t *testing.T, txCfg client.TxConfig, app *baseapp.BaseApp, header tmproto.Header, msgs []sdk.Msg,
	chainID string, accNums, accSeqs []uint64, priv ...cryptotypes.PrivKey,
) (sdk.GasInfo, *sdk.Result, error) {
	tx, err := GenTx(
		txCfg,
		msgs,
		sdk.Coins{sdk.NewInt64Coin("loki", 2500)},
		DefaultGenTxGas,
		chainID,
		accNums,
		accSeqs,
		priv...,
	)
	require.NoError(t, err)

	// Simulate a sending a transaction and committing a block
	app.BeginBlock(abci.RequestBeginBlock{Header: header, Hash: header.AppHash})
	gInfo, res, err := app.SimDeliver(txCfg.TxEncoder(), tx)

	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	return gInfo, res, err
}
*/
