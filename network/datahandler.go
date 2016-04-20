package network

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/pit"
  "github.com/aabbcc1241/goutils/log"
)

func OnDataPacketReceived(in_packet packet.DataPacket_s) (err error) {
  log.Info.Println("received data packet", in_packet)
  /*
   *    1. lookup PIT, forward to pending ports
   *    2. store in CS if allow cache
   *    3. update FIB if needed
   */

  /* 1. lookup PIT */
  err = pit.OnDataPacketReceived(in_packet)
  if err != nil {
    log.Error.Println("PIT failed to handle data packet", in_packet, err)
  }

  /* 2. store in CS if allow cache */
  log.Error.Println("not impl")
  return errortype.NotImpl
}
