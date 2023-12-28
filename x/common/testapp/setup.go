package testapp

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/cli"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	owasm "github.com/odin-protocol/go-owasm/api"
	"github.com/spf13/viper"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	odinapp "github.com/ODIN-PROTOCOL/odin-core/app"
	me "github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
)

const (
	DefaultBondDenom   = "loki"
	DefaultHelperDenom = "minigeo"
)

// nolint
var (
	Owner              Account
	Treasury           Account
	FeePayer           Account
	Peter              Account
	Alice              Account
	Bob                Account
	Carol              Account
	OraclePoolProvider Account
	FeePoolProvider    Account
	NotBondedPool      Account
	Validators         []Account
	DataSources        []oracletypes.DataSource
	OracleScripts      []oracletypes.OracleScript
	OwasmVM            *owasm.Vm
)

func init() {
	odinapp.SetBech32AddressPrefixesAndBip44CoinType(sdk.GetConfig())
	r := rand.New(rand.NewSource(time.Now().Unix()))
	Owner = createArbitraryAccount(r)
	Treasury = createArbitraryAccount(r)
	FeePayer = createArbitraryAccount(r)
	Peter = createArbitraryAccount(r)
	Alice = createArbitraryAccount(r)
	Bob = createArbitraryAccount(r)
	Carol = createArbitraryAccount(r)
	OraclePoolProvider = createArbitraryAccount(r)
	FeePoolProvider = createArbitraryAccount(r)
	NotBondedPool = createArbitraryAccount(r)
	for i := 0; i < 3; i++ {
		Validators = append(Validators, createArbitraryAccount(r))
	}
	owasmVM, err := owasm.NewVm(10)
	if err != nil {
		panic(err)
	}
	OwasmVM = owasmVM
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

// NewSimApp creates instance of our app using in test.
func NewSimApp(chainID string, logger log.Logger) *odinapp.OdinApp {
	// Set HomeFlag to a temp folder for simulation run.
	dir, err := os.MkdirTemp("", "odind")
	if err != nil {
		panic(err)
	}
	viper.Set(cli.HomeFlag, dir)

	db := dbm.NewMemDB()
	encCdc := odinapp.MakeEncodingConfig()
	app := odinapp.NewOdinApp(logger, db, nil, true, map[int64]bool{}, dir, 0, encCdc, EmptyAppOptions{}, false, 0, baseapp.SetChainID("ODINCHAIN"))

	genesis := odinapp.NewDefaultGenesisState()
	acc := []authtypes.GenesisAccount{
		&authtypes.BaseAccount{Address: Owner.Address.String()},
		&authtypes.BaseAccount{Address: FeePayer.Address.String()},
		&authtypes.BaseAccount{Address: Peter.Address.String()},
		&authtypes.BaseAccount{Address: Alice.Address.String()},
		&authtypes.BaseAccount{Address: Bob.Address.String()},
		&authtypes.BaseAccount{Address: Validators[0].Address.String()},
		&authtypes.BaseAccount{Address: Validators[1].Address.String()},
		&authtypes.BaseAccount{Address: Validators[2].Address.String()},
		&authtypes.BaseAccount{Address: Carol.Address.String()},
		&authtypes.BaseAccount{Address: OraclePoolProvider.Address.String()},
		&authtypes.BaseAccount{Address: FeePoolProvider.Address.String()},
	}
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), acc)
	genesis[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(Validators))
	signingInfos := make([]slashingtypes.SigningInfo, 0, len(Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(Validators))
	bamt := []sdk.Int{Coins100000000loki[0].Amount, Coins1000000loki[0].Amount, Coins99999999loki[0].Amount}
	bamtSum := 0
	for _, i := range bamt {
		bamtSum += int(i.Int64())
	}

	// bondAmt := sdk.NewInt(1000000)
	for idx, val := range Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
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
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}
		consAddr, err := validator.GetConsAddr()
		validatorSigningInfo := slashingtypes.NewValidatorSigningInfo(consAddr, 0, 0, time.Unix(0, 0), false, 0)
		if err != nil {
			panic(err)
		}

		validators = append(validators, validator)
		signingInfos = append(signingInfos, slashingtypes.SigningInfo{Address: consAddr.String(), ValidatorSigningInfo: validatorSigningInfo})
		delegations = append(delegations, stakingtypes.NewDelegation(acc[4+idx].GetAddress(), val.Address.Bytes(), sdk.OneDec()))
	}
	// set validators and delegations
	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = DefaultBondDenom
	stakingGenesis := stakingtypes.NewGenesisState(stakingParams, validators, delegations)
	genesis[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	slashingParams := slashingtypes.DefaultParams()
	slashingGenesis := slashingtypes.NewGenesisState(slashingParams, signingInfos, nil)
	genesis[slashingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(slashingGenesis)

	// Fund seed accounts and validators with 1000000odin and 100000000odin initially.
	balances := []banktypes.Balance{
		{
			Address: Owner.Address.String(),
			Coins:   Coins1000000loki,
		},
		{Address: FeePayer.Address.String(), Coins: Coins100000000loki},
		{Address: Peter.Address.String(), Coins: Coins1000000loki},
		{Address: Alice.Address.String(), Coins: Coins1000000loki.Add(Coin100000000minigeo)},
		{Address: Bob.Address.String(), Coins: Coins1000000loki},
		{Address: Carol.Address.String(), Coins: Coins1000000loki},
		{Address: OraclePoolProvider.Address.String(), Coins: DefaultDataProvidersPool},
		{Address: FeePoolProvider.Address.String(), Coins: DefaultCommunityPool},
		{Address: Validators[0].Address.String(), Coins: Coins100000000loki},
		{Address: Validators[1].Address.String(), Coins: Coins100000000loki},
		{Address: Validators[2].Address.String(), Coins: Coins100000000loki},
	}
	/*totalSupply := sdk.NewCoins()
	for idx := 0; idx < len(balances); idx++ {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(balances[idx].Coins...)
	}*/ /*
		for idx := 0; idx < len(balances)-len(validators); idx++ {
			// add genesis acc tokens and delegated tokens to total supply
			totalSupply = totalSupply.Add(balances[idx].Coins...)
		}
		for idx := 0; idx < len(validators); idx++ {
			// add genesis acc tokens and delegated tokens to total supply
			totalSupply = totalSupply.Add(balances[idx+len(balances)-len(validators)].Coins.Add(sdk.NewCoin(DefaultBondDenom, bamt[idx]))...)
		}*/
	totalSupply := sdk.NewCoins()
	for idx := 0; idx < len(balances)-len(validators); idx++ {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(balances[idx].Coins...)
	}
	for idx := 0; idx < len(validators); idx++ {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(balances[idx+len(balances)-len(validators)].Coins.Add(sdk.NewCoin("loki", bamt[idx]))...)
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin("loki", sdk.NewInt(int64(bamtSum)))},
	})

	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesis[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	// Add genesis data sources and oracle scripts
	oracleGenesis := oracletypes.DefaultGenesisState()

	oracleGenesis.DataSources = getGenesisDataSources(dir)
	oracleGenesis.OracleScripts = getGenesisOracleScripts(dir)

	genesis[oracletypes.ModuleName] = app.AppCodec().MustMarshalJSON(oracleGenesis)
	stateBytes, err := json.MarshalIndent(genesis, "", " ")

	// Initialize the sim blockchain. We are ready for testing!
	app.InitChain(abci.RequestInitChain{
		ChainId:       chainID,
		Validators:    []abci.ValidatorUpdate{},
		AppStateBytes: stateBytes,
	})

	return app
}

