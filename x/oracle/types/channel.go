package types

func NewIBCChannel(portId, channelId string) IBCChannel {
	return IBCChannel{
		PortId:    portId,
		ChannelId: channelId,
	}
}

func NewIBCSource(portId, channelId string) IBCSource {
	return IBCSource{
		SourcePort:    portId,
		SourceChannel: channelId,
	}
}
