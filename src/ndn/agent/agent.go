package agent

import (
	"fmt"
	"ndn"
	"net"
	"strconv"
	"sync"
)

func handleInterestConnection(conn net.Conn, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	return
}
func handleDataConnection(conn net.Conn, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	return
}

func Init(config ndn.Config, wg *sync.WaitGroup) (err error) {
	fmt.Println("agent init start")
	/* start socket service for client */
	server := config.Agent.Self
	interestLn, err := net.Listen(server.Mode, ndn.JoinHostPort(server.Host, server.InterestPort))
	if err != nil {
		fmt.Println("failed to listen on interest port", server.InterestPort, err)
		return
	}
	wg.Add(1)
	go func(ln net.Listener, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			if conn, err := ln.Accept(); err != nil {
				fmt.Println("failed to accept interest connection", err)
				return err
			} else {
				wg.Add(1)
				handleInterestConnection(conn, wg)
			}
		}
	}()
	ln, err := net.Listen(config.Agent.Mode, ":"+strconv.Itoa(config.Agent.Port))
	if err != nil {
		fmt.Printf("failed to listen on port %v, %v\n", config.Agent.Port, err)
		return
	}
	fmt.Println("fork to handle incoming connections")
	wg.Add(1)
	go func(ln net.Listener, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			fmt.Println("waiting for incoming connection")
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("failed to accect incoming connection", err)
			}
			fmt.Printf("accepted incoming connection, remote address = %v\n", conn.RemoteAddr())
			wg.Add(1)
			go handlePeerConnection(conn, wg)
		}
	}(ln, wg)
	/* contact peer / neighbour from config */
	for _, peer := range config.Agent.Peers {
		fmt.Println("connecting to peer", peer)
		interestConn, err := net.Dial(peer.Mode, ndn.JoinHostPort(peer.Host, peer.InterestPort))
		if err != nil {
			fmt.Println("failed to start interest connection to", peer, err)
			continue
		} else {
			wg.Add(1)
			go handleInterestConnection(interestConn, wg)
		}
		dataConn, err := net.Dial(peer.Mode, ndn.JoinHostPort(peer.Host, peer.DataPort))
		if err != nil {
			fmt.Println("failed to start data connection to", peer, err)
			continue
		} else {
			wg.Add(1)
			go handleDataConnection(dataConn, wg)
		}
	}
	fmt.Println("agent init finished")
	return
}
