package helpers

import (
  "fmt"
  // "reflect"
  "time"
  "os"
)


type Main struct {}

const (
  LogTimeFormat = "02 Jan 2006 15:04:05"
)


func (self *Main) LogWarning(err error) {
  fmt.Printf(
    "%s | Warning: %s\n",
    time.Now().Format(LogTimeFormat),
    err.Error(),
  )
}

func (self *Main) LogError(err error) {
  fmt.Printf(
    "%s | Error: %s\n",
    time.Now().Format(LogTimeFormat),
    err.Error(),
  )
  os.Exit(1)
}

func (self *Main) LogSimpleMessage(message string) {
  fmt.Printf(
    "%s | %s\n",
    time.Now().Format(LogTimeFormat),
    message,
  )
}
