package packet

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
  . "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "bitbucket.org/polyu-named-data-network/ndn/packet/returncode"
  "crypto/rsa"
  "github.com/aabbcc1241/goutils/lang"
  "github.com/aabbcc1241/goutils/log"
  "math/big"
  "time"
)

type PublicKey_s struct {
  N string
  E int
}
type InterestPacket_s struct {
  InterestPacket_s   bool
  ContentName        ContentName_s
  SeqNum             uint64
  AllowCache         bool
  PublisherPublicKey PublicKey_s
  DataPort           int
}

/* NAK */
type InterestReturnPacket_s struct {
  InterestReturnPacket_s bool
  ContentName            ContentName_s
  SeqNum                 uint64
  ReturnCode             returncode.Base
}

type DataPacket_s struct {
  DataPacket_s       bool
  ContentName        ContentName_s
  SeqNum             uint64
  ExpireTime         time.Time
  AllowCache         bool
  PublisherPublicKey PublicKey_s
  ContentData        []byte
}
type ServiceProviderPacket_s struct {
  ServiceProviderPacket_s bool
  ContentName             ContentName_s
  PublicKey               PublicKey_s
}

func (p PublicKey_s) ToPublicKey() rsa.PublicKey {
  publicKey := rsa.PublicKey{}
  publicKey.N = big.NewInt(0)
  publicKey.N.UnmarshalText([]byte(p.N))
  publicKey.E = p.E
  return publicKey
}
func ToPublicKey_s(p rsa.PublicKey) (publicKey PublicKey_s, err error) {
  if bs, err2 := p.N.MarshalText(); err != nil {
    err = err2
  } else {
    publicKey = PublicKey_s{
      N: string(bs),
      E: p.E,
    }
  }
  return
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

func (s *PublicKey_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    if err := lang.SetField(s, k, v); err != nil {
      return err
    }
  }
  return nil
}
func (s *InterestPacket_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    var err error
    switch k {
    case "ContentName":
      contentName := ContentName_s{}
      err = contentName.FillStruct(v)
      if err == nil {
        err = lang.SetField(s, k, contentName)
      }
      break
    case "PublicKey":
    case "PublisherPublicKey":
      publicKey := PublicKey_s{}
      err = publicKey.FillStruct(v)
      if err == nil {
        err = lang.SetField(s, k, publicKey)
      }
      break
    default:
      err = lang.SetField(s, k, v)
    }
    if err != nil {
      return err
    }
  }
  return nil
}
func (s *InterestReturnPacket_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    var err error
    switch k {
    case "ContentName":
      contentName := ContentName_s{}
      err = contentName.FillStruct(v)
      if err == nil {
        err = lang.SetField(s, k, contentName)
      }
      break
    case "PublicKey":
    case "PublisherPublicKey":
      publicKey := PublicKey_s{}
      err = publicKey.FillStruct(v)
      if err == nil {
        err = lang.SetField(s, k, publicKey)
      }
      break
    default:
      err = lang.SetField(s, k, v)
    }
    if err != nil {
      return err
    }
  }
  return nil
}
func (s *DataPacket_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    var err error
    switch k {
    case "ContentName":
      contentName := ContentName_s{}
      err = contentName.FillStruct(v)
      if err == nil {
        err = lang.SetField(s, k, contentName)
      }
      break
    case "PublicKey":
    case "PublisherPublicKey":
      publicKey := PublicKey_s{}
      err = publicKey.FillStruct(v)
      if err == nil {
        err = lang.SetField(s, k, publicKey)
      }
      break
    default:
      err = lang.SetField(s, k, v)
    }
    if err != nil {
      return err
    }
  }
  return nil
}
func (s *ServiceProviderPacket_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    var err error
    switch k {
    case "ContentName":
      contentName := ContentName_s{}
      err = contentName.FillStruct(v)
      if err == nil {
        err = lang.SetField(s, k, contentName)
      }
      break
    case "PublicKey":
    case "PublisherPublicKey":
      publicKey := PublicKey_s{}
      err = publicKey.FillStruct(v)
      if err == nil {
        err = lang.SetField(s, k, publicKey)
      }
      break
    default:
      err = lang.SetField(s, k, v)
    }
    if err != nil {
      return err
    }
  }
  log.Debug.Println("sevice packet parsed", s)
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
