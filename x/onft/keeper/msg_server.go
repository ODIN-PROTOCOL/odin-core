package keeper

import (
	"bytes"
	"context"
	"fmt"

	"cosmossdk.io/x/nft"
	"github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the mint MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) CreateNFTClass(ctx context.Context, msg *types.MsgCreateNFTClass) (*types.MsgCreateNFTClassResponse, error) {
	classID, err := m.NFTClassID.Next(ctx)
	if err != nil {
		return nil, err
	}

	sender, err := m.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s%d", types.NFTClassPrefix, classID)

	// TODO: add checks
	err = m.nftKeeper.SaveClass(ctx, nft.Class{
		Id:          id,
		Name:        msg.Name,
		Symbol:      msg.Symbol,
		Description: msg.Description,
		Uri:         msg.Uri,
		UriHash:     msg.UriHash,
		Data:        msg.Data,
	})
	if err != nil {
		return nil, err
	}

	err = m.ClassOwners.Set(ctx, id, sender)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateNFTClass,
			sdk.NewAttribute(types.AttributeKeyClassID, id),
			sdk.NewAttribute(types.AttributeKeyUri, msg.Uri),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	)

	return &types.MsgCreateNFTClassResponse{Id: id}, nil
}

func (m msgServer) TransferClassOwnership(ctx context.Context, msg *types.MsgTransferClassOwnership) (*types.MsgTransferClassOwnershipResponse, error) {
	if !m.nftKeeper.HasClass(ctx, msg.ClassId) {
		return nil, types.ErrClassNotFound
	}

	owner, err := m.ClassOwners.Get(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}

	sender, err := m.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(sender, owner) {
		return nil, types.ErrSenderNotAuthorized
	}

	newOwner, err := m.addressCodec.StringToBytes(msg.NewOwner)
	if err != nil {
		return nil, err
	}

	err = m.ClassOwners.Set(ctx, msg.ClassId, newOwner)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransferClassOwnership,
			sdk.NewAttribute(types.AttributeKeyClassID, msg.ClassId),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.NewOwner),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	)

	return &types.MsgTransferClassOwnershipResponse{}, nil
}

func (m msgServer) MintNFT(ctx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	if !m.nftKeeper.HasClass(ctx, msg.ClassId) {
		return nil, types.ErrClassNotFound
	}

	owner, err := m.ClassOwners.Get(ctx, msg.ClassId)
	if err != nil {
		return nil, err
	}

	sender, err := m.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(sender, owner) {
		return nil, types.ErrSenderNotAuthorized
	}

	receiver, err := m.addressCodec.StringToBytes(msg.Receiver)
	if err != nil {
		return nil, err
	}

	nftID := fmt.Sprintf("%d", m.nftKeeper.GetTotalSupply(ctx, msg.ClassId))

	err = m.nftKeeper.Mint(ctx, nft.NFT{
		ClassId: msg.ClassId,
		Id:      nftID,
		Uri:     msg.Uri,
		UriHash: msg.UriHash,
		Data:    msg.Data,
	}, receiver)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMintNFT,
			sdk.NewAttribute(types.AttributeKeyNFTID, nftID),
			sdk.NewAttribute(types.AttributeKeyClassID, msg.ClassId),
			sdk.NewAttribute(types.AttributeKeyUri, msg.Uri),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	)

	return &types.MsgMintNFTResponse{Id: nftID}, nil
}
