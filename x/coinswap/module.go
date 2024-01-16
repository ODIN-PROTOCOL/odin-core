package coinswap

import (
	"context"
	"encoding/json"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	coinswapcli "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/client/cli"
	coinswapkeeper "github.com/ODIN-PROTOCOL/odin-core/x/coinswap/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/coinswap/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is Odin Oracle's module basic object.
type AppModuleBasic struct{}

// Name returns this module's name - "coinswap" (SDK AppModuleBasic interface).
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the staking module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (b AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns the default genesis state as raw bytes.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(coinswapkeeper.DefaultGenesisState())
}

// ValidateGenesis checks of the Genesis
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	err := cdc.UnmarshalJSON(bz, &gs)
	if err != nil {
		return err
	}
	return gs.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the oracle module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns cobra CLI command to send txs for this module (SDK AppModuleBasic interface).
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return coinswapcli.GetTxCmd()
}

// GetQueryCmd returns cobra CLI command to query chain state (SDK AppModuleBasic interface).
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return coinswapcli.GetQueryCmd()
}

// AppModule represents the AppModule for this module.
type AppModule struct {
	AppModuleBasic
	keeper coinswapkeeper.Keeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(k coinswapkeeper.Keeper) AppModule {
	return AppModule{
		keeper: k,
	}
}

// RegisterInvariants is a noop function to satisfy SDK AppModule interface.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// QuerierRoute returns the oracle module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), coinswapkeeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), coinswapkeeper.Querier{Keeper: am.keeper})
}

// BeginBlock processes ABCI begin block message for this oracle module (SDK AppModule interface).
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

// EndBlock processes ABCI end block message for this oracle module (SDK AppModule interface).
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// InitGenesis performs genesis initialization for the oracle module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	coinswapkeeper.InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the current state as genesis raw bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := coinswapkeeper.ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(&gs)
}

// ConsensusVersion returns the current module store definitions version
func (am AppModule) ConsensusVersion() uint64 {
	return types.ModuleVersion
}
