/* Pending Interest Table */
package pit

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "bitbucket.org/polyu-named-data-network/ndn/portmaps"
  "bitbucket.org/polyu-named-data-network/ndn/utils"
  "crypto/rsa"
  "github.com/aabbcc1241/goutils/log"
)

type pending_interest_s struct {
  SeqNum             uint64
  AllowCache         bool
  PublisherPublicKey rsa.PublicKey
  InterestPort       int
}

var exactMatchTable = make(map[string][]pending_interest_s)

func Register(port int, packet packet.InterestPacket_s) {
  contentName := packet.ContentName
  log.Debug.Println("register interest packet, port:", port, "contentName:", contentName)
  switch contentName.ContentType {
  case contentname.ExactMatch:
    pendingInterests := exactMatchTable[contentName.Name]
    pendingInterests = append(pendingInterests, pending_interest_s{
      SeqNum:             packet.SeqNum,
      AllowCache:         packet.AllowCache,
      PublisherPublicKey: packet.PublisherPublicKey,
      InterestPort:       port,
    })
    log.Debug.Println("new list", pendingInterests)
    exactMatchTable[contentName.Name] = pendingInterests
    break
  default:

  }
}
func HandleDataPacket(in_packet packet.DataPacket_s) (err error) {
  switch in_packet.ContentName.ContentType {
  case contentname.ExactMatch:
    pendingInterests := exactMatchTable[in_packet.ContentName.Name]
    for i := len(pendingInterests) - 1; i >= 0; i-- {
      current := pendingInterests[i]
      /* keyMatched && seqMatched */
      if (current.PublisherPublicKey == utils.ZeroKey || current.PublisherPublicKey == in_packet.PublisherPublicKey) && (current.AllowCache || current.SeqNum == in_packet.SeqNum) {
        /* delete from PIT */
        pendingInterests = append(pendingInterests[:i], pendingInterests[i+1:]...)
        /* do forward */
        err = portmaps.SendDataPacket(current.InterestPort, in_packet.New(current.SeqNum))
        if err != nil {
          log.Error.Printf("failed to send data packet to port %v in PIT, %v", current.InterestPort, err)
        }
      }
    }
    return
  default:
    return errortype.ContentTypeNotSupported
  }
}
