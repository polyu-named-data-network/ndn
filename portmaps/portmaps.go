package portmaps

import (
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "encoding/json"
  "github.com/aabbcc1241/goutils/log"
  "sync"
)

var interest_encoder_map = make(map[int]*json.Encoder)
var interest_lock = sync.Mutex{}
var interest_lock_map = make(map[int]*sync.Mutex)

var data_encoder_map = make(map[int]*json.Encoder)
var data_lock = sync.Mutex{}
var data_lock_map = make(map[int]*sync.Mutex)

func AddInterestPacketEncoder(port int, encoder *json.Encoder) {
  log.Debug.Println("add interest packet encoder on port", port)
  interest_lock.Lock()
  defer interest_lock.Unlock()
  interest_encoder_map[port] = encoder
  interest_lock_map[port] = &sync.Mutex{}
}
func RemoveInterestPacketEncoder(port int) {
  log.Debug.Println("remove interest packet encoder on port", port)
  interest_lock.Lock()
  defer interest_lock.Unlock()
  delete(interest_encoder_map, port)
  delete(interest_lock_map, port)
}
func SendInterestPacket(port int, packet packet.InterestPacket_s) error {
  interest_lock_map[port].Lock()
  defer interest_lock_map[port].Unlock()
  return interest_encoder_map[port].Encode(packet)
}
func BroadcastInterestPacket(excludePort int, packet packet.InterestPacket_s) {
  wg := sync.WaitGroup{}
  log.Debug.Println("broadcase interest")
  //log.Debug.Println("interest_lock_map", interest_lock_map)
  for port := range interest_lock_map {
    //log.Debug.Printf("excludePort:%v, port:%v\n", excludePort, port)
    if port == excludePort {
      continue
    }
    wg.Add(1)
    go func(port int) {
      defer wg.Done()
      log.Debug.Println("forwarding interest packet to port", port)
      if err := SendInterestPacket(port, packet); err != nil {
        log.Error.Println(err)
      }
    }(port)
  }
  wg.Wait()
}

func AddDataPacketEncoder(port int, encoder *json.Encoder) {
  data_lock.Lock()
  defer data_lock.Unlock()
  data_encoder_map[port] = encoder
  data_lock_map[port] = &sync.Mutex{}
}
func RemoveDataPacketEncoder(port int) {
  data_lock.Lock()
  defer data_lock.Unlock()
  delete(data_encoder_map, port)
}
func SendDataPacket(port int, packet packet.DataPacket_s) error {
  data_lock_map[port].Lock()
  defer data_lock_map[port].Unlock()
  return data_encoder_map[port].Encode(packet)
}
