package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/ODIN-PROTOCOL/odin-core/testing/testapp"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

func testRequest(
	t *testing.T,
	k keeper.Keeper,
	ctx sdk.Context,
	rid types.RequestID,
	resolveStatus types.ResolveStatus,
	reportCount uint64,
	hasRequest bool,
) {
	if resolveStatus == types.RESOLVE_STATUS_OPEN {
		hasResult, err := k.HasResult(ctx, rid)
		require.NoError(t, err)
		require.False(t, hasResult)
	} else {
		r, err := k.GetResult(ctx, rid)
		require.NoError(t, err)
		require.NotNil(t, r)
		require.Equal(t, resolveStatus, r.ResolveStatus)
	}

	reportCount, err := k.GetReportCount(ctx, rid)
	require.NoError(t, err)
	require.Equal(t, reportCount, reportCount)

	hasRequest1, err := k.HasRequest(ctx, rid)
	require.NoError(t, err)
	require.Equal(t, hasRequest, hasRequest1)
}

func TestHasRequest(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// We should not have a request ID 42 without setting it.
	hasRequest, err := k.HasRequest(ctx, 42)
	require.NoError(t, err)
	require.False(t, hasRequest)
	// After we set it, we should be able to find it.
	err = k.SetRequest(ctx, 42, types.NewRequest(1, BasicCalldata, nil, 1, 1, testapp.ParseTime(0), "", nil, nil, 0))
	require.NoError(t, err)

	hasRequest, err = k.HasRequest(ctx, 42)
	require.NoError(t, err)
	require.True(t, hasRequest)
}

func TestDeleteRequest(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// After we set it, we should be able to find it.
	err := k.SetRequest(ctx, 42, types.NewRequest(1, BasicCalldata, nil, 1, 1, testapp.ParseTime(0), "", nil, nil, 0))
	require.NoError(t, err)

	hasRequest, err := k.HasRequest(ctx, 42)
	require.NoError(t, err)
	require.True(t, hasRequest)
	// After we delete it, we should not find it anymore.
	err = k.DeleteRequest(ctx, 42)
	require.NoError(t, err)

	hasRequest, err = k.HasRequest(ctx, 42)
	require.NoError(t, err)
	require.False(t, hasRequest)
	_, err = k.GetRequest(ctx, 42)
	require.ErrorIs(t, err, types.ErrRequestNotFound)
	require.Panics(t, func() { _ = k.MustGetRequest(ctx, 42) })
}

func TestSetterGetterRequest(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Getting a non-existent request should return error.
	_, err := k.GetRequest(ctx, 42)
	require.ErrorIs(t, err, types.ErrRequestNotFound)
	require.Panics(t, func() { _ = k.MustGetRequest(ctx, 42) })
	// Creates some basic requests.
	req1 := types.NewRequest(1, BasicCalldata, nil, 1, 1, testapp.ParseTime(0), "", nil, nil, 0)
	req2 := types.NewRequest(2, BasicCalldata, nil, 1, 1, testapp.ParseTime(0), "", nil, nil, 0)
	// Sets id 42 with request 1 and id 42 with request 2.
	err = k.SetRequest(ctx, 42, req1)
	require.NoError(t, err)

	err = k.SetRequest(ctx, 43, req2)
	require.NoError(t, err)

	// Checks that Get and MustGet perform correctly.
	req1Res, err := k.GetRequest(ctx, 42)
	require.Nil(t, err)
	require.Equal(t, req1, req1Res)
	require.Equal(t, req1, k.MustGetRequest(ctx, 42))
	req2Res, err := k.GetRequest(ctx, 43)
	require.Nil(t, err)
	require.Equal(t, req2, req2Res)
	require.Equal(t, req2, k.MustGetRequest(ctx, 43))
	// Replaces id 42 with another request.
	err = k.SetRequest(ctx, 42, req2)
	require.NoError(t, err)

	require.NotEqual(t, req1, k.MustGetRequest(ctx, 42))
	require.Equal(t, req2, k.MustGetRequest(ctx, 42))
}

