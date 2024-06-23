package odin

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"
	"cosmossdk.io/x/circuit"
	circuittypes "cosmossdk.io/x/circuit/types"
	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/nft"
	nftmodule "cosmossdk.io/x/nft/module"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	icagenesistypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/genesis/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransafertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/ODIN-PROTOCOL/odin-core/app/upgrades/v2_6"
	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// GenesisState defines a type alias for the Odin genesis application state.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	cdc := MakeEncodingConfig().Marshaler
	denom := "loki"
	// Get default genesis states of the modules we are to override.
	authGenesis := authtypes.DefaultGenesisState()
	stakingGenesis := stakingtypes.DefaultGenesisState()
	distrGenesis := distrtypes.DefaultGenesisState()
	mintGenesis := minttypes.DefaultGenesisState()
	govGenesis := govv1.DefaultGenesisState()
	crisisGenesis := crisistypes.DefaultGenesisState()
	slashingGenesis := slashingtypes.DefaultGenesisState()
	oracleGenesis := oracletypes.DefaultGenesisState()
	icaGenesis := icagenesistypes.DefaultGenesis()
	// Override the genesis parameters.
	authGenesis.Params.TxSizeCostPerByte = 5
	stakingGenesis.Params.BondDenom = denom
	stakingGenesis.Params.HistoricalEntries = 1000
	mintGenesis.Params.BlocksPerYear = 10519200 // target 3-second block time
	mintGenesis.Params.MintDenom = denom
	govGenesis.Params.MinDeposit = sdk.NewCoins(
		sdk.NewCoin(denom, sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction)),
	)
	crisisGenesis.ConstantFee = sdk.NewCoin(denom, sdk.TokensFromConsensusPower(10000, sdk.DefaultPowerReduction))
	slashingGenesis.Params.SignedBlocksWindow = 30000                                // approximately 1 day
	slashingGenesis.Params.MinSignedPerWindow = math.LegacyNewDecWithPrec(5, 2)      // 5%
	slashingGenesis.Params.DowntimeJailDuration = 60 * 10 * time.Second              // 10 minutes
	slashingGenesis.Params.SlashFractionDoubleSign = math.LegacyNewDecWithPrec(5, 2) // 5%
	slashingGenesis.Params.SlashFractionDowntime = math.LegacyNewDecWithPrec(1, 4)   // 0.01%

	icaGenesis.HostGenesisState.Params = icahosttypes.Params{
		HostEnabled:   true,
		AllowMessages: v2_6.ICAAllowMessages,
	}
	//oracleGenesis.ModuleCoinsAccount = "odin1lqf6hm3nfunmhppmjhgrme9jp9d8vle90hjy5m"

	return GenesisState{
		authtypes.ModuleName:         cdc.MustMarshalJSON(authGenesis),
		genutiltypes.ModuleName:      genutil.AppModuleBasic{}.DefaultGenesis(cdc),
		banktypes.ModuleName:         bank.AppModuleBasic{}.DefaultGenesis(cdc),
		capabilitytypes.ModuleName:   capability.AppModuleBasic{}.DefaultGenesis(cdc),
		stakingtypes.ModuleName:      cdc.MustMarshalJSON(stakingGenesis),
		minttypes.ModuleName:         cdc.MustMarshalJSON(mintGenesis),
		distrtypes.ModuleName:        cdc.MustMarshalJSON(distrGenesis),
		govtypes.ModuleName:          cdc.MustMarshalJSON(govGenesis),
		crisistypes.ModuleName:       cdc.MustMarshalJSON(crisisGenesis),
		slashingtypes.ModuleName:     cdc.MustMarshalJSON(slashingGenesis),
		ibcexported.ModuleName:       ibc.AppModuleBasic{}.DefaultGenesis(cdc),
		upgradetypes.ModuleName:      upgrade.AppModuleBasic{}.DefaultGenesis(cdc),
		evidencetypes.ModuleName:     evidence.AppModuleBasic{}.DefaultGenesis(cdc),
		authz.ModuleName:             authzmodule.AppModuleBasic{}.DefaultGenesis(cdc),
		feegrant.ModuleName:          feegrantmodule.AppModuleBasic{}.DefaultGenesis(cdc),
		group.ModuleName:             groupmodule.AppModuleBasic{}.DefaultGenesis(cdc),
		ibctransafertypes.ModuleName: ibctransfer.AppModuleBasic{}.DefaultGenesis(cdc),
		icatypes.ModuleName:          cdc.MustMarshalJSON(icaGenesis),
		oracletypes.ModuleName:       cdc.MustMarshalJSON(oracleGenesis),
		nft.ModuleName:               nftmodule.AppModuleBasic{}.DefaultGenesis(cdc),
		circuittypes.ModuleName:      circuit.AppModuleBasic{}.DefaultGenesis(cdc),
		//wasmtypes.ModuleName:         wasm.AppModuleBasic{}.DefaultGenesis(cdc),
	}
}
