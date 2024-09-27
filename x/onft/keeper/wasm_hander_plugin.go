package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"

	"cosmossdk.io/x/nft"
	"github.com/ODIN-PROTOCOL/odin-core/x/onft/types"
	wasmkeeper "github.com/ODIN-PROTOCOL/wasmd/x/wasm/keeper"
	wasmtypes "github.com/ODIN-PROTOCOL/wasmd/x/wasm/types"
	wasmvmtypes "github.com/ODIN-PROTOCOL/wasmvm/v2/types"
	abci "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ wasmkeeper.Messenger = MessageHandlerFunc(nil)

// MessageHandlerFunc is a helper to construct a function based message handler.
type MessageHandlerFunc func(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, msgResponses [][]*codectypes.Any, err error)

// DispatchMsg delegates dispatching of provided message into the MessageHandlerFunc.
func (m MessageHandlerFunc) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, msgResponses [][]*codectypes.Any, err error) {
	return m(ctx, contractAddr, contractIBCPortID, msg)
}

func NewMintNFTMessageHandler(onftKeeper Keeper) MessageHandlerFunc {
	return func(ctx sdk.Context, contractAddr sdk.AccAddress, _ string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, msgResponses [][]*codectypes.Any, err error) {
		if len(msg.Any.Value) > 0 && msg.Any.TypeURL == "onft/mint" {
			var msgNFT types.MsgMintNFT
			err := json.Unmarshal(msg.Any.Value, &msgNFT)
			if err != nil {
				return nil, nil, nil, err
			}

			if !onftKeeper.nftKeeper.HasClass(ctx, msgNFT.ClassId) {
				return nil, nil, nil, types.ErrClassNotFound
			}

			owner, err := onftKeeper.ClassOwners.Get(ctx, msgNFT.ClassId)
			if err != nil {
				return nil, nil, nil, err
			}

			if !bytes.Equal(contractAddr, owner) {
				return nil, nil, nil, types.ErrSenderNotAuthorized
			}

			nftID := fmt.Sprintf("%d", onftKeeper.nftKeeper.GetTotalSupply(ctx, msgNFT.ClassId))

			receiver, err := sdk.AccAddressFromBech32(msgNFT.Receiver)
			if err != nil {
				return nil, nil, nil, err
			}

			err = onftKeeper.nftKeeper.Mint(ctx, nft.NFT{
				ClassId: msgNFT.ClassId,
				Id:      nftID,
				Uri:     msgNFT.Uri,
				UriHash: msgNFT.UriHash,
				Data:    nil,
			}, receiver)
			if err != nil {
				return nil, nil, nil, err
			}

			ctx.Logger().
				With("module", fmt.Sprintf("x/%s", wasmtypes.ModuleName)).
				Info("Minted NFT", "id", nftID)

			events = []sdk.Event{
				{
					Type: "nft",
					Attributes: []abci.EventAttribute{
						{
							Key:   "nft_id",
							Value: nftID,
						},
					},
				},
			}

			return events, nil, nil, nil
		}

		return nil, nil, nil, wasmtypes.ErrUnknownMsg
	}
}
