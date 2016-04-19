package main

import (
  "bitbucket.org/polyu-named-data-network/ndn/agent"
  "bitbucket.org/polyu-named-data-network/ndn/config"
  "bitbucket.org/polyu-named-data-network/ndn/proxy"
  "fmt"
  mlog "github.com/aabbcc1241/goutils/log"
  "log"
  "sync"
)

func test() {
  fmt.Println("-----------------")

  //fmt.Println("now is ", time.Now())
  //fmt.Println("nano", time.Now().Nanosecond())
  //fmt.Println("unix nano", time.Now().UnixNano())

  //var xs = []int{1, 2, 3}
  //fmt.Println("xs", xs)
  //xs = append(xs, 4, 5)
  //fmt.Println("xs", xs)

  //p1 := packet.InterestPacket_s{}
  //p2 := packet.InterestPacket_s{}
  //fmt.Println("is zero value always equal?", p1 == p2)
  //fmt.Println("is zero value equal to null?", p1 == nil)

  fmt.Println("-----------------")
  //panic(0)
}

func init() {
  mlog.Init(true, true, true, log.Ltime|log.Lshortfile)
}

func main() {
  test()
  fmt.Println("NDN service initializing")
  config, err := config.CreateConfigFromFile("config.json")
  if err != nil {
    fmt.Println("failed to load config", err)
    return
  }
  wg := sync.WaitGroup{}
  agent.Init(config, &wg)
  proxy.Init(config, &wg)
  fmt.Println("NDN service running")
  wg.Wait()
  fmt.Println("NDN service stopped")
}
