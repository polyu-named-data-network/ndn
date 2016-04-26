package proxy

import (
  "bitbucket.org/polyu-named-data-network/ndn/config"
  "bitbucket.org/polyu-named-data-network/ndn/fib"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/portmaps"
  "bitbucket.org/polyu-named-data-network/ndn/utils"
  "encoding/json"
  "github.com/aabbcc1241/goutils/log"
  "io"
  "net"
  "strconv"
  "sync"
)

func Init(config config.Config, wg *sync.WaitGroup) (err error) {
  log.Info.Println("proxy init start")
  server := config.Proxy.ServiceProvider
  if providerLn, err := net.Listen(server.Mode, utils.JoinHostPort(server.Host, server.Port)); err != nil {
    log.Error.Println("failed to listen on Service Provider port", err)
  } else {
    log.Info.Printf("listening for service provider socket connection (%v)", server.Port)
    wg.Add(1)
    go func() {
      defer wg.Done()
      for {
        if providerConn, err := providerLn.Accept(); err != nil {
          log.Error.Println("failed to listen on incoming provider socker", err)
        } else {
          defer providerConn.Close()
          log.Info.Println("client connected to provider service", providerConn.RemoteAddr().Network(), providerConn.RemoteAddr().String())
          if _, port_string, err := net.SplitHostPort(providerConn.RemoteAddr().String()); err != nil {
            log.Error.Println("failed to parse port from remote address", err)
            return
          } else {
            port, err := strconv.Atoi(port_string)
            if err != nil {
              log.Error.Println("failed to parse port from string", err)
              return
            } else {
              defer fib.UnRegister(port)
              serviceProviderPacketDecoder := json.NewDecoder(providerConn)
              var serviceProviderPacket packet.ServiceProviderPacket_s
              portmaps.AddInterestPacketEncoder(port, json.NewEncoder(providerConn))
              var err error
              wg.Add(1)
              go func() {
                defer portmaps.RemoveInterestPacketEncoder(port)
                defer wg.Done()
                for err == nil {
                  err = serviceProviderPacketDecoder.Decode(&serviceProviderPacket)
                  if err != nil {
                    if err != io.EOF {
                      log.Error.Println("failed to decode, not service provider packet?", err)
                    } else {
                      log.Debug.Printf("closed service provider socket (%v)\n", port)
                    }
                  } else {
                    log.Debug.Println("received service provider packet", serviceProviderPacket)
                    fib.Register(port, serviceProviderPacket)
                  }
                }
              }()
            }
          }
        }
      }
    }()
  }

  log.Info.Println("proxy init finished")
  return
}
