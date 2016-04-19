package agent

import (
  "bitbucket.org/polyu-named-data-network/ndn/config"
  "bitbucket.org/polyu-named-data-network/ndn/utils"
  "github.com/aabbcc1241/goutils/log"
  "net"
  "sync"
)

/*
  1. start interest server
  2. start data server
  3. connect to peer interest server
  4. connect to peer data server
*/
func Init(config config.Config, wg *sync.WaitGroup) (err error) {
  log.Info.Println("agent init start")
  server := config.Agent.Self

  var interestHandler, dataHandler utils.ConnectionHandler
  interestHandler = interestHandler_s{wg}
  dataHandler = dataHandler_s{wg}

  /* 1. start interest server */
  if interestLn, err := net.Listen(server.Mode, utils.JoinHostPort(server.Host, server.InterestPort)); err != nil {
    log.Error.Println("failed to listen on interest port", err)
  } else {
    // fork and wait for handle incoming connection
    wg.Add(1)
    go func() {
      defer wg.Done()
      log.Info.Println("listening for incoming interest socket connection")
      utils.LoopWaitHandleConnection(interestLn, interestHandler)
    }()
  }

  /* 2. start data server */
  if dataLn, err := net.Listen(server.Mode, utils.JoinHostPort(server.Host, server.DataPort)); err != nil {
    log.Error.Println("failed to listen on data port", err)
  } else {
    // fork and wait for handle incoming connection
    wg.Add(1)
    go func() {
      defer wg.Done()
      log.Info.Println("listening for incoming data socket connection")
      utils.LoopWaitHandleConnection(dataLn, dataHandler)
    }()
  }

  for _, peer := range config.Agent.Peers {
    log.Info.Println("connecting to peer", peer)

    /* 3. connect to peer interest server */
    if conn, err := net.Dial(peer.Mode, utils.JoinHostPort(peer.Host, peer.InterestPort)); err != nil {
      log.Error.Printf("failed to connect to peer %v for interst (%v)\n", peer.Host, peer.InterestPort)
    } else {
      //wg.Add(1)
      //defer wg.Done()
      go interestHandler.HandleConnection(conn)
    }

    /* 4. connect to peer data server */
    if conn, err := net.Dial(peer.Mode, utils.JoinHostPort(peer.Host, peer.DataPort)); err != nil {
      log.Error.Printf("failed to connect to peer %v for interst (%v)\n", peer.Host, peer.DataPort)
    } else {
      //wg.Add(1)
      //defer wg.Done()
      go dataHandler.HandleConnection(conn)
    }
  }

  log.Info.Println("agent init finished")
  return
}
