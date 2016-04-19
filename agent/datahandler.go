package agent

import (
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/pit"
  "encoding/json"
  "fmt"
  "github.com/aabbcc1241/goutils/log"
  "io"
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
  log.Info.Println("client connected to data service", conn.RemoteAddr().Network(), conn.RemoteAddr().String())

  go func() {
    defer p.wg.Done()
    in := json.NewDecoder(conn)
    var in_packet packet.DataPacket_s
    var err error
    for err != nil {
      err = in.Decode(&in_packet)
      if err != nil {
        if err != io.EOF {
          log.Error.Println("failed to decode incoming data packet", err)
        }
      } else {
        log.Info.Println("received data packet", in_packet)
        /*
         *    1. lookup PIT
         *    2. forward to pending ports
         *    3. store in CS if allow cache
         *    4. update FIB if needed
         */

        /* 1. lookup PIT */
        var ports []int
        ports, err = pit.GetPendingPorts(in_packet.ContentName)
        if err != nil {
          log.Error.Println("failed to get pending ports", err)
        } else {
          for _,port:=range ports{
            pit.
          }
        }
      }
    }
  }()
  //TODO
}
