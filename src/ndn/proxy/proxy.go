package proxy

import (
  //"encoding/json"
  "fmt"
  "ndn"
  //"ndn/packet"
  "encoding/json"
  "io"
  "ndn/packet"
  "net"
  "sync"
)

func Init(config ndn.Config, wg *sync.WaitGroup) (err error) {
  fmt.Println("proxy init start")
  server := config.Proxy.ServiceProvider
  if providerLn, err := net.Listen(server.Mode, ndn.JoinHostPort(server.Host, server.Port)); err != nil {
    ndn.ErrorLogger.Println("failed to listen on Service Provider port", err)
  } else {
    fmt.Println("listening for incoming service provider socket connection")
    wg.Add(1)
    go func() {
      defer wg.Done()
      for {
        if conn, err := providerLn.Accept(); err != nil {
          ndn.ErrorLogger.Println("failed to listen on incoming provider socker", err)
        } else {
          fmt.Println("client connected to provider service", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
          //TODO
          wg.Add(1)
          go func(conn net.Conn, wg *sync.WaitGroup) {
            defer wg.Done()
            decoder := json.NewDecoder(conn)
            var packet packet.ServiceProviderPacket_s
            for err == nil {
              //TODO
              err = decoder.Decode(&packet)
              if err != nil && err != io.EOF {
                fmt.Println("failed to decode content, not service provider packet?", err)
              } else {
                fmt.Println("received a servier provider packet", packet)
              }
            }
          }(conn, wg)
        }
      }
    }()
  }

  fmt.Println("proxy init finished")
  return
}
