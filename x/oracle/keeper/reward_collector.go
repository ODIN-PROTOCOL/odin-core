package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	oracletypes "github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

type rewardCollector struct {
	oracleKeeper Keeper
	bankKeeper   oracletypes.BankKeeper
	collected    sdk.Coins
}

func (r rewardCollector) Collect(ctx sdk.Context, coins sdk.Coins, address sdk.AccAddress) error {
	r.collected = r.collected.Add(coins...)
	return r.oracleKeeper.SetDataProviderAccumulatedReward(ctx, address, coins)
}

func (r rewardCollector) Collected() sdk.Coins {
	return r.collected
}

func (r rewardCollector) CalculateReward(data []byte, pricePerByte sdk.Coins) sdk.Coins {
	price := sdk.NewDecCoinsFromCoins(pricePerByte...)
	reward, _ := price.MulDec(math.LegacyNewDecFromInt(math.NewInt(int64(len(data))))).TruncateDecimal()
	return reward
}

func newRewardCollector(oracleKeeper Keeper, bankKeeper oracletypes.BankKeeper) RewardCollector {
	return &rewardCollector{
		oracleKeeper: oracleKeeper,
		bankKeeper:   bankKeeper,
		collected:    sdk.NewCoins(),
	}
}