func TestSetterGettterPendingResolveList(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)

	// Initially, we should get an empty list of pending resolves.
	pendingResolveList, err := k.GetPendingResolveList(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingResolveList, []types.RequestID(nil))

	// After we set something, we should get that thing back.
	err = k.SetPendingResolveList(ctx, []types.RequestID{5, 6, 7, 8})
	require.NoError(t, err)

	pendingResolveList, err = k.GetPendingResolveList(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingResolveList, []types.RequestID{5, 6, 7, 8})

	// Let's also try setting it back to empty list.
	err = k.SetPendingResolveList(ctx, []types.RequestID(nil))
	require.NoError(t, err)

	pendingResolveList, err = k.GetPendingResolveList(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingResolveList, []types.RequestID(nil))

	// Nil should also works.
	err = k.SetPendingResolveList(ctx, nil)
	require.NoError(t, err)

	pendingResolveList, err = k.GetPendingResolveList(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingResolveList, []types.RequestID(nil))
}

func TestAddDataSourceBasic(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)

	// We start by setting an oracle request available at ID 42.
	err := k.SetOracleScript(ctx, 42, types.NewOracleScript(
		testapp.Owner.Address, BasicName, BasicDesc, BasicFilename, BasicSchema, BasicSourceCodeURL,
	))
	require.NoError(t, err)

	// Adding the first request should return ID 1.
	id, err := k.AddRequest(
		ctx,
		types.NewRequest(42, BasicCalldata, []sdk.ValAddress{}, 1, 1, testapp.ParseTime(0), "", nil, nil, 0),
	)
	require.NoError(t, err)
	require.Equal(t, id, types.RequestID(1))

	// Adding another request should return ID 2.
	id, err = k.AddRequest(
		ctx,
		types.NewRequest(42, BasicCalldata, []sdk.ValAddress{}, 1, 1, testapp.ParseTime(0), "", nil, nil, 0),
	)
	require.NoError(t, err)
	require.Equal(t, id, types.RequestID(2))
}

func TestAddPendingResolveList(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)

	// Initially, we should get an empty list of pending resolves.
	pendingResolveList, err := k.GetPendingResolveList(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingResolveList, []types.RequestID(nil))

	// Everytime we append a new request ID, it should show up.
	err = k.AddPendingRequest(ctx, 42)
	require.NoError(t, err)

	pendingResolveList, err = k.GetPendingResolveList(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingResolveList, []types.RequestID{42})

	err = k.AddPendingRequest(ctx, 43)
	require.NoError(t, err)

	pendingResolveList, err = k.GetPendingResolveList(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingResolveList, []types.RequestID{42, 43})
}

