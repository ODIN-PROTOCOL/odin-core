package keepers

import (
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	oraclekeeper "github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v8/keeper"
	wasmlckeeper "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/keeper"

	mintkeeper "github.com/ODIN-PROTOCOL/odin-core/x/mint/keeper"
	wasmkeeper "github.com/ODIN-PROTOCOL/wasmd/x/wasm/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	odinbankkeeper "github.com/ODIN-PROTOCOL/odin-core/x/bank/keeper"
)

type AppKeepers struct {
	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       odinbankkeeper.WrappedBankKeeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    *stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     *crisiskeeper.Keeper
	UpgradeKeeper    *upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	// IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCKeeper             *ibckeeper.Keeper
	ICAHostKeeper         icahostkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	FeegrantKeeper        feegrantkeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	NFTKeeper             nftkeeper.Keeper
	CircuitKeeper         circuitkeeper.Keeper
	OracleKeeper          oraclekeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	WasmKeeper            wasmkeeper.Keeper
	WasmClientKeeper      wasmlckeeper.Keeper
	ICQKeeper             icqkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper  capabilitykeeper.ScopedKeeper
	ScopedOracleKeeper   capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper
	ScopedICQKeeper      capabilitykeeper.ScopedKeeper
}
