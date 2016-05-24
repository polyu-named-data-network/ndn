/* Content Store */
package cs

import (
	"bitbucket.org/polyu-named-data-network/ndn/packet"
	"bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
	"github.com/beenotung/goutils/log"
	"sync"
)

var exactMatch_map = make(map[string]packet.DataPacket_s)
var exactMatch_lock = sync.Mutex{}

/* TODO load store limit from config */

func Store(packet packet.DataPacket_s) {
	contentName := packet.ContentName
	log.Debug.Printf("store %v into CS\n", contentName)
	switch contentName.ContentType {
	case contentname.ExactMatch:
		exactMatch_lock.Lock()
		defer exactMatch_lock.Unlock()
		exactMatch_map[contentName.Name] = packet
		break
	default:
		log.Error.Println("not impl")
	}
}
func Get(contentName contentname.ContentName_s) (packet.DataPacket_s, bool) {
	switch contentName.ContentType {
	case contentname.ExactMatch:
		p, b := exactMatch_map[contentName.Name]
		return p, b
	default:
		log.Error.Println("not impl")
		return packet.DataPacket_s{}, false
	}
}
