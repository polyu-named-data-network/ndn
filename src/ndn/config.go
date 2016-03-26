package ndn

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Address struct {
	InterestPort int
	DataPort     int
	Host         string
	Mode         string
}
type Config struct {
	Agent struct {
		Self  Address
		Peers []Address
	}
}

func JoinHostPort(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}
func CreateConfigFromFile(filename string) (config Config, err error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("failed to open file", err)
		return
	}
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		fmt.Println("failed to decode config file", err)
	}
	return
}
