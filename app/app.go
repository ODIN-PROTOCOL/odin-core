package odin

import (
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	// "github.com/althea-net/bech32-ibc/x/bech32ibc"

	// "github.com/althea-net/bech32-ibc/x/bech32ibc"
	// bech32ibckeeper "github.com/althea-net/bech32-ibc/x/bech32ibc/keeper"
	// bech32ibctypes "github.com/althea-net/bech32-ibc/x/bech32ibc/types"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"

	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	"github.com/gorilla/mux"
	owasm "github.com/odin-protocol/go-owasm/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	distr "github.com/cosmos/cosmos-sdk/x/distribution"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	odinappparams "github.com/ODIN-PROTOCOL/odin-core/app/params"
	"github.com/ODIN-PROTOCOL/odin-core/x/auction"
	auctionkeeper "github.com/ODIN-PROTOCOL/odin-core/x/auction/keeper"
	auctiontypes "github.com/ODIN-PROTOCOL/odin-core/x/auction/types"
	odinbank "github.com/ODIN-PROTOCOL/odin-core/x/bank"
	bandbankkeeper "github.com/ODIN-PROTOCOL/odin-core/x/bank/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/coinswap"
	coinswapkeeper "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/keeper"
	coinswaptypes "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"

	nodeservice "github.com/ODIN-PROTOCOL/odin-core/client/grpc/node"
	odinmint "github.com/ODIN-PROTOCOL/odin-core/x/mint"
	odinmintkeeper "github.com/ODIN-PROTOCOL/odin-core/x/mint/keeper"
	odinminttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle"
	bandante "github.com/ODIN-PROTOCOL/odin-core/x/oracle/ante"
	oraclekeeper "github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
	"github.com/ODIN-PROTOCOL/odin-core/x/telemetry"
	telemetrykeeper "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/keeper"
	telemetrytypes "github.com/ODIN-PROTOCOL/odin-core/x/telemetry/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	ibcfeekeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"

	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

const (
	appName          = "OdinApp"
	Bech32MainPrefix = "odin"
	Bip44CoinType    = 118
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		odinmint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
				ibcclientclient.UpdateClientProposalHandler,
				ibcclientclient.UpgradeProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		vesting.AppModuleBasic{},
		oracle.AppModuleBasic{},
		coinswap.AppModuleBasic{},
		auction.AppModuleBasic{},
		telemetry.AppModuleBasic{},
		transfer.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		ica.AppModuleBasic{},
		groupmodule.AppModuleBasic{},
		consensus.AppModuleBasic{},
		wasm.AppModuleBasic{},
		ibctm.AppModuleBasic{},
	)
	// module account permissions
	maccPerms = map[string][]string{
		oracletypes.ModuleName:         nil,
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		odinminttypes.ModuleName:       {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		transfertypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:            nil,
		wasm.ModuleName:                {authtypes.Burner},
	}
	// module accounts that are allowed to receive tokens.
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}
)

var (
	_ runtime.AppI            = (*OdinApp)(nil)
	_ servertypes.Application = (*OdinApp)(nil)
)

// OdinApp is the application of BandChain, extended base ABCI application.
type OdinApp struct {
	*baseapp.BaseApp

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint
	// keys to access the substores.
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bandbankkeeper.WrappedBankKeeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            odinmintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper          ibcfeekeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	OracleKeeper          oraclekeeper.Keeper
	CoinswapKeeper        coinswapkeeper.Keeper
	AuctionKeeper         auctionkeeper.Keeper
	TelemetryKeeper       telemetrykeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	ICAHostKeeper         icahostkeeper.Keeper
	WasmKeeper            wasmkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedOracleKeeper   capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper  capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper

	// Module manager.
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// Deliver context, set during InitGenesis/BeginBlock and cleared during Commit. It allows
	// anyone with access to OdinApp to read/mutate consensus state anytime. USE WITH CARE!
	DeliverContext sdk.Context

	// List of hooks
	hooks []Hook

	// the configurator
	configurator module.Configurator
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		stdlog.Println("Failed to get home dir %2", err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".odin")
}

