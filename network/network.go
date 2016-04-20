package network

import (
  "bitbucket.org/polyu-named-data-network/ndn/config"
  "bitbucket.org/polyu-named-data-network/ndn/packet"
  "bitbucket.org/polyu-named-data-network/ndn/portmaps"
  "bitbucket.org/polyu-named-data-network/ndn/utils"
  "encoding/json"
  "github.com/aabbcc1241/goutils/log"
  "net"
  "reflect"
  "strconv"
  "sync"
)

func handleConnection(conn net.Conn, wg *sync.WaitGroup) (err error) {
  _, port_string, err := net.SplitHostPort(conn.RemoteAddr().String())
  if err != nil {
    log.Error.Println("failed to split port from remote address", conn.RemoteAddr(), err)
    return err
  }
  port, err := strconv.Atoi(port_string)
  if err != nil {
    log.Error.Println("failed to parse port", port_string, err)
    return err
  }
  decoder := json.NewDecoder(conn)
  portmaps.AddEncoder(port, json.NewEncoder(conn))
  wg.Add(1)
  go func() {
    defer portmaps.RemoveEncoder(port)
    defer conn.Close()
    defer wg.Done()
    for err == nil {
      log.Debug.Println("waiting for incoming packet")
      var in_packet interface{}
      err = decoder.Decode(in_packet)
      if err != nil {
        log.Error.Println("failed to parse incoming message into json", err)
        break
      }
      packet_i, err2 := packet.ParsePacket(in_packet)
      if err2 != nil {
        err = err2
        log.Error.Println("failed to parse packet into packet struct", err)
        break
      }
      packetType := reflect.ValueOf(packet_i).Type()
      switch packetType {
      case reflect.ValueOf(packet.InterestPacket_s{}).Type():
        err = OnInterestPacketReceived(port, packet_i.(packet.InterestPacket_s))
        break
      case reflect.ValueOf(packet.InterestReturnPacket_s{}).Type():
        log.Error.Println("not impl", packetType)
        break
      case reflect.ValueOf(packet.DataPacket_s{}).Type():
        err = OnDataPacketReceived(packet_i.(packet.DataPacket_s))
        break
      case reflect.ValueOf(packet.ServiceProviderPacket_s{}).Type():
        OnServicePacketReceived(port, packet_i.(packet.ServiceProviderPacket_s))
        break
      default:
        log.Error.Println("not impl", packetType)
      }
    }
  }()
  return
}

/*
  1. start socket server
  2. connect to peer (socket server)
*/
func Init(config config.Config, wg *sync.WaitGroup) (err error) {
  log.Info.Println("network init start")
  server := config.Self

  /* 1. start socket server */
  var ln net.Listener
  ln, err = net.Listen(server.Mode, utils.JoinHostPort(server.Host, server.Port))
  if err == nil {
    wg.Add(1)
    go func() {
      defer wg.Done()
      var err error
      for err == nil {
        var conn net.Conn
        conn, err = ln.Accept()
        if err != nil {
          log.Error.Println(err)
        } else {
          err = handleConnection(conn, wg)
        }
      }
    }()
  } else {
    log.Error.Println("failed to bind server socket", err)
  }

  /* 2. connect to peer (socket server) */
  for _, peer := range config.Peers {
    log.Info.Println("connecting to peer", peer)
    var conn net.Conn
    if conn, err = net.Dial(peer.Mode, utils.JoinHostPort(peer.Host, peer.Port)); err != nil {
      log.Error.Printf("failed to connect to peer, host:%v, port%v", peer.Host, peer.Port)
    } else {
      handleConnection(conn, wg)
    }
  }

  log.Info.Println("network init finished")
  return
}
