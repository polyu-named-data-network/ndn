package errortype

import "errors"

var (
  ContentNameNotFound     = errors.New("Content Name not found")
  ContentTypeNotSupported = errors.New("Content Type not supported")
  PortNotRegistered       = errors.New("Port not registered")
)
