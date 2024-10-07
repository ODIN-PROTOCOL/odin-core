package keeper

import (
	"context"

	"cosmossdk.io/x/nft"
	onfttypes "github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ onfttypes.QueryServer = Keeper{}

func (k Keeper) ClassOwner(ctx context.Context, r *onfttypes.QueryClassOwnerRequest) (*onfttypes.QueryClassOwnerResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	if len(r.ClassId) == 0 {
		return nil, nft.ErrEmptyClassID
	}

	owner, err := k.ClassOwners.Get(ctx, r.ClassId)
	if err != nil {
		return nil, err
	}

	ownerStr, err := k.addressCodec.BytesToString(owner)
	if err != nil {
		return nil, err
	}

	return &onfttypes.QueryClassOwnerResponse{Owner: ownerStr}, nil
}

func (k Keeper) NFTs(ctx context.Context, r *onfttypes.QueryNFTsRequest) (*onfttypes.QueryNFTsResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	nfts, err := k.nftKeeper.NFTs(ctx, &nft.QueryNFTsRequest{Pagination: r.Pagination})
	if err != nil {
		return nil, err
	}

	nftsOwners := make([]*onfttypes.NFT, len(nfts.Nfts))
	for i, n := range nfts.Nfts {
		ownerStr := ""
		owner := k.nftKeeper.GetOwner(ctx, n.ClassId, n.Id)
		if !owner.Empty() {
			ownerStr, _ = k.addressCodec.BytesToString(owner)
		}

		nftsOwners[i] = &onfttypes.NFT{
			Id:      n.Id,
			ClassId: n.ClassId,
			Uri:     n.Uri,
			UriHash: n.UriHash,
			Data:    n.Data,
			Owner:   ownerStr,
		}
	}

	return &onfttypes.QueryNFTsResponse{
		Nfts:       nftsOwners,
		Pagination: nfts.Pagination,
	}, nil
}

func (k Keeper) NFT(ctx context.Context, r *onfttypes.QueryNFTRequest) (*onfttypes.QueryNFTResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	if len(r.ClassId) == 0 {
		return nil, nft.ErrEmptyClassID
	}
	if len(r.Id) == 0 {
		return nil, nft.ErrEmptyNFTID
	}

	n, has := k.nftKeeper.GetNFT(ctx, r.ClassId, r.Id)
	if !has {
		return nil, nft.ErrNFTNotExists.Wrapf("not found nft: class: %s, id: %s", r.ClassId, r.Id)
	}

	owner := k.nftKeeper.GetOwner(ctx, r.ClassId, r.Id)
	ownerStr, err := k.addressCodec.BytesToString(owner)
	if err != nil {
		return nil, err
	}

	resp := &onfttypes.QueryNFTResponse{
		Nft: &onfttypes.NFT{
			ClassId: n.ClassId,
			Id:      n.Id,
			Uri:     n.Uri,
			UriHash: n.UriHash,
			Owner:   ownerStr,
			Data:    n.Data,
		},
	}

	return resp, nil
}

func (k Keeper) Class(ctx context.Context, r *onfttypes.QueryClassRequest) (*onfttypes.QueryClassResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	if len(r.ClassId) == 0 {
		return nil, nft.ErrEmptyClassID
	}

	class, has := k.nftKeeper.GetClass(ctx, r.ClassId)
	if !has {
		return nil, nft.ErrClassNotExists.Wrapf("not found class: %s", r.ClassId)
	}

	owner, err := k.ClassOwners.Get(ctx, r.ClassId)
	if err != nil {
		return nil, err
	}

	ownerStr, err := k.addressCodec.BytesToString(owner)
	if err != nil {
		return nil, err
	}

	resp := &onfttypes.QueryClassResponse{
		Class: &onfttypes.Class{
			Id:          class.Id,
			Name:        class.Name,
			Symbol:      class.Symbol,
			Description: class.Description,
			Uri:         class.Uri,
			UriHash:     class.UriHash,
			Data:        class.Data,
			Owner:       ownerStr,
		},
	}

	return resp, nil
}

func (k Keeper) Classes(ctx context.Context, r *onfttypes.QueryClassesRequest) (*onfttypes.QueryClassesResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	classes, err := k.nftKeeper.Classes(ctx, &nft.QueryClassesRequest{Pagination: r.Pagination})
	if err != nil {
		return nil, err
	}

	classesOwners := make([]*onfttypes.Class, len(classes.Classes))
	for i, class := range classes.Classes {
		ownerStr := ""
		owner, err := k.ClassOwners.Get(ctx, class.Id)
		if err == nil {
			ownerStr, _ = k.addressCodec.BytesToString(owner)
		}

		classesOwners[i] = &onfttypes.Class{
			Id:          class.Id,
			Name:        class.Name,
			Symbol:      class.Symbol,
			Description: class.Description,
			Uri:         class.Uri,
			UriHash:     class.UriHash,
			Data:        class.Data,
			Owner:       ownerStr,
		}
	}

	return &onfttypes.QueryClassesResponse{
		Classes:    classesOwners,
		Pagination: classes.Pagination,
	}, nil
}
