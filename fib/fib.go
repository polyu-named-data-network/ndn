/* Forwarding Information Base */
package fib

import (
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "bitbucket.org/polyu-named-data-network/ndn/portmaps"
  "bitbucket.org/polyu-named-data-network/ndn/utils"
  "crypto/rsa"
  "encoding/json"
  "github.com/aabbcc1241/goutils/log"
  "net"
  "strconv"
  "sync"
)

type publicKeyPortsMap_t map[rsa.PublicKey][]int

var lock = sync.Mutex{}
var exactMatchTable = make(map[string]publicKeyPortsMap_t)

func UnRegister(conn net.Conn) {
  _, port_string, _ := net.SplitHostPort(conn.RemoteAddr().String())
  port, _ := strconv.Atoi(port_string)
  /* delete from name map */
  for name, publicKeyPortsMap := range exactMatchTable {
    for publicKey, ports := range publicKeyPortsMap {
      for k, v := range ports {
        if v == port {
          ports = append(ports[:k], ports[k+1:]...)
          if len(ports) == 0 {
            delete(publicKeyPortsMap, publicKey)
            if len(publicKeyPortsMap) == 0 {
              delete(exactMatchTable, name)
            }
          }
        }
      }
    }
  }
}
func Register(port int, packet packet.ServiceProviderPacket_s) {
  switch packet.ContentName.Type {
  case contentname.ExactMatch:
    var publicKeyPortsMap publicKeyPortsMap_t
    var found bool
    lock.Lock()
    defer lock.Unlock()
    publicKeyPortsMap, found = exactMatchTable[packet.ContentName.Name]
    if !found {
      publicKeyPortsMap = make(publicKeyPortsMap_t)
      exactMatchTable[packet.ContentName.Name] = publicKeyPortsMap
    }
    ports := publicKeyPortsMap[packet.PublicKey]
    ports = append(ports, port)
    publicKeyPortsMap[packet.PublicKey] = ports
    return
  default:
    log.Error.Println("not impl")
    return
  }
}
func Lookup(contentName contentname.ContentName_s, publicKey rsa.PublicKey) (port int, found bool) {
  lock.Lock()
  defer lock.Unlock()
  switch contentName.Type {
  case contentname.ExactMatch:
    var publicKeyPortsMap publicKeyPortsMap_t
    publicKeyPortsMap, found = exactMatchTable[contentName.Name]
    if !found {
      log.Debug.Println("the name is not found:", contentName.Name)
      return
    }
    /*  check PublicKey
     *    if is defined, lookup by PublicKey
     *    if not defined, lookup by any //TODO implement a rating algorithm
     */
    ports := make([]int, 0)
    if publicKey == utils.ZeroKey {
      for _, v := range publicKeyPortsMap {
        ports = append(ports, v...)
      }
    } else {
      ports, found = publicKeyPortsMap[publicKey]
    }
    if !found {
      log.Debug.Println("the public key is not found:", publicKey)
      return
    }
    found = len(ports) > 0
    if !found {
      log.Debug.Println("the port is not found")
      return
    }
    //TODO implement a rating algorithm according to loading (last/avg responding time), rather than just pick the first one
    //  similar to the equation for round trip time { (1-alpha) * oldVal + alpha * newVal }
    port = ports[0]
    break
  case contentname.LongestMatch:
    break
  case contentname.FuzzyMatch:
    break
  case contentname.Custom:
  default:

  }
  return
}
func Forward(port int, packet packet.InterestPacket_s) (err error) {
  var encoder *json.Encoder
  if encoder, err = portmaps.GetInterestPacketEncoder(port); err == nil {
    encoder.Encode(packet)
  }
  return
}
