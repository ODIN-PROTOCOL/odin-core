package types

import (
	"time"

	"github.com/ODIN-PROTOCOL/wasmvm/v2/types"
)

var (
	_ types.EnvInterface = (*PrepareEnv)(nil)
	_ types.EnvInterface = (*ExecuteEnv)(nil)
)

// BaseEnv combines shared functions used in prepare and execution Owasm program,
type BaseEnv struct {
	request     Request
	maxSpanSize int64
}

// GetSpanSize implements Owasm ExecEnv interface.
func (env *BaseEnv) GetSpanSize() int64 {
	return env.maxSpanSize
}

// GetCalldata implements Owasm ExecEnv interface.
func (env *BaseEnv) GetCalldata() []byte {
	return env.request.Calldata
}

// SetReturnData implements Owasm ExecEnv interface.
func (env *BaseEnv) SetReturnData(data []byte) error {
	return types.ErrWrongPeriodAction
}

// GetAskCount implements Owasm ExecEnv interface.
func (env *BaseEnv) GetAskCount() int64 {
	return int64(len(env.request.RequestedValidators))
}

// GetMinCount implements Owasm ExecEnv interface.
func (env *BaseEnv) GetMinCount() int64 {
	return int64(env.request.MinCount)
}

// GetPrepareTime implements Owasm ExecEnv interface.
func (env *BaseEnv) GetPrepareTime() int64 {
	return env.request.RequestTime
}

// GetExecuteTime implements Owasm ExecEnv interface.
func (env *BaseEnv) GetExecuteTime() (int64, error) {
	return 0, types.ErrWrongPeriodAction
}

// GetAnsCount implements Owasm ExecEnv interface.
func (env *BaseEnv) GetAnsCount() (int64, error) {
	return 0, types.ErrWrongPeriodAction
}

// AskExternalData implements Owasm ExecEnv interface.
func (env *BaseEnv) AskExternalData(eid int64, did int64, data []byte) error {
	return types.ErrWrongPeriodAction
}

// GetExternalDataStatus implements Owasm ExecEnv interface.
func (env *BaseEnv) GetExternalDataStatus(eid int64, vid int64) (int64, error) {
	return 0, types.ErrWrongPeriodAction
}

// GetExternalData implements Owasm ExecEnv interface.
func (env *BaseEnv) GetExternalData(eid int64, vid int64) ([]byte, error) {
	return nil, types.ErrWrongPeriodAction
}

// PrepareEnv implements ExecEnv interface only expected function and panic on non-prepare functions.
type PrepareEnv struct {
	BaseEnv
	maxCalldataSize int64
	maxRawRequests  int64
	rawRequests     []RawRequest
}

// NewPrepareEnv creates a new environment instance for prepare period.
func NewPrepareEnv(req Request, maxCalldataSize int64, maxRawRequests int64, spanSize int64) *PrepareEnv {
	return &PrepareEnv{
		BaseEnv: BaseEnv{
			request:     req,
			maxSpanSize: spanSize,
		},
		maxCalldataSize: maxCalldataSize,
		maxRawRequests:  maxRawRequests,
	}
}

// AskExternalData implements Owasm ExecEnv interface.
func (env *PrepareEnv) AskExternalData(eid int64, did int64, data []byte) error {
	if int64(len(data)) > env.maxCalldataSize {
		return types.ErrSpanTooSmall
	}
	if int64(len(env.rawRequests)) >= env.maxRawRequests {
		return types.ErrTooManyExternalData
	}
	for _, raw := range env.rawRequests {
		if raw.ExternalID == ExternalID(eid) {
			return types.ErrDuplicateExternalID
		}
	}
	env.rawRequests = append(env.rawRequests, NewRawRequest(
		ExternalID(eid), DataSourceID(did), data,
	))
	return nil
}

// GetRawRequests returns the list of raw requests made during Owasm prepare run.
func (env *PrepareEnv) GetRawRequests() []RawRequest {
	return env.rawRequests
}

// ExecuteEnv implements ExecEnv interface only expected function and panic on prepare related functions.
type ExecuteEnv struct {
	BaseEnv
	reports     map[string]map[ExternalID]RawReport
	Retdata     []byte
	ExecuteTime int64
}

// NewExecuteEnv creates a new environment instance for execution period.
func NewExecuteEnv(req Request, reports []Report, executeTime time.Time, spanSize int64) *ExecuteEnv {
	envReports := make(map[string]map[ExternalID]RawReport)
	for _, report := range reports {
		valReports := make(map[ExternalID]RawReport)
		for _, each := range report.RawReports {
			valReports[each.ExternalID] = each
		}
		envReports[report.Validator] = valReports
	}
	return &ExecuteEnv{
		BaseEnv: BaseEnv{
			request:     req,
			maxSpanSize: spanSize,
		},
		reports:     envReports,
		ExecuteTime: executeTime.Unix(),
	}
}

// GetExecuteTime implements Owasm ExecEnv interface.
func (env *ExecuteEnv) GetExecuteTime() (int64, error) {
	return env.ExecuteTime, nil
}

// GetAnsCount implements Owasm ExecEnv interface.
func (env *ExecuteEnv) GetAnsCount() (int64, error) {
	return int64(len(env.reports)), nil
}

// SetReturnData implements Owasm ExecEnv interface.
func (env *ExecuteEnv) SetReturnData(data []byte) error {
	if env.Retdata != nil {
		return types.ErrRepeatSetReturnData
	}
	env.Retdata = data
	return nil
}

func (env *ExecuteEnv) getExternalDataFull(eid int64, valIdx int64) ([]byte, int64, error) {
	if valIdx < 0 || valIdx >= int64(len(env.request.RequestedValidators)) {
		return nil, 0, types.ErrBadValidatorIndex
	}
	valAddr := env.request.RequestedValidators[valIdx]
	valReports, ok := env.reports[valAddr]
	if !ok {
		return nil, -1, nil
	}
	valReport, ok := valReports[ExternalID(eid)]
	if !ok {
		return nil, 0, types.ErrBadExternalID
	}
	return valReport.Data, int64(valReport.ExitCode), nil
}

// GetExternalDataStatus implements Owasm ExecEnv interface.
func (env *ExecuteEnv) GetExternalDataStatus(eid int64, vid int64) (int64, error) {
	_, status, err := env.getExternalDataFull(eid, vid)
	return status, err
}

// GetExternalData implements Owasm ExecEnv interface.
func (env *ExecuteEnv) GetExternalData(eid int64, vid int64) ([]byte, error) {
	data, _, err := env.getExternalDataFull(eid, vid)
	return data, err
}
