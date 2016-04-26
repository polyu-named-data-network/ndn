package packet

import (
	. "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
	"bitbucket.org/polyu-named-data-network/ndn/packet/packettype"
	"bitbucket.org/polyu-named-data-network/ndn/packet/returncode"
	"crypto/rsa"
	"encoding/json"
	"time"
)

type GenericPacket_s struct {
	PacketType packettype.Base
	Payload    []byte
}
type InterestPacket_s struct {
	ContentName        ContentName_s
	SeqNum             uint64
	AllowCache         bool
	PublisherPublicKey rsa.PublicKey
}

/* NAK */
type InterestReturnPacket_s struct {
	ContentName ContentName_s
	SeqNum      uint64
	ReturnCode  returncode.Base
}

type DataPacket_s struct {
	ContentName        ContentName_s
	SeqNum             uint64
	ExpireTime         time.Time
	AllowCache         bool
	PublisherPublicKey rsa.PublicKey
	ContentData        []byte
}
type ServiceProviderPacket_s struct {
	ContentName ContentName_s
	PublicKey   rsa.PublicKey
}

func (p DataPacket_s) New(seqNum uint64) DataPacket_s {
	return DataPacket_s{
		ContentName:        p.ContentName,
		SeqNum:             seqNum,
		ExpireTime:         p.ExpireTime,
		AllowCache:         p.AllowCache,
		PublisherPublicKey: p.PublisherPublicKey,
		ContentData:        p.ContentData,
	}
}

func (p InterestPacket_s) ToGenericPacket() (gp GenericPacket_s, err error) {
	var bs []byte
	if bs, err = json.Marshal(p); err == nil {
		gp = GenericPacket_s{PacketType: packettype.InterestPacket_c, Payload: bs}
	}
	return
}
func (p InterestReturnPacket_s) ToGenericPacket() (gp GenericPacket_s, err error) {
	var bs []byte
	if bs, err = json.Marshal(p); err == nil {
		gp = GenericPacket_s{PacketType: packettype.InterestReturnPacket_c, Payload: bs}
	}
	return
}
func (p DataPacket_s) ToGenericPacket() (gp GenericPacket_s, err error) {
	var bs []byte
	if bs, err = json.Marshal(p); err == nil {
		gp = GenericPacket_s{PacketType: packettype.DataPacket_c, Payload: bs}
	}
	return
}
func (p ServiceProviderPacket_s) ToGenericPacket() (gp GenericPacket_s, err error) {
	var bs []byte
	if bs, err = json.Marshal(p); err == nil {
		gp = GenericPacket_s{PacketType: packettype.ServiceProviderPacket_c, Payload: bs}
	}
	return
}
