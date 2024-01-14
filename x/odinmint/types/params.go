package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMintDenom            = []byte("MintDenom")
	KeyInflationRateChange  = []byte("InflationRateChange")
	KeyInflationMax         = []byte("InflationMax")
	KeyInflationMin         = []byte("InflationMin")
	KeyGoalBonded           = []byte("GoalBonded")
	KeyBlocksPerYear        = []byte("BlocksPerYear")
	KeyMintAir              = []byte("MintAir")
	KeyIntegrationAddresses = []byte("IntegrationAddresses")
	KeyMaxWithdrawalPerTime = []byte("MaxWithdrawalPerTime")
	KeyEligibleAccountsPool = []byte("EligibleAccountsPool")
	KeyMaxAllowedMintVolume = []byte("MaxAllowedMintVolume")
	KeyAllowedMintDenoms    = []byte("AllowedMintDenoms")
	KeyAllowedMinter        = []byte("AllowedMinter")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string,
	inflationRateChange, inflationMax, inflationMin, goalBonded math.LegacyDec,
	MaxWithdrawalPerTime sdk.Coins,
	blocksPerYear uint64,
	mintAir bool,
	integrationAddresses map[string]string,
	eligibleAccountsPool []string,
	maxAllowedMintVolume sdk.Coins,
	allowedMintDenoms []*AllowedDenom,
	AllowedMinter []string,
) Params {
	return Params{
		MintDenom:            mintDenom,
		InflationRateChange:  inflationRateChange,
		InflationMax:         inflationMax,
		InflationMin:         inflationMin,
		GoalBonded:           goalBonded,
		BlocksPerYear:        blocksPerYear,
		MintAir:              mintAir,
		IntegrationAddresses: integrationAddresses,
		MaxWithdrawalPerTime: MaxWithdrawalPerTime,
		EligibleAccountsPool: eligibleAccountsPool,
		MaxAllowedMintVolume: maxAllowedMintVolume,
		AllowedMintDenoms:    allowedMintDenoms,
		AllowedMinter:        AllowedMinter,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:            sdk.DefaultBondDenom,
		InflationRateChange:  math.LegacyNewDecWithPrec(13, 2),
		InflationMax:         math.LegacyNewDecWithPrec(20, 2),
		InflationMin:         math.LegacyNewDecWithPrec(7, 2),
		GoalBonded:           math.LegacyNewDecWithPrec(67, 2),
		BlocksPerYear:        uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
		MintAir:              false,
		IntegrationAddresses: map[string]string{}, // default value (might be invalid for actual use)
		MaxWithdrawalPerTime: sdk.Coins{sdk.NewCoin("loki", math.NewInt(100))},
		EligibleAccountsPool: []string{"odin1cgfdwtrqfdrzh4z8rkcyx8g4jv22v8wgs39amj"},
		MaxAllowedMintVolume: sdk.Coins{sdk.NewCoin("minigeo", math.NewInt(100000000))},
		AllowedMintDenoms:    []*AllowedDenom{{"loki", "odin"}, {"minigeo", "geo"}},
		AllowedMinter:        []string{"odin1cgfdwtrqfdrzh4z8rkcyx8g4jv22v8wgs39amj"},
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationRateChange(p.InflationRateChange); err != nil {
		return err
	}
	if err := validateInflationMax(p.InflationMax); err != nil {
		return err
	}
	if err := validateInflationMin(p.InflationMin); err != nil {
		return err
	}
	if err := validateGoalBonded(p.GoalBonded); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}
	if err := validateMintAir(p.MintAir); err != nil {
		return err
	}
	if err := validateIntegrationAddresses(p.IntegrationAddresses); err != nil {
		return err
	}
	if err := validateMaxWithdrawalPerTime(p.MaxWithdrawalPerTime); err != nil {
		return err
	}
	if err := validateEligibleAccountsPool(p.EligibleAccountsPool); err != nil {
		return err
	}
	if err := validateMaxAllowedMintVolume(p.MaxAllowedMintVolume); err != nil {
		return err
	}
	if err := validateAllowedMintDenoms(p.AllowedMintDenoms); err != nil {
		return err
	}
	if err := validateAllowedMinter(p.AllowedMinter); err != nil {
		return err
	}
	if p.InflationMax.LT(p.InflationMin) {
		return fmt.Errorf(
			"max inflation (%s) must be greater than or equal to min inflation (%s)",
			p.InflationMax, p.InflationMin,
		)
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	return fmt.Sprintf(`Minting Params:
  Mint Denom:             	%s
  Inflation Rate Change:  	%s
  Inflation Max:          	%s
  Inflation Min:          	%s
  Goal Bonded:            	%s
  Blocks Per Year:        	%d
  Integration Addresses: 	%s
  Max Withdrawal Per Time:	%s
  Eligible Accounts Pool: 	%s
`,
		p.MintDenom, p.InflationRateChange, p.InflationMax, p.InflationMin, p.GoalBonded,
		p.BlocksPerYear, p.IntegrationAddresses, p.MaxWithdrawalPerTime, p.EligibleAccountsPool,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyInflationRateChange, &p.InflationRateChange, validateInflationRateChange),
		paramtypes.NewParamSetPair(KeyInflationMax, &p.InflationMax, validateInflationMax),
		paramtypes.NewParamSetPair(KeyInflationMin, &p.InflationMin, validateInflationMin),
		paramtypes.NewParamSetPair(KeyGoalBonded, &p.GoalBonded, validateGoalBonded),
		paramtypes.NewParamSetPair(KeyBlocksPerYear, &p.BlocksPerYear, validateBlocksPerYear),
		paramtypes.NewParamSetPair(KeyMintAir, &p.MintAir, validateMintAir),
		paramtypes.NewParamSetPair(KeyIntegrationAddresses, &p.IntegrationAddresses, validateIntegrationAddresses),
		paramtypes.NewParamSetPair(KeyMaxWithdrawalPerTime, &p.MaxWithdrawalPerTime, validateMaxWithdrawalPerTime),
		paramtypes.NewParamSetPair(KeyEligibleAccountsPool, &p.EligibleAccountsPool, validateEligibleAccountsPool),
		paramtypes.NewParamSetPair(KeyMaxAllowedMintVolume, &p.MaxAllowedMintVolume, validateMaxAllowedMintVolume),
		paramtypes.NewParamSetPair(KeyAllowedMintDenoms, &p.AllowedMintDenoms, validateAllowedMintDenoms),
		paramtypes.NewParamSetPair(KeyAllowedMinter, &p.AllowedMinter, validateAllowedMinter),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.Wrap(ErrInvalidMintDenom, "mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateInflationRateChange(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation rate change cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("inflation rate change too large: %s", v)
	}

	return nil
}

func validateEligibleAccountsPool(i interface{}) error {
	_, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateInflationMax(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("max inflation cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("max inflation too large: %s", v)
	}

	return nil
}

func validateInflationMin(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("min inflation cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("min inflation too large: %s", v)
	}

	return nil
}

func validateGoalBonded(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("goal bonded cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("goal bonded too large: %s", v)
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}

func validateMaxWithdrawalPerTime(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !v.IsValid() {
		return fmt.Errorf("max withdrawal per time parameter is not valid: %s", v)
	}
	if v.IsAnyNegative() {
		return fmt.Errorf("max withdrawal per time cannot be negative: %s", v)
	}

	return nil
}

func validateIntegrationAddresses(i interface{}) error {
	_, ok := i.(map[string]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateMintAir(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMaxAllowedMintVolume(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !v.IsValid() {
		return fmt.Errorf("max allowed mint volume parameter is not valid: %s", v)
	}
	if v.IsAnyNegative() {
		return fmt.Errorf("max allowed mint volume cannot be negative: %s", v)
	}

	return nil
}

func validateAllowedMintDenoms(i interface{}) error {
	_, ok := i.([]*AllowedDenom)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateAllowedMinter(i interface{}) error {
	_, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
