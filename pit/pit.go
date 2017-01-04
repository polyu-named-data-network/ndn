/* Pending Interest Table */
package pit

import (
	"github.com/polyu-named-data-network/ndn/errortype"
	"github.com/polyu-named-data-network/ndn/packet"
	"github.com/polyu-named-data-network/ndn/packet/contentname"
	"github.com/polyu-named-data-network/ndn/portmaps"
	"github.com/polyu-named-data-network/ndn/utils"
	"crypto/rsa"
	"github.com/beenotung/goutils/log"
	"sync"
)

type pending_interest_s struct {
	SeqNum             uint64
	AllowCache         bool
	PublisherPublicKey rsa.PublicKey
	InterestPort       int
}

var exactMatchTable = make(map[string][]pending_interest_s)
var exactMatchTable_lock = sync.Mutex{}

func Register(port int, packet packet.InterestPacket_s) {
	contentName := packet.ContentName
	log.Debug.Println("register interest packet, port:", port, "contentName:", contentName)
	switch contentName.ContentType {
	case contentname.ExactMatch:
		exactMatchTable_lock.Lock()
		defer exactMatchTable_lock.Unlock()
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

/* used when peer disconnect */
func UnRegister(port int) {
	log.Info.Println("unregister from PIT on port", port)
	exactMatchTable_lock.Lock()
	defer exactMatchTable_lock.Unlock()
	for name, ps := range exactMatchTable {
		for i := len(ps) - 1; i > 0; i-- {
			if ps[i].InterestPort == port {
				/* delete from matched record */
				ps = append(ps[:i], ps[i+1:]...)
			}
		}
		if len(ps) == 0 {
			delete(exactMatchTable, name)
		}
	}
}
func AddToListIfExist(in_port int, in_packet packet.InterestPacket_s) bool {
	contentName := in_packet.ContentName
	switch contentName.ContentType {
	case contentname.ExactMatch:
		exactMatchTable_lock.Lock()
		defer exactMatchTable_lock.Unlock()
		ps, found := exactMatchTable[contentName.Name]
		if !found {
			return false
		}
		p := pending_interest_s{
			SeqNum:             in_packet.SeqNum,
			AllowCache:         in_packet.AllowCache,
			PublisherPublicKey: in_packet.PublisherPublicKey,
			InterestPort:       in_port,
		}
		exactMatchTable[contentName.Name] = append(ps, p)
		return true
	default:
		log.Error.Println("not impl")
		return false
	}
}
func HandleDataPacket(in_packet packet.DataPacket_s) (err error) {
	switch in_packet.ContentName.ContentType {
	case contentname.ExactMatch:
		exactMatchTable_lock.Lock()
		defer exactMatchTable_lock.Unlock()
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
		if len(pendingInterests) == 0 {
			delete(exactMatchTable, in_packet.ContentName.Name)
		} else {
			exactMatchTable[in_packet.ContentName.Name] = pendingInterests
		}
		return
	default:
		return errortype.ContentTypeNotSupported
	}
}
