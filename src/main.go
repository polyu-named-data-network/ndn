package main

import ()
import (
  "fmt"
  "ndn"
  "ndn/agent"
  "ndn/proxy"
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
