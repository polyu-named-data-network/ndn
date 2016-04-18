package utils

import (
  "fmt"
  "log"
  "net"
  "os"
  "strconv"
)

func JoinHostPort(host string, port int) string {
  return host + ":" + strconv.Itoa(port)
}

type ConnectionHandler interface {
  HandleConnection(net.Conn)
  HandleError(error)
}

/* REMARK: this function is blocking */
func LoopWaitHandleConnection(ln net.Listener, handler ConnectionHandler) {
  for {
    if conn, err := ln.Accept(); err != nil {
      handler.HandleError(err)
    } else {
      handler.HandleConnection(conn)
    }
  }
}

var errorLogger = log.New(os.Stderr, "", log.LstdFlags)
var InfoVerbose = true
var DebugVerbose = true
var ErrorVerbose = true
var WarnVerbose = true

func Info(xs ...interface{}) {
  if InfoVerbose {
    fmt.Println(xs)
  }
}
func Debug(xs ...interface{}) {
  if DebugVerbose {
    fmt.Println(xs)
  }
}
func Error(xs ...interface{}) {
  if ErrorVerbose {
    errorLogger.Println(xs)
  }
}
func Warn(xs ...interface{}) {
  if WarnVerbose {
    errorLogger.Println(xs)
  }
}
