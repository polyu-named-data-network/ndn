package config

import (
  "encoding/json"
  "fmt"
  "os"
)

type Address struct {
  Host string
  Port int
  Mode string
}
type Config struct {
  Self  Address
  Peers []Address
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
