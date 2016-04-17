package packet

import (
  "crypto/rsa"
  "ndn/packet/contentname"
  "ndn/packet/returncode"
  "time"
)

type ContentName_s struct {
  Name  string
  Type  contentname.Base
  Param interface{}
}

type InterestPacket_s struct {
  ContentName        ContentName_s
  SeqNum             int64
  AllowCache         bool
  PublisherPublicKey rsa.PublicKey
  //Selector    struct{}
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
  ExpireTime  time.Time
  AllowCache  bool
}
type ServiceProviderPacket_s struct {
  ContentName ContentName_s
  PublicKey   rsa.PublicKey
}

func test() {
}
