package agent

import (
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/pit"
  "bitbucket.org/polyu-named-data-network/ndn/portmaps"
  "encoding/json"
  "fmt"
  "github.com/aabbcc1241/goutils/log"
  "io"
  "net"
  "strconv"
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
  log.Info.Println("client connected to data service", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
  _, port_string, err := net.SplitHostPort(conn.RemoteAddr().String())
  if err != nil {
    conn.Close()
    return
  }
  port, err := strconv.Atoi(port_string)
  if err != nil {
    conn.Close()
    return
  }
  p.wg.Add(1)
  go func() {
    defer p.wg.Done()
    defer conn.Close()
    in := json.NewDecoder(conn)
    portmaps.AddDataPacketEncoder(port, json.NewEncoder(conn))
    defer portmaps.RemoveDataPacketEncoder(port)
    var in_packet packet.DataPacket_s
    var err error
    for err == nil {
      err = in.Decode(&in_packet)
      if err != nil {
        if err != io.EOF {
          log.Error.Println("failed to decode incoming data packet", err)
        }
      } else {
        log.Info.Println("received data packet", in_packet)
        /*
         *    1. lookup PIT, forward to pending ports
         *    2. store in CS if allow cache
         *    3. update FIB if needed
         */

        /* 1. lookup PIT */
        err = pit.OnDataPacketReceived(in_packet)
        if err != nil {
          log.Error.Println("PIT failed to handle data packet", in_packet, err)
        }

        /* 2. store in CS if allow cache */
        if in_packet.AllowCache {
          log.Error.Println("not impl: store into content store")
        }
      }
    }
  }()
  //TODO
}
