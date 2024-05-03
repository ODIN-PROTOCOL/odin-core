package keeper_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/ODIN-PROTOCOL/odin-core/testing/testapp"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

func TestGetSetParams(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	expectedParams := types.Params{
		MaxRawRequestCount:       1,
		MaxAskCount:              10,
		MaxCalldataSize:          256,
		MaxReportDataSize:        512,
		ExpirationBlockCount:     30,
		BaseOwasmGas:             50000,
		PerValidatorRequestGas:   3000,
		SamplingTryCount:         3,
		OracleRewardPercentage:   50,
		InactivePenaltyDuration:  1000,
		IBCRequestEnabled:        true,
		RewardDecreasingFraction: math.LegacyNewDec(10),
	}
	err := k.SetParams(ctx, expectedParams)
	require.NoError(t, err)

	params, err := k.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedParams, params)

	expectedParams = types.Params{
		MaxRawRequestCount:       2,
		MaxAskCount:              20,
		MaxCalldataSize:          512,
		MaxReportDataSize:        256,
		ExpirationBlockCount:     40,
		BaseOwasmGas:             150000,
		PerValidatorRequestGas:   30000,
		SamplingTryCount:         5,
		OracleRewardPercentage:   80,
		InactivePenaltyDuration:  10000,
		IBCRequestEnabled:        false,
		RewardDecreasingFraction: math.LegacyNewDec(10),
	}
	err = k.SetParams(ctx, expectedParams)
	require.NoError(t, err)

	params, err = k.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedParams, params)

	expectedParams = types.Params{
		MaxRawRequestCount:       2,
		MaxAskCount:              20,
		MaxCalldataSize:          512,
		MaxReportDataSize:        256,
		ExpirationBlockCount:     40,
		BaseOwasmGas:             0,
		PerValidatorRequestGas:   0,
		SamplingTryCount:         5,
		OracleRewardPercentage:   0,
		InactivePenaltyDuration:  0,
		IBCRequestEnabled:        false,
		RewardDecreasingFraction: math.LegacyNewDec(100),
	}
	err = k.SetParams(ctx, expectedParams)
	require.NoError(t, err)

	params, err = k.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedParams, params)

	expectedParams = types.Params{
		MaxRawRequestCount:       0,
		MaxAskCount:              20,
		MaxCalldataSize:          512,
		MaxReportDataSize:        256,
		ExpirationBlockCount:     40,
		BaseOwasmGas:             150000,
		PerValidatorRequestGas:   30000,
		SamplingTryCount:         5,
		OracleRewardPercentage:   80,
		InactivePenaltyDuration:  10000,
		IBCRequestEnabled:        false,
		RewardDecreasingFraction: math.LegacyNewDec(1),
	}
	err = k.SetParams(ctx, expectedParams)
	require.EqualError(t, fmt.Errorf("max raw request count must be positive: 0"), err.Error())
}
