package ndn

import (
  "log"
  "net"
  "os"
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

var ErrorLogger = log.New(os.Stderr, "", log.LstdFlags)
