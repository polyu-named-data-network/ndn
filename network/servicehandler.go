package network

import (
	"github.com/polyu-named-data-network/ndn/fib"
	"github.com/polyu-named-data-network/ndn/packet"
	"github.com/beenotung/goutils/log"
)

func OnServicePacketReceived(port int, packet packet.ServiceProviderPacket_s) error {
	log.Info.Println("received service provider packet form port", port)
	return fib.Register(port, packet)
}
