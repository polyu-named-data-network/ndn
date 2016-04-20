package contentname

import (
  "github.com/aabbcc1241/goutils/lang"
)

const (
  ExactMatch int = 1 + iota
  LongestMatch
  FuzzyMatch
  Custom
)

type ContentName_s struct {
  Name           string
  ContentType    int
  ContentParam   interface{}
  AcceptDataType []string
}

func (s *ContentName_s) FillStruct(i interface{}) error {
  m := i.(map[string]interface{})
  for k, v := range m {
    if v == nil {
      //log.Debug.Println("skipped",k,v)
      continue
    }
    if err := lang.SetField(s, k, v); err != nil {
      return err
    }
    //log.Debug.Println("set success: k", k, "v", v)
  }
  return nil
}

func (original *ContentName_s) isMatch(test ContentName_s) bool {
  switch test.ContentType {
  case ExactMatch:
    return original.Name == test.Name
  default:
    return false
  }
}
