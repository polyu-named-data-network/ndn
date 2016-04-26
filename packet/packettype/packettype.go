package packettype

type Base int

const (
	InterestPacket_c = 1 + iota
	InterestReturnPacket_c
	DataPacket_c
	ServiceProviderPacket_c
)
