# Coinswap Module

This document introduces the [Concepts](#concepts), [Parameters](#parameters), [Transactions](#transactions), [Queries](#queries), [Fees](#fees), [Rates Balancing](#rates-balancing), Go code [Exchange Example](#exchange-example) and [Possible Changes](#further-possible-module-changes) of the Coinswap module. The module provides the logic to exchange coins.

## Concepts

Exchanging logic is partially inspired by Osmosis pools (however, there is only one exchanging pool in the Coinswap module [for now](#further-possible-module-changes)).<br>
Pool operates with decimal coins for better exchange rates balancing and collecting fees. Exchange rates performed via a number of "balances". Each balance represents a unit of value, so they are equal to each other.<br>
Coins that the user receives as a result of the exchange are stored in the distribution pool.

## Parameters

Parameters define exchange rates and fees. They can be modified using the Governance module.

### Example

Possible default parameters in `genesis.json`:

```json
{
  "params": {
    "exchange_rates": [
      {
        "denom": "loki",
        "amount": "2.000000000000000000"
      },
      {
        "denom": "minigeo",
        "amount": "1.000000000000000000"
      }
    ],
    "fee": "0.001000000000000000"
  }
}
```

## Transactions

The module has a single `exchange` transaction.

### Usage

```bash
odind tx coinswap exchange [amount] [to-denom] [flags]
```

### Example

```bash
odind tx coinswap exchange 2000000loki minigeo --from <your-wallet-name>
```

## Queries

The module has `params` and `rate` queries.

### Usage

```bash
odind query coinswap params [flags]
```

```bash
odind query coinswap rate [from-denom] [to-denom] [flags]
```

## Fees

Fees are taken from coins for exchange and sent to the distribution pool.

## Rates Balancing

The module parameters have the `ExchangeRates` field, which stores an array of decimal coins.<br>
Rates might look like `2.000000000000000000loki,1.000000000000000000minigeo`, which means `1 loki = 1 / 2 minigeo = 0.5 minigeo` or `1 minigeo = 2 / 1 loki = 2 loki`.<br>
When exchanging coins, the multiplier for the number of coins for exchange is calculated by the formula `m = x / y`, where `x` is received coins rate, `y` is exchanged coins rate.

## Exchange Example

Let's illustrate simplified 1000 minigeo to loki (and backwards) exchanging logic via Go code:

```
func Exchange(params Params, amount sdk.DecCoin, toDenom string) (exchanged sdk.DecCoin) {
	// take fee
	fee := sdk.NewDecCoinFromDec(amount.Denom, amount.Amount.Mul(params.Fee))
	amount = amount.Sub(fee)

	// exchange coins
	multiplier := params.ExchangeRates.AmountOf(toDenom).Quo(params.ExchangeRates.AmountOf(amount.Denom))
	exchangedCoins := sdk.NewDecCoinFromDec(toDenom, amount.Amount.Mul(multiplier))

	return exchangedCoins
}

func TestNewExchange(t *testing.T) {
	// set module params
	var moduleParams Params
	moduleParams.ExchangeRates = sdk.NewDecCoins(sdk.NewInt64DecCoin("minigeo", 1), sdk.NewInt64DecCoin("loki", 2))
	moduleParams.Fee = sdk.MustNewDecFromStr("0.001000000000000000")

	// exchange coins
	exchanged := Exchange(moduleParams, sdk.NewInt64DecCoin("minigeo", 1000), "loki")
	fmt.Println(exchanged)

	// and backwards
	exchanged = Exchange(moduleParams, exchanged, "minigeo")
	fmt.Println(exchanged)
}
```

Output:

```
1998.000000000000000000loki
998.001000000000000000minigeo
```

So, as we can see, module provides simple and governable rates balancing mechanism.

## Further Possible Module Changes

* Adding different fees for different coins for exchange;
* Adding isolated coinswap pools (so some coins cannot be exchanged for others), possibly with different exchange rates for some cases;
  * Moving fee parameter to isolated pool structure.