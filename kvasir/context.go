package kvasir

import (
	"sync/atomic"
	"time"

	"github.com/ODIN-PROTOCOL/odin-core/kvasir/executor"
	wasmtypes "github.com/ODIN-PROTOCOL/wasmd/x/wasm/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
)

type Request struct {
	ChosenValidators []string  `json:"chosen_validators"`
	Metadata         string    `json:"metadata"`
	PartsReceived    uint32    `json:"parts_received"`
	PartsRequested   uint32    `json:"parts_requested"`
	Payed            sdk.Coins `json:"payed"`
	RequestHeight    uint64    `json:"request_Height"`
	RequestID        uint64    `json:"request_id"`
	Sender           string    `json:"sender"`
	Status           string    `json:"status"`
}

type RequestKey struct {
	ContractAddress string
	RequestID       uint64
}

type ReportMsgWithKey struct {
	result          []byte
	contractAddress string
	msg             *wasmtypes.RawContractMessage
	execVersion     string
	keyIndex        int64
	request         Request
}

type Context struct {
	client           rpcclient.Client
	validator        sdk.ValAddress
	validatorAccAddr sdk.AccAddress
	gasPrices        string
	keys             []*keyring.Record
	executor         executor.Executor
	broadcastTimeout time.Duration
	maxTry           uint64
	rpcPollInterval  time.Duration
	maxReport        uint64
	ipfs             string
	grpc             *grpc.ClientConn

	pendingMsgs        chan ReportMsgWithKey
	freeKeys           chan int64
	keyRoundRobinIndex int64 // Must use in conjunction with sync/atomic

	pendingRequests map[RequestKey]bool

	contracts []string

	metricsEnabled bool
	handlingGauge  int64
	pendingGauge   int64
	errorCount     int64
	submittedCount int64
	home           string
}

func (c *Context) nextKeyIndex() int64 {
	keyIndex := atomic.AddInt64(&c.keyRoundRobinIndex, 1) % int64(len(c.keys))
	return keyIndex
}

func (c *Context) updateHandlingGauge(amount int64) {
	if c.metricsEnabled {
		atomic.AddInt64(&c.handlingGauge, amount)
	}
}

func (c *Context) updatePendingGauge(amount int64) {
	if c.metricsEnabled {
		atomic.AddInt64(&c.pendingGauge, amount)
	}
}

func (c *Context) updateErrorCount(amount int64) {
	if c.metricsEnabled {
		atomic.AddInt64(&c.errorCount, amount)
	}
}

func (c *Context) updateSubmittedCount(amount int64) {
	if c.metricsEnabled {
		atomic.AddInt64(&c.submittedCount, amount)
	}
}
