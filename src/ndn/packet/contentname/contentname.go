package contentname

type Base int

const (
	ExactMatch Base = 1 + iota
	LongestMatch
	FuzzyMatch
	Custom
)
