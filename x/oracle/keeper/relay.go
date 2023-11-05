package oraclekeeper

import (
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data types.OracleRequestPacketData,
	relayer sdk.AccAddress,
) (types.RequestID, error) {
	if err := data.ValidateBasic(); err != nil {
		return 0, err
	}
	ibcSource := types.NewIBCSource(packet.DestinationPort, packet.DestinationChannel)

	return k.PrepareRequest(ctx, &data, relayer, &ibcSource)
}
