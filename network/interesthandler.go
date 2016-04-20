package network

import (
  "bitbucket.org/polyu-named-data-network/ndn/errortype"
  "bitbucket.org/polyu-named-data-network/ndn/fib"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/pit"
  "encoding/json"
  "github.com/aabbcc1241/goutils/log"
  "io"
  "net"
  "strconv"
  "sync"
)

type interestHandler_s struct {
  wg *sync.WaitGroup
}

func (interestHandler_s) HandleError(err error) {
  log.Error.Println("failed to listen on incoming interest socket", err)
}

/* REMARK: this function is blocking */
func (p interestHandler_s) HandleConnection(conn net.Conn) {
  p.wg.Add(1)
  log.Info.Println("client connected to interest service", conn.RemoteAddr().Network(), conn.RemoteAddr().String())

  if _, in_port_string, err := net.SplitHostPort(conn.RemoteAddr().String()); err != nil {
    log.Error.Println("failed to parse port from connection", conn)
    conn.Close()
    p.wg.Done()
  } else {
    if in_port, err := strconv.Atoi(in_port_string); err != nil {
      log.Error.Println("failed to parse port", in_port_string)
      conn.Close()
      p.wg.Done()
    } else {
      go func() {
        defer conn.Close()
        defer p.wg.Done()
        in := json.NewDecoder(conn)
        var in_packet packet.InterestPacket_s
        var err error
        for err == nil {
          log.Info.Println("wait for incoming interest packet")
          err = in.Decode(&in_packet)
          if err != nil {
            if err != io.EOF {
              log.Error.Println("failed to decode incoming interest packet", err)
            }
          } else {
            OnInterestPacketReceived(in_port, in_packet)
          }
        }
      }()
    }
  }
  //TODO
}

func OnInterestPacketReceived(port int, in_packet packet.InterestPacket_s) (err error) {
  log.Info.Println("received interest port", port, "packet", in_packet)
  /*  find data, response if found, otherwise do forwarding
   *    1. lookup CS
   *    2. lookup PIT
   *    3. lookup FIB
   *    4. calculate forwarding port according to algorithm for unknown cases
   */

  /* 1. lookup CS */
  log.Debug.Println("checking CS")
  csFound := false
  log.Debug.Println("not found in CS")

  /* 2. lookup PIT */
  log.Debug.Println("checking PIT")
  pitFound := false
  log.Debug.Println("not found in PIT")

  if csFound {

  } else if pitFound {

  } else {
    /* 3. lookup FIB */
    log.Debug.Println("checking FIB")
    port, fibFound := fib.Lookup(in_packet.ContentName, in_packet.PublisherPublicKey)
    if fibFound {
      log.Debug.Println("found in FIB, port:", port)
      if err := fib.Forward(port, in_packet); err != nil {
        log.Debug.Println("failed to forward on port", port, err)
      } else {
        pit.Register(port, in_packet)
      }
    } else {
      log.Debug.Println("not found in FIB")
      log.Error.Println("not impl")
      err = errortype.NotImpl
    }
  }
  return
}