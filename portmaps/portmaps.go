package portmaps

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/packet/packettype"
  "encoding/json"
  "sync"
)

var locksLock = sync.Mutex{}
var portLocks = make(map[int]sync.Mutex)
var encoders = make(map[int]*json.Encoder)

func AddEncoder(port int, encoder *json.Encoder) {
  locksLock.Lock()
  defer locksLock.Unlock()
  portLocks[port] = sync.Mutex{}
  encoders[port] = encoder
}
func RemoveEncoder(port int) {
  locksLock.Lock()
  defer locksLock.Unlock()
  delete(portLocks, port)
  delete(encoders, port)
}
func AllPorts() []int {
  locksLock.Lock()
  defer locksLock.Unlock()
  ports := make([]int, len(portLocks))
  i := 0
  for k, _ := range portLocks {
    //log.Debug.Printf("k:%v\n", k)
    ports[i] = k
    i += 1
  }
  return ports
}
func SendInterestPacket(port int, packet packet.InterestPacket_s) error {
  return sendGenericPacket(port, packettype.InterestPacket_c, packet)
}
func SendInterestReturnPacket(port int, packet packet.InterestReturnPacket_s) error {
  return sendGenericPacket(port, packettype.InterestReturnPacket_c, packet)
}
func SendDataPacket(port int, packet packet.DataPacket_s) error {
  return sendGenericPacket(port, packettype.DataPacket_c, packet)
}
func SendServiceProviderPacket(port int, packet packet.ServiceProviderPacket_s) error {
  return sendGenericPacket(port, packettype.ServiceProviderPacket_c, packet)
}
func sendGenericPacket(port int, packetType packettype.Base, out_packet interface{}) error {
  locksLock.Lock()
  lock := portLocks[port]
  locksLock.Unlock()
  lock.Lock()
  defer lock.Unlock()
  encoder, found := encoders[port]
  if !found {
    return errortype.PortNotRegistered
  }
  bs, err := json.Marshal(out_packet)
  if err != nil {
    return err
  }
  return encoder.Encode(packet.GenericPacket_s{
    PacketType: packetType,
    Payload:    bs,
  })
}