// CreateTestInput creates a new test environment for unit tests.
// params[0] - activate;
// params[1] - fund pools;
// Deprecated
//   - use TestAppBuilder instead
func CreateTestInput(params ...bool) (*odinapp.OdinApp, sdk.Context, me.Keeper) {
	app := NewSimApp("ODINCHAIN", log.NewNopLogger())
	ctx := app.NewContext(false, tmproto.Header{Height: app.LastBlockHeight()})
	if len(params) > 0 && params[0] {
		app.OracleKeeper.Activate(ctx, Validators[0].ValAddress)
		app.OracleKeeper.Activate(ctx, Validators[1].ValAddress)
		app.OracleKeeper.Activate(ctx, Validators[2].ValAddress)
	}

	if len(params) > 1 && params[1] {
		app.DistrKeeper.FundCommunityPool(ctx, DefaultCommunityPool, FeePoolProvider.Address)
		accumulatedPaymentsForData := app.OracleKeeper.GetAccumulatedPaymentsForData(ctx)
		accumulatedPaymentsForData.AccumulatedAmount = accumulatedPaymentsForData.AccumulatedAmount.Add(DefaultDataProvidersPool...)

		app.OracleKeeper.SetAccumulatedPaymentsForData(ctx, accumulatedPaymentsForData)

		ctx = app.NewContext(false, tmproto.Header{})
	}
	return app, ctx, app.OracleKeeper
}
