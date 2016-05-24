package network

import (
	"bitbucket.org/polyu-named-data-network/ndn/cs"
	"bitbucket.org/polyu-named-data-network/ndn/fib"
	"bitbucket.org/polyu-named-data-network/ndn/packet"
	"bitbucket.org/polyu-named-data-network/ndn/packet/returncode"
	"bitbucket.org/polyu-named-data-network/ndn/pit"
	"bitbucket.org/polyu-named-data-network/ndn/portmaps"
	"encoding/json"
	"github.com/beenotung/goutils/log"
	"io"
	"net"
	"strconv"
	"sync"
)

type interestHandler_s struct {
	wg *sync.WaitGroup
}

func (interestHandler_s) HandleError(err error) {
	log.Error.Println("failed to listen on incoming interest socket", err)
}

/* REMARK: this function is blocking */
func (p interestHandler_s) HandleConnection(conn net.Conn) {
	p.wg.Add(1)
	log.Info.Println("client connected to interest service", conn.RemoteAddr().Network(), conn.RemoteAddr().String())

	if _, in_port_string, err := net.SplitHostPort(conn.RemoteAddr().String()); err != nil {
		log.Error.Println("failed to parse port from connection", conn)
		conn.Close()
		p.wg.Done()
	} else {
		if in_port, err := strconv.Atoi(in_port_string); err != nil {
			log.Error.Println("failed to parse port", in_port_string)
			conn.Close()
			p.wg.Done()
		} else {
			go func() {
				defer conn.Close()
				defer p.wg.Done()
				in := json.NewDecoder(conn)
				var in_packet packet.InterestPacket_s
				var err error
				for err == nil {
					log.Info.Println("wait for incoming interest packet")
					err = in.Decode(&in_packet)
					if err != nil {
						if err != io.EOF {
							log.Error.Println("failed to decode incoming interest packet", err)
						}
					} else {
						OnInterestPacketReceived(in_port, in_packet)
					}
				}
			}()
		}
	}
	//TODO
}

func OnInterestPacketReceived(in_port int, in_packet packet.InterestPacket_s) (err error) {
	log.Info.Println("received interest port", in_port, "packet", in_packet)
	/*  find data, response if found, otherwise do forwarding
	 *    1. lookup CS
	 *    2. lookup PIT
	 *    3. lookup FIB
	 *    4. calculate forwarding port according to algorithm for unknown cases
	 */

	/* 1. lookup CS */
	log.Debug.Println("checking CS")
	if cachedPacket, csFound := cs.Get(in_packet.ContentName); csFound {
		log.Debug.Println("found in CS")
		portmaps.SendDataPacket(in_port, cachedPacket.New(in_packet.SeqNum))
		return
	}

	log.Debug.Println("not found in CS")

	/* 2. lookup PIT */
	if in_packet.AllowCache {
		log.Debug.Println("checking PIT")
		if pitFound := pit.AddToListIfExist(in_port, in_packet); pitFound {
			log.Debug.Println("found in PIT")
			return
		}
		log.Debug.Println("not found in PIT")
	} else {
		log.Debug.Println("not allow cache, skipped PIT checking")
	}

	/* 3. lookup FIB */
	log.Debug.Println("checking FIB")
	if out_port, fibFound := fib.Lookup(in_packet.ContentName, in_packet.PublisherPublicKey); fibFound {
		log.Debug.Println("found in FIB, port:", out_port)
		if err := fib.Forward(out_port, in_packet); err != nil {
			log.Error.Println("failed to forward on port", out_port, err)
		} else {
			pit.Register(in_port, in_packet)
		}
	} else {
		log.Debug.Println("not found in FIB")
		//TODO replace by the strategy in excel
		sentCount := 0
		ports := portmaps.AllPorts()
		//log.Debug.Println("ports", ports)
		for _, out_port := range ports {
			//log.Debug.Println("out_port", out_port)
			if out_port == in_port {
				continue
			}
			if err := fib.Forward(out_port, in_packet); err != nil {
				log.Error.Println("failed to forward on port", out_port, err)
			} else {
				sentCount += 1
			}
		}
		if sentCount > 0 {
			log.Debug.Printf("forwarded interest to %v peer(s)\n", sentCount)
			//portmaps.SendInterestPacket(in_port, in_packet)
			pit.Register(in_port, in_packet)
		} else {
			log.Debug.Println("interest not resolved, no available peer, sending NAK")
			portmaps.SendInterestReturnPacket(in_port, packet.InterestReturnPacket_s{
				ContentName: in_packet.ContentName,
				SeqNum:      in_packet.SeqNum,
				ReturnCode:  returncode.NoRoute,
			})
		}
		//log.Error.Println("not impl")
		//err = errortype.NotImpl
	}
	return
}
