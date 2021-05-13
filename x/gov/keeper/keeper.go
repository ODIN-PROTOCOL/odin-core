package keeper

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	govtypes "github.com/GeoDB-Limited/odin-core/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Keeper defines the governance module Keeper
type Keeper struct {
	// The reference to the Paramstore to get and set gov specific params
	paramSpace govtypes.ParamSubspace

	authKeeper govtypes.AccountKeeper
	bankKeeper govtypes.BankKeeper

	// The reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	stakingKeeper govtypes.StakingKeeper

	mintKeeper govtypes.MintKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc codec.BinaryMarshaler

	// Proposal router
	router govtypes.Router
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
//
// CONTRACT: the parameter Subspace must have the param key table already initialized
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace govtypes.ParamSubspace,
	authKeeper govtypes.AccountKeeper, bankKeeper govtypes.BankKeeper, sk govtypes.StakingKeeper, mintKeeper govtypes.MintKeeper, rtr govtypes.Router,
) Keeper {

	// ensure governance module account is set
	if addr := authKeeper.GetModuleAddress(govtypes.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", govtypes.ModuleName))
	}

	// It is vital to seal the governance proposal router here as to not allow
	// further handlers to be registered after the keeper is created since this
	// could create invalid or non-deterministic behavior.
	rtr.Seal()

	return Keeper{
		storeKey:      key,
		paramSpace:    paramSpace,
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: sk,
		mintKeeper:    mintKeeper,
		cdc:           cdc,
		router:        rtr,
	}
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+govtypes.ModuleName)
}

// Router returns the gov Keeper's Router
func (keeper Keeper) Router() govtypes.Router {
	return keeper.router
}

// GetGovernanceAccount returns the governance ModuleAccount
func (keeper Keeper) GetGovernanceAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return keeper.authKeeper.GetModuleAccount(ctx, govtypes.ModuleName)
}

// ProposalQueues

// InsertActiveProposalQueue inserts a ProposalID into the active proposal queue at endTime
func (keeper Keeper) InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := govtypes.GetProposalIDBytes(proposalID)
	store.Set(govtypes.ActiveProposalQueueKey(proposalID, endTime), bz)
}

// RemoveFromActiveProposalQueue removes a proposalID from the Active Proposal Queue
func (keeper Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(govtypes.ActiveProposalQueueKey(proposalID, endTime))
}

// InsertInactiveProposalQueue Inserts a ProposalID into the inactive proposal queue at endTime
func (keeper Keeper) InsertInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := govtypes.GetProposalIDBytes(proposalID)
	store.Set(govtypes.InactiveProposalQueueKey(proposalID, endTime), bz)
}

// RemoveFromInactiveProposalQueue removes a proposalID from the Inactive Proposal Queue
func (keeper Keeper) RemoveFromInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(govtypes.InactiveProposalQueueKey(proposalID, endTime))
}

// Iterators

// IterateActiveProposalsQueue iterates over the proposals in the active proposal queue
// and performs a callback function
func (keeper Keeper) IterateActiveProposalsQueue(ctx sdk.Context, endTime time.Time, cb func(proposal govtypes.Proposal) (stop bool)) {
	iterator := keeper.ActiveProposalQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := govtypes.SplitActiveProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// IterateInactiveProposalsQueue iterates over the proposals in the inactive proposal queue
// and performs a callback function
func (keeper Keeper) IterateInactiveProposalsQueue(ctx sdk.Context, endTime time.Time, cb func(proposal govtypes.Proposal) (stop bool)) {
	iterator := keeper.InactiveProposalQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := govtypes.SplitInactiveProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// ActiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Active Queue that expire by endTime
func (keeper Keeper) ActiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(govtypes.ActiveProposalQueuePrefix, sdk.PrefixEndBytes(govtypes.ActiveProposalByTimeKey(endTime)))
}

// InactiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Inactive Queue that expire by endTime
func (keeper Keeper) InactiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(govtypes.InactiveProposalQueuePrefix, sdk.PrefixEndBytes(govtypes.InactiveProposalByTimeKey(endTime)))
}
