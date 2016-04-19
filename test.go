package main

import (
  "crypto/rsa"
  "fmt"
  "reflect"
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
func main() {
}
