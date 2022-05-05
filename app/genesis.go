package odin

import (
	"encoding/json"
	"github.com/ODIN-PROTOCOL/odin-core/x/auction"
	auctiontypes "github.com/ODIN-PROTOCOL/odin-core/x/auction/types"
	"github.com/ODIN-PROTOCOL/odin-core/x/coinswap"
	coinswaptypes "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
	"time"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v2/modules/core"
	ibchost "github.com/cosmos/ibc-go/v2/modules/core/24-host"

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
	govGenesis := govtypes.DefaultGenesisState()
	crisisGenesis := crisistypes.DefaultGenesisState()
	slashingGenesis := slashingtypes.DefaultGenesisState()
	oracleGenesis := oracletypes.DefaultGenesisState()
	transferGenesis := ibctypes.DefaultGenesisState()
	// Override the genesis parameters.
	authGenesis.Params.TxSizeCostPerByte = 5
	stakingGenesis.Params.BondDenom = denom
	stakingGenesis.Params.HistoricalEntries = 1000
	distrGenesis.Params.BaseProposerReward = sdk.NewDecWithPrec(3, 2)   // 3%
	distrGenesis.Params.BonusProposerReward = sdk.NewDecWithPrec(12, 2) // 12%
	mintGenesis.Params.BlocksPerYear = 10519200                         // target 3-second block time
	mintGenesis.Params.MintDenom = denom
	mintGenesis.Params.MintAir = true
	govGenesis.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(denom, sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction)))
	crisisGenesis.ConstantFee = sdk.NewCoin(denom, sdk.TokensFromConsensusPower(10000, sdk.DefaultPowerReduction))
	slashingGenesis.Params.SignedBlocksWindow = 30000                         // approximately 1 day
	slashingGenesis.Params.MinSignedPerWindow = sdk.NewDecWithPrec(5, 2)      // 5%
	slashingGenesis.Params.DowntimeJailDuration = 60 * 10 * time.Second       // 10 minutes
	slashingGenesis.Params.SlashFractionDoubleSign = sdk.NewDecWithPrec(5, 2) // 5%
	slashingGenesis.Params.SlashFractionDowntime = sdk.NewDecWithPrec(1, 4)   // 0.01%
	mintGenesis.ModuleCoinsAccount = "odin13jp4udqlxknzrpsk9jkr3hpmp6gy242xm0s2kq"
	oracleGenesis.ModuleCoinsAccount = "odin1lqf6hm3nfunmhppmjhgrme9jp9d8vle90hjy5m"
	return GenesisState{
		authtypes.ModuleName:       cdc.MustMarshalJSON(authGenesis),
		genutiltypes.ModuleName:    genutil.AppModuleBasic{}.DefaultGenesis(cdc),
		banktypes.ModuleName:       bank.AppModuleBasic{}.DefaultGenesis(cdc),
		capabilitytypes.ModuleName: capability.AppModuleBasic{}.DefaultGenesis(cdc),
		stakingtypes.ModuleName:    cdc.MustMarshalJSON(stakingGenesis),
		minttypes.ModuleName:       cdc.MustMarshalJSON(mintGenesis),
		distrtypes.ModuleName:      cdc.MustMarshalJSON(distrGenesis),
		govtypes.ModuleName:        cdc.MustMarshalJSON(govGenesis),
		crisistypes.ModuleName:     cdc.MustMarshalJSON(crisisGenesis),
		slashingtypes.ModuleName:   cdc.MustMarshalJSON(slashingGenesis),
		ibchost.ModuleName:         ibc.AppModuleBasic{}.DefaultGenesis(cdc),
		upgradetypes.ModuleName:    upgrade.AppModuleBasic{}.DefaultGenesis(cdc),
		evidencetypes.ModuleName:   evidence.AppModuleBasic{}.DefaultGenesis(cdc),
		oracletypes.ModuleName:     cdc.MustMarshalJSON(oracleGenesis),
		coinswaptypes.ModuleName:   coinswap.AppModuleBasic{}.DefaultGenesis(cdc),
		auctiontypes.ModuleName:    auction.AppModuleBasic{}.DefaultGenesis(cdc),
		ibctypes.ModuleName:        cdc.MustMarshalJSON(transferGenesis),
	}
}
