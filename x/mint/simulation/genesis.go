package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

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
func GenInflation(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationRateChange randomized InflationRateChange
func GenInflationRateChange(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationMax randomized InflationMax
func GenInflationMax(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(20, 2)
}

// GenInflationMin randomized InflationMin
func GenInflationMin(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(7, 2)
}

// GenGoalBonded randomized GoalBonded
func GenGoalBonded(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(67, 2)
}

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, Inflation, &inflation, simState.Rand,
		func(r *rand.Rand) { inflation = GenInflation(r) },
	)

	// params
	var inflationRateChange sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InflationRateChange, &inflationRateChange, simState.Rand,
		func(r *rand.Rand) { inflationRateChange = GenInflationRateChange(r) },
	)

	var inflationMax sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InflationMax, &inflationMax, simState.Rand,
		func(r *rand.Rand) { inflationMax = GenInflationMax(r) },
	)

	var inflationMin sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InflationMin, &inflationMin, simState.Rand,
		func(r *rand.Rand) { inflationMin = GenInflationMin(r) },
	)

	var goalBonded sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, GoalBonded, &goalBonded, simState.Rand,
		func(r *rand.Rand) { goalBonded = GenGoalBonded(r) },
	)

	mintDenom := sdk.DefaultBondDenom
	blocksPerYear := uint64(60 * 60 * 8766 / 5)
	MaxWithdrawalPerTime := sdk.Coins{}
	mintAir := false

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	integrationAddresses := map[string]string{"eth": hexutil.Encode(privateKeyBytes)}
	eligibleAccountsPool := []string{"odin1cgfdwtrqfdrzh4z8rkcyx8g4jv22v8wgs39amj"}
	maxAllowedMintVolume := sdk.Coins{sdk.NewCoin("minigeo", sdk.NewInt(100000000))}
	allowedMintDenoms := []*minttypes.AllowedDenom{{TokenUnitDenom: "loki", TokenDenom: "odin"}, {TokenUnitDenom: "minigeo", TokenDenom: "geo"}}
	allowedMinter := []string{"odin1cgfdwtrqfdrzh4z8rkcyx8g4jv22v8wgs39amj"}

	params := minttypes.NewParams(
		mintDenom,
		inflationRateChange,
		inflationMax,
		inflationMin,
		goalBonded,
		MaxWithdrawalPerTime,
		blocksPerYear,
		mintAir,
		integrationAddresses,
		eligibleAccountsPool,
		maxAllowedMintVolume,
		allowedMintDenoms,
		allowedMinter,
	)
	mintGenesis := minttypes.NewGenesisState(minttypes.InitialMinter(inflation), params, minttypes.InitialMintPool(), sdk.AccAddress("odin13jp4udqlxknzrpsk9jkr3hpmp6gy242xm0s2kq"))

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[minttypes.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
