package keeper

import (
	"cosmossdk.io/collections"
	addresscodec "cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	"github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService storetypes.KVStoreService
	AuthKeeper   authtypes.AccountKeeper
	nftKeeper    nftkeeper.Keeper
	addressCodec addresscodec.Codec

	Schema      collections.Schema
	NFTClassID  collections.Sequence
	ClassOwners collections.Map[string, []byte]

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	authKeeper authtypes.AccountKeeper,
	nftKeeper nftkeeper.Keeper,
	authority string,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		AuthKeeper:   authKeeper,
		nftKeeper:    nftKeeper,
		addressCodec: authKeeper.AddressCodec(),
		authority:    authority,

		NFTClassID:  collections.NewSequence(sb, types.NFTClassIDStoreKeyPrefix, "nft_class_id"),
		ClassOwners: collections.NewMap(sb, types.ClassOwnersStoreKeyPrefix, "class_owners", collections.StringKey, collections.BytesValue),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the x/oracle module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}
