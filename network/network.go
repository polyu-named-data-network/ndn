package network

import (
	"bitbucket.org/polyu-named-data-network/ndn/config"
	"bitbucket.org/polyu-named-data-network/ndn/errortype"
	"bitbucket.org/polyu-named-data-network/ndn/fib"
	"bitbucket.org/polyu-named-data-network/ndn/packet"
	"bitbucket.org/polyu-named-data-network/ndn/packet/packettype"
	"bitbucket.org/polyu-named-data-network/ndn/pit"
	"bitbucket.org/polyu-named-data-network/ndn/portmaps"
	"bitbucket.org/polyu-named-data-network/ndn/utils"
	"encoding/json"
	"github.com/beenotung/goutils/log"
	"io"
	"net"
	"strconv"
	"sync"
)

func handleConnection(conn net.Conn, wg *sync.WaitGroup) (err error) {
	log.Info.Println("connected with", conn.RemoteAddr().String())
	_, port_string, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Error.Println("failed to split port from remote address", conn.RemoteAddr(), err)
		return err
	}
	port, err := strconv.Atoi(port_string)
	if err != nil {
		log.Error.Println("failed to parse port", port_string, err)
		return err
	}
	decoder := json.NewDecoder(conn)
	portmaps.AddEncoder(port, json.NewEncoder(conn))
	wg.Add(1)
	go func() {
		defer func() {
			log.Debug.Println("close connection", conn.RemoteAddr().String())
			fib.UnRegister(port)
			pit.UnRegister(port)
			portmaps.RemoveEncoder(port)
			conn.Close()
			wg.Done()
		}()
		for err == nil {
			log.Debug.Println("waiting for incoming packet")
			var in_packet packet.GenericPacket_s
			err = decoder.Decode(&in_packet)
			if err != nil {
				if err != io.EOF {
					log.Error.Println("failed to parse incoming message into json", in_packet, err)
					//continue
				}
				break
			}
			if err := OnGenericPacketReceived(port, in_packet); err != nil {
				log.Error.Println("failed to handle generic packet", err)
			}
		}
	}()
	return
}

func OnGenericPacketReceived(in_port int, in_packet packet.GenericPacket_s) (err error) {
	//log.Info.Println("received packet from port",in_port)
	switch in_packet.PacketType {
	case packettype.InterestPacket_c:
		var packet packet.InterestPacket_s
		err = json.Unmarshal(in_packet.Payload, &packet)
		if err != nil {
			return
		}
		return OnInterestPacketReceived(in_port, packet)
	case packettype.InterestReturnPacket_c:
		var packet packet.InterestReturnPacket_s
		err = json.Unmarshal(in_packet.Payload, &packet)
		if err != nil {
			return
		}
		log.Error.Println("not impl: received interest return packet from port", in_port)
		return errortype.NotImpl
	case packettype.DataPacket_c:
		var packet packet.DataPacket_s
		err = json.Unmarshal(in_packet.Payload, &packet)
		if err != nil {
			return
		}
		return OnDataPacketReceived(in_port, packet)
	case packettype.ServiceProviderPacket_c:
		var packet packet.ServiceProviderPacket_s
		err = json.Unmarshal(in_packet.Payload, &packet)
		if err != nil {
			return
		}
		return OnServicePacketReceived(in_port, packet)
	default:
		log.Error.Printf("packettype (%v) not suppored: %v\n", in_packet.PacketType, in_packet)
		return errortype.NotImpl
	}
}

/*
  1. start socket server
  2. connect to peer (socket server)
*/
func Init(config config.Config, wg *sync.WaitGroup) (err error) {
	log.Info.Println("network init start")
	server := config.Self

	/* 1. start socket server */
	var ln net.Listener
	ln, err = net.Listen(server.Mode, utils.JoinHostPort(server.Host, server.Port))
	if err == nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			for err == nil {
				var conn net.Conn
				conn, err = ln.Accept()
				if err != nil {
					log.Error.Println(err)
				} else {
					err = handleConnection(conn, wg)
				}
			}
		}()
	} else {
		log.Error.Println("failed to bind server socket", err)
	}

	/* 2. connect to peer (socket server) */
	for _, peer := range config.Peers {
		log.Info.Println("connecting to peer", peer)
		var conn net.Conn
		if conn, err = net.Dial(peer.Mode, utils.JoinHostPort(peer.Host, peer.Port)); err != nil {
			log.Error.Printf("failed to connect to peer, host:%v, port:%v, %v", peer.Host, peer.Port, err)
		} else {
			handleConnection(conn, wg)
		}
	}

	log.Info.Println("network init finished")
	return
}