func TestProcessExpiredRequests(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	params, err := k.GetParams(ctx)
	require.NoError(t, err)
	params.ExpirationBlockCount = 3
	err = k.SetParams(ctx, params)
	require.NoError(t, err)

	// Set some initial requests. All requests are asked to validators 1 & 2.
	req1 := defaultRequest()
	req1.RequestHeight = 5
	req2 := defaultRequest()
	req2.RequestHeight = 6
	req3 := defaultRequest()
	req3.RequestHeight = 6
	req4 := defaultRequest()
	req4.RequestHeight = 10
	_, err = k.AddRequest(ctx, req1)
	require.NoError(t, err)
	_, err = k.AddRequest(ctx, req2)
	require.NoError(t, err)
	_, err = k.AddRequest(ctx, req3)
	require.NoError(t, err)
	_, err = k.AddRequest(ctx, req4)
	require.NoError(t, err)

	// Initially all validators are active.
	validatorStatus, err := k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.True(t, validatorStatus.IsActive)
	validatorStatus, err = k.GetValidatorStatus(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.True(t, validatorStatus.IsActive)

	// Validator 1 reports all requests. Validator 2 misses request#3.
	rawReports := []types.RawReport{types.NewRawReport(42, 0, BasicReport), types.NewRawReport(43, 0, BasicReport)}
	err = k.AddReport(ctx, 1, testapp.Validators[0].ValAddress, false, rawReports)
	require.NoError(t, err)
	err = k.AddReport(ctx, 2, testapp.Validators[0].ValAddress, true, rawReports)
	require.NoError(t, err)
	err = k.AddReport(ctx, 3, testapp.Validators[0].ValAddress, false, rawReports)
	require.NoError(t, err)
	err = k.AddReport(ctx, 4, testapp.Validators[0].ValAddress, true, rawReports)
	require.NoError(t, err)
	err = k.AddReport(ctx, 1, testapp.Validators[1].ValAddress, true, rawReports)
	require.NoError(t, err)
	err = k.AddReport(ctx, 2, testapp.Validators[1].ValAddress, true, rawReports)
	require.NoError(t, err)
	err = k.AddReport(ctx, 4, testapp.Validators[1].ValAddress, true, rawReports)
	require.NoError(t, err)

	// Request 1, 2 and 4 gets resolved. Request 3 does not.
	err = k.ResolveSuccess(ctx, 1, BasicResult, 1234)
	require.NoError(t, err)
	err = k.ResolveFailure(ctx, 2, "ARBITRARY_REASON")
	require.NoError(t, err)
	err = k.ResolveSuccess(ctx, 4, BasicResult, 1234)
	require.NoError(t, err)

	// Initially, last expired request ID should be 0.
	requestLastExpired, err := k.GetRequestLastExpired(ctx)
	require.NoError(t, err)
	require.Equal(t, types.RequestID(0), requestLastExpired)

	// At block 7, nothing should happen.
	ctx = ctx.WithBlockHeight(7).WithBlockTime(testapp.ParseTime(7000)).WithEventManager(sdk.NewEventManager())
	err = k.ProcessExpiredRequests(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.Events{}, ctx.EventManager().Events())

	requestLastExpired, err = k.GetRequestLastExpired(ctx)
	require.NoError(t, err)
	require.Equal(t, types.RequestID(0), requestLastExpired)

	validatorStatus, err = k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.True(t, validatorStatus.IsActive)
	validatorStatus, err = k.GetValidatorStatus(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.True(t, validatorStatus.IsActive)

	testRequest(t, k, ctx, types.RequestID(1), types.RESOLVE_STATUS_SUCCESS, 2, true)
	testRequest(t, k, ctx, types.RequestID(2), types.RESOLVE_STATUS_FAILURE, 2, true)
	testRequest(t, k, ctx, types.RequestID(3), types.RESOLVE_STATUS_OPEN, 1, true)
	testRequest(t, k, ctx, types.RequestID(4), types.RESOLVE_STATUS_SUCCESS, 2, true)

	// At block 8, now last request ID should move to 1. No events should be emitted.
	ctx = ctx.WithBlockHeight(8).WithBlockTime(testapp.ParseTime(8000)).WithEventManager(sdk.NewEventManager())
	err = k.ProcessExpiredRequests(ctx)
	require.NoError(t, err)

	require.Equal(t, sdk.Events{}, ctx.EventManager().Events())

	requestLastExpired, err = k.GetRequestLastExpired(ctx)
	require.NoError(t, err)
	require.Equal(t, types.RequestID(1), requestLastExpired)

	validatorStatus, err = k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.True(t, validatorStatus.IsActive)
	validatorStatus, err = k.GetValidatorStatus(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.True(t, validatorStatus.IsActive)

	testRequest(t, k, ctx, types.RequestID(1), types.RESOLVE_STATUS_SUCCESS, 0, false)
	testRequest(t, k, ctx, types.RequestID(2), types.RESOLVE_STATUS_FAILURE, 2, true)
	testRequest(t, k, ctx, types.RequestID(3), types.RESOLVE_STATUS_OPEN, 1, true)
	testRequest(t, k, ctx, types.RequestID(4), types.RESOLVE_STATUS_SUCCESS, 2, true)

	// At block 9, request#3 is expired and validator 2 becomes inactive.
	ctx = ctx.WithBlockHeight(9).WithBlockTime(testapp.ParseTime(9000)).WithEventManager(sdk.NewEventManager())
	err = k.ProcessExpiredRequests(ctx)
	require.NoError(t, err)

	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, "3"),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, "3"),
	), sdk.NewEvent(
		types.EventTypeDeactivate,
		sdk.NewAttribute(types.AttributeKeyValidator, testapp.Validators[1].ValAddress.String()),
	)}, ctx.EventManager().Events())

	requestLastExpired, err = k.GetRequestLastExpired(ctx)
	require.NoError(t, err)
	require.Equal(t, types.RequestID(3), requestLastExpired)

	validatorStatus, err = k.GetValidatorStatus(ctx, testapp.Validators[0].ValAddress)
	require.NoError(t, err)
	require.True(t, validatorStatus.IsActive)
	validatorStatus, err = k.GetValidatorStatus(ctx, testapp.Validators[1].ValAddress)
	require.NoError(t, err)
	require.False(t, validatorStatus.IsActive)

	require.Equal(t, types.NewResult(
		BasicClientID, req3.OracleScriptID, req3.Calldata, uint64(len(req3.RequestedValidators)), req3.MinCount,
		3, 1, req3.RequestTime, testapp.ParseTime(9000).Unix(),
		types.RESOLVE_STATUS_EXPIRED, nil,
	), k.MustGetResult(ctx, 3))
	testRequest(t, k, ctx, types.RequestID(1), types.RESOLVE_STATUS_SUCCESS, 0, false)
	testRequest(t, k, ctx, types.RequestID(2), types.RESOLVE_STATUS_FAILURE, 0, false)
	testRequest(t, k, ctx, types.RequestID(3), types.RESOLVE_STATUS_EXPIRED, 0, false)
	testRequest(t, k, ctx, types.RequestID(4), types.RESOLVE_STATUS_SUCCESS, 2, true)

	// At block 10, nothing should happen
	ctx = ctx.WithBlockHeight(10).WithBlockTime(testapp.ParseTime(10000)).WithEventManager(sdk.NewEventManager())
	err = k.ProcessExpiredRequests(ctx)
	require.NoError(t, err)

	require.Equal(t, sdk.Events{}, ctx.EventManager().Events())

	requestLastExpired, err = k.GetRequestLastExpired(ctx)
	require.NoError(t, err)
	require.Equal(t, types.RequestID(3), requestLastExpired)

	testRequest(t, k, ctx, types.RequestID(1), types.RESOLVE_STATUS_SUCCESS, 0, false)
	testRequest(t, k, ctx, types.RequestID(2), types.RESOLVE_STATUS_FAILURE, 0, false)
	testRequest(t, k, ctx, types.RequestID(3), types.RESOLVE_STATUS_EXPIRED, 0, false)
	testRequest(t, k, ctx, types.RequestID(4), types.RESOLVE_STATUS_SUCCESS, 2, true)

	// At block 13, last expired request becomes 4.
	ctx = ctx.WithBlockHeight(13).WithBlockTime(testapp.ParseTime(13000)).WithEventManager(sdk.NewEventManager())
	err = k.ProcessExpiredRequests(ctx)
	require.NoError(t, err)

	require.Equal(t, sdk.Events{}, ctx.EventManager().Events())

	requestLastExpired, err = k.GetRequestLastExpired(ctx)
	require.NoError(t, err)
	require.Equal(t, types.RequestID(4), requestLastExpired)
	testRequest(t, k, ctx, types.RequestID(1), types.RESOLVE_STATUS_SUCCESS, 0, false)
	testRequest(t, k, ctx, types.RequestID(2), types.RESOLVE_STATUS_FAILURE, 0, false)
	testRequest(t, k, ctx, types.RequestID(3), types.RESOLVE_STATUS_EXPIRED, 0, false)
	testRequest(t, k, ctx, types.RequestID(4), types.RESOLVE_STATUS_SUCCESS, 0, false)
}
