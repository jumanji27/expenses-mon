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


func (self *Main) CreateEvent(title string, message string) {
  if title == "Log" {
    fmt.Printf(
      "%s | %s\n",
      time.Now().Format(LogTimeFormat),
      message,
    )
  } else {
    fmt.Printf(
      "%s | %s: %s\n",
      time.Now().Format(LogTimeFormat),
      title,
      message,
    )
  }

  if title == "Error" {
    os.Exit(1)
  }
}
