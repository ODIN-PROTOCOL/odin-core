# Coinswap Module

This document introduces the [Parameters](#parameters), [Transactions](#transactions), [Queries](#queries), [Rates Balancing](#rates-balancing) and [Fees](#fees) of the Coinswap module. The module provides the logic to exchange coins.

## Parameters

Parameters define exchange rates and fees. They can be modified using the Governance module.

### Example

Default parameters in `genesis.json`:

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
odind tx coinswap exchange [from-denom] [to-denom] [amount] [flags]
```

### Example

```bash
odind tx coinswap exchange loki minigeo 2000000loki --from <your-wallet-name>
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

## Rates Balancing

The module parameters have the `ExchangeRates` field, which stores an array of decimal coins.<br>
Rates might look like `2.000000000000000000loki,1.000000000000000000minigeo`, which means `1 loki = 1 / 2 minigeo = 0.5 minigeo` or `1 minigeo = 2 / 1 loki = 2 loki`.<br>
When exchanging coins, the multiplier for the number of coins for exchange is calculated by the formula `m = x / y`, where `x` - received coins rate, `y` - exchanged coins rate.

## Fees

Fees are taken from coins for exchange and sent to the distribution pool. Coins that the user receives as a result of the exchange are also stored in the distribution pool.