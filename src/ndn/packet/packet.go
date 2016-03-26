package packet

type Packet interface {
	GetType() string
}

type InterestPacket struct {
	Packet
	Type string
}

type DataPacket struct {
	Packet
	Type string
}

func (p *InterestPacket) GetType() string {
	return p.Type
}

func (p *DataPacket) GetType() string {
	return p.Type
}

func Create() Packet {
	p := Packet{Type: "???"}
	return p
}
func IsInterestPacket(packet Packet) bool {
	return packet.GetType().Eq
}
