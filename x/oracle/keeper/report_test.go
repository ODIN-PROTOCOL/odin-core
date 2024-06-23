package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ODIN-PROTOCOL/odin-core/testing/testapp"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

func defaultRequest() types.Request {
	return types.NewRequest(
		1, BasicCalldata,
		[]sdk.ValAddress{testapp.Validators[0].ValAddress, testapp.Validators[1].ValAddress},
		2, 0, testapp.ParseTime(0),
		BasicClientID, []types.RawRequest{
			types.NewRawRequest(42, 1, BasicCalldata),
			types.NewRawRequest(43, 2, BasicCalldata),
		}, nil, 0,
	)
}

func TestHasReport(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// We should not have a report to request ID 42 from Alice without setting it.
	hasReport, err := k.HasReport(ctx, 42, testapp.Alice.ValAddress)
	require.NoError(t, err)
	require.False(t, hasReport)
	// After we set it, we should be able to find it.
	err = k.SetReport(ctx, 42, types.NewReport(testapp.Alice.ValAddress, true, nil))
	require.NoError(t, err)

	hasReport, err = k.HasReport(ctx, 42, testapp.Alice.ValAddress)
	require.NoError(t, err)
	require.True(t, hasReport)
}

func TestAddReportSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	err := k.SetRequest(ctx, 1, defaultRequest())
	require.NoError(t, err)
	err = k.AddReport(ctx, 1,
		testapp.Validators[0].ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		},
	)
	require.NoError(t, err)

	reports, err := k.GetReports(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, []types.Report{
		types.NewReport(testapp.Validators[0].ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		}),
	}, reports)
}

func TestReportOnNonExistingRequest(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	err := k.AddReport(ctx, 1,
		testapp.Validators[0].ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		},
	)
	require.ErrorIs(t, err, types.ErrRequestNotFound)
}

func TestReportByNotRequestedValidator(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	err := k.SetRequest(ctx, 1, defaultRequest())
	require.NoError(t, err)

	err = k.AddReport(ctx, 1,
		testapp.Alice.ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		},
	)
	require.ErrorIs(t, err, types.ErrValidatorNotRequested)
}

func TestDuplicateReport(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	err := k.SetRequest(ctx, 1, defaultRequest())
	require.NoError(t, err)

	err = k.AddReport(ctx, 1,
		testapp.Validators[0].ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		},
	)
	require.NoError(t, err)
	err = k.AddReport(ctx, 1,
		testapp.Validators[0].ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		},
	)
	require.ErrorIs(t, err, types.ErrValidatorAlreadyReported)
}

func TestReportInvalidDataSourceCount(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	err := k.SetRequest(ctx, 1, defaultRequest())
	require.NoError(t, err)

	err = k.AddReport(ctx, 1,
		testapp.Validators[0].ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
		},
	)
	require.ErrorIs(t, err, types.ErrInvalidReportSize)
}

func TestReportInvalidExternalIDs(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	err := k.SetRequest(ctx, 1, defaultRequest())
	require.NoError(t, err)

	err = k.AddReport(ctx, 1,
		testapp.Validators[0].ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(44, 1, []byte("data2/1")), // BAD EXTERNAL ID!
		},
	)
	require.ErrorIs(t, err, types.ErrRawRequestNotFound)
}

func TestGetReportCount(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// We start by setting some aribrary reports.
	err := k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Alice.ValAddress, true, []types.RawReport{}))
	require.NoError(t, err)
	err = k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Bob.ValAddress, true, []types.RawReport{}))
	require.NoError(t, err)
	err = k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Alice.ValAddress, true, []types.RawReport{}))
	require.NoError(t, err)
	err = k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Bob.ValAddress, true, []types.RawReport{}))
	require.NoError(t, err)
	err = k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Carol.ValAddress, true, []types.RawReport{}))
	require.NoError(t, err)
	// GetReportCount should return the correct values.
	reportCount, err := k.GetReportCount(ctx, types.RequestID(1))
	require.NoError(t, err)
	require.Equal(t, uint64(2), reportCount)

	reportCount, err = k.GetReportCount(ctx, types.RequestID(2))
	require.NoError(t, err)
	require.Equal(t, uint64(3), reportCount)
}

func TestDeleteReports(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// We start by setting some arbitrary reports.
	err := k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Alice.ValAddress, true, []types.RawReport{}))
	require.NoError(t, err)
	err = k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Bob.ValAddress, true, []types.RawReport{}))
	require.NoError(t, err)
	err = k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Alice.ValAddress, true, []types.RawReport{}))
	require.NoError(t, err)
	err = k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Bob.ValAddress, true, []types.RawReport{}))
	require.NoError(t, err)
	err = k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Carol.ValAddress, true, []types.RawReport{}))
	require.NoError(t, err)
	// All reports should exist on the state.
	hasReport, err := k.HasReport(ctx, types.RequestID(1), testapp.Alice.ValAddress)
	require.NoError(t, err)
	require.True(t, hasReport)
	hasReport, err = k.HasReport(ctx, types.RequestID(1), testapp.Bob.ValAddress)
	require.NoError(t, err)
	require.True(t, hasReport)
	hasReport, err = k.HasReport(ctx, types.RequestID(2), testapp.Alice.ValAddress)
	require.NoError(t, err)
	require.True(t, hasReport)
	hasReport, err = k.HasReport(ctx, types.RequestID(2), testapp.Bob.ValAddress)
	require.NoError(t, err)
	require.True(t, hasReport)
	hasReport, err = k.HasReport(ctx, types.RequestID(2), testapp.Carol.ValAddress)
	require.NoError(t, err)
	require.True(t, hasReport)
	// After we delete reports related to request#1, they must disappear.
	err = k.DeleteReports(ctx, types.RequestID(1))
	require.NoError(t, err)

	hasReport, err = k.HasReport(ctx, types.RequestID(1), testapp.Alice.ValAddress)
	require.NoError(t, err)
	require.False(t, hasReport)
	hasReport, err = k.HasReport(ctx, types.RequestID(1), testapp.Bob.ValAddress)
	require.NoError(t, err)
	require.False(t, hasReport)
	hasReport, err = k.HasReport(ctx, types.RequestID(2), testapp.Alice.ValAddress)
	require.NoError(t, err)
	require.True(t, hasReport)
	hasReport, err = k.HasReport(ctx, types.RequestID(2), testapp.Bob.ValAddress)
	require.NoError(t, err)
	require.True(t, hasReport)
	hasReport, err = k.HasReport(ctx, types.RequestID(2), testapp.Carol.ValAddress)
	require.NoError(t, err)
	require.True(t, hasReport)
}
