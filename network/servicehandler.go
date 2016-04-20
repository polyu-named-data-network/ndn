package network

import (
  "bitbucket.org/polyu-named-data-network/ndn/fib"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
)

func OnServicePacketReceived(port int, packet packet.ServiceProviderPacket_s) {
  fib.Register(port, packet)
}
