/* Forwarding Information Base */
package fib

import (
  "crypto/rsa"
  "sync"
  "bitbucket.org/polyu-named-data-network/ndn/packet/contentname"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
)

var lock = sync.Mutex{}

type publicKeyPortsMap map[rsa.PublicKey][]int

var exactMatchTable = make(map[string]map[rsa.PublicKey][]int)

func Register(contentName packet.ContentName_s, publicKey rsa.PublicKey, port int) {
  lock.Lock()
  defer lock.Unlock()
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
}

func Lookup(contentName packet.ContentName_s, publicKey rsa.PublicKey) (port int, found bool) {
  lock.Lock()
  defer lock.Unlock()
  switch contentName.Type {
  case contentname.ExactMatch:
    publicKeyPortsMap, found := exactMatchTable[contentName.Name]
    if !found {
      return port, found
    }
    /*  check PublicKey
     *    if is defined, lookup by PublicKey
     *    if not defined, lookup by any //TODO implement a rating algorithm
     */
    ports, found := publicKeyPortsMap[publicKey]
    if !found {
      return port, found
    }
    found = len(ports) > 0
    if !found {
      return port, found
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
