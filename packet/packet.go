package packet

import (
  . "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "bitbucket.org/polyu-named-data-network/ndn/packet/returncode"
  "crypto/rsa"
  "go/printer"
  "time"
)

type InterestPacket_s struct {
  ContentName        ContentName_s
  SeqNum             int64
  AllowCache         bool
  PublisherPublicKey rsa.PublicKey
  DataPort           int
  InterestReturnPort int
}

/* NAK */
type InterestReturnPacket_s struct {
  ContentName ContentName_s
  SeqNum      int64
  ReturnCode  returncode.Base
}

type DataPacket_s struct {
  ContentName        ContentName_s
  SeqNum             int64
  ExpireTime         time.Time
  AllowCache         bool
  PublisherPublicKey rsa.PublicKey
  ContentData        []byte
}
type ServiceProviderPacket_s struct {
  ContentName ContentName_s
  PublicKey   rsa.PublicKey
}

func (p DataPacket_s) New(seqNum int64) DataPacket_s {
  return DataPacket_s{
    ContentName:        p.ContentName,
    SeqNum:             seqNum,
    ExpireTime:         p.ExpireTime,
    AllowCache:         p.AllowCache,
    PublisherPublicKey: p.PublisherPublicKey,
    ContentData:        p.ContentData,
  }
}

//TODO
func (p InterestPacket_s) New(dataPort int) InterestPacket_s {
  return InterestPacket_s{
    ContentName:        p.ContentName,
    SeqNum:             p.SeqNum,
    AllowCache:         p.AllowCache,
    PublisherPublicKey: p.PublisherPublicKey,
    DataPort:           dataPort,
    InterestReturnPort: p.InterestReturnPort,
  }
}
