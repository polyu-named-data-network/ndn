package main

import (
  "crypto/rsa"
  "encoding/json"
  "fmt"
  "github.com/aabbcc1241/goutils/lang"
  "net"
  "reflect"
  "sync"
)

type PublicKey interface {
  E() string
}
type s1 struct {
  e string
  f string
}
type s2 struct {
  e string
  f string
}

//func (s s1)String() (string) {
//  return s.E()
//}
func (s s1) E() string {
  return s.e
}
func (s s2) E() string {
  return s.e
}
func compare_pointer() {
  fmt.Println("----compare by value or pointer?----")

  var p1, p2, p3 PublicKey
  p1 = s1{"1", "2"}
  p2 = s2{"1", "2"}
  p3 = s1{"1", "12"}

  fmt.Println("p1", reflect.TypeOf(p1), p1.E(), p1)
  fmt.Println("p2", reflect.TypeOf(p2), p2.E(), p2)
  fmt.Println("p3", reflect.TypeOf(p3), p3.E(), p3)
  fmt.Println("not same struct, is same?", p1 == p2)
  fmt.Println("is same value?", p1.E() == p2.E())
  fmt.Println("same struct, is same?", p1 == p3)
  fmt.Println("is same value?", p1.E() == p3.E())
  fmt.Println("self address compare", &p1 == &p1, &p1, &p1)
  fmt.Println("other address compare", &p1 == &p2, &p1, &p2)
}
func compare_string() {
  s1 := "1"
  s2 := "1"
  fmt.Println("is string same?", s1 == s2)
}
func compare_zero() {
  fmt.Println("is same", rsa.PublicKey{} == rsa.PublicKey{})
}

type InterestPacket_s struct {
  //lang.GenericStruct
  DataPort int
  Name     string
}
type DataPacket_s struct {
  //lang.GenericStruct
  Name      string
  Publisher string
}

func (s *InterestPacket_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    if err := lang.SetField(s, k, v); err != nil {
      return err
    }
  }
  return nil
}
func (s *DataPacket_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    if err := lang.SetField(s, k, v); err != nil {
      return err
    }
  }
  return nil
}
func generic_json_socket() {
  const (
    server_port = "8123"
  )

  serverLn, err := net.Listen("tcp", "127.0.0.1:"+server_port)
  if err != nil {
    panic(err)
  }
  wg := sync.WaitGroup{}
  wg.Add(1)
  /* server */
  go func() {
    defer wg.Done()
    serverConn, err := serverLn.Accept()
    if err != nil {
      panic(err)
    }
    var in_pack interface{}
    err = json.NewDecoder(serverConn).Decode(&in_pack)
    if err != nil {
      panic(err)
    }
    fmt.Println("in_pack", in_pack)
    fmt.Println("type", reflect.TypeOf(in_pack))
    p1 := InterestPacket_s{}
    p2 := DataPacket_s{}
    var b1, b2 bool
    if err := p1.FillStruct(in_pack); err != nil {
      fmt.Println("not interest pack", err)
      b1 = false
    } else {
      b1 = true
    }
    if err := p2.FillStruct(in_pack); err != nil {
      fmt.Println("not data pack", err)
      b2 = false
    } else {
      b2 = true
    }
    fmt.Println("is interest?", b1, p1)
    fmt.Println("is data?", b2, p2)
  }()
  /* client */
  wg.Add(1)
  go func() {
    defer wg.Done()
    clientConn, err := net.Dial("tcp", "127.0.0.1:"+server_port)
    if err != nil {
      panic(err)
    }
    out_pack1 := InterestPacket_s{
      DataPort: 123,
      Name:     "foo",
    }
    out_pack2 := DataPacket_s{
      Name:      "bar",
      Publisher: "peter",
    }
    err = json.NewEncoder(clientConn).Encode(out_pack1)
    fmt.Println("sent p1")
    if err != nil {
      panic(err)
    }
    err = json.NewEncoder(clientConn).Encode(out_pack2)
    fmt.Println("sent p2")
    if err != nil {
      panic(err)
    }
  }()
  wg.Wait()
}
func main() {
  generic_json_socket()
}
