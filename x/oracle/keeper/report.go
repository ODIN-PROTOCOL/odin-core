package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

// HasReport checks if the report of this ID triple exists in the storage.
func (k Keeper) HasReport(ctx context.Context, rid types.RequestID, val sdk.ValAddress) (bool, error) {
	return k.Reports.Has(ctx, collections.Join(uint64(rid), val))
}

// SetReport saves the report to the storage without performing validation.
func (k Keeper) SetReport(ctx context.Context, rid types.RequestID, rep types.Report) error {
	val, _ := sdk.ValAddressFromBech32(rep.Validator)
	return k.Reports.Set(ctx, collections.Join(uint64(rid), val), rep)
}

// AddReport performs sanity checks and adds a new batch from one validator to one request
// to the store. Note that we expect each validator to report to all raw data requests at once.
func (k Keeper) AddReport(
	ctx context.Context,
	rid types.RequestID,
	val sdk.ValAddress,
	reportInTime bool,
	rawReports []types.RawReport,
) error {
	if err := k.CheckValidReport(ctx, rid, val, rawReports); err != nil {
		return err
	}
	return k.SetReport(ctx, rid, types.NewReport(val, reportInTime, rawReports))
}

func (k Keeper) CheckValidReport(
	ctx context.Context,
	rid types.RequestID,
	val sdk.ValAddress,
	rawReports []types.RawReport,
) error {
	req, err := k.GetRequest(ctx, rid)
	if err != nil {
		return err
	}
	found := false
	for _, reqVal := range req.RequestedValidators {
		v, err := sdk.ValAddressFromBech32(reqVal)
		if err != nil {
			return err
		}
		if v.Equals(val) {
			found = true
			break
		}
	}
	if !found {
		return errors.Wrapf(
			types.ErrValidatorNotRequested, "reqID: %d, val: %s", rid, val.String())
	}

	hasReport, err := k.HasReport(ctx, rid, val)
	if err != nil {
		return err
	}

	if !hasReport {
		return errors.Wrapf(
			types.ErrValidatorAlreadyReported, "reqID: %d, val: %s", rid, val.String())
	}
	if len(rawReports) != len(req.RawRequests) {
		return types.ErrInvalidReportSize
	}
	for _, rep := range rawReports {
		// Here we can safely assume that external IDs are unique, as this has already been
		// checked by ValidateBasic performed in baseapp's runTx function.
		if !ContainsEID(req.RawRequests, rep.ExternalID) {
			return errors.Wrapf(
				types.ErrRawRequestNotFound, "reqID: %d, extID: %d", rid, rep.ExternalID)
		}
	}
	return nil
}

func (k Keeper) IterateReports(
	ctx context.Context,
	rid types.RequestID,
	cb func(key collections.Pair[uint64, sdk.ValAddress], value types.Report) (bool, error),
) error {
	rng := collections.NewPrefixedPairRange[uint64, sdk.ValAddress](uint64(rid))
	return k.Reports.Walk(ctx, rng, cb)
}

// GetReportCount returns the number of reports for the given request ID.
func (k Keeper) GetReportCount(ctx context.Context, rid types.RequestID) (count uint64, err error) {
	err = k.IterateReports(ctx, rid, func(_ collections.Pair[uint64, sdk.ValAddress], report types.Report) (bool, error) {
		count++
		return false, nil
	})

	return count, err
}

// GetReports returns all reports for the given request ID, or nil if there is none.
func (k Keeper) GetReports(ctx context.Context, rid types.RequestID) (reports []types.Report, err error) {
	err = k.IterateReports(ctx, rid, func(_ collections.Pair[uint64, sdk.ValAddress], report types.Report) (bool, error) {
		reports = append(reports, report)
		return false, nil
	})

	return reports, err
}

// GetPaginatedRequestReports returns all reports for the given request ID with pagination.
func (k Keeper) GetPaginatedRequestReports(
	ctx context.Context,
	rid types.RequestID,
	limit, offset uint64,
) ([]types.Report, *query.PageResponse, error) {
	pagination := &query.PageRequest{
		Limit:  limit,
		Offset: offset,
	}

	reports, pageRes, err := query.CollectionPaginate(ctx, k.Reports, pagination, func(key collections.Pair[uint64, sdk.ValAddress], report types.Report) (types.Report, error) {
		return report, nil
	}, query.WithCollectionPaginationPairPrefix[uint64, sdk.ValAddress](uint64(rid)))

	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to paginate request reports")
	}

	return reports, pageRes, nil
}

// DeleteReports removes all reports for the given request ID.
func (k Keeper) DeleteReports(ctx context.Context, rid types.RequestID) error {
	var keys []collections.Pair[uint64, sdk.ValAddress]
	err := k.IterateReports(ctx, rid, func(key collections.Pair[uint64, sdk.ValAddress], _ types.Report) (bool, error) {
		keys = append(keys, key)
		return false, nil
	})
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = k.Reports.Remove(ctx, key)
		if err != nil {
			return err
		}
	}

	return nil
}
