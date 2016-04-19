package contentname

type Base int

const (
  ExactMatch Base = 1 + iota
  LongestMatch
  FuzzyMatch
  Custom
)

type ContentName_s struct {
  Name           string
  Type           Base
  Param          interface{}
  AcceptDataType []string
}

func (original *ContentName_s) isMatch(test ContentName_s) bool {
  switch test.Type {
  case ExactMatch:
    return original.Name == test.Name
  default:
    return false
  }
}
