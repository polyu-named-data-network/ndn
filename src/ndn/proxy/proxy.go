package proxy

import (
  "fmt"
  "ndn"
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
          go func(wg *sync.WaitGroup) {
            defer wg.Done()
            for {
              fmt.Println("waiting for packet")
            }
          }(&wg)
        }
      }
    }()
  }

  fmt.Println("proxy init finished")
  return
}
