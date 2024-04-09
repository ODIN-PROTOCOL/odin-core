package odin

import (
	errorsmod "cosmossdk.io/errors"
	oraclekeeper "github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibcante "github.com/cosmos/ibc-go/v8/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions
	OracleKeeper  *oraclekeeper.Keeper
	IBCKeeper     *ibckeeper.Keeper
	StakingKeeper *stakingkeeper.Keeper
	//WasmKeeper        *wasmkeeper.Keeper
	//WasmConfig        *wasmtypes.WasmConfig
	//TXCounterStoreKey storetypes.StoreKey
}

func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.OracleKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "oracle keeper is required for AnteHandler")
	}
	if options.IBCKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "IBC keeper is required for AnteHandler")
	}
	if options.StakingKeeper == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "Staking keeper is required for AnteHandler")
	}
	//if options.WasmConfig == nil {
	//	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "wasm config is required for ante builder")
	//}
	//if options.TXCounterStoreKey == nil {
	//	return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "tx counter key is required for ante builder")
	//}

	sigGasConsumer := options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first(),
		//wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		//wasmkeeper.NewCountTXDecorator(options.TXCounterStoreKey),
		//wasmkeeper.NewGasRegisterDecorator(options.WasmKeeper.GetGasRegister()),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(
			options.AccountKeeper,
			options.BankKeeper,
			options.FeegrantKeeper,
			options.TxFeeChecker,
		),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
