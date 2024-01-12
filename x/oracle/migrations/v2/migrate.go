package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/exported"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

const (
	ModuleName = "oracle"
)

// Migrate migrates the x/oracle module state from the consensus version 1 to
// version 2. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/oracle
// module state.
func Migrate(
	ctx sdk.Context,
	store sdk.KVStore,
	legacySubspace exported.Subspace,
	cdc codec.BinaryCodec,
) error {
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&currParams)
	store.Set(types.ParamsKeyPrefix, bz)

	return nil
}
