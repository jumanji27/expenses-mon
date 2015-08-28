package helpers

import (
    "fmt"
    "time"
    // "reflect"
)


type Main struct {}

const (
    LogTimeFormat = "02 Jan 2006 15:04:05"
)


func (self *Main) LogError(err error) {
    fmt.Printf(
        "%s | Critical Error: %s",
        time.Now().Format(LogTimeFormat),
        err.Error(),
    )
}

func (self *Main) LogSimpleMessage(message string) {
    fmt.Printf(
        "%s | %s\n",
        time.Now().Format(LogTimeFormat),
        message,
    )
}
