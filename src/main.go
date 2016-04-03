package main

import (
)
import (
	"fmt"
	"ndn"
	"sync"
	"ndn/agent"
)

func main() {
	fmt.Println("NDN service initializing")
	config, err := ndn.CreateConfigFromFile("config.json")
	if err != nil {
		fmt.Println("failed to load config", err)
		return
	}
	wg := sync.WaitGroup{}
	agent.Init(config, &wg)
	fmt.Println("NDN service running")
	wg.Wait()
	fmt.Println("NDN service stopped")
}
