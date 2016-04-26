package utils

import (
  "crypto/rsa"
  "net"
  "strconv"
)

func JoinHostPort(host string, port int) string {
  return host + ":" + strconv.Itoa(port)
}
func RemotePort(conn net.Conn) int {
  _, port_string, _ := net.SplitHostPort(conn)
  port, _ := strconv.Atoi(port_string)
  return port
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

var ZeroKey = rsa.PublicKey{}
