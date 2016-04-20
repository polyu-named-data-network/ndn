package packet

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
  . "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "bitbucket.org/polyu-named-data-network/ndn/packet/returncode"
  "crypto/rsa"
  "github.com/aabbcc1241/goutils/lang"
  "time"
)

type InterestPacket_s struct {
  InterestPacket_s   bool
  ContentName        ContentName_s
  SeqNum             int64
  AllowCache         bool
  PublisherPublicKey rsa.PublicKey
  DataPort           int
}

/* NAK */
type InterestReturnPacket_s struct {
  InterestReturnPacket_s bool
  ContentName            ContentName_s
  SeqNum                 int64
  ReturnCode             returncode.Base
}

type DataPacket_s struct {
  DataPacket_s       bool
  ContentName        ContentName_s
  SeqNum             int64
  ExpireTime         time.Time
  AllowCache         bool
  PublisherPublicKey rsa.PublicKey
  ContentData        []byte
}
type ServiceProviderPacket_s struct {
  ServiceProviderPacket_s bool
  ContentName             ContentName_s
  PublicKey               rsa.PublicKey
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

func (s *InterestPacket_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    if err := lang.SetField(s, k, v); err != nil {
      return err
    }
  }
  return nil
}
func (s *InterestReturnPacket_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    if err := lang.SetField(s, k, v); err != nil {
      return err
    }
  }
  return nil
}
func (s *DataPacket_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    if err := lang.SetField(s, k, v); err != nil {
      return err
    }
  }
  return nil
}
func (s *ServiceProviderPacket_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    if err := lang.SetField(s, k, v); err != nil {
      return err
    }
  }
  return nil
}
func ParsePacket(i interface{}) (interface{}, error) {
  m := i.(map[string]interface{})
  _, found := m["InterestPacket_s"]
  if found {
    packet := InterestPacket_s{}
    err := packet.FillStruct(i)
    return packet, err
  }
  _, found = m["InterestReturnPacket_s"]
  if found {
    packet := InterestReturnPacket_s{}
    err := packet.FillStruct(i)
    return packet, err
  }
  _, found = m["DataPacket_s"]
  if found {
    packet := DataPacket_s{}
    err := packet.FillStruct(i)
    return packet, err
  }
  _, found = m["ServiceProviderPacket_s"]
  if found {
    packet := ServiceProviderPacket_s{}
    err := packet.FillStruct(i)
    return packet, err
  }
  return nil, errortype.UnknownDataType
}