// SetBech32AddressPrefixesAndBip44CoinTypeAndSeal sets the global Bech32 prefixes and HD wallet coin type and seal config.
func SetBech32AddressPrefixesAndBip44CoinTypeAndSeal(config *sdk.Config) {
	accountPrefix := Bech32MainPrefix
	validatorPrefix := Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	consensusPrefix := Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	config.SetCoinType(Bip44CoinType)
	config.SetBech32PrefixForAccount(accountPrefix, accountPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(validatorPrefix, validatorPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(consensusPrefix, consensusPrefix+sdk.PrefixPublic)

	config.Seal()
}

type MyAppVersionGetter struct {
	App *OdinApp
}

func (m MyAppVersionGetter) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return strconv.FormatUint((m.App.BaseApp.AppVersion()), 10), true
}

// SetBech32AddressPrefixesAndBip44CoinType sets the global Bech32 prefixes and HD wallet coin type.
func SetBech32AddressPrefixesAndBip44CoinType(config *sdk.Config) {
	accountPrefix := Bech32MainPrefix
	validatorPrefix := Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	consensusPrefix := Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	config.SetBech32PrefixForAccount(accountPrefix, accountPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(validatorPrefix, validatorPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(consensusPrefix, consensusPrefix+sdk.PrefixPublic)
	config.SetCoinType(Bip44CoinType)
}

// NewOdinApp returns a reference to an initialized OdinApp.
func NewOdinApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig odinappparams.EncodingConfig, appOpts servertypes.AppOptions,
	disableFeelessReports bool, owasmCacheSize uint32, baseAppOptions ...func(*baseapp.BaseApp),
) *OdinApp {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		odinminttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, capabilitytypes.StoreKey, oracletypes.StoreKey,
		coinswaptypes.StoreKey, auctiontypes.StoreKey, transfertypes.StoreKey,
		feegrant.StoreKey, authzkeeper.StoreKey, icahosttypes.StoreKey, consensusparamtypes.StoreKey,
		crisistypes.StoreKey, ibcexported.StoreKey, group.StoreKey,
		wasm.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &OdinApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}
	owasmVM, err := owasm.NewVm(owasmCacheSize)
	if err != nil {
		panic(err)
	}
	// Initialize params keeper and module subspaces.
	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		keys[consensusparamtypes.StoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	bApp.SetParamStore(&app.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(transfertypes.ModuleName)
	scopedOracleKeeper := app.CapabilityKeeper.ScopeToModule(oracletypes.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)

	// Add keepers.
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], authtypes.ProtoBaseAccount, maccPerms, Bech32MainPrefix, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	// wrappedBankerKeeper overrides burn token behavior to instead transfer to community pool.
	app.BankKeeper = bandbankkeeper.NewWrappedBankKeeperBurnToCommunityPool(
		bankkeeper.NewBaseKeeper(
			appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.BlockedAddrs(), authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		),
		app.AccountKeeper,
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	groupConfig := group.DefaultConfig()
	/*
		Example of setting group params:
		groupConfig.MaxMetadataLen = 1000
	*/
	app.GroupKeeper = groupkeeper.NewKeeper(
		keys[group.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		groupConfig,
	)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, authtypes.NewModuleAddress(stakingtypes.ModuleName).String(),
	)

	app.MintKeeper = odinmintkeeper.NewKeeper(appCodec, keys[odinminttypes.StoreKey], app.GetSubspace(odinminttypes.ModuleName), app.StakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.AccountKeeper, app.BankKeeper,
		app.StakingKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(distrtypes.ModuleName).String(),
	)

	// DistrKeeper must be set afterward due to the circular reference between banker-staking-distr.
	app.BankKeeper.SetDistrKeeper(&app.DistrKeeper)
	app.BankKeeper.SetMintKeeper(&app.MintKeeper)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, legacyAmino, keys[slashingtypes.StoreKey], app.StakingKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		keys[crisistypes.StoreKey],
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)

	//app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// upgrade handlers
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())

	// app.UpgradeKeeper.SetUpgradeHandler("v0.5.5", func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	// 	var pz odinminttypes.Params
	// 	for _, pair := range pz.ParamSetPairs() {
	// 		if bytes.Equal(pair.Key, odinminttypes.KeyAllowedMinter) {
	// 			pz.AllowedMinter = make([]string, 0)
	// 		} else if bytes.Equal(pair.Key, odinminttypes.KeyAllowedMintDenoms) {
	// 			pz.AllowedMintDenoms = make([]*odinminttypes.AllowedDenom, 0)
	// 		} else if bytes.Equal(pair.Key, odinminttypes.KeyMaxAllowedMintVolume) {
	// 			pz.MaxAllowedMintVolume = sdk.Coins{}
	// 		} else {
	// 			app.GetSubspace(odinminttypes.ModuleName).Get(ctx, pair.Key, pair.Value)
	// 		}
	// 	}
	// 	app.MintKeeper.SetParams(ctx, pz)

	// 	minter := app.MintKeeper.GetMinter(ctx)
	// 	minter.CurrentMintVolume = sdk.Coins{}
	// 	app.MintKeeper.SetMinter(ctx, minter)

	// 	return fromVM, nil
	// })

	// create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)
	// IBC Fee Module keeper
	app.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec, keys[ibcfeetypes.StoreKey],
		app.IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper, app.AccountKeeper, app.BankKeeper,
	)

	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[transfertypes.StoreKey], app.GetSubspace(transfertypes.ModuleName), app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper, app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)
	transferModuleIBC := transfer.NewIBCModule(app.TransferKeeper)

	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)

	icaModule := ica.NewAppModule(nil, &app.ICAHostKeeper)
	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	wasmDir := filepath.Join(homePath, "wasm")

	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "iterator,staking,stargate"
	wasmOpts := GetWasmOpts(appOpts)

	// app.WasmKeeper = wasmkeeper.NewKeeper(appCodec, authtypes.NewModuleAddress(wasmtypes.ModuleName).String(),
	// 	app.AccountKeeper,
	// 	app.BankKeeper,
	// 	app.StakingKeeper,
	// 	app.DistrKeeper, wasmRouter, app.MsgServiceRouter(),
	// 	app.GRPCQueryRouter(), wasmDir, wasmConfig, supportedFeatures, nil, nil, wasmOpts...)

	app.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		keys[wasmtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		distrkeeper.NewQuerier(app.DistrKeeper),
		app.IBCFeeKeeper, // ISC4 Wrapper: fee IBC middleware
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)

	// register the proposal types.
	// govRouter := govtypes.NewRouter()
	// govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
	// 	AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
	// 	AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
	// 	AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
	// 	AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
	// 	AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, wasm.EnableAllProposals))

	// app.GovKeeper = govkeeper.NewKeeper(
	// 	appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
	// 	stakingKeeper, govRouter,
	// )

	govConfig := govtypes.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	app.OracleKeeper = oraclekeeper.NewKeeper(
		appCodec, keys[oracletypes.StoreKey], app.GetSubspace(oracletypes.ModuleName), filepath.Join(homePath, "files"),
		authtypes.FeeCollectorName, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.DistrKeeper,
		app.IBCKeeper.ChannelKeeper, app.IBCKeeper.PortKeeper, scopedOracleKeeper, owasmVM,
	)

	// app.OracleKeeper = oraclekeeper.NewKeeper(
	// 	appCodec,
	// 	keys[oracletypes.StoreKey],
	// 	filepath.Join(homePath, "files"),
	// 	authtypes.FeeCollectorName,
	// 	app.AccountKeeper,
	// 	app.BankKeeper,
	// 	app.StakingKeeper,
	// 	app.DistrKeeper,
	// 	app.AuthzKeeper,
	// 	app.IBCKeeper.ChannelKeeper,
	// 	&app.IBCKeeper.PortKeeper,
	// 	scopedOracleKeeper,
	// 	owasmVM,
	// )

	app.CoinswapKeeper = coinswapkeeper.NewKeeper(
		appCodec,
		keys[coinswaptypes.StoreKey],
		app.GetSubspace(coinswaptypes.ModuleName),
		app.BankKeeper,
		app.DistrKeeper,
		app.OracleKeeper,
	)
	app.AuctionKeeper = auctionkeeper.NewKeeper(
		appCodec,
		keys[auctiontypes.StoreKey],
		app.GetSubspace(auctiontypes.ModuleName),
		app.OracleKeeper,
		app.CoinswapKeeper,
	)

	app.TelemetryKeeper = telemetrykeeper.NewKeeper(appCodec, encodingConfig.TxConfig, app.BankKeeper, app.StakingKeeper, app.DistrKeeper)

	oracleModule := oracle.NewAppModule(app.OracleKeeper)
	oracleModuleIBC := oracle.NewIBCModule(app.OracleKeeper)

	appVersionGetter := MyAppVersionGetter{App: app}

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(transfertypes.ModuleName, transferModuleIBC)
	ibcRouter.AddRoute(oracletypes.ModuleName, oracleModuleIBC)
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)
	// ibcRouter.AddRoute(wasm.ModuleName, wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper))
	ibcRouter.AddRoute(wasm.ModuleName, wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper, appVersionGetter))
	app.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router.
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper
	/****  Module Options ****/

	/****  Module Options ****/
	skipGenesisInvariants := false
	opt := appOpts.Get(crisis.FlagSkipGenesisInvariants)
	if opt, ok := opt.(bool); ok {
		skipGenesisInvariants = opt
	}

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		auth.NewAppModule(
			appCodec,
			app.AccountKeeper,
			nil,
			app.GetSubspace(authtypes.ModuleName),
		),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		odinbank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		odinmint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		oracleModule,
		coinswap.NewAppModule(app.CoinswapKeeper),
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		auction.NewAppModule(app.AuctionKeeper),
		telemetry.NewAppModule(app.TelemetryKeeper),
		transferModule,
		icaModule,
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
	)
	// NOTE: Oracle module must occur before distr as it takes some fee to distribute to active oracle validators.
	// NOTE: During begin block slashing happens after distr.BeginBlocker so that there is nothing left
	// over in the validator fee pool, so as to keep the CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName, capabilitytypes.ModuleName, odinminttypes.ModuleName, oracletypes.ModuleName, distrtypes.ModuleName,
		auctiontypes.ModuleName, slashingtypes.ModuleName, evidencetypes.ModuleName, stakingtypes.ModuleName, ibcexported.ModuleName, icatypes.ModuleName,
		authz.ModuleName, feegrant.ModuleName, group.ModuleName, paramstypes.ModuleName, vestingtypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName,
		govtypes.ModuleName, crisistypes.ModuleName, genutiltypes.ModuleName, transfertypes.ModuleName, telemetrytypes.ModuleName,
		coinswaptypes.ModuleName, consensusparamtypes.ModuleName, wasmtypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName, govtypes.ModuleName, stakingtypes.ModuleName, oracletypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName,
		govtypes.ModuleName, capabilitytypes.ModuleName, telemetrytypes.ModuleName, coinswaptypes.ModuleName, transfertypes.ModuleName,
		paramstypes.ModuleName, vestingtypes.ModuleName, evidencetypes.ModuleName, distrtypes.ModuleName, auctiontypes.ModuleName,
		authz.ModuleName, feegrant.ModuleName, group.ModuleName, slashingtypes.ModuleName, genutiltypes.ModuleName, ibcexported.ModuleName, icatypes.ModuleName, odinminttypes.ModuleName, upgradetypes.ModuleName,
		consensusparamtypes.ModuleName, wasmtypes.ModuleName,
	)
	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName, odinminttypes.ModuleName, oracletypes.ModuleName,
		distrtypes.ModuleName, stakingtypes.ModuleName, slashingtypes.ModuleName, govtypes.ModuleName, crisistypes.ModuleName,
		ibcexported.ModuleName, icatypes.ModuleName, genutiltypes.ModuleName, evidencetypes.ModuleName, coinswaptypes.ModuleName, auctiontypes.ModuleName,
		transfertypes.ModuleName, authz.ModuleName, feegrant.ModuleName, group.ModuleName, paramstypes.ModuleName, upgradetypes.ModuleName, vestingtypes.ModuleName, consensusparamtypes.ModuleName,
		telemetrytypes.ModuleName, wasm.ModuleName, icatypes.ModuleName,
	)
	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions

	overrideModules := map[string]module.AppModuleSimulation{}
	app.sm = module.NewSimulationManagerFromAppModules(app.mm.Modules, overrideModules)

	app.sm.RegisterStoreDecoders()

	// Initialize stores.
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp.
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper:         app.IBCKeeper,
			TxCounterStoreKey: keys[wasm.StoreKey],
			WasmConfig:        wasmConfig,
			Cdc:               appCodec,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create ante handler: %s", err))
	}
	if !disableFeelessReports {
		anteHandler = bandante.NewFeelessReportsAnteHandler(anteHandler, app.OracleKeeper)
	}
	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	// app.UpgradeKeeper.SetUpgradeHandler("v0.6.0", func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	// 	fromVM[icatypes.ModuleName] = icaModule.ConsensusVersion()
	// 	// create ICS27 Controller submodule params
	// 	controllerParams := icacontrollertypes.Params{
	// 		ControllerEnabled: true,
	// 	}
	// 	// create ICS27 Host submodule params
	// 	hostParams := icahosttypes.Params{
	// 		HostEnabled: true,
	// 		AllowMessages: []string{
	// 			"/cosmos.authz.v1beta1.MsgExec",
	// 			"/cosmos.authz.v1beta1.MsgGrant",
	// 			"/cosmos.authz.v1beta1.MsgRevoke",
	// 			"/cosmos.bank.v1beta1.MsgSend",
	// 			"/cosmos.bank.v1beta1.MsgMultiSend",
	// 			"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress",
	// 			"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
	// 			"/cosmos.distribution.v1beta1.MsgFundCommunityPool",
	// 			"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
	// 			"/cosmos.feegrant.v1beta1.MsgGrantAllowance",
	// 			"/cosmos.feegrant.v1beta1.MsgRevokeAllowance",
	// 			"/cosmos.gov.v1beta1.MsgVoteWeighted",
	// 			"/cosmos.gov.v1beta1.MsgSubmitProposal",
	// 			"/cosmos.gov.v1beta1.MsgDeposit",
	// 			"/cosmos.gov.v1beta1.MsgVote",
	// 			"/cosmos.staking.v1beta1.MsgEditValidator",
	// 			"/cosmos.staking.v1beta1.MsgDelegate",
	// 			"/cosmos.staking.v1beta1.MsgUndelegate",
	// 			"/cosmos.staking.v1beta1.MsgBeginRedelegate",
	// 			"/cosmos.staking.v1beta1.MsgCreateValidator",
	// 			"/cosmos.vesting.v1beta1.MsgCreateVestingAccount",
	// 			"/ibc.applications.transfer.v1.MsgTransfer",
	// 			sdk.MsgTypeURL(&wasmtypes.MsgStoreCode{}),
	// 			sdk.MsgTypeURL(&wasmtypes.MsgInstantiateContract{}),
	// 			sdk.MsgTypeURL(&wasmtypes.MsgExecuteContract{}),
	// 		},
	// 	}

	// 	ctx.Logger().Info("start to init interchainaccount module...")
	// 	// initialize ICS27 module
	// 	icaModule.InitModule(ctx, controllerParams, hostParams)
	// 	ctx.Logger().Info("start to run module migrations...")

	// 	bech32ibc.InitGenesis(ctx, *app.Bech32IbcKeeper, bech32ibctypes.GenesisState{
	// 		NativeHRP:     "odin",
	// 		HrpIBCRecords: []bech32ibctypes.HrpIbcRecord{},
	// 	})

	// 	app.mm.OrderMigrations = make([]string, 0)
	// 	// app.mm.OrderMigrations = append(app.mm.OrderMigrations, gravitytypes.ModuleName)
	// 	app.mm.OrderMigrations = append(app.mm.OrderMigrations, wasm.ModuleName)

	// 	return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	// })

	// app.UpgradeKeeper.SetUpgradeHandler(
	// 	"v0.6.2",
	// 	func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	// 		app.mm.OrderMigrations = make([]string, 0)
	// 		app.mm.OrderMigrations = append(app.mm.OrderMigrations, authz.ModuleName)

	// 		consensusParams := app.GetConsensusParams(ctx)
	// 		consensusParams.Block.MaxGas = 100000000
	// 		consensusParams.Block.MaxBytes = 5020096
	// 		app.StoreConsensusParams(ctx, consensusParams)

	// 		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	// 	},
	// )

	// upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	// }

	// if upgradeInfo.Name == "v0.6.0" && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
	// 	storeUpgrades := storetypes.StoreUpgrades{
	// 		Added: []string{icahosttypes.StoreKey, wasm.StoreKey, gravitytypes.StoreKey, bech32ibctypes.StoreKey},
	// 	}

	// 	// configure store loader that checks if version == upgradeHeight and applies store upgrades
	// 	app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	// }

	// if upgradeInfo.Name == "v0.6.2" && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
	// 	storeUpgrades := storetypes.StoreUpgrades{
	// 		Added: []string{authzkeeper.StoreKey},
	// 	}

	// 	// configure store loader that checks if version == upgradeHeight and applies store upgrades
	// 	app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	// }

	if manager := app.SnapshotManager(); manager != nil {
		err = manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic("failed to register snapshot extension: " + err.Error())
		}
	}

	if loadLatest {
		err := app.LoadLatestVersion()
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedOracleKeeper = scopedOracleKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper
	app.ScopedWasmKeeper = scopedWasmKeeper

	return app
}

