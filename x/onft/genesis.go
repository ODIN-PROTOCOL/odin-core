package onft

import (
	"github.com/ODIN-PROTOCOL/odin-core/x/onft/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis performs genesis initialization for the oracle module.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data *types.GenesisState) {
	err := k.NFTClassID.Set(ctx, data.NftClassId)
	if err != nil {
		panic(err)
	}

	codec := k.AuthKeeper.AddressCodec()

	for class, owner := range data.ClassOwners {
		ownerAddr, err := codec.StringToBytes(owner)
		if err != nil {
			panic(err)
		}

		err = k.ClassOwners.Set(ctx, class, ownerAddr)
		if err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (*types.GenesisState, error) {
	nftClassID, err := k.NFTClassID.Peek(ctx)
	if err != nil {
		return nil, err
	}

	codec := k.AuthKeeper.AddressCodec()

	classOwners := make(map[string]string)

	err = k.ClassOwners.Walk(ctx, nil, func(class string, owner []byte) (stop bool, err error) {
		ownerAddr, err := codec.BytesToString(owner)
		if err != nil {
			return false, err
		}

		classOwners[class] = ownerAddr

		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.GenesisState{
		NftClassId:  nftClassID,
		ClassOwners: classOwners,
	}, nil
}
