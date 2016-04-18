package agent

import (
  "fmt"
  "net"
  "sync"
)

type dataHandler_s struct {
  wg *sync.WaitGroup
}

func (dataHandler_s) HandleError(err error) {
  fmt.Println("failed to listen on incoming data socket", err)
}

/* REMARK: this function is blocking */
func (p dataHandler_s) HandleConnection(conn net.Conn) {
  p.wg.Add(1)
  defer p.wg.Done()
  fmt.Println("client connected to data service", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
  //TODO
}
