package proxy

import (
  "encoding/json"
  "fmt"
  "io"
  "net"
  "strconv"
  "sync"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/fib"
  "bitbucket.org/polyu-named-data-network/ndn/config"
  "bitbucket.org/polyu-named-data-network/ndn/utils"
)

func Init(config config.Config, wg *sync.WaitGroup) (err error) {
  fmt.Println("proxy init start")
  server := config.Proxy.ServiceProvider
  if providerLn, err := net.Listen(server.Mode, utils.JoinHostPort(server.Host, server.Port)); err != nil {
    utils.ErrorLogger.Println("failed to listen on Service Provider port", err)
  } else {
    fmt.Println("listening for incoming service provider socket connection")
    wg.Add(1)
    go func() {
      defer wg.Done()
      for {
        if conn, err := providerLn.Accept(); err != nil {
          utils.ErrorLogger.Println("failed to listen on incoming provider socker", err)
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
              if err != nil {
                if err != io.EOF {
                  fmt.Println("failed to decode content, not service provider packet?", err)
                } else {
                  fmt.Println("client disconnect from provider service", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
                }
              } else {
                fmt.Println("received a servier provider packet", packet)
                if _, port, err := net.SplitHostPort(conn.RemoteAddr().String()); err != nil {
                  fmt.Println("failed to parse port from remote address", err)
                } else {
                  port, err := strconv.Atoi(port)
                  if err != nil {
                    fmt.Println("failed to parse port from string", err)
                  } else {
                    fib.Register(packet.ContentName, packet.PublicKey, port)
                  }
                }
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
