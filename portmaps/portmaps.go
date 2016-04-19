package portmaps

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
  "encoding/json"
  "sync"
)

var interestPacket_encoder_map = make(map[int]*json.Encoder)
var interestPacket_encoder_lock = sync.Mutex{}

var dataPacket_encoder_map = make(map[int]*json.Encoder)
var dataPacket_encoder_lock = sync.Mutex{}

func AddInterestPacketEncoder(port int, encoder *json.Encoder) {
  interestPacket_encoder_lock.Lock()
  defer interestPacket_encoder_lock.Unlock()
  interestPacket_encoder_map[port] = encoder
}
func RemoveInterestPacketEncoder(port int) {
  interestPacket_encoder_lock.Lock()
  defer interestPacket_encoder_lock.Unlock()
  delete(interestPacket_encoder_map, port)
}
func GetInterestPacketEncoder(port int) (encoder *json.Encoder, err error) {
  interestPacket_encoder_lock.Lock()
  defer interestPacket_encoder_lock.Unlock()
  if encoder, found := interestPacket_encoder_map[port]; found {
    return encoder, nil
  } else {
    return nil, errortype.PortNotRegistered
  }
}

func AddDataPacketEncoder(port int, encoder *json.Encoder) {
  dataPacket_encoder_lock.Lock()
  defer dataPacket_encoder_lock.Unlock()
  dataPacket_encoder_map[port] = encoder
}
func RemoveDataPacketEncoder(port int) {
  dataPacket_encoder_lock.Lock()
  defer dataPacket_encoder_lock.Unlock()
  delete(dataPacket_encoder_map, port)
}
func GetDataPacketEncoder(port int) (encoder *json.Encoder, err error) {
  dataPacket_encoder_lock.Lock()
  defer dataPacket_encoder_lock.Unlock()
  if encoder, found := dataPacket_encoder_map[port]; found {
    return encoder, nil
  } else {
    return nil, errortype.PortNotRegistered
  }
}
