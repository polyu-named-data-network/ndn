/* Pending Interest Table */
package pit

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "crypto/rsa"
  "github.com/aabbcc1241/goutils/log"
)

type name_interest_map_t map[string]packet.InterestPacket_s
type pending_interest_s struct {
  SeqNum             int64
  PublisherPublicKey rsa.PublicKey
}

//type seqnum_
var exactMatchTable = make(map[string][]int)

func Register(port int, packet packet.InterestPacket_s) {
  contentName := packet.ContentName
  log.Debug.Println("register interest packet, port:", port, "contentName:", contentName)
  switch contentName.Type {
  case contentname.ExactMatch:
    ports, found := exactMatchTable[contentName.Name]
    if found {
      ports = append(ports, port)
      log.Debug.Println("added to port list", ports)
    } else {
      ports = []int{port}
      log.Debug.Println("created new port list")
    }
    exactMatchTable[contentName.Name] = ports
    break
  default:

  }
}
func GetPendingPorts(contentName contentname.ContentName_s) (ports []int, err error) {
  switch contentName.Type {
  case contentname.ExactMatch:
    var found bool
    ports, found = exactMatchTable[contentName.Name]
    if !found {
      err = errortype.ContentNameNotFound
    }
    exactMatchTable[contentName.Name] = ports
    return
  default:
    err = errortype.ContentTypeNotSupported
    return
  }
}
func Forward(ports []int, packet packet.DataPacket_s) {

}
