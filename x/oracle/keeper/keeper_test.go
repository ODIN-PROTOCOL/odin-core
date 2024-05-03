package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ODIN-PROTOCOL/odin-core/testing/testapp"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

func TestGetSetRequestCount(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Initially request count must be 0.
	requestCount, err := k.GetRequestCount(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(0), requestCount)
	// After we set the count manually, it should be reflected.
	err = k.SetRequestCount(ctx, 42)
	require.NoError(t, err)
	requestCount, err = k.GetRequestCount(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(42), requestCount)
}

func TestGetDataSourceCount(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	err := k.SetDataSourceCount(ctx, 42)
	require.Nil(t, err)
	dataSourceCount, err := k.GetDataSourceCount(ctx)
	require.Nil(t, err)
	require.Equal(t, uint64(42), dataSourceCount)
}

func TestGetSetOracleScriptCount(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	err := k.SetOracleScriptCount(ctx, 42)
	require.Nil(t, err)
	oracleScriptCount, err := k.GetOracleScriptCount(ctx)
	require.Nil(t, err)
	require.Equal(t, uint64(42), oracleScriptCount)
}

func TestGetSetRollingSeed(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	err := k.SetRollingSeed(ctx, []byte("HELLO_WORLD"))
	require.Nil(t, err)
	rollingSeed, err := k.GetRollingSeed(ctx)
	require.Nil(t, err)
	require.Equal(t, []byte("HELLO_WORLD"), rollingSeed)
}

func TestGetNextRequestID(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// First request id must be 1.
	nextRequestID, err := k.GetNextRequestID(ctx)
	require.Nil(t, err)
	require.Equal(t, types.RequestID(1), nextRequestID)

	// After we add new requests, the request count must increase accordingly.
	requestCount, err := k.GetRequestCount(ctx)
	require.Nil(t, err)
	require.Equal(t, uint64(1), requestCount)

	nextRequestID, err = k.GetNextRequestID(ctx)
	require.Nil(t, err)
	require.Equal(t, types.RequestID(2), nextRequestID)

	nextRequestID, err = k.GetNextRequestID(ctx)
	require.Nil(t, err)
	require.Equal(t, types.RequestID(3), nextRequestID)

	nextRequestID, err = k.GetNextRequestID(ctx)
	require.Nil(t, err)
	require.Equal(t, types.RequestID(4), nextRequestID)

	requestCount, err = k.GetRequestCount(ctx)
	require.Nil(t, err)
	require.Equal(t, uint64(4), requestCount)
}

func TestGetNextDataSourceID(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	initialID, err := k.GetDataSourceCount(ctx)
	require.NoError(t, err)
	nextDataSourceID, err := k.GetNextDataSourceID(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DataSourceID(initialID+1), nextDataSourceID)
	nextDataSourceID, err = k.GetNextDataSourceID(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DataSourceID(initialID+2), nextDataSourceID)
	nextDataSourceID, err = k.GetNextDataSourceID(ctx)
	require.NoError(t, err)
	require.Equal(t, types.DataSourceID(initialID+3), nextDataSourceID)
}

func TestGetNextOracleScriptID(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	initialID, err := k.GetOracleScriptCount(ctx)
	require.NoError(t, err)
	nextOracleScriptID, err := k.GetNextOracleScriptID(ctx)
	require.NoError(t, err)
	require.Equal(t, types.OracleScriptID(initialID+1), nextOracleScriptID)
	nextOracleScriptID, err = k.GetNextOracleScriptID(ctx)
	require.NoError(t, err)
	require.Equal(t, types.OracleScriptID(initialID+2), nextOracleScriptID)
	nextOracleScriptID, err = k.GetNextOracleScriptID(ctx)
	require.NoError(t, err)
	require.Equal(t, types.OracleScriptID(initialID+3), nextOracleScriptID)
}

func TestGetSetRequestLastExpiredID(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Initially last expired request must be 0.
	requestLastExpired, err := k.GetRequestLastExpired(ctx)
	require.NoError(t, err)
	require.Equal(t, types.RequestID(0), requestLastExpired)
	err = k.SetRequestLastExpired(ctx, 20)
	require.NoError(t, err)
	requestLastExpired, err = k.GetRequestLastExpired(ctx)
	require.NoError(t, err)
	require.Equal(t, types.RequestID(20), requestLastExpired)
}
