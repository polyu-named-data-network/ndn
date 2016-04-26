package portmaps

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
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
  for i := range portLocks {
    ports = append(ports, i)
  }
  return ports
}
func Encode(port int, packet interface{}) error {
  locksLock.Lock()
  lock := portLocks[port]
  locksLock.Unlock()
  lock.Lock()
  defer lock.Unlock()
  encoder, found := encoders[port]
  if !found {
    return errortype.PortNotRegistered
  }
  encoder.Encode(packet)
  return nil
}
