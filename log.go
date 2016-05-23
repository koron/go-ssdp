package ssdp

import (
	"io/ioutil"
	"log"
)

// Logger is default logger for SSDP module.
var Logger = log.New(ioutil.Discard, "", log.LstdFlags|log.Lmicroseconds)

func logf(s string, a ...interface{}) {
	Logger.Printf(s, a...)
}
