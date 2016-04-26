package datatype

type Base int

const (
  HTML Base = 1 + iota
  PDF
  JSON
  XML
  RAW
)
