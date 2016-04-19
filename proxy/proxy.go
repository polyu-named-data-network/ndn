package proxy

import (
  "bitbucket.org/polyu-named-data-network/ndn/config"
  "bitbucket.org/polyu-named-data-network/ndn/fib"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/utils"
  "encoding/json"
  "github.com/aabbcc1241/goutils/log"
  "io"
  "net"
  "sync"
)

func Init(config config.Config, wg *sync.WaitGroup) (err error) {
  log.Info.Println("proxy init start")
  server := config.Proxy.ServiceProvider
  if providerLn, err := net.Listen(server.Mode, utils.JoinHostPort(server.Host, server.Port)); err != nil {
    log.Error.Println("failed to listen on Service Provider port", err)
  } else {
    log.Info.Println("listening for incoming service provider socket connection")
    wg.Add(1)
    go func() {
      defer wg.Done()
      for {
        if providerConn, err := providerLn.Accept(); err != nil {
          log.Error.Println("failed to listen on incoming provider socker", err)
        } else {
          log.Info.Println("client connected to provider service", providerConn.RemoteAddr().Network(), providerConn.RemoteAddr().String())
          //TODO
          wg.Add(1)
          go func(conn net.Conn, wg *sync.WaitGroup) {
            defer func() {
              wg.Done()
              conn.Close()
              fib.UnRegister(conn)
            }()
            decoder := json.NewDecoder(conn)
            var packet packet.ServiceProviderPacket_s
            for err == nil {
              //TODO
              err = decoder.Decode(&packet)
              if err != nil {
                if err != io.EOF {
                  log.Error.Println("failed to decode content, not service provider packet?", err)
                } else {
                  log.Info.Println("client disconnect from provider service", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
                }
              } else {
                log.Info.Println("received a servier provider packet", packet)
                fib.Register(packet.ContentName, packet.PublicKey, conn)

              }
            }
          }(providerConn, wg)
        }
      }
    }()
  }

  log.Info.Println("proxy init finished")
  return
}
