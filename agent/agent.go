package agent

import (
  "fmt"
  "ndn"
  "net"
  "sync"
)

/*
  1. start interest server
  2. start data server
  3. connect to peer interest server
  4. connect to peer data server
*/
func Init(config ndn.Config, wg *sync.WaitGroup) (err error) {
  fmt.Println("agent init start")
  server := config.Agent.Self

  interestHandler := &interestHandler_s{wg}
  dataHandler := &dataHandler_s{wg}

  /* 1. start interest server */
  if interestLn, err := net.Listen(server.Mode, ndn.JoinHostPort(server.Host, server.InterestPort)); err != nil {
    ndn.ErrorLogger.Println("failed to listen on interest port", err)
  } else {
    // fork and wait for handle incoming connection
    wg.Add(1)
    go func() {
      defer wg.Done()
      fmt.Println("listening for incoming interest socket connection")
      ndn.LoopWaitHandleConnection(interestLn, interestHandler)
    }()
  }

  /* 2. start data server */
  if dataLn, err := net.Listen(server.Mode, ndn.JoinHostPort(server.Host, server.DataPort)); err != nil {
    ndn.ErrorLogger.Println("failed to listen on data port", err)
  } else {
    // fork and wait for handle incoming connection
    wg.Add(1)
    go func() {
      defer wg.Done()
      fmt.Println("listening for incoming data socket connection")
      ndn.LoopWaitHandleConnection(dataLn, dataHandler)
    }()
  }

  for _, peer := range config.Agent.Peers {
    fmt.Println("connecting to peer", peer)

    /* 3. connect to peer interest server */
    if conn, err := net.Dial(peer.Mode, ndn.JoinHostPort(peer.Host, peer.InterestPort)); err != nil {
      fmt.Printf("failed to connect to peer %v for interst (%v)\n", peer.Host, peer.InterestPort)
    } else {
      //wg.Add(1)
      //defer wg.Done()
      go interestHandler.HandleConnection(conn)
    }

    /* 4. connect to peer data server */
    if conn, err := net.Dial(peer.Mode, ndn.JoinHostPort(peer.Host, peer.DataPort)); err != nil {
      fmt.Printf("failed to connect to peer %v for interst (%v)\n", peer.Host, peer.DataPort)
    } else {
      //wg.Add(1)
      //defer wg.Done()
      go dataHandler.HandleConnection(conn)
    }
  }

  fmt.Println("agent init finished")
  return
}
