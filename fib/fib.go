/* Forwarding Information Base */
package fib

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "crypto/rsa"
  "encoding/json"
  "github.com/aabbcc1241/goutils/log"
  "net"
  "strconv"
  "sync"
)

type publicKeyPortsMap_t map[rsa.PublicKey][]int

var lock = sync.Mutex{}
var portEncoderMap = make(map[int]json.Encoder)
var exactMatchTable = make(map[string]publicKeyPortsMap_t)

func UnRegister(conn net.Conn) {
  _, port_string, _ := net.SplitHostPort(conn.RemoteAddr().String())
  port, _ := strconv.Atoi(port_string)
  /* delete from port map */
  delete(portEncoderMap, port)
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
func Register(contentName packet.ContentName_s, publicKey rsa.PublicKey, conn net.Conn) (err error) {
  if _, port, err := net.SplitHostPort(conn.RemoteAddr().String()); err != nil {
    log.Error.Println("failed to parse port from remote address", err)
    return err
  } else {
    port, err := strconv.Atoi(port)
    if err != nil {
      log.Error.Println("failed to parse port from string", err)
      return err
    } else {
      lock.Lock()
      defer lock.Unlock()

      /* save connection */
      //TODO detect old encoder, remove it, discard it from PIT
      log.Info.Println("FIB registered port:", port, "contentname:", contentName)
      portEncoderMap[port] = *json.NewEncoder(conn)

      /* save name */
      switch contentName.Type {
      case contentname.ExactMatch:

        publicKeyPortsMap, found := exactMatchTable[contentName.Name]
        if !found {
          publicKeyPortsMap = make(map[rsa.PublicKey][]int)
          exactMatchTable[contentName.Name] = publicKeyPortsMap
        }
        ports, found := publicKeyPortsMap[publicKey]
        if !found {
          publicKeyPortsMap[publicKey] = []int{port}
        } else {
          publicKeyPortsMap[publicKey] = append(ports, port)
        }
        break
      case contentname.LongestMatch:
        break
      case contentname.FuzzyMatch:
        break
      case contentname.Custom:
      default:

      }
      return err
    }
  }
}

func Lookup(contentName packet.ContentName_s, publicKey rsa.PublicKey) (port int, found bool) {
  lock.Lock()
  defer lock.Unlock()
  zeroKey := rsa.PublicKey{}
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
    if publicKey == zeroKey {
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
    //TODO implement a rating algorithm according to loading (last/avg responding time)
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
  if encoder, found := portEncoderMap[port]; found {
    //callback.Apply(packet)
    log.Debug.Println("found encoder on port", port)
    encoder.Encode(packet)
    return
  } else {
    return errortype.PortNotRegistered
  }
}
