package main

import (
	"fmt"
	"ndn"
	"ndn/agent"
)

func main() {
	fmt.Println("program start")
	if config, err := ndn.CreateConfigFromFile("config.json"); err != nil {
		fmt.Println("failed to load config", err)
	} else {
		agent.Init(config)
	}
}
