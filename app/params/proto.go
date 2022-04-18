package params

import (
	soft "github.com/YuriyNasretdinov/golang-soft-mocks"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func MakeEncodingConfig() EncodingConfig {
	closeFunc := (*codec.LegacyAmino).UnmarshalJSON
	p := int32(111111)
	soft.RegisterFunc(closeFunc, &p)
	soft.Mock(closeFunc, func(amino *codec.LegacyAmino, bz []byte, ptr interface{}) error {
		if bz == nil {
			ptr = nil
			return nil
		}
		res, _ := soft.CallOriginal(closeFunc, amino, bz, ptr)[0].(error)
		return res
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
