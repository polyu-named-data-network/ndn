/* Forwarding Information Base */
package fib

import (
  "crypto/rsa"
  "ndn/packet"
  "ndn/packet/contentname"
  "sync"
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
