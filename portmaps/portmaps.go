package portmaps

import (
  "encoding/json"
  "sync"
)

var interestPacket_encoder_map = make(map[int]*json.Encoder)
var interestPacketEncoderLock = sync.Mutex{}

func AddInterestPacketEncoder(port int, encoder *json.Encoder) {
  interestPacketEncoderLock.Lock()
  defer interestPacketEncoderLock.Unlock()
  interestPacket_encoder_map[port] = encoder
}
func RemoveInterestEncoder(port int) {
  interestPacketEncoderLock.Lock()
  defer interestPacketEncoderLock.Unlock()
  delete(interestPacket_encoder_map, port)
}
func GetInterestEncoder(port int) (encoder *json.Encoder, found bool) {
  interestPacketEncoderLock.Lock()
  defer interestPacketEncoderLock.Unlock()
  encoder, found = interestPacket_encoder_map[port]
  return
}
