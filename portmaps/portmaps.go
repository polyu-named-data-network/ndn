package portmaps

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
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

var interestReturn_encoder_map = make(map[int]*json.Encoder)
var interestReturn_lock = sync.Mutex{}
var interestReturn_lock_map = make(map[int]*sync.Mutex)

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
  //TODO risky to dead lock
  if lock, found := interest_lock_map[port]; found {
    lock.Lock()
  }
  delete(interest_encoder_map, port)
  log.Debug.Println("before remove", interest_lock_map)
  delete(interest_lock_map, port)
  log.Debug.Println("after remove", interest_lock_map)
}
func SendInterestPacket(port int, packet packet.InterestPacket_s) error {
  if lock, found := interest_lock_map[port]; found {
    lock.Lock()
    defer lock.Unlock()
    if encoder, found := interest_encoder_map[port]; found {
      return encoder.Encode(packet)
    } else {
      log.Error.Println("interest encoder not found on port", port)
      return errortype.PortNotRegistered
    }
  } else {
    log.Error.Println("interest lock not found on port", port)
    return errortype.PortNotRegistered
  }
}

/* return false if no peer available to forward */
func BroadcastInterestPacket(excludePort int, packet packet.InterestPacket_s) bool {
  wg := sync.WaitGroup{}
  log.Debug.Println("broadcase interest")
  log.Debug.Println("interest_lock_map", interest_lock_map)
  sentCount := 0
  for port := range interest_lock_map {
    log.Debug.Printf("excludePort:%v, port:%v\n", excludePort, port)
    if port == excludePort {
      continue
    }
    wg.Add(1)
    go func(port int) {
      defer wg.Done()
      log.Debug.Println("forwarding interest packet to port", port)
      if err := SendInterestPacket(port, packet); err != nil {
        log.Error.Println(err)
      } else {
        sentCount += 1
      }
    }(port)
  }
  wg.Wait()
  return sentCount > 0
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
  //TODO risky to dead lock
  if lock, found := data_lock_map[port]; found {
    lock.Lock()
  }
  delete(data_encoder_map, port)
  delete(data_lock_map, port)
}
func SendDataPacket(port int, packet packet.DataPacket_s) error {
  data_lock_map[port].Lock()
  defer data_lock_map[port].Unlock()
  return data_encoder_map[port].Encode(packet)
}

func AddInterestReturnEncoder(port int, encoder *json.Encoder) {
  interestReturn_lock.Lock()
  defer interestReturn_lock.Unlock()
  interestReturn_encoder_map[port] = encoder
  interestReturn_lock_map[port] = &sync.Mutex{}
}
func RemoveInterestReturnEncoder(port int) {
  interestReturn_lock.Lock()
  defer interestReturn_lock.Unlock()
  if lock, found := interestReturn_lock_map[port]; found {
    lock.Lock()
  }
  delete(interestReturn_encoder_map, port)
  delete(interestReturn_lock_map, port)
}
func SendInterestReturn(port int, packet packet.InterestReturnPacket_s) error {
  interestReturn_lock_map[port].Lock()
  defer interestReturn_lock_map[port].Unlock()
  return interestReturn_encoder_map[port].Encode(packet)
}
