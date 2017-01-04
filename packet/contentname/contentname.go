package contentname

import "github.com/polyu-named-data-network/ndn/packet/datatype"

type Base int

const (
	ExactMatch Base = 1 + iota
	LongestMatch
	FuzzyMatch
	Custom
)

type ContentName_s struct {
	Name           string
	ContentType    Base
	ContentParam   interface{}
	AcceptDataType []datatype.Base
}

func (original *ContentName_s) isMatch(test ContentName_s) bool {
	switch test.ContentType {
	case ExactMatch:
		return original.Name == test.Name
	default:
		return false
	}
}
