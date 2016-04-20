package utils

import (
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "net"
  "strconv"
)

func JoinHostPort(host string, port int) string {
  return host + ":" + strconv.Itoa(port)
}

type ConnectionHandler interface {
  HandleConnection(net.Conn)
  HandleError(error)
}

/* REMARK: this function is blocking */
func LoopWaitHandleConnection(ln net.Listener, handler ConnectionHandler) {
  for {
    if conn, err := ln.Accept(); err != nil {
      handler.HandleError(err)
    } else {
      handler.HandleConnection(conn)
    }
  }
}

var ZeroKey = packet.PublicKey_s{}
