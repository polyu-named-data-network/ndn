package packet

import "ndn/packet/returncode"

type ContentName_s struct {
}

type InterestPacket_s struct {
	ContentName ContentName_s
	SeqNum      int64
	AllowCache  bool
	Selector    struct{}
}

/* NAK */
type InterestReturnPacket_s struct {
	ContentName ContentName_s
	SeqNum      int64
	ReturnCode  returncode.Base
}

type DataPacket_s struct {
	ContentName ContentName_s
	SeqNum      int64
	AllowCache  bool
}

func test() {
}
