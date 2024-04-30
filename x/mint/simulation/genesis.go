package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	minttypes "github.com/ODIN-PROTOCOL/odin-core/x/mint/types"
)

// Simulation parameter constants
const (
	Inflation           = "inflation"
	InflationRateChange = "inflation_rate_change"
	InflationMax        = "inflation_max"
	InflationMin        = "inflation_min"
	GoalBonded          = "goal_bonded"
)

// GenInflation randomized Inflation
func GenInflation(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationRateChange randomized InflationRateChange
func GenInflationRateChange(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationMax randomized InflationMax
func GenInflationMax(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(20, 2)
}

// GenInflationMin randomized InflationMin
func GenInflationMin(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(7, 2)
}

// GenGoalBonded randomized GoalBonded
func GenGoalBonded(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(67, 2)
}

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation math.LegacyDec
	simState.AppParams.GetOrGenerate(
		Inflation, &inflation, simState.Rand,
		func(r *rand.Rand) { inflation = GenInflation(r) },
	)

	// params
	var inflationRateChange math.LegacyDec
	simState.AppParams.GetOrGenerate(
		InflationRateChange, &inflationRateChange, simState.Rand,
		func(r *rand.Rand) { inflationRateChange = GenInflationRateChange(r) },
	)

	var inflationMax math.LegacyDec
	simState.AppParams.GetOrGenerate(
		InflationMax, &inflationMax, simState.Rand,
		func(r *rand.Rand) { inflationMax = GenInflationMax(r) },
	)

	var inflationMin math.LegacyDec
	simState.AppParams.GetOrGenerate(
		InflationMin, &inflationMin, simState.Rand,
		func(r *rand.Rand) { inflationMin = GenInflationMin(r) },
	)

	var goalBonded math.LegacyDec
	simState.AppParams.GetOrGenerate(
		GoalBonded, &goalBonded, simState.Rand,
		func(r *rand.Rand) { goalBonded = GenGoalBonded(r) },
	)

	mintDenom := simState.BondDenom
	blocksPerYear := uint64(60 * 60 * 8766 / 5)
	maxWithdrawalPerTime := sdk.Coins{sdk.NewCoin("loki", math.NewInt(10))}
	mintAir := false

	eligibleAccountsPool := []string{"odin1cgfdwtrqfdrzh4z8rkcyx8g4jv22v8wgs39amj"}
	maxAllowedMintVolume := sdk.Coins{sdk.NewCoin("minigeo", math.NewInt(100000000))}
	allowedMintDenoms := []*minttypes.AllowedDenom{{TokenUnitDenom: "loki", TokenDenom: "odin"}, {TokenUnitDenom: "minigeo", TokenDenom: "geo"}}
	allowedMinter := []string{"odin1cgfdwtrqfdrzh4z8rkcyx8g4jv22v8wgs39amj"}

	params := minttypes.NewParams(
		mintDenom,
		inflationRateChange,
		inflationMax,
		inflationMin,
		goalBonded,
		maxWithdrawalPerTime,
		blocksPerYear,
		mintAir,
		eligibleAccountsPool,
		maxAllowedMintVolume,
		allowedMintDenoms,
		allowedMinter,
	)
	mintGenesis := minttypes.NewGenesisState(minttypes.InitialMinter(inflation), params, minttypes.InitialMintPool())

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[minttypes.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
