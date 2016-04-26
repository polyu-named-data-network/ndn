package errortype

import "errors"

var (
	ContentNameNotFound     = errors.New("Content Name not found")
	ContentTypeNotSupported = errors.New("Content Type not supported")
	PortNotRegistered       = errors.New("Port not registered")
	UnknownDataType         = errors.New("Unknown Data Type")
	NotImpl                 = errors.New("Not impl")
)
