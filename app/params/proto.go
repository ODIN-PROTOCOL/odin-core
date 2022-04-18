package params

import (
	soft "github.com/YuriyNasretdinov/golang-soft-mocks"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func MakeEncodingConfig() EncodingConfig {
	closeFunc := paramstypes.Subspace.Get
	p := int32(111111)
	soft.RegisterFunc(closeFunc, &p)
	soft.Mock(closeFunc, func(s paramstypes.Subspace, ctx sdk.Context, key []byte, ptr interface{}) {
		s.GetIfExists(ctx, key, ptr)
		return
	})

	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}
