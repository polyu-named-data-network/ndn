package agent

import (
  "encoding/json"
  "fmt"
  "io"
  "net"
  "sync"
  "bitbucket.org/polyu-named-data-network/ndn/fib"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
)

type interestHandler_s struct {
  wg *sync.WaitGroup
}

func (interestHandler_s) HandleError(err error) {
  fmt.Println("falied to listen on incoming interest socket", err)
}

/* REMARK: this function is blocking */
func (p interestHandler_s) HandleConnection(conn net.Conn) {
  p.wg.Add(1)
  fmt.Println("client connected to interest service", conn.RemoteAddr().Network(), conn.RemoteAddr().String())

  go func() {
    defer p.wg.Done()
    var err error
    decoder := json.NewDecoder(conn)
    var in_packet packet.InterestPacket_s
    for err == nil {
      fmt.Println("wait for incoming interest packet")
      err = decoder.Decode(&in_packet)
      if err != nil {
        if err != io.EOF {
          fmt.Println("failed to decode incoming interest packet", err)
        }
      } else {
        /*  find data, response if found, otherwise do forwarding
         *    1. lookup CS
         *    2. lookup PIT
         *    3. lookup FIB
         *    4. calculate forwarding port according to algorithm for unknown cases
         */
        csFound := false
        pitFound := false
        if csFound {

        } else if pitFound {

        } else {
          /* 3. lookup FIB */
          fib.Lookup(in_packet.ContentName, in_packet.PublisherPublicKey)
        }
      }
    }
  }()
  //TODO
}
