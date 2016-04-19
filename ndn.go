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

func init() {
  mlog.Init(true, true, true, log.Ltime|log.Lshortfile)
}

func main() {
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
