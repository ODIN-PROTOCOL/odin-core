package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"sort"
)

type Keeper struct {
	cdc        codec.BinaryMarshaler
	bankKeeper bankkeeper.ViewKeeper
}

func NewKeeper(cdc codec.BinaryMarshaler, bk bankkeeper.ViewKeeper) Keeper {
	return Keeper{
		cdc:        cdc,
		bankKeeper: bk,
	}
}

func (k Keeper) GetPaginatedBalances(ctx sdk.Context, denom string, limit, offset uint64) []banktypes.Balance {
	balances := make([]banktypes.Balance, 0)

	var currOffset uint64 = 0
	var currLimit uint64 = 0
	k.bankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, balance sdk.Coin) bool {
		if currOffset < offset {
			currOffset++
			return false
		}
		currLimit++
		if currLimit > limit {
			return true
		}
		if balance.GetDenom() != denom {
			return false
		}

		accountBalance := banktypes.Balance{
			Address: addr.String(),
			Coins:   sdk.NewCoins(balance),
		}
		balances = append(balances, accountBalance)
		return false
	})

	sort.Slice(balances, func(i, j int) bool {
		return balances[i].GetCoins().AmountOf(denom).LT(balances[j].GetCoins().AmountOf(denom))
	})

	return balances
}
