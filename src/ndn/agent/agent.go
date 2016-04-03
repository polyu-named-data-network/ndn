package agent

import (
	"fmt"
	"ndn"
	"net"
	"sync"
	"log"
)

type interestHandler_s struct {
	wg *sync.WaitGroup
}
type dataHandler_s struct {
	wg *sync.WaitGroup
}

func (interestHandler_s) HandleError(err error) {
	fmt.Println("falied to listen on incoming interest socket", err)
}
func (dataHandler_s)HandleError(err error) {
	fmt.Println("failed to listen on incoming data socket", err)
}

/* REMARK: this function is blocking */
func (p interestHandler_s)HandleConnection(conn net.Conn) {
	p.wg.Add(1)
	defer p.wg.Done()
	//TODO
}

/* REMARK: this function is blocking */
func (p dataHandler_s)HandleConnection(conn net.Conn) {
	p.wg.Add(1)
	defer p.wg.Done()
	//TODO
}

/*
  1. start interest server
  2. start data server
  3. connect to peer interest server
  4. connect to peer data server
*/
func Init(config ndn.Config, wg *sync.WaitGroup) (err error) {
	fmt.Println("agent init start")
	server := config.Agent.Self

	/* 1. start interest server */
	interestLn, err := net.Listen(server.Mode, ndn.JoinHostPort(server.Host, server.InterestPort))
	if err != nil {
		log.Fatalf("failed to start interest socker:%v\nserver=%v\nport=%v\n", err, server.Host, server.InterestPort)
	}
	// fork and wait for handle incoming connection
	interestHandler := &interestHandler_s{wg}
	wg.Add(1)
	go ndn.LoopWaitHandleConnection(interestLn, interestHandler)

	/* 2. start data server */
	dataLn, err := net.Listen(server.Mode, ndn.JoinHostPort(server.Host, server.DataPort))
	if err != nil {
		log.Fatalln("falied to listen on data port:", err)
	}
	dataHandler := &dataHandler_s{wg}
	wg.Add(1)
	go ndn.LoopWaitHandleConnection(dataLn, dataHandler)

	for _, peer := range config.Agent.Peers {
		fmt.Println("connecting to peer", peer)
		/* 3. connect to peer interest server */
		if conn, err := net.Dial(peer.Mode, ndn.JoinHostPort(peer.Host, peer.InterestPort)); err != nil {
			fmt.Printf("failed to connect to peer %v for interst (%v)\n", peer.Host, peer.InterestPort)
		}else {
			wg.Add(1)
			go interestHandler.HandleConnection(conn)
		}
		/* 4. connect to peer data server */
		if conn, err := net.Dial(peer.Mode, ndn.JoinHostPort(peer.Host, peer.DataPort)); err != nil {
			fmt.Printf("failed to connect to peer %v for interst (%v)\n", peer.Host, peer.DataPort)
		}else {
			wg.Add(1)
			go dataHandler.HandleConnection(conn)
		}
	}

	fmt.Println("agent init finished")
	return
}