// RegisterUpgradeHandlers returns upgrade handlers
func (app *OdinApp) RegisterUpgradeHandlers(cfg module.Configurator) {
	// Don't need upgrade handlers as of now
	//app.UpgradeKeeper.SetUpgradeHandler(v7.UpgradeName, v7.CreateUpgradeHandler(*app.mm, cfg, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MintKeeper))
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// Gaia. It is useful for tests and clients who do not want to construct the
// full gaia application
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}

// Name returns the name of the App.
func (app *OdinApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block.
func (app *OdinApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.DeliverContext = ctx
	res := app.mm.BeginBlock(ctx, req)
	for _, hook := range app.hooks {
		hook.AfterBeginBlock(ctx, req, res)
	}
	return res
}

// EndBlocker application updates every end block.
func (app *OdinApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	res := app.mm.EndBlock(ctx, req)
	for _, hook := range app.hooks {
		hook.AfterEndBlock(ctx, req, res)
	}
	return res
}

// Commit overrides the default BaseApp's ABCI commit by adding DeliverContext clearing.
func (app *OdinApp) Commit() (res abci.ResponseCommit) {
	for _, hook := range app.hooks {
		hook.BeforeCommit()
	}
	app.DeliverContext = sdk.Context{}
	return app.BaseApp.Commit()
}

// InitChainer application update at chain initialization
func (app *OdinApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	res := app.mm.InitGenesis(ctx, app.appCodec, genesisState)

	for _, hook := range app.hooks {
		hook.AfterInitChain(ctx, req, res)
	}
	return res
}

// DeliverTx overwrite DeliverTx to apply afterDeliverTx hook
func (app *OdinApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	res := app.BaseApp.DeliverTx(req)
	for _, hook := range app.hooks {
		hook.AfterDeliverTx(app.DeliverContext, req, res)
	}
	return res
}

func (app *OdinApp) Query(req abci.RequestQuery) abci.ResponseQuery {
	hookReq := req

	// when a client did not provide a query height, manually inject the latest
	if hookReq.Height == 0 {
		hookReq.Height = app.LastBlockHeight()
	}

	for _, hook := range app.hooks {
		res, stop := hook.ApplyQuery(hookReq)
		if stop {
			return res
		}
	}
	return app.BaseApp.Query(req)
}

// LoadHeight loads a particular height
func (app *OdinApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *OdinApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *OdinApp) BlockedAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}
	return blacklistedAddrs
}

// LegacyAmino returns OdinApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *OdinApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns Band's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *OdinApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Band's InterfaceRegistry
func (app *OdinApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *OdinApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *OdinApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *OdinApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *OdinApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *OdinApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *OdinApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for additional services
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *OdinApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *OdinApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}

// RegisterNodeService registers all additional services.
func (app *OdinApp) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccPerms returns a mapping of the application's module account permissions.
func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(odinminttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(transfertypes.ModuleName)
	paramsKeeper.Subspace(oracletypes.ModuleName)
	paramsKeeper.Subspace(coinswaptypes.ModuleName)
	paramsKeeper.Subspace(auctiontypes.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)

	return paramsKeeper
}

// AddHook appends hook that will be call after process abci request
func (app *OdinApp) AddHook(hook Hook) {
	app.hooks = append(app.hooks, hook)
}

func GetWasmOpts(appOpts servertypes.AppOptions) []wasm.Option {
	var wasmOpts []wasm.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}

	wasmOpts = append(wasmOpts, wasmkeeper.WithGasRegister(NewWasmGasRegister()))

	return wasmOpts
}
