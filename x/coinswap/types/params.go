package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultFromExchange = "geo"
	DefaultToExchange   = "odin"
)

var KeyExchanges = []byte("Exchanges")

// ParamKeyTable param table for coinswap module.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyExchanges, &p.Exchanges, validateExchanges),
	}
}

func DefaultParams() Params {
	return Params{
		Exchanges: []Exchange{
			{
				From:           DefaultFromExchange,
				To:             DefaultToExchange,
				RateMultiplier: sdk.NewDec(1),
			},
		},
	}
}

func validateExchanges(i interface{}) error {
	exchanges, ok := i.([]Exchange)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, ex := range exchanges {
		if !ex.RateMultiplier.IsPositive() && !ex.RateMultiplier.IsZero() {
			return fmt.Errorf("rate multiplier %s must be positive or zero", ex)
		}

		if ex.From == "" || ex.To == "" {
			return fmt.Errorf("one or both denoms are empty. From: %s, To: %s", ex.From, ex.To)
		}
	}

	return nil
}
