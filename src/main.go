package main

import ()
import (
  "fmt"
  "ndn"
  "ndn/agent"
  "ndn/proxy"
  "sync"
  "time"
)

func test() {
  fmt.Println("-----------------")

  fmt.Println("now is ", time.Now())
  fmt.Println("nano", time.Now().Nanosecond())
  fmt.Println("unix nano", time.Now().UnixNano())

  fmt.Println("-----------------")
}

func main() {
  test()
  fmt.Println("NDN service initializing")
  config, err := ndn.CreateConfigFromFile("config.json")
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
