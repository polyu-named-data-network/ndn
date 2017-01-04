package main

import (
	"github.com/polyu-named-data-network/ndn/config"
	"github.com/polyu-named-data-network/ndn/network"
	"fmt"
	mlog "github.com/beenotung/goutils/log"
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
	network.Init(config, &wg)
	fmt.Println("NDN service running")
	wg.Wait()
	fmt.Println("NDN service stopped")
}
